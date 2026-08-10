// Challenge 26: Regular Expression Text Processor

package regex

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	phoneRegex = regexp.MustCompile(`^\(\d{3}\) \d{3}-\d{4}$`) // format (XXX) XXX-XXXX
	urlRegex   = regexp.MustCompile(`https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(?:/[^\s]*)?`)
	logPattern = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}) (\d{2}:\d{2}:\d{2}) (\w+) (.+)$`)
)

func ExtractEmails(text string) []string {
	res := emailRegex.FindAllString(text, -1)
	if res == nil {
		return []string{}
	}
	return res
}
func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	// TODO: Implement this function
	// 1. Create a regular expression to identify the parts of the card number to mask
	// 2. Use ReplaceAllString or similar method to perform the replacement
	// 3. Return the masked card number
	re := regexp.MustCompile(`\d`)
	digitsSeen := 0
	for i := len(cardNumber) - 1; i >= 0; i-- {
		if cardNumber[i] >= '0' && cardNumber[i] <= '9' {
			digitsSeen++
		}
	}
	remaining := digitsSeen
	return re.ReplaceAllStringFunc(cardNumber, func(s string) string {
		if remaining > 4 {
			remaining--
			return "X"
		}
		remaining--
		return s
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
	// logLine := "2023-11-15 14:23:45 INFO Server started on port 8080",
	// Define regex to parse log entries
	matches := logPattern.FindStringSubmatch(logLine)
	if len(matches) != 5 {
		return nil
	}
	return map[string]string{
		"date":    matches[1],
		"time":    matches[2],
		"level":   matches[3],
		"message": matches[4],
	}
}
func ExtractURLs(text string) []string {
	urlRegex := regexp.MustCompile(
		`https?://[^\s<>'"]+`,
	)

	urls := urlRegex.FindAllString(text, -1)
	if urls == nil {
		return []string{}
	}
	for i := range urls {
		urls[i] = strings.TrimRight(
			urls[i],
			".,;:!?)]}",
		)
	}

	return urls
}
