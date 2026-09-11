package services

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"web-ui/internal/models"
)

// ReleaseService loads the releases/ track from disk and runs its challenges.
//
// Everything is directory-driven: dropping a new releases/<version>/<feature>/
// directory in publishes it. There is no registry to update.
type ReleaseService struct {
	releasesPath string
	cached       []*models.Release
}

func NewReleaseService() *ReleaseService {
	return &ReleaseService{releasesPath: "../releases"} // relative to web-ui/
}

// Load reads every release from disk. Called once at start-up.
func (s *ReleaseService) Load() error {
	s.cached = nil
	releases := s.GetReleases()

	features, challenges := 0, 0
	for _, r := range releases {
		features += r.AvailableCount()
		challenges += r.ChallengeCount()
	}
	log.Printf("Loaded %d release(s), %d documented feature(s), %d challenge(s)",
		len(releases), features, challenges)
	return nil
}

// GetReleases returns all releases, newest version first.
func (s *ReleaseService) GetReleases() []*models.Release {
	if s.cached != nil {
		return s.cached
	}

	entries, err := os.ReadDir(s.releasesPath)
	if err != nil {
		log.Printf("releases: cannot read %s: %v", s.releasesPath, err)
		s.cached = []*models.Release{}
		return s.cached
	}

	var out []*models.Release
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if r := s.loadRelease(filepath.Join(s.releasesPath, e.Name())); r != nil {
			out = append(out, r)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return versionLess(out[j].Version, out[i].Version) // descending
	})

	s.cached = out
	return s.cached
}

func (s *ReleaseService) GetRelease(version string) *models.Release {
	for _, r := range s.GetReleases() {
		if r.Version == version || "go"+r.Version == version {
			return r
		}
	}
	return nil
}

func (s *ReleaseService) GetFeature(version, slug string) *models.ReleaseFeature {
	r := s.GetRelease(version)
	if r == nil {
		return nil
	}
	for _, f := range r.Features {
		if f.Slug == slug {
			return f
		}
	}
	return nil
}

func (s *ReleaseService) GetChallenge(version, feature, challenge string) *models.ReleaseChallenge {
	f := s.GetFeature(version, feature)
	if f == nil {
		return nil
	}
	for _, c := range f.Challenges {
		if c.Slug == challenge {
			return c
		}
	}
	return nil
}

// ── loading ──────────────────────────────────────────────────────────────────

func (s *ReleaseService) loadRelease(dir string) *models.Release {
	raw, err := os.ReadFile(filepath.Join(dir, "release.json"))
	if err != nil {
		return nil // not a release directory (e.g. the track README)
	}

	var rel models.Release
	if err := json.Unmarshal(raw, &rel); err != nil {
		log.Printf("releases: bad release.json in %s: %v", dir, err)
		return nil
	}

	for i, slug := range rel.FeatureSlugs {
		fdir := filepath.Join(dir, slug)
		f := s.loadFeature(fdir, slug, &rel)
		if f == nil {
			// Listed in release.json but not written yet.
			f = &models.ReleaseFeature{
				Slug:             slug,
				Title:            titleFromSlug(slug),
				Kind:             "language",
				ShortDescription: "Explainer, visualization and challenges are on the way.",
				Status:           "coming-soon",
				Icon:             "bi-hourglass-split",
				ReleaseVersion:   rel.Version,
				ReleaseName:      rel.DisplayName,
				Toolchain:        rel.Toolchain,
			}
		}
		if f.Order == 0 {
			f.Order = i + 1
		}
		rel.Features = append(rel.Features, f)
	}

	sort.SliceStable(rel.Features, func(i, j int) bool { return rel.Features[i].Order < rel.Features[j].Order })
	return &rel
}

