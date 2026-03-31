package tui

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aniketjoshi/aitap/internal/model"
)

type providerTheme struct {
	accent   lipgloss.Color
	border   lipgloss.Color
	softBg   lipgloss.Color
	softText lipgloss.Color
}

var (
	appBg             = lipgloss.Color("#05070A")
	surfaceBg         = lipgloss.Color("#090C10")
	surfaceBgAlt      = lipgloss.Color("#0D1117")
	panelBorder       = lipgloss.Color("#1A2530")
	panelBorderStrong = lipgloss.Color("#27415B")
	textPrimary       = lipgloss.Color("#F2F5F7")
	textSecondary     = lipgloss.Color("#B6C2CF")
	textMuted         = lipgloss.Color("#6F8194")
	brandMint         = lipgloss.Color("#66E3B3")
	brandCyan         = lipgloss.Color("#59C8F7")
	brandAmber        = lipgloss.Color("#FFB35C")
	brandSlate        = lipgloss.Color("#AAB8C9")
	brandViolet       = lipgloss.Color("#A78BFA")
	successColor      = lipgloss.Color("#90F0C7")
	warningColor      = lipgloss.Color("#FFD089")
	dangerColor       = lipgloss.Color("#FF8C82")
)

var providerThemes = map[model.Provider]providerTheme{
	model.ProviderOpenAI: {
		accent:   lipgloss.Color("#66E3B3"),
		border:   lipgloss.Color("#1C4C4A"),
		softBg:   lipgloss.Color("#12302C"),
		softText: lipgloss.Color("#A5F1D0"),
	},
	model.ProviderAnthropic: {
		accent:   lipgloss.Color("#FFB35C"),
		border:   lipgloss.Color("#5A3B1D"),
		softBg:   lipgloss.Color("#332111"),
		softText: lipgloss.Color("#FFD8A0"),
	},
	model.ProviderGoogle: {
		accent:   lipgloss.Color("#59C8F7"),
		border:   lipgloss.Color("#1A4260"),
		softBg:   lipgloss.Color("#10283A"),
		softText: lipgloss.Color("#A3E4FF"),
	},
	model.ProviderOllama: {
		accent:   lipgloss.Color("#C7D4E2"),
		border:   lipgloss.Color("#324455"),
		softBg:   lipgloss.Color("#1B2836"),
		softText: lipgloss.Color("#E6EEF8"),
	},
	model.ProviderOpenRouter: {
		accent:   lipgloss.Color("#C3A8FF"),
		border:   lipgloss.Color("#41306A"),
		softBg:   lipgloss.Color("#1A1832"),
		softText: lipgloss.Color("#E3D7FF"),
	},
	model.ProviderUnknown: {
		accent:   lipgloss.Color("#93A4B5"),
		border:   lipgloss.Color("#2C4154"),
		softBg:   lipgloss.Color("#182532"),
		softText: lipgloss.Color("#D5E0EA"),
	},
}

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
			}
		case "down", "j":
			if m.cursor < len(m.session.Calls)-1 {
				m.cursor++
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
	if m.width == 0 || m.height == 0 {
		return "Starting aitap..."
	}

	frameWidth := max(60, m.width-2)

	header := m.renderHeader(frameWidth)
	bodyHeight := max(14, m.height-lipgloss.Height(header)-5)

	var body string
	if len(m.session.Calls) == 0 {
		body = m.renderEmptyState(frameWidth, bodyHeight)
	} else {
		body = m.renderDashboard(frameWidth, bodyHeight)
	}

	footer := m.renderFooter(frameWidth)

	frame := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		footer,
	)

	return lipgloss.NewStyle().
		Padding(1, 1).
		Background(appBg).
		Render(frame)
}

