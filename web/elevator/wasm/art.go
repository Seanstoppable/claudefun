package main

import (
	"math/rand"
	"strings"
)

// ArtGenerator produces procedural ASCII album art.
type ArtGenerator struct {
	rng *rand.Rand
}

// NewArtGenerator creates an ArtGenerator with the given seed.
func NewArtGenerator(seed int64) *ArtGenerator {
	return &ArtGenerator{rng: rand.New(rand.NewSource(seed))}
}

// GenerateArt returns an ASCII album cover (~12-15 lines, ~20-25 chars wide).
// Style is selected based on genre, with variation from the seed.
func (ag *ArtGenerator) GenerateArt(genre Genre, mood string) string {
	switch genre {
	case AmbientLobby, CorporateZen, SuburbanAmbient, AcousticBeige:
		return ag.pickFrom([]func() string{ag.geometric, ag.diamondGrid})
	case WaitingRoomJazz, LoFiElevator:
		return ag.pickFrom([]func() string{ag.minimalistScene, ag.nightSky})
	case HoldMusicDeluxe, SmoothBureaucracy, DentistOfficeCore:
		return ag.pickFrom([]func() string{ag.wavePattern, ag.soundBars})
	case MidtempoMalaise, SoftRockPurgatory, ParkingGarageWave:
		return ag.pickFrom([]func() string{ag.gradient, ag.voidStare})
	default:
		// Unknown genre — pick a random style
		all := []func() string{
			ag.geometric, ag.minimalistScene, ag.wavePattern,
			ag.gradient, ag.diamondGrid, ag.nightSky,
			ag.soundBars, ag.voidStare,
		}
		return ag.pickFrom(all)
	}
}

func (ag *ArtGenerator) pickFrom(fns []func() string) string {
	return fns[ag.rng.Intn(len(fns))]()
}

// Style 1: Geometric/Abstract
func (ag *ArtGenerator) geometric() string {
	shapes := []string{"◇", "○", "□", "△", "◆"}
	s := shapes[ag.rng.Intn(len(shapes))]

	var b strings.Builder
	b.WriteString("┌──────────────────┐\n")
	b.WriteString("│                  │\n")
	b.WriteString("│    " + s + "      " + s + "     │\n")
	b.WriteString("│      " + s + "  " + s + "       │\n")
	b.WriteString("│        " + s + "        │\n")
	b.WriteString("│      " + s + "  " + s + "       │\n")
	b.WriteString("│    " + s + "      " + s + "     │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│  ░░░░░░░░░░░░░░  │\n")
	b.WriteString("│  ░░░░░░░░░░░░░░  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("└──────────────────┘\n")
	return b.String()
}

// Style 2: Minimalist Scene
func (ag *ArtGenerator) minimalistScene() string {
	labels := []string{"LOBBY", "FLOOR 7", "SUITE B", "LEVEL 3", "ATRIUM"}
	label := labels[ag.rng.Intn(len(labels))]
	padded := centerPad(label, 9)

	var b strings.Builder
	b.WriteString("┌──────────────────┐\n")
	b.WriteString("│    *  .    .  *  │\n")
	b.WriteString("│  .    *  .       │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│ ♪  ♫  ♪  ♫  ♪   │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│  ___________     │\n")
	b.WriteString("│ /           \\    │\n")
	b.WriteString("│| ♪" + padded + "♪ |   │\n")
	b.WriteString("│ \\___________/    │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│ ═══════════════  │\n")
	b.WriteString("└──────────────────┘\n")
	return b.String()
}

// Style 3: Wave/Pattern
func (ag *ArtGenerator) wavePattern() string {
	waves := []string{"∿", "~", "≈"}
	w := waves[ag.rng.Intn(len(waves))]
	line := strings.Repeat(w, 14)
	offset := " " + strings.Repeat(w, 13)

	var b strings.Builder
	b.WriteString("┌──────────────────┐\n")
	b.WriteString("│                  │\n")
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			b.WriteString("│ " + line + "  │\n")
		} else {
			b.WriteString("│" + offset + "   │\n")
		}
		b.WriteString("│                  │\n")
	}
	b.WriteString("└──────────────────┘\n")
	return b.String()
}

