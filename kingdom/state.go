package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// AllFactions is the canonical list of factions.
var AllFactions = []Faction{Farmers, Merchants, Nobles, Scholars, Jesters}

// Kingdom holds the entire mutable state of a single game.
type Kingdom struct {
	Name        string
	Turn        int
	Treasury    int             // gold coins (can go negative = debt)
	Population  int             // total citizens
	Happiness   int             // 0-100
	Military    int             // 0-100 (defense strength)
	Culture     int             // 0-100 (arts & learning)
	Food        int             // 0-100 (food supply)
	Reputation  int             // 0-100 (how other kingdoms see you)
	FactionMood map[Faction]int // -50 to 50 per faction
	RulerTitle  string          // changes based on performance
	GameOver    bool
	GameOverMsg string
	Victory     bool

	// internal tracking for game-over conditions
	zeroHappinessTurns int // consecutive turns with Happiness == 0
	highStatsTurns     int // consecutive turns with all stats > 80
}

var funKingdomNames = []string{
	"Absurdistan",
	"Chaoswick",
	"Blunderheim",
	"Follyton",
	"Muddle-upon-Thames",
	"Fort Questionable",
	"New Indecisiopolis",
	"Bumbleshire",
	"Daftcastle",
	"Pandemonium Heights",
}

// NewKingdom creates a fresh kingdom with balanced starting stats.
func NewKingdom(name string) *Kingdom {
	if name == "" {
		name = funKingdomNames[rand.IntN(len(funKingdomNames))]
	}

	k := &Kingdom{
		Name:       name,
		Turn:       1,
		Treasury:   100,
		Population: 100,
		Happiness:  55,
		Military:   50,
		Culture:    50,
		Food:       60,
		Reputation: 50,
		FactionMood: map[Faction]int{
			Farmers:   0,
			Merchants: 0,
			Nobles:    0,
			Scholars:  0,
			Jesters:   0,
		},
	}
	k.RulerTitle = k.RulerRating()
	return k
}

// clampStat keeps a value within [min, max].
func clampStat(val, lo, hi int) int {
	if val < lo {
		return lo
	}
	if val > hi {
		return hi
	}
	return val
}

// ApplyEffect mutates the kingdom by adding an Effect's deltas.
func (k *Kingdom) ApplyEffect(e Effect) {
	k.Treasury += e.Treasury
	if k.Treasury > 999999 {
		k.Treasury = 999999
	}
	k.Population += e.Population
	if k.Population < 0 {
		k.Population = 0
	}
	if k.Population > 999999 {
		k.Population = 999999
	}
	k.Happiness = clampStat(k.Happiness+e.Happiness, 0, 100)
	k.Military = clampStat(k.Military+e.Military, 0, 100)
	k.Culture = clampStat(k.Culture+e.Culture, 0, 100)
	k.Food = clampStat(k.Food+e.Food, 0, 100)
	k.Reputation = clampStat(k.Reputation+e.Reputation, 0, 100)

	if e.FactionEffects != nil {
		for faction, delta := range e.FactionEffects {
			k.FactionMood[faction] = clampStat(k.FactionMood[faction]+delta, -50, 50)
		}
	}

	k.RulerTitle = k.RulerRating()
}

// averageStat returns the mean of the five 0-100 stats.
func (k *Kingdom) averageStat() int {
	return (k.Happiness + k.Military + k.Culture + k.Food + k.Reputation) / 5
}

// IsStable returns true when no stats are critically low.
func (k *Kingdom) IsStable() bool {
	return k.Happiness > 10 &&
		k.Military > 10 &&
		k.Culture > 10 &&
		k.Food > 10 &&
		k.Reputation > 10 &&
		k.Treasury > -200 &&
		k.Population > 20
}

// allStatsAbove returns true if every 0-100 stat exceeds the threshold.
func (k *Kingdom) allStatsAbove(threshold int) bool {
	return k.Happiness > threshold &&
		k.Military > threshold &&
		k.Culture > threshold &&
		k.Food > threshold &&
		k.Reputation > threshold
}

var (
	titlesExcellent = []string{"The Magnificent", "The Beloved", "The Wise"}
	titlesGood      = []string{"The Capable", "The Steady", "The Fair"}
	titlesMediocre  = []string{"The Unremarkable", "The Confused", "The Indecisive"}
	titlesPoor      = []string{"The Questionable", "The Bewildered", "The Chaotic"}
	titlesTerrible  = []string{"The Catastrophic", "The Infamous", "The Absurd"}
)

