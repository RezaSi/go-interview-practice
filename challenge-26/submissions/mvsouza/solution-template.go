package regex

import (
	"regexp"
)

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	// TODO: Implement this function
	// 1. Create a regular expression to match email addresses
	// 2. Find all matches in the input text
	// 3. Return the matched emails as a slice of strings
	re := regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9-]+(?:\.[A-Za-z0-9-]+)+`)
	matches := re.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}
	return matches
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	// TODO: Implement this function
	// 1. Create a regular expression to match the specified phone format
	// 2. Check if the input string matches the pattern
	// 3. Return true if it's a match, false otherwise
	b, er := regexp.MatchString(`^\([0-9]{3}\) [0-9]{3}\-[0-9]{4}$`, phone)
	if er != nil {
		return false
	}
	return b
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	countReg := regexp.MustCompile(`\d{4}`)
	matches := countReg.FindAllString(cardNumber, -1)
	totalGroups := len(matches)

	if totalGroups <= 1 {
		return cardNumber
	}
	digitIndex := 0
	maskReg := regexp.MustCompile(`\d{4}`)

	return maskReg.ReplaceAllStringFunc(cardNumber, func(m string) string {
		digitIndex++
		if digitIndex <= totalGroups-1 {
			return "XXXX"
		}
		return m
	})
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	logRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}) (\d{2}:\d{2}:\d{2}) (\w+) (.*)$`)
	matches := logRegex.FindStringSubmatch(logLine)
	if len(matches) < 5 {
		return nil
	}
	logMap := make(map[string]string)
	logMap["date"] = matches[1]
	logMap["time"] = matches[2]
	logMap["level"] = matches[3]
	logMap["message"] = matches[4]
	return logMap
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	// TODO: Implement this function
	// 1. Create a regular expression to match URLs (both http and https)
	// 2. Find all matches in the input text
	// 3. Return the matched URLs as a slice of strings
	re := regexp.MustCompile(`https?://[\w:@-]+(\.[a-zA-Z0-9]{2,4})*(:\d+)?(/[\w.?#&=]*)*`)
	matches := re.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}
	return matches
}