func (m Model) renderHeader(width int) string {
	brand := lipgloss.NewStyle().
		Bold(true).
		Foreground(textPrimary).
		Background(lipgloss.Color("#0E2434")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorderStrong).
		Padding(0, 1).
		Render("aitap")

	subtitle := lipgloss.NewStyle().
		Foreground(textSecondary).
		Render("local-first llm traffic inspector")

	left := lipgloss.JoinVertical(lipgloss.Left, brand, subtitle)

	stats := []string{
		m.renderHeaderPill(fmt.Sprintf(":%d", m.port), brandCyan, lipgloss.Color("#0F2234")),
		m.renderHeaderPill(fmt.Sprintf("%d calls", len(m.session.Calls)), brandMint, lipgloss.Color("#10281F")),
		m.renderHeaderPill(fmt.Sprintf("%s tok", formatTokens(m.session.TotalIn+m.session.TotalOut)), warningColor, lipgloss.Color("#2E2210")),
		m.renderHeaderPill(formatCost(m.session.TotalCost), successColor, lipgloss.Color("#133129")),
	}
	if m.currentFilter() != "" {
		stats = append(stats, m.renderHeaderPill("filter "+m.currentFilter(), textPrimary, lipgloss.Color("#1A2136")))
	}
	right := lipgloss.JoinHorizontal(lipgloss.Left, stats...)

	top := fillHorizontal(left, right, width)

	hint := lipgloss.NewStyle().
		Foreground(textMuted).
		Render("j/k navigate  enter inspect  g/G jump  q quit")

	return lipgloss.JoinVertical(lipgloss.Left, top, hint)
}

func (m Model) renderHeaderPill(label string, fg, bg lipgloss.Color) string {
	return lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Padding(0, 1).
		MarginLeft(1).
		Bold(true).
		Render(label)
}

func (m Model) renderDashboard(width, height int) string {
	if width < 110 {
		stack := []string{
			m.renderCallsPanel(width, max(12, height/2)),
			m.renderInspectorPanel(width, max(14, height/2)),
		}
		return lipgloss.JoinVertical(lipgloss.Left, stack...)
	}

	listWidth := max(42, width*44/100)
	detailWidth := max(52, width-listWidth-2)

	left := m.renderCallsPanel(listWidth, height)
	right := m.renderInspectorPanel(detailWidth, height)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
}

func (m Model) renderCallsPanel(width, height int) string {
	panelInnerWidth := boxedInnerWidth(width)

	headerLeft := lipgloss.NewStyle().Bold(true).Foreground(textPrimary).Render("Session Flow")
	headerRight := lipgloss.NewStyle().Foreground(textMuted).Render(fmt.Sprintf("%d captured", len(m.session.Calls)))
	panelHeader := fillHorizontal(headerLeft, headerRight, panelInnerWidth)

	available := max(1, height-7)
	cardHeight := 6
	itemsPerPage := max(1, available/cardHeight)
	start := 0
	if m.cursor >= itemsPerPage {
		start = m.cursor - itemsPerPage + 1
	}
	end := min(len(m.session.Calls), start+itemsPerPage)

	var cards []string
	for i := start; i < end; i++ {
		cards = append(cards, renderCallCard(m.session.Calls[i], panelInnerWidth, i == m.cursor))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{panelHeader}, cards...)...)

	return lipgloss.NewStyle().
		Width(panelInnerWidth).
		Height(max(boxedInnerHeight(height), lipgloss.Height(content))).
		Background(surfaceBg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorder).
		Padding(1).
		Render(content)
}

func renderCallCard(c *model.Call, width int, selected bool) string {
	contentWidth := boxedInnerWidth(width)
	theme := themeFor(c.Provider)

	cardBg := surfaceBg
	cardBorder := panelBorder
	if selected {
		cardBg = surfaceBgAlt
		cardBorder = theme.accent
	}

	idPill := lipgloss.NewStyle().
		Foreground(textMuted).
		Bold(true).
		Render(fmt.Sprintf("#%d", c.ID))

	providerPill := lipgloss.NewStyle().
		Foreground(theme.accent).
		Bold(true).
		Render(strings.ToUpper(string(c.Provider)))

	statusPill := renderStatusPill(c)
	top := lineWithBg(fillHorizontal(
		lipgloss.JoinHorizontal(lipgloss.Left, idPill, " ", providerPill),
		statusPill,
		contentWidth,
	), contentWidth, cardBg, textSecondary)

	modelLine := lipgloss.NewStyle().
		Bold(true).
		Foreground(textPrimary).
		Render(truncateRunes(fallbackModel(c.Model), contentWidth))
	modelLine = lineWithBg(modelLine, contentWidth, cardBg, textPrimary)

	preview := lipgloss.NewStyle().
		Foreground(textSecondary).
		Render(truncateRunes(firstPreview(c), contentWidth))
	preview = lineWithBg(preview, contentWidth, cardBg, textSecondary)

	metrics := lipgloss.JoinHorizontal(
		lipgloss.Left,
		renderInlineMetric("STATUS", fmt.Sprintf("%d", c.StatusCode), textSecondary),
		renderInlineMetric("TOK", fmt.Sprintf("%s>%s", formatTokens(c.InputTokens), formatTokens(c.OutputTokens)), warningColor),
		renderInlineMetric("LAT", formatLatency(c.Latency), brandCyan),
		renderInlineMetric("COST", formatCost(c.EstimatedCost), successColor),
	)
	metrics = lineWithBg(metrics, contentWidth, cardBg, textSecondary)

	card := lipgloss.JoinVertical(lipgloss.Left, top, modelLine, preview, metrics)

	return lipgloss.NewStyle().
		Width(contentWidth).
		Background(cardBg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cardBorder).
		Padding(0, 1).
		Render(card)
}

