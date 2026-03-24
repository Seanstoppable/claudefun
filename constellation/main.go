package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	text := flag.String("text", "", "input sentence to transform into a constellation")
	svgPath := flag.String("svg", "", "output path for SVG export")
	width := flag.Int("width", 60, "terminal render width")
	height := flag.Int("height", 30, "terminal render height")
	svgWidth := flag.Int("svg-width", 800, "SVG width")
	svgHeight := flag.Int("svg-height", 600, "SVG height")
	noMyth := flag.Bool("no-myth", false, "skip mythology output")
	minimal := flag.Bool("minimal", false, "show only the constellation art")

	flag.Parse()

	// Resolve input: -text flag takes priority, then first positional arg.
	input := *text
	if input == "" && flag.NArg() > 0 {
		input = strings.Join(flag.Args(), " ")
	}
	if input == "" {
		fmt.Println()
		fmt.Println("  Feed me words and I'll give you stars! ✨")
		fmt.Println()
		fmt.Println("  Usage:")
		fmt.Println(`    constellation "hello world"`)
		fmt.Println(`    constellation -text "hello world" -svg out.svg`)
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(1)
	}

	// --- Generate ---
	sm := GenerateStarMap(input)
	cons := Connect(sm)
	name := sm.ConstellationName()

	termRenderer := NewTerminalRenderer(*width, *height)
	art := termRenderer.Render(cons)

	// --- Styles ---
	purple := lipgloss.Color("#7B68EE")
	gold := lipgloss.Color("#FFD700")
	dim := lipgloss.Color("#6C6C6C")
	mythTitle := lipgloss.Color("#C8A2C8")

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#E0D0FF")).
		Background(purple).
		Padding(0, 2).
		Align(lipgloss.Center)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purple).
		Padding(0, 1).
		Align(lipgloss.Center)

	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(gold)

	inputStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#A0A0D0"))

	dimStyle := lipgloss.NewStyle().
		Foreground(dim)

	sepStyle := lipgloss.NewStyle().
		Foreground(purple)

	mythHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(mythTitle).
		MarginTop(1)

	moralStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#C0C0FF"))

	viewingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB"))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90EE90"))

	separator := sepStyle.Render(strings.Repeat("━", *width))

	// --- Minimal mode ---
	if *minimal {
		fmt.Println(art)
		if *svgPath != "" {
			exportSVG(cons, *svgPath, *svgWidth, *svgHeight, successStyle)
		}
		return
	}

	// --- Full output ---
	fmt.Println()

	// Header box
	header := headerStyle.Render("✨  CONSTELLATION COMPOSER  ✨")
	fmt.Println(boxStyle.Width(*width).Render(header))
	fmt.Println()

	// Constellation name bar
	nameBar := fmt.Sprintf("━━━━━━━━━━━ ☆ %s ☆ ━━━━━━━━━━━", name)
	fmt.Println(sepStyle.Render(nameBar))
	fmt.Println()

	// Input info
	info := fmt.Sprintf("%s  →  %s",
		inputStyle.Render(fmt.Sprintf("%q", input)),
		dimStyle.Render(fmt.Sprintf("%d stars, %d edges", cons.StarCount(), cons.EdgeCount())),
	)
	fmt.Println("  " + info)
	fmt.Println()

	// Constellation art
	fmt.Println(art)
	fmt.Println()
	fmt.Println(separator)

	// Mythology
	if !*noMyth {
		mg := NewMythologyGenerator(sm.Seed)
		myth := mg.Generate(name, cons.StarCount())

		fmt.Println()
		fmt.Println(mythHeaderStyle.Render(fmt.Sprintf("  THE MYTH OF %s", strings.ToUpper(myth.Name))))
		fmt.Println()

		// Word-wrap the story
		for _, line := range wordWrap(myth.Story, *width-4) {
			fmt.Println("  " + nameStyle.Copy().UnsetBold().Foreground(lipgloss.Color("#D0D0E0")).Render(line))
		}
		fmt.Println()

		fmt.Println("  " + moralStyle.Render(fmt.Sprintf("✦ Moral: %q", myth.Moral)))
		fmt.Println()
		fmt.Println("  " + viewingStyle.Render(fmt.Sprintf("🔭 Best viewed: %s", myth.BestViewing)))
		fmt.Println()
		fmt.Println(separator)
	}

	// SVG export
	if *svgPath != "" {
		fmt.Println()
		exportSVG(cons, *svgPath, *svgWidth, *svgHeight, successStyle)
	}

	fmt.Println()
}

func exportSVG(cons *Constellation, path string, w, h int, style lipgloss.Style) {
	svgR := NewSVGRenderer(w, h)
	if err := svgR.RenderToFile(cons, path); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing SVG: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(style.Render(fmt.Sprintf("  ✨ SVG saved to: %s", path)))
}

func wordWrap(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		maxWidth = 60
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	current := words[0]
	for _, w := range words[1:] {
		if len(current)+1+len(w) > maxWidth {
			lines = append(lines, current)
			current = w
		} else {
			current += " " + w
		}
	}
	lines = append(lines, current)
	return lines
}
