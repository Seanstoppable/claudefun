package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── CLI styles ────────────────────────────────────────────────────────

var (
	headerBoxStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#f0d9b5")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#d4a843")).
			Padding(0, 2).
			Align(lipgloss.Center).
			Width(56)

	generatingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c8d6d6")).
			Italic(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b6b6b"))

	loreSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#f0d9b5"))

	svgSavedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#82b74b")).
			Bold(true)

	usageHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f0d9b5")).
			Italic(true)
)

// ── main ──────────────────────────────────────────────────────────────

func main() {
	svgPath := flag.String("svg", "", "export SVG to file path")
	width := flag.Int("w", 40, "map grid width")
	height := flag.Int("h", 25, "map grid height")
	svgWidth := flag.Int("svg-width", 1000, "SVG pixel width")
	svgHeight := flag.Int("svg-height", 700, "SVG pixel height")
	noLore := flag.Bool("no-lore", false, "skip lore output")
	minimal := flag.Bool("minimal", false, "show only the map, no legend/lore")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}
	placeName := strings.Join(args, " ")

	// ── header ────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println(indent(headerBoxStyle.Render("🗺️  IMAGINARY CARTOGRAPHER  🗺️")))
	fmt.Println()
	fmt.Println(indent(generatingStyle.Render(fmt.Sprintf("Generating map for: %q", placeName))))
	fmt.Println()

	// ── generate data ─────────────────────────────────────────────
	terrain := GenerateTerrain(placeName, *width, *height)

	lmGen := NewLandmarkGenerator(placeName)
	landmarks := lmGen.PlaceLandmarks(terrain.Width, terrain.Height, terrain.IsLand)

	loreGen := NewLoreGenerator(placeName)
	lore := loreGen.Generate(placeName)

	// ── terminal rendering ────────────────────────────────────────
	renderer := NewTermMapRenderer()

	fmt.Println(indent(renderer.RenderCompass()))
	fmt.Println()
	fmt.Println(indent(renderer.RenderMap(terrain, landmarks)))
	fmt.Println()

	if !*minimal {
		if legend := renderer.RenderLegend(landmarks); legend != "" {
			fmt.Println(indent(legend))
			fmt.Println()
		}

		fmt.Println(indent(renderSeparator(56)))
		fmt.Println()

		if !*noLore {
			title := loreSectionStyle.Render(
				fmt.Sprintf("📜 THE LORE OF %s", strings.ToUpper(placeName)),
			)
			fmt.Println(indent(title))
			fmt.Println()
			fmt.Println(indent(renderer.RenderLore(lore, placeName)))
			fmt.Println()
			fmt.Println(indent(renderSeparator(56)))
		}
	}

	// ── SVG export ────────────────────────────────────────────────
	if *svgPath != "" {
		svg := NewSVGMapRenderer(*svgWidth, *svgHeight)
		if err := svg.RenderToFile(terrain, landmarks, placeName, *svgPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving SVG: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println(indent(svgSavedStyle.Render(fmt.Sprintf("🗺️  SVG saved to: %s", *svgPath))))
	}

	fmt.Println()
}

// ── helpers ───────────────────────────────────────────────────────────

func printUsage() {
	fmt.Println()
	fmt.Println(indent(headerBoxStyle.Render("🗺️  IMAGINARY CARTOGRAPHER  🗺️")))
	fmt.Println()
	fmt.Println(indent("Usage: imaginary-cartographer [flags] <place-name>"))
	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println(indent(usageHintStyle.Render("🗺️  Name a place, and I shall map it!")))
	fmt.Println()
}

func indent(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = "  " + l
	}
	return strings.Join(lines, "\n")
}

func renderSeparator(width int) string {
	return separatorStyle.Render(strings.Repeat("━", width))
}
