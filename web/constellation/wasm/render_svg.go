package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"

	"strings"
)

// SVGRenderer turns a Constellation into a self-contained SVG document.
type SVGRenderer struct {
	Width  int // SVG width in pixels (default 800)
	Height int // SVG height in pixels (default 600)
}

// NewSVGRenderer creates a renderer with the given pixel dimensions.
func NewSVGRenderer(width, height int) *SVGRenderer {
	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 600
	}
	return &SVGRenderer{Width: width, Height: height}
}

// Render produces a complete SVG string for the constellation.
func (r *SVGRenderer) Render(c *Constellation) string {
	var b strings.Builder
	w := float64(r.Width)
	h := float64(r.Height)
	name := c.Stars.ConstellationName()

	// Deterministic RNG for background stars.
	rng := seedRNG(c.Stars.Seed)

	// Header + metadata comment.
	b.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`,
		r.Width, r.Height, r.Width, r.Height))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  <!-- Constellation: %q -->\n", c.Stars.Seed))

	// ── defs: filters & gradients ──────────────────────────────────────
	r.writeDefs(&b)

	// ── background ─────────────────────────────────────────────────────
	b.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#0a0a2e"/>`, r.Width, r.Height))
	b.WriteString("\n")

	// Subtle vignette overlay.
	b.WriteString(fmt.Sprintf(
		`  <rect width="%d" height="%d" fill="url(#vignette)" opacity="0.4"/>`,
		r.Width, r.Height))
	b.WriteString("\n")

	// ── background stars ───────────────────────────────────────────────
	b.WriteString("  <!-- background stars -->\n")
	b.WriteString(`  <g opacity="0.6">` + "\n")
	bgCount := 50 + rng.Intn(31) // 50-80
	for i := 0; i < bgCount; i++ {
		bx := rng.Float64() * w
		by := rng.Float64() * h
		br := 0.5 + rng.Float64()*1.5 // 0.5-2px
		op := 0.15 + rng.Float64()*0.4
		fill := "#ffffff"
		if rng.Float64() < 0.3 {
			fill = "#b0c4de" // light steel blue tint
		}
		b.WriteString(fmt.Sprintf(
			`    <circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" opacity="%.2f"/>`,
			bx, by, br, fill, op))
		b.WriteString("\n")
	}
	b.WriteString("  </g>\n")

	// ── constellation lines ────────────────────────────────────────────
	b.WriteString("  <!-- constellation lines -->\n")
	b.WriteString(`  <g stroke="#4a6fa5" stroke-width="1.5" stroke-linecap="round" opacity="0.6">` + "\n")
	for _, e := range c.Edges {
		s1 := c.Stars.Stars[e.From]
		s2 := c.Stars.Stars[e.To]
		x1, y1 := s1.X*w, s1.Y*h
		x2, y2 := s2.X*w, s2.Y*h
		b.WriteString(fmt.Sprintf(
			`    <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" filter="url(#lineGlow)"/>`,
			x1, y1, x2, y2))
		b.WriteString("\n")
	}
	b.WriteString("  </g>\n")

	// ── stars with glow ────────────────────────────────────────────────
	b.WriteString("  <!-- stars -->\n")
	for i, s := range c.Stars.Stars {
		sx := s.X * w
		sy := s.Y * h
		radius := 2.0 + s.Size*6.0 // 2-8px
		fill := starColor(s.Brightness)
		glowRadius := radius * 3.0
		glowOpacity := 0.1 + s.Brightness*0.25

		// Outer glow circle.
		b.WriteString(fmt.Sprintf(
			`  <circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" opacity="%.2f" filter="url(#glow)"/>`,
			sx, sy, glowRadius, fill, glowOpacity))
		b.WriteString("\n")

		// Core star.
		b.WriteString(fmt.Sprintf(
			`  <circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" id="star%d"/>`,
			sx, sy, radius, fill, i))
		b.WriteString("\n")

		// Bright star spikes (cross effect for brightness >= 0.75).
		if s.Brightness >= 0.75 {
			spikeLen := radius * 2.5
			spikeOp := 0.3 + (s.Brightness-0.75)*1.0 // 0.3-0.55
			b.WriteString(fmt.Sprintf(
				`  <g stroke="%s" stroke-width="0.8" opacity="%.2f">`, fill, spikeOp))
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf(
				`    <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f"/>`,
				sx-spikeLen, sy, sx+spikeLen, sy))
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf(
				`    <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f"/>`,
				sx, sy-spikeLen, sx, sy+spikeLen))
			b.WriteString("\n")
			b.WriteString("  </g>\n")
		}
	}

	// ── title text ─────────────────────────────────────────────────────
	b.WriteString("  <!-- title -->\n")
	b.WriteString(fmt.Sprintf(
		`  <text x="%.1f" y="%.1f" `+
			`font-family="Georgia, 'Times New Roman', serif" `+
			`font-size="14" fill="#6b7ea8" `+
			`text-anchor="middle" opacity="0.7" `+
			`letter-spacing="3">%s</text>`,
		w/2, h-20, escapeXML(name)))
	b.WriteString("\n")

	b.WriteString("</svg>\n")
	return b.String()
}


