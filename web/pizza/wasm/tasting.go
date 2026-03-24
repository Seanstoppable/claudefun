package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// TastingNoteGenerator produces sommelier-style tasting notes for absurd pizza combinations.
type TastingNoteGenerator struct {
	rng *rand.Rand
}

func NewTastingNoteGenerator(rng *rand.Rand) *TastingNoteGenerator {
	return &TastingNoteGenerator{rng: rng}
}

// GenerateNote creates a full tasting note for a pizza recipe.
func (t *TastingNoteGenerator) GenerateNote(ingredients []Ingredient) string {
	var parts []string
	parts = append(parts, t.opening(ingredients))
	parts = append(parts, t.flavorDescriptions(ingredients)...)
	parts = append(parts, t.closing())
	return strings.Join(parts, " ")
}

// GeneratePairing suggests an absurd beverage pairing.
func (t *TastingNoteGenerator) GeneratePairing() string {
	return pairings[t.rng.Intn(len(pairings))]
}

// GenerateChefQuote generates a fake quote from a fictional chef.
func (t *TastingNoteGenerator) GenerateChefQuote() string {
	name := t.chefName()
	restaurant := t.restaurantName()
	quote := chefQuoteTemplates[t.rng.Intn(len(chefQuoteTemplates))]
	return fmt.Sprintf("%s, %s: %s", name, restaurant, quote)
}

// --- internal helpers --------------------------------------------------------

func (t *TastingNoteGenerator) pick(pool []string) string {
	return pool[t.rng.Intn(len(pool))]
}

func (t *TastingNoteGenerator) pickIngredient(ingredients []Ingredient) Ingredient {
	return ingredients[t.rng.Intn(len(ingredients))]
}

func (t *TastingNoteGenerator) pickFlavor(ingredients []Ingredient) string {
	ing := t.pickIngredient(ingredients)
	if ing.Flavor != "" {
		return ing.Flavor
	}
	return fallbackFlavors[t.rng.Intn(len(fallbackFlavors))]
}

func (t *TastingNoteGenerator) opening(ingredients []Ingredient) string {
	idx := t.rng.Intn(len(openingTemplates))
	switch idx {
	case 0: // "A bold reimagining … [category]."
		cat := string(t.pickIngredient(ingredients).Category)
		if cat == "" {
			cat = t.pick(ingredientCategories)
		}
		return fmt.Sprintf(
			"A bold reimagining of the pizza form that challenges everything you thought you knew about %s.", cat)
	case 1:
		return "This pizza doesn't ask for permission. It arrives, and you are changed."
	case 2: // "At first glance … [ingredient]."
		return fmt.Sprintf(
			"At first glance, one might question the %s. By the third bite, one questions everything else.",
			t.pickIngredient(ingredients).Name)
	case 3: // "A symphony of [flavor1] and [flavor2]…"
		return fmt.Sprintf(
			"A symphony of %s and %s that somehow transcends the sum of its parts.",
			t.pickFlavor(ingredients), t.pickFlavor(ingredients))
	case 4:
		return "The culinary equivalent of a standing ovation."
	case 5:
		return "If Picasso made pizza, he would weep at this creation. Then he would order two."
	case 6:
		return "This is not merely a pizza. This is a thesis statement."
	case 7:
		return "Controversial? Perhaps. Revolutionary? Undeniably."
	case 8:
		return fmt.Sprintf(
			"Close your eyes. Take a bite. Open them. You're in %s.", t.pick(fancyPlaces))
	case 9:
		return "This pizza has more layers than your therapist could unpack."
	default:
		return "This pizza has more layers than your therapist could unpack."
	}
}

func (t *TastingNoteGenerator) flavorDescriptions(ingredients []Ingredient) []string {
	count := 2 + t.rng.Intn(2) // 2 or 3
	used := make(map[int]bool)
	var out []string

	for len(out) < count {
		idx := t.rng.Intn(len(flavorTemplates))
		if used[idx] {
			continue
		}
		used[idx] = true
		out = append(out, t.renderFlavorTemplate(idx, ingredients))
	}
	return out
}

func (t *TastingNoteGenerator) renderFlavorTemplate(idx int, ingredients []Ingredient) string {
	switch idx {
	case 0: // "The umami of the [ingredient] creates … with the [ingredient]."
		return fmt.Sprintf(
			"The umami of the %s creates an almost spiritual resonance with the %s.",
			t.pickIngredient(ingredients).Name, t.pickIngredient(ingredients).Name)
	case 1: // "Notes of [flavor] dance…"
		return fmt.Sprintf(
			"Notes of %s dance across the palate with reckless abandon.",
			t.pickFlavor(ingredients))
	case 2: // "The [ingredient] provides a textural counterpoint…"
		return fmt.Sprintf(
			"The %s provides a textural counterpoint that borders on the philosophical.",
			t.pickIngredient(ingredients).Name)
	case 3: // "One detects hints of [random thing]…"
		return fmt.Sprintf(
			"One detects hints of %s — or perhaps that's just the %s asserting dominance.",
			t.pick(randomThings), t.pickIngredient(ingredients).Name)
	case 4: // "The marriage of [ingredient] and [ingredient]…"
		return fmt.Sprintf(
			"The marriage of %s and %s is unexpected, yet inevitable in hindsight.",
			t.pickIngredient(ingredients).Name, t.pickIngredient(ingredients).Name)
	case 5: // "A whisper of [flavor] gives way…"
		return fmt.Sprintf(
			"A whisper of %s gives way to a crescendo of %s.",
			t.pickFlavor(ingredients), t.pickFlavor(ingredients))
	case 6: // "The [ingredient] brings a je ne sais quoi…"
		return fmt.Sprintf(
			"The %s brings a je ne sais quoi that translates roughly to 'what is happening in my mouth?'",
			t.pickIngredient(ingredients).Name)
	case 7: // "There's a minerality here…"
		return fmt.Sprintf(
			"There's a minerality here that suggests the %s has been on a personal journey.",
			t.pickIngredient(ingredients).Name)
	default:
		return fmt.Sprintf(
			"Notes of %s dance across the palate with reckless abandon.",
			t.pickFlavor(ingredients))
	}
}

