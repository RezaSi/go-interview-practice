package regex

import "regexp"

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	re := regexp.MustCompile(`[\w._%+\-]+@[\w.\-]+\.[a-zA-Z]{2,}`)
	results := re.FindAllString(text, -1)
	if results == nil {
		return []string{}
	}
	return results
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	return regexp.MustCompile(`^\(\d{3}\)\s\d{3}-\d{4}$`).MatchString(phone)
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	re := regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-(\d{4})$`)
	if re.MatchString(cardNumber) {
		return re.ReplaceAllString(cardNumber, "XXXX-XXXX-XXXX-${1}")
	}

	re = regexp.MustCompile(`^\d{12}(\d{4})$`)
	if re.MatchString(cardNumber) {
		return re.ReplaceAllString(cardNumber, "XXXXXXXXXXXX${1}")
	}

	re = regexp.MustCompile(`^\d{4}-(\d{4})$`)
	if re.MatchString(cardNumber) {
		return re.ReplaceAllString(cardNumber, "XXXX-${1}")
	}

	return cardNumber
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})\s(\d{2}:\d{2}:\d{2})\s([A-Z]+)\s(.+)$`)
	matches := re.FindStringSubmatch(logLine)
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

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	re := regexp.MustCompile(
		`https?://` + 
		`([\w~\-]+:[\w~\-]+@)?` +
		`([a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}|[a-zA-Z0-9\-]+)` +
		`(:[\d]{1,5})?` +
		`(/[\w\-\.]+)*` +
		`(\?[\w\-&=]+)*` +
		`(#[\w\-]+)*`)
	results := re.FindAllString(text, -1)
	if results == nil {
		return []string{}
	}
	return results
}
