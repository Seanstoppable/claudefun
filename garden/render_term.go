package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TermGardenRenderer renders a Garden as colored terminal output using lipgloss.
type TermGardenRenderer struct{}

// NewTermGardenRenderer creates a new terminal renderer.
func NewTermGardenRenderer() *TermGardenRenderer {
	return &TermGardenRenderer{}
}

// Render draws the garden grid as colored ASCII art inside a border.
func (r *TermGardenRenderer) Render(g *Garden) string {
	var sb strings.Builder

	// Build each row of the garden grid
	var rows []string
	for y := 0; y < g.Height; y++ {
		var row strings.Builder
		for x := 0; x < g.Width; x++ {
			plot := g.Plots[y][x]
			sym := plot.Element.Symbol()
			color := plot.Element.Color()
			if plot.Element == Grass {
				// For grass, show the plant symbol instead
				sym = plot.Plant.Symbol()
				color = plot.Plant.Color()
			}
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			row.WriteString(style.Render(sym))
			row.WriteString(" ")
		}
		rows = append(rows, row.String())
	}

	// Season + Weather header inside the box
	header := fmt.Sprintf("  %s · %s", g.Season, g.Weather)
	rows = append([]string{header, ""}, rows...)

	content := strings.Join(rows, "\n")

	// Build a bordered box with the title centered
	titleLine := fmt.Sprintf(" 🌱 %s 🌱 ", g.Title)
	boxWidth := g.Width*3 + 4
	if boxWidth < len(titleLine)+4 {
		boxWidth = len(titleLine) + 4
	}

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#228B22")).
		Padding(0, 1).
		Width(boxWidth)

	box := border.Render(content)

	// Center the title in the top border area
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#228B22"))

	sb.WriteString(titleStyle.Render(centerText(titleLine, boxWidth+2)))
	sb.WriteString("\n")
	sb.WriteString(box)

	return sb.String()
}

// RenderStats produces a summary panel with health bar, stats, and a whimsical report.
func (r *TermGardenRenderer) RenderStats(g *Garden, stats *CodebaseStats) string {
	var sb strings.Builder

	// Health bar
	healthPct := int(g.HealthScore)
	filled := healthPct / 10
	empty := 10 - filled
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	healthColor := "#228B22"
	if healthPct < 40 {
		healthColor = "#FF0000"
	} else if healthPct < 70 {
		healthColor = "#DAA520"
	}
	healthStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(healthColor))

	sb.WriteString(fmt.Sprintf("  🌱 Garden Health: %s %d/100\n", healthStyle.Render(bar), healthPct))
	sb.WriteString("\n")

	// Weather + Season
	sb.WriteString(fmt.Sprintf("  %s · %s\n", g.Weather, g.Season))
	sb.WriteString("\n")

	// Element stats
	treeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#006400"))
	flowerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))
	fenceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B4513"))
	bugStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	butterflyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB"))
	rockStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))

	sb.WriteString(fmt.Sprintf("  %s %d Trees · %s %d Flowers · %s %d Fences\n",
		treeStyle.Render("🌳"), g.Stats.Trees,
		flowerStyle.Render("✿"), g.Stats.Flowers,
		fenceStyle.Render("┃"), g.Stats.Fences,
	))
	sb.WriteString(fmt.Sprintf("  %s %d Bugs · %s %d Butterflies · %s %d Rocks\n",
		bugStyle.Render("🐛"), g.Stats.Bugs,
		butterflyStyle.Render("🦋"), g.Stats.Butterflies,
		rockStyle.Render("●"), g.Stats.Rocks,
	))
	sb.WriteString("\n")

	// Code summary
	langNames := make([]string, 0, len(stats.Languages))
	for _, l := range stats.Languages {
		langNames = append(langNames, l.Name)
	}
	langs := strings.Join(langNames, ", ")
	if langs == "" {
		langs = "none detected"
	}

	sb.WriteString(fmt.Sprintf("  📊 %s files · %s lines · %s\n",
		formatNumber(stats.TotalFiles),
		formatNumber(stats.TotalLines),
		langs,
	))
	sb.WriteString("\n")

	// Whimsical garden report
	report := generateReport(g, stats)
	reportStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#228B22"))
	sb.WriteString(fmt.Sprintf("  %s\n", reportStyle.Render(fmt.Sprintf(`"%s"`, report))))

	return sb.String()
}

