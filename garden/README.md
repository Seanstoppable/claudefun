# 🌱 Code Garden

```
        🌱  C O D E   G A R D E N  🌱

    ψ ♣ 🦋 ♠ · · ┃ ♣ · 🌱 · ≈
    · ✿ ╂ ♠ · ♤ · · ψ · 🌱 ·
    · · 🐛 · ψ · · · · ψ · ·
    🌱 · · ♣ 🦋 · ┃ ♠ · ≈ · ·
    · · ψ · · ♤ · ψ 🌳 · · ·
    · ≈ · ψ · · · 🐛 · · ⚑ ·

    "Your garden is thriving!"
```

**Plant your codebase and watch it grow.** Code Garden turns source code into a whimsical ASCII garden — functions become plants, tests become fences, TODOs become literal bugs, and comments flutter by as butterflies.

---

## 🌿 Quick Start

```bash
# Run on any project directory
go run . /path/to/your/project

# Run on the current directory
go run .
```

### Install as a binary

```bash
go install github.com/ssmith/code-garden@latest
code-garden /path/to/your/project
```

---

## 🚩 CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--svg <file>` | Export garden as SVG image | _(none)_ |
| `--svg-width <px>` | SVG output width | `800` |
| `--svg-height <px>` | SVG output height | `600` |
| `--detailed` | Show per-file breakdown table | `false` |
| `--json` | Output raw stats as JSON | `false` |

```bash
# Export a pretty SVG
go run . --svg garden.svg --svg-width 1200 /path/to/project

# Get raw metrics as JSON
go run . --json /path/to/project

# See every file's stats
go run . --detailed /path/to/project
```

---

## 🌳 Garden Element Legend

### Plants (from your functions)

| Symbol | Name | Meaning |
|--------|------|---------|
| 🌱 | Seedling | Tiny functions (1–5 lines) |
| ψ | Sprout | Small functions (6–15 lines) |
| ♣ | Bush | Medium functions (16–30 lines) |
| ♠ | Tree | Large functions (31–60 lines) |
| 🌳 | Oak Tree | Huge functions (60+ lines) |
| ✿ | Flower | Type/struct definitions |
| ~ | Vine | Import/dependency |

### Elements (from your code quality)

| Symbol | Name | What triggers it |
|--------|------|------------------|
| · | Grass | Empty garden space |
| ┃ | Fence | Test functions protecting the garden |
| 🐛 | Bug | TODO, FIXME, HACK, XXX, or BUG markers |
| 🦋 | Butterfly | Comments drifting through the code |
| ⌇ | Weed | Untested or uncommented code |
| ● | Rock | Complexity hotspots (lots of if/for/switch) |
| ♤ | Mushroom | Type definitions popping up |
| ≈ | Pond | Breathing room between code sections |
| ⚑ | Gnome | Easter egg! Appears in very healthy gardens (90+) |
| ◎ | Tumbleweed | Rolls through abandoned code (health < 20) |
| ◈ | Beehive | Heavily imported modules (10+ imports) |
| ⌂ | Birdhouse | Well-documented functions (>30% comments) |
| ╂ | Scarecrow | Error handling code standing watch |

---

## 🌦️ Weather System

The garden's weather reflects your codebase health:

| Weather | Health Score | Condition |
|---------|-------------|-----------|
| 🌈 Rainbow | 80+ with tests | The gold standard! |
| ☀️ Sunny | 80+ | Clear skies, healthy code |
| ⛅ Partly Cloudy | 60–79 | Mostly good, some clouds |
| ☁️ Cloudy | 40–59 | Needs attention |
| 🌧️ Rainy | 20–39 | Storm's coming — refactor time |
| ⛈️ Stormy | Below 20 | Code emergency! |

---

## 📊 Health Score Breakdown

Your garden health is a score from 0 to 100, calculated from:

| Factor | Effect |
|--------|--------|
| Base score | Start at 70 |
| Test coverage (tests/functions > 20%) | +10 |
| Comment ratio (comments/lines > 10%) | +5 |
| Multiple languages | +5 |
| Few TODOs (< 5) | +10 |
| No tests at all | −15 |
| Each TODO/FIXME (up to 4) | −5 each |
| High avg complexity (> 10 per function) | −10 |

---

## 🌸 Seasons

The season is determined by your code's maturity:

| Season | Meaning |
|--------|---------|
| 🌸 Spring | Small, focused functions — fresh and growing |
| ☀️ Summer | Large codebase with many functions — in full bloom |
| 🍂 Autumn | High TODO ratio — leaves are falling |
| ❄️ Winter | Everything else — dormant, waiting for spring |

---

## 🎨 Example

```bash
$ go run . ~/projects/my-api

╭──────────────────────────────────────────────╮
│          🌱  my-api  🌱                       │
│                                              │
│  Spring 🌸 · Rainbow 🌈                      │
│                                              │
│  ψ ♣ 🦋 ♠ · · ┃ ♣ · 🌱 · ≈                 │
│  · ✿ ╂ ♠ · ♤ · · ψ · 🌱 ·                  │
│  · · · · ψ · · · · ψ · ·                    │
│  🌱 · · ♣ 🦋 · ┃ ♠ · ≈ · ·                 │
╰──────────────────────────────────────────────╯

  🌱 Garden Health: ████████░░ 80/100

  Rainbow 🌈 · Spring 🌸

  🌳 12 Trees · ✿ 3 Flowers · ┃ 5 Fences
  🐛 0 Bugs · 🦋 4 Butterflies · ● 1 Rocks

  📊 8 files · 342 lines · Go

  "Your garden is thriving! The oaks stand tall.
   The fences are sturdy."
```

---

## 🌐 Web Version

A WebAssembly version is available at [claudefun/garden](../web/garden/) — paste code into a textarea and watch it bloom in your browser. No file system access needed!

---

## 📋 Supported Languages

Go · Python · JavaScript · TypeScript · Java · Rust · C · C++ · Ruby · HTML · CSS

---

## 🌻 Contributing

Found a weed? Plant a PR! All gardeners welcome.

---

*"Every codebase is a garden. Some just need a little more weeding."* 🌱
