package privacy

import (
	"testing"
)

func TestPIIDetector_Anonymize(t *testing.T) {
	detector := NewPIIDetector()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Email Detection",
			input:    "Please contact me at user@example.com for more information.",
			expected: "Please contact me at [EMAIL] for more information.",
		},
		{
			name:     "Phone Number Detection",
			input:    "Call me at +1-123-456-7890 or 123.456.7890.",
			expected: "Call me at [PHONE] or [PHONE].",
		},
		{
			name:     "IP Address Detection",
			input:    "Server IP: 192.168.1.1 and 2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: "Server IP: [IP_ADDRESS] and [IP_ADDRESS]",
		},
		{
			name:     "Credit Card Detection",
			input:    "Payment with card: 4111-1111-1111-1111 or 4111111111111111",
			expected: "Payment with card: [CREDIT_CARD] or [CREDIT_CARD]",
		},
		{
			name:     "SSN Detection",
			input:    "SSN: 123-45-6789 or 123456789",
			expected: "SSN: [SSN] or [SSN]",
		},
		{
			name:     "Multiple PII Types",
			input:    "User: john@example.com, Phone: 555-123-4567, IP: 10.0.0.1, CC: 5555-5555-5555-4444, SSN: 987-65-4321",
			expected: "User: [EMAIL], Phone: [PHONE], IP: [IP_ADDRESS], CC: [CREDIT_CARD], SSN: [SSN]",
		},
		{
			name:     "No PII",
			input:    "This text contains no personal information.",
			expected: "This text contains no personal information.",
		},
		{
			name:     "Empty String",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.Anonymize(tt.input)
			if result != tt.expected {
				t.Errorf("Anonymize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPIIDetector_EmailDetection(t *testing.T) {
	detector := NewPIIDetector()

	emails := []string{
		"user@example.com",
		"first.last@example.co.uk",
		"user+tag@example.org",
		"user.name@subdomain.example.com",
	}

	for _, email := range emails {
		result := detector.Anonymize(email)
		if result != "[EMAIL]" {
			t.Errorf("Failed to detect email %s, got %s", email, result)
		}
	}
}

func TestPIIDetector_PhoneDetection(t *testing.T) {
	detector := NewPIIDetector()

	phones := []string{
		"+1-123-456-7890",
		"123.456.7890",
		"(123) 456-7890",
		"123 456 7890",
		"+886 912 345 678", // 台灣手機號碼格式
	}

	for _, phone := range phones {
		result := detector.Anonymize(phone)
		if result != "[PHONE]" {
			t.Errorf("Failed to detect phone %s, got %s", phone, result)
		}
	}
}

func TestPIIDetector_IPAddressDetection(t *testing.T) {
	detector := NewPIIDetector()

	ips := []string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2001:db8::1",
	}

	for _, ip := range ips {
		result := detector.Anonymize(ip)
		if result != "[IP_ADDRESS]" {
			t.Errorf("Failed to detect IP %s, got %s", ip, result)
		}
	}
}

func TestPIIDetector_CreditCardDetection(t *testing.T) {
	detector := NewPIIDetector()

	cards := []string{
		"4111-1111-1111-1111",
		"4111111111111111",
		"5555-5555-5555-4444",
		"5555555555554444",
		"3782-822463-10005",
		"378282246310005",
	}

	for _, card := range cards {
		result := detector.Anonymize(card)
		if result != "[CREDIT_CARD]" {
			t.Errorf("Failed to detect credit card %s, got %s", card, result)
		}
	}
}

func TestPIIDetector_SSNDetection(t *testing.T) {
	detector := NewPIIDetector()

	ssns := []string{
		"123-45-6789",
		"123456789",
	}

	for _, ssn := range ssns {
		result := detector.Anonymize(ssn)
		if result != "[SSN]" {
			t.Errorf("Failed to detect SSN %s, got %s", ssn, result)
		}
	}
}
