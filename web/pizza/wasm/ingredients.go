package main

import "math/rand"

// IngredientCategory represents a layer of the pizza construction.
type IngredientCategory string

const (
	Base    IngredientCategory = "base"
	Sauce   IngredientCategory = "sauce"
	Cheese  IngredientCategory = "cheese"
	Protein IngredientCategory = "protein"
	Topping IngredientCategory = "topping"
	Garnish IngredientCategory = "garnish"
	Drizzle IngredientCategory = "drizzle"
)

// Ingredient is a single pizza component with sommelier-grade metadata.
type Ingredient struct {
	Name        string
	Category    IngredientCategory
	Flavor      string // e.g., "umami", "sweet", "tangy", "chaotic"
	Pretension  int    // 1-5 scale of how pretentious this sounds
	Description string // short sommelier-esque descriptor
}

// AllIngredients returns the full ingredient database keyed by category.
func AllIngredients() map[IngredientCategory][]Ingredient {
	return map[IngredientCategory][]Ingredient{
		Base:    bases(),
		Sauce:   sauces(),
		Cheese:  cheeses(),
		Protein: proteins(),
		Topping: toppings(),
		Garnish: garnishes(),
		Drizzle: drizzles(),
	}
}

// RandomFrom picks a random ingredient from the given category.
func RandomFrom(category IngredientCategory, rng *rand.Rand) Ingredient {
	all := AllIngredients()
	items := all[category]
	return items[rng.Intn(len(items))]
}

func bases() []Ingredient {
	return []Ingredient{
		{"Neapolitan Hand-Stretched", Base, "classic", 3, "Whispers of volcanic Campanian soil in every leopard-spotted char bubble"},
		{"Sourdough 72-Hour Ferment", Base, "tangy", 4, "Three days of existential yeast contemplation, yielding a crumb that weeps with complexity"},
		{"Cauliflower Deconstructed", Base, "earthy", 3, "A brassica manifesto pressed into service, structurally aspirational"},
		{"Deep-Dish Chicago Fortress", Base, "buttery", 2, "An architectural buttress of dough engineered to contain a soup's worth of ambition"},
		{"NY Thin Crisp", Base, "classic", 2, "Foldable, unapologetic, and thinner than your landlord's patience"},
		{"Pretzel Crust", Base, "salty", 2, "Bavarian heritage twisted into circular submission with a lye-kissed attitude"},
		{"Waffle Iron-Pressed", Base, "sweet", 1, "Grid-embossed chaos that holds syrup and marinara with equal conviction"},
		{"Croissant Laminated", Base, "buttery", 5, "128 layers of Franco-Italian diplomacy, each one flakier than your last excuse"},
		{"Rice Paper Translucent", Base, "delicate", 4, "A ghostly veil of carbohydrate, more concept than crust"},
		{"Naan-Inspired", Base, "smoky", 2, "Tandoor-kissed and pillowy, harboring quiet garlic secrets"},
		{"Cornbread Rustic", Base, "sweet", 1, "Southern hospitality in structural form, crumbly yet unyielding"},
		{"Tortilla Flash-Baked", Base, "neutral", 1, "A thin-crust speedrun, unapologetically efficient"},
		{"Bagel Dense Ring", Base, "malty", 2, "Boiled then baked, impossibly chewy, spiritually from the Lower East Side"},
		{"Lavender Shortbread", Base, "floral", 5, "A Provençal fever dream masquerading as pizza infrastructure"},
		{"Activated Charcoal Midnight", Base, "chaotic", 4, "Void-black and Instagram-ready, tastes like wellness culture"},
		{"Potato Rösti Foundation", Base, "earthy", 3, "Shredded, golden, and structurally dubious — a Swiss engineer's nightmare"},
	}
}

