package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xeonx/timeago"
)

type tickMsg struct{}

type activitiesMsg struct {
	activities Activities
	err        error
	fetchedAt  time.Time
}

type model struct {
	activities Activities
	err        error
	width      int
	height     int
	ready      bool
	updatedAt  time.Time
	viewport   viewport.Model
}

var (
	docStyle = lipgloss.NewStyle().
			Padding(0, 2)

	bannerStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 2).
			Foreground(lipgloss.Color("#E2E8F0")).
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#38BDF8"))

	bannerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F8FAFC"))

	bannerSubtleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#94A3B8"))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#334155")).
			Padding(1, 2)

	cardLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8"))

	cardValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F8FAFC"))

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#334155")).
			Padding(1, 2)

	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E2E8F0"))

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#CBD5E1"))

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8"))

	liveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#22C55E"))

	staleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F59E0B"))

	highRiskStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F97316"))

	criticalRiskStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#EF4444"))

	okStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ADE80"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F87171")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#94A3B8"))
)

func newModel() model {
	return model{}
}

func fetchActivitiesCmd() tea.Cmd {
	return func() tea.Msg {
		activities, err := GetActivities()
		return activitiesMsg{
			activities: activities,
			err:        err,
			fetchedAt:  time.Now(),
		}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchActivitiesCmd(), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.viewport = viewport.New(m.contentWidth(), m.viewportHeight())
		m.viewport.SetContent(m.renderContent())
	case tickMsg:
		return m, tea.Batch(fetchActivitiesCmd(), tickCmd())
	case activitiesMsg:
		m.updatedAt = msg.fetchedAt
		m.err = msg.err
		if msg.err == nil {
			m.activities = msg.activities
		}
		if m.ready {
			yOffset := m.viewport.YOffset
			m.viewport.SetContent(m.renderContent())
			m.viewport.SetYOffset(yOffset)
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading dashboard..."
	}

	footer := m.renderFooter()
	body := docStyle.Width(max(0, m.width)).Render(m.viewport.View())

	return lipgloss.JoinVertical(lipgloss.Left, body, footer)
}

func (m model) renderContent() string {
	sections := []string{
		"",
		m.renderBanner(),
		m.renderStats(),
		m.renderBody(),
	}

	if m.err != nil {
		sections = append(sections, errorStyle.Render("Refresh error: "+m.err.Error()))
	}

	sections = append(sections, "")

	return strings.Join(sections, "\n\n")
}

func (m model) renderBanner() string {
	title := bannerTitleStyle.Render("SERVER WATCHDOG")
	subtitle := bannerSubtleStyle.Render("Live SSH session monitoring")

	status := liveStyle.Render("LIVE FEED")
	if m.err != nil {
		status = errorStyle.Render("DEGRADED")
	}

	lastUpdated := "syncing"
	if !m.updatedAt.IsZero() {
		lastUpdated = m.updatedAt.Format("15:04:05")
	}

	left := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
	)

	right := lipgloss.JoinVertical(
		lipgloss.Right,
		status,
		bannerSubtleStyle.Render("Updated "+lastUpdated),
	)

	width := m.contentWidth()
	if width < 48 {
		return bannerStyle.Width(width).Render(lipgloss.JoinVertical(lipgloss.Left, left, right))
	}

	return bannerStyle.Width(width).Render(lipgloss.JoinHorizontal(lipgloss.Top, left, spacer(width-lipgloss.Width(left)-lipgloss.Width(right)-8), right))
}

