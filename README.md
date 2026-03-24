```
         ___
       /     \
      |       |
      | ♥ ‿ ♥ |
       \_____/
        |||||
     _  |||||  _
    / \ ||||| / \
   |   \|||||/   |
    \   |||||   /
     \  |||||  /
      \ || || /
       \|   |/
        \   /
         \ /
          V
```

# 🐙 Mood Octopus

A whimsical terminal pet that reacts to your emotions. Type anything and watch your octopus friend come alive with animated ASCII art, expressive eyes, and unsolicited ocean-themed life advice.

## Features

- **8 Emotional Arms** — Joy, Sadness, Anger, Fear, Curiosity, Sleepy, Silly, and Love — each with unique arm poses and eye expressions
- **Animated ASCII Art** — smooth frame-based animations with squish transitions between moods
- **Mood Detection** — keyword-based sentiment analysis that supports mixed emotions
- **Unsolicited Wisdom** — 30+ absurd-but-oddly-wise ocean-themed life advice quotes
- **Mood History** — persists your mood timeline across sessions with emoji sparklines
- **Startup Greetings** — your octopus remembers how you felt last time

## Install & Run

```bash
go build -o mood-octopus .
./mood-octopus
```

Or run directly:

```bash
go run .
```

## How It Works

Type a sentence and press Enter. The octopus analyzes your words for emotional keywords and reacts:

| Try typing...                  | Octopus reacts with... |
|-------------------------------|----------------------|
| "I'm so happy today!"         | 😊 Wiggly joy dance   |
| "Everything is terrible"      | 😢 Droopy sad arms    |
| "This makes me furious!"      | 😡 Coiled angry pose  |
| "I'm nervous about tomorrow"  | 😨 Curled-up fear     |
| "How does this even work?"    | 🤔 Curious reaching   |
| "So tired... need sleep"      | 😴 Sleepy with z's    |
| "lol that's bonkers"          | 🤪 Tangled silly arms |
| "I love you"                  | 🥰 Heart-shaped arms  |

## Controls

- **Enter** — Feed your words to the octopus
- **Ctrl+C / Esc** — Say goodbye

## Mood History

Your mood history is saved to `~/.mood-octopus/history.json`. When you return, the octopus greets you based on how you were feeling last time.

## Built With

- [Bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- Pure Go, zero external API dependencies

## Philosophy

> "The ocean doesn't care about your deadlines. Be like the ocean." — Your Octopus

---

*Made with 🐙 and three hearts*
