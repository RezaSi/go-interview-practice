package regex

import (
	"regexp"
	"strings"
)

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	// TODO: Implement this function
	// 1. Create a regular expression to match email addresses
	// 2. Find all matches in the input text
	// 3. Return the matched emails as a slice of strings

	// Email validation (simplified)
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	result := re.FindAllString(text, -1)

	if result == nil {
		return []string{}
	}
	return result
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	// TODO: Implement this function
	// 1. Create a regular expression to match the specified phone format
	// 2. Check if the input string matches the pattern
	// 3. Return true if it's a match, false otherwise

	// Phone number validation (US format)
	re := regexp.MustCompile(`^\(\d{3}\) \d{3}-\d{4}$`)
	result := re.MatchString(phone) // true

	return result
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	// TODO: Implement this function
	// 1. Create a regular expression to identify the parts of the card number to mask
	// 2. Use ReplaceAllString or similar method to perform the replacement
	// 3. Return the masked card number

	re := regexp.MustCompile(`^(.*)(\d{4})([^\d]*)$`)
	return re.ReplaceAllStringFunc(cardNumber, func(s string) string {
		parts := re.FindStringSubmatch(s)
		if len(parts) != 4 {
			return s
		}

		masked := regexp.MustCompile(`\d`).ReplaceAllString(parts[1], "X")
		return masked + parts[2] + parts[3]
	})
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	// TODO: Implement this function
	// 1. Create a regular expression with capture groups for each component
	// 2. Use FindStringSubmatch to extract the components
	// 3. Populate a map with the extracted values
	// 4. Return the populated map

	// Define regex to parse log entries
	logPattern := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}) (\d{2}:\d{2}:\d{2}) (\w+) (.+)$`)

	matches := logPattern.FindStringSubmatch(logLine)
	if len(matches) == 5 {
		date := matches[1]
		time := matches[2]
		level := matches[3]
		message := matches[4]

		return map[string]string{
			"date":    date,
			"time":    time,
			"level":   level,
			"message": message,
		}

	}

	return nil
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	// TODO: Implement this function
	// 1. Create a regular expression to match URLs (both http and https)
	// 2. Find all matches in the input text
	// 3. Return the matched URLs as a slice of strings

	re := regexp.MustCompile(`https?://[^\s<>"']+`)
	matches := re.FindAllString(text, -1)
	result := make([]string, 0, len(matches))
	for _, url := range matches {
		url = strings.TrimRight(url, ".,;:!?')]}\"")
		result = append(result, url)
	}

	return result

}
