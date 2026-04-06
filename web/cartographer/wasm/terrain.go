package main

import (
	"hash/fnv"
	"math"
)

// Biome represents a terrain type on the fantasy map.
type Biome int

const (
	Ocean Biome = iota
	ShallowWater
	Beach
	Plains
	Forest
	DenseForest
	Hills
	Mountains
	SnowPeak
	Desert
	Swamp
	Tundra
)

func (b Biome) String() string {
	switch b {
	case Ocean:
		return "Ocean"
	case ShallowWater:
		return "ShallowWater"
	case Beach:
		return "Beach"
	case Plains:
		return "Plains"
	case Forest:
		return "Forest"
	case DenseForest:
		return "DenseForest"
	case Hills:
		return "Hills"
	case Mountains:
		return "Mountains"
	case SnowPeak:
		return "SnowPeak"
	case Desert:
		return "Desert"
	case Swamp:
		return "Swamp"
	case Tundra:
		return "Tundra"
	default:
		return "Unknown"
	}
}

func (b Biome) Symbol() string {
	switch b {
	case Ocean:
		return "≈"
	case ShallowWater:
		return "~"
	case Beach:
		return "∷"
	case Plains:
		return "·"
	case Forest:
		return "♣"
	case DenseForest:
		return "♠"
	case Hills:
		return "∧"
	case Mountains:
		return "▲"
	case SnowPeak:
		return "△"
	case Desert:
		return "∴"
	case Swamp:
		return "≋"
	case Tundra:
		return "∘"
	default:
		return "?"
	}
}

func (b Biome) Color() string {
	switch b {
	case Ocean:
		return "#1a5276"
	case ShallowWater:
		return "#2e86c1"
	case Beach:
		return "#f0d9b5"
	case Plains:
		return "#82b74b"
	case Forest:
		return "#2d6a2f"
	case DenseForest:
		return "#1a4314"
	case Hills:
		return "#a0785a"
	case Mountains:
		return "#6b6b6b"
	case SnowPeak:
		return "#e8e8e8"
	case Desert:
		return "#d4a843"
	case Swamp:
		return "#4a6b3a"
	case Tundra:
		return "#c8d6d6"
	default:
		return "#000000"
	}
}

// TerrainCell holds the properties of a single map cell.
type TerrainCell struct {
	Biome    Biome
	Height   float64 // 0-1 (sea level ~0.35)
	Moisture float64 // 0-1
	X, Y     int
}

// TerrainMap is a 2D grid of terrain cells generated from a place name.
type TerrainMap struct {
	Width, Height int
	Cells         [][]TerrainCell
	Seed          string
	SeaLevel      float64
}

// IsLand returns true if the cell at (x, y) is above sea level.
func (t *TerrainMap) IsLand(x, y int) bool {
	if x < 0 || x >= t.Width || y < 0 || y >= t.Height {
		return false
	}
	return t.Cells[y][x].Height >= t.SeaLevel
}

// BiomeAt returns the biome at the given coordinates.
func (t *TerrainMap) BiomeAt(x, y int) Biome {
	if x < 0 || x >= t.Width || y < 0 || y >= t.Height {
		return Ocean
	}
	return t.Cells[y][x].Biome
}

// --- deterministic hash-based noise ---

// hashNoise produces a deterministic float64 in [0,1) from integer coords and a seed.
func hashNoise(x, y int, seed uint64) float64 {
	h := fnv.New64a()
	buf := [24]byte{}
	putUint64(buf[0:8], uint64(x))
	putUint64(buf[8:16], uint64(y))
	putUint64(buf[16:24], seed)
	h.Write(buf[:])
	return float64(h.Sum64()&0x7FFFFFFFFFFFFFFF) / float64(0x7FFFFFFFFFFFFFFF)
}