func (t *TastingNoteGenerator) closing() string {
	return closingLines[t.rng.Intn(len(closingLines))]
}

func (t *TastingNoteGenerator) chefName() string {
	return t.pick(chefFirstNames) + " " + t.pick(chefLastNames)
}

func (t *TastingNoteGenerator) restaurantName() string {
	idx := t.rng.Intn(len(restaurantPatterns))
	switch idx {
	case 0:
		return "Maison du " + t.pick(restaurantNouns)
	case 1:
		return "The Gilded " + t.pick(restaurantNouns)
	case 2:
		return "Chez " + t.pick(chefFirstNames)
	case 3:
		return t.pick(restaurantNouns) + " & " + t.pick(restaurantNouns)
	case 4:
		return "Le Petit " + t.pick(restaurantNouns)
	default:
		return "The Gilded " + t.pick(restaurantNouns)
	}
}

// --- word pools --------------------------------------------------------------

// openingTemplates is only used for its length to pick an index; the actual
// text lives in the opening() switch for template interpolation.
var openingTemplates = [10]string{
	"category", "permission", "first glance", "symphony",
	"standing ovation", "picasso", "thesis", "controversial",
	"close your eyes", "layers",
}

var flavorTemplates = [8]string{
	"umami", "notes", "textural", "hints",
	"marriage", "whisper", "je ne sais quoi", "minerality",
}

var closingLines = []string{
	"Pairs magnificently with existential dread and a complete suspension of disbelief.",
	"Best enjoyed at 3 AM, preferably while questioning your life choices.",
	"Not for the faint of heart, nor the weak of stomach.",
	"A pizza that demands to be eaten alone, in the dark, with no witnesses.",
	"Would order again. Would not explain to friends.",
	"This pizza got a standing ovation from my taste buds and a restraining order from my cardiologist.",
	"Michelin stars are too pedestrian for this creation.",
}

var pairings = []string{
	"A 2019 Châteauneuf-du-Pape, served at exactly the temperature of regret",
	"Mountain Dew Code Red, decanted for 45 minutes",
	"A double espresso with oat milk and a single tear",
	"Chocolate milk, but make it pretentious — Valrhona single-origin cacao in Normandy whole milk",
	"Tap water, room temperature, to remind yourself of humility",
	"A craft IPA so hoppy it has its own passport",
	"Champagne, because you've earned it (you haven't, but the pizza has)",
	"Thai iced tea from a place that doesn't have a sign",
	"A kombucha that your friend Karen won't stop talking about",
	"Sparkling water with a lemon wedge, pretending it's enough",
	"Gas station coffee at 2 AM — the only honest pairing",
	"A smoothie made of hope and questionable protein powder",
}

var chefQuoteTemplates = []string{
	"\"I wept when I first tasted this. Then I wept again because I didn't think of it first.\"",
	"\"This pizza is the reason I went to culinary school. And also the reason I sometimes question it.\"",
	"\"Bold. Unhinged. Perfect.\"",
}

var chefFirstNames = []string{
	"Marco", "Valentina", "Jean-Pierre", "Yuki", "Björk",
	"Isabella", "Dmitri", "Fatima", "Alejandro", "Solange",
	"Thierry", "Kenji", "Anya", "Luciano", "Priya",
}

var chefLastNames = []string{
	"Fontaine", "Castellano", "Nakamura", "Volkov", "Okafor",
	"Delacroix", "Petrosyan", "Lindqvist", "Mendes", "Tanaka",
}

var restaurantPatterns = [5]string{
	"Maison du", "The Gilded", "Chez", "& combo", "Le Petit",
}

var restaurantNouns = []string{
	"Fromage", "Canard", "Truffle", "Radish", "Pigeon",
	"Anchovy", "Fennel", "Brioche", "Saffron", "Foie Gras",
}

var fancyPlaces = []string{
	"a small trattoria on the Amalfi Coast",
	"a Parisian rooftop at sunset",
	"a hidden izakaya in Shibuya",
	"a vineyard in Tuscany that doesn't appear on any map",
	"your grandmother's kitchen, but fancier",
	"a Michelin-starred food truck",
	"a monastery where monks have perfected the dough for centuries",
}

var randomThings = []string{
	"childhood nostalgia", "a summer thunderstorm", "velvet",
	"the color blue", "a forgotten dream", "pencil shavings",
	"old library books", "a jazz saxophone solo", "campfire smoke",
	"fresh laundry", "the ocean at dawn",
}

var fallbackFlavors = []string{
	"smoke", "brine", "earth", "citrus", "funk",
	"char", "sweetness", "acid", "butter", "spice",
}

var ingredientCategories = []string{
	"carbohydrates", "dairy", "fermentation", "cured meats",
	"nightshades", "stone fruits", "cruciferous vegetables",
	"umami", "dessert toppings", "condiments", "fungi",
}
