package main

import (
	"encoding/json"
	"syscall/js"
	"time"
)

func ingredientJSON(ing Ingredient) map[string]interface{} {
	return map[string]interface{}{
		"name":        ing.Name,
		"category":    string(ing.Category),
		"flavor":      ing.Flavor,
		"pretension":  ing.Pretension,
		"description": ing.Description,
	}
}

func ingredientsJSON(ings []Ingredient) []interface{} {
	out := make([]interface{}, len(ings))
	for i, ing := range ings {
		out[i] = ingredientJSON(ing)
	}
	return out
}

func generatePizza(_ js.Value, args []js.Value) interface{} {
	seed := time.Now().UnixNano()
	if len(args) > 0 {
		seed = int64(args[0].Float())
	}

	gen := NewRecipeGenerator(seed)
	recipe := gen.Generate()

	result := map[string]interface{}{
		"name":              recipe.Name,
		"base":              ingredientJSON(recipe.Base),
		"sauce":             ingredientJSON(recipe.Sauce),
		"cheeses":           ingredientsJSON(recipe.Cheeses),
		"proteins":          ingredientsJSON(recipe.Proteins),
		"toppings":          ingredientsJSON(recipe.Toppings),
		"garnish":           ingredientJSON(recipe.Garnish),
		"drizzle":           ingredientJSON(recipe.Drizzle),
		"tastingNote":       recipe.TastingNote,
		"pairing":           recipe.Pairing,
		"chefQuote":         recipe.ChefQuote,
		"pretensionRating":  recipe.PretensionRating(),
		"pretension":        recipe.Pretension,
		"seed":              seed,
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		js.Global().Get("console").Call("error", "JSON marshal error:", err.Error())
		return `{"error":"internal encoding error"}`
	}
	return string(jsonBytes)
}

func generateMenu(_ js.Value, args []js.Value) interface{} {
	count := 3
	if len(args) > 0 {
		count = int(args[0].Float())
	}
	if count < 1 {
		count = 1
	}
	if count > 10 {
		count = 10
	}

	baseSeed := time.Now().UnixNano()
	if len(args) > 1 {
		baseSeed = int64(args[1].Float())
	}

	gen := NewRecipeGenerator(baseSeed)
	pizzas := make([]interface{}, count)
	for i := 0; i < count; i++ {
		recipe := gen.Generate()
		pizzas[i] = map[string]interface{}{
			"name":              recipe.Name,
			"base":              ingredientJSON(recipe.Base),
			"sauce":             ingredientJSON(recipe.Sauce),
			"cheeses":           ingredientsJSON(recipe.Cheeses),
			"proteins":          ingredientsJSON(recipe.Proteins),
			"toppings":          ingredientsJSON(recipe.Toppings),
			"garnish":           ingredientJSON(recipe.Garnish),
			"drizzle":           ingredientJSON(recipe.Drizzle),
			"tastingNote":       recipe.TastingNote,
			"pairing":           recipe.Pairing,
			"chefQuote":         recipe.ChefQuote,
			"pretensionRating":  recipe.PretensionRating(),
			"pretension":        recipe.Pretension,
		}
	}

	result := map[string]interface{}{
		"seed":   baseSeed,
		"pizzas": pizzas,
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		js.Global().Get("console").Call("error", "JSON marshal error:", err.Error())
		return `{"error":"internal encoding error"}`
	}
	return string(jsonBytes)
}

func main() {
	js.Global().Set("generatePizza", js.FuncOf(generatePizza))
	js.Global().Set("generateMenu", js.FuncOf(generateMenu))

	select {}
}