func (s *ReleaseService) loadFeature(dir, slug string, rel *models.Release) *models.ReleaseFeature {
	raw, err := os.ReadFile(filepath.Join(dir, "feature.json"))
	if err != nil {
		return nil
	}

	var f models.ReleaseFeature
	if err := json.Unmarshal(raw, &f); err != nil {
		log.Printf("releases: bad feature.json in %s: %v", dir, err)
		return nil
	}

	f.Slug = slug
	f.Status = "available"
	f.ReleaseVersion = rel.Version
	f.ReleaseName = rel.DisplayName
	f.Toolchain = rel.Toolchain
	if f.Kind == "" {
		f.Kind = "language"
	}
	if f.Icon == "" {
		f.Icon = "bi-stars"
	}

	f.SummaryHTML = RenderInline(f.Summary)
	f.ShortDescriptionHTML = RenderInline(f.ShortDescription)

	if md, err := os.ReadFile(filepath.Join(dir, "explainer.md")); err == nil {
		f.ExplainerHTML, f.Sections = RenderMarkdown(string(md))
	}

	// proposal.md is the optional deep dive: design history and internals.
	if md := readFile(filepath.Join(dir, "proposal.md")); md != "" {
		f.ProposalHTML, _ = RenderMarkdown(md)
		f.HasProposal = true
	}

	// diagram.svg is a hand-authored schematic shown above the explainer.
	if d, err := os.ReadFile(filepath.Join(dir, "diagram.svg")); err == nil && len(strings.TrimSpace(string(d))) > 0 {
		f.DiagramSVG = template.HTML(string(d))
		f.HasDiagram = true
	}

	// visual.html is a trusted, self-contained fragment authored in this repo.
	if v, err := os.ReadFile(filepath.Join(dir, "visual.html")); err == nil && len(strings.TrimSpace(string(v))) > 0 {
		f.VisualHTML = template.HTML(string(v))
		f.HasVisual = true
	}

	for i, cslug := range f.ChallengeSlugs {
		c := s.loadChallenge(filepath.Join(dir, cslug), cslug, &f, rel)
		if c == nil {
			continue
		}
		if c.Order == 0 {
			c.Order = i + 1
		}
		f.Challenges = append(f.Challenges, c)
	}
	sort.SliceStable(f.Challenges, func(i, j int) bool { return f.Challenges[i].Order < f.Challenges[j].Order })

	return &f
}

func (s *ReleaseService) loadChallenge(dir, slug string, f *models.ReleaseFeature, rel *models.Release) *models.ReleaseChallenge {
	raw, err := os.ReadFile(filepath.Join(dir, "metadata.json"))
	if err != nil {
		return nil
	}

	var c models.ReleaseChallenge
	if err := json.Unmarshal(raw, &c); err != nil {
		log.Printf("releases: bad metadata.json in %s: %v", dir, err)
		return nil
	}

	c.Slug = slug
	c.FeatureSlug = f.Slug
	c.FeatureTitle = f.Title
	c.ReleaseVersion = rel.Version
	if c.GoVersion == "" {
		c.GoVersion = strings.TrimPrefix(rel.Toolchain, "go")
	}
	if c.Icon == "" {
		c.Icon = "bi-code-square"
	}

	c.Template = readFile(filepath.Join(dir, "solution-template.go"))
	c.TestFile = readFile(filepath.Join(dir, "solution-template_test.go"))
	c.GoMod = readFile(filepath.Join(dir, "go.mod"))

	if md := readFile(filepath.Join(dir, "README.md")); md != "" {
		c.ReadmeHTML, _ = RenderMarkdown(md)
	}
	if md := readFile(filepath.Join(dir, "learning.md")); md != "" {
		c.LearningHTML, _ = RenderMarkdown(md)
	}
	if md := readFile(filepath.Join(dir, "hints.md")); md != "" {
		c.HintsHTML, _ = RenderMarkdown(md)
		c.Hints = splitHints(md)
	}

	return &c
}

// ── running ──────────────────────────────────────────────────────────────────

// RunnerEnabled reports whether in-browser test running is switched on.
//
// Like the local ExecutionService this compiles and runs submitted code in a
// temporary directory on the host. That is fine for local development. Anywhere
// untrusted users can reach it, set RELEASES_RUNNER=off and serve the content
// read-only, or route this through the sandboxed execution engine.
func (s *ReleaseService) RunnerEnabled() bool {
	return !strings.EqualFold(os.Getenv("RELEASES_RUNNER"), "off")
}