// Style 4: Gradient/Blocks
func (ag *ArtGenerator) gradient() string {
	fills := [][]string{
		{"█", "▓", "▒", "░"},
		{"▓", "▒", "░", " "},
		{"█", "▓", "░", " "},
	}
	fill := fills[ag.rng.Intn(len(fills))]

	var b strings.Builder
	b.WriteString("┌──────────────────┐\n")
	for _, ch := range fill {
		row := " " + strings.Repeat(ch, 16) + " "
		b.WriteString("│" + row + "│\n")
		b.WriteString("│" + row + "│\n")
	}
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("└──────────────────┘\n")
	return b.String()
}

// Style 5: Diamond grid
func (ag *ArtGenerator) diamondGrid() string {
	glyphs := []string{"◆", "◇", "●", "○"}
	a := glyphs[ag.rng.Intn(len(glyphs))]
	c := glyphs[ag.rng.Intn(len(glyphs))]

	var b strings.Builder
	b.WriteString("┌──────────────────┐\n")
	b.WriteString("│                  │\n")
	for row := 0; row < 4; row++ {
		if row%2 == 0 {
			b.WriteString("│  " + a + "  " + c + "  " + a + "  " + c + "  " + a + "   │\n")
		} else {
			b.WriteString("│    " + c + "  " + a + "  " + c + "  " + a + "    │\n")
		}
	}
	b.WriteString("│                  │\n")
	b.WriteString("│  ┄┄┄┄┄┄┄┄┄┄┄┄┄  │\n")
	b.WriteString("│   E L E V A T E  │\n")
	b.WriteString("│  ┄┄┄┄┄┄┄┄┄┄┄┄┄  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("└──────────────────┘\n")
	return b.String()
}

// Style 6: Night sky
func (ag *ArtGenerator) nightSky() string {
	var b strings.Builder
	b.WriteString("┌──────────────────┐\n")

	// 8 rows of stars
	for i := 0; i < 8; i++ {
		row := []byte("│                  │\n")
		numStars := 2 + ag.rng.Intn(4)
		for s := 0; s < numStars; s++ {
			pos := 1 + ag.rng.Intn(18)
			if row[pos] == ' ' {
				stars := []byte{'*', '.', '+', '`'}
				row[pos] = stars[ag.rng.Intn(len(stars))]
			}
		}
		b.Write(row)
	}

	b.WriteString("│   ♪ ═══════ ♪    │\n")
	b.WriteString("│    ~ night ~     │\n")
	b.WriteString("│   ♪ ═══════ ♪    │\n")
	b.WriteString("│                  │\n")
	b.WriteString("└──────────────────┘\n")
	return b.String()
}

// Style 7: Sound bars / equalizer
func (ag *ArtGenerator) soundBars() string {
	var b strings.Builder
	b.WriteString("┌──────────────────┐\n")
	b.WriteString("│                  │\n")

	// Generate 8 bar heights (1-8)
	heights := make([]int, 8)
	for i := range heights {
		heights[i] = 1 + ag.rng.Intn(8)
	}

	for row := 8; row >= 1; row-- {
		line := make([]byte, 18)
		for i := range line {
			line[i] = ' '
		}
		for i, h := range heights {
			if h >= row {
				col := 1 + i*2
				if col < len(line) {
					line[col] = '#'
				}
			}
		}
		b.WriteString("│" + string(line) + "│\n")
	}

	b.WriteString("│ ──────────────── │\n")
	b.WriteString("│   ♪  NOW  ♪     │\n")
	b.WriteString("│   PLAYING       │\n")
	b.WriteString("└──────────────────┘\n")
	return b.String()
}

// Style 8: Void stare — existential minimalism
func (ag *ArtGenerator) voidStare() string {
	msgs := []string{
		"  this is fine.  ",
		" floor unknown.  ",
		"   going down.   ",
		"  please wait.   ",
		"    forever.     ",
	}
	msg := msgs[ag.rng.Intn(len(msgs))]

	var b strings.Builder
	b.WriteString("┌──────────────────┐\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│ " + msg + "│\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("│                  │\n")
	b.WriteString("└──────────────────┘\n")
	return b.String()
}

// centerPad centers s within a field of width w.
func centerPad(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	left := (w - len(s)) / 2
	right := w - len(s) - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
