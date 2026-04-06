package main

import (
	"fmt"
	"os"
	"strings"
)

// SVGMapRenderer produces fantasy-style SVG maps from terrain and landmark data.
type SVGMapRenderer struct {
	Width, Height int // SVG pixel dimensions
	CellSize      int // computed; not set directly
}

// NewSVGMapRenderer creates a renderer with the given SVG pixel dimensions.
// Defaults to 800×600 if either dimension is <= 0.
func NewSVGMapRenderer(width, height int) *SVGMapRenderer {
	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 600
	}
	return &SVGMapRenderer{Width: width, Height: height}
}

// layout constants (fractions of total SVG size)
const (
	titleFraction  = 0.08 // top 8% reserved for title cartouche
	legendFraction = 0.12 // bottom 12% reserved for legend
	marginFraction = 0.03 // left/right margin
)

// Render produces a complete SVG string for the given terrain, landmarks, and place name.
func (r *SVGMapRenderer) Render(t *TerrainMap, landmarks []Landmark, placeName string) string {
	var b strings.Builder
	b.Grow(64 * 1024)

	// Map area geometry
	marginX := int(float64(r.Width) * marginFraction)
	titleH := int(float64(r.Height) * titleFraction)
	legendH := int(float64(r.Height) * legendFraction)
	mapX := marginX
	mapY := titleH + 10
	mapW := r.Width - 2*marginX
	mapH := r.Height - titleH - legendH - 20

	cellW := float64(mapW) / float64(t.Width)
	cellH := float64(mapH) / float64(t.Height)

	// SVG header
	fmt.Fprintf(&b, `<?xml version="1.0" encoding="UTF-8"?>
<!-- Imaginary Cartographer: %s -->
<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
`, xmlEscape(placeName), r.Width, r.Height, r.Width, r.Height)

	r.writeDefs(&b, mapX, mapY, mapW, mapH)

	// Parchment background
	fmt.Fprintf(&b, `<rect width="%d" height="%d" fill="#f4e8c1"/>
`, r.Width, r.Height)

	// Aged vignette overlay
	fmt.Fprintf(&b, `<rect width="%d" height="%d" fill="url(#vignette)" opacity="0.4"/>
`, r.Width, r.Height)

	// Decorative double-line border
	r.writeBorder(&b)

	// Title cartouche
	r.writeTitle(&b, placeName, titleH)

	// Terrain cells
	r.writeTerrain(&b, t, mapX, mapY, cellW, cellH)

	// Grid lines (very subtle)
	r.writeGrid(&b, t, mapX, mapY, mapW, mapH, cellW, cellH)

	// Landmarks
	r.writeLandmarks(&b, landmarks, mapX, mapY, cellW, cellH)

	// Compass rose
	r.writeCompass(&b, mapX+mapW-60, mapY+mapH-60)

	// Legend box
	r.writeLegend(&b, t, landmarks, marginX, r.Height-legendH-5, r.Width-2*marginX, legendH)

	b.WriteString("</svg>\n")
	return b.String()
}

// RenderToFile writes the SVG to the given file path.
func (r *SVGMapRenderer) RenderToFile(t *TerrainMap, landmarks []Landmark, placeName string, path string) error {
	svg := r.Render(t, landmarks, placeName)
	return os.WriteFile(path, []byte(svg), 0644)
}

// --- internal rendering helpers ---

func (r *SVGMapRenderer) writeDefs(b *strings.Builder, mapX, mapY, mapW, mapH int) {
	b.WriteString("<defs>\n")

	// Radial gradient for vignette / aged edges
	b.WriteString(`<radialGradient id="vignette" cx="50%" cy="50%" r="70%">
  <stop offset="0%" stop-color="#f4e8c1" stop-opacity="0"/>
  <stop offset="80%" stop-color="#c4a66a" stop-opacity="0.15"/>
  <stop offset="100%" stop-color="#8b6914" stop-opacity="0.5"/>
</radialGradient>
`)

	// Wave pattern for water cells
	b.WriteString(`<pattern id="waves" x="0" y="0" width="12" height="6" patternUnits="userSpaceOnUse">
  <path d="M0 3 Q3 0 6 3 Q9 6 12 3" fill="none" stroke="rgba(255,255,255,0.25)" stroke-width="0.8"/>
</pattern>
`)

	// Clip path for map area
	fmt.Fprintf(b, `<clipPath id="mapClip">
  <rect x="%d" y="%d" width="%d" height="%d"/>
</clipPath>
`, mapX, mapY, mapW, mapH)

	b.WriteString("</defs>\n")
}

