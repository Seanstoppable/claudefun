# 🎵 Elevator Music Composer

```
  ▐█  ▐█  ▐█  ▐█  ▐█  ▐█  ▐█  ▐█
  ██  ██  ██  ██  ██  ██  ██  ██
  ██  ██  ██  ██  ██  ██  ██  ██
  ██  ██  ██  ██  ██  ██  ██  ██
  ██  ██  ██▌ ██  ██▌ ██  ██▌ ██
  ██▌ ██▌ ██▌ ██▌ ██▌ ██▌ ██▌ ██▌
  ███ ███ ███ ███ ███ ███ ███ ███
  ▀▀▀ ▀▀▀ ▀▀▀ ▀▀▀ ▀▀▀ ▀▀▀ ▀▀▀ ▀▀▀
  ♪ E L E V A T O R  M U S I C ♪
```

**Describe a mood. Get the elevator music album you never knew you needed.**

A procedural elevator music album generator that creates fictional albums complete with artist bios, track listings, ASCII cover art, critic reviews, and monthly listener counts — all deterministically generated from a single mood string.

---

## 🔧 Install & Run

```bash
# Generate a single album
go run . "melancholy"

# Generate a full discography (multiple albums)
go run . -discography "corporate despair"
```

### CLI Flags

| Flag | Description |
|------|-------------|
| `-discography` | Generate 3–5 albums instead of one |
| `-tracks N` | Override the number of tracks (default: 8–13) |
| `-seed N` | Set a specific RNG seed for reproducibility |

---

## 🎶 Genres

Elevator music is a rich and nuanced art form. We support 12 sub-genres:

| Genre | Vibe |
|-------|------|
| **Ambient Lobby** | The sound of a building that's been waiting for you |
| **Corporate Zen** | Synergy in E♭. Team-building, but make it smooth |
| **Waiting Room Jazz** | Old magazines. Older saxophones |
| **Hold Music Deluxe** | Your call is important to us. Estimated wait: forever |
| **Dentist Office Core** | Open wide (softly). Novocaine dreams |
| **Parking Garage Wave** | Concrete echo. Lost tickets. Found vibes |
| **Soft Rock Purgatory** | The hotel breakfast buffet has a soundtrack now |
| **Acoustic Beige** | Off-white serenade. Khaki afternoon. Eggshell emotions |
| **Lo-Fi Elevator** | Crackle & ascend. Dusty lobby tape vibes |
| **Smooth Bureaucracy** | Form B-7 (Smooth Mix). Permit pending groove |
| **Midtempo Malaise** | Not sad, just indoors |
| **Suburban Ambient** | Cul-de-sac calm. HOA-approved vibes only |

---

## 🎵 Sample Track Names

Some of our finest compositions:

- *Waiting for the Elevator to Arrive (But Not Urgently)*
- *Soft Realization That It's Only Tuesday*
- *The Feeling of Being on Hold But It's Fine*
- *Contemplating the Snack Machine (Interlude)*
- *Quiet Enthusiasm for Carpet Samples*
- *Ambient Spreadsheet Energy*
- *The Emotional Range of a Thermostat*
- *Politely Declining a Meeting Invite*
- *Afternoon Slump in D Minor*
- *Saxophone Solo for No Occasion*
- *Vibes (Room Temperature)*
- *Elevator Doors Closing (Take 7)*

Featured artists include **Kenneth (from accounting)**, **The Saxophone of Ambiguity**, and **DJ Smooth Transition**.

---

## ★ Sample Reviews

> "I put this on in the elevator at work and three people complimented the ambiance. Career highlight."
> — *@elevator_enthusiast* ★★★★★

> "Too engaging. I caught myself nodding along during a meeting. Unprofessional."
> — *@cubicle_carl* ★☆☆☆☆

> "Perfectly adequate. Would listen to again if it happened to be playing somewhere I was."
> — *@background_patricia* ★★★½☆

> "This has a saxophone solo that demands attention. That defeats the entire purpose."
> — *@smooth_jazz_steve* ★★☆☆☆

---

## 🌐 Web Version

Also available as a [WebAssembly-powered web app](../web/elevator/) — enter a mood and get your album with ASCII cover art, track listing, and reviews in your browser.

---

## 🏗️ How It Works

Everything is deterministic — the same mood always produces the same album. The mood string is hashed to create a seed, which drives all the random number generators for:

1. **Album generation** — title, artist, genre, label, year
2. **Track listing** — titles from a pool of ~50 options, BPM, key, duration
3. **ASCII cover art** — 8 visual styles selected by genre
4. **Reviews & stats** — critic quotes, play counts, similar artists

No AI, no network calls, no databases. Just vibes (room temperature).

---

*Now playing: Beige Sunset — from the album "Tepid Feelings Vol. 3" by Gentle Smoothington*
*♪ Your call is important to us. Please continue to hold. ♪*
