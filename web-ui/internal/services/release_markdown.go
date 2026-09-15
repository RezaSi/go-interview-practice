package services

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"web-ui/internal/models"
)

// A small, dependency-free Markdown renderer for the releases track.
//
// The repo already has a regex-based converter in internal/utils, but it predates
// this track and does not handle tables, blockquotes, or code blocks that contain
// asterisks, all of which the release explainers rely on heavily. Rather than
// change behaviour for every existing page, the releases track renders its own
// Markdown here.
//
// Release content is authored in this repository, so inline HTML in it is trusted
// and passes through. Text inside code spans and code fences is always escaped.

var (
	reHeading   = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	reOrdered   = regexp.MustCompile(`^\d+\.\s+(.*)$`)
	reTableSep  = regexp.MustCompile(`^\|?\s*:?-{2,}:?\s*(\|\s*:?-{2,}:?\s*)*\|?$`)
	reLink      = regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)\)`)
	reBold      = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	reItalic    = regexp.MustCompile(`(^|[^*\w])\*([^*\n]+)\*([^*\w]|$)`)
	reNonSlug   = regexp.MustCompile(`[^a-z0-9]+`)
	reCodeSpan  = regexp.MustCompile("`([^`]+)`")
	reTryItLead = regexp.MustCompile(`^\s*\*\*Try it:\*\*\s*`)
)

// RenderMarkdown converts release Markdown to HTML and returns the `##` headings
// it found, in order, for building a table of contents.
func RenderMarkdown(src string) (template.HTML, []models.FeatureSection) {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	lines := strings.Split(src, "\n")

	var out strings.Builder
	var sections []models.FeatureSection
	seen := map[string]int{}

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		switch {
		case trimmed == "":
			i++

		// ── fenced code ───────────────────────────────────────────────────────
		case strings.HasPrefix(trimmed, "```"):
			lang := strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			i++
			var body []string
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
				body = append(body, lines[i])
				i++
			}
			i++ // closing fence
			code := template.HTMLEscapeString(strings.Join(body, "\n"))
			switch lang {
			case "text", "console", "output":
				// Compiler and test output, not source. Rendered as a terminal.
				out.WriteString(`<pre class="rl-term"><code>` + code + `</code></pre>` + "\n")
			case "":
				out.WriteString(`<pre><code>` + code + `</code></pre>` + "\n")
			default:
				out.WriteString(`<pre><code class="language-` + template.HTMLEscapeString(lang) + `">` + code + "</code></pre>\n")
			}

		// ── raw HTML block (e.g. <details> in hints) ──────────────────────────
		case strings.HasPrefix(trimmed, "<"):
			out.WriteString(line + "\n")
			i++

		// ── heading ───────────────────────────────────────────────────────────
		case reHeading.MatchString(trimmed):
			m := reHeading.FindStringSubmatch(trimmed)
			level := len(m[1])
			text := m[2]
			id := slugify(text)
			if n, dup := seen[id]; dup {
				seen[id] = n + 1
				id = fmt.Sprintf("%s-%d", id, n+1)
			} else {
				seen[id] = 1
			}
			if level == 2 {
				sections = append(sections, models.FeatureSection{ID: id, Title: stripInline(text)})
			}
			out.WriteString(fmt.Sprintf(`<h%d id="%s" class="rl-h%d">%s</h%d>`+"\n",
				level, id, level, inline(text), level))
			i++

		// ── horizontal rule ───────────────────────────────────────────────────
		case trimmed == "---" || trimmed == "***":
			out.WriteString("<hr>\n")
			i++

		// ── table ─────────────────────────────────────────────────────────────
		case strings.HasPrefix(trimmed, "|") && i+1 < len(lines) && reTableSep.MatchString(strings.TrimSpace(lines[i+1])):
			header := splitRow(trimmed)
			i += 2
			var rows [][]string
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "|") {
				rows = append(rows, splitRow(strings.TrimSpace(lines[i])))
				i++
			}
			out.WriteString(`<div class="table-responsive"><table class="table rl-table">` + "\n<thead><tr>")
			for _, h := range header {
				out.WriteString("<th>" + inline(h) + "</th>")
			}
			out.WriteString("</tr></thead>\n<tbody>")
			for _, r := range rows {
				out.WriteString("<tr>")
				for _, c := range r {
					out.WriteString("<td>" + inline(c) + "</td>")
				}
				out.WriteString("</tr>")
			}
			out.WriteString("</tbody></table></div>\n")

		// ── blockquote ────────────────────────────────────────────────────────
		case strings.HasPrefix(trimmed, ">"):
			var body []string
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), ">") {
				body = append(body, strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[i]), ">")), " "))
				i++
			}
			text := strings.Join(body, " ")
			if reTryItLead.MatchString(text) {
				text = reTryItLead.ReplaceAllString(text, "")
				out.WriteString(`<div class="rl-tryit"><span class="rl-tryit-tag"><i class="bi bi-play-fill"></i> Try it</span>` +
					`<div>` + inline(text) +
					`<div class="rl-tryit-go"><a href="#practice" class="rl-tryit-btn">Start the first challenge <i class="bi bi-arrow-right"></i></a></div>` +
					`</div></div>` + "\n")
			} else {
				out.WriteString(`<blockquote class="rl-quote">` + inline(text) + `</blockquote>` + "\n")
			}

		// ── unordered list ────────────────────────────────────────────────────
		case strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* "):
			out.WriteString("<ul class=\"rl-list\">\n")
			for i < len(lines) {
				t := strings.TrimSpace(lines[i])
				if !strings.HasPrefix(t, "- ") && !strings.HasPrefix(t, "* ") {
					break
				}
				item := t[2:]
				i++
				// fold continuation lines into the same bullet
				for i < len(lines) {
					n := strings.TrimSpace(lines[i])
					if n == "" || strings.HasPrefix(n, "- ") || strings.HasPrefix(n, "* ") ||
						strings.HasPrefix(n, "#") || strings.HasPrefix(n, "```") || strings.HasPrefix(n, "|") {
						break
					}
					item += " " + n
					i++
				}
				out.WriteString("<li>" + inline(item) + "</li>\n")
			}
			out.WriteString("</ul>\n")

		// ── ordered list ──────────────────────────────────────────────────────
		case reOrdered.MatchString(trimmed):
			out.WriteString("<ol class=\"rl-list\">\n")
			for i < len(lines) {
				t := strings.TrimSpace(lines[i])
				m := reOrdered.FindStringSubmatch(t)
				if m == nil {
					break
				}
				item := m[1]
				i++
				for i < len(lines) {
					n := strings.TrimSpace(lines[i])
					if n == "" || reOrdered.MatchString(n) || strings.HasPrefix(n, "#") ||
						strings.HasPrefix(n, "```") || strings.HasPrefix(n, "|") {
						break
					}
					item += " " + n
					i++
				}
				out.WriteString("<li>" + inline(item) + "</li>\n")
			}
			out.WriteString("</ol>\n")

		// ── paragraph ─────────────────────────────────────────────────────────
		default:
			var body []string
			for i < len(lines) {
				t := strings.TrimSpace(lines[i])
				if t == "" || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "```") ||
					strings.HasPrefix(t, ">") || strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* ") ||
					strings.HasPrefix(t, "|") || strings.HasPrefix(t, "<") || reOrdered.MatchString(t) {
					break
				}
				body = append(body, t)
				i++
			}
			if len(body) > 0 {
				out.WriteString("<p>" + inline(strings.Join(body, " ")) + "</p>\n")
			}
		}
	}

	return template.HTML(out.String()), sections
}

