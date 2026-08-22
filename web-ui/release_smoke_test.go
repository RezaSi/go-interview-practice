package main

import (
	"html/template"
	"strings"
	"testing"

	"web-ui/internal/services"
	"web-ui/internal/utils"
)

// Local smoke test for the releases track: everything on disk loads coherently and
// every release template parses against base.html with the shared func map.
func TestReleasesLoadAndTemplatesParse(t *testing.T) {
	svc := services.NewReleaseService()
	releases := svc.GetReleases()
	if len(releases) == 0 {
		t.Fatal("no releases loaded from ../releases")
	}

	for _, r := range releases {
		t.Logf("release %s: %d documented feature(s), %d challenge(s)",
			r.Version, r.AvailableCount(), r.ChallengeCount())

		if r.Status == "" || r.Toolchain == "" || r.Headline == "" {
			t.Errorf("release %s is missing status, toolchain or headline", r.Version)
		}

		for _, f := range r.Features {
			if f.Status != "available" {
				continue
			}
			if len(f.Sections) < 3 {
				t.Errorf("%s/%s: explainer produced %d sections, expected a real table of contents",
					r.Version, f.Slug, len(f.Sections))
			}
			if !strings.Contains(string(f.ExplainerHTML), "<h2") {
				t.Errorf("%s/%s: explainer HTML has no headings", r.Version, f.Slug)
			}
			if !f.HasDiagram {
				t.Errorf("%s/%s: no diagram.svg", r.Version, f.Slug)
			}
			if len(f.Challenges) == 0 {
				t.Errorf("%s/%s: no challenges", r.Version, f.Slug)
			}
			for _, c := range f.Challenges {
				if c.Template == "" || c.TestFile == "" || c.GoMod == "" {
					t.Errorf("%s/%s/%s: missing template, test file or go.mod", r.Version, f.Slug, c.Slug)
				}
				if !strings.Contains(c.GoMod, r.Version) {
					t.Errorf("%s/%s/%s: go.mod does not pin %s: %q",
						r.Version, f.Slug, c.Slug, r.Version, c.GoMod)
				}
				if c.ReadmeHTML == "" || c.LearningHTML == "" || c.HintsHTML == "" {
					t.Errorf("%s/%s/%s: missing README, learning or hints", r.Version, f.Slug, c.Slug)
				}
			}
		}
	}

	for _, page := range []string{"releases", "release_detail", "release_feature", "release_challenge"} {
		if _, err := template.New("").Funcs(utils.GetTemplateFuncs()).
			ParseFiles("templates/base.html", "templates/release_common.html", "templates/"+page+".html"); err != nil {
			t.Errorf("template %s failed to parse: %v", page, err)
		}
	}
}