func (m Model) renderInspectorPanel(width, height int) string {
	call := m.session.Calls[m.cursor]
	theme := themeFor(call.Provider)
	panelInnerWidth := boxedInnerWidth(width)

	header := m.renderInspectorHeader(call, panelInnerWidth, theme)
	overview := m.renderOverviewBar(call, panelInnerWidth)

	request := renderSectionPanel(
		"Request",
		renderMessages(call, panelInnerWidth-4),
		panelInnerWidth,
	)

	response := renderSectionPanel(
		"Response",
		renderResponse(call, panelInnerWidth-4),
		panelInnerWidth,
	)

	var sections []string
	sections = append(sections, header, overview)
	if panelInnerWidth > 66 && m.expanded {
		leftWidth := (panelInnerWidth - 1) / 2
		rightWidth := panelInnerWidth - 1 - leftWidth
		left := renderSectionPanel("Request", renderMessages(call, boxedInnerWidth(leftWidth)), leftWidth)
		right := renderSectionPanel("Response", renderResponse(call, boxedInnerWidth(rightWidth)), rightWidth)
		sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right))
	} else {
		sections = append(sections, request, response)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return lipgloss.NewStyle().
		Width(panelInnerWidth).
		Height(max(boxedInnerHeight(height), lipgloss.Height(content))).
		Background(surfaceBg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorderStrong).
		Padding(1).
		Render(content)
}

func (m Model) renderInspectorHeader(c *model.Call, width int, theme providerTheme) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(textPrimary).
		Render(truncateRunes(fallbackModel(c.Model), max(18, width-16)))

	provider := lipgloss.NewStyle().
		Foreground(theme.accent).
		Bold(true).
		Render(strings.ToUpper(string(c.Provider)))

	top := lineWithBg(fillHorizontal(title, provider, width), width, surfaceBg, textPrimary)

	subtitle := lipgloss.NewStyle().
		Foreground(textMuted).
		Render(truncateRunes(endpointLabel(c.Endpoint), width))
	subtitle = lineWithBg(subtitle, width, surfaceBg, textMuted)

	return lipgloss.JoinVertical(lipgloss.Left, top, subtitle)
}

func (m Model) renderOverviewBar(c *model.Call, width int) string {
	line1 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		renderInlineMetric("STATUS", fmt.Sprintf("%d", c.StatusCode), textPrimary),
		renderInlineMetric("STREAM", yesNo(c.IsStreaming), brandMint),
		renderInlineMetric("MSGS", fmt.Sprintf("%d", len(c.Messages)), brandCyan),
		renderInlineMetric("TOOLS", fmt.Sprintf("%d", c.Tools), warningColor),
	)
	line2 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		renderInlineMetric("INPUT", formatTokens(c.InputTokens), warningColor),
		renderInlineMetric("OUTPUT", formatTokens(c.OutputTokens), brandCyan),
		renderInlineMetric("LATENCY", formatLatency(c.Latency), textSecondary),
		renderInlineMetric("COST", formatCost(c.EstimatedCost), successColor),
	)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lineWithBg(line1, width, surfaceBg, textSecondary),
		lineWithBg(line2, width, surfaceBg, textSecondary),
	)
}