func sauces() []Ingredient {
	return []Ingredient{
		{"San Marzano Crush", Sauce, "tangy", 3, "DOP-certified tomatoes crushed by hand, probably while someone played mandolin"},
		{"White Truffle Béchamel", Sauce, "umami", 5, "A velvet roux elevated to aristocratic heights by fungal opulence"},
		{"Sriracha-Honey Fusion", Sauce, "spicy-sweet", 2, "The duality of man, bottled and drizzled with reckless abandon"},
		{"Pesto Genovese", Sauce, "herbaceous", 3, "Basil, pine nuts, and Parmigiano in a mortar-pounded love triangle"},
		{"BBQ Bourbon Reduction", Sauce, "smoky", 2, "Hickory smoke and Kentucky courage simmered past the point of sobriety"},
		{"Alfredo Cloud", Sauce, "rich", 3, "Butter and cream in a codependent relationship, thickened by denial"},
		{"Salsa Verde Taqueria", Sauce, "tangy", 2, "Tomatillo and jalapeño arguing in a blender at 2 AM"},
		{"Curry Tikka Masala", Sauce, "spicy", 2, "A British-Indian détente of tomato, cream, and aromatic confusion"},
		{"Peanut Satay Drizzle", Sauce, "nutty", 3, "Southeast Asian street wisdom liquefied into spreadable swagger"},
		{"Ranch Elevated", Sauce, "tangy", 1, "Hidden Valley sent to finishing school and given a tiny herb garden"},
		{"Cranberry Compote", Sauce, "tart", 3, "Thanksgiving's revenge, served year-round without remorse"},
		{"Miso-Ginger Glaze", Sauce, "umami", 4, "Fermented soybean enlightenment meets rhizome heat in quiet harmony"},
		{"Nutella Hazelnut Spread", Sauce, "sweet", 1, "Italy's most exported coping mechanism, now on pizza"},
		{"Blue Cheese Fondue", Sauce, "pungent", 4, "Liquid Roquefort courage for those who fear nothing"},
		{"Hot Honey Cascade", Sauce, "spicy-sweet", 3, "Capsaicin-infused nectar falling in slow, cinematic rivulets"},
		{"Gravy Sunday Roast", Sauce, "savory", 1, "Your nan's gravy, repurposed with zero regret and maximum comfort"},
		{"Marshmallow Fluff Base", Sauce, "sweet", 1, "A structural atrocity that tastes like childhood and poor decisions"},
	}
}

func cheeses() []Ingredient {
	return []Ingredient{
		{"Fresh Mozzarella di Bufala", Cheese, "milky", 4, "Hand-pulled from water buffalo milk by someone who takes eye contact seriously"},
		{"Aged Gruyère Cave-Ripened", Cheese, "nutty", 5, "Twelve months in a Swiss cave, emerging with the gravitas of a philosophy professor"},
		{"Vegan Cashew Cream", Cheese, "neutral", 3, "Nuts, soaked and blended into something that believes deeply in itself"},
		{"Kraft Singles Ironic", Cheese, "processed", 1, "Individually wrapped American exceptionalism, deployed with postmodern sincerity"},
		{"Smoked Gouda", Cheese, "smoky", 3, "Dutch cheese that survived a campfire and came out more interesting"},
		{"Crumbled Gorgonzola", Cheese, "pungent", 4, "Italian blue that crumbles like a dynasty and hits twice as hard"},
		{"Halloumi Grilled Slabs", Cheese, "salty", 3, "Cypriot squeaky cheese that refuses to melt, a contrarian on principle"},
		{"Cream Cheese Schmear", Cheese, "tangy", 1, "New York deli energy in thick, unapologetic dollops"},
		{"Cheez Whiz Artisanal", Cheese, "chaotic", 1, "Pressurized dairy product rebranded with a straight face and a farm logo"},
		{"Manchego 12-Month", Cheese, "nutty", 4, "Sheep's milk aged to develop opinions about wine pairings"},
		{"Ricotta Clouds", Cheese, "milky", 3, "Soft, white, and dropped in dreamy spoonfuls like dairy cumulus"},
		{"Brie Wheel Melted", Cheese, "buttery", 4, "An entire wheel, sacrificed to heat, oozing with Gallic indifference"},
		{"Cottage Cheese Controversial", Cheese, "tangy", 1, "Lumpy, divisive, and fully aware of the discourse it generates"},
		{"Pepper Jack Lava", Cheese, "spicy", 2, "Monterey Jack radicalized by jalapeño fragments into something dangerous"},
		{"String Cheese Pulled", Cheese, "milky", 1, "Peeled strand by strand in a meditative act of snack deconstruction"},
		{"Parmesan Reggiano Shaved", Cheese, "umami", 4, "24-month crystals of glutamic acid, each flake a tiny flavor grenade"},
	}
}

