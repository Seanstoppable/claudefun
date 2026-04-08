package main

import (
	"hash/fnv"
	"math"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------------------------
// Plant types – mapped from function size
// ---------------------------------------------------------------------------

type PlantType int

const (
	Seedling PlantType = iota // tiny function  (1-5 lines)
	Sprout                    // small function (6-15 lines)
	Bush                      // medium function (16-30 lines)
	Tree                      // large function  (31-60 lines)
	OakTree                   // huge function   (60+ lines)
	Flower                    // type / struct definition
	Vine                      // import / dependency
)

func (p PlantType) Symbol() string {
	switch p {
	case Seedling:
		return "🌱"
	case Sprout:
		return "ψ"
	case Bush:
		return "♣"
	case Tree:
		return "♠"
	case OakTree:
		return "🌳"
	case Flower:
		return "✿"
	case Vine:
		return "~"
	default:
		return "?"
	}
}

func (p PlantType) Color() string {
	switch p {
	case Seedling:
		return "#90EE90" // light green
	case Sprout:
		return "#228B22" // green
	case Bush:
		return "#228B22" // green
	case Tree:
		return "#006400" // dark green
	case OakTree:
		return "#004D00" // dark green
	case Flower:
		return "#FF69B4" // pink
	case Vine:
		return "#90EE90" // light green
	default:
		return "#FFFFFF"
	}
}

// ---------------------------------------------------------------------------
// Garden elements – decorations, critters, terrain
// ---------------------------------------------------------------------------

type GardenElement int

const (
	Grass      GardenElement = iota
	Fence                    // tests
	Bug                      // TODOs / FIXMEs
	Butterfly                // comments
	Weed                     // dead code indicators
	Rock                     // complexity hotspots
	Mushroom                 // type definitions
	Pond                     // blank space / breathing room
	Gnome                    // easter egg (healthy gardens)
	Tumbleweed               // unhealthy / abandoned code
	Scarecrow                // error handling code
	Birdhouse                // well-documented function
	Beehive                  // heavily imported module
)

func (e GardenElement) Symbol() string {
	switch e {
	case Grass:
		return "·"
	case Fence:
		return "┃"
	case Bug:
		return "🐛"
	case Butterfly:
		return "🦋"
	case Weed:
		return "⌇"
	case Rock:
		return "●"
	case Mushroom:
		return "♤"
	case Pond:
		return "≈"
	case Gnome:
		return "⚑"
	case Tumbleweed:
		return "◎"
	case Scarecrow:
		return "╂"
	case Birdhouse:
		return "⌂"
	case Beehive:
		return "◈"
	default:
		return " "
	}
}

func (e GardenElement) Color() string {
	switch e {
	case Grass:
		return "#228B22" // green
	case Fence:
		return "#8B4513" // brown
	case Bug:
		return "#FF0000" // red
	case Butterfly:
		return "#87CEEB" // light blue
	case Weed:
		return "#DAA520" // yellow-brown
	case Rock:
		return "#808080" // gray
	case Mushroom:
		return "#FF4500" // red/white
	case Pond:
		return "#4169E1" // blue
	case Gnome:
		return "#FF0000" // red
	case Tumbleweed:
		return "#8B4513" // brown
	case Scarecrow:
		return "#8B4513" // brown
	case Birdhouse:
		return "#8B4513" // brown
	case Beehive:
		return "#FFD700" // yellow
	default:
		return "#FFFFFF"
	}
}

// ---------------------------------------------------------------------------
// Garden data structures
// ---------------------------------------------------------------------------

type GardenPlot struct {
	X, Y    int
	Element GardenElement
	Plant   PlantType
	Label   string  // file / function name
	Size    int     // for scaling
	Health  float64 // 0-1 per-plot health
}

type Garden struct {
	Width, Height int
	Plots         [][]GardenPlot
	Weather       string // "Sunny ☀️", "Cloudy ☁️", etc.
	Season        string // "Spring", "Summer", "Autumn", "Winter"
	Title         string // garden name (from directory name)
	HealthScore   float64
	Stats         GardenSummary
}

type GardenSummary struct {
	Trees       int
	Flowers     int
	Fences      int
	Bugs        int
	Butterflies int
	Weeds       int
	Rocks       int
}

// ---------------------------------------------------------------------------
// BuildGarden – turns boring numbers into a charming garden
// ---------------------------------------------------------------------------

func BuildGarden(stats *CodebaseStats) *Garden {
	// 1. Calculate garden dimensions
	totalFiles := stats.TotalFiles
	if totalFiles < 1 {
		totalFiles = 1
	}
	width := int(math.Sqrt(float64(totalFiles)) * 4)
	if width < 10 {
		width = 10
	}
	height := int(float64(width) * 0.6)
	if height < 6 {
		height = 6
	}

	// 2. Fill with Grass
	plots := make([][]GardenPlot, height)
	for y := range plots {
		plots[y] = make([]GardenPlot, width)
		for x := range plots[y] {
			plots[y][x] = GardenPlot{
				X:       x,
				Y:       y,
				Element: Grass,
				Plant:   Seedling,
				Health:  1.0,
			}
		}
	}

	summary := GardenSummary{}

	// 3. Allocate plot areas for each file
	plotsPerRow := width / 4
	if plotsPerRow < 1 {
		plotsPerRow = 1
	}

	for i, file := range stats.Files {
		// Compute the top-left corner of this file's plot area
		col := i % plotsPerRow
		row := i / plotsPerRow
		baseX := col * 4
		baseY := row * 3
		if baseY >= height || baseX >= width {
			continue
		}

		fileName := filepath.Base(file.Path)

		// --- Functions → Plants ---
		funcsToPlace := file.Functions
		if funcsToPlace < 1 {
			funcsToPlace = 1
		}
		avgLines := 0
		if file.Functions > 0 {
			avgLines = file.Lines / file.Functions
		}
		plant := plantFromSize(avgLines)

		placed := 0
		for dy := 0; dy < 3 && placed < funcsToPlace; dy++ {
			for dx := 0; dx < 4 && placed < funcsToPlace; dx++ {
				px, py := baseX+dx, baseY+dy
				if py < height && px < width {
					plots[py][px].Plant = plant
					plots[py][px].Label = fileName
					plots[py][px].Size = avgLines
					plots[py][px].Health = fileHealth(file)
					placed++
					summary.Trees++
				}
			}
		}

		// --- Types → Flowers / Mushrooms ---
		for t := 0; t < file.Types; t++ {
			px, py := wrapCoord(baseX+t%4, baseY+1, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Mushroom
				plots[py][px].Plant = Flower
				plots[py][px].Label = fileName
				summary.Flowers++
			}
		}

		// --- Tests → Fences (edges of file plot) ---
		for t := 0; t < file.Tests; t++ {
			// Place fences along the top edge
			px := baseX + t%4
			py := baseY
			px, py = wrapCoord(px, py, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Fence
				plots[py][px].Label = fileName
				summary.Fences++
			}
		}

		// --- TODOs → Bugs ---
		for t := 0; t < file.TODOs; t++ {
			px, py := wrapCoord(baseX+1+t%3, baseY+2, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Bug
				plots[py][px].Label = fileName
				summary.Bugs++
			}
		}

		// --- Comments → Butterflies ---
		butterflyCount := file.Comments / 5
		if file.Comments > 0 && butterflyCount == 0 {
			butterflyCount = 1
		}
		for b := 0; b < butterflyCount; b++ {
			px, py := wrapCoord(baseX+b%4, baseY, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Butterfly
				plots[py][px].Label = fileName
				summary.Butterflies++
			}
		}

		// --- High complexity → Rocks ---
		if file.Complexity > 10 {
			rockCount := file.Complexity / 10
			for r := 0; r < rockCount; r++ {
				px, py := wrapCoord(baseX+2+r%2, baseY+1+r%2, width, height)
				if plots[py][px].Element == Grass {
					plots[py][px].Element = Rock
					plots[py][px].Label = fileName
					summary.Rocks++
				}
			}
		}

		// --- Dead code indicators → Weeds ---
		if file.TODOs > 3 || (file.Functions > 0 && file.Tests == 0 && file.Comments == 0) {
			px, py := wrapCoord(baseX+3, baseY+2, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Weed
				summary.Weeds++
			}
		}

		// --- Beehive for heavily imported modules (>10 imports) ---
		if file.Imports > 10 {
			px, py := wrapCoord(baseX, baseY, width, height)
			plots[py][px].Element = Beehive
			plots[py][px].Label = fileName
		}

		// --- Birdhouse for well-documented functions (comment ratio > 0.3) ---
		if file.Lines > 0 && float64(file.Comments)/float64(file.Lines) > 0.3 {
			px, py := wrapCoord(baseX+3, baseY, width, height)
			plots[py][px].Element = Birdhouse
			plots[py][px].Label = fileName
		}

		// --- Scarecrow for error handling code ---
		hasErrors := strings.Contains(strings.ToLower(file.Path), "error") ||
			strings.Contains(strings.ToLower(file.Path), "err") ||
			file.Complexity > 5
		if hasErrors {
			px, py := wrapCoord(baseX+1, baseY+1, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Scarecrow
				plots[py][px].Label = fileName
			}
		}
	}

	// 4. Determine Weather from health score
	weather := determineWeather(stats.HealthScore, stats.TotalTests)

	// 5. Determine Season
	season := determineSeason(stats)

	// 6. Special elements
	if stats.HealthScore > 90 {
		// Gnome easter egg!
		gx, gy := width/2, height/2
		plots[gy][gx].Element = Gnome
		plots[gy][gx].Label = "Garden Gnome 🎉"
	}
	if stats.HealthScore < 20 {
		// Tumbleweed for sad gardens
		tx, ty := width/3, height/3
		plots[ty][tx].Element = Tumbleweed
		plots[ty][tx].Label = "Abandoned..."
	}

	// 7. Scatter Ponds for visual breathing room
	scatterPonds(plots, width, height, stats.RootPath)

	title := filepath.Base(stats.RootPath)
	if title == "." || title == "" {
		title = "My Code Garden"
	}

	return &Garden{
		Width:       width,
		Height:      height,
		Plots:       plots,
		Weather:     weather,
		Season:      season,
		Title:       title,
		HealthScore: stats.HealthScore,
		Stats:       summary,
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func plantFromSize(avgLines int) PlantType {
	switch {
	case avgLines <= 5:
		return Seedling
	case avgLines <= 15:
		return Sprout
	case avgLines <= 30:
		return Bush
	case avgLines <= 60:
		return Tree
	default:
		return OakTree
	}
}

func fileHealth(f FileStats) float64 {
	h := 1.0
	if f.Tests == 0 && f.Functions > 0 {
		h -= 0.2
	}
	if f.TODOs > 0 {
		h -= 0.1 * math.Min(float64(f.TODOs), 3)
	}
	if f.Complexity > 20 {
		h -= 0.2
	}
	if f.Comments == 0 && f.Lines > 50 {
		h -= 0.1
	}
	if h < 0 {
		h = 0
	}
	return h
}

func determineWeather(healthScore float64, totalTests int) string {
	switch {
	case healthScore >= 80 && totalTests > 0:
		return "Rainbow 🌈"
	case healthScore >= 80:
		return "Sunny ☀️"
	case healthScore >= 60:
		return "Partly Cloudy ⛅"
	case healthScore >= 40:
		return "Cloudy ☁️"
	case healthScore >= 20:
		return "Rainy 🌧️"
	default:
		return "Stormy ⛈️"
	}
}

func determineSeason(stats *CodebaseStats) string {
	if stats.TotalFiles == 0 {
		return "Winter"
	}

	avgFuncSize := 0
	if stats.TotalFuncs > 0 {
		avgFuncSize = stats.TotalLines / stats.TotalFuncs
	}

	todoRatio := float64(stats.TotalTODOs) / float64(stats.TotalFiles)

	switch {
	case avgFuncSize <= 10 && stats.TotalFuncs > stats.TotalFiles:
		return "Spring 🌸"
	case stats.TotalLines > 1000 && stats.TotalFuncs > 20:
		return "Summer ☀️"
	case todoRatio > 1.0:
		return "Autumn 🍂"
	default:
		return "Winter ❄️"
	}
}

// wrapCoord keeps coordinates within the garden bounds.
func wrapCoord(x, y, w, h int) (int, int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= w {
		x = w - 1
	}
	if y >= h {
		y = h - 1
	}
	return x, y
}

// scatterPonds places small ponds at deterministic positions based on
// a hash of the garden title so the layout is stable across runs.
func scatterPonds(plots [][]GardenPlot, w, h int, seed string) {
	hsh := fnv.New32a()
	hsh.Write([]byte(seed))
	v := int(hsh.Sum32())

	count := (w * h) / 40 // ~2.5 % coverage
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		px := ((v + i*17) % w)
		py := ((v + i*31) % h)
		if px < 0 {
			px = -px
		}
		if py < 0 {
			py = -py
		}
		px = px % w
		py = py % h
		if plots[py][px].Element == Grass {
			plots[py][px].Element = Pond
		}
	}
}
