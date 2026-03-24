package main

import (
	"hash/fnv"
	"math"
	"strings"
	"unicode"
)

// Star represents a single star in the constellation.
type Star struct {
	X, Y       float64 // position in 0-1 normalized space
	Brightness float64 // 0-1 (brighter = more prominent)
	Size       float64 // 0-1 (maps to visual size)
	Char       rune    // the original character
	Index      int     // position in the input string
}

// StarMap holds the full constellation derived from an input sentence.
type StarMap struct {
	Stars  []Star
	Width  float64 // logical width (always 1.0)
	Height float64 // logical height (always 1.0)
	Seed   string  // the original input sentence
}

// GenerateStarMap maps an input sentence to a deterministic set of star positions.
func GenerateStarMap(input string) *StarMap {
	sm := &StarMap{
		Width:  1.0,
		Height: 1.0,
		Seed:   input,
	}

	if len(input) == 0 {
		return sm
	}

	runes := []rune(input)
	freq := charFrequency(runes)
	maxFreq := maxFrequencyValue(freq)

	// Generate raw star positions for letters and digits.
	for i, ch := range runes {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			continue
		}

		x, y := starPosition(ch, i, runes)
		brightness := starBrightness(ch, freq, maxFreq)
		size := starSize(ch)

		sm.Stars = append(sm.Stars, Star{
			X:          x,
			Y:          y,
			Brightness: brightness,
			Size:       size,
			Char:       ch,
			Index:      i,
		})
	}

	// Repulsion pass — push apart stars that are too close.
	applyRepulsion(sm.Stars, 0.08, 3)

	return sm
}

// starPosition computes a deterministic (x, y) in [0,1] for a character at a given index.
func starPosition(ch rune, index int, runes []rune) (float64, float64) {
	// Base angle from character code — spreads around a circle.
	baseAngle := float64(unicode.ToLower(ch)) * 2.4 // golden-angle-ish spacing

	// Radius grows with position index so later chars fan outward.
	totalChars := float64(len(runes))
	radius := 0.15 + 0.30*(float64(index)/math.Max(totalChars-1, 1))

	// Deterministic hash scatter using FNV.
	h := fnv.New64a()
	h.Write([]byte{byte(ch >> 8), byte(ch), byte(index >> 8), byte(index)})
	hash := h.Sum64()

	scatter := float64(hash%10000) / 10000.0
	angleJitter := (scatter - 0.5) * 1.2
	radiusJitter := (float64((hash/10000)%10000)/10000.0 - 0.5) * 0.15

	// Neighbour-based organic jitter.
	neighbourInfluence := 0.0
	if index > 0 {
		neighbourInfluence += float64(runes[index-1]) * 0.001
	}
	if index < len(runes)-1 {
		neighbourInfluence += float64(runes[index+1]) * 0.0007
	}

	angle := baseAngle + angleJitter + neighbourInfluence
	r := radius + radiusJitter

	// Convert polar → cartesian, centre at (0.5, 0.5), then clamp to [0.05, 0.95].
	x := clamp(0.5+r*math.Cos(angle), 0.05, 0.95)
	y := clamp(0.5+r*math.Sin(angle), 0.05, 0.95)
	return x, y
}

// charFrequency counts how often each lower-cased letter/digit appears.
func charFrequency(runes []rune) map[rune]int {
	freq := make(map[rune]int)
	for _, ch := range runes {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			freq[unicode.ToLower(ch)]++
		}
	}
	return freq
}

func maxFrequencyValue(freq map[rune]int) int {
	m := 1
	for _, v := range freq {
		if v > m {
			m = v
		}
	}
	return m
}

// starBrightness: more frequent characters are brighter; uppercase gets a boost.
func starBrightness(ch rune, freq map[rune]int, maxFreq int) float64 {
	base := float64(freq[unicode.ToLower(ch)]) / float64(maxFreq)
	// Scale into [0.3, 0.9] so even rare chars are visible.
	b := 0.3 + base*0.6
	if unicode.IsUpper(ch) {
		b = math.Min(b+0.1, 1.0)
	}
	return b
}

// starSize: vowels large, digits medium, consonants small.
func starSize(ch rune) float64 {
	if isVowel(ch) {
		return 0.8 + 0.2*float64(unicode.ToLower(ch)%5)/4.0 // 0.8-1.0
	}
	if unicode.IsDigit(ch) {
		return 0.5
	}
	return 0.2 + 0.2*float64(unicode.ToLower(ch)%7)/6.0 // 0.2-0.4
}

func isVowel(ch rune) bool {
	switch unicode.ToLower(ch) {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

// applyRepulsion pushes apart stars closer than minDist over several iterations.
func applyRepulsion(stars []Star, minDist float64, iterations int) {
	for iter := 0; iter < iterations; iter++ {
		for i := 0; i < len(stars); i++ {
			for j := i + 1; j < len(stars); j++ {
				dx := stars[j].X - stars[i].X
				dy := stars[j].Y - stars[i].Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist < minDist && dist > 1e-9 {
					// Push each star half the overlap distance apart.
					overlap := (minDist - dist) / 2.0
					nx := dx / dist
					ny := dy / dist
					stars[i].X = clamp(stars[i].X-nx*overlap, 0.05, 0.95)
					stars[i].Y = clamp(stars[i].Y-ny*overlap, 0.05, 0.95)
					stars[j].X = clamp(stars[j].X+nx*overlap, 0.05, 0.95)
					stars[j].Y = clamp(stars[j].Y+ny*overlap, 0.05, 0.95)
				}
			}
		}
	}
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

// ConstellationName generates a Latin-sounding celestial name from the input.
func (s *StarMap) ConstellationName() string {
	if s.Seed == "" {
		return "Vacuus"
	}

	suffixes := []string{"us", "ae", "um", "is", "ium", "ara", "onis", "oris"}

	// Collect only letters from the seed.
	var letters []rune
	for _, ch := range s.Seed {
		if unicode.IsLetter(ch) {
			letters = append(letters, ch)
		}
	}
	if len(letters) == 0 {
		return "Numeris"
	}

	// Take 3-5 characters as the stem.
	stemLen := 3
	if len(letters) >= 5 {
		stemLen = 5
	} else if len(letters) >= 4 {
		stemLen = 4
	}
	stem := string(letters[:stemLen])
	stem = strings.ToLower(stem)

	// If stem starts with a vowel, prefix a consonant cluster.
	if isVowel(rune(stem[0])) {
		prefixes := []string{"str", "cr", "pl", "fl", "tr"}
		// Deterministic pick based on the second character.
		pick := int(letters[0]) % len(prefixes)
		stem = prefixes[pick] + stem
	}

	// Pick suffix deterministically from a hash of the full input.
	h := fnv.New64a()
	h.Write([]byte(s.Seed))
	suffixIdx := int(h.Sum64() % uint64(len(suffixes)))
	suffix := suffixes[suffixIdx]

	// Drop trailing vowel from stem if suffix starts with a vowel to avoid awkward runs.
	if len(stem) > 0 && isVowel(rune(stem[len(stem)-1])) && isVowel(rune(suffix[0])) {
		stem = stem[:len(stem)-1]
	}

	// Capitalise.
	name := strings.ToUpper(stem[:1]) + stem[1:] + suffix
	return name
}
