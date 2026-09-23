package regex

import "regexp"

var (
	emailRE      = regexp.MustCompile(`[a-zA-Z0-9.%+-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,}`)
	phoneRE      = regexp.MustCompile(`^\([0-9]{3}\) [0-9]{3}-[0-9]{4}$`)
	creditCardRE = regexp.MustCompile(`\d`)
    logEntryRE   = regexp.MustCompile(`^(?P<date>\d{4}-\d{2}-\d{2}) (?P<time>\d{2}:\d{2}:\d{2}) (?P<level>[A-Z]+) (?P<message>.+)$`)
	urlRE        = regexp.MustCompile(`https?://[a-zA-Z0-9_:%%@~#=+\/?\.\-\&!$;]+`)
)

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	matches := emailRE.FindAllString(text, -1)
	if matches == nil {
	    return []string{}
	}

	return matches
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	return phoneRE.MatchString(phone)
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	numDigits := len(creditCardRE.FindAllString(cardNumber, -1))

    digitCount := 0
	return creditCardRE.ReplaceAllStringFunc(cardNumber, func(digit string) string {
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
	matches := logEntryRE.FindStringSubmatch(logLine)
	if matches == nil {
	    return nil
	}
	
	parseMap := make(map[string]string, 4)
	for i, name := range logEntryRE.SubexpNames() {
	    if name != "" {
	        parseMap[name] = matches[i]
	    }
	}
	
	return parseMap
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	matches := urlRE.FindAllString(text, -1)
	if matches == nil {
	    return []string{}
	}

	return matches
}
