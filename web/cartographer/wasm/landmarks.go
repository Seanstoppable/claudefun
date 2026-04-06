package main

import (
	"fmt"
	"hash/fnv"
	"math/rand"
)

// LandmarkType represents the category of a landmark.
type LandmarkType int

const (
	Town LandmarkType = iota
	Village
	Ruins
	Tower
	Cave
	Bridge
	Tavern
	Monolith
	Shrine
	Lighthouse
	Mine
	WitchHut
	DragonLair
	Library
	Market
)

// Landmark is a named point of interest placed on the terrain map.
type Landmark struct {
	X, Y   int
	Type   LandmarkType
	Name   string
	Symbol string // for terminal rendering
	Icon   string // for SVG rendering (emoji or shape description)
}

// LandmarkGenerator places and names landmarks on a terrain map.
type LandmarkGenerator struct {
	rng *rand.Rand
}

// NewLandmarkGenerator creates a generator seeded from the given place-name string.
func NewLandmarkGenerator(seed string) *LandmarkGenerator {
	h := fnv.New64a()
	h.Write([]byte(seed))
	return &LandmarkGenerator{
		rng: rand.New(rand.NewSource(int64(h.Sum64()))),
	}
}

// ── name-generation tables ────────────────────────────────────────────

var namePrefixes = []string{
	"Bram", "Thorn", "Oak", "Iron", "Silver",
	"Dusk", "Storm", "Frost", "Ember", "Moss",
	"Hollow", "Raven", "Copper", "Moon", "Stone",
	"Glen", "Ash", "Wren", "Fox", "Eld",
}

var nameMiddles = []string{
	"ble", "worth", "mere", "dale", "wick",
	"brook", "ford", "gate", "haven", "wood",
	"fell", "bury", "ton", "stead", "crest",
}

var nameSuffixes = []string{
	"'s Rest", "'s Folly", " Point", " Crossing", " Keep",
	" Hollow", " Watch", " Depths", " End", " Landing",
	"", "", "",
}

var tavernAdjs = []string{
	"Sleepy", "Crooked", "Ambitious", "Weeping", "Dancing",
	"Suspicious", "Bewildered", "Resigned", "Enthusiastic", "Philosophical",
}

var tavernNouns = []string{
	"Badger", "Heron", "Turnip", "Compass", "Anchor",
	"Mushroom", "Owl", "Kettle", "Boot", "Dragon",
}

var ruinNouns = []string{
	"Keep", "Tower", "Hall", "Temple", "Citadel",
	"Gate", "Throne", "Beacon", "Bridge", "Spire",
}

var monolithAdjs = []string{
	"Whispering", "Leaning", "Singing", "Judgmental", "Indifferent",
}

// ── symbol / icon tables ──────────────────────────────────────────────

var landmarkSymbols = map[LandmarkType]string{
	Town:       "■",
	Village:    "□",
	Ruins:      "⌘",
	Tower:      "↑",
	Cave:       "◎",
	Bridge:     "═",
	Tavern:     "⌂",
	Monolith:   "◆",
	Shrine:     "✟",
	Lighthouse: "☆",
	Mine:       "⛏",
	WitchHut:   "⚗",
	DragonLair: "⊛",
	Library:    "◫",
	Market:     "◈",
}

var landmarkIcons = map[LandmarkType]string{
	Town:       "🏘️",
	Village:    "🏡",
	Ruins:      "🏚️",
	Tower:      "🗼",
	Cave:       "🕳️",
	Bridge:     "🌉",
	Tavern:     "🍺",
	Monolith:   "🗿",
	Shrine:     "⛩️",
	Lighthouse: "🏠",
	Mine:       "⛏️",
	WitchHut:   "🧙",
	DragonLair: "🐉",
	Library:    "📚",
	Market:     "🛒",
}

var landmarkTypeNames = map[LandmarkType]string{
	Town:       "Town",
	Village:    "Village",
	Ruins:      "Ruins",
	Tower:      "Tower",
	Cave:       "Cave",
	Bridge:     "Bridge",
	Tavern:     "Tavern",
	Monolith:   "Monolith",
	Shrine:     "Shrine",
	Lighthouse: "Lighthouse",
	Mine:       "Mine",
	WitchHut:   "Witch Hut",
	DragonLair: "Dragon Lair",
	Library:    "Library",
	Market:     "Market",
}

// String returns the human-readable name of a LandmarkType.
func (lt LandmarkType) String() string {
	if s, ok := landmarkTypeNames[lt]; ok {
		return s
	}
	return fmt.Sprintf("LandmarkType(%d)", int(lt))
}

// Symbol returns the single-glyph terminal symbol for a LandmarkType.
func (lt LandmarkType) Symbol() string {
	if s, ok := landmarkSymbols[lt]; ok {
		return s
	}
	return "?"
}

// ── name generators ───────────────────────────────────────────────────

func (lg *LandmarkGenerator) pick(list []string) string {
	return list[lg.rng.Intn(len(list))]
}

// standardName produces a syllable-assembled fantasy place name.
func (lg *LandmarkGenerator) standardName() string {
	return lg.pick(namePrefixes) + lg.pick(nameMiddles) + lg.pick(nameSuffixes)
}

// tavernName produces a name like "The Sleepy Badger".
func (lg *LandmarkGenerator) tavernName() string {
	return "The " + lg.pick(tavernAdjs) + " " + lg.pick(tavernNouns)
}

