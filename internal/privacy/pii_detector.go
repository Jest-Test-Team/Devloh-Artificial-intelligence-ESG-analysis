package privacy

import (
	"regexp"
	"strings"
)

var (
	// Regular expressions for various PII data types
	emailRegex    = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	phoneRegex    = regexp.MustCompile(`(\+\d{1,3}[- ]?)?\(?\d{3}\)?[- ]?\d{3}[- ]?\d{4}`)
	ssnRegex      = regexp.MustCompile(`\d{3}-\d{2}-\d{4}`)
	ipRegex       = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	creditRegex   = regexp.MustCompile(`\b(?:\d{4}[- ]?){3}\d{4}\b`)
	addressRegex  = regexp.MustCompile(`\d+\s+[A-Za-z0-9\s,'-]+\b(?:street|st|avenue|ave|road|rd|boulevard|blvd|drive|dr|court|ct|lane|ln|way)\b`)
	nameRegex     = regexp.MustCompile(`\b(?:[A-Z][a-z]+\s+){1,2}[A-Z][a-z]+\b`)
	passportRegex = regexp.MustCompile(`\b[A-Z]{1,2}\d{6,9}\b`)
)

// DetectPII searches for PII data in the given text
func DetectPII(text string) map[string][]string {
	result := make(map[string][]string)

	emails := emailRegex.FindAllString(text, -1)
	if len(emails) > 0 {
		result["emails"] = emails
	}

	phones := phoneRegex.FindAllString(text, -1)
	if len(phones) > 0 {
		result["phones"] = phones
	}

	ssns := ssnRegex.FindAllString(text, -1)
	if len(ssns) > 0 {
		result["ssns"] = ssns
	}

	ips := ipRegex.FindAllString(text, -1)
	if len(ips) > 0 {
		result["ips"] = ips
	}

	creditCards := creditRegex.FindAllString(text, -1)
	if len(creditCards) > 0 {
		result["credit_cards"] = creditCards
	}

	addresses := addressRegex.FindAllString(text, -1)
	if len(addresses) > 0 {
		result["addresses"] = addresses
	}

	names := nameRegex.FindAllString(text, -1)
	if len(names) > 0 {
		result["names"] = names
	}

	passports := passportRegex.FindAllString(text, -1)
	if len(passports) > 0 {
		result["passports"] = passports
	}

	return result
}

// AnonymizePII replaces PII data in the given text with placeholders
func AnonymizePII(text string) string {
	// Detect PII data
	piiData := DetectPII(text)

	// Replace each type of PII with a placeholder
	anonymizedText := text

	for piiType, instances := range piiData {
		for _, instance := range instances {
			var placeholder string

			switch piiType {
			case "emails":
				placeholder = "[EMAIL]"
			case "phones":
				placeholder = "[PHONE_NUMBER]"
			case "ssns":
				placeholder = "[SSN]"
			case "ips":
				placeholder = "[IP_ADDRESS]"
			case "credit_cards":
				placeholder = "[CREDIT_CARD]"
			case "addresses":
				placeholder = "[ADDRESS]"
			case "names":
				placeholder = "[NAME]"
			case "passports":
				placeholder = "[PASSPORT]"
			default:
				placeholder = "[REDACTED]"
			}

			anonymizedText = strings.Replace(anonymizedText, instance, placeholder, -1)
		}
	}

	return anonymizedText
}
