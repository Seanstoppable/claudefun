package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// SVGGardenRenderer renders a Garden as an SVG image.
type SVGGardenRenderer struct {
	Width, Height int
}

// NewSVGGardenRenderer creates a new SVG renderer with the given dimensions.
func NewSVGGardenRenderer(w, h int) *SVGGardenRenderer {
	return &SVGGardenRenderer{Width: w, Height: h}
}

// Render produces a complete SVG string representing the garden.
func (r *SVGGardenRenderer) Render(g *Garden) string {
	var sb strings.Builder

	w := r.Width
	h := r.Height

	// Compute cell sizing
	skyHeight := h / 5
	groundTop := skyHeight
	legendHeight := h / 8
	gardenHeight := h - skyHeight - legendHeight
	cellW := float64(w) / float64(g.Width+2)
	cellH := float64(gardenHeight) / float64(g.Height+2)
	if cellW > cellH {
		cellW = cellH
	}

	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, w, h, w, h))
	sb.WriteString("\n")

	// Definitions: gradients
	sb.WriteString(`<defs>`)
	sb.WriteString(`<linearGradient id="sky" x1="0" y1="0" x2="0" y2="1">`)
	sb.WriteString(`<stop offset="0%" stop-color="#87CEEB"/>`)
	sb.WriteString(`<stop offset="100%" stop-color="#E0F0FF"/>`)
	sb.WriteString(`</linearGradient>`)
	sb.WriteString(`<linearGradient id="grass-bg" x1="0" y1="0" x2="0" y2="1">`)
	sb.WriteString(`<stop offset="0%" stop-color="#4CAF50"/>`)
	sb.WriteString(`<stop offset="100%" stop-color="#2E7D32"/>`)
	sb.WriteString(`</linearGradient>`)
	sb.WriteString(`</defs>`)
	sb.WriteString("\n")

	// Sky background
	sb.WriteString(fmt.Sprintf(`<rect x="0" y="0" width="%d" height="%d" fill="url(#sky)"/>`, w, skyHeight))
	sb.WriteString("\n")

	// Weather elements in the sky
	r.renderWeather(&sb, g.Weather, w, skyHeight)

	// Ground background
	sb.WriteString(fmt.Sprintf(`<rect x="0" y="%d" width="%d" height="%d" fill="url(#grass-bg)"/>`, groundTop, w, gardenHeight+legendHeight))
	sb.WriteString("\n")

	// Grass texture dots
	for i := 0; i < 40; i++ {
		gx := (i*97 + 13) % w
		gy := groundTop + (i*53+7)%gardenHeight
		sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="1.5" fill="#388E3C" opacity="0.3"/>`, gx, gy))
		sb.WriteString("\n")
	}

	// Title
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="30" text-anchor="middle" font-family="sans-serif" font-size="20" font-weight="bold" fill="#1B5E20">🌱 %s 🌱</text>`, w/2, escapeXML(g.Title)))
	sb.WriteString("\n")

	// Season + Weather text
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="55" text-anchor="middle" font-family="sans-serif" font-size="12" fill="#555">%s · %s</text>`, w/2, escapeXML(g.Season), escapeXML(g.Weather)))
	sb.WriteString("\n")

	// Garden grid offset to center it
	offsetX := (float64(w) - cellW*float64(g.Width)) / 2
	offsetY := float64(groundTop) + 10

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			plot := g.Plots[y][x]
			cx := offsetX + float64(x)*cellW + cellW/2
			cy := offsetY + float64(y)*cellH + cellH/2

			r.renderPlot(&sb, plot, cx, cy, cellW, cellH)
		}
	}

	// Stats/legend at bottom
	legendY := h - legendHeight + 20
	r.renderLegend(&sb, g, w, legendY)

	sb.WriteString("</svg>\n")
	return sb.String()
}

// RenderToFile writes the SVG to a file.
func (r *SVGGardenRenderer) RenderToFile(g *Garden, path string) error {
	svg := r.Render(g)
	return os.WriteFile(path, []byte(svg), 0644)
}

func (r *SVGGardenRenderer) renderWeather(sb *strings.Builder, weather string, w, skyH int) {
	switch {
	case strings.Contains(weather, "Sunny"), strings.Contains(weather, "Rainbow"):
		// Sun
		sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="50" r="25" fill="#FFD700"/>`, w-60))
		for i := 0; i < 8; i++ {
			angle := float64(i) * 45
			x1 := float64(w-60) + 30*cos(angle)
			y1 := 50 + 30*sin(angle)
			x2 := float64(w-60) + 40*cos(angle)
			y2 := 50 + 40*sin(angle)
			sb.WriteString(fmt.Sprintf(`<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#FFD700" stroke-width="2"/>`, x1, y1, x2, y2))
		}
		sb.WriteString("\n")
		if strings.Contains(weather, "Rainbow") {
			colors := []string{"#FF0000", "#FF7F00", "#FFFF00", "#00FF00", "#0000FF", "#8B00FF"}
			for i, c := range colors {
				radius := 80 + i*8
				sb.WriteString(fmt.Sprintf(`<path d="M %d %d A %d %d 0 0 1 %d %d" fill="none" stroke="%s" stroke-width="4" opacity="0.6"/>`,
					w/2-radius, skyH-10, radius, radius, w/2+radius, skyH-10, c))
			}
			sb.WriteString("\n")
		}
	case strings.Contains(weather, "Cloudy"), strings.Contains(weather, "Partly"):
		// Clouds
		for i := 0; i < 3; i++ {
			cx := 100 + i*200
			cy := 40 + (i%2)*20
			sb.WriteString(fmt.Sprintf(`<ellipse cx="%d" cy="%d" rx="40" ry="20" fill="white" opacity="0.8"/>`, cx, cy))
			sb.WriteString(fmt.Sprintf(`<ellipse cx="%d" cy="%d" rx="30" ry="18" fill="white" opacity="0.8"/>`, cx-25, cy+5))
			sb.WriteString(fmt.Sprintf(`<ellipse cx="%d" cy="%d" rx="30" ry="18" fill="white" opacity="0.8"/>`, cx+25, cy+5))
		}
		if strings.Contains(weather, "Partly") {
			sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="45" r="20" fill="#FFD700"/>`, w-80))
		}
		sb.WriteString("\n")
	case strings.Contains(weather, "Rainy"), strings.Contains(weather, "Stormy"):
		// Dark clouds + rain drops
		for i := 0; i < 3; i++ {
			cx := 80 + i*220
			cy := 35 + (i%2)*15
			sb.WriteString(fmt.Sprintf(`<ellipse cx="%d" cy="%d" rx="50" ry="22" fill="#9E9E9E" opacity="0.9"/>`, cx, cy))
			sb.WriteString(fmt.Sprintf(`<ellipse cx="%d" cy="%d" rx="35" ry="18" fill="#9E9E9E" opacity="0.9"/>`, cx-30, cy+5))
			sb.WriteString(fmt.Sprintf(`<ellipse cx="%d" cy="%d" rx="35" ry="18" fill="#9E9E9E" opacity="0.9"/>`, cx+30, cy+5))
			// Rain drops
			for d := 0; d < 5; d++ {
				dx := cx - 30 + d*15
				dy := cy + 25 + d*8
				sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#4FC3F7" stroke-width="1.5" opacity="0.7"/>`, dx, dy, dx-3, dy+10))
			}
		}
		sb.WriteString("\n")
	}
}

