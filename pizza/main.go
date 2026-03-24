package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ── lipgloss styles ─────────────────────────────────────────────────────────

var (
	headerBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 4).
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("229"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Italic(true)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("117")).
			Width(9)

	ingredientStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229"))

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)

	sectionHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("117"))

	noteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("251"))

	ratingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))

	courseStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Align(lipgloss.Center)

	divider = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("  " + strings.Repeat("━", 52))
)

const pizzaArt = `        _____
       /     \
      | () () |
      |  __   |
       \_____/`

func main() {
	n := flag.Int("n", 1, "number of pizzas to generate")
	seed := flag.Int64("seed", 0, "random seed (0 = random)")
	flag.Parse()

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	gen := NewRecipeGenerator(*seed)

	// Header
	fmt.Println()
	fmt.Println(headerBox.Render("🍕  INFINITE PIZZA GENERATOR  🍕"))
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(pizzaArt))
	fmt.Println()

	if *n > 1 {
		fmt.Println(courseStyle.Width(56).Render(
			fmt.Sprintf("— TASTING MENU: %d COURSES —", *n)))
		fmt.Println()
	}

	for i := 0; i < *n; i++ {
		recipe := gen.Generate()
		if *n > 1 {
			fmt.Println(courseStyle.Width(56).Render(
				fmt.Sprintf("· Course %d of %d ·", i+1, *n)))
			fmt.Println()
		}
		printRecipe(recipe)
		if i < *n-1 {
			fmt.Println()
		}
	}
}

func printRecipe(r Recipe) {
	fmt.Println(divider)
	fmt.Println()

	// Pizza name
	fmt.Printf("  %s\n\n", titleStyle.Render("\u201c"+r.Name+"\u201d"))

	// Ingredients
	printIngredient("BASE", r.Base)
	printIngredient("SAUCE", r.Sauce)
	for i, c := range r.Cheeses {
		lab := "CHEESE"
		if i > 0 {
			lab = ""
		}
		printIngredient(lab, c)
	}
	for i, p := range r.Proteins {
		lab := "PROTEIN"
		if i > 0 {
			lab = ""
		}
		printIngredient(lab, p)
	}
	for i, t := range r.Toppings {
		lab := "TOPPINGS"
		if i > 0 {
			lab = ""
		}
		printIngredient(lab, t)
	}
	printIngredient("GARNISH", r.Garnish)
	printIngredient("DRIZZLE", r.Drizzle)
	fmt.Println()

	// Tasting notes
	fmt.Println(divider)
	fmt.Println()
	fmt.Printf("  %s\n", sectionHeader.Render("TASTING NOTES"))
	printWrapped(r.TastingNote, 52, "  ")
	fmt.Println()

	// Pairing
	fmt.Printf("  %s\n", sectionHeader.Render("PAIRS WITH"))
	printWrapped(r.Pairing, 52, "  ")
	fmt.Println()

	// Chef quote
	fmt.Printf("  %s\n", sectionHeader.Render("CHEF'S WORD"))
	printWrapped(r.ChefQuote, 52, "  ")
	fmt.Println()

	// Rating
	fmt.Printf("  %s %s\n",
		sectionHeader.Render("Pretension Level:"),
		ratingStyle.Render(r.PretensionRating()))
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()
}

func printIngredient(label string, ing Ingredient) {
	fmt.Printf("  %s %s\n", labelStyle.Render(label), ingredientStyle.Render(ing.Name))
	fmt.Printf("  %s %s\n", strings.Repeat(" ", 9), descStyle.Render(ing.Description))
}

// printWrapped writes text word-wrapped to the given width with a prefix on every line.
func printWrapped(text string, width int, prefix string) {
	words := strings.Fields(text)
	line := prefix
	for _, w := range words {
		// +1 for the space before the word
		if len(line)+1+len(w) > width+len(prefix) && line != prefix {
			fmt.Println(line)
			line = prefix + w
		} else {
			if line == prefix {
				line += w
			} else {
				line += " " + w
			}
		}
	}
	if line != prefix {
		fmt.Println(line)
	}
}
