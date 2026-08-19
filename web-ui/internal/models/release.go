package models

import "html/template"

// Release is one Go release (e.g. Go 1.27) and the features we teach from it.
// It is loaded from releases/<version>/release.json.
type Release struct {
	Version      string   `json:"version"`
	DisplayName  string   `json:"display_name"`
	Status       string   `json:"status"`    // upcoming | rc | released
	Toolchain    string   `json:"toolchain"` // e.g. go1.27rc2 — written into each challenge's go.mod
	Expected     string   `json:"expected"`
	Headline     string   `json:"headline"`
	Summary      string   `json:"summary"`
	DocsURL      string   `json:"docs_url"`
	FeatureSlugs []string `json:"features"`

	// Populated from disk.
	Features []*ReleaseFeature `json:"feature_details,omitempty"`
}

// AvailableCount reports how many of this release's features have content on disk.
func (r *Release) AvailableCount() int {
	n := 0
	for _, f := range r.Features {
		if f.Status == "available" {
			n++
		}
	}
	return n
}

// ChallengeCount reports the total number of hands-on challenges in the release.
func (r *Release) ChallengeCount() int {
	n := 0
	for _, f := range r.Features {
		n += len(f.Challenges)
	}
	return n
}

// ReleaseFeature is a single feature of a release: an explainer, an optional
// interactive visualization, and a list of challenges.
type ReleaseFeature struct {
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Kind             string   `json:"kind"` // language | stdlib | toolchain | runtime
	ShortDescription string   `json:"short_description"`
	Summary          string   `json:"summary"`
	Difficulty       string   `json:"difficulty"`
	EstimatedTime    string   `json:"estimated_time"`
	Icon             string   `json:"icon"`
	Tags             []string `json:"tags"`
	ProposalURL      string   `json:"proposal_url"`
	DiagramAlt       string   `json:"diagram_caption"`
	DocsURL          string   `json:"docs_url"`
	Order            int      `json:"order"`
	ChallengeSlugs   []string `json:"challenges"`

	// Populated from disk.
	Status               string              `json:"status"` // available | coming-soon
	SummaryHTML          template.HTML       `json:"-"`
	ShortDescriptionHTML template.HTML       `json:"-"`
	ReleaseVersion       string              `json:"release_version"`
	ReleaseName          string              `json:"release_name"`
	Toolchain            string              `json:"toolchain"`
	ExplainerHTML        template.HTML       `json:"-"`
	Sections             []FeatureSection    `json:"sections,omitempty"`
	VisualHTML           template.HTML       `json:"-"`
	HasVisual            bool                `json:"has_visual"`
	ProposalHTML         template.HTML       `json:"-"`
	HasProposal          bool                `json:"has_proposal"`
	DiagramSVG           template.HTML       `json:"-"`
	HasDiagram           bool                `json:"has_diagram"`
	DiagramCaption       string              `json:"diagram_caption"`
	Challenges           []*ReleaseChallenge `json:"challenge_details,omitempty"`
}

// FeatureSection is one `##` heading of the explainer, used to build its
// table of contents.
type FeatureSection struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// ReleaseChallenge is one hands-on exercise attached to a feature.
type ReleaseChallenge struct {
	Slug                string   `json:"slug"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	ShortDescription    string   `json:"short_description"`
	Difficulty          string   `json:"difficulty"`
	EstimatedTime       string   `json:"estimated_time"`
	GoVersion           string   `json:"go_version"`
	LearningObjectives  []string `json:"learning_objectives"`
	Prerequisites       []string `json:"prerequisites"`
	Tags                []string `json:"tags"`
	Requirements        []string `json:"requirements"`
	BonusPoints         []string `json:"bonus_points"`
	RealWorldConnection string   `json:"real_world_connection"`
	Icon                string   `json:"icon"`
	Order               int      `json:"order"`

	// Populated from disk.
	FeatureSlug    string        `json:"feature_slug"`
	FeatureTitle   string        `json:"feature_title"`
	ReleaseVersion string        `json:"release_version"`
	Template       string        `json:"template"`
	TestFile       string        `json:"test_file"`
	GoMod          string        `json:"-"`
	ReadmeHTML     template.HTML `json:"-"`
	LearningHTML   template.HTML `json:"-"`
	HintsHTML      template.HTML `json:"-"`
	Hints          []Hint        `json:"-"`
}

// ReleaseRunResult is the outcome of running a release challenge's tests.
type ReleaseRunResult struct {
	Passed      bool   `json:"passed"`
	Output      string `json:"output"`
	ExecutionMs int64  `json:"executionMs"`
	Toolchain   string `json:"toolchain"`
}

// Hint is one step of a challenge's progressive hints. The site reveals them one
// at a time, so hints.md is split on its "## Hint N: title" headings, the same
// convention the packages/ challenges use.
type Hint struct {
	Title string
	HTML  template.HTML
}
