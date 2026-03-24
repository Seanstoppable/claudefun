package main

import "math/rand"

// Recipe represents a complete, absurdly generated pizza.
type Recipe struct {
	Name        string
	Base        Ingredient
	Sauce       Ingredient
	Cheeses     []Ingredient
	Proteins    []Ingredient
	Toppings    []Ingredient
	Garnish     Ingredient
	Drizzle     Ingredient
	TastingNote string
	Pairing     string
	ChefQuote   string
	Pretension  float64
}

// RecipeGenerator assembles complete pizza recipes from random ingredients.
type RecipeGenerator struct {
	rng    *rand.Rand
	namer  *NameGenerator
	taster *TastingNoteGenerator
}

// NewRecipeGenerator creates a RecipeGenerator seeded for reproducible chaos.
func NewRecipeGenerator(seed int64) *RecipeGenerator {
	rng := rand.New(rand.NewSource(seed))
	return &RecipeGenerator{
		rng:    rng,
		namer:  NewNameGenerator(rng),
		taster: NewTastingNoteGenerator(rng),
	}
}

// Generate assembles one complete pizza recipe.
func (r *RecipeGenerator) Generate() Recipe {
	rec := Recipe{
		Base:    RandomFrom(Base, r.rng),
		Sauce:   RandomFrom(Sauce, r.rng),
		Garnish: RandomFrom(Garnish, r.rng),
		Drizzle: RandomFrom(Drizzle, r.rng),
	}

	// 1-2 cheeses
	numCheeses := 1 + r.rng.Intn(2)
	for i := 0; i < numCheeses; i++ {
		rec.Cheeses = append(rec.Cheeses, RandomFrom(Cheese, r.rng))
	}

	// 1-2 proteins
	numProteins := 1 + r.rng.Intn(2)
	for i := 0; i < numProteins; i++ {
		rec.Proteins = append(rec.Proteins, RandomFrom(Protein, r.rng))
	}

	// 2-4 toppings
	numToppings := 2 + r.rng.Intn(3)
	for i := 0; i < numToppings; i++ {
		rec.Toppings = append(rec.Toppings, RandomFrom(Topping, r.rng))
	}

	all := rec.AllIngredients()
	rec.Name = r.namer.Generate(all)
	rec.TastingNote = r.taster.GenerateNote(all)
	rec.Pairing = r.taster.GeneratePairing()
	rec.ChefQuote = r.taster.GenerateChefQuote()
	rec.Pretension = rec.avgPretension()

	return rec
}

// AllIngredients returns a flattened list of every ingredient in the recipe.
func (r *Recipe) AllIngredients() []Ingredient {
	out := []Ingredient{r.Base, r.Sauce}
	out = append(out, r.Cheeses...)
	out = append(out, r.Proteins...)
	out = append(out, r.Toppings...)
	out = append(out, r.Garnish, r.Drizzle)
	return out
}

// PretensionRating returns a fun emoji rating based on the average pretension score.
func (r *Recipe) PretensionRating() string {
	p := r.Pretension
	switch {
	case p >= 4.5:
		return "🍕🍕🍕🍕🍕 Transcendent (you need a second mortgage)"
	case p >= 4.0:
		return "🍕🍕🍕🍕 Insufferable (the waiter judges you)"
	case p >= 3.0:
		return "🍕🍕🍕 Pretentious (requires a reservation)"
	case p >= 2.0:
		return "🍕🍕 Elevated (Instagram-worthy)"
	default:
		return "🍕 Approachable (your mom would eat this)"
	}
}

func (r *Recipe) avgPretension() float64 {
	all := r.AllIngredients()
	if len(all) == 0 {
		return 0
	}
	sum := 0
	for _, ing := range all {
		sum += ing.Pretension
	}
	return float64(sum) / float64(len(all))
}
