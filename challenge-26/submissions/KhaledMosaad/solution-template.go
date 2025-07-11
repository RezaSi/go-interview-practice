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

    emailReg := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
    
	emails := emailReg.FindAllString(text, -1)

	if emails == nil {
		return []string{}
	}

	return emails
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	// TODO: Implement this function
	// 1. Create a regular expression to match the specified phone format
	// 2. Check if the input string matches the pattern
	// 3. Return true if it's a match, false otherwise
	
	phoneReg := regexp.MustCompile(`^\(\d{3}\)\s\d{3}-\d{4}$`)

	return phoneReg.MatchString(phone)
}

// It removes the empty strings from array of strings
func RemoveEmptyElements(slice []string) []string {
	result := make([]string, 0, len(slice))

	for _, element := range slice {
		if element != "" {
			result = append(result, element)
		}
	}

	return result
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	// TODO: Implement this function
	// 1. Create a regular expression to identify the parts of the card number to mask
	// 2. Use ReplaceAllString or similar method to perform the replacement
	// 3. Return the masked card number
	cardReg := regexp.MustCompile(`^(\d{4}|XXXX)(-)?(\d{4}|XXXX)?(-)?(\d{4}|XXXX)?(-)?(\d{4})?$`)

	matches := cardReg.FindStringSubmatch(cardNumber)
	matches = RemoveEmptyElements(matches)
	// first and last substring in the array won't be replaced

	matchesLen := len(matches)
	if matchesLen == 0 {
		return ""
	}

	result := strings.Builder{}

	for i := 1; i < matchesLen-1; i++ {
		if matches[i] == "-" {
			result.WriteString("-")
		} else {
			result.WriteString("XXXX")
		}
	}

	result.WriteString(matches[matchesLen-1])
	return result.String()
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
	logReg := regexp.MustCompile(`(?P<date>\d{4}-\d{2}-\d{2})\s(?P<time>\d{2}:\d{2}:\d{2})\s(?P<level>\w+)\s(?P<message>.*)`)
	matches := logReg.FindStringSubmatch(logLine)
	if len(matches) > 0 {
		result := make(map[string]string)
		for i, name := range logReg.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = matches[i]
			}
		}
		return result
	}

	return nil
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	// TODO: Implement this function
	// 1. Create a regular expression to match URLs (both http and https)
	// 2. Find all matches in the input text
	// 3. Return the matched URLs as a slice of strings
	urlRegex := regexp.MustCompile(`(https?://)` + // Protocol
		`(?:[a-zA-Z0-9-._~%]+(?::[a-zA-Z0-9-._~%]+)?@)?` + // Optional userinfo
		`(?:` + // Start of host (domain or IP)
		`(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+` + // Domain labels
		`[a-zA-Z]{2,63}` + // TLD (up to 63 chars for new TLDs)
		`|` + // OR
		`localhost` + // Match localhost
		`|` + // OR
		`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}` + // IPv4
		`)(?::\d{1,5})?` + // Optional port
		`(?:` + // Start of optional path/query/fragment
		`/?` + // Optional leading slash for path
		`(?:[a-zA-Z0-9\-._~:/?#\[@!$&(*+;=%]*)` + // Characters for path, query, fragment
		`)`)
	urls := urlRegex.FindAllString(text, -1)

	if urls == nil {
		return []string{}
	}

	return urls
}