func (r *SVGGardenRenderer) renderPlot(sb *strings.Builder, plot GardenPlot, cx, cy, cellW, cellH float64) {
	radius := cellW * 0.35
	if radius < 3 {
		radius = 3
	}

	switch plot.Element {
	case Grass:
		r.renderPlant(sb, plot.Plant, cx, cy, radius)
	case Fence:
		// Brown horizontal line
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="2"/>`,
			cx-radius, cy, cx+radius, cy))
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="1.5"/>`,
			cx-radius*0.3, cy-radius*0.6, cx-radius*0.3, cy+radius*0.3))
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="1.5"/>`,
			cx+radius*0.3, cy-radius*0.6, cx+radius*0.3, cy+radius*0.3))
	case Bug:
		// Small red body with legs
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#FF0000"/>`, cx, cy, radius*0.5, radius*0.35))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#CC0000"/>`, cx-radius*0.3, cy, radius*0.2))
		for i := -1; i <= 1; i++ {
			lx := cx + float64(i)*radius*0.3
			sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#333" stroke-width="0.5"/>`,
				lx, cy-radius*0.35, lx-radius*0.15, cy-radius*0.6))
			sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#333" stroke-width="0.5"/>`,
				lx, cy+radius*0.35, lx-radius*0.15, cy+radius*0.6))
		}
	case Butterfly:
		// Blue wing shapes
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#87CEEB" opacity="0.8"/>`, cx-radius*0.3, cy, radius*0.4, radius*0.25))
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#87CEEB" opacity="0.8"/>`, cx+radius*0.3, cy, radius*0.4, radius*0.25))
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#555" stroke-width="0.7"/>`, cx, cy-radius*0.15, cx, cy+radius*0.15))
	case Weed:
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#DAA520" stroke-width="1"/>`, cx, cy+radius*0.5, cx, cy-radius*0.5))
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#DAA520" stroke-width="1"/>`, cx, cy-radius*0.2, cx-radius*0.4, cy-radius*0.5))
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#DAA520" stroke-width="1"/>`, cx, cy, cx+radius*0.4, cy-radius*0.3))
	case Rock:
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#808080"/>`, cx, cy, radius*0.5, radius*0.35))
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#999" opacity="0.5"/>`, cx-radius*0.1, cy-radius*0.1, radius*0.3, radius*0.2))
	case Mushroom:
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#DEB887" stroke-width="2"/>`, cx, cy, cx, cy+radius*0.5))
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#FF4500"/>`, cx, cy-radius*0.1, radius*0.45, radius*0.3))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="white" opacity="0.7"/>`, cx-radius*0.15, cy-radius*0.2, radius*0.08))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="white" opacity="0.7"/>`, cx+radius*0.1, cy-radius*0.05, radius*0.06))
	case Pond:
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#4169E1" opacity="0.7"/>`, cx, cy, radius*0.6, radius*0.35))
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#5B9BD5" opacity="0.4"/>`, cx-radius*0.15, cy-radius*0.05, radius*0.25, radius*0.12))
	case Gnome:
		// Red hat triangle + body
		sb.WriteString(fmt.Sprintf(`<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="#FF0000"/>`,
			cx, cy-radius*0.7, cx-radius*0.3, cy-radius*0.1, cx+radius*0.3, cy-radius*0.1))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#FFDAB9"/>`, cx, cy+radius*0.1, radius*0.25))
	case Beehive:
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#FFD700"/>`, cx, cy, radius*0.4, radius*0.5))
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#B8860B" stroke-width="0.7"/>`, cx-radius*0.35, cy-radius*0.15, cx+radius*0.35, cy-radius*0.15))
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#B8860B" stroke-width="0.7"/>`, cx-radius*0.35, cy+radius*0.15, cx+radius*0.35, cy+radius*0.15))
	case Birdhouse:
		sb.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#8B4513"/>`, cx-radius*0.3, cy-radius*0.2, radius*0.6, radius*0.6))
		sb.WriteString(fmt.Sprintf(`<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="#A0522D"/>`,
			cx-radius*0.45, cy-radius*0.2, cx, cy-radius*0.6, cx+radius*0.45, cy-radius*0.2))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#333"/>`, cx, cy+radius*0.05, radius*0.1))
	case Scarecrow:
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="1.5"/>`, cx, cy-radius*0.5, cx, cy+radius*0.5))
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="1.5"/>`, cx-radius*0.5, cy-radius*0.1, cx+radius*0.5, cy-radius*0.1))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#DEB887"/>`, cx, cy-radius*0.4, radius*0.2))
	case Tumbleweed:
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke="#8B4513" stroke-width="1"/>`, cx, cy, radius*0.4))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke="#8B4513" stroke-width="0.5" opacity="0.6"/>`, cx, cy, radius*0.25))
	default:
		// Fallback: small dot
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#228B22" opacity="0.3"/>`, cx, cy, radius*0.2))
	}
}