func pickTitle(titles []string) string {
	return titles[rand.IntN(len(titles))]
}

// RulerRating returns a fun ruler title based on overall performance.
func (k *Kingdom) RulerRating() string {
	avg := k.averageStat()
	switch {
	case avg >= 80:
		return pickTitle(titlesExcellent)
	case avg >= 60:
		return pickTitle(titlesGood)
	case avg >= 40:
		return pickTitle(titlesMediocre)
	case avg >= 20:
		return pickTitle(titlesPoor)
	default:
		return pickTitle(titlesTerrible)
	}
}

func moodEmoji(mood int) string {
	switch {
	case mood >= 30:
		return "\U0001f60d"
	case mood >= 10:
		return "\U0001f642"
	case mood <= -30:
		return "\U0001f92c"
	case mood <= -10:
		return "\U0001f620"
	default:
		return "\U0001f610"
	}
}

// TurnSummary returns a brief text summary of the kingdom state.
func (k *Kingdom) TurnSummary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "=== %s \u2014 Turn %d ===\n", k.Name, k.Turn)
	fmt.Fprintf(&b, "Ruler: %s\n", k.RulerTitle)
	fmt.Fprintf(&b, "Treasury: %d gold | Population: %d\n", k.Treasury, k.Population)
	fmt.Fprintf(&b, "Happiness: %d | Military: %d | Culture: %d\n", k.Happiness, k.Military, k.Culture)
	fmt.Fprintf(&b, "Food: %d | Reputation: %d\n", k.Food, k.Reputation)

	factionParts := make([]string, 0, len(AllFactions))
	for _, f := range AllFactions {
		mood := k.FactionMood[f]
		factionParts = append(factionParts, fmt.Sprintf("%s %s(%+d)", string(f), moodEmoji(mood), mood))
	}
	fmt.Fprintf(&b, "Factions: %s\n", strings.Join(factionParts, " | "))

	if !k.IsStable() {
		b.WriteString("\u26a0\ufe0f  The kingdom teeters on the edge!\n")
	}
	return b.String()
}

// CheckGameOver evaluates win/lose conditions and sets GameOver fields.
func (k *Kingdom) CheckGameOver() bool {
	if k.GameOver {
		return true
	}

	// --- Lose conditions ---
	if k.Treasury < -500 {
		k.GameOver = true
		k.GameOverMsg = "The kingdom is bankrupt! Creditors seize the castle."
		return true
	}
	if k.Population < 10 {
		k.GameOver = true
		k.GameOverMsg = "There's no one left to rule. The tumbleweeds judge you."
		return true
	}

	// Track consecutive zero-happiness turns.
	if k.Happiness == 0 {
		k.zeroHappinessTurns++
	} else {
		k.zeroHappinessTurns = 0
	}
	if k.zeroHappinessTurns >= 3 {
		k.GameOver = true
		k.GameOverMsg = "The citizens have revolted! You're exiled to a very small island."
		return true
	}

	// --- Win conditions ---
	if k.allStatsAbove(80) {
		k.highStatsTurns++
	} else {
		k.highStatsTurns = 0
	}
	if k.highStatsTurns >= 5 {
		k.GameOver = true
		k.Victory = true
		k.GameOverMsg = "Your kingdom achieves legendary status!"
		return true
	}

	if k.Turn >= 30 {
		k.GameOver = true
		k.Victory = true
		k.GameOverMsg = "Against all odds, you've survived 30 turns of governance!"
		return true
	}

	return false
}

// AdvanceTurn increments the turn counter and applies natural drift.
func (k *Kingdom) AdvanceTurn() {
	k.Turn++

	// Food decreases by 2-5 (people eat)
	k.Food = clampStat(k.Food-(rand.IntN(4)+2), 0, 100)

	// Treasury decreases by 3-8 (expenses)
	k.Treasury -= rand.IntN(6) + 3

	// Small random fluctuations to happiness (+/-3)
	k.Happiness = clampStat(k.Happiness+(rand.IntN(7)-3), 0, 100)

	// Population grows slightly if food and happiness are high
	if k.Food > 60 && k.Happiness > 50 {
		growth := rand.IntN(5) + 1
		k.Population += growth
	} else if k.Food < 20 {
		// Starvation causes population loss
		k.Population -= rand.IntN(5) + 1
		if k.Population < 0 {
			k.Population = 0
		}
	}

	k.RulerTitle = k.RulerRating()
	k.CheckGameOver()
}
