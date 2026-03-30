package parser

import (
	"testing"

	"github.com/aniketjoshi/aitap/internal/model"
)

func TestParseOpenAISSE(t *testing.T) {
	chunks := []string{
		`{"id":"chatcmpl-123","model":"gpt-4o","choices":[{"delta":{"role":"assistant","content":""}}]}`,
		`{"id":"chatcmpl-123","model":"gpt-4o","choices":[{"delta":{"content":"Hello"}}]}`,
		`{"id":"chatcmpl-123","model":"gpt-4o","choices":[{"delta":{"content":" world"}}]}`,
		`{"id":"chatcmpl-123","model":"gpt-4o","choices":[{"delta":{}}],"usage":{"prompt_tokens":10,"completion_tokens":5}}`,
	}

	call := &model.Call{Provider: model.ProviderOpenAI}
	ParseSSEChunks(model.ProviderOpenAI, chunks, call)

	if call.Model != "gpt-4o" {
		t.Errorf("expected model gpt-4o, got %s", call.Model)
	}
	if call.ResponseText != "Hello world" {
		t.Errorf("expected 'Hello world', got '%s'", call.ResponseText)
	}
	if call.InputTokens != 10 {
		t.Errorf("expected 10 input tokens, got %d", call.InputTokens)
	}
	if call.OutputTokens != 5 {
		t.Errorf("expected 5 output tokens, got %d", call.OutputTokens)
	}
}

func TestParseOpenAISSEWithoutUsage(t *testing.T) {
	// When stream_options.include_usage is not set, usage is absent
	chunks := []string{
		`{"model":"gpt-4o-mini","choices":[{"delta":{"content":"Hi there!"}}]}`,
	}

	call := &model.Call{Provider: model.ProviderOpenAI}
	ParseSSEChunks(model.ProviderOpenAI, chunks, call)

	if call.ResponseText != "Hi there!" {
		t.Errorf("expected 'Hi there!', got '%s'", call.ResponseText)
	}
	// Should estimate tokens from text
	if call.OutputTokens == 0 {
		t.Error("expected estimated output tokens > 0")
	}
}

func TestParseAnthropicSSE(t *testing.T) {
	chunks := []string{
		`{"type":"message_start","message":{"model":"claude-sonnet-4-20250514","usage":{"input_tokens":20}}}`,
		`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
		`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}`,
		`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" from Claude"}}`,
		`{"type":"content_block_stop","index":0}`,
		`{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":8}}`,
		`{"type":"message_stop"}`,
	}

	call := &model.Call{Provider: model.ProviderAnthropic}
	ParseSSEChunks(model.ProviderAnthropic, chunks, call)

	if call.Model != "claude-sonnet-4-20250514" {
		t.Errorf("expected claude-sonnet-4-20250514, got %s", call.Model)
	}
	if call.ResponseText != "Hello from Claude" {
		t.Errorf("expected 'Hello from Claude', got '%s'", call.ResponseText)
	}
	if call.InputTokens != 20 {
		t.Errorf("expected 20 input tokens, got %d", call.InputTokens)
	}
	if call.OutputTokens != 8 {
		t.Errorf("expected 8 output tokens, got %d", call.OutputTokens)
	}
}

func TestParseOllamaSSE(t *testing.T) {
	chunks := []string{
		`{"model":"llama3","message":{"role":"assistant","content":"Hello"},"done":false}`,
		`{"model":"llama3","message":{"role":"assistant","content":" there"},"done":false}`,
		`{"model":"llama3","message":{"role":"assistant","content":""},"done":true,"prompt_eval_count":15,"eval_count":30}`,
	}

	call := &model.Call{Provider: model.ProviderOllama}
	ParseSSEChunks(model.ProviderOllama, chunks, call)

	if call.Model != "llama3" {
		t.Errorf("expected llama3, got %s", call.Model)
	}
	if call.ResponseText != "Hello there" {
		t.Errorf("expected 'Hello there', got '%s'", call.ResponseText)
	}
	if call.InputTokens != 15 {
		t.Errorf("expected 15 input tokens, got %d", call.InputTokens)
	}
	if call.OutputTokens != 30 {
		t.Errorf("expected 30 output tokens, got %d", call.OutputTokens)
	}
}

func TestParseGoogleSSE(t *testing.T) {
	chunks := []string{
		`{"candidates":[{"content":{"parts":[{"text":"Hello "}]}}],"usageMetadata":{"promptTokenCount":5,"candidatesTokenCount":2}}`,
		`{"candidates":[{"content":{"parts":[{"text":"world"}]}}],"usageMetadata":{"promptTokenCount":5,"candidatesTokenCount":4}}`,
	}

	call := &model.Call{Provider: model.ProviderGoogle}
	ParseSSEChunks(model.ProviderGoogle, chunks, call)

	if call.ResponseText != "Hello world" {
		t.Errorf("expected 'Hello world', got '%s'", call.ResponseText)
	}
	// Google sends cumulative — last one wins
	if call.InputTokens != 5 {
		t.Errorf("expected 5 input tokens, got %d", call.InputTokens)
	}
	if call.OutputTokens != 4 {
		t.Errorf("expected 4 output tokens, got %d", call.OutputTokens)
	}
}

func TestEstimateTokens(t *testing.T) {
	if estimateTokens("") != 0 {
		t.Error("empty string should have 0 tokens")
	}
	if estimateTokens("Hi") != 1 {
		t.Error("very short string should have at least 1 token")
	}
	tokens := estimateTokens("This is a longer string for token estimation")
	if tokens < 5 || tokens > 20 {
		t.Errorf("unexpected token estimate: %d", tokens)
	}
}

func TestParseSSEEmptyChunks(t *testing.T) {
	call := &model.Call{Provider: model.ProviderOpenAI}
	ParseSSEChunks(model.ProviderOpenAI, nil, call)

	if call.ResponseText != "" {
		t.Error("expected empty response text for nil chunks")
	}
}

func TestParseSSEInvalidJSON(t *testing.T) {
	chunks := []string{"not json at all", "{malformed"}
	call := &model.Call{Provider: model.ProviderOpenAI}
	ParseSSEChunks(model.ProviderOpenAI, chunks, call)

	// Should not panic and should have empty results
	if call.Model != "" {
		t.Error("expected empty model for invalid JSON chunks")
	}
}