func (r *SVGGardenRenderer) renderPlant(sb *strings.Builder, plant PlantType, cx, cy, radius float64) {
	switch plant {
	case Seedling:
		// Small green dot
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#90EE90"/>`, cx, cy, radius*0.2))
	case Sprout:
		// Small stalk with a leaf
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#228B22" stroke-width="1.5"/>`, cx, cy+radius*0.4, cx, cy-radius*0.3))
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#228B22"/>`, cx+radius*0.15, cy-radius*0.15, radius*0.15, radius*0.08))
	case Bush:
		// Rounded green blob
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#228B22"/>`, cx, cy, radius*0.45, radius*0.35))
		sb.WriteString(fmt.Sprintf(`<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="#2E8B2E" opacity="0.7"/>`, cx-radius*0.15, cy-radius*0.08, radius*0.3, radius*0.22))
	case Tree:
		// Brown trunk + green circle canopy
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="2"/>`, cx, cy+radius*0.5, cx, cy-radius*0.1))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#006400"/>`, cx, cy-radius*0.3, radius*0.35))
	case OakTree:
		// Bigger trunk + larger canopy
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#8B4513" stroke-width="3"/>`, cx, cy+radius*0.5, cx, cy-radius*0.2))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#004D00"/>`, cx, cy-radius*0.35, radius*0.5))
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#006400" opacity="0.6"/>`, cx-radius*0.2, cy-radius*0.25, radius*0.3))
	case Flower:
		// Colored petals around a center
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#228B22" stroke-width="1"/>`, cx, cy+radius*0.4, cx, cy))
		for i := 0; i < 5; i++ {
			angle := float64(i) * 72
			px := cx + radius*0.2*cos(angle)
			py := cy - radius*0.15 + radius*0.2*sin(angle)
			sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#FF69B4" opacity="0.8"/>`, px, py, radius*0.1))
		}
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#FFD700"/>`, cx, cy-radius*0.15, radius*0.08))
	case Vine:
		// Wavy line
		sb.WriteString(fmt.Sprintf(`<path d="M %.1f %.1f Q %.1f %.1f %.1f %.1f Q %.1f %.1f %.1f %.1f" fill="none" stroke="#90EE90" stroke-width="1"/>`,
			cx-radius*0.4, cy,
			cx-radius*0.2, cy-radius*0.3,
			cx, cy,
			cx+radius*0.2, cy+radius*0.3,
			cx+radius*0.4, cy))
	}
}

