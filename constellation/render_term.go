package main

import (
	"hash/fnv"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TerminalRenderer converts a Constellation into colored Unicode art.
type TerminalRenderer struct {
	Width  int // character columns (default 60)
	Height int // character rows   (default 30)
}

// NewTerminalRenderer creates a renderer with the given grid dimensions.
func NewTerminalRenderer(width, height int) *TerminalRenderer {
	if width <= 0 {
		width = 60
	}
	if height <= 0 {
		height = 30
	}
	return &TerminalRenderer{Width: width, Height: height}
}

// cell holds the content and style for one grid position.
type cell struct {
	ch    string
	style lipgloss.Style
	layer int // 0 = background, 1 = edge, 2 = star
}

// ── styles ─────────────────────────────────────────────────────────────

var (
	styleBrightStar = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")) // gold/yellow
	styleMedStar    = lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB")) // light blue
	styleDimStar    = lipgloss.NewStyle().Foreground(lipgloss.Color("#5B7EA0")) // dim blue
	styleFaintStar  = lipgloss.NewStyle().Foreground(lipgloss.Color("#3A3A4A")) // dark gray

	styleEdge1 = lipgloss.NewStyle().Foreground(lipgloss.Color("#6A5ACD")) // slate blue
	styleEdge2 = lipgloss.NewStyle().Foreground(lipgloss.Color("#7B68EE")) // medium slate blue
	styleEdge3 = lipgloss.NewStyle().Foreground(lipgloss.Color("#9370DB")) // medium purple

	styleBgDot = lipgloss.NewStyle().Foreground(lipgloss.Color("#2A2A35")) // very dim
)

// Render produces a colored Unicode string of the constellation.
func (r *TerminalRenderer) Render(c *Constellation) string {
	w, h := r.Width, r.Height
	grid := make([][]cell, h)
	for y := range grid {
		grid[y] = make([]cell, w)
		for x := range grid[y] {
			grid[y][x] = cell{ch: " ", style: lipgloss.NewStyle(), layer: 0}
		}
	}

	// 1. Background stars — deterministic ~5 % density.
	seedBackground(grid, w, h, c.Stars.Seed)

	// 2. Edges (Bresenham lines).
	for _, e := range c.Edges {
		if e.From < 0 || e.From >= len(c.Stars.Stars) ||
			e.To < 0 || e.To >= len(c.Stars.Stars) {
			continue
		}
		a := c.Stars.Stars[e.From]
		b := c.Stars.Stars[e.To]
		drawLine(grid, w, h, a.X, a.Y, b.X, b.Y, e.From+e.To)
	}

	// 3. Stars — drawn last so they sit on top.
	for _, s := range c.Stars.Stars {
		gx := clampInt(int(math.Round(s.X*float64(w-1))), 0, w-1)
		gy := clampInt(int(math.Round(s.Y*float64(h-1))), 0, h-1)
		ch, sty := starGlyph(s.Brightness)
		grid[gy][gx] = cell{ch: ch, style: sty, layer: 2}
	}

	// 4. Compose output.
	var buf strings.Builder
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := grid[y][x]
			buf.WriteString(c.style.Render(c.ch))
		}
		if y < h-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

// ── background ─────────────────────────────────────────────────────────

func seedBackground(grid [][]cell, w, h int, seed string) {
	hasher := fnv.New64a()
	hasher.Write([]byte(seed))
	base := hasher.Sum64()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Simple deterministic hash per cell.
			v := base ^ uint64(y*7919+x*6271)
			v ^= v >> 17
			v *= 0xbf58476d1ce4e5b9
			v ^= v >> 31
			if v%20 == 0 { // ~5 % density
				grid[y][x] = cell{ch: "·", style: styleBgDot, layer: 0}
			}
		}
	}
}

// ── star glyphs ────────────────────────────────────────────────────────

func starGlyph(brightness float64) (string, lipgloss.Style) {
	switch {
	case brightness > 0.8:
		return "✦", styleBrightStar
	case brightness > 0.5:
		return "★", styleMedStar
	case brightness > 0.3:
		return "✧", styleDimStar
	default:
		return "·", styleFaintStar
	}
}

// ── Bresenham line drawing ─────────────────────────────────────────────

func drawLine(grid [][]cell, w, h int, x0, y0, x1, y1 float64, seed int) {
	gx0 := clampInt(int(math.Round(x0*float64(w-1))), 0, w-1)
	gy0 := clampInt(int(math.Round(y0*float64(h-1))), 0, h-1)
	gx1 := clampInt(int(math.Round(x1*float64(w-1))), 0, w-1)
	gy1 := clampInt(int(math.Round(y1*float64(h-1))), 0, h-1)

	dx := absInt(gx1 - gx0)
	dy := absInt(gy1 - gy0)
	sx := 1
	if gx0 > gx1 {
		sx = -1
	}
	sy := 1
	if gy0 > gy1 {
		sy = -1
	}

	// Direction from start to end for choosing line chars.
	dirX := gx1 - gx0
	dirY := gy1 - gy0

	err := dx - dy
	cx, cy := gx0, gy0

	steps := dx + dy
	if steps == 0 {
		steps = 1
	}

	for i := 0; ; i++ {
		if grid[cy][cx].layer < 1 {
			ch := lineChar(dirX, dirY, cx, cy, gx0, gy0, gx1, gy1)
			sty := edgeStyle(seed, i, steps)
			grid[cy][cx] = cell{ch: ch, style: sty, layer: 1}
		}

		if cx == gx1 && cy == gy1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			cx += sx
		}
		if e2 < dx {
			err += dx
			cy += sy
		}
	}
}

// lineChar picks the correct line-drawing rune based on local direction.
func lineChar(dirX, dirY, cx, cy, x0, y0, x1, y1 int) string {
	adx := absInt(dirX)
	ady := absInt(dirY)

	if ady == 0 {
		return "─"
	}
	if adx == 0 {
		return "│"
	}

	// Ratio determines dominant axis.
	ratio := float64(adx) / float64(ady)

	switch {
	case ratio > 2.0:
		return "─"
	case ratio < 0.5:
		return "│"
	default:
		// Diagonal — choose slash direction.
		// Positive slope in screen coords (down-right) ⇒ ╲
		// Negative slope in screen coords (up-right)   ⇒ ╱
		if (dirX > 0 && dirY > 0) || (dirX < 0 && dirY < 0) {
			return "╲"
		}
		return "╱"
	}
}

// edgeStyle picks a subtle gradient color for edge segments.
func edgeStyle(seed, step, totalSteps int) lipgloss.Style {
	t := float64(step) / float64(totalSteps)
	// Mix based on position along the edge and a per-edge seed.
	idx := int(math.Round(t*2.0+float64(seed%3))) % 3
	switch idx {
	case 0:
		return styleEdge1
	case 1:
		return styleEdge2
	default:
		return styleEdge3
	}
}

// ── integer helpers ────────────────────────────────────────────────────

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
