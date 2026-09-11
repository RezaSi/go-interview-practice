# Go Releases Track

Every Go release ships features that most developers read about once in the release
notes and then never actually *use*. This track exists to close that gap. For each
notable feature we publish three things, together, on release day:

1. **An explainer**: written by a human, not a bullet-point summary of the release
   notes. What changed, why it changed, what it lets you do that you couldn't do
   before, and where the sharp edges are.
2. **A visualization**: an interactive diagram of the mechanism. Compilers and type
   systems are much easier to understand when you can see the pieces move.
3. **Hands-on challenges**: real code, real tests, run in the browser against the
   actual toolchain that ships the feature.

## Directory layout

```
releases/
  <release>/                       e.g. go1.27
    release.json                   release metadata + the feature list
    <feature>/                     e.g. generic-methods
      feature.json                 feature metadata
      explainer.md                 the written explainer
      visual.html                  the interactive visualization (optional)
      <challenge>/                 e.g. challenge-1-your-first-generic-method
        metadata.json
        README.md
        learning.md
        hints.md
        go.mod
        solution-template.go
        solution-template_test.go
        run_tests.sh
```

The web UI discovers everything on disk. There is no registry to update and no code
to write to publish a new feature, drop the directory in and it appears.

## Adding a new release

1. `mkdir releases/go1.28` and write `release.json`:

```json
{
  "version": "1.28",
  "display_name": "Go 1.28",
  "status": "upcoming",
  "toolchain": "go1.28rc1",
  "expected": "February 2027",
  "headline": "One sentence a developer would actually repeat to a colleague.",
  "summary": "A short paragraph for the release landing page.",
  "docs_url": "https://tip.golang.org/doc/go1.28",
  "features": ["some-feature", "another-feature"]
}
```

`status` is one of `upcoming`, `rc`, `released`.

`toolchain` is the Go toolchain the challenges are compiled with. It is written
verbatim into each challenge's `go.mod`, and Go's `GOTOOLCHAIN=auto` (the default)
downloads it on demand, so a release-candidate feature can be practised before the
final release, on a machine whose system Go is a version behind.

2. List feature slugs in `features`. A slug with no directory on disk renders as
   **Coming soon**: which is how you publish the release landing page on day one and
   fill in the features over the following days.

## Adding a feature

`mkdir releases/go1.28/some-feature` and write `feature.json`:

```json
{
  "title": "Some Feature",
  "kind": "language",
  "short_description": "One line for the feature card.",
  "summary": "A paragraph for the top of the feature page.",
  "difficulty": "Intermediate",
  "estimated_time": "35 min",
  "icon": "bi-braces",
  "tags": ["generics", "types"],
  "proposal_url": "https://github.com/golang/go/issues/12345",
  "docs_url": "https://tip.golang.org/doc/go1.28#some-feature",
  "order": 1,
  "challenges": ["challenge-1-something"]
}
```

`kind` is one of `language`, `stdlib`, `toolchain`, `runtime`, it only sets the
badge colour and grouping.

Then write `explainer.md`, optionally `visual.html`, and the challenges.

### explainer.md

Plain Markdown with `##` section headings. The page builds its own table of contents
from those headings, so keep them short and meaningful. Fenced code blocks are
syntax-highlighted; use ```go.

Two conventions the renderer understands:

- A blockquote starting with `> **Try it:**` is styled as a call-to-action box.
- A fenced block tagged ```text is rendered as terminal output (used for compiler
  errors), not as Go source.

### visual.html

A self-contained HTML fragment, no `<html>`/`<body>` wrapper, no external requests.
Scope every selector and every CSS custom property to a unique root class so two
visualizations on different features can never collide. Bootstrap 5.3 and
bootstrap-icons are already loaded by the page.

### Challenges

Same shape as the `packages/` challenges, with one addition: `metadata.json` may
carry a `"go_version"` that overrides the release toolchain for that challenge.

`solution-template.go` must **compile** as shipped, leave TODOs that return zero
values rather than leaving the file syntactically incomplete. The learner's first
run should be a clean, readable test failure, not a wall of parser errors.

`solution-template_test.go` is the graded test file. Write assertions that fail
loudly and explain themselves; the failure message is the teaching moment.

## Local development

```bash
cd web-ui && go run main.go
```

Then open <http://localhost:8080/releases>.
