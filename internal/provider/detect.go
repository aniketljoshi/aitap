package provider

import (
	"strings"

	"github.com/aniketljoshi/aitap/internal/model"
)

// DetectProvider identifies the LLM provider from the request host.
func DetectProvider(host string) model.Provider {
	h := strings.ToLower(host)
	switch {
	case strings.Contains(h, "api.openai.com"):
		return model.ProviderOpenAI
	case strings.Contains(h, "api.anthropic.com"):
		return model.ProviderAnthropic
	case strings.Contains(h, "generativelanguage.googleapis.com"):
		return model.ProviderGoogle
	case strings.Contains(h, "openrouter.ai"):
		return model.ProviderOpenRouter
	case strings.Contains(h, "localhost:11434"), strings.Contains(h, "127.0.0.1:11434"):
		return model.ProviderOllama
	default:
		return model.ProviderUnknown
	}
}

// Price per 1M tokens (input, output) in USD.
// Updated March 2026. Keep this file current.
type ModelPricing struct {
	InputPer1M  float64
	OutputPer1M float64
}

var Pricing = map[string]ModelPricing{
	// OpenAI
	"gpt-4o":          {2.50, 10.00},
	"gpt-4o-mini":     {0.15, 0.60},
	"gpt-4.1":         {2.00, 8.00},
	"gpt-4.1-mini":    {0.40, 1.60},
	"gpt-4.1-nano":    {0.10, 0.40},
	"o3":              {2.00, 8.00},
	"o3-mini":         {1.10, 4.40},
	"o4-mini":         {1.10, 4.40},

	// Anthropic
	"claude-sonnet-4-20250514":    {3.00, 15.00},
	"claude-opus-4-20250514":     {15.00, 75.00},
	"claude-haiku-4-5-20251001":  {0.80, 4.00},

	// Google
	"gemini-2.5-pro":   {1.25, 10.00},
	"gemini-2.5-flash": {0.15, 0.60},
}

// EstimateCost calculates the cost for a call based on model and token counts.
func EstimateCost(modelName string, inputTokens, outputTokens int) float64 {
	// Try exact match first
	if p, ok := Pricing[modelName]; ok {
		return (float64(inputTokens) * p.InputPer1M / 1_000_000) +
			(float64(outputTokens) * p.OutputPer1M / 1_000_000)
	}

	// Try prefix match (e.g., "gpt-4o-2025-01-01" matches "gpt-4o")
	for name, p := range Pricing {
		if strings.HasPrefix(modelName, name) {
			return (float64(inputTokens) * p.InputPer1M / 1_000_000) +
				(float64(outputTokens) * p.OutputPer1M / 1_000_000)
		}
	}

	return 0 // Unknown model
}