func (m Model) renderEmptyState(width, height int) string {
	if width < 96 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			renderSectionPanel(
				"aitap",
				"Point your SDK at localhost and inspect prompts, tokens, latency, and costs as they stream.",
				width-4,
			),
			renderSectionPanel(
				"Forward Proxy",
				fmt.Sprintf(
					"OPENAI_BASE_URL=http://localhost:%d/openai/v1\nANTHROPIC_BASE_URL=http://localhost:%d/anthropic\nGOOGLE_API_BASE=http://localhost:%d/google\nOPENROUTER_BASE_URL=http://localhost:%d/openrouter/api/v1\nOLLAMA_HOST=http://localhost:%d/ollama",
					m.port, m.port, m.port, m.port, m.port,
				),
				width-4,
			),
			renderSectionPanel(
				"HTTP Proxy",
				fmt.Sprintf("HTTP_PROXY=http://127.0.0.1:%d", m.port),
				width-4,
			),
		)
		return lipgloss.NewStyle().
			Width(boxedInnerWidth(width)).
			Height(boxedInnerHeight(height)).
			Background(surfaceBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(panelBorder).
			Padding(1).
			Render(content)
	}

	leftWidth := max(40, width*45/100)
	rightWidth := max(34, width-leftWidth-2)

	heroTitle := lipgloss.NewStyle().Bold(true).Foreground(textPrimary).Render("Ready For Live Traffic")
	heroCopy := lipgloss.NewStyle().
		Foreground(textSecondary).
		Render("aitap sits between your app and the model provider so you can inspect every request, every streamed chunk, and every cost signal without shipping data to a dashboard.")

	heroBadges := lipgloss.JoinHorizontal(
		lipgloss.Left,
		renderMetricPill("modes", "forward + http proxy", brandCyan, lipgloss.Color("#152A3C")),
		renderMetricPill("providers", "openai anthropic google ollama", brandMint, lipgloss.Color("#143126")),
	)

	left := lipgloss.NewStyle().
		Width(boxedInnerWidth(leftWidth)).
		Height(boxedInnerHeight(height)).
		Background(surfaceBg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorderStrong).
		Padding(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, heroTitle, heroCopy, heroBadges))

	rightContent := lipgloss.JoinVertical(
		lipgloss.Left,
		renderSectionPanel(
			"Forward Proxy",
			fmt.Sprintf(
				"OPENAI_BASE_URL=http://localhost:%d/openai/v1\nANTHROPIC_BASE_URL=http://localhost:%d/anthropic\nGOOGLE_API_BASE=http://localhost:%d/google\nOPENROUTER_BASE_URL=http://localhost:%d/openrouter/api/v1\nOLLAMA_HOST=http://localhost:%d/ollama",
				m.port, m.port, m.port, m.port, m.port,
			),
			rightWidth-4,
		),
		renderSectionPanel(
			"HTTP Proxy",
			fmt.Sprintf("HTTP_PROXY=http://127.0.0.1:%d", m.port),
			rightWidth-4,
		),
	)

	right := lipgloss.NewStyle().
		Width(boxedInnerWidth(rightWidth)).
		Height(boxedInnerHeight(height)).
		Background(surfaceBg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorder).
		Padding(1).
		Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
}

func (m Model) renderFooter(width int) string {
	left := lipgloss.NewStyle().
		Foreground(textMuted).
		Render("bubble tea tui  |  stream-aware parsing  |  local-first inspection")
	right := lipgloss.NewStyle().
		Foreground(brandMint).
		Render("waiting for next call")
	if len(m.session.Calls) > 0 {
		right = lipgloss.NewStyle().
			Foreground(brandMint).
			Render(fmt.Sprintf("selected #%d of %d", m.cursor+1, len(m.session.Calls)))
	}
	return fillHorizontal(left, right, width)
}

func renderSectionPanel(title, body string, width int) string {
	contentWidth := boxedInnerWidth(width)
	titleStyle := lipgloss.NewStyle().
		Foreground(textMuted).
		Bold(true).
		Render(strings.ToUpper(title))

	bodyStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Width(max(16, contentWidth)).
		Background(surfaceBg)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lineWithBg(titleStyle, contentWidth, surfaceBg, textMuted),
		bodyStyle.Render(body),
	)

	return lipgloss.NewStyle().
		Width(contentWidth).
		Background(surfaceBg).
		Border(lipgloss.NormalBorder(), true, false, false, true).
		BorderForeground(panelBorder).
		Padding(0, 1).
		Render(content)
}

func renderMessages(c *model.Call, width int) string {
	var blocks []string

	if c.SystemPrompt != "" {
		blocks = append(blocks, renderMessageBlock("system", c.SystemPrompt, width, lipgloss.Color("#F18E94"), lipgloss.Color("#351A20")))
	}

	if len(c.Messages) == 0 && c.RequestBody != "" {
		return truncateParagraph(normalizeWhitespace(c.RequestBody), max(120, width*8))
	}

	for _, msg := range c.Messages {
		fg := brandCyan
		bg := lipgloss.Color("#18283A")
		switch msg.Role {
		case "assistant":
			fg = brandMint
			bg = lipgloss.Color("#162C26")
		case "system":
			fg = lipgloss.Color("#F18E94")
			bg = lipgloss.Color("#351A20")
		}
		blocks = append(blocks, renderMessageBlock(msg.Role, msg.Content, width, fg, bg))
	}

	if len(blocks) == 0 {
		return "No request payload captured yet."
	}

	return lipgloss.JoinVertical(lipgloss.Left, blocks...)
}

