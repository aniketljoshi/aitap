package provider

import (
	"testing"

	"github.com/aniketljoshi/aitap/internal/model"
)

func TestDetectProvider(t *testing.T) {
	tests := []struct {
		host     string
		expected model.Provider
	}{
		{"api.openai.com", model.ProviderOpenAI},
		{"api.openai.com:443", model.ProviderOpenAI},
		{"api.anthropic.com", model.ProviderAnthropic},
		{"api.anthropic.com:443", model.ProviderAnthropic},
		{"generativelanguage.googleapis.com", model.ProviderGoogle},
		{"openrouter.ai", model.ProviderOpenRouter},
		{"localhost:11434", model.ProviderOllama},
		{"127.0.0.1:11434", model.ProviderOllama},
		{"example.com", model.ProviderUnknown},
		{"", model.ProviderUnknown},
	}

	for _, tt := range tests {
		got := DetectProvider(tt.host)
		if got != tt.expected {
			t.Errorf("DetectProvider(%q) = %s, want %s", tt.host, got, tt.expected)
		}
	}
}

func TestEstimateCostExactMatch(t *testing.T) {
	// gpt-4o: $2.50/1M input, $10.00/1M output
	cost := EstimateCost("gpt-4o", 1000, 500)
	expected := (1000.0 * 2.50 / 1_000_000) + (500.0 * 10.00 / 1_000_000)
	if cost != expected {
		t.Errorf("expected cost %f, got %f", expected, cost)
	}
}

func TestEstimateCostUnknownModel(t *testing.T) {
	cost := EstimateCost("unknown-model-xyz", 1000, 500)
	if cost != 0 {
		t.Errorf("expected 0 cost for unknown model, got %f", cost)
	}
}

func TestEstimateCostZeroTokens(t *testing.T) {
	cost := EstimateCost("gpt-4o", 0, 0)
	if cost != 0 {
		t.Errorf("expected 0 cost for zero tokens, got %f", cost)
	}
}

func TestEstimateCostAnthropicModel(t *testing.T) {
	// claude-sonnet-4-20250514: $3.00/1M input, $15.00/1M output
	cost := EstimateCost("claude-sonnet-4-20250514", 10000, 5000)
	expected := (10000.0 * 3.00 / 1_000_000) + (5000.0 * 15.00 / 1_000_000)
	if cost != expected {
		t.Errorf("expected cost %f, got %f", expected, cost)
	}
}

func TestPricingTableCompleteness(t *testing.T) {
	// Ensure all pricing entries have non-zero values
	for model, pricing := range Pricing {
		if pricing.InputPer1M <= 0 {
			t.Errorf("model %s has zero/negative input pricing", model)
		}
		if pricing.OutputPer1M <= 0 {
			t.Errorf("model %s has zero/negative output pricing", model)
		}
	}
}