func (r *SVGMapRenderer) writeBorder(b *strings.Builder) {
	// Outer border
	fmt.Fprintf(b, `<rect x="4" y="4" width="%d" height="%d" fill="none" stroke="#5a3e1b" stroke-width="3" rx="4"/>
`, r.Width-8, r.Height-8)
	// Inner border
	fmt.Fprintf(b, `<rect x="10" y="10" width="%d" height="%d" fill="none" stroke="#5a3e1b" stroke-width="1" rx="2"/>
`, r.Width-20, r.Height-20)
}

func (r *SVGMapRenderer) writeTitle(b *strings.Builder, placeName string, titleH int) {
	cx := r.Width / 2
	cy := titleH/2 + 2

	// Cartouche background
	boxW := len(placeName)*14 + 80
	if boxW < 200 {
		boxW = 200
	}
	if boxW > r.Width-60 {
		boxW = r.Width - 60
	}
	boxH := titleH - 10

	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" rx="8" fill="#f4e8c1" stroke="#5a3e1b" stroke-width="2"/>
`, cx-boxW/2, cy-boxH/2, boxW, boxH)
	// Inner cartouche line
	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" rx="5" fill="none" stroke="#8b6914" stroke-width="0.7"/>
`, cx-boxW/2+4, cy-boxH/2+4, boxW-8, boxH-8)
	// Corner flourishes
	for _, dx := range []int{-boxW/2 + 12, boxW/2 - 12} {
		for _, dy := range []int{-boxH/2 + 8, boxH/2 - 8} {
			fmt.Fprintf(b, `<circle cx="%d" cy="%d" r="2" fill="#8b6914"/>
`, cx+dx, cy+dy)
		}
	}
	// Title text
	fmt.Fprintf(b, `<text x="%d" y="%d" text-anchor="middle" dominant-baseline="central" font-family="Georgia, 'Times New Roman', serif" font-size="20" font-weight="bold" fill="#3e2712" letter-spacing="3">%s</text>
`, cx, cy, xmlEscape(strings.ToUpper(placeName)))
}

func (r *SVGMapRenderer) writeTerrain(b *strings.Builder, t *TerrainMap, mapX, mapY int, cellW, cellH float64) {
	b.WriteString(fmt.Sprintf(`<g clip-path="url(#mapClip)">
`))
	for y := 0; y < t.Height; y++ {
		for x := 0; x < t.Width; x++ {
			cell := t.Cells[y][x]
			px := float64(mapX) + float64(x)*cellW
			py := float64(mapY) + float64(y)*cellH
			baseColor := cell.Biome.Color()

			// Height-based brightness variation: slightly lighten high cells, darken low ones
			brightness := (cell.Height - 0.5) * 0.15 // range roughly -0.075..+0.075
			fill := adjustBrightness(baseColor, brightness)

			fmt.Fprintf(b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>
`, px, py, cellW+0.5, cellH+0.5, fill)

			// Water cells get a wave pattern overlay
			if cell.Biome == Ocean || cell.Biome == ShallowWater {
				fmt.Fprintf(b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="url(#waves)"/>
`, px, py, cellW+0.5, cellH+0.5)
			}
		}
	}
	b.WriteString("</g>\n")
}

func (r *SVGMapRenderer) writeGrid(b *strings.Builder, t *TerrainMap, mapX, mapY, mapW, mapH int, cellW, cellH float64) {
	b.WriteString(`<g stroke="#5a3e1b" stroke-width="0.3" opacity="0.12">
`)
	for x := 0; x <= t.Width; x++ {
		px := float64(mapX) + float64(x)*cellW
		fmt.Fprintf(b, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d"/>
`, px, mapY, px, mapY+mapH)
	}
	for y := 0; y <= t.Height; y++ {
		py := float64(mapY) + float64(y)*cellH
		fmt.Fprintf(b, `<line x1="%d" y1="%.1f" x2="%d" y2="%.1f"/>
`, mapX, py, mapX+mapW, py)
	}
	b.WriteString("</g>\n")
}

func (r *SVGMapRenderer) writeLandmarks(b *strings.Builder, landmarks []Landmark, mapX, mapY int, cellW, cellH float64) {
	if len(landmarks) == 0 {
		return
	}
	b.WriteString(`<g>
`)
	for _, lm := range landmarks {
		cx := float64(mapX) + (float64(lm.X)+0.5)*cellW
		cy := float64(mapY) + (float64(lm.Y)+0.5)*cellH
		sz := cellW * 0.35
		if sz < 4 {
			sz = 4
		}
		if sz > 10 {
			sz = 10
		}
		r.writeLandmarkShape(b, lm.Type, cx, cy, sz)

		// Label
		labelX := cx + sz + 3
		labelY := cy + 3
		fmt.Fprintf(b, `<text x="%.1f" y="%.1f" font-family="Georgia, serif" font-size="8" fill="#3e2712" stroke="#f4e8c1" stroke-width="2" paint-order="stroke">%s</text>
`, labelX, labelY, xmlEscape(lm.Name))
	}
	b.WriteString("</g>\n")
}

func (r *SVGMapRenderer) writeLandmarkShape(b *strings.Builder, lt LandmarkType, cx, cy, sz float64) {
	markerColor := "#2c1810"
	markerStroke := "#f4e8c1"
	switch lt {
	case Town:
		// Filled circle
		fmt.Fprintf(b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="%s" stroke-width="1"/>
`, cx, cy, sz, markerColor, markerStroke)
	case Village:
		// Smaller filled circle
		fmt.Fprintf(b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="%s" stroke-width="0.8"/>
`, cx, cy, sz*0.7, markerColor, markerStroke)
	case Ruins:
		// Hollow circle
		fmt.Fprintf(b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke="%s" stroke-width="1.5"/>
`, cx, cy, sz, markerColor)
	case Tower:
		// Triangle pointing up
		fmt.Fprintf(b, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s" stroke-width="0.8"/>
`, cx, cy-sz, cx-sz*0.8, cy+sz*0.6, cx+sz*0.8, cy+sz*0.6, markerColor, markerStroke)
	case Cave:
		// Diamond
		fmt.Fprintf(b, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s" stroke-width="0.8"/>
`, cx, cy-sz, cx+sz, cy, cx, cy+sz, cx-sz, cy, markerColor, markerStroke)
	case Tavern:
		// Small house shape (pentagon)
		hw := sz * 0.8
		hh := sz * 0.6
		roofH := sz * 0.7
		fmt.Fprintf(b, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s" stroke-width="0.8"/>
`,
			cx, cy-hh-roofH, // roof peak
			cx+hw, cy-hh, // top-right
			cx+hw, cy+hh, // bottom-right
			cx-hw, cy+hh, // bottom-left
			cx-hw, cy-hh, // top-left
			markerColor, markerStroke)
	default:
		// Generic filled circle for all other types
		fmt.Fprintf(b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="%s" stroke-width="0.8"/>
`, cx, cy, sz*0.6, markerColor, markerStroke)
	}
}

func (r *SVGMapRenderer) writeCompass(b *strings.Builder, cx, cy int) {
	// Decorative circle
	fmt.Fprintf(b, `<g transform="translate(%d,%d)">
`, cx, cy)
	b.WriteString(`<circle cx="0" cy="0" r="22" fill="#f4e8c1" fill-opacity="0.85" stroke="#5a3e1b" stroke-width="1.5"/>
<circle cx="0" cy="0" r="18" fill="none" stroke="#8b6914" stroke-width="0.5"/>
`)
	// North arrow
	b.WriteString(`<polygon points="0,-16 3,-6 -3,-6" fill="#5a3e1b"/>
<polygon points="0,16 3,6 -3,6" fill="#8b6914"/>
<line x1="-16" y1="0" x2="16" y2="0" stroke="#8b6914" stroke-width="0.5"/>
<line x1="0" y1="-16" x2="0" y2="16" stroke="#8b6914" stroke-width="0.5"/>
`)
	// Direction labels
	b.WriteString(`<text x="0" y="-8" text-anchor="middle" font-family="Georgia, serif" font-size="7" font-weight="bold" fill="#3e2712">N</text>
<text x="0" y="13" text-anchor="middle" font-family="Georgia, serif" font-size="6" fill="#5a3e1b">S</text>
<text x="10" y="2" text-anchor="middle" font-family="Georgia, serif" font-size="6" fill="#5a3e1b">E</text>
<text x="-10" y="2" text-anchor="middle" font-family="Georgia, serif" font-size="6" fill="#5a3e1b">W</text>
`)
	// Diagonal ticks
	for _, angle := range []string{"45", "135", "225", "315"} {
		fmt.Fprintf(b, `<line x1="0" y1="-14" x2="0" y2="-11" stroke="#8b6914" stroke-width="0.5" transform="rotate(%s)"/>
`, angle)
	}
	b.WriteString("</g>\n")
}

