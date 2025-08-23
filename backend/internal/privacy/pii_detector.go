package privacy

import (
	"regexp"
)

// PIIDetector 用於檢測和匿名化個人識別信息（PII）
type PIIDetector struct {
	patterns map[string]*regexp.Regexp
}

// NewPIIDetector 創建一個新的 PII 檢測器
func NewPIIDetector() *PIIDetector {
	detector := &PIIDetector{
		patterns: make(map[string]*regexp.Regexp),
	}
	
	// 編譯常見 PII 的正則表達式模式
	detector.patterns["EMAIL"] = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`)
	detector.patterns["PHONE"] = regexp.MustCompile(`\b(\+\d{1,3}[-.\s]?)?(\(?\d{1,4}\)?[-.\s]?)?(\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9})\b`)
	detector.patterns["CREDIT_CARD"] = regexp.MustCompile(`\b(\d{4}[-\s]?){3}\d{4}\b|\b\d{16}\b`)
	detector.patterns["SSN"] = regexp.MustCompile(`\b\d{3}[-\s]?\d{2}[-\s]?\d{4}\b|\b\d{9}\b`)
	detector.patterns["IP_ADDRESS"] = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b|\b([0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4}\b`)
	
	return detector
}

// Anonymize 匿名化字符串中的 PII
func (d *PIIDetector) Anonymize(text string) string {
	result := text
	
	// 對每種 PII 類型進行匿名化
	for piiType, pattern := range d.patterns {
		result = pattern.ReplaceAllString(result, "["+piiType+"]")
	}
	
	return result
}

// DetectPII 檢測字符串中的 PII 類型
func (d *PIIDetector) DetectPII(text string) []string {
	var detectedTypes []string
	
	// 檢查每種 PII 類型
	for piiType, pattern := range d.patterns {
		if pattern.MatchString(text) {
			detectedTypes = append(detectedTypes, piiType)
		}
	}
	
	return detectedTypes
}

// IsPII 檢查文本是否包含任何 PII
func (d *PIIDetector) IsPII(text string) bool {
	for _, pattern := range d.patterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}

// DetectAndAnonymize 檢測並匿名化文本中的 PII
func (d *PIIDetector) DetectAndAnonymize(text string) (string, []string) {
	detectedTypes := d.DetectPII(text)
	anonymizedText := d.Anonymize(text)
	return anonymizedText, detectedTypes
}

// AnonymizeLogEntry 匿名化日誌條目中的 PII
func (d *PIIDetector) AnonymizeLogEntry(entry map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	// 複製所有字段
	for k, v := range entry {
		result[k] = v
	}
	
	// 匿名化消息字段
	if message, ok := entry["message"].(string); ok {
		result["message"] = d.Anonymize(message)
	}
	
	// 匿名化其他可能包含 PII 的字段
	if details, ok := entry["details"].(string); ok {
		result["details"] = d.Anonymize(details)
	}
	
	if userInfo, ok := entry["user_info"].(string); ok {
		result["user_info"] = d.Anonymize(userInfo)
	}
	
	return result
}

// AnonymizePII 是一個方便的函數，用於匿名化文本中的 PII
func AnonymizePII(text string) string {
	detector := NewPIIDetector()
	return detector.Anonymize(text)
}