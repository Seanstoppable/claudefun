# 🏰 Tiny Kingdom Simulator

```
                 ♛
                ╱ ╲
               ╱ ♦ ╲
              ╱ ♦ ♦ ╲
             ╱───────╲
            │ ♦  ♦  ♦ │
            │  TINY   │
            │ KINGDOM │
     ┌──────┴─────────┴──────┐
     │ ▓▓▓ │ ░░░░░ │ ▓▓▓▓▓ │
     │ ▓▓▓ │ ░░░░░ │ ▓▓▓▓▓ │
     │ ▓▓▓ │ ░░░░░ │ ▓▓▓▓▓ │
     │ ▓▓▓ │ ░░ ░░ │ ▓▓▓▓▓ │
     │ ▓▓▓ │       │ ▓▓▓▓▓ │
     └─────┴───────┴───────┘
    ═══════════════════════════
```

> *"Hear ye, hear ye! A kingdom awaits thy dubious leadership!"*

A **Go + Bubbletea** terminal game where you rule a medieval kingdom through
increasingly absurd policy decisions—while a bard narrates your every blunder
in rhyming verse. Will you be remembered as **The Magnificent**… or **The Catastrophic**?

---

## ⚔️ Install & Run

```bash
go run .
```

That's it. No moats to cross, no drawbridges to lower. Just `go run .` and
claim your throne.

---

## 🎮 How to Play

Each turn, the kingdom presents you with a pressing policy dilemma.
Press **A** or **B** to make your royal decree.

Survive **30 turns** of governance without your kingdom collapsing—or achieve
**legendary status** by getting all stats above 80 for 5 consecutive turns.

Your kingdom starts with a random name plucked from the finest cartographers of
the absurd: *Absurdistan*, *Chaoswick*, *Bumbleshire*, and more.

---

## 📊 The Royal Ledger (Stats)

| Stat | What It Tracks | Starting Value |
|---|---|---|
| 💰 **Treasury** | Gold coins (can go negative — hello, debt!) | 100 |
| 👥 **Population** | How many souls tolerate your rule | 100 |
| 😊 **Happiness** | General mood (0 = pitchforks, 100 = parades) | 55 |
| ⚔️ **Military** | Army strength & readiness | 50 |
| 🎭 **Culture** | Arts, learning, and entertainment | 50 |
| 🍞 **Food** | Granary supply (people eat every turn!) | 60 |
| 🏅 **Reputation** | How neighboring kingdoms see you | 50 |

Stats naturally drift each turn: people eat food, the treasury bleeds expenses,
and happiness wobbles like a jester on a unicycle. Keep everything balanced or
watch it all crumble.

---

## 👥 The Five Factions

Five citizen factions react to your policies with moods ranging from **-50**
(grabbing pitchforks) to **+50** (naming children after you):

| Faction | Cares About |
|---|---|
| 🌾 **Farmers** | Land, crops, farming taxes |
| 🪙 **Merchants** | Trade, markets, economic policy |
| 👑 **Nobles** | Hierarchy, luxury, bowing frequency |
| 📜 **Scholars** | Research, knowledge, funding |
| 🤡 **Jesters** | Entertainment, cultural policies, fun |

Neglect a faction for too long and you'll hear about it. Loudly.

---

## 📜 Sample Royal Decrees

Here are a few of the 46 absurd policies awaiting your wisdom:

> **"Tax all left-handed citizens?"**
> - **A:** *Yes, they're suspiciously dexterous* → 💰+30, 😊−15, 🏅−10
> - **B:** *No, that's absurd* → 😊+5

> **"Make cheese the national currency?"**
> - **A:** *Embrace the fromage economy* → 💰−20, 🎭+15, 😊+10, 🍞−10
> - **B:** *Keep boring old gold* → 💰+5

> **"Court jester demands a seat on the council"**
> - **A:** *Grant the fool a chair* → 🎭+20, 🏅−10, Jesters 🤡+30
> - **B:** *Deny the fool* → Jesters 🤡−20, Nobles 👑+10

> **"Teach the army to dance?"**
> - **A:** *Pirouettes build morale!* → ⚔️+5, 🎭+15, 😊+10
> - **B:** *Soldiers fight, not foxtrot* → ⚔️+15

---

## 🏆 Win Conditions

- **Survival Victory** — Reach Turn 30 without everything collapsing.
  *"Against all odds, you've survived 30 turns of governance!"*
- **Legendary Victory** — Keep Happiness, Military, Culture, Food, and
  Reputation all above 80 for 5 consecutive turns.
  *"Your kingdom achieves legendary status!"*

## 💀 Lose Conditions

- **Bankruptcy** — Treasury drops below −500. Creditors seize the castle.
- **Depopulation** — Population falls below 10. The tumbleweeds judge you.
- **Uprising** — Happiness stays at 0 for 3 consecutive turns. You're exiled
  to a very small island.

---

## 🎵 The Bard

Every decision, every triumph, every spectacular failure is narrated by your
kingdom's resident **bard**—in dramatic rhyming verse.

> *"The ruler pondered, stroked their chin,*
> *And chose: 'Embrace the fromage economy' — let it begin!*
> *The bard records this fateful call,*
> *May it not lead to kingdom's fall!"*

The bard comments on your stats, announces each turn like a chapter in an epic
saga, and delivers your victory speech (or your eulogy) at game's end. You
cannot fire the bard. The bard is eternal.

---

## 👑 Ruler Titles

Your title changes based on your average performance:

| Average Stat | Possible Titles |
|---|---|
| ≥ 80 | The Magnificent, The Beloved, The Wise |
| ≥ 60 | The Capable, The Steady, The Fair |
| ≥ 40 | The Unremarkable, The Confused, The Indecisive |
| ≥ 20 | The Questionable, The Bewildered, The Chaotic |
| < 20 | The Catastrophic, The Infamous, The Absurd |

---

## 🏗️ Built With

- [Go](https://go.dev/)
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — terminal UI framework

---

*Go forth, O Ruler of Dubious Competence! May thy treasury overflow, thy citizens
rejoice, and thy jester never seize the throne. Should thy kingdom fall to ruin,
fear not—for the bard shall ensure thy failures echo through the ages in
magnificent verse.*

*Fare thee well. The crown awaits.* 👑
