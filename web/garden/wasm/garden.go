package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"strings"
)

// ---------------------------------------------------------------------------
// Plant types – mapped from function size
// ---------------------------------------------------------------------------

type PlantType int

const (
	Seedling PlantType = iota
	Sprout
	Bush
	Tree
	OakTree
	Flower
	Vine
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
		return "#90EE90"
	case Sprout:
		return "#228B22"
	case Bush:
		return "#228B22"
	case Tree:
		return "#006400"
	case OakTree:
		return "#004D00"
	case Flower:
		return "#FF69B4"
	case Vine:
		return "#90EE90"
	default:
		return "#FFFFFF"
	}
}

// ---------------------------------------------------------------------------
// Garden elements
// ---------------------------------------------------------------------------

type GardenElement int

const (
	Grass      GardenElement = iota
	Fence
	Bug
	Butterfly
	Weed
	Rock
	Mushroom
	Pond
	Gnome
	Tumbleweed
	Scarecrow
	Birdhouse
	Beehive
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
		return "#228B22"
	case Fence:
		return "#8B4513"
	case Bug:
		return "#FF0000"
	case Butterfly:
		return "#87CEEB"
	case Weed:
		return "#DAA520"
	case Rock:
		return "#808080"
	case Mushroom:
		return "#FF4500"
	case Pond:
		return "#4169E1"
	case Gnome:
		return "#FF0000"
	case Tumbleweed:
		return "#8B4513"
	case Scarecrow:
		return "#8B4513"
	case Birdhouse:
		return "#8B4513"
	case Beehive:
		return "#FFD700"
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
	Label   string
	Size    int
	Health  float64
}

type Garden struct {
	Width, Height int
	Plots         [][]GardenPlot
	Weather       string
	Season        string
	Title         string
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
// BuildGarden
// ---------------------------------------------------------------------------

func BuildGarden(stats *CodebaseStats) *Garden {
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

	plots := make([][]GardenPlot, height)
	for y := range plots {
		plots[y] = make([]GardenPlot, width)
		for x := range plots[y] {
			plots[y][x] = GardenPlot{
				X: x, Y: y,
				Element: Grass,
				Plant:   Seedling,
				Health:  1.0,
			}
		}
	}

	summary := GardenSummary{}
	plotsPerRow := width / 4
	if plotsPerRow < 1 {
		plotsPerRow = 1
	}

	for i, file := range stats.Files {
		col := i % plotsPerRow
		row := i / plotsPerRow
		baseX := col * 4
		baseY := row * 3
		if baseY >= height || baseX >= width {
			continue
		}

		fileName := file.Path

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

		for t := 0; t < file.Types; t++ {
			px, py := wrapCoord(baseX+t%4, baseY+1, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Mushroom
				plots[py][px].Plant = Flower
				plots[py][px].Label = fileName
				summary.Flowers++
			}
		}

		for t := 0; t < file.Tests; t++ {
			px := baseX + t%4
			py := baseY
			px, py = wrapCoord(px, py, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Fence
				plots[py][px].Label = fileName
				summary.Fences++
			}
		}

		for t := 0; t < file.TODOs; t++ {
			px, py := wrapCoord(baseX+1+t%3, baseY+2, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Bug
				plots[py][px].Label = fileName
				summary.Bugs++
			}
		}

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

		if file.TODOs > 3 || (file.Functions > 0 && file.Tests == 0 && file.Comments == 0) {
			px, py := wrapCoord(baseX+3, baseY+2, width, height)
			if plots[py][px].Element == Grass {
				plots[py][px].Element = Weed
				summary.Weeds++
			}
		}

		if file.Imports > 10 {
			px, py := wrapCoord(baseX, baseY, width, height)
			plots[py][px].Element = Beehive
			plots[py][px].Label = fileName
		}

		if file.Lines > 0 && float64(file.Comments)/float64(file.Lines) > 0.3 {
			px, py := wrapCoord(baseX+3, baseY, width, height)
			plots[py][px].Element = Birdhouse
			plots[py][px].Label = fileName
		}

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

	weather := determineWeather(stats.HealthScore, stats.TotalTests)
	season := determineSeason(stats)

	if stats.HealthScore > 90 {
		gx, gy := width/2, height/2
		plots[gy][gx].Element = Gnome
		plots[gy][gx].Label = "Garden Gnome 🎉"
	}
	if stats.HealthScore < 20 {
		tx, ty := width/3, height/3
		plots[ty][tx].Element = Tumbleweed
		plots[ty][tx].Label = "Abandoned..."
	}

	scatterPonds(plots, width, height, stats.RootPath)

	title := stats.RootPath
	if title == "" || title == "." {
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
		return "Winter ❄️"
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

func scatterPonds(plots [][]GardenPlot, w, h int, seed string) {
	hsh := fnv.New32a()
	hsh.Write([]byte(seed))
	v := int(hsh.Sum32())

	count := (w * h) / 40
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
