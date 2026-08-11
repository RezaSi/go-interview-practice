package regex

import "regexp"

var emailRegexp = regexp.MustCompile(`[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}`)
var phoneRegexp = regexp.MustCompile(`^\(\d{3}\) \d{3}-\d{4}$`)
var creditCardWithHyphenRegexp = regexp.MustCompile(`(\d{4})(\D)`)
var creditCardWithoutHyphenRegexp = regexp.MustCompile(`^(\d{12})(\d{4})$`)
var logRegexp = regexp.MustCompile(`^(?<date>\d{4}-\d{2}-\d{2}) (?<time>\d{2}:\d{2}:\d{2}) (?<level>[a-zA-Z]+) (?<message>.+)$`)
var urlRegexp = regexp.MustCompile(`(?<protocol>https?)(://)(?<low_level_domain>[\w:@-]+)(?<top_level_domain>\.\w{2,63})*(?<port>:\d+)*(?<path>/\w*)*(?<query>\?[\w=#]+)*(?<extension>\.\w+)*`)

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	// 1. Create a regular expression to match email addresses
	// 2. Find all matches in the input text
	// 3. Return the matched emails as a slice of strings
	if found := emailRegexp.FindAllString(text, -1); len(found) != 0 {
	    return found
	}
	return []string{}
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	// 1. Create a regular expression to match the specified phone format
	// 2. Check if the input string matches the pattern
	// 3. Return true if it's a match, false otherwise
	return phoneRegexp.Match([]byte(phone))
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	// 1. Create a regular expression to identify the parts of the card number to mask
	// 2. Use ReplaceAllString or similar method to perform the replacement
	// 3. Return the masked card number
	if creditCardWithHyphenRegexp.Match([]byte(cardNumber)) {
	    return string(creditCardWithHyphenRegexp.ReplaceAll([]byte(cardNumber), []byte("XXXX$2")))
	}
	
	return string(creditCardWithoutHyphenRegexp.ReplaceAll([]byte(cardNumber), []byte("XXXXXXXXXXXX$2")))
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	// 1. Create a regular expression with capture groups for each component
	// 2. Use FindStringSubmatch to extract the components
	// 3. Populate a map with the extracted values
	// 4. Return the populated map
    parts := logRegexp.FindStringSubmatch(logLine)

    if len(parts) < 5 {
        return nil
    }

	return map[string]string{
	    "date": parts[1],
	    "time": parts[2],
	    "level": parts[3],
	    "message": parts[4],
	}
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	// 1. Create a regular expression to match URLs (both http and https)
	// 2. Find all matches in the input text
	// 3. Return the matched URLs as a slice of strings
	if found := urlRegexp.FindAllString(text, -1); len(found) != 0 {
	    return found
	}
	return []string{}
}
