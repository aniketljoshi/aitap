package redact

import (
	"strings"
	"testing"
)

func TestRedactOpenAIKey(t *testing.T) {
	input := "My key is sk-abcdefghijklmnopqrstuvwxyz123456"
	result := Redact(input)
	if strings.Contains(result, "abcdefghijklmnopqrstuvwxyz") {
		t.Error("OpenAI key was not redacted")
	}
	if !strings.Contains(result, "sk-a") {
		t.Error("first 4 chars should be preserved")
	}
}

func TestRedactAnthropicKey(t *testing.T) {
	input := "key: sk-ant-api03-abcdefghijklmnopqrstuvwxyz"
	result := Redact(input)
	if strings.Contains(result, "abcdefghijklmnopqrstuvwxyz") {
		t.Error("Anthropic key was not redacted")
	}
}

func TestRedactGoogleKey(t *testing.T) {
	input := "key AIzaSyD12345678901234567890123456789012"
	result := Redact(input)
	if strings.Contains(result, "12345678901234567890123456789012") {
		t.Error("Google API key was not redacted")
	}
}

func TestRedactBearerToken(t *testing.T) {
	input := "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.payload"
	result := Redact(input)
	if strings.Contains(result, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9") {
		t.Error("Bearer token was not redacted")
	}
}

func TestRedactGenericSecret(t *testing.T) {
	input := "password: mysecretpassword123"
	result := Redact(input)
	if strings.Contains(result, "mysecretpassword123") {
		t.Error("generic password was not redacted")
	}
}

func TestRedactNoSecrets(t *testing.T) {
	input := "This is a normal string with no secrets."
	result := Redact(input)
	if result != input {
		t.Errorf("normal string should not be modified, got: %s", result)
	}
}

func TestRedactMultipleSecrets(t *testing.T) {
	input := "key1=sk-abcdefghijklmnopqrstuvwxyz key2=sk-ant-api03-abcdefghijklmno"
	result := Redact(input)
	// Both should be partially masked
	if !strings.Contains(result, "****") || !strings.Contains(result, "sk-a") {
		t.Error("multiple secrets should be redacted with partial masking")
	}
}

func TestRedactPreservesPartialContent(t *testing.T) {
	input := "sk-abcdefghijklmnopqrst"
	result := Redact(input)
	// Should keep first 4 and last 4 chars
	if !strings.HasPrefix(result, "sk-a") {
		t.Error("should preserve first 4 chars")
	}
	if !strings.HasSuffix(result, "qrst") {
		t.Error("should preserve last 4 chars")
	}
	if !strings.Contains(result, "****") {
		t.Error("middle should be masked")
	}
}
