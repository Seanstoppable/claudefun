```
            ╔═══════════════════════════════════╗
            ║    HERE BE DRAGONS (probably)      ║
            ╚═══════════════════════════════════╝

                         N
                         ▲
                    ╱ · ·│· · ╲
                  ╱· ⛰️ ·│· 🌲 ·╲
                ╱ · · · ·│· · · · ╲
           W ◄──── 🏰 ···+··· 🐉 ────► E
                ╲ · · · ·│· · · · ╱
                  ╲· 🏚️ ·│· 🌊 ·╱
                    ╲ · ·│· · ╱
                         ▼
                         S

       ~~~ The Imaginary Cartographer ~~~
```

# 🗺️ The Imaginary Cartographer

**Chart the uncharted! Name a place, and watch a world unfold.**

*The Imaginary Cartographer* is a Go CLI that conjures entire procedural fantasy worlds from nothing more than a place name. Mountains erupt from syllables. Forests sprout from vowels. Ancient lore writes itself in the margins. Every name is a seed — speak it, and a world answers.

You are the cartographer. The map was always there. You simply… *found* it.

---

## 🧭 Quickstart

```bash
# Clone the repository
git clone https://github.com/your-org/imaginary-cartographer.git
cd imaginary-cartographer

# Generate your first map
go run . "Fumbleshire"
```

A full-color terminal map will appear, complete with terrain, landmarks, and the folklore of a land that never existed — until now.

### Export to SVG

```bash
go run . -svg "Fumbleshire"
```

Produces a beautiful parchment-style SVG with labeled landmarks, a compass rose, and a proper legend. Perfect for printing and pinning to your dungeon master's screen.

---

## ⚙️ CLI Flags

| Flag | Default | Description |
|---|---|---|
| `-svg` | `false` | Export the map as an SVG file instead of terminal output |
| `-w` | `80` | Map width in terminal columns |
| `-h` | `40` | Map height in terminal rows |
| `-svg-width` | `1200` | SVG canvas width in pixels |
| `-svg-height` | `800` | SVG canvas height in pixels |
| `-no-lore` | `false` | Skip folklore generation (just the map, ma'am) |
| `-minimal` | `false` | Minimal output — no colors, no lipgloss, just raw terrain |

```bash
# A grand 120×60 map of the Dreadmoors, SVG only, no lore
go run . -svg -w 120 -h 60 -no-lore "The Dreadmoors"

# A tiny pocket map for quick reference
go run . -w 40 -h 20 -minimal "Gnomevale"
```

---

## 🌍 Biome Table

The cartographer layers Perlin noise, simplex noise, and elevation masks to produce **12 distinct terrain biomes**. Each one renders as colored Unicode in the terminal:

| Symbol | Biome | Description |
|:---:|---|---|
| `🌊` | **Deep Ocean** | The abyss. Leviathans stir below. |
| `≈` | **Shallow Sea** | Turquoise shallows where merchant ships anchor. |
| `░` | **Beach** | Sun-bleached sand, littered with shells and secrets. |
| `♣` | **Forest** | Dense canopy. Paths vanish within three steps. |
| `♠` | **Dark Forest** | Even the trees hold their breath here. |
| `▲` | **Mountains** | Granite spires that scrape the belly of the sky. |
| `△` | **Foothills** | Rolling highlands where shepherds lose their flocks. |
| `∴` | **Desert** | Endless dunes. Mirages sell you water for gold. |
| `≋` | **Swamp** | Mud that remembers every boot that entered. |
| `❋` | **Tundra** | Frozen silence. The wind is the only citizen. |
| `✿` | **Meadow** | Wildflowers and tall grass as far as the eye can see. |
| `⌇` | **Volcanic** | Cracked earth. The ground glows between the fissures. |

---

## 🏰 Landmark Types

Each map generates **8–15 named landmarks**, procedurally placed in appropriate terrain. You may encounter:

- 🏘️ **Towns & Villages** — *"Bridlecross"*, *"Tumblewell"* — where NPCs argue about cheese
- 🏚️ **Ruins** — *"The Sunken Nave"*, *"Pillars of Alsoth"* — crumbling and full of loot (and traps)
- 🍺 **Taverns** — *"The Drowsy Chimera"*, *"Flagon & Fable"* — every adventure starts (and ends) here
- 🗿 **Monoliths** — *"The Wailing Stone"*, *"Obelisk of Thren"* — nobody knows who built them
- 🐉 **Dragon Lairs** — *"Char Hollow"*, *"The Ember Throne"* — marked on the map with a polite skull
- ⛪ **Shrines** — *"Moonwell of Sylara"*, *"The Thorn Altar"* — leave an offering or don't; it's your funeral
- ⚔️ **Battlefields** — *"The Red Meadow"*, *"Ashwalker's Stand"* — the grass never grew back
- 🌀 **Portals** — *"The Flickering Gate"*, *"Rift of Umbrax"* — step through if you're feeling lucky

---

## 📜 Procedural Folklore

Every generated world comes with its own mythology. The lore engine weaves **creation myths**, **local monsters**, **cursed locations**, and **legendary heroes** from the same seed as the map itself.

### A Creation Myth

> *In the age before names, the sky was a river and the ground was a song. The First Cartographer walked barefoot across the silence and, with a burnt stick, drew the coast of Fumbleshire on the back of a sleeping turtle. The turtle woke. The land stayed.*

### A Local Monster

> **The Gravel Wight of Tumblewell**
> A creature assembled from river stones and old grudges. It rises at dusk from the creek bed near Tumblewell, dragging itself toward anyone who has ever broken a promise. Townsfolk leave jars of honey at the bridge to appease it. This works about 40% of the time.

### A Cursed Location

> **The Whispering Cistern** *(beneath the Ruins of Alsoth)*
> Anyone who drinks from the cistern gains the ability to understand the speech of insects — permanently. This sounds useful until you learn what ants think about you. Several adventurers have begged to have the curse reversed. None have succeeded.

---

## 🎲 Determinism

**The same name always produces the same map.** Every pixel, landmark, and scrap of folklore is derived deterministically from the input string. "Fumbleshire" will always be Fumbleshire — the same mountains, the same taverns, the same disgruntled Gravel Wight.

This means you can:
- **Share maps by name** — just tell a friend to run `"Fumbleshire"` and they'll see your world
- **Use names as seeds** — explore variations like `"Fumbleshire North"` or `"Old Fumbleshire"`
- **Reproduce bugs** — if the dragon lair spawned in the ocean, we can see it too

---

## 🛠️ Building

```bash
go build -o cartographer .
./cartographer "The Shattered Expanse"
```

Requires **Go 1.21+**. Dependencies are managed via `go.mod`.

---

## 📦 Dependencies

- [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) — terminal styling that would make a bard weep
- [`ojrac/opensimplex-go`](https://github.com/ojrac/opensimplex-go) — noise generation for terrain layers
- [`ajstarks/svgo`](https://github.com/ajstarks/svgo) — SVG export with the elegance of a royal decree

---

## 🏴‍☠️ A Note from the Cartographer's Guild

> *We do not claim these lands are real. We do not claim they are fake. We claim only that they are* ***mapped***, *and that is enough.*
>
> *Go forth. Name the unnamed. Chart the uncharted. And if you find a tavern called "The Drowsy Chimera," order the stew — trust us.*
>
> — The Imaginary Cartographer's Guild, Third Era, probably

---

<p align="center">
  <code>~ Every name is a world waiting to be found ~</code>
</p>
