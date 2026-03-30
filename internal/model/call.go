package model

import "time"

// Provider represents a known LLM API provider.
type Provider string

const (
	ProviderOpenAI     Provider = "openai"
	ProviderAnthropic  Provider = "anthropic"
	ProviderGoogle     Provider = "google"
	ProviderOllama     Provider = "ollama"
	ProviderOpenRouter Provider = "openrouter"
	ProviderUnknown    Provider = "unknown"
)

// Call represents a single intercepted LLM API call.
type Call struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Provider  Provider  `json:"provider"`
	Model     string    `json:"model"`
	Endpoint  string    `json:"endpoint"`

	// Request
	RequestBody    string    `json:"request_body"`
	Messages       []Message `json:"messages,omitempty"`
	SystemPrompt   string    `json:"system_prompt,omitempty"`
	Tools          int       `json:"tools_count,omitempty"`

	// Response
	ResponseBody   string `json:"response_body"`
	ResponseText   string `json:"response_text,omitempty"`
	StatusCode     int    `json:"status_code"`
	IsStreaming     bool   `json:"is_streaming"`

	// Metrics
	InputTokens    int           `json:"input_tokens"`
	OutputTokens   int           `json:"output_tokens"`
	Latency        time.Duration `json:"latency_ms"`
	EstimatedCost  float64       `json:"estimated_cost_usd"`

	// Internal
	StartTime time.Time `json:"-"`
	Completed bool      `json:"-"`
}

// Message represents a chat message in the request.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Session holds all calls in the current aitap run.
type Session struct {
	StartTime  time.Time `json:"start_time"`
	Calls      []*Call   `json:"calls"`
	TotalCost  float64   `json:"total_cost_usd"`
	TotalIn    int       `json:"total_input_tokens"`
	TotalOut   int       `json:"total_output_tokens"`
}

func NewSession() *Session {
	return &Session{
		StartTime: time.Now(),
		Calls:     make([]*Call, 0),
	}
}

func (s *Session) Add(c *Call) {
	c.ID = len(s.Calls) + 1
	s.Calls = append(s.Calls, c)
	s.TotalCost += c.EstimatedCost
	s.TotalIn += c.InputTokens
	s.TotalOut += c.OutputTokens
}
