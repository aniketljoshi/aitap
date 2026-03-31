package tui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aniketljoshi/aitap/internal/model"
)

func TestViewEmptyState(t *testing.T) {
	m := New(model.NewSession(), make(chan *model.Call), 9119)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 140, Height: 42})
	m = updated.(Model)

	view := m.View()

	for _, want := range []string{
		"Ready For Live Traffic",
		"OPENAI_BASE_URL=http://localhost:9119/openai/v1",
		"HTTP_PROXY=http://127.0.0.1:9119",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected empty-state view to contain %q", want)
		}
	}
}

func TestViewWithCapturedCall(t *testing.T) {
	m := New(model.NewSession(), make(chan *model.Call), 9119)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 148, Height: 44})
	m = updated.(Model)

	call := &model.Call{
		Provider:      model.ProviderOpenAI,
		Model:         "gpt-4.1",
		Endpoint:      "/v1/chat/completions",
		Messages:      []model.Message{{Role: "user", Content: "Summarize the release notes and pull requests."}},
		ResponseText:  "Here are the notable changes and what to watch next.",
		StatusCode:    200,
		InputTokens:   1240,
		OutputTokens:  410,
		Latency:       1200 * time.Millisecond,
		EstimatedCost: 0.0064,
		Completed:     true,
	}

	updated, _ = m.Update(NewCallMsg{Call: call})
	m = updated.(Model)

	view := m.View()

	for _, want := range []string{
		"Session Flow",
		"gpt-4.1",
		"REQUEST",
		"Here are the notable changes",
		"STATUS 200",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected populated view to contain %q", want)
		}
	}
}
