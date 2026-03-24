package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Game phases ─────────────────────────────────────────────────────────────

type gamePhase int

const (
	phaseWelcome  gamePhase = iota // show welcome narration
	phasePolicy                    // presenting a policy choice
	phaseResult                    // showing result of choice + bard narration
	phaseEvent                     // showing random event (if any)
	phaseGameOver                  // game over screen
)

// ── Model ───────────────────────────────────────────────────────────────────

// Model is the top-level bubbletea model for the game.
type Model struct {
	kingdom  *Kingdom
	policies *PolicyGenerator
	events   *EventGenerator
	bard     *Bard
	phase    gamePhase

	currentPolicy Policy
	currentEvent  *Event
	narration     string // current bard narration text
	resultText    string // what happened after a choice

	width, height int
	quitting      bool
}

// NewGameModel creates a fresh game model ready to play.
func NewGameModel() Model {
	seed := time.Now().UnixNano()
	k := NewKingdom("")
	b := NewBard(seed+2, k.Name)
	return Model{
		kingdom:  k,
		policies: NewPolicyGenerator(seed),
		events:   NewEventGenerator(seed + 1),
		bard:     b,
		phase:    phaseWelcome,
		narration: b.NarrateWelcome(k.Name),
	}
}

// ── Styles ──────────────────────────────────────────────────────────────────

var (
	gold   = lipgloss.Color("#FFD700")
	amber  = lipgloss.Color("#FFBF00")
	cream  = lipgloss.Color("#FFF8DC")
	red    = lipgloss.Color("#FF4444")
	green  = lipgloss.Color("#44FF77")
	yellow = lipgloss.Color("#FFFF55")
	dim    = lipgloss.Color("#888888")
	white  = lipgloss.Color("#FFFFFF")

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gold).
			Foreground(gold).
			Bold(true).
			Align(lipgloss.Center).
			Padding(0, 2)

	statLabelStyle = lipgloss.NewStyle().Bold(true).Foreground(white)
	statGreen      = lipgloss.NewStyle().Foreground(green).Bold(true)
	statYellow     = lipgloss.NewStyle().Foreground(yellow).Bold(true)
	statRed        = lipgloss.NewStyle().Foreground(red).Bold(true)

	factionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#C9B1FF"))

	divider = lipgloss.NewStyle().Foreground(dim)

	policyBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#88AAFF")).
			Padding(0, 1)

	bardStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(amber)

	optionAStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#55FF55"))
	optionBStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF8855"))

	footerStyle = lipgloss.NewStyle().Foreground(dim)

	gameOverVictoryStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(gold).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(gold).
				Padding(1, 3).
				Align(lipgloss.Center)

	gameOverDefeatStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(red).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(red).
				Padding(1, 3).
				Align(lipgloss.Center)

	welcomeStyle = lipgloss.NewStyle().
			Foreground(cream).
			Padding(1, 2)

	resultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AADDFF")).
			Padding(0, 1)
)

// ── Bubbletea interface ─────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// Global quit
		if key == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.phase {
		case phaseWelcome:
			return m.handleWelcome(key)
		case phasePolicy:
			return m.handlePolicy(key)
		case phaseResult:
			return m.handleResult(key)
		case phaseEvent:
			return m.handleEvent(key)
		case phaseGameOver:
			return m.handleGameOver(key)
		}
	}
	return m, nil
}

func (m Model) handleWelcome(key string) (tea.Model, tea.Cmd) {
	if key == "q" {
		m.quitting = true
		return m, tea.Quit
	}
	// Any key → go to first policy
	m.phase = phasePolicy
	m.currentPolicy = m.policies.NextPolicy()
	m.narration = m.bard.NarrateTurnStart(m.kingdom.Turn)
	return m, nil
}

func (m Model) handlePolicy(key string) (tea.Model, tea.Cmd) {
	if key == "q" {
		m.quitting = true
		return m, tea.Quit
	}

	switch key {
	case "a", "A":
		m.kingdom.ApplyEffect(m.currentPolicy.EffectA)
		m.resultText = m.currentPolicy.FlavorA
		m.narration = m.bard.NarratePolicy(
			m.currentPolicy.Question,
			m.currentPolicy.OptionA,
			m.currentPolicy.FlavorA,
		)
		m.phase = phaseResult
	case "b", "B":
		m.kingdom.ApplyEffect(m.currentPolicy.EffectB)
		m.resultText = m.currentPolicy.FlavorB
		m.narration = m.bard.NarratePolicy(
			m.currentPolicy.Question,
			m.currentPolicy.OptionB,
			m.currentPolicy.FlavorB,
		)
		m.phase = phaseResult
	}
	return m, nil
}

