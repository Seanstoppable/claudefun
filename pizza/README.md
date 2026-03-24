# 🍕 Infinite Pizza Generator

```
        ___
       /   \
      | 🫒  |
     /  🍄  🧀 \
    | 🌶️  🍍  🥓 |
    |  🦐  🍫  🥑 |
     \  🧅  🍋  /
      \_________/
     /___________\
    (=============)
```

> _"Ah yes, the Deconstructed Sardine Rhapsody — a bold choice."_

**Infinite Pizza Generator** is a whimsical Go CLI that generates absurd pizza
recipes with confident sommelier-style tasting notes, pretentious names, and
deeply questionable beverage pairings. It takes pizza very seriously so you
don't have to.

## Install & Run

```bash
# Run directly
go run .

# Or build a binary
go build -o pizza-gen .
./pizza-gen
```

## CLI Flags

| Flag    | Description                    | Default |
|---------|--------------------------------|---------|
| `-n`    | Number of pizzas to generate   | `1`     |
| `-seed` | Seed for reproducible results  | random  |

```bash
# Generate 5 absurd pizzas
./pizza-gen -n 5

# Recreate that one masterpiece you saw earlier
./pizza-gen -seed 42
```

## Example Output

```
═══════════════════════════════════════════════════
  🍕  THE MELANCHOLY GORGONZOLA PARADOX  🍕
═══════════════════════════════════════════════════

  Base:     Charcoal-infused sourdough, extra thicc
  Sauce:    Mango-habanero reduction with a whisper of regret
  Cheese:   Triple-cream brie, crumbled feta, emergency mozzarella
  Toppings: Candied jalapeños, pickled watermelon rind,
            caramelized onion dust, one defiant anchovy

  🍷 TASTING NOTES
  "An audacious composition. The brie opens with
   a velvet curtain of umami, before the anchovy
   arrives uninvited — yet somehow welcome. The
   watermelon rind whispers of summer regrets."

  🥂 SUGGESTED PAIRING
  A lukewarm strawberry Yakult, served in a
  repurposed jam jar.
═══════════════════════════════════════════════════
```

## What's Inside

- **100+ ingredients** across 7 categories (bases, sauces, cheeses, proteins,
  vegetables, wildcards, and "why would you do this")
- **5 name generation strategies** — from faux-Italian to existential crisis
- **Procedural tasting notes** that sound disturbingly plausible
- **Absurd beverage pairings** that no sommelier would endorse

## Part of the Fun

This lives in the [`claudefun`](../) repo alongside other delightfully
pointless projects. Life is too short for boring side projects.

---

> _"In the grand pizzeria of life, we are all just toppings trying to find our
> crust."_ — Ancient Proverb (citation needed)
