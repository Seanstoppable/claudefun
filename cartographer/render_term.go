package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TermMapRenderer renders fantasy maps and lore to the terminal using lipgloss.
type TermMapRenderer struct {
	Width int // override map width for display (0 = use map's width)
}

// NewTermMapRenderer creates a renderer with default settings.
func NewTermMapRenderer() *TermMapRenderer {
	return &TermMapRenderer{}
}

// ── styles ────────────────────────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#f0d9b5")).
			Background(lipgloss.Color("#2d2d2d")).
			Padding(0, 1)

	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b6b6b"))

	landmarkOverlayStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ffd700"))

	legendHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#c8d6d6")).
				Underline(true)

	legendSymbolStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ffd700"))

	legendNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f0d9b5"))

	legendTypeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b6b6b")).
			Italic(true)

	compassStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c8d6d6"))

	compassStarStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ffd700"))

	loreHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#f0d9b5")).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6b6b6b"))

	mottoStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#82b74b"))

	mythPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#2e86c1")).
			Padding(0, 1).
			Width(62)

	legendTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#e8e8e8"))

	legendCategoryStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#2e86c1")).
				Italic(true)

	warningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#e07060")).
			Padding(0, 1)

	warningLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ff6b6b")).
				Background(lipgloss.Color("#2d2d2d")).
				Padding(0, 1)
)

// ── RenderMap ─────────────────────────────────────────────────────────

// RenderMap renders the terrain grid with landmark overlays and a border.
func (r *TermMapRenderer) RenderMap(t *TerrainMap, landmarks []Landmark) string {
	w := t.Width
	if r.Width > 0 {
		w = r.Width
		if w > t.Width {
			w = t.Width
		}
	}
	h := t.Height

	// Build a lookup of landmark positions.
	lmAt := make(map[[2]int]Landmark)
	for _, lm := range landmarks {
		if lm.X >= 0 && lm.X < w && lm.Y >= 0 && lm.Y < h {
			lmAt[[2]int{lm.X, lm.Y}] = lm
		}
	}

	// Title bar — centered in the border width.
	placeName := t.Seed
	innerWidth := w // each cell is one rune wide
	title := titleStyle.Render(placeName)
	titleLen := lipgloss.Width(title)
	pad := 0
	if innerWidth > titleLen {
		pad = (innerWidth - titleLen) / 2
	}
	titleLine := strings.Repeat(" ", pad) + title

	// Top border.
	topBorder := borderStyle.Render("╭" + strings.Repeat("─", innerWidth) + "╮")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", innerWidth) + "╯")
	vbar := borderStyle.Render("│")

	var sb strings.Builder
	sb.WriteString(titleLine + "\n")
	sb.WriteString(topBorder + "\n")

	for y := 0; y < h; y++ {
		sb.WriteString(vbar)
		for x := 0; x < w; x++ {
			if lm, ok := lmAt[[2]int{x, y}]; ok {
				sb.WriteString(landmarkOverlayStyle.Render(lm.Symbol))
			} else {
				cell := t.Cells[y][x]
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(cell.Biome.Color()))
				sb.WriteString(style.Render(cell.Biome.Symbol()))
			}
		}
		sb.WriteString(vbar + "\n")
	}

	sb.WriteString(bottomBorder)
	return sb.String()
}

// ── RenderLegend ──────────────────────────────────────────────────────

// RenderLegend renders a two-column key of the landmarks.
func (r *TermMapRenderer) RenderLegend(landmarks []Landmark) string {
	if len(landmarks) == 0 {
		return ""
	}

	header := legendHeaderStyle.Render("Landmarks")

	// Build formatted entries.
	entries := make([]string, len(landmarks))
	for i, lm := range landmarks {
		sym := legendSymbolStyle.Render(lm.Symbol)
		name := legendNameStyle.Render(lm.Name)
		typ := legendTypeStyle.Render("(" + lm.Type.String() + ")")
		entries[i] = fmt.Sprintf("%s %s %s", sym, name, typ)
	}

	// Two-column layout if >= 4 entries.
	var body string
	if len(entries) >= 4 {
		mid := (len(entries) + 1) / 2
		col1 := entries[:mid]
		col2 := entries[mid:]

		// Find max visible width in col1 for alignment.
		maxW := 0
		for _, e := range col1 {
			w := lipgloss.Width(e)
			if w > maxW {
				maxW = w
			}
		}
		colWidth := maxW + 3

		var lines []string
		for i := 0; i < mid; i++ {
			left := col1[i]
			leftW := lipgloss.Width(left)
			padding := ""
			if colWidth > leftW {
				padding = strings.Repeat(" ", colWidth-leftW)
			}
			right := ""
			if i < len(col2) {
				right = col2[i]
			}
			lines = append(lines, left+padding+right)
		}
		body = strings.Join(lines, "\n")
	} else {
		body = strings.Join(entries, "\n")
	}

	return header + "\n" + body
}

// ── RenderCompass ─────────────────────────────────────────────────────

// RenderCompass renders a small compass rose.
func (r *TermMapRenderer) RenderCompass() string {
	n := compassStyle.Render("N")
	s := compassStyle.Render("S")
	w := compassStyle.Render("W")
	e := compassStyle.Render("E")
	star := compassStarStyle.Render("✦")

	return fmt.Sprintf("    %s\n  %s %s %s\n    %s", n, w, star, e, s)
}

// ── RenderLore ────────────────────────────────────────────────────────

// RenderLore formats the lore with styled headers, panels, and word-wrapped text.
func (r *TermMapRenderer) RenderLore(lore *Lore, placeName string) string {
	if lore == nil {
		return ""
	}

	var sections []string

	// Place name header.
	sections = append(sections, loreHeaderStyle.Render("⚑ "+placeName))

	// Motto in italic.
	sections = append(sections, mottoStyle.Render("\""+lore.Motto+"\""))

	// Creation myth in a bordered panel.
	wrapped := wordWrap(lore.CreationMyth, 58)
	mythLabel := legendTitleStyle.Render("Creation Myth")
	mythBody := mythPanelStyle.Render(wrapped)
	sections = append(sections, mythLabel+"\n"+mythBody)

	// Legends.
	for _, leg := range lore.Legends {
		title := legendTitleStyle.Render(leg.Title)
		cat := legendCategoryStyle.Render("[" + leg.Category + "]")
		story := wordWrap(leg.Story, 60)
		sections = append(sections, fmt.Sprintf("%s  %s\n%s", title, cat, story))
	}

	// Warning.
	label := warningLabelStyle.Render("⚠ WARNING")
	warnText := warningStyle.Render(wordWrap(lore.Warning, 58))
	sections = append(sections, label+"\n"+warnText)

	return strings.Join(sections, "\n\n")
}

// ── helpers ───────────────────────────────────────────────────────────

// wordWrap breaks text into lines of at most maxWidth characters, splitting on spaces.
func wordWrap(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	line := words[0]
	for _, w := range words[1:] {
		if len(line)+1+len(w) > maxWidth {
			lines = append(lines, line)
			line = w
		} else {
			line += " " + w
		}
	}
	lines = append(lines, line)
	return strings.Join(lines, "\n")
}
