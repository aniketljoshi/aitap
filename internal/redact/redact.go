package redact

import (
	"regexp"
	"strings"
)

var patterns = []*regexp.Regexp{
	regexp.MustCompile(`sk-[a-zA-Z0-9]{20,}`),                    // OpenAI keys
	regexp.MustCompile(`sk-ant-[a-zA-Z0-9\-]{20,}`),              // Anthropic keys
	regexp.MustCompile(`AIza[a-zA-Z0-9\-_]{35}`),                 // Google API keys
	regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`),                    // GitHub PAT
	regexp.MustCompile(`AKIA[A-Z0-9]{16}`),                       // AWS access key
	regexp.MustCompile(`Bearer\s+[a-zA-Z0-9\-_.]{20,}`),          // Bearer tokens
	regexp.MustCompile(`(?i)(password|secret|token)\s*[:=]\s*\S+`), // Generic secrets
}

// Redact replaces detected secrets with [REDACTED].
func Redact(s string) string {
	for _, p := range patterns {
		s = p.ReplaceAllStringFunc(s, func(match string) string {
			if len(match) <= 8 {
				return "[REDACTED]"
			}
			return match[:4] + strings.Repeat("*", len(match)-8) + match[len(match)-4:] 
		})
	}
	return s
}
