package utils

import (
	"fmt"
	"html/template"
	"reflect"
	"regexp"
	"strings"
)

// Simple markdown to HTML converter
func markdownToHTML(markdown string) string {
	html := markdown

	// Convert headers
	html = regexp.MustCompile(`(?m)^#{6}\s+(.+)$`).ReplaceAllString(html, "<h6>$1</h6>")
	html = regexp.MustCompile(`(?m)^#{5}\s+(.+)$`).ReplaceAllString(html, "<h5>$1</h5>")
	html = regexp.MustCompile(`(?m)^#{4}\s+(.+)$`).ReplaceAllString(html, "<h4>$1</h4>")
	html = regexp.MustCompile(`(?m)^#{3}\s+(.+)$`).ReplaceAllString(html, "<h3>$1</h3>")
	html = regexp.MustCompile(`(?m)^#{2}\s+(.+)$`).ReplaceAllString(html, "<h2>$1</h2>")
	html = regexp.MustCompile(`(?m)^#{1}\s+(.+)$`).ReplaceAllString(html, "<h1>$1</h1>")

	// Convert bold text
	html = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(html, "<strong>$1</strong>")

	// Convert italic text
	html = regexp.MustCompile(`\*([^*]+)\*`).ReplaceAllString(html, "<em>$1</em>")

	// Convert code blocks
	html = regexp.MustCompile("(?s)```(\\w*)\\n(.*?)```").ReplaceAllStringFunc(html, func(match string) string {
		parts := regexp.MustCompile("(?s)```(\\w*)\\n(.*?)```").FindStringSubmatch(match)
		if len(parts) >= 3 {
			language := parts[1]
			code := strings.TrimSpace(parts[2])
			if language != "" {
				return fmt.Sprintf(`<pre><code class="language-%s">%s</code></pre>`, language, template.HTMLEscapeString(code))
			}
			return fmt.Sprintf(`<pre><code>%s</code></pre>`, template.HTMLEscapeString(code))
		}
		return match
	})

	// Convert inline code
	html = regexp.MustCompile("`([^`]+)`").ReplaceAllString(html, "<code>$1</code>")

	// Convert links
	html = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`).ReplaceAllString(html, `<a href="$2" target="_blank">$1</a>`)

	// Convert lists (simple implementation)
	lines := strings.Split(html, "\n")
	var result []string
	inList := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle numbered lists
		if matched, _ := regexp.MatchString(`^\d+\.\s+`, trimmed); matched {
			if !inList {
				result = append(result, "<ol>")
				inList = true
			}
			content := regexp.MustCompile(`^\d+\.\s+`).ReplaceAllString(trimmed, "")
			result = append(result, fmt.Sprintf("<li>%s</li>", content))
		} else if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			if !inList {
				result = append(result, "<ul>")
				inList = true
			}
			content := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* ")
			result = append(result, fmt.Sprintf("<li>%s</li>", content))
		} else {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			if trimmed != "" {
				// Check if the line is already HTML (header, code block, etc.)
				if strings.HasPrefix(trimmed, "<h") || strings.HasPrefix(trimmed, "<pre") ||
					strings.HasPrefix(trimmed, "<div") || strings.HasPrefix(trimmed, "<code") ||
					strings.HasPrefix(trimmed, "<blockquote") || strings.HasPrefix(trimmed, "<hr") {
					result = append(result, trimmed)
				} else {
					result = append(result, fmt.Sprintf("<p>%s</p>", trimmed))
				}
			}
		}
	}

	if inList {
		result = append(result, "</ul>")
	}

	return strings.Join(result, "\n")
}

// GetTemplateFuncs returns the template functions used across the application
func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"lower": strings.ToLower,
		"truncateDescription": func(s string) string {
			// Extract first paragraph that is not a heading or link
			lines := strings.Split(s, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "[") {
					continue
				}
				// Found an actual paragraph
				if len(line) > 150 {
					return line[:150] + "..."
				}
				return line
			}

			// Fallback to simple truncation
			if len(s) > 150 {
				return s[:150] + "..."
			}
			return s
		},
		"add": func(a, b int) int {
			return a + b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"calculateProgress": func(attemptedCount, totalCount int) int {
			if totalCount == 0 {
				return 0
			}
			return (attemptedCount * 100) / totalCount
		},
		"calculatePercentage": func(passed, total int) int {
			if total == 0 {
				return 0
			}
			return (passed * 100) / total
		},
		"countPackageAttempts": func(userAttempts interface{}) int {
			if userAttempts == nil {
				return 0
			}

			// Use reflection to access the AttemptedIDs field
			v := reflect.ValueOf(userAttempts)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			if v.Kind() == reflect.Struct {
				field := v.FieldByName("AttemptedIDs")
				if field.IsValid() && field.Kind() == reflect.Map {
					count := 0
					for _, key := range field.MapKeys() {
						if key.Kind() == reflect.Int && key.Int() < 0 {
							value := field.MapIndex(key)
							if value.IsValid() && value.Kind() == reflect.Bool && value.Bool() {
								count++
							}
						}
					}
					return count
				}
			}
			return 0
		},
		"extractTitle": func(description string) string {
			// Extract title from markdown content
			lines := strings.Split(description, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "# ") {
					return strings.TrimPrefix(line, "# ")
				}
			}
			return ""
		},
		"js": func(s string) template.JS {
			// Safely escape backticks and other special characters for JavaScript
			// Replace backticks with HTML entity
			s = strings.Replace(s, "`", "\\`", -1)
			// Replace dollar signs that might interfere with template literals
			s = strings.Replace(s, "${", "\\${", -1)
			return template.JS(s)
		},
		"replace": func(old, new, str string) string {
			return strings.Replace(str, old, new, -1)
		},
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"markdown": func(s string) template.HTML {
			return template.HTML(markdownToHTML(s))
		},
		"formatStars": func(stars int) string {
			if stars >= 1000000 {
				return fmt.Sprintf("%.1fM", float64(stars)/1000000)
			} else if stars >= 1000 {
				return fmt.Sprintf("%.1fk", float64(stars)/1000)
			}
			return fmt.Sprintf("%d", stars)
		},
		"truncate": func(length int, s string) string {
			if len(s) <= length {
				return s
			}
			if length <= 3 {
				return s[:length]
			}
			return s[:length-3] + "..."
		},
	}
}
