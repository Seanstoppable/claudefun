package main

import "math/rand"

// Event represents a random occurrence between turns.
type Event struct {
	Name        string
	Description string
	Effect      Effect
}

// EventGenerator produces random events each turn.
type EventGenerator struct {
	rng *rand.Rand
}

// NewEventGenerator creates an EventGenerator with the given seed.
func NewEventGenerator(seed int64) *EventGenerator {
	return &EventGenerator{rng: rand.New(rand.NewSource(seed))}
}

// conditionalEvent pairs a kingdom-state predicate with a higher-probability event.
type conditionalEvent struct {
	condition func(k *Kingdom) bool
	event     Event
}

func baseEventPool() []Event {
	return []Event{
		// ── Positive ──
		{
			Name:        "Bountiful Harvest",
			Description: "The fields overflow with produce!",
			Effect:      Effect{Food: 20, FactionEffects: map[Faction]int{Farmers: 10}},
		},
		{
			Name:        "Trade Caravan",
			Description: "Merchants from afar bring exotic goods!",
			Effect:      Effect{Treasury: 25, FactionEffects: map[Faction]int{Merchants: 10}},
		},
		{
			Name:        "Cultural Festival",
			Description: "The arts flourish!",
			Effect:      Effect{Culture: 15, Happiness: 10},
		},
		{
			Name:        "Military Parade",
			Description: "The troops look sharp!",
			Effect:      Effect{Military: 10, Reputation: 10},
		},
		{
			Name:        "Baby Boom",
			Description: "Must be something in the water...",
			Effect:      Effect{Population: 15},
		},
		{
			Name:        "Diplomatic Gift",
			Description: "A neighboring kingdom sends gold as tribute!",
			Effect:      Effect{Treasury: 30, Reputation: 10},
		},
		{
			Name:        "Scholar's Discovery",
			Description: "New knowledge is unlocked!",
			Effect:      Effect{Culture: 20, FactionEffects: map[Faction]int{Scholars: 15}},
		},
		{
			Name:        "Jester's Performance",
			Description: "The kingdom laughs!",
			Effect:      Effect{Happiness: 15, FactionEffects: map[Faction]int{Jesters: 10}},
		},

		// ── Negative ──
		{
			Name:        "Plague of Mice",
			Description: "Mice devour the grain stores!",
			Effect:      Effect{Food: -20, FactionEffects: map[Faction]int{Farmers: -10}},
		},
		{
			Name:        "Tax Revolt",
			Description: "Citizens refuse to pay!",
			Effect:      Effect{Treasury: -20, Happiness: -10},
		},
		{
			Name:        "Border Skirmish",
			Description: "Bandits raid the frontier!",
			Effect:      Effect{Military: -10, Population: -5},
		},
		{
			Name:        "Noble Scandal",
			Description: "The duke was found wearing a chicken costume",
			Effect:      Effect{Reputation: -15, FactionEffects: map[Faction]int{Nobles: -15}},
		},
		{
			Name:        "Drought",
			Description: "The rains won't come...",
			Effect:      Effect{Food: -15, Happiness: -10},
		},
		{
			Name:        "Merchant Fraud",
			Description: "A con artist swindled the treasury!",
			Effect:      Effect{Treasury: -25, FactionEffects: map[Faction]int{Merchants: -15}},
		},
		{
			Name:        "Jester's Insult",
			Description: "The jester roasted the wrong diplomat...",
			Effect:      Effect{Reputation: -20, FactionEffects: map[Faction]int{Jesters: -10}},
		},
		{
			Name:        "Library Fire",
			Description: "The great library burns!",
			Effect:      Effect{Culture: -20, FactionEffects: map[Faction]int{Scholars: -15}},
		},

		// ── Weird ──
		{
			Name:        "Mysterious Fog",
			Description: "A thick fog descends. When it lifts, all the statues have moved.",
			Effect:      Effect{Culture: 5, Happiness: -5},
		},
		{
			Name:        "Singing Fish",
			Description: "Fish in the river start singing. Scholars are baffled.",
			Effect:      Effect{Culture: 10, FactionEffects: map[Faction]int{Scholars: 5}},
		},
		{
			Name:        "Rain of Cheese",
			Description: "It's raining cheese! Literally!",
			Effect:      Effect{Food: 15, Happiness: 10, Reputation: -5},
		},
		{
			Name:        "Time Loop Tuesday",
			Description: "Everyone relives Tuesday. Twice.",
			Effect:      Effect{},
		},
		{
			Name:        "Sentient Turnip",
			Description: "A turnip in the royal garden appears to be sentient.",
			Effect:      Effect{Culture: 5, FactionEffects: map[Faction]int{Farmers: 10}},
		},
		{
			Name:        "Goose Uprising",
			Description: "The geese have organized. They have demands.",
			Effect:      Effect{Military: -5, Happiness: 10},
		},
		{
			Name:        "Wandering Wizard",
			Description: "Offers to enchant the castle. It turns slightly to the left.",
			Effect:      Effect{Culture: 10, Reputation: -5},
		},
		{
			Name:        "Ghost Accountant",
			Description: "The treasury is audited by a ghost. It's surprisingly thorough.",
			Effect:      Effect{Treasury: 15},
		},
		{
			Name:        "Competitive Sheep",
			Description: "The sheep have started racing each other. The kingdom watches.",
			Effect:      Effect{Happiness: 15},
		},
	}
}

func conditionalEvents() []conditionalEvent {
	return []conditionalEvent{
		{
			condition: func(k *Kingdom) bool { return k.Food < 30 },
			event: Event{
				Name:        "Bread Riots",
				Description: "The hungry citizens demand food!",
				Effect:      Effect{Happiness: -15, Population: -5},
			},
		},
		{
			condition: func(k *Kingdom) bool { return k.Happiness < 20 },
			event: Event{
				Name:        "Brewing Rebellion",
				Description: "Whispers of revolt grow louder...",
				Effect:      Effect{Military: -10, Reputation: -10},
			},
		},
		{
			condition: func(k *Kingdom) bool { return k.Treasury > 200 },
			event: Event{
				Name:        "Thieves' Guild",
				Description: "Your wealth attracts... professionals.",
				Effect:      Effect{Treasury: -40},
			},
		},
		{
			condition: func(k *Kingdom) bool { return k.Culture > 80 },
			event: Event{
				Name:        "Cultural Renaissance",
				Description: "Artists flock to your kingdom!",
				Effect:      Effect{Population: 10, Culture: 10},
			},
		},
	}
}

// MaybeEvent returns a random event (40% base chance) or nil.
// Conditional events that match the current kingdom state fire at 80%.
func (eg *EventGenerator) MaybeEvent(k *Kingdom) *Event {
	// Check conditional events first (80% trigger chance each).
	for _, ce := range conditionalEvents() {
		if ce.condition(k) && eg.rng.Float64() < 0.80 {
			e := ce.event
			return &e
		}
	}

	// 40% chance of a normal event.
	if eg.rng.Float64() >= 0.40 {
		return nil
	}

	pool := baseEventPool()
	e := pool[eg.rng.Intn(len(pool))]
	return &e
}