func (m Model) handleResult(key string) (tea.Model, tea.Cmd) {
	// Check for random event
	event := m.events.MaybeEvent(m.kingdom)
	if event != nil {
		m.currentEvent = event
		m.kingdom.ApplyEffect(event.Effect)
		m.narration = m.bard.NarrateEvent(event.Description)
		m.phase = phaseEvent
		return m, nil
	}

	// No event — advance turn
	return m.advanceTurn()
}

func (m Model) handleEvent(key string) (tea.Model, tea.Cmd) {
	return m.advanceTurn()
}

func (m Model) handleGameOver(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "Q":
		m.quitting = true
		return m, tea.Quit
	case "r", "R":
		return NewGameModel(), nil
	}
	return m, nil
}

func (m Model) advanceTurn() (tea.Model, tea.Cmd) {
	m.kingdom.AdvanceTurn()

	if m.kingdom.GameOver {
		m.narration = m.bard.NarrateGameOver(m.kingdom.Victory, m.kingdom.GameOverMsg)
		m.phase = phaseGameOver
		return m, nil
	}

	if !m.policies.HasMore() {
		m.kingdom.GameOver = true
		m.kingdom.Victory = true
		m.kingdom.GameOverMsg = "You've answered every possible policy question! The kingdom is baffled but impressed."
		m.narration = m.bard.NarrateGameOver(true, m.kingdom.GameOverMsg)
		m.phase = phaseGameOver
		return m, nil
	}

	m.currentPolicy = m.policies.NextPolicy()
	m.narration = m.bard.NarrateTurnStart(m.kingdom.Turn)
	m.phase = phasePolicy
	return m, nil
}

