package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ─── colour palette ────────────────────────────────────────
var (
	gold    = lipgloss.Color("#FFD700")
	amber   = lipgloss.Color("#FFBF00")
	seafoam = lipgloss.Color("#20B2AA")
	coral   = lipgloss.Color("#FF6F61")
	silver  = lipgloss.Color("#C0C0C0")
	dimGray = lipgloss.Color("#696969")
	moonlit = lipgloss.Color("#B0C4DE")
	crimson = lipgloss.Color("#DC143C")
)

// ─── lipgloss styles ───────────────────────────────────────
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(gold).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(amber).
			Align(lipgloss.Center)

	tempoStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(dimGray)

	sectionLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(seafoam)

	verseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	chorusStyle = lipgloss.NewStyle().
			Foreground(coral)

	bridgeStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(moonlit)

	codaStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(silver)

	bandStyle = lipgloss.NewStyle().
			Foreground(dimGray)

	headerBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(gold).
			Padding(1, 2).
			Align(lipgloss.Center)

	songTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(gold)
)

// ─── ASCII art by shanty type ──────────────────────────────

const shipArt = `        |    |    |
       )_)  )_)  )_)
      )___))___))___)\
     )____)____)_____)\\
   _____|____|____|____\\\__
---\                   /----
    \_________________/
      ~~~~~~~~~~~~~~~`

const anchorArt = `    ⚓
   /|\
  / | \
    |
   /|\
  /_|_\`

const skullArt = `    _____
   /     \
  | () () |
  |   ^   |
  |  \_/  |
   \_____/
   /|   |\`

const moonArt = `     _..._
   .:::::::.
  :::::::::::
  '::::::::'
    ':::::'
      ':` + "`"

const starBurst = `         .  *  .    .   *    .  *  .
      *    \  |  /    .    *    .
          -- ☆ --
      *    /  |  \    .    *    .
         .  *  .    .   *    .  *  .`

// ─── ShantyRenderer ────────────────────────────────────────

// ShantyRenderer formats a Shanty for beautiful terminal output.
type ShantyRenderer struct {
	Width int // terminal width (default 60)
}

// NewShantyRenderer returns a renderer with the given width.
// If width <= 0 it defaults to 60.
func NewShantyRenderer(width int) *ShantyRenderer {
	if width <= 0 {
		width = 60
	}
	return &ShantyRenderer{Width: width}
}

// artForType picks ASCII art that matches the shanty mood.
func artForType(t ShantyType) string {
	switch t {
	case CelebrationJig, EpicSaga:
		return shipArt
	case MournfulBallad, WorkSong:
		return anchorArt
	case MutinyAnthem:
		return skullArt
	case LullabyCalm:
		return moonArt
	default:
		return shipArt
	}
}

// ─── public methods ────────────────────────────────────────

// RenderHeader produces the main "GIT SHANTY" banner with star art.
func (r *ShantyRenderer) RenderHeader() string {
	banner := titleStyle.Width(r.Width).Render("🏴\u200d☠️  GIT SHANTY  🏴\u200d☠️")
	tagline := subtitleStyle.Width(r.Width).Render(`"Your git history, in song"`)

	box := headerBoxStyle.Width(r.Width).Render(
		lipgloss.JoinVertical(lipgloss.Center, banner, tagline),
	)

	stars := lipgloss.NewStyle().
		Foreground(amber).
		Align(lipgloss.Center).
		Width(r.Width).
		Render(starBurst)

	return lipgloss.JoinVertical(lipgloss.Center, box, "", stars)
}

// RenderDivider returns a nautical-themed divider that spans the width.
func (r *ShantyRenderer) RenderDivider() string {
	// ┄┄┄┄┄┄┄┄ ⚓ ┄┄┄┄┄┄┄┄┄┄┄ ⚓ ┄┄┄┄┄┄┄┄
	seg := (r.Width - 6) / 2 // space for two anchors + padding
	if seg < 4 {
		seg = 4
	}
	left := strings.Repeat("┄", seg)
	right := strings.Repeat("┄", seg)
	line := fmt.Sprintf("%s ⚓ %s", left, right)

	return lipgloss.NewStyle().
		Foreground(seafoam).
		Width(r.Width).
		Align(lipgloss.Center).
		Render(line)
}

