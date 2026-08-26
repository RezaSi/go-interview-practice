package handlers

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"web-ui/internal/services"
	"web-ui/internal/utils"
)

// ReleaseHandler serves the "What's new in Go" track: release landing pages,
// per-feature explainers with their visualizations, and the hands-on challenges.
type ReleaseHandler struct {
	content        embed.FS
	releaseService *services.ReleaseService
}

func NewReleaseHandler(content embed.FS, releaseService *services.ReleaseService) *ReleaseHandler {
	return &ReleaseHandler{content: content, releaseService: releaseService}
}

// Route dispatches everything under /releases.
//
//	/releases                                  → all releases
//	/releases/1.27                             → one release, its features
//	/releases/1.27/generic-methods             → explainer + visualization
//	/releases/1.27/generic-methods/challenge-1 → hands-on editor
func (h *ReleaseHandler) Route(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/releases"), "/")
	if path == "" {
		h.indexPage(w, r)
		return
	}

	parts := strings.Split(path, "/")
	parts[0] = strings.TrimPrefix(parts[0], "go") // tolerate /releases/go1.27

	switch len(parts) {
	case 1:
		h.releasePage(w, r, parts[0])
	case 2:
		h.featurePage(w, r, parts[0], parts[1])
	case 3:
		h.challengePage(w, r, parts[0], parts[1], parts[2])
	default:
		http.NotFound(w, r)
	}
}

func (h *ReleaseHandler) indexPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, "templates/releases.html", map[string]interface{}{
		"Releases": h.releaseService.GetReleases(),
	})
}

func (h *ReleaseHandler) releasePage(w http.ResponseWriter, r *http.Request, version string) {
	rel := h.releaseService.GetRelease(version)
	if rel == nil {
		http.NotFound(w, r)
		return
	}
	h.render(w, "templates/release_detail.html", map[string]interface{}{
		"Release": rel,
	})
}

func (h *ReleaseHandler) featurePage(w http.ResponseWriter, r *http.Request, version, slug string) {
	rel := h.releaseService.GetRelease(version)
	feature := h.releaseService.GetFeature(version, slug)
	if rel == nil || feature == nil || feature.Status != "available" {
		http.NotFound(w, r)
		return
	}
	h.render(w, "templates/release_feature.html", map[string]interface{}{
		"Release": rel,
		"Feature": feature,
	})
}

func (h *ReleaseHandler) challengePage(w http.ResponseWriter, r *http.Request, version, slug, cslug string) {
	rel := h.releaseService.GetRelease(version)
	feature := h.releaseService.GetFeature(version, slug)
	challenge := h.releaseService.GetChallenge(version, slug, cslug)
	if rel == nil || feature == nil || challenge == nil {
		http.NotFound(w, r)
		return
	}

	var prev, next string
	for i, c := range feature.Challenges {
		if c.Slug != cslug {
			continue
		}
		if i > 0 {
			prev = feature.Challenges[i-1].Slug
		}
		if i < len(feature.Challenges)-1 {
			next = feature.Challenges[i+1].Slug
		}
	}

	h.render(w, "templates/release_challenge.html", map[string]interface{}{
		"Release":       rel,
		"Feature":       feature,
		"Challenge":     challenge,
		"PrevChallenge": prev,
		"NextChallenge": next,
		"RunnerEnabled": h.releaseService.RunnerEnabled(),
	})
}

func (h *ReleaseHandler) render(w http.ResponseWriter, page string, data map[string]interface{}) {
	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).
		ParseFS(h.content, "templates/base.html", "templates/release_common.html", page)
	if err != nil {
		log.Printf("releases: template parse error (%s): %v", page, err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("releases: template execute error (%s): %v", page, err)
	}
}

// ── API ──────────────────────────────────────────────────────────────────────

type releaseRunRequest struct {
	Release   string `json:"release"`
	Feature   string `json:"feature"`
	Challenge string `json:"challenge"`
	Code      string `json:"code"`
}

// RunChallenge compiles and tests a submitted solution for a release challenge.
func (h *ReleaseHandler) RunChallenge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req releaseRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	challenge := h.releaseService.GetChallenge(
		strings.TrimPrefix(req.Release, "go"), req.Feature, req.Challenge)
	if challenge == nil {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	result := h.releaseService.RunChallenge(req.Code, challenge)
	json.NewEncoder(w).Encode(result)
}
