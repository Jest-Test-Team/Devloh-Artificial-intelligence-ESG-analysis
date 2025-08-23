package privacy

import (
	"regexp"
)

// PIIDetector detects and anonymizes Personally Identifiable Information
type PIIDetector struct {
	emailRegex    *regexp.Regexp
	phoneRegex    *regexp.Regexp
	creditCardRegex *regexp.Regexp
	ssnRegex      *regexp.Regexp
}

// NewPIIDetector creates a new PII detector instance
func NewPIIDetector() *PIIDetector {
	return &PIIDetector{
		emailRegex:    regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		phoneRegex:    regexp.MustCompile(`(\+?1[-.\s]?)?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})`),
		creditCardRegex: regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`),
		ssnRegex:      regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
	}
}

// DetectAndAnonymize scans text for PII and replaces it with anonymized versions
func (p *PIIDetector) DetectAndAnonymize(text string) string {
	result := text
	
	// Anonymize emails
	result = p.emailRegex.ReplaceAllString(result, "[EMAIL_REDACTED]")
	
	// Anonymize phone numbers
	result = p.phoneRegex.ReplaceAllString(result, "[PHONE_REDACTED]")
	
	// Anonymize credit card numbers
	result = p.creditCardRegex.ReplaceAllString(result, "[CREDIT_CARD_REDACTED]")
	
	// Anonymize SSNs
	result = p.ssnRegex.ReplaceAllString(result, "[SSN_REDACTED]")
	
	return result
}

// DetectPII scans text and returns a list of detected PII types
func (p *PIIDetector) DetectPII(text string) []string {
	var detected []string
	
	if p.emailRegex.MatchString(text) {
		detected = append(detected, "email")
	}
	
	if p.phoneRegex.MatchString(text) {
		detected = append(detected, "phone")
	}
	
	if p.creditCardRegex.MatchString(text) {
		detected = append(detected, "credit_card")
	}
	
	if p.ssnRegex.MatchString(text) {
		detected = append(detected, "ssn")
	}
	
	return detected
}

// IsPII checks if a specific text contains PII
func (p *PIIDetector) IsPII(text string) bool {
	return len(p.DetectPII(text)) > 0
}

// AnonymizeLogEntry anonymizes sensitive fields in a log entry
func (p *PIIDetector) AnonymizeLogEntry(logEntry map[string]interface{}) map[string]interface{} {
	anonymized := make(map[string]interface{})
	
	for key, value := range logEntry {
		switch v := value.(type) {
		case string:
			anonymized[key] = p.DetectAndAnonymize(v)
		case map[string]interface{}:
			anonymized[key] = p.AnonymizeLogEntry(v)
		default:
			anonymized[key] = v
		}
	}
	
	return anonymized
}
