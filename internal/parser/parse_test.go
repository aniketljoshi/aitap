package parser

import (
	"testing"

	"github.com/aniketjoshi/aitap/internal/model"
)

func TestParseOpenAIRequest(t *testing.T) {
	body := []byte(`{
		"model": "gpt-4o",
		"messages": [
			{"role": "system", "content": "You are helpful."},
			{"role": "user", "content": "Hello"}
		],
		"tools": [{"type": "function", "function": {"name": "get_weather"}}],
		"stream": true
	}`)

	call := &model.Call{}
	ParseRequest(model.ProviderOpenAI, body, call)

	if call.Model != "gpt-4o" {
		t.Errorf("expected model gpt-4o, got %s", call.Model)
	}
	if len(call.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(call.Messages))
	}
	if call.Messages[0].Role != "system" {
		t.Errorf("expected first message role system, got %s", call.Messages[0].Role)
	}
	if call.Tools != 1 {
		t.Errorf("expected 1 tool, got %d", call.Tools)
	}
	if !call.IsStreaming {
		t.Error("expected IsStreaming to be true")
	}
}

func TestParseOpenAIResponse(t *testing.T) {
	body := []byte(`{
		"choices": [
			{"message": {"content": "Hello! How can I help you?"}}
		],
		"usage": {
			"prompt_tokens": 25,
			"completion_tokens": 8
		}
	}`)

	call := &model.Call{Model: "gpt-4o"}
	ParseResponse(model.ProviderOpenAI, body, call)

	if call.InputTokens != 25 {
		t.Errorf("expected 25 input tokens, got %d", call.InputTokens)
	}
	if call.OutputTokens != 8 {
		t.Errorf("expected 8 output tokens, got %d", call.OutputTokens)
	}
	if call.ResponseText != "Hello! How can I help you?" {
		t.Errorf("unexpected response text: %s", call.ResponseText)
	}
}

func TestParseAnthropicRequest(t *testing.T) {
	body := []byte(`{
		"model": "claude-sonnet-4-20250514",
		"system": "You are a helpful assistant.",
		"messages": [
			{"role": "user", "content": "What is Go?"}
		],
		"stream": false
	}`)

	call := &model.Call{}
	ParseRequest(model.ProviderAnthropic, body, call)

	if call.Model != "claude-sonnet-4-20250514" {
		t.Errorf("expected claude-sonnet-4-20250514, got %s", call.Model)
	}
	if call.SystemPrompt != "You are a helpful assistant." {
		t.Errorf("unexpected system prompt: %s", call.SystemPrompt)
	}
	if len(call.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(call.Messages))
	}
}

func TestParseAnthropicResponse(t *testing.T) {
	body := []byte(`{
		"content": [
			{"type": "text", "text": "Go is a programming language."}
		],
		"usage": {
			"input_tokens": 15,
			"output_tokens": 12
		}
	}`)

	call := &model.Call{Model: "claude-sonnet-4-20250514"}
	ParseResponse(model.ProviderAnthropic, body, call)

	if call.InputTokens != 15 {
		t.Errorf("expected 15 input tokens, got %d", call.InputTokens)
	}
	if call.OutputTokens != 12 {
		t.Errorf("expected 12 output tokens, got %d", call.OutputTokens)
	}
	if call.ResponseText != "Go is a programming language." {
		t.Errorf("unexpected response text: %s", call.ResponseText)
	}
}

func TestParseAnthropicContentBlocks(t *testing.T) {
	// Test Anthropic content as array of blocks (not just string)
	body := []byte(`{
		"model": "claude-sonnet-4-20250514",
		"messages": [
			{"role": "user", "content": [{"type": "text", "text": "Hello from blocks"}]}
		]
	}`)

	call := &model.Call{}
	ParseRequest(model.ProviderAnthropic, body, call)

	if len(call.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(call.Messages))
	}
	if call.Messages[0].Content != "Hello from blocks" {
		t.Errorf("expected 'Hello from blocks', got '%s'", call.Messages[0].Content)
	}
}

func TestParseOllamaRequest(t *testing.T) {
	body := []byte(`{
		"model": "llama3",
		"messages": [
			{"role": "user", "content": "Explain Go."}
		],
		"stream": false
	}`)

	call := &model.Call{}
	ParseRequest(model.ProviderOllama, body, call)

	if call.Model != "llama3" {
		t.Errorf("expected llama3, got %s", call.Model)
	}
	if call.IsStreaming {
		t.Error("expected IsStreaming to be false")
	}
}

func TestParseOllamaResponse(t *testing.T) {
	body := []byte(`{
		"message": {"content": "Go is great."},
		"prompt_eval_count": 10,
		"eval_count": 50
	}`)

	call := &model.Call{Model: "llama3"}
	ParseResponse(model.ProviderOllama, body, call)

	if call.InputTokens != 10 {
		t.Errorf("expected 10 input tokens, got %d", call.InputTokens)
	}
	if call.OutputTokens != 50 {
		t.Errorf("expected 50 output tokens, got %d", call.OutputTokens)
	}
}

func TestParseGoogleRequest(t *testing.T) {
	body := []byte(`{
		"contents": [
			{"role": "user", "parts": [{"text": "What is Go?"}]}
		]
	}`)

	call := &model.Call{}
	ParseRequest(model.ProviderGoogle, body, call)

	if len(call.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(call.Messages))
	}
	if call.Messages[0].Content != "What is Go?" {
		t.Errorf("unexpected content: %s", call.Messages[0].Content)
	}
}

func TestParseGoogleResponse(t *testing.T) {
	body := []byte(`{
		"candidates": [
			{"content": {"parts": [{"text": "Go is a language by Google."}]}}
		],
		"usageMetadata": {
			"promptTokenCount": 5,
			"candidatesTokenCount": 10
		}
	}`)

	call := &model.Call{Model: "gemini-2.5-pro"}
	ParseResponse(model.ProviderGoogle, body, call)

	if call.InputTokens != 5 {
		t.Errorf("expected 5 input tokens, got %d", call.InputTokens)
	}
	if call.OutputTokens != 10 {
		t.Errorf("expected 10 output tokens, got %d", call.OutputTokens)
	}
}

func TestTruncate(t *testing.T) {
	short := "hello"
	if truncate(short, 10) != "hello" {
		t.Error("short string should not be truncated")
	}

	long := "This is a very long string that exceeds the limit"
	result := truncate(long, 20)
	if len(result) != 23 { // 20 + "..."
		t.Errorf("expected truncated length 23, got %d", len(result))
	}
	if result[len(result)-3:] != "..." {
		t.Error("truncated string should end with ...")
	}
}

func TestParseInvalidJSON(t *testing.T) {
	// Should not panic on invalid JSON
	call := &model.Call{}
	ParseRequest(model.ProviderOpenAI, []byte("not json"), call)
	if call.Model != "" {
		t.Error("model should be empty for invalid JSON")
	}

	ParseResponse(model.ProviderOpenAI, []byte("not json"), call)
	if call.InputTokens != 0 {
		t.Error("tokens should be 0 for invalid JSON")
	}
}