func proteins() []Ingredient {
	return []Ingredient{
		{"Pepperoni Classic", Protein, "spicy", 1, "Crimson discs that curl into tiny grease chalices of pure joy"},
		{"Prosciutto di Parma", Protein, "salty", 5, "Air-dried ham so thin it's practically a pork rumor"},
		{"Smoked Salmon Lox", Protein, "briny", 4, "Silky cured fish draped over dough like an expensive silk negligée"},
		{"Fried Chicken Tenders", Protein, "savory", 1, "Breaded, golden, and absolutely unaware they've been promoted to pizza duty"},
		{"Tofu Silk Cubes", Protein, "neutral", 3, "Pressed soybean potential, absorbing the personality of everything around it"},
		{"Bacon Maple-Glazed", Protein, "sweet-smoky", 2, "Pork belly strips lacquered in tree sap, the breakfast-dinner treaty"},
		{"Anchovy Umami Bombs", Protein, "umami", 3, "Tiny salt-cured fish delivering weaponized savory from the Cantabrian deep"},
		{"Chorizo Iberian", Protein, "spicy", 3, "Paprika-stained pork with the swagger of a flamenco dancer"},
		{"Lamb Merguez", Protein, "spicy", 4, "North African lamb sausage radiating cumin and harissa conviction"},
		{"Shrimp Tempura", Protein, "delicate", 3, "Crustaceans in crispy batter, a textural plot twist on every slice"},
		{"Pulled Pork Shoulder", Protein, "smoky", 2, "Low-and-slow smoked pork, shredded like your last relationship"},
		{"Gummy Bears Protein-Adjacent", Protein, "sweet", 1, "Gelatin-based wildlife technically containing amino acids, spiritually containing chaos"},
		{"Vienna Sausage Retro", Protein, "salty", 1, "Canned cylindrical nostalgia from the back of your grandfather's pantry"},
		{"Cricket Flour Crumble", Protein, "nutty", 4, "Sustainable entomophagy sprinkled with the confidence of a TED talk"},
		{"Spam Musubi-Style", Protein, "salty", 2, "Canned pork legend, sliced and seared with Hawaiian-Japanese reverence"},
		{"Egg Sunny-Side-Up", Protein, "rich", 3, "A wobbling golden eye staring back at you from a sea of cheese"},
	}
}

func toppings() []Ingredient {
	return []Ingredient{
		{"Pineapple Controversial", Topping, "sweet", 1, "The fruit that launched a thousand arguments and zero resolutions"},
		{"Jalapeño Fire Rings", Topping, "spicy", 1, "Green circles of Scoville-rated menace, sliced with surgical intent"},
		{"Black Olive Mediterranean", Topping, "briny", 2, "Kalamata darkness scattered like punctuation across a cheesy sentence"},
		{"Sun-Dried Tomato Intense", Topping, "tangy", 3, "Tomatoes that stared at the sun and emerged leathery and concentrated"},
		{"Caramelized Onion Silk", Topping, "sweet", 3, "Two hours of patience rewarded with jammy, translucent allium ribbons"},
		{"Roasted Garlic Cloves", Topping, "savory", 3, "Entire heads roasted until they spread like butter and repel vampires"},
		{"Artichoke Hearts", Topping, "tangy", 3, "The tender core of a thistle, sacrificed for your Mediterranean fantasy"},
		{"Mushroom Wild Foraged", Topping, "earthy", 4, "Chanterelles and porcini gathered at dawn by someone in linen pants"},
		{"Arugula Peppery", Topping, "bitter", 3, "Rocket leaves added post-bake with the confidence of a final garnish"},
		{"Corn Sweet Kernels", Topping, "sweet", 1, "Golden niblets that have no idea they're on a pizza and don't care"},
		{"Pickled Red Onion", Topping, "tangy", 2, "Crimson crescents quick-pickled into electric pink submission"},
		{"French Fries Loaded", Topping, "salty", 1, "Carbs on carbs: an act of defiance against the food pyramid"},
		{"Doritos Crushed", Topping, "chaotic", 1, "Nacho cheese triangles reduced to rubble and deployed without shame"},
		{"Marshmallows Torched", Topping, "sweet", 1, "Campfire nostalgia brûléed onto pizza by someone who simply does not care"},
		{"Edible Flowers", Topping, "floral", 5, "Violas and nasturtiums arranged by someone who journals in cursive"},
		{"Banana Sliced", Topping, "sweet", 1, "Cavendish rounds that polarize harder than pineapple ever could"},
		{"Pop-Tart Crumbled", Topping, "sweet", 1, "Frosted pastry shrapnel scattered across cheese like breakfast shrapnel"},
		{"Kimchi Fermented", Topping, "spicy-tangy", 3, "Napa cabbage fermented into pungent, probiotic glory"},
		{"Truffle Shavings", Topping, "earthy", 5, "Paper-thin fungal currency, each shaving worth more than the pizza beneath it"},
		{"Rainbow Sprinkles", Topping, "sweet", 1, "Technicolor sugar pellets declaring that rules are a social construct"},
	}
}