func splitRow(line string) []string {
	line = strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(line), "|"), "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// inline applies code spans, links, bold and italic. Code span contents are
// HTML-escaped and shielded from the other transforms.
func inline(s string) string {
	var codes []string
	s = reCodeSpan.ReplaceAllStringFunc(s, func(m string) string {
		inner := template.HTMLEscapeString(strings.Trim(m, "`"))
		codes = append(codes, inner)
		return fmt.Sprintf("\x00c%d\x00", len(codes)-1)
	})

	s = reLink.ReplaceAllString(s, `<a href="$2" target="_blank" rel="noopener">$1</a>`)
	s = reBold.ReplaceAllString(s, `<strong>$1</strong>`)
	s = reItalic.ReplaceAllString(s, `$1<em>$2</em>$3`)

	for i, c := range codes {
		s = strings.ReplaceAll(s, fmt.Sprintf("\x00c%d\x00", i), "<code>"+c+"</code>")
	}
	return s
}

// stripInline removes Markdown punctuation for use in a table of contents.
func stripInline(s string) string {
	s = strings.ReplaceAll(s, "`", "")
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "*", "")
	return strings.TrimSpace(s)
}

func slugify(s string) string {
	s = strings.ToLower(stripInline(s))
	s = reNonSlug.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// RenderInline renders a single line of Markdown, code spans, links, bold and
// italic, for short strings like a feature summary that live in JSON metadata.
func RenderInline(s string) template.HTML {
	return template.HTML(inline(s))
}