func (r *SVGGardenRenderer) renderLegend(sb *strings.Builder, g *Garden, w, y int) {
	sb.WriteString(fmt.Sprintf(`<rect x="10" y="%d" width="%d" height="55" rx="8" fill="white" opacity="0.85"/>`, y-15, w-20))

	healthPct := int(g.HealthScore)
	sb.WriteString(fmt.Sprintf(`<text x="20" y="%d" font-family="sans-serif" font-size="11" fill="#333">🌱 Health: %d/100  ·  🌳 %d Trees  ·  ✿ %d Flowers  ·  ┃ %d Fences  ·  🐛 %d Bugs  ·  🦋 %d Butterflies  ·  ● %d Rocks</text>`,
		y+2, healthPct,
		g.Stats.Trees, g.Stats.Flowers, g.Stats.Fences,
		g.Stats.Bugs, g.Stats.Butterflies, g.Stats.Rocks))

	// Health bar
	barX := 20
	barY := y + 12
	barW := w - 40
	filledW := barW * healthPct / 100

	barColor := "#4CAF50"
	if healthPct < 40 {
		barColor = "#F44336"
	} else if healthPct < 70 {
		barColor = "#FF9800"
	}
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="8" rx="4" fill="#E0E0E0"/>`, barX, barY, barW))
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="8" rx="4" fill="%s"/>`, barX, barY, filledW, barColor))

	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-family="sans-serif" font-size="9" fill="#666">🌱 = seedling · ψ = sprout · ♣ = bush · ♠ = tree · 🌳 = oak · ✿ = flower · ┃ = fence(test) · 🐛 = bug(TODO) · 🦋 = butterfly(comment)</text>`,
		w/2, y+35))
}

// Trig helpers that take degrees.
func cos(deg float64) float64 { return math.Cos(deg * math.Pi / 180) }
func sin(deg float64) float64 { return math.Sin(deg * math.Pi / 180) }

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