func (m model) renderStats() string {
	cards := []string{
		m.renderStatCard("Tracked IPs", strconv.Itoa(len(m.activities.Attempts)), "Historical and live"),
		m.renderStatCard("Live Sessions", strconv.Itoa(m.liveCount()), "Currently connected"),
		m.renderStatCard("Whitelist", strconv.Itoa(len(m.activities.WhitelistedIPs)), "Trusted sources"),
		m.renderStatCard("Risk Level", m.riskSummary(), "Based on duration"),
	}

	if m.width >= 120 {
		return lipgloss.JoinHorizontal(lipgloss.Top, cards...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, cards...)
}

func (m model) renderStatCard(label, value, hint string) string {
	width := m.contentWidth()
	if m.width >= 120 {
		width = max(18, (m.contentWidth()-3)/4)
	}

	valueStyle := cardValueStyle
	switch label {
	case "Live Sessions":
		valueStyle = liveStyle
	case "Risk Level":
		valueStyle = m.riskStyle()
	}

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		cardLabelStyle.Render(strings.ToUpper(label)),
		valueStyle.Render(value),
		metaStyle.Render(hint),
	)

	return cardStyle.Width(width).Render(body)
}

func (m model) renderBody() string {
	mainPanel := m.renderAttemptsTable()
	sidePanel := m.renderSidebar()

	if m.width >= 130 {
		leftWidth := max(60, m.contentWidth()-34)
		rightWidth := max(28, m.contentWidth()-leftWidth-2)
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			panelStyle.Width(leftWidth).Render(mainPanel),
			panelStyle.Width(rightWidth).Render(sidePanel),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		panelStyle.Width(m.contentWidth()).Render(mainPanel),
		panelStyle.Width(m.contentWidth()).Render(sidePanel),
	)
}

func (m model) renderAttemptsTable() string {
	headers := []string{"#", "IP", "State", "Last Seen", "Country", "Sessions"}
	widths := []int{4, 18, 10, 14, 16, 26}

	lines := []string{
		panelTitleStyle.Render("ATTEMPT FEED"),
		metaStyle.Render("Sorted with live sessions first, then repeat offenders."),
		"",
		m.renderTableRow(headers, widths, tableHeaderStyle),
	}

	if len(m.activities.Attempts) == 0 {
		lines = append(lines, "", warningStyle.Render("No SSH attempts recorded yet."))
		return strings.Join(lines, "\n")
	}

	for idx, attempt := range m.activities.Attempts {
		row := []string{
			strconv.Itoa(idx + 1),
			attempt.IP,
			attempt.Status,
			m.lastSeenText(attempt),
			attempt.Country,
			fmt.Sprintf("%d attempts (%s)", len(attempt.Sessions), FormatDuration(GetTotalDuration(attempt.Sessions))),
		}
		lines = append(lines, m.renderAttemptRow(row, widths, attempt))
	}

	return strings.Join(lines, "\n")
}

func (m model) renderSidebar() string {
	parts := []string{
		panelTitleStyle.Render("SIDEBAR"),
		metaStyle.Render("Quick operational context."),
		"",
		m.renderWhitelistBlock(),
		"",
		m.renderHotspotsBlock(),
	}

	return strings.Join(parts, "\n")
}

func (m model) renderFooter() string {
	parts := []string{
		"j/k or arrows: scroll",
		"pgup/pgdn: jump",
		"g/G: top/bottom",
		"q: quit",
	}

	if m.err != nil {
		parts = append(parts, "status: degraded")
	} else {
		parts = append(parts, "status: live")
	}

	return helpStyle.Width(max(0, m.width)).Render(strings.Join(parts, "  |  "))
}

func (m model) renderWhitelistBlock() string {
	lines := []string{panelTitleStyle.Render("WHITELIST")}
	if len(m.activities.WhitelistedIPs) == 0 {
		lines = append(lines, warningStyle.Render("No whitelisted IPs"))
		return strings.Join(lines, "\n")
	}

	for _, ip := range m.activities.WhitelistedIPs {
		lines = append(lines, okStyle.Render("OK ")+ip)
	}

	return strings.Join(lines, "\n")
}