// RenderShipArt returns the ASCII art decoration appropriate for the
// given shanty type. When called on the renderer alone (no Shanty
// context), it defaults to the ship.
func (r *ShantyRenderer) RenderShipArt() string {
	return r.renderArt(CelebrationJig)
}

func (r *ShantyRenderer) renderArt(t ShantyType) string {
	art := artForType(t)
	return lipgloss.NewStyle().
		Foreground(amber).
		Align(lipgloss.Center).
		Width(r.Width).
		Render(art)
}

// Render produces the complete, decorated terminal output for a Shanty.
func (r *ShantyRenderer) Render(s Shanty) string {
	var b strings.Builder

	// Header
	b.WriteString(r.RenderHeader())
	b.WriteString("\n\n")

	// Type-specific art
	b.WriteString(r.renderArt(s.Type))
	b.WriteString("\n\n")

	// Divider
	b.WriteString(r.RenderDivider())
	b.WriteString("\n\n")

	// Song title
	b.WriteString(songTitleStyle.Render(fmt.Sprintf("  🎵  %q", s.Title)))
	b.WriteString("\n\n")

	// Tempo marking
	if s.Tempo != "" {
		b.WriteString(tempoStyle.Render(fmt.Sprintf("  ♩ = %s", s.Tempo)))
		b.WriteString("\n\n")
	}

	// Thin divider before lyrics
	b.WriteString(r.thinDivider())
	b.WriteString("\n\n")

	// Verses
	for i, v := range s.Verses {
		label := sectionLabelStyle.Render(fmt.Sprintf("  [Verse %d]", i+1))
		b.WriteString(label)
		b.WriteString("\n")
		b.WriteString(verseStyle.Render(indent(v, 2)))
		b.WriteString("\n\n")

		// Chorus after every verse (if present)
		if s.Chorus != "" {
			b.WriteString(sectionLabelStyle.Render("  [Chorus]"))
			b.WriteString("\n")
			b.WriteString(chorusStyle.Render(indent(s.Chorus, 2)))
			b.WriteString("\n\n")
		}
	}

	// Bridge
	if s.Bridge != "" {
		b.WriteString(sectionLabelStyle.Render("  [Bridge]"))
		b.WriteString("\n")
		b.WriteString(bridgeStyle.Render(indent(s.Bridge, 2)))
		b.WriteString("\n\n")
	}

	// Coda
	if s.Coda != "" {
		b.WriteString(sectionLabelStyle.Render("  [Coda]"))
		b.WriteString("\n")
		b.WriteString(codaStyle.Render(indent(s.Coda, 2)))
		b.WriteString("\n\n")
	}

	// Divider after lyrics
	b.WriteString(r.thinDivider())
	b.WriteString("\n\n")

	// Band name
	if s.BandName != "" {
		b.WriteString(bandStyle.Render(fmt.Sprintf("  Performed by: %s", s.BandName)))
		b.WriteString("\n\n")
	}

	// Footer border
	b.WriteString(r.footerBorder())
	b.WriteString("\n")

	return b.String()
}

// ─── helpers ───────────────────────────────────────────────

func (r *ShantyRenderer) thinDivider() string {
	line := strings.Repeat("┄", r.Width-4)
	return lipgloss.NewStyle().
		Foreground(dimGray).
		Render("  " + line)
}

func (r *ShantyRenderer) footerBorder() string {
	inner := strings.Repeat("─", r.Width-4)
	line := fmt.Sprintf("  ⚓%s⚓", inner)
	return lipgloss.NewStyle().
		Foreground(gold).
		Render(line)
}

func indent(text string, spaces int) string {
	pad := strings.Repeat(" ", spaces)
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = pad + l
	}
	return strings.Join(lines, "\n")
}