func garnishes() []Ingredient {
	return []Ingredient{
		{"Fresh Basil Chiffonade", Garnish, "herbaceous", 3, "Ribboned leaves releasing volatile oils and Neapolitan street cred"},
		{"Microgreens Pretentious", Garnish, "earthy", 5, "Infant vegetables harvested before they could develop a personality"},
		{"Edible Gold Leaf", Garnish, "neutral", 5, "Actual gold, contributing nothing to flavor and everything to your ego"},
		{"Crushed Red Pepper", Garnish, "spicy", 1, "Shaker-dispensed heat from every pizzeria counter since 1953"},
		{"Everything Bagel Seasoning", Garnish, "savory", 2, "Sesame, poppy, garlic, onion, and salt in a commitment-free blend"},
		{"Maldon Sea Salt Flakes", Garnish, "salty", 4, "Pyramidal crystals hand-harvested from Essex waters with geometric precision"},
		{"Lemon Zest Curls", Garnish, "bright", 3, "Citrus microplaned into aromatic confetti that makes everything sing"},
		{"Fresh Dill Fronds", Garnish, "herbaceous", 3, "Feathery Scandinavian herb energy, whispering of gravlax and open sandwiches"},
		{"Toasted Sesame Seeds", Garnish, "nutty", 2, "Tiny golden seeds that punch above their weight in the aroma department"},
		{"Candy Corn Seasonal", Garnish, "sweet", 1, "Waxy tricolor abominations deployed exclusively to upset purists"},
		{"Gummy Worm Arrangement", Garnish, "sweet", 1, "Gelatin invertebrates draped artfully across the pie by a chaos agent"},
		{"Cotton Candy Wisps", Garnish, "sweet", 2, "Spun sugar dissolving on contact, a fleeting garnish for the philosophical"},
		{"Ranch Dust", Garnish, "tangy", 1, "Dehydrated ranch packet sprinkled like the spice of suburban dreams"},
		{"Tajín Rim", Garnish, "tangy-spicy", 2, "Chili-lime salt rimming each slice like a Michelada's pizza cousin"},
		{"Balsamic Pearls", Garnish, "tangy", 5, "Spherified vinegar orbs that burst with the drama of molecular gastronomy"},
		{"Fennel Pollen", Garnish, "floral", 5, "The most expensive spice you've never heard of, dusted by a trembling hand"},
	}
}

func drizzles() []Ingredient {
	return []Ingredient{
		{"Extra Virgin Olive Oil", Drizzle, "fruity", 3, "Cold-pressed Tuscan liquid gold, finishing every pie with a silken shimmer"},
		{"White Truffle Oil", Drizzle, "earthy", 5, "A few drops that cost more than the pizza and smell like forest floor royalty"},
		{"Hot Honey", Drizzle, "spicy-sweet", 3, "Inferno-kissed nectar that makes cheese weep and taste buds cheer"},
		{"Balsamic Reduction", Drizzle, "tangy", 4, "Vinegar simmered into syrupy darkness, zigzagged with sommelier precision"},
		{"Ranch Dressing", Drizzle, "tangy", 1, "America's condiment, liberated from salads and applied with zero restraint"},
		{"Sriracha Mayo", Drizzle, "spicy", 1, "Rooster sauce and mayo in a harmonious union of heat and cream"},
		{"Nutella Drizzle", Drizzle, "sweet", 1, "Hazelnut chocolate lava flowing across pizza like a dessert insurgency"},
		{"Condensed Milk", Drizzle, "sweet", 2, "Viscous sweetened dairy poured in spirals of deliberate, sticky madness"},
		{"Fish Sauce Caramel", Drizzle, "umami-sweet", 4, "Vietnamese-inspired alchemy turning fermented anchovies into liquid obsession"},
		{"Maple Syrup", Drizzle, "sweet", 2, "Canadian tree blood drizzled to blur the line between breakfast and dinner"},
		{"Chili Crisp Oil", Drizzle, "spicy", 3, "Lao Gan Ma's crunchy, oily gift to humanity, one spoonful at a time"},
		{"Blue Cheese Dressing", Drizzle, "pungent", 2, "Chunky, divisive, and aggressively aromatic — the wing night crossover"},
		{"Melted Butter Garlic", Drizzle, "rich", 2, "Clarified courage infused with allium, pooling in golden rivulets"},
		{"Tahini Stream", Drizzle, "nutty", 3, "Sesame paste thinned into a pale, earthy ribbon of Middle Eastern wisdom"},
		{"Chocolate Ganache", Drizzle, "sweet", 3, "Dark chocolate and cream melted into decadent, pizza-questioning rivers"},
		{"Liquid Smoke Essence", Drizzle, "smoky", 2, "Condensed campfire in a bottle, deployed in terrifying micro-drops"},
	}
}
