package regex

import "regexp"

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	if text == "" {
		return []string{}
	}
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	res := re.FindAllString(text, -1)
	if res == nil {
		return []string{}
	}
	return res
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	if phone == "" {
		return false
	}
	re := regexp.MustCompile(`^\(\d{3}\) \d{3}-\d{4}$`)
	return re.MatchString(phone)
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	if cardNumber == "" {
		return ""
	}
	re := regexp.MustCompile(`\d`)
	positions := re.FindAllStringIndex(cardNumber, -1)
	keepFrom := len(positions) - 4
	if keepFrom <= 0 {
		return cardNumber
	}
	result := []byte(cardNumber)
	for i, pos := range positions {
		if i < keepFrom {
			result[pos[0]] = 'X'
		}
	}
	return string(result)
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}) (\d{2}:\d{2}:\d{2}) (\w+) (.+)$`)
	m := re.FindStringSubmatch(logLine)
	if m == nil {
		return nil
	}
	return map[string]string{
		"date":    m[1],
		"time":    m[2],
		"level":   m[3],
		"message": m[4],
	}
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	if text == "" {
		return []string{}
	}
	// Optional userinfo (user:pass@) before host; path/query/fragment exclude quotes and brackets.
	re := regexp.MustCompile(`https?://(?:[^\s@/?#]+@)?[a-zA-Z0-9.-]+(?::\d+)?(?:/[^\s?#'")\]>]*)?(?:\?[^\s#'")\]>]*)?(?:#[^\s'")\]>]*)?`)
	result := re.FindAllString(text, -1)
	if result == nil {
		return []string{}
	}
	return result
}
