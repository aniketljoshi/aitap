package parser

import (
	"encoding/json"
	"strings"

	"github.com/aniketljoshi/aitap/internal/model"
)

// ParseSSEChunks processes accumulated SSE data chunks from a streaming response
// and extracts the model name, token counts, and response text.
func ParseSSEChunks(p model.Provider, chunks []string, call *model.Call) {
	switch p {
	case model.ProviderOpenAI, model.ProviderOpenRouter:
		parseOpenAISSE(chunks, call)
	case model.ProviderAnthropic:
		parseAnthropicSSE(chunks, call)
	case model.ProviderOllama:
		parseOllamaSSE(chunks, call)
	case model.ProviderGoogle:
		parseGoogleSSE(chunks, call)
	}
}

// --- OpenAI SSE format ---
// Each chunk: {"id":"...","object":"chat.completion.chunk","model":"gpt-4o","choices":[{"delta":{"content":"Hello"}}]}
// Final chunk may include usage: {"usage":{"prompt_tokens":10,"completion_tokens":20,"total_tokens":30}}

func parseOpenAISSE(chunks []string, call *model.Call) {
	var textParts []string

	for _, chunk := range chunks {
		var data struct {
			Model   string `json:"model"`
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
					Role    string `json:"role"`
				} `json:"delta"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
			Usage *struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
			} `json:"usage"`
		}
		if err := json.Unmarshal([]byte(chunk), &data); err != nil {
			continue
		}

		// Capture model from first chunk
		if call.Model == "" && data.Model != "" {
			call.Model = data.Model
		}

		// Accumulate content deltas
		for _, choice := range data.Choices {
			if choice.Delta.Content != "" {
				textParts = append(textParts, choice.Delta.Content)
			}
		}

		// Capture usage from final chunk (OpenAI includes it when stream_options.include_usage=true)
		if data.Usage != nil {
			call.InputTokens = data.Usage.PromptTokens
			call.OutputTokens = data.Usage.CompletionTokens
		}
	}

	call.ResponseText = truncate(strings.Join(textParts, ""), 500)

	// Estimate tokens from text length if usage wasn't included in stream
	if call.OutputTokens == 0 && len(textParts) > 0 {
		fullText := strings.Join(textParts, "")
		call.OutputTokens = estimateTokens(fullText)
	}
}

// --- Anthropic SSE format ---
// Events: message_start (model, usage.input_tokens), content_block_delta (text delta),
// message_delta (usage.output_tokens), message_stop

func parseAnthropicSSE(chunks []string, call *model.Call) {
	var textParts []string

	for _, chunk := range chunks {
		// Try message_start (contains model and input tokens)
		var msgStart struct {
			Type    string `json:"type"`
			Message struct {
				Model string `json:"model"`
				Usage struct {
					InputTokens int `json:"input_tokens"`
				} `json:"usage"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(chunk), &msgStart); err == nil && msgStart.Type == "message_start" {
			call.Model = msgStart.Message.Model
			call.InputTokens = msgStart.Message.Usage.InputTokens
			continue
		}

		// Try content_block_delta (text chunks)
		var delta struct {
			Type  string `json:"type"`
			Delta struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"delta"`
		}
		if err := json.Unmarshal([]byte(chunk), &delta); err == nil && delta.Type == "content_block_delta" {
			if delta.Delta.Text != "" {
				textParts = append(textParts, delta.Delta.Text)
			}
			continue
		}

		// Try message_delta (output tokens at the end)
		var msgDelta struct {
			Type  string `json:"type"`
			Usage struct {
				OutputTokens int `json:"output_tokens"`
			} `json:"usage"`
		}
		if err := json.Unmarshal([]byte(chunk), &msgDelta); err == nil && msgDelta.Type == "message_delta" {
			if msgDelta.Usage.OutputTokens > 0 {
				call.OutputTokens = msgDelta.Usage.OutputTokens
			}
			continue
		}
	}

	call.ResponseText = truncate(strings.Join(textParts, ""), 500)

	if call.OutputTokens == 0 && len(textParts) > 0 {
		fullText := strings.Join(textParts, "")
		call.OutputTokens = estimateTokens(fullText)
	}
}

// --- Ollama SSE format ---
// Each chunk: {"model":"llama3","message":{"role":"assistant","content":"Hi"},"done":false}
// Final chunk: {"model":"llama3","done":true,"prompt_eval_count":26,"eval_count":283}

func parseOllamaSSE(chunks []string, call *model.Call) {
	var textParts []string

	for _, chunk := range chunks {
		var data struct {
			Model   string `json:"model"`
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Done            bool `json:"done"`
			PromptEvalCount int  `json:"prompt_eval_count"`
			EvalCount       int  `json:"eval_count"`
		}
		if err := json.Unmarshal([]byte(chunk), &data); err != nil {
			continue
		}

		if call.Model == "" && data.Model != "" {
			call.Model = data.Model
		}

		if data.Message.Content != "" {
			textParts = append(textParts, data.Message.Content)
		}

		if data.Done {
			call.InputTokens = data.PromptEvalCount
			call.OutputTokens = data.EvalCount
		}
	}

	call.ResponseText = truncate(strings.Join(textParts, ""), 500)
}

// --- Google SSE format ---
// Each chunk: {"candidates":[{"content":{"parts":[{"text":"Hello"}]}}],"usageMetadata":{...}}

func parseGoogleSSE(chunks []string, call *model.Call) {
	var textParts []string

	for _, chunk := range chunks {
		var data struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
			UsageMetadata struct {
				PromptTokenCount     int `json:"promptTokenCount"`
				CandidatesTokenCount int `json:"candidatesTokenCount"`
			} `json:"usageMetadata"`
		}
		if err := json.Unmarshal([]byte(chunk), &data); err != nil {
			continue
		}

		for _, cand := range data.Candidates {
			for _, part := range cand.Content.Parts {
				if part.Text != "" {
					textParts = append(textParts, part.Text)
				}
			}
		}

		// Google sends cumulative usage in each chunk — last one wins
		if data.UsageMetadata.PromptTokenCount > 0 {
			call.InputTokens = data.UsageMetadata.PromptTokenCount
		}
		if data.UsageMetadata.CandidatesTokenCount > 0 {
			call.OutputTokens = data.UsageMetadata.CandidatesTokenCount
		}
	}

	if call.Model == "" {
		call.Model = "gemini"
	}

	call.ResponseText = truncate(strings.Join(textParts, ""), 500)
}

// estimateTokens gives a rough token count based on character length.
// ~4 chars per token is a common English approximation.
func estimateTokens(text string) int {
	if len(text) == 0 {
		return 0
	}
	tokens := len(text) / 4
	if tokens == 0 {
		tokens = 1
	}
	return tokens
}