// ── View ────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.quitting {
		return "Farewell, O Ruler! The bard shall remember your... attempt.\n"
	}

	w := m.width
	if w == 0 {
		w = 80
	}

	var sections []string

	// Header — always shown
	sections = append(sections, m.viewHeader(w))

	switch m.phase {
	case phaseWelcome:
		sections = append(sections, m.viewWelcome(w))
	case phasePolicy:
		sections = append(sections, m.viewStats(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewPolicy(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewBard(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewFooterPolicy())
	case phaseResult:
		sections = append(sections, m.viewStats(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewResultText(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewBard(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewFooterContinue())
	case phaseEvent:
		sections = append(sections, m.viewStats(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewEventText(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewBard(w))
		sections = append(sections, m.viewDivider(w))
		sections = append(sections, m.viewFooterContinue())
	case phaseGameOver:
		sections = append(sections, m.viewGameOver(w))
	}

	return strings.Join(sections, "\n\n")
}

// ── View helpers ────────────────────────────────────────────────────────────

func (m Model) viewHeader(w int) string {
	title := "👑 TINY KINGDOM SIMULATOR 👑"
	subtitle := fmt.Sprintf(`"%s"`, m.kingdom.Name)
	info := fmt.Sprintf("Turn %d — %s", m.kingdom.Turn, m.kingdom.RulerTitle)
	content := fmt.Sprintf("%s\n%s\n%s", title, subtitle, info)
	return headerStyle.Width(min(w-4, 58)).Render(content)
}

func (m Model) viewWelcome(w int) string {
	text := welcomeStyle.Width(min(w-4, 60)).Render(m.narration)
	footer := footerStyle.Render("  Press any key to begin your reign...")
	return text + "\n\n" + footer
}

func (m Model) viewStats(w int) string {
	k := m.kingdom

	stats := []struct {
		emoji string
		label string
		value int
		max   int
	}{
		{"💰", "Treasury", k.Treasury, 0},
		{"👥", "Population", k.Population, 0},
		{"😊", "Happiness", k.Happiness, 100},
		{"⚔️ ", "Military", k.Military, 100},
		{"🎭", "Culture", k.Culture, 100},
		{"🍞", "Food", k.Food, 100},
		{"🌟", "Reputation", k.Reputation, 100},
	}

	var rows []string
	for i := 0; i < len(stats); i += 2 {
		left := formatStat(stats[i].emoji, stats[i].label, stats[i].value, stats[i].max)
		right := ""
		if i+1 < len(stats) {
			right = formatStat(stats[i+1].emoji, stats[i+1].label, stats[i+1].value, stats[i+1].max)
		}
		if right != "" {
			rows = append(rows, fmt.Sprintf("  %-36s%s", left, right))
		} else {
			rows = append(rows, fmt.Sprintf("  %s", left))
		}
	}

	// Factions
	fParts := make([]string, 0, len(AllFactions))
	for _, f := range AllFactions {
		mood := k.FactionMood[f]
		fParts = append(fParts, fmt.Sprintf("%s %+d", string(f), mood))
	}
	factionLine := factionStyle.Render("  Factions: " + strings.Join(fParts, " │ "))

	return strings.Join(rows, "\n") + "\n\n" + factionLine
}

func formatStat(emoji, label string, value, max int) string {
	var valStr string
	if max > 0 {
		bar := renderBar(value, max)
		valStr = fmt.Sprintf("%d  %s", value, bar)
	} else {
		valStr = fmt.Sprintf("%d", value)
	}

	styled := colorByValue(valStr, value)
	return fmt.Sprintf("%s %s: %s",
		emoji,
		statLabelStyle.Render(label),
		styled,
	)
}

func renderBar(val, max int) string {
	const barLen = 10
	filled := val * barLen / max
	if filled < 0 {
		filled = 0
	}
	if filled > barLen {
		filled = barLen
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)
}

func colorByValue(text string, val int) string {
	switch {
	case val > 60:
		return statGreen.Render(text)
	case val >= 30:
		return statYellow.Render(text)
	default:
		return statRed.Render(text)
	}
}

func (m Model) viewDivider(w int) string {
	dw := min(w-4, 58)
	if dw < 10 {
		dw = 58
	}
	return divider.Render("  " + strings.Repeat("━", dw))
}

func (m Model) viewPolicy(w int) string {
	p := m.currentPolicy
	question := fmt.Sprintf("📜 ROYAL DECREE REQUIRED:\n\n\"%s\"", p.Question)
	opts := fmt.Sprintf(
		"%s %s\n%s %s",
		optionAStyle.Render("[A]"), p.OptionA,
		optionBStyle.Render("[B]"), p.OptionB,
	)
	content := question + "\n\n" + opts
	return policyBoxStyle.Width(min(w-4, 58)).Render(content)
}

func (m Model) viewBard(w int) string {
	if m.narration == "" {
		return ""
	}
	header := statLabelStyle.Render("  🎭 The Bard Speaks:")
	body := bardStyle.Render(indentLines(m.narration, "  "))
	return header + "\n" + body
}

func (m Model) viewResultText(w int) string {
	header := statLabelStyle.Render("  📜 What Happened:")
	body := resultStyle.Render(indentLines(m.resultText, "  "))
	return header + "\n" + body
}

func (m Model) viewEventText(w int) string {
	if m.currentEvent == nil {
		return ""
	}
	header := statLabelStyle.Render(fmt.Sprintf("  ⚡ EVENT: %s", m.currentEvent.Name))
	body := resultStyle.Render(indentLines(m.currentEvent.Description, "  "))
	return header + "\n" + body
}

func (m Model) viewFooterPolicy() string {
	return footerStyle.Render("  [A/B] Choose  │  [Q] Abdicate")
}

func (m Model) viewFooterContinue() string {
	return footerStyle.Render("  Press any key to continue  │  [Q] Abdicate")
}

func (m Model) viewGameOver(w int) string {
	k := m.kingdom
	var box string
	if k.Victory {
		title := "🏆 VICTORY! 🏆"
		content := fmt.Sprintf("%s\n\n%s\n\nFinal Turn: %d  │  Treasury: %d  │  Population: %d",
			title, k.GameOverMsg, k.Turn, k.Treasury, k.Population)
		box = gameOverVictoryStyle.Width(min(w-4, 58)).Render(content)
	} else {
		title := "💀 GAME OVER 💀"
		content := fmt.Sprintf("%s\n\n%s\n\nFinal Turn: %d  │  Treasury: %d  │  Population: %d",
			title, k.GameOverMsg, k.Turn, k.Treasury, k.Population)
		box = gameOverDefeatStyle.Width(min(w-4, 58)).Render(content)
	}

	bard := ""
	if m.narration != "" {
		bard = "\n\n" + bardStyle.Render(indentLines(m.narration, "  "))
	}

	footer := footerStyle.Render("  [R] Restart  │  [Q] Quit")

	return box + bard + "\n\n" + footer
}

// ── Utilities ───────────────────────────────────────────────────────────────

func indentLines(s string, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
