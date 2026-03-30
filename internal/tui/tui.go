package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aniketjoshi/aitap/internal/model"
)

// Styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#3D3D3D"))

	providerColors = map[model.Provider]lipgloss.Color{
		model.ProviderOpenAI:     lipgloss.Color("#10A37F"),
		model.ProviderAnthropic:  lipgloss.Color("#D97757"),
		model.ProviderGoogle:     lipgloss.Color("#4285F4"),
		model.ProviderOllama:     lipgloss.Color("#FFFFFF"),
		model.ProviderOpenRouter: lipgloss.Color("#B4A0FF"),
		model.ProviderUnknown:    lipgloss.Color("#888888"),
	}

	dimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	costStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#85E89D"))
	latStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#79B8FF"))
	roleSystem = lipgloss.NewStyle().Foreground(lipgloss.Color("#F97583")).Bold(true)
	roleUser   = lipgloss.NewStyle().Foreground(lipgloss.Color("#79B8FF")).Bold(true)
	roleAssist = lipgloss.NewStyle().Foreground(lipgloss.Color("#85E89D")).Bold(true)
	hintStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#B4A0FF"))
)

// NewCallMsg is sent when the proxy captures a new call.
type NewCallMsg struct{ Call *model.Call }

// Model is the Bubble Tea model for the TUI.
type Model struct {
	session  *model.Session
	cursor   int
	expanded bool
	width    int
	height   int
	callChan <-chan *model.Call
	port     int
}

// New creates a new TUI model.
func New(session *model.Session, callChan <-chan *model.Call, port int) Model {
	return Model{
		session:  session,
		callChan: callChan,
		port:     port,
	}
}

func (m Model) Init() tea.Cmd {
	return waitForCall(m.callChan)
}

func waitForCall(ch <-chan *model.Call) tea.Cmd {
	return func() tea.Msg {
		call := <-ch
		return NewCallMsg{Call: call}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.expanded = false
			}
		case "down", "j":
			if m.cursor < len(m.session.Calls)-1 {
				m.cursor++
				m.expanded = false
			}
		case "enter", " ":
			m.expanded = !m.expanded
		case "G":
			if len(m.session.Calls) > 0 {
				m.cursor = len(m.session.Calls) - 1
			}
		case "g":
			m.cursor = 0
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case NewCallMsg:
		m.session.Add(msg.Call)
		m.cursor = len(m.session.Calls) - 1
		return m, waitForCall(m.callChan)
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Starting aitap..."
	}

	var b strings.Builder

	// Header
	header := headerStyle.Render(fmt.Sprintf(
		" aitap  :%d  |  %d calls  |  %s tokens  |  %s ",
		m.port,
		len(m.session.Calls),
		formatTokens(m.session.TotalIn+m.session.TotalOut),
		formatCost(m.session.TotalCost),
	))
	b.WriteString(header + "\n\n")

	if len(m.session.Calls) == 0 {
		b.WriteString(dimStyle.Render("  Waiting for LLM API calls...\n\n"))
		b.WriteString(hintStyle.Render("  Forward proxy mode (recommended):\n"))
		b.WriteString(dimStyle.Render(fmt.Sprintf("    OPENAI_BASE_URL=http://localhost:%d/openai/v1\n", m.port)))
		b.WriteString(dimStyle.Render(fmt.Sprintf("    ANTHROPIC_BASE_URL=http://localhost:%d/anthropic\n", m.port)))
		b.WriteString(dimStyle.Render(fmt.Sprintf("    OLLAMA_HOST=http://localhost:%d/ollama\n\n", m.port)))
		b.WriteString(hintStyle.Render("  HTTP proxy mode:\n"))
		b.WriteString(dimStyle.Render(fmt.Sprintf("    HTTP_PROXY=http://127.0.0.1:%d\n", m.port)))
		return b.String()
	}

	// Call list — show as many as fit
	listHeight := m.height - 8 // reserve for header + detail
	if m.expanded {
		listHeight = min(len(m.session.Calls), 5)
	}

	start := 0
	if m.cursor >= listHeight {
		start = m.cursor - listHeight + 1
	}
	end := min(start+listHeight, len(m.session.Calls))

	for i := start; i < end; i++ {
		call := m.session.Calls[i]
		line := renderCallRow(call, m.width)
		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		b.WriteString(line + "\n")
	}

	// Detail pane
	if m.expanded && m.cursor < len(m.session.Calls) {
		b.WriteString("\n")
		b.WriteString(renderDetail(m.session.Calls[m.cursor], m.width))
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  j/k navigate  enter expand  q quit"))

	return b.String()
}

func renderCallRow(c *model.Call, width int) string {
	provColor := providerColors[c.Provider]
	prov := lipgloss.NewStyle().Foreground(provColor).Render(
		fmt.Sprintf("%-10s", c.Provider),
	)

	modelName := c.Model
	if len(modelName) > 22 {
		modelName = modelName[:19] + "..."
	}

	tokens := fmt.Sprintf("%s>%s", formatTokens(c.InputTokens), formatTokens(c.OutputTokens))
	cost := formatCost(c.EstimatedCost)
	lat := formatLatency(c.Latency)

	stream := " "
	if c.IsStreaming {
		stream = "~"
	}

	return fmt.Sprintf("  %s %3d | %s | %-22s | %10s | %7s | %6s",
		stream, c.ID, prov, modelName, tokens, costStyle.Render(cost), latStyle.Render(lat))
}

func renderDetail(c *model.Call, width int) string {
	var b strings.Builder

	detailWidth := min(width-4, 100)

	// Messages
	b.WriteString("  -- Request --\n")
	if c.SystemPrompt != "" {
		b.WriteString(fmt.Sprintf("  %s %s\n", roleSystem.Render("system:"), wrapText(c.SystemPrompt, detailWidth-12)))
	}
	for _, m := range c.Messages {
		style := roleUser
		if m.Role == "assistant" {
			style = roleAssist
		}
		b.WriteString(fmt.Sprintf("  %s %s\n", style.Render(m.Role+":"), wrapText(m.Content, detailWidth-12)))
	}

	if c.Tools > 0 {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  [%d tools attached]\n", c.Tools)))
	}

	// Response
	b.WriteString("\n  -- Response --\n")
	if c.ResponseText != "" {
		b.WriteString(fmt.Sprintf("  %s %s\n", roleAssist.Render("assistant:"), wrapText(c.ResponseText, detailWidth-12)))
	}
	b.WriteString(fmt.Sprintf("  %s\n", dimStyle.Render(fmt.Sprintf(
		"status=%d  in=%d  out=%d  cost=%s  latency=%s",
		c.StatusCode, c.InputTokens, c.OutputTokens,
		formatCost(c.EstimatedCost), formatLatency(c.Latency),
	))))

	return b.String()
}

// --- Formatting helpers ---

func formatTokens(n int) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}

func formatCost(c float64) string {
	if c == 0 {
		return "free"
	}
	if c < 0.01 {
		return fmt.Sprintf("$%.4f", c)
	}
	return fmt.Sprintf("$%.3f", c)
}

func formatLatency(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

func wrapText(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	var lines []string
	for len(s) > maxWidth {
		cut := maxWidth
		if idx := strings.LastIndex(s[:cut], " "); idx > maxWidth/2 {
			cut = idx
		}
		lines = append(lines, s[:cut])
		s = s[cut:]
		s = strings.TrimLeft(s, " ")
	}
	if len(s) > 0 {
		lines = append(lines, s)
	}
	indent := "\n" + strings.Repeat(" ", 12)
	return strings.Join(lines, indent)
}
