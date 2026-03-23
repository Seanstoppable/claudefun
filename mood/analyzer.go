package mood

import (
	"strings"
	"unicode"

	"github.com/ssmith/mood-octopus/octopus"
)

// MoodResult pairs an emotion with a confidence score between 0 and 1.
type MoodResult struct {
	Emotion    octopus.Emotion
	Confidence float64
}

// Analyzer detects emotions from text input using keyword matching.
type Analyzer struct {
	// keywordIndex maps each lowercase keyword to its emotion.
	keywordIndex map[string]octopus.Emotion
}

// NewAnalyzer builds the keyword lookup index from the emotion registry.
func NewAnalyzer() *Analyzer {
	idx := make(map[string]octopus.Emotion)
	for _, e := range octopus.AllEmotions() {
		for _, kw := range e.Info().Keywords {
			idx[kw] = e
		}
	}
	return &Analyzer{keywordIndex: idx}
}

// tokenize splits input into lowercase words, stripping punctuation.
func tokenize(input string) []string {
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return ' '
	}, strings.ToLower(input))
	return strings.Fields(cleaned)
}

// Analyze returns all detected emotions with confidence scores.
// Confidence is the fraction of matched keywords belonging to each emotion.
// Returns Curiosity with confidence 1.0 when no keywords match.
func (a *Analyzer) Analyze(input string) []MoodResult {
	words := tokenize(input)
	if len(words) == 0 {
		return []MoodResult{{Emotion: octopus.Curiosity, Confidence: 1.0}}
	}

	hits := make(map[octopus.Emotion]int)
	totalHits := 0

	for _, w := range words {
		if e, ok := a.keywordIndex[w]; ok {
			hits[e]++
			totalHits++
		}
	}

	if totalHits == 0 {
		return []MoodResult{{Emotion: octopus.Curiosity, Confidence: 1.0}}
	}

	results := make([]MoodResult, 0, len(hits))
	for e, count := range hits {
		results = append(results, MoodResult{
			Emotion:    e,
			Confidence: float64(count) / float64(totalHits),
		})
	}

	// Sort by confidence descending for stable output.
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Confidence > results[i].Confidence ||
				(results[j].Confidence == results[i].Confidence && results[j].Emotion < results[i].Emotion) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results
}

// DominantMood returns the single emotion with the highest confidence.
func (a *Analyzer) DominantMood(input string) octopus.Emotion {
	results := a.Analyze(input)
	return results[0].Emotion
}