// ruinName produces "The Fallen {noun}" or "Old {standardName}".
func (lg *LandmarkGenerator) ruinName() string {
	if lg.rng.Intn(2) == 0 {
		return "The Fallen " + lg.pick(ruinNouns)
	}
	return "Old " + lg.pick(namePrefixes) + lg.pick(nameMiddles)
}

// monolithName produces "The {adj} Stone".
func (lg *LandmarkGenerator) monolithName() string {
	return "The " + lg.pick(monolithAdjs) + " Stone"
}

// nameFor generates a thematic name for the given landmark type.
func (lg *LandmarkGenerator) nameFor(lt LandmarkType) string {
	switch lt {
	case Tavern:
		return lg.tavernName()
	case Ruins:
		return lg.ruinName()
	case Monolith:
		return lg.monolithName()
	default:
		return lg.standardName()
	}
}

// ── placement logic ───────────────────────────────────────────────────

// adjacentToWater reports whether (x, y) is a land cell with at least
// one water neighbour (useful for lighthouses, bridges, etc.).
func adjacentToWater(x, y, w, h int, isLand func(int, int) bool) bool {
	for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		nx, ny := x+d[0], y+d[1]
		if nx >= 0 && nx < w && ny >= 0 && ny < h && !isLand(nx, ny) {
			return true
		}
	}
	return false
}

// farEnough returns true when (x, y) is at least minDist cells from
// every already-placed landmark (Manhattan distance).
func farEnough(x, y int, placed []Landmark, minDist int) bool {
	for _, lm := range placed {
		dx := x - lm.X
		dy := y - lm.Y
		if dx < 0 {
			dx = -dx
		}
		if dy < 0 {
			dy = -dy
		}
		if dx+dy < minDist {
			return false
		}
	}
	return true
}

// landmarkOrder defines the sequence of types to attempt placing.
// The first entry is always Town (the capital).
var landmarkOrder = []LandmarkType{
	Town, Village, Tavern, Library, Market, Shrine,
	Tower, Ruins, WitchHut, DragonLair,
	Cave, Mine,
	Bridge, Lighthouse,
	Monolith, Village,
}

// PlaceLandmarks scatters 8-15 named landmarks on land cells.
func (lg *LandmarkGenerator) PlaceLandmarks(width, height int, isLand func(x, y int) bool) []Landmark {
	// Determine target count: 8-15 scaled by map area.
	area := width * height
	count := 8 + lg.rng.Intn(8) // 8..15
	if area < 400 {
		count = 8
	} else if area > 2000 {
		count = 12 + lg.rng.Intn(4) // 12..15
	}

	const minDist = 3
	var placed []Landmark

	// Collect candidate land cells.
	type cell struct{ x, y int }
	var landCells []cell
	var waterEdgeCells []cell
	var interiorCells []cell

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !isLand(x, y) {
				continue
			}
			landCells = append(landCells, cell{x, y})
			if adjacentToWater(x, y, width, height, isLand) {
				waterEdgeCells = append(waterEdgeCells, cell{x, y})
			} else {
				interiorCells = append(interiorCells, cell{x, y})
			}
		}
	}

	if len(landCells) == 0 {
		return nil
	}

	// shuffle helpers
	shuffleCells := func(cs []cell) {
		lg.rng.Shuffle(len(cs), func(i, j int) { cs[i], cs[j] = cs[j], cs[i] })
	}
	shuffleCells(landCells)
	shuffleCells(waterEdgeCells)
	shuffleCells(interiorCells)

	// tryPlace picks the first suitable candidate from a list.
	tryPlace := func(candidates []cell) (int, int, bool) {
		for _, c := range candidates {
			if farEnough(c.x, c.y, placed, minDist) {
				return c.x, c.y, true
			}
		}
		return 0, 0, false
	}

	// preferredCandidates returns the ideal cell pool for a type.
	preferredCandidates := func(lt LandmarkType) []cell {
		switch lt {
		case Town, Village, Bridge:
			// Prefer near water (low elevation feel)
			return waterEdgeCells
		case Lighthouse:
			// Must be on water edge
			return waterEdgeCells
		case Tower, Ruins:
			// Prefer elevated / interior
			return interiorCells
		case Cave, Mine:
			// Prefer interior (mountain feel)
			return interiorCells
		default:
			return landCells
		}
	}

	typeIdx := 0
	for len(placed) < count && typeIdx <= len(landmarkOrder)+len(landCells) {
		var lt LandmarkType
		if typeIdx < len(landmarkOrder) {
			lt = landmarkOrder[typeIdx]
		} else {
			// Ran out of scripted order; pick random type.
			lt = LandmarkType(lg.rng.Intn(int(Market) + 1))
		}
		typeIdx++

		// Try preferred candidates first, then fall back to all land.
		x, y, ok := tryPlace(preferredCandidates(lt))
		if !ok {
			x, y, ok = tryPlace(landCells)
		}
		if !ok {
			continue // map too crowded
		}

		name := lg.nameFor(lt)
		placed = append(placed, Landmark{
			X:      x,
			Y:      y,
			Type:   lt,
			Name:   name,
			Symbol: lt.Symbol(),
			Icon:   landmarkIcons[lt],
		})
	}

	return placed
}
