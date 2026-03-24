# Constellation Composer

```
                          ✦
                 ✧ ─ ─ ─ ─ ─ ─ ★
                ╱                 ╲
          ✧    ╱    ·              ✦
           ╲  ╱         ·        ╱
            ★        ✧     ✧   ╱
           ╱  ╲      │   ╱    ╱
     ·    ╱    ╲     │  ╱    ╱
         ✦      ★ ─ ─ ★    ╱    ·
          ╲           │    ╱
           ╲    ·     │   ╱
            ╲         │  ╱
     ✧       ★ ─ ─ ─ ✦ ╱
                        ★        ·
                  ·          ✧
```

*Every sentence hides a constellation. This tool finds it.*

**Constellation Composer** is a Go CLI that transforms plain text into
deterministic star maps — complete with connected constellation lines,
beautiful terminal rendering, SVG export, and procedurally generated
mythology from fictional cultures who have been watching your stars for
centuries.

The same words will always produce the same sky.

---

## ✦ Quick Start

```bash
go run . "the owl stretches its wings at dusk"
```

A constellation blooms in your terminal: stars placed by your letters,
connected by minimum spanning trees with artistic flourishes, named in
the style of ancient Latin star catalogs, and accompanied by a myth
about how the stars came to be.

---

## ★ Installation

```bash
git clone https://github.com/ssmith/constellation-composer
cd constellation-composer
go build -o constellation-composer .
```

Or simply run directly:

```bash
go run . "your sentence here"
```

---

## ✧ CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-svg` | `false` | Export constellation as an SVG file |
| `-width` | `80` | Terminal output width (columns) |
| `-height` | `24` | Terminal output height (rows) |
| `-svg-width` | `800` | SVG canvas width (px) |
| `-svg-height` | `600` | SVG canvas height (px) |
| `-no-myth` | `false` | Skip the mythology backstory |
| `-minimal` | `false` | Stars only — no lines, no myth, no name |

```bash
# Render wide in the terminal
go run . -width 120 -height 40 "a cathedral of light"

# Export glowing SVG
go run . -svg -svg-width 1200 "the river remembers"

# Just the stars, nothing else
go run . -minimal "silence"
```

---

## ✦ How It Works

Every constellation follows the same celestial pipeline:

```
  "the owl stretches its wings at dusk"
                  │
          ┌───────▼────────┐
          │   Text → Stars  │   Extract letters & digits.
          │                 │   Each character gets a position via
          │   starmap.go    │   golden-angle placement, FNV-64a hashing,
          │                 │   neighbour-influenced jitter, and a
          │                 │   repulsion pass to avoid overlap.
          └───────┬────────┘
                  │
          ┌───────▼────────┐
          │  Star → Lines   │   Prim's MST connects the stars.
          │                 │   Then 1–2 artistic edges are added
          │  connect.go     │   (short, deterministic) to create
          │                 │   triangles and visual depth.
          └───────┬────────┘
                  │
       ┌──────────┼──────────┐
       ▼                     ▼
 ┌───────────┐       ┌────────────┐
 │  Terminal  │       │    SVG     │
 │            │       │            │
 │ ✦ ★ ✧ ·   │       │  Glow fx   │
 │ Bresenham  │       │  Gradient  │
 │ lines      │       │  Spikes    │
 │            │       │  #0a0a2e   │
 └───────────┘       └────────────┘
       │                     │
       └──────────┬──────────┘
                  ▼
          ┌───────────────┐
          │   Mythology    │   SHA-256 seeded storytelling.
          │                │   A fictional culture, a name,
          │  mythology.go  │   a conflict, a transformation
          │                │   into starlight.
          └───────────────┘
```

**Star brightness** is driven by character frequency and case — uppercase
letters burn brighter. **Star size** follows phonetics: vowels are large
and luminous, consonants are compact, digits are faint sparks.

---

## ★ Example

```bash
$ go run . "the fox who learned to fly"
```

```
  Constellation: Thefoum

                    ✧
           ·    ✦ ─ ─ ─ ─ ★
               ╱           ╱
        ·     ╱    ✧      ╱
             ★ ─ ─ ─ ─ ─ ✦       ·
            ╱ ╲         ╱
     ✧     ╱   ✧      ╱
          ✦ ─ ─ ─ ─ ★
              ·            ✧

  ─── Mythology ───────────────────────────

  The Aethervolk, nomadic starlight navigators of the
  upper winds, tell of a dancer who once challenged
  the darkness to a game of names. When the dancer
  could not win, they became light itself — scattered
  across the sky in the shape we now call Thefoum.

  The elders say: "What you name, you can never lose."

  Best viewed facing south in late autumn, after the
  second meal, when the horizon still remembers the sun.
```

---

## ✧ The Cultures

Every constellation is claimed by one of ten fictional civilizations,
each with their own relationship to the sky:

| Culture | Who They Are |
|---------|-------------|
| **Selenari** | Moon-worshippers of the Silver Coast |
| **Aethervolk** | Nomadic starlight navigators of the upper winds |
| **Thalassians** | Deep-sea dwellers who see the sky as an inverted ocean |
| **Caelidrae** | Mountain folk who carve star-maps into glaciers |
| **Ignari** | Desert astronomers who chart the sky during sandstorms |
| **Arborites** | Forest-dwellers who read constellations in leaf shadows |
| **Pelagians** | Island singers who sing to the stars at equinoxes |
| **Umbrani** | Subterranean dreamers of a sky they have never seen |
| **Stratosi** | Cloud-city dwellers who live close enough to touch the stars |
| **Errantines** | Wandering nomads whose only homeland is the night sky |

The culture, story, characters, and moral are all seeded from your input
text — so the same sentence will always conjure the same myth.

---

## ✦ Determinism

Constellation Composer is **fully deterministic**. Given the same input
string, the output is identical — same star positions, same connections,
same name, same mythology, same viewing instructions. Every random-looking
element is derived from FNV-64a hashes or SHA-256 seeded RNGs.

The sky doesn't change. Only your words do.

---

## ★ Rendering Details

### Terminal

Stars are rendered as Unicode glyphs based on brightness:

| Brightness | Glyph | Color |
|-----------|-------|-------|
| > 0.8 | `✦` | Gold `#FFD700` |
| 0.5 – 0.8 | `★` | Sky blue `#87CEEB` |
| 0.3 – 0.5 | `✧` | Dim blue `#5B7EA0` |
| < 0.3 | `·` | Gray `#3A3A4A` |

Lines are drawn with Bresenham's algorithm using box-drawing characters
(`─ │ ╱ ╲`). A scattering of faint `·` dots fills the background sky.

### SVG

Exported SVGs feature a deep navy sky (`#0a0a2e`), Gaussian-blur glow
effects on stars and lines, cross-shaped spikes on bright stars, and a
vignette that fades to black at the edges. Background stars are sprinkled
for depth. The constellation name is set in Georgia at the bottom.

---

## ✧ Project Structure

```
constellation/
├── main.go          # CLI entry point & flag parsing
├── starmap.go       # Text → deterministic star positions
├── connect.go       # Prim's MST + artistic edge selection
├── render_term.go   # Unicode terminal renderer
├── render_svg.go    # SVG export with glow effects
└── mythology.go     # Procedural myth generation
```

---

## ✦ License

Do whatever the stars tell you to.

---

*Built because every sentence deserves its own sky.*