func renderMessageBlock(role, content string, width int, fg, bg lipgloss.Color) string {
	roleBadge := lipgloss.NewStyle().
		Foreground(fg).
		Bold(true).
		Render(strings.ToUpper(role))

	text := lipgloss.NewStyle().
		Foreground(textSecondary).
		Width(max(16, width)).
		Background(surfaceBg).
		Render(truncateParagraph(normalizeWhitespace(content), max(140, width*10)))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lineWithBg(roleBadge, width, surfaceBg, fg),
		text,
	)
}

func renderResponse(c *model.Call, width int) string {
	response := normalizeWhitespace(c.ResponseText)
	if response == "" {
		response = normalizeWhitespace(c.ResponseBody)
	}
	if response == "" {
		return "No response body captured yet."
	}
	return truncateParagraph(response, max(180, width*12))
}

func renderMetricPill(label, value string, fg, bg lipgloss.Color) string {
	return lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Padding(0, 1).
		MarginRight(1).
		Render(strings.ToUpper(label) + " " + value)
}

func renderStatusPill(c *model.Call) string {
	if c.IsStreaming {
		return lipgloss.NewStyle().
			Foreground(brandMint).
			Bold(true).
			Render("STREAM")
	}

	statusColor := textSecondary
	if c.StatusCode >= 400 {
		statusColor = dangerColor
	}
	return lipgloss.NewStyle().
		Foreground(statusColor).
		Bold(true).
		Render(fmt.Sprintf("HTTP %d", c.StatusCode))
}

func renderInlineMetric(label, value string, valueColor lipgloss.Color) string {
	labelStyle := lipgloss.NewStyle().Foreground(textMuted).Bold(true).Render(label)
	valueStyle := lipgloss.NewStyle().Foreground(valueColor).Bold(true).Render(value)
	return labelStyle + " " + valueStyle + "   "
}

func fillHorizontal(left, right string, width int) string {
	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	if leftW+rightW >= width {
		if width <= rightW+1 {
			return truncateRunes(left, width)
		}
		left = truncateRunes(left, width-rightW-1)
		leftW = lipgloss.Width(left)
	}
	gap := max(1, width-leftW-rightW)
	return left + strings.Repeat(" ", gap) + right
}

func truncateParagraph(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	return truncateRunes(s, maxRunes-1) + "..."
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	rs := []rune(s)
	return string(rs[:maxRunes])
}

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func firstPreview(c *model.Call) string {
	if len(c.Messages) > 0 {
		return normalizeWhitespace(c.Messages[len(c.Messages)-1].Content)
	}
	if c.SystemPrompt != "" {
		return normalizeWhitespace(c.SystemPrompt)
	}
	if c.ResponseText != "" {
		return normalizeWhitespace(c.ResponseText)
	}
	if c.Endpoint != "" {
		return endpointLabel(c.Endpoint)
	}
	return "Waiting for payload details..."
}

func endpointLabel(endpoint string) string {
	if endpoint == "" {
		return "no endpoint captured"
	}
	return strings.TrimPrefix(endpoint, "/")
}

func fallbackModel(modelName string) string {
	if strings.TrimSpace(modelName) == "" {
		return "model pending"
	}
	return modelName
}

func themeFor(provider model.Provider) providerTheme {
	if theme, ok := providerThemes[provider]; ok {
		return theme
	}
	return providerThemes[model.ProviderUnknown]
}

func yesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

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
	if d <= 0 {
		return "pending"
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

func (m Model) currentFilter() string {
	if len(m.session.Calls) == 0 {
		return ""
	}
	return ""
}

func boxedInnerWidth(totalWidth int) int {
	return max(16, totalWidth-4)
}

func boxedInnerHeight(totalHeight int) int {
	return max(6, totalHeight-4)
}

func lineWithBg(content string, width int, bg, fg lipgloss.Color) string {
	return lipgloss.NewStyle().
		Width(max(1, width)).
		Background(bg).
		Foreground(fg).
		Render(content)
}
