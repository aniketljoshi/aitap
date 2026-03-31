package main

import (
	"strings"
	"time"

	"github.com/aniketljoshi/aitap/internal/model"
)

func startDemoFeed(callChan chan<- *model.Call, filterProvider string) {
	go func() {
		base := time.Now().Add(-2 * time.Minute)
		for i, call := range demoCalls() {
			call.Timestamp = base.Add(time.Duration(i) * 8 * time.Second)
			call.StartTime = call.Timestamp
			call.Completed = true

			if filterProvider != "" && string(call.Provider) != filterProvider {
				continue
			}

			callChan <- cloneCall(call)
			time.Sleep(450 * time.Millisecond)
		}
	}()
}

func demoCalls() []*model.Call {
	return []*model.Call{
		{
			Provider:      model.ProviderAnthropic,
			Model:         "claude-sonnet-4-20250514",
			Endpoint:      "/v1/messages",
			SystemPrompt:  "You are a precise release engineer who summarizes pull requests and risks.",
			Messages:      []model.Message{{Role: "user", Content: "Summarize the latest PR changes, mention any risks, and suggest the next verification steps."}},
			ResponseText:  "The latest changes improve parser coverage, add GitHub community docs, and redesign the terminal UI. The main risk is visual regressions on narrow terminals, so verify the layout at multiple widths and with streaming responses.",
			StatusCode:    200,
			InputTokens:   1820,
			OutputTokens:  710,
			Latency:       2300 * time.Millisecond,
			EstimatedCost: 0.0161,
			IsStreaming:   true,
			Tools:         2,
		},
		{
			Provider:      model.ProviderOpenAI,
			Model:         "gpt-4.1",
			Endpoint:      "/v1/chat/completions",
			Messages:      []model.Message{{Role: "user", Content: "Generate a release summary for this sprint and keep it concise."}},
			ResponseText:  "This sprint focused on proxy fidelity, export safety, and a more premium terminal experience for inspecting LLM traffic.",
			StatusCode:    200,
			InputTokens:   1240,
			OutputTokens:  420,
			Latency:       1200 * time.Millisecond,
			EstimatedCost: 0.0058,
		},
		{
			Provider:      model.ProviderGoogle,
			Model:         "gemini-2.5-pro",
			Endpoint:      "/v1beta/models/gemini-2.5-pro:generateContent",
			Messages:      []model.Message{{Role: "user", Content: "Compare all open issues and group them by root cause."}},
			ResponseText:  "Most issues cluster around parser edge cases, pricing drift, and developer-experience polish. The highest leverage work is automated regression coverage for streamed responses.",
			StatusCode:    200,
			InputTokens:   2100,
			OutputTokens:  980,
			Latency:       4700 * time.Millisecond,
			EstimatedCost: 0.0124,
			IsStreaming:   true,
		},
		{
			Provider:      model.ProviderOpenRouter,
			Model:         "openai/gpt-4.1-mini",
			Endpoint:      "/api/v1/chat/completions",
			Messages:      []model.Message{{Role: "user", Content: "Draft the GitHub description, topics, and release notes for this repo."}},
			ResponseText:  "Use a short description, precise topics, and a terminal-first positioning statement that emphasizes live visibility into LLM traffic.",
			StatusCode:    200,
			InputTokens:   890,
			OutputTokens:  350,
			Latency:       1600 * time.Millisecond,
			EstimatedCost: 0.0021,
		},
		{
			Provider:      model.ProviderOllama,
			Model:         "llama3:8b",
			Endpoint:      "/api/chat",
			Messages:      []model.Message{{Role: "user", Content: "Say hello from local demo mode."}},
			ResponseText:  "Hello from aitap demo mode. Your local TUI is rendering sample traffic without calling any external provider.",
			StatusCode:    200,
			InputTokens:   220,
			OutputTokens:  96,
			Latency:       860 * time.Millisecond,
			EstimatedCost: 0,
		},
	}
}

func cloneCall(src *model.Call) *model.Call {
	dst := *src
	if len(src.Messages) > 0 {
		dst.Messages = append([]model.Message(nil), src.Messages...)
	}
	dst.RequestBody = strings.Clone(src.RequestBody)
	dst.ResponseBody = strings.Clone(src.ResponseBody)
	dst.ResponseText = strings.Clone(src.ResponseText)
	dst.SystemPrompt = strings.Clone(src.SystemPrompt)
	dst.Model = strings.Clone(src.Model)
	dst.Endpoint = strings.Clone(src.Endpoint)
	return &dst
}
