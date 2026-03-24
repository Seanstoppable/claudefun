package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// NameGenerator produces pretentious, artisanal pizza names.
type NameGenerator struct {
	rng *rand.Rand
}

// NewNameGenerator creates a NameGenerator with the given random source.
func NewNameGenerator(rng *rand.Rand) *NameGenerator {
	return &NameGenerator{rng: rng}
}

var (
	adjectives = []string{
		"Midnight", "Forbidden", "Velvet", "Savage", "Celestial",
		"Haunted", "Reckless", "Smoldering", "Whispering", "Defiant",
		"Volcanic", "Serpentine", "Golden", "Phantom", "Atomic",
	}

	dramaWords = []string{
		"Reckoning", "Awakening", "Prophecy", "Rebellion", "Symphony",
		"Paradox", "Rhapsody", "Inferno", "Odyssey", "Eclipse",
		"Tempest", "Revelation", "Requiem", "Crescendo", "Epiphany",
	}

	italianConnectors = []string{"e", "con", "alla"}

	frenchWords = []string{
		"Mystère", "Triomphe", "Sérénade", "Danse", "Rêve",
		"Catastrophe", "Fantaisie", "Passion", "Audace", "Folie",
	}

	tastingPhrases = []string{
		"A Study in Contrast", "Deconstructed Dreams",
		"The One That Got Away", "Controlled Chaos",
		"Beautiful Disaster", "Organized Madness", "Elegant Confusion",
	}

	fictionalNames = []string{
		"Giovanni", "Margherita", "Salvatore", "Francesca",
		"Dimitri", "Adelaide", "Bartholomew", "Clementine",
	}

	legacyWords = []string{
		"Emotion", "Legacy", "Secret", "Regret", "Folly", "Triumph", "Downfall",
	}
)

// Generate creates a pizza name, optionally incorporating ingredient names.
func (n *NameGenerator) Generate(ingredients []Ingredient) string {
	strategy := n.rng.Intn(5)
	switch strategy {
	case 0:
		return n.adjectiveIngredientDrama(ingredients)
	case 1:
		return n.italianPretension(ingredients)
	case 2:
		return n.frenchDe(ingredients)
	case 3:
		return n.tastingMenu()
	case 4:
		return n.fictionalPerson()
	default:
		return n.adjectiveIngredientDrama(ingredients)
	}
}

// Strategy 1: "The [Adjective] [Ingredient] [Drama]"
func (n *NameGenerator) adjectiveIngredientDrama(ingredients []Ingredient) string {
	adj := adjectives[n.rng.Intn(len(adjectives))]
	drama := dramaWords[n.rng.Intn(len(dramaWords))]
	ing := n.pickIngredient(ingredients)
	return fmt.Sprintf("The %s %s %s", adj, ing, drama)
}

// Strategy 2: "[Ingredient] e [Ingredient]" (Italian pretension)
func (n *NameGenerator) italianPretension(ingredients []Ingredient) string {
	connector := italianConnectors[n.rng.Intn(len(italianConnectors))]
	a := n.pickIngredient(ingredients)
	b := n.pickIngredientExcluding(ingredients, a)
	if connector == "alla" {
		return fmt.Sprintf("%s %s Nonna", a, connector)
	}
	return fmt.Sprintf("%s %s %s", a, connector, b)
}

// Strategy 3: "Le/La [French word] de [Ingredient]"
func (n *NameGenerator) frenchDe(ingredients []Ingredient) string {
	word := frenchWords[n.rng.Intn(len(frenchWords))]
	ing := n.pickIngredient(ingredients)
	article := "La"
	if n.rng.Intn(2) == 0 {
		article = "Le"
	}
	return fmt.Sprintf("%s %s de %s", article, word, ing)
}

// Strategy 4: "No. [random number]: [Dramatic phrase]"
func (n *NameGenerator) tastingMenu() string {
	num := n.rng.Intn(99) + 1
	phrase := tastingPhrases[n.rng.Intn(len(tastingPhrases))]
	return fmt.Sprintf("No. %d: %s", num, phrase)
}

// Strategy 5: "[Name]'s [Legacy word]"
func (n *NameGenerator) fictionalPerson() string {
	name := fictionalNames[n.rng.Intn(len(fictionalNames))]
	word := legacyWords[n.rng.Intn(len(legacyWords))]
	return fmt.Sprintf("%s's %s", name, word)
}

func (n *NameGenerator) pickIngredient(ingredients []Ingredient) string {
	if len(ingredients) > 0 {
		return titleCase(ingredients[n.rng.Intn(len(ingredients))].Name)
	}
	// Fallback absurd ingredients when none provided
	fallback := []string{
		"Gummy Bear", "Marshmallow", "Truffle Dust",
		"Dragon Fruit", "Unicorn Tear", "Prosciutto",
		"Dark Chocolate", "Ghost Pepper", "Gold Leaf",
	}
	return fallback[n.rng.Intn(len(fallback))]
}

func (n *NameGenerator) pickIngredientExcluding(ingredients []Ingredient, exclude string) string {
	if len(ingredients) > 1 {
		for attempts := 0; attempts < 10; attempts++ {
			pick := titleCase(ingredients[n.rng.Intn(len(ingredients))].Name)
			if pick != exclude {
				return pick
			}
		}
	}
	// Fallback: return a different absurd ingredient
	fallback := []string{
		"Burrata", "Fig Jam", "Caviar", "Honeycomb",
		"Balsamic Reduction", "Edible Flowers", "Wagyu",
	}
	return fallback[n.rng.Intn(len(fallback))]
}

func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