// ── defs ───────────────────────────────────────────────────────────────

func (r *SVGRenderer) writeDefs(b *strings.Builder) {
	b.WriteString("  <defs>\n")

	// Star glow filter.
	b.WriteString(`    <filter id="glow" x="-50%" y="-50%" width="200%" height="200%">` + "\n")
	b.WriteString(`      <feGaussianBlur in="SourceGraphic" stdDeviation="4" result="blur"/>` + "\n")
	b.WriteString(`      <feMerge>` + "\n")
	b.WriteString(`        <feMergeNode in="blur"/>` + "\n")
	b.WriteString(`        <feMergeNode in="SourceGraphic"/>` + "\n")
	b.WriteString(`      </feMerge>` + "\n")
	b.WriteString(`    </filter>` + "\n")

	// Subtle line glow filter.
	b.WriteString(`    <filter id="lineGlow" x="-20%" y="-20%" width="140%" height="140%">` + "\n")
	b.WriteString(`      <feGaussianBlur in="SourceGraphic" stdDeviation="2" result="blur"/>` + "\n")
	b.WriteString(`      <feMerge>` + "\n")
	b.WriteString(`        <feMergeNode in="blur"/>` + "\n")
	b.WriteString(`        <feMergeNode in="SourceGraphic"/>` + "\n")
	b.WriteString(`      </feMerge>` + "\n")
	b.WriteString(`    </filter>` + "\n")

	// Vignette radial gradient.
	b.WriteString(`    <radialGradient id="vignette" cx="50%" cy="50%" r="70%">` + "\n")
	b.WriteString(`      <stop offset="0%" stop-color="transparent"/>` + "\n")
	b.WriteString(`      <stop offset="100%" stop-color="#000000"/>` + "\n")
	b.WriteString(`    </radialGradient>` + "\n")

	b.WriteString("  </defs>\n")
}

// ── helpers ────────────────────────────────────────────────────────────

// starColor maps brightness 0-1 to a color from blue-gray (dim) to warm white (bright).
func starColor(brightness float64) string {
	// Dim: #7b8ea8 → Mid: #c8d6e5 → Bright: #fff8e7
	if brightness < 0.5 {
		t := brightness / 0.5
		r := lerp(0x7b, 0xc8, t)
		g := lerp(0x8e, 0xd6, t)
		b := lerp(0xa8, 0xe5, t)
		return fmt.Sprintf("#%02x%02x%02x", int(r), int(g), int(b))
	}
	t := (brightness - 0.5) / 0.5
	r := lerp(0xc8, 0xff, t)
	g := lerp(0xd6, 0xf8, t)
	b := lerp(0xe5, 0xe7, t)
	return fmt.Sprintf("#%02x%02x%02x", int(r), int(g), int(b))
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// seedRNG creates a deterministic *rand.Rand from a string.
func seedRNG(seed string) *rand.Rand {
	h := sha256.Sum256([]byte(seed))
	s := int64(binary.BigEndian.Uint64(h[:8]))
	return rand.New(rand.NewSource(s))
}

// escapeXML escapes the five XML special characters.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}


