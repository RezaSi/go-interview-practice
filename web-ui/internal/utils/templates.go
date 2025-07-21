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
		// New function to count attempts for a specific package
		"countPackageAttemptsForPackage": func(userAttempts interface{}, pkg interface{}) int {
			if userAttempts == nil || pkg == nil {
				return 0
			}

			// Get package name
			pkgValue := reflect.ValueOf(pkg)
			if pkgValue.Kind() == reflect.Ptr {
				pkgValue = pkgValue.Elem()
			}

			var packageName string
			if pkgValue.Kind() == reflect.Struct {
				nameField := pkgValue.FieldByName("Name")
				if nameField.IsValid() && nameField.Kind() == reflect.String {
					packageName = nameField.String()
				}
			}

			if packageName == "" {
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

					// Get learning path to know how many challenges this package has
					learningPathField := pkgValue.FieldByName("LearningPath")
					if !learningPathField.IsValid() {
						return 0
					}

					learningPath := learningPathField.Interface()
					if pathSlice, ok := learningPath.([]string); ok {
						// Check each challenge in the learning path
						for i := range pathSlice {
							// Generate the same unique ID as in the web handler
							packageChallengeID := -(1000 + i*10 + len(packageName))

							key := reflect.ValueOf(packageChallengeID)
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
		// New template functions for dynamic package rendering
		"getChallengeInfo": func(pkg interface{}, challengeID string) map[string]interface{} {
			// Extract challenge information dynamically from package
			v := reflect.ValueOf(pkg)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			challengeDetailsField := v.FieldByName("ChallengeDetails")
			if !challengeDetailsField.IsValid() || challengeDetailsField.Kind() != reflect.Map {
				return map[string]interface{}{
					"Title":       "Coming Soon",
					"Description": "Challenge content will be available soon",
					"Difficulty":  "Beginner",
					"Status":      "coming-soon",
					"Icon":        "bi-clock",
				}
			}

			key := reflect.ValueOf(challengeID)
			challenge := challengeDetailsField.MapIndex(key)

			if !challenge.IsValid() {
				return map[string]interface{}{
					"Title":       "Coming Soon",
					"Description": "Challenge content will be available soon",
					"Difficulty":  "Beginner",
					"Status":      "coming-soon",
					"Icon":        "bi-clock",
				}
			}

			// Convert challenge to map for template use
			challengeInfo := make(map[string]interface{})
			challengeValue := challenge.Elem()
			challengeType := challengeValue.Type()

			for i := 0; i < challengeValue.NumField(); i++ {
				field := challengeValue.Field(i)
				fieldName := challengeType.Field(i).Name
				challengeInfo[fieldName] = field.Interface()
			}

			return challengeInfo
		},
		"isPackageActive": func(pkg interface{}) bool {
			// Check if package has any available challenges
			v := reflect.ValueOf(pkg)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			challengeDetailsField := v.FieldByName("ChallengeDetails")
			if !challengeDetailsField.IsValid() || challengeDetailsField.Kind() != reflect.Map {
				return false
			}

			// Check if there are any challenges with "available" status
			for _, key := range challengeDetailsField.MapKeys() {
				challenge := challengeDetailsField.MapIndex(key)
				if challenge.IsValid() && challenge.Kind() == reflect.Ptr {
					statusField := challenge.Elem().FieldByName("Status")
					if statusField.IsValid() && statusField.String() == "available" {
						return true
					}
				}
			}

			return false
		},
		"getPackageChallenges": func(pkg interface{}) []map[string]interface{} {
			// Extract challenges from package and sort by order
			v := reflect.ValueOf(pkg)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			challengeDetailsField := v.FieldByName("ChallengeDetails")
			learningPathField := v.FieldByName("LearningPath")

			if !challengeDetailsField.IsValid() || !learningPathField.IsValid() {
				return []map[string]interface{}{}
			}

			var challenges []map[string]interface{}

			// Use learning path to maintain order
			learningPath := learningPathField.Interface()
			if pathSlice, ok := learningPath.([]string); ok {
				for _, challengeID := range pathSlice {
					key := reflect.ValueOf(challengeID)
					challenge := challengeDetailsField.MapIndex(key)

					if challenge.IsValid() && challenge.Kind() == reflect.Ptr {
						challengeInfo := make(map[string]interface{})
						challengeValue := challenge.Elem()
						challengeType := challengeValue.Type()

						for i := 0; i < challengeValue.NumField(); i++ {
							field := challengeValue.Field(i)
							fieldName := challengeType.Field(i).Name
							challengeInfo[fieldName] = field.Interface()
						}

						challenges = append(challenges, challengeInfo)
					}
				}
			}

			return challenges
		},
		"getDifficultyBadgeClass": func(difficulty string) string {
			switch strings.ToLower(difficulty) {
			case "beginner":
				return "bg-success"
			case "intermediate":
				return "bg-warning"
			case "advanced":
				return "bg-danger"
			default:
				return "bg-secondary"
			}
		},
		"getCategoryIcon": func(category string) string {
			switch strings.ToLower(category) {
			case "web":
				return "bi-globe2"
			case "database":
				return "bi-database"
			case "cli":
				return "bi-terminal"
			default:
				return "bi-lightning"
			}
		},
		"getCategoryGradient": func(category string) string {
			switch strings.ToLower(category) {
			case "web":
				return "bg-gradient-primary"
			case "database":
				return "bg-gradient-secondary"
			case "cli":
				return "bg-gradient-warning"
			default:
				return "bg-gradient-info"
			}
		},
		"isComingSoon": func(challengeInfo map[string]interface{}) bool {
			if status, ok := challengeInfo["Status"].(string); ok {
				return status == "coming-soon"
			}
			return false
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