// RunChallenge compiles the submitted code together with the challenge's test
// file and reports the result.
func (s *ReleaseService) RunChallenge(code string, c *models.ReleaseChallenge) models.ReleaseRunResult {
	start := time.Now()
	toolchain := "go" + c.GoVersion

	if !s.RunnerEnabled() {
		return models.ReleaseRunResult{
			Output:    "In-browser test running is disabled on this instance (RELEASES_RUNNER=off).\nClone the repo and run: go test -v ./...",
			Toolchain: toolchain,
		}
	}

	tmp, err := os.MkdirTemp("", "release-exec-")
	if err != nil {
		return models.ReleaseRunResult{Output: fmt.Sprintf("Failed to create temp dir: %v", err), Toolchain: toolchain}
	}
	defer os.RemoveAll(tmp)

	// Use the challenge's own go.mod so an in-browser run and a local
	// `go test` are compiling against exactly the same toolchain directive.
	gomod := c.GoMod
	if strings.TrimSpace(gomod) == "" {
		gomod = fmt.Sprintf("module %s\n\ngo %s\n", moduleName(c.Slug), c.GoVersion)
	}

	files := map[string]string{
		"go.mod":                    gomod,
		"solution-template.go":      code,
		"solution-template_test.go": c.TestFile,
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte(body), 0o644); err != nil {
			return models.ReleaseRunResult{Output: fmt.Sprintf("Failed to write %s: %v", name, err), Toolchain: toolchain}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "test", "-v", "./...")
	cmd.Dir = tmp
	// The go.mod may name a toolchain newer than the one installed (a release
	// candidate, for instance). GOTOOLCHAIN=auto lets the go command fetch it.
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=auto", "GOFLAGS=")

	out, runErr := cmd.CombinedOutput()
	res := models.ReleaseRunResult{
		Output:      string(out),
		ExecutionMs: time.Since(start).Milliseconds(),
		Toolchain:   toolchain,
	}

	if ctx.Err() == context.DeadlineExceeded {
		res.Output += "\n\nTimed out after 3 minutes."
		return res
	}
	if runErr == nil {
		res.Passed = true
		return res
	}
	if _, isExit := runErr.(*exec.ExitError); !isExit {
		res.Output = fmt.Sprintf("Failed to run tests: %v\n%s", runErr, res.Output)
	}
	return res
}

// ── helpers ──────────────────────────────────────────────────────────────────

func readFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func moduleName(slug string) string {
	var b strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "challenge"
	}
	return b.String()
}

func titleFromSlug(slug string) string {
	parts := strings.Split(slug, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

// versionLess compares dotted numeric versions such as "1.27" and "1.9".
func versionLess(a, b string) bool {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(as) || i < len(bs); i++ {
		var ai, bi int
		if i < len(as) {
			ai, _ = strconv.Atoi(as[i])
		}
		if i < len(bs) {
			bi, _ = strconv.Atoi(bs[i])
		}
		if ai != bi {
			return ai < bi
		}
	}
	return false
}

// reHintHeading matches the "## Hint 3: Something" headings that separate one hint
// from the next, matching the convention used by the packages/ challenges.
var reHintHeading = regexp.MustCompile(`(?im)^##\s+Hint\s+(\d+)\s*[:.-]?\s*(.*)$`)

// splitHints breaks a hints.md file into its individual hints so the page can
// reveal them one at a time.
func splitHints(md string) []models.Hint {
	locs := reHintHeading.FindAllStringSubmatchIndex(md, -1)
	if len(locs) == 0 {
		return nil
	}

	var out []models.Hint
	for i, loc := range locs {
		title := strings.TrimSpace(md[loc[4]:loc[5]])
		start := loc[1]
		end := len(md)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		body := strings.TrimSpace(md[start:end])
		html, _ := RenderMarkdown(body)
		out = append(out, models.Hint{Title: title, HTML: html})
	}
	return out
}
