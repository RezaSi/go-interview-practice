package regex

import (
	"regexp"
	"unicode"
)

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	// 1. Create a regular expression to match email addresses
	re := regexp.MustCompile(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}`)

	// 2. Find all matches in the input text
	allEmails := re.FindAllString(text, -1)
	if allEmails == nil {
		return []string{}
	}

	// 3. Return the matched emails as a slice of strings
	return allEmails
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	// 1. Create a regular expression to match the specified phone format
	re := regexp.MustCompile(`^\(\d{3}\)\s{1}\d{3}-\d{4}$`)
	// 2. Check if the input string matches the pattern
	// 3. Return true if it's a match, false otherwise
	return re.MatchString(phone)

}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	digitCount := 0
	for _, c := range cardNumber {
		if unicode.IsDigit(c) {
			digitCount += 1
		}
	}

	// 1. Create a regular expression to identify the parts of the card number to mask
	re := regexp.MustCompile(`\d`)

	// 2. Use ReplaceAllString or similar method to perform the replacement
	seen := 0
	maskedCard := re.ReplaceAllStringFunc(cardNumber, func(s string) string {
		seen += 1
		if seen > digitCount-4 {
			return s
		}
		return "X"
	})

	// 3. Return the masked card number
	return maskedCard
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	// 1. Create a regular expression with capture groups for each component
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})\s(\d{2}:\d{2}:\d{2})\s(\w+)\s(.*)$`)
	// 2. Use FindStringSubmatch to extract the components
	matches := re.FindStringSubmatch(logLine)
	if len(matches) == 0 {
		return nil
	}

	// 3. Populate a map with the extracted values
	// 4. Return the populated map
	return map[string]string{
		"date":    matches[1],
		"time":    matches[2],
		"level":   matches[3],
		"message": matches[4],
	}
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	// 1. Create a regular expression to match URLs (both http and https)
	re := regexp.MustCompile(`https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.?[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9@:%_\+.~#?&\/=]*)`)
	// 2. Find all matches in the input text
	URLs := re.FindAllString(text, -1)
	if len(URLs) == 0 {
		return []string{}
	}
	// 3. Return the matched URLs as a slice of strings

	return URLs
}
