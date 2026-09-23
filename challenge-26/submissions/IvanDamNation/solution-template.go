package regex

import "regexp"

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	re := regexp.MustCompile(`[a-zA-Z0-9.%+-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,}`)
	matches := re.FindAllString(text, -1)
	if matches == nil {
	    return []string{}
	}

	return matches
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	re := regexp.MustCompile(`^\([0-9]{3}\) [0-9]{3}-[0-9]{4}$`)
	return re.MatchString(phone)
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	re := regexp.MustCompile(`\d`)
	numDigits := len(re.FindAllString(cardNumber, -1))

    digitCount := 0
	return re.ReplaceAllStringFunc(cardNumber, func(digit string) string {
	    digitCount++
	    if digitCount <= numDigits - 4 {
	        return "X"
	    }
	    return digit
	})
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	re := regexp.MustCompile(`^(?P<date>\d{4}-\d{2}-\d{2}) (?P<time>\d{2}:\d{2}:\d{2}) (?P<level>[A-Z]+) (?P<message>.+)$`)
	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
	    return nil
	}
	
	parseMap := make(map[string]string, 4)
	for i, name := range re.SubexpNames() {
	    if name != "" {
	        parseMap[name] = matches[i]
	    }
	}
	
	return parseMap
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	re := regexp.MustCompile(`https?://[a-zA-Z0-9_:%%@~#=+\/?\.\-\&]+`)
	matches := re.FindAllString(text, -1)
	if matches == nil {
	    return []string{}
	}

	return matches
}
