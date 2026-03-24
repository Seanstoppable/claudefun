package main

// model.go — Bubbletea TUI model for Mood Octopus.
// Lives in package main (not octopus) to avoid an import cycle:
// octopus ← mood ← octopus.

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ssmith/mood-octopus/mood"
	"github.com/ssmith/mood-octopus/octopus"
)

// Model is the top-level bubbletea model for Mood Octopus.
type Model struct {
	textInput textinput.Model
	animation *octopus.AnimationState
	advisor   *octopus.Advisor
	analyzer  *mood.Analyzer
	history   *mood.History

	currentMood  octopus.Emotion
	moodResults  []mood.MoodResult
	adviceBubble string
	greeting     string
	width        int
	height       int
	quitting     bool
}

// NewModel initialises the TUI model with all subsystems.
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Type something to feed the octopus..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	history, err := mood.NewHistory()
	if err != nil {
		history = nil
	}

	var greeting string
	if history != nil {
		greeting = history.StartupGreeting()
	} else {
		greeting = "Hello! I'm your new octopus friend! 🐙 Type anything and watch me react!"
	}

	return Model{
		textInput:   ti,
		animation:   octopus.NewAnimationState(),
		advisor:     octopus.NewAdvisor(),
		analyzer:    mood.NewAnalyzer(),
		history:     history,
		currentMood: octopus.Curiosity,
		greeting:    greeting,
	}
}

// Init starts the animation tick loop and the text-input cursor blink.
func (m Model) Init() tea.Cmd {
	return tea.Batch(octopus.AnimTick(), textinput.Blink)
}

// Update handles all incoming messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			return m.handleInput()
		}

	case octopus.AnimTickMsg:
		return m.handleAnimTick()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = max(20, m.width-8)
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) handleInput() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.textInput.Value())
	if input == "" {
		return m, nil
	}

	m.moodResults = m.analyzer.Analyze(input)
	m.currentMood = m.moodResults[0].Emotion
	m.animation.SetEmotion(m.currentMood)

	if m.history != nil {
		_ = m.history.Record(m.currentMood, input)
	}

	if m.advisor.ShouldGiveAdvice() {
		m.adviceBubble = m.advisor.GetAdvice(m.currentMood)
	} else {
		m.adviceBubble = ""
	}

	m.greeting = ""
	m.textInput.Reset()

	return m, nil
}

func (m Model) handleAnimTick() (tea.Model, tea.Cmd) {
	_ = m.animation.Tick()

	var cmd tea.Cmd
	if m.animation.IsTransitioning() {
		cmd = octopus.AnimTickFast()
	} else {
		cmd = octopus.AnimTick()
	}
	return m, cmd
}

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF69B4"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	inputBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	greetingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Italic(true)

	adviceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1ABC9C"))
)

func moodStyle(e octopus.Emotion) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(e.Info().Color).Bold(true)
}

// View renders the full TUI.
func (m Model) View() string {
	if m.quitting {
		return "Goodbye from your octopus friend! 🐙👋\n"
	}

	termW := m.width
	if termW <= 0 {
		termW = 80
	}
	termH := m.height
	if termH <= 0 {
		termH = 24
	}

	var sections []string

	// Greeting or advice bubble.
	if m.greeting != "" {
		g := greetingStyle.Render(m.greeting)
		sections = append(sections, lipgloss.PlaceHorizontal(termW, lipgloss.Center, g))
	} else if m.adviceBubble != "" {
		a := adviceStyle.Render(m.adviceBubble)
		sections = append(sections, lipgloss.PlaceHorizontal(termW, lipgloss.Center, a))
	}

	// Octopus ASCII frame (read-only — Tick advances in Update).
	frame := m.animation.CurrentFrame()
	sections = append(sections, lipgloss.PlaceHorizontal(termW, lipgloss.Center, frame))

	// Mood indicator.
	if len(m.moodResults) > 0 {
		info := m.currentMood.Info()
		conf := m.moodResults[0].Confidence
		moodLine := moodStyle(m.currentMood).Render(
			fmt.Sprintf("Current mood: %s %s (confidence: %.2f)", info.Emoji, info.Name, conf),
		)
		sections = append(sections, lipgloss.PlaceHorizontal(termW, lipgloss.Center, moodLine))
	}

	// Mood sparkline.
	if m.history != nil {
		spark := m.history.MoodSparkline(20)
		if spark != "" {
			sparkLine := subtitleStyle.Render("Mood history: ") + spark
			sections = append(sections, lipgloss.PlaceHorizontal(termW, lipgloss.Center, sparkLine))
		}
	}

	sections = append(sections, "")

	// Text input.
	inputW := min(termW-6, 60)
	if inputW < 10 {
		inputW = 10
	}
	inputBox := inputBorderStyle.Width(inputW).Render(m.textInput.View())
	sections = append(sections, lipgloss.PlaceHorizontal(termW, lipgloss.Center, inputBox))

	// Footer.
	footer := titleStyle.Render("🐙 Mood Octopus") +
		subtitleStyle.Render("  │  ctrl+c to quit")
	sections = append(sections, lipgloss.PlaceHorizontal(termW, lipgloss.Center, footer))

	body := strings.Join(sections, "\n")

	bodyLines := strings.Count(body, "\n") + 1
	padTop := (termH - bodyLines) / 2
	if padTop < 0 {
		padTop = 0
	}

	return strings.Repeat("\n", padTop) + body + "\n"
}
