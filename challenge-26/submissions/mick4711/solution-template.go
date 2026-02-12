package regex

import (
	"regexp"
	"strings"
)

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	// 1. Create a regular expression to match email addresses
	reEmail := regexp.MustCompile(`(\w+(\+|\-|\.)?\w+)+@(\w+(\+|\-|\.)?)+`)

	// 2. Find all matches in the input text
	matches := reEmail.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}

	// 3. Return the matched emails as a slice of strings
	return matches
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	// 1. Create a regular expression to match the specified phone format
	rePhone := regexp.MustCompile(`^\(\d{3}\)\s\d{3}\-\d{4}$`)

	// 2. Check if the input string matches the pattern
	// 3. Return true if it's a match, false otherwise
	return rePhone.MatchString(phone)
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	// 1. Create a regular expression to identify the parts of the card number to mask
	re := regexp.MustCompile(`\w{4}`)
	groups := re.FindAllString(cardNumber, -1)

	// 2. Use ReplaceAllString or similar method to perform the replacement
	mask := "XXXX"
	for i, group := range groups {
		if i == len(groups)-1 {
			break
		}
		if group == mask {
			continue
		}
		cardNumber = strings.Replace(cardNumber, group, mask, 1)
	}

	// 3. Return the masked card number
	return cardNumber
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	// 1. Create a regular expression with capture groups for each component
	re := regexp.MustCompile(`^(\d{4}\-\d{2}\-\d{2})\s(\d{2}:\d{2}:\d{2})\s(\w+)\s(.+)$`)

	// 2. Use FindStringSubmatch to extract the components
	subs := re.FindStringSubmatch(logLine)
	if len(subs) != 5 {
		return nil
	}

	// 3. Populate a map with the extracted values
	comps := make(map[string]string)
	comps["date"] = subs[1]
	comps["time"] = subs[2]
	comps["level"] = subs[3]
	comps["message"] = subs[4]

	// 4. Return the populated map
	return comps
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	// 1. Create a regular expression to match URLs (both http and https)
	re := regexp.MustCompile(`(http(s?):\/\/)([^\s^\)^\]^\,^\']+)`)

	// 2. Find all matches in the input text
	matches := re.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}

	// 3. Return the matched URLs as a slice of strings
	return matches
}