// RenderLegend shows what each symbol means.
func (r *TermGardenRenderer) RenderLegend() string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#228B22"))
	sb.WriteString(titleStyle.Render("  📖 Garden Legend"))
	sb.WriteString("\n\n")

	type entry struct {
		symbol string
		color  string
		label  string
	}

	plants := []entry{
		{"🌱", "#90EE90", "Seedling (tiny functions, 1-5 lines)"},
		{"ψ", "#228B22", "Sprout (small functions, 6-15 lines)"},
		{"♣", "#228B22", "Bush (medium functions, 16-30 lines)"},
		{"♠", "#006400", "Tree (large functions, 31-60 lines)"},
		{"🌳", "#004D00", "Oak Tree (huge functions, 60+ lines)"},
		{"✿", "#FF69B4", "Flower (type/struct definition)"},
		{"~", "#90EE90", "Vine (import/dependency)"},
	}

	elements := []entry{
		{"·", "#228B22", "Grass (empty space)"},
		{"┃", "#8B4513", "Fence (tests)"},
		{"🐛", "#FF0000", "Bug (TODO/FIXME)"},
		{"🦋", "#87CEEB", "Butterfly (comments)"},
		{"⌇", "#DAA520", "Weed (dead code)"},
		{"●", "#808080", "Rock (complexity hotspot)"},
		{"♤", "#FF4500", "Mushroom (type definition)"},
		{"≈", "#4169E1", "Pond (breathing room)"},
		{"⚑", "#FF0000", "Gnome (healthy garden!)"},
		{"◈", "#FFD700", "Beehive (heavily imported)"},
		{"⌂", "#8B4513", "Birdhouse (well-documented)"},
		{"╂", "#8B4513", "Scarecrow (error handling)"},
	}

	sb.WriteString("  Plants:\n")
	for _, e := range plants {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(e.color))
		sb.WriteString(fmt.Sprintf("    %s  %s\n", style.Render(e.symbol), e.label))
	}

	sb.WriteString("\n  Elements:\n")
	for _, e := range elements {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(e.color))
		sb.WriteString(fmt.Sprintf("    %s  %s\n", style.Render(e.symbol), e.label))
	}

	return sb.String()
}

// generateReport creates a whimsical description of the garden state.
func generateReport(g *Garden, stats *CodebaseStats) string {
	healthPct := int(g.HealthScore)

	switch {
	case healthPct >= 90:
		return "Your garden is magnificent! The oaks stand tall, the fences are sturdy, butterflies dance among the flowers, and not a weed in sight."
	case healthPct >= 75:
		parts := []string{"Your garden is thriving!"}
		if g.Stats.Trees > 5 {
			parts = append(parts, "The oaks stand tall.")
		}
		if g.Stats.Fences > 3 {
			parts = append(parts, "The fences are sturdy.")
		}
		if g.Stats.Bugs > 0 {
			parts = append(parts, fmt.Sprintf("Only %d bug(s) lurk in the undergrowth.", g.Stats.Bugs))
		}
		if g.Stats.Butterflies > 5 {
			parts = append(parts, "Butterflies flutter through the comments.")
		}
		return strings.Join(parts, " ")
	case healthPct >= 50:
		parts := []string{"Your garden is growing nicely, but could use some tending."}
		if g.Stats.Bugs > 3 {
			parts = append(parts, "Watch out for the bugs!")
		}
		if g.Stats.Weeds > 0 {
			parts = append(parts, "Some weeds are creeping in.")
		}
		if g.Stats.Fences == 0 {
			parts = append(parts, "Consider planting some fences (tests).")
		}
		return strings.Join(parts, " ")
	case healthPct >= 30:
		return "Your garden needs attention. The weeds are spreading, bugs are multiplying, and the fences are thin. Time to refactor and add tests!"
	default:
		return "Your garden is a wilderness! Tumbleweeds roll through untested code. Rally the gardeners — it's time for a major cleanup."
	}
}

// formatNumber adds comma separators to integers.
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

// centerText centers text within the given width.
func centerText(text string, width int) string {
	textLen := len([]rune(text))
	if textLen >= width {
		return text
	}
	pad := (width - textLen) / 2
	return strings.Repeat(" ", pad) + text
}