func (r *SVGMapRenderer) writeLegend(b *strings.Builder, t *TerrainMap, landmarks []Landmark, lx, ly, lw, lh int) {
	// Semi-transparent legend panel
	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" rx="4" fill="#f4e8c1" fill-opacity="0.9" stroke="#5a3e1b" stroke-width="1"/>
`, lx, ly, lw, lh)

	// Collect unique biomes present on the map
	biomeSet := map[Biome]bool{}
	for y := 0; y < t.Height; y++ {
		for x := 0; x < t.Width; x++ {
			biomeSet[t.Cells[y][x].Biome] = true
		}
	}
	allBiomes := []Biome{Ocean, ShallowWater, Beach, Plains, Forest, DenseForest, Hills, Mountains, SnowPeak, Desert, Swamp, Tundra}
	var presentBiomes []Biome
	for _, bi := range allBiomes {
		if biomeSet[bi] {
			presentBiomes = append(presentBiomes, bi)
		}
	}

	// Title
	fmt.Fprintf(b, `<text x="%d" y="%d" font-family="Georgia, serif" font-size="9" font-weight="bold" fill="#3e2712">LEGEND</text>
`, lx+8, ly+14)

	// Biome swatches in two rows
	swatchSize := 8
	colWidth := 90
	startX := lx + 8
	startY := ly + 22
	for i, bi := range presentBiomes {
		col := i % ((lw - 16) / colWidth)
		row := i / ((lw - 16) / colWidth)
		sx := startX + col*colWidth
		sy := startY + row*14
		if sy+14 > ly+lh {
			break
		}
		fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s" stroke="#5a3e1b" stroke-width="0.5"/>
`, sx, sy, swatchSize, swatchSize, bi.Color())
		fmt.Fprintf(b, `<text x="%d" y="%d" font-family="Georgia, serif" font-size="7" fill="#3e2712">%s</text>
`, sx+swatchSize+3, sy+7, bi.String())
	}

	// Landmark symbols on the right side
	lmSectionX := lx + lw/2 + 40
	if lmSectionX > lx+lw-150 {
		lmSectionX = lx + lw - 150
	}
	fmt.Fprintf(b, `<text x="%d" y="%d" font-family="Georgia, serif" font-size="8" font-weight="bold" fill="#3e2712">Landmarks</text>
`, lmSectionX, ly+14)

	// Collect unique landmark types
	lmTypeSet := map[LandmarkType]bool{}
	for _, lm := range landmarks {
		lmTypeSet[lm.Type] = true
	}
	displayTypes := []struct {
		lt   LandmarkType
		desc string
	}{
		{Town, "Town"}, {Village, "Village"}, {Ruins, "Ruins"},
		{Tower, "Tower"}, {Cave, "Cave"}, {Tavern, "Tavern"},
	}

	idx := 0
	for _, dt := range displayTypes {
		if !lmTypeSet[dt.lt] {
			continue
		}
		col := idx % 2
		row := idx / 2
		sx := float64(lmSectionX + col*70)
		sy := float64(startY + row*14)
		if int(sy)+14 > ly+lh {
			break
		}
		r.writeLandmarkShape(b, dt.lt, sx+5, sy+4, 4)
		fmt.Fprintf(b, `<text x="%.0f" y="%.0f" font-family="Georgia, serif" font-size="7" fill="#3e2712">%s</text>
`, sx+14, sy+7, dt.desc)
		idx++
	}
}

// --- utility functions ---

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

// adjustBrightness modifies a hex color (#rrggbb) by a brightness delta (-1..+1).
func adjustBrightness(hex string, delta float64) string {
	if len(hex) != 7 || hex[0] != '#' {
		return hex
	}
	r, g, bl := hexToRGB(hex)
	r = clampInt(r+int(delta*255), 0, 255)
	g = clampInt(g+int(delta*255), 0, 255)
	bl = clampInt(bl+int(delta*255), 0, 255)
	return fmt.Sprintf("#%02x%02x%02x", r, g, bl)
}

func hexToRGB(hex string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return r, g, b
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