func (m model) renderHotspotsBlock() string {
	lines := []string{
		panelTitleStyle.Render("HOTSPOTS"),
	}

	if len(m.activities.Attempts) == 0 {
		lines = append(lines, metaStyle.Render("No traffic yet"))
		return strings.Join(lines, "\n")
	}

	limit := min(3, len(m.activities.Attempts))
	for i := 0; i < limit; i++ {
		attempt := m.activities.Attempts[i]
		duration := FormatDuration(GetTotalDuration(attempt.Sessions))
		label := fmt.Sprintf("%s  %s", attempt.IP, m.riskBadge(attempt))
		lines = append(lines, truncate(label, 28))
		lines = append(lines, metaStyle.Render(truncate(duration, 28)))
	}

	return strings.Join(lines, "\n")
}

func (m model) renderAttemptRow(values []string, widths []int, attempt SSHAttempt) string {
	r, g, b := GetWarmth(attempt.Sessions)
	rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b)))

	rendered := make([]string, 0, len(values))
	for i, value := range values {
		if i == 2 {
			value = m.renderStatus(attempt.Status)
		}
		cellStyle := rowStyle.Width(widths[i]).MaxWidth(widths[i])
		rendered = append(rendered, cellStyle.Render(truncate(value, widths[i])))
	}

	return strings.Join(rendered, " ")
}

func (m model) renderTableRow(values []string, widths []int, style lipgloss.Style) string {
	rendered := make([]string, 0, len(values))
	for i, value := range values {
		rendered = append(rendered, style.Width(widths[i]).MaxWidth(widths[i]).Render(truncate(value, widths[i])))
	}
	return strings.Join(rendered, " ")
}

func (m model) renderStatus(status string) string {
	if status == LIVE {
		return liveStyle.Render("LIVE")
	}
	return staleStyle.Render("STALE")
}

func (m model) lastSeenText(attempt SSHAttempt) string {
	if attempt.Status == LIVE {
		return "now"
	}

	if len(attempt.Sessions) == 0 {
		return "n/a"
	}

	lastSession := attempt.Sessions[len(attempt.Sessions)-1]
	return timeago.English.Format(lastSession.End)
}

func (m model) liveCount() int {
	count := 0
	for _, attempt := range m.activities.Attempts {
		if attempt.Status == LIVE {
			count++
		}
	}
	return count
}

func (m model) riskSummary() string {
	maxSeverity := 0
	for _, attempt := range m.activities.Attempts {
		if severity := riskSeverity(attempt.Sessions); severity > maxSeverity {
			maxSeverity = severity
		}
	}

	switch maxSeverity {
	case 3:
		return "Critical"
	case 2:
		return "High"
	case 1:
		return "Elevated"
	default:
		return "Low"
	}
}

func (m model) riskStyle() lipgloss.Style {
	switch m.riskSummary() {
	case "Critical":
		return criticalRiskStyle
	case "High":
		return highRiskStyle
	case "Elevated":
		return staleStyle
	default:
		return okStyle
	}
}

func (m model) riskBadge(attempt SSHAttempt) string {
	switch riskSeverity(attempt.Sessions) {
	case 3:
		return criticalRiskStyle.Render("CRITICAL")
	case 2:
		return highRiskStyle.Render("HIGH")
	case 1:
		return staleStyle.Render("ELEVATED")
	default:
		return okStyle.Render("LOW")
	}
}

func (m model) contentWidth() int {
	if m.width <= 0 {
		return 80
	}

	width := m.width - 6
	if width < 40 {
		return 40
	}

	return width
}

func (m model) viewportHeight() int {
	height := m.height - 1
	if height < 8 {
		return 8
	}
	return height
}

func riskSeverity(sessions []AttemptSession) int {
	total := GetTotalDuration(sessions)
	switch {
	case total >= 300:
		return 3
	case total >= 60:
		return 2
	case total > 0:
		return 1
	default:
		return 0
	}
}

func truncate(value string, width int) string {
	if lipgloss.Width(value) <= width {
		return value
	}

	if width <= 1 {
		runes := []rune(value)
		if len(runes) == 0 {
			return ""
		}
		return string(runes[:1])
	}

	runes := []rune(value)
	if len(runes) <= width {
		return value
	}

	return string(runes[:width-1]) + "…"
}

func spacer(width int) string {
	if width <= 0 {
		return ""
	}
	return strings.Repeat(" ", width)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
