package main

import (
	"strings"
	"unicode"
)

// MoodResult pairs an emotion with a confidence score between 0 and 1.
type MoodResult struct {
	Emotion    Emotion
	Confidence float64
}

// Analyzer detects emotions from text input using keyword matching.
type Analyzer struct {
	keywordIndex map[string]Emotion
}

// NewAnalyzer builds the keyword lookup index from the emotion registry.
func NewAnalyzer() *Analyzer {
	idx := make(map[string]Emotion)
	for _, e := range AllEmotions() {
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
func (a *Analyzer) Analyze(input string) []MoodResult {
	words := tokenize(input)
	if len(words) == 0 {
		return []MoodResult{{Emotion: Curiosity, Confidence: 1.0}}
	}

	hits := make(map[Emotion]int)
	totalHits := 0

	for _, w := range words {
		if e, ok := a.keywordIndex[w]; ok {
			hits[e]++
			totalHits++
		}
	}

	if totalHits == 0 {
		return []MoodResult{{Emotion: Curiosity, Confidence: 1.0}}
	}

	results := make([]MoodResult, 0, len(hits))
	for e, count := range hits {
		results = append(results, MoodResult{
			Emotion:    e,
			Confidence: float64(count) / float64(totalHits),
		})
	}

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