func putUint64(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

// smoothNoise applies bilinear interpolation over the lattice hash values.
func smoothNoise(x, y float64, seed uint64) float64 {
	ix := int(math.Floor(x))
	iy := int(math.Floor(y))
	fx := x - float64(ix)
	fy := y - float64(iy)

	// hermite curve for smoother interpolation
	fx = fx * fx * (3 - 2*fx)
	fy = fy * fy * (3 - 2*fy)

	n00 := hashNoise(ix, iy, seed)
	n10 := hashNoise(ix+1, iy, seed)
	n01 := hashNoise(ix, iy+1, seed)
	n11 := hashNoise(ix+1, iy+1, seed)

	nx0 := n00*(1-fx) + n10*fx
	nx1 := n01*(1-fx) + n11*fx

	return nx0*(1-fy) + nx1*fy
}

// octaveNoise layers multiple frequencies of smooth noise for natural terrain.
func octaveNoise(x, y float64, seed uint64, octaves int) float64 {
	var total, maxAmp float64
	freq := 1.0
	amp := 1.0
	for i := 0; i < octaves; i++ {
		total += smoothNoise(x*freq, y*freq, seed+uint64(i)*7919) * amp
		maxAmp += amp
		freq *= 2.0
		amp *= 0.5
	}
	return total / maxAmp
}

// seedFromName creates a deterministic uint64 seed from a place name.
func seedFromName(name string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(name))
	return h.Sum64()
}

// GenerateTerrain builds a TerrainMap from a place name deterministically.
func GenerateTerrain(name string, width, height int) *TerrainMap {
	seed := seedFromName(name)
	seaLevel := 0.35

	tm := &TerrainMap{
		Width:    width,
		Height:   height,
		Cells:    make([][]TerrainCell, height),
		Seed:     name,
		SeaLevel: seaLevel,
	}

	// noise scale controls how zoomed-in the terrain features are
	scale := 0.05

	// Generate height and moisture, then assign biomes
	for y := 0; y < height; y++ {
		tm.Cells[y] = make([]TerrainCell, width)
		for x := 0; x < width; x++ {
			nx := float64(x) * scale
			ny := float64(y) * scale

			h := octaveNoise(nx, ny, seed, 4)
			m := octaveNoise(nx, ny, seed+123456, 4)

			// Island shaping: radial gradient that pushes edges toward ocean
			cx := float64(x)/float64(width) - 0.5
			cy := float64(y)/float64(height) - 0.5
			dist := math.Sqrt(cx*cx+cy*cy) * 2.0 // 0 at center, ~1.41 at corners
			falloff := 1.0 - math.Pow(clamp(dist, 0, 1), 2.0)
			h *= falloff

			h = clamp(h, 0, 1)
			m = clamp(m, 0, 1)

			biome := classifyBiome(h, m, seaLevel, y, height)

			tm.Cells[y][x] = TerrainCell{
				Biome:    biome,
				Height:   h,
				Moisture: m,
				X:        x,
				Y:        y,
			}
		}
	}

	return tm
}

func classifyBiome(height, moisture, seaLevel float64, y, mapHeight int) Biome {
	// Tundra at extreme latitudes for land cells
	edgeFraction := 0.12
	atEdge := float64(y) < float64(mapHeight)*edgeFraction ||
		float64(y) > float64(mapHeight)*(1-edgeFraction)

	if height < 0.30 {
		return Ocean
	}
	if height < seaLevel {
		return ShallowWater
	}
	if height < 0.38 {
		return Beach
	}

	// High-altitude biomes take priority
	if height >= 0.85 {
		return SnowPeak
	}
	if height >= 0.7 {
		return Mountains
	}
	if height >= 0.55 {
		return Hills
	}

	// Low-to-mid altitude: latitude can override to tundra
	if atEdge && moisture < 0.5 {
		return Tundra
	}

	if moisture < 0.25 {
		return Desert
	}
	if moisture < 0.45 {
		return Plains
	}
	if moisture < 0.65 {
		return Forest
	}
	// Very wet lowlands
	if moisture >= 0.65 {
		if height < 0.45 {
			return Swamp
		}
		return DenseForest
	}

	return Plains
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
