package main

import (
	"math/rand"
	"time"
)

var generalAdvice = []string{
	"Have you tried being a sea cucumber instead? Less stress.",
	"The ocean doesn't care about your deadlines. Be like the ocean.",
	"Ink first, ask questions later.",
	"Remember: even the kraken has bad days.",
	"You're doing swimmingly. Trust me, I have eight arms.",
	"Fun fact: octopuses have three hearts. I'm giving you all three right now.",
	"The deepest part of the ocean is still part of the ocean. You belong.",
	"When in doubt, squirt ink and swim away.",
	"You know what's great about having eight arms? Eight hugs at once.",
	"Be the tide you wish to see in the world.",
	"Coral doesn't grow overnight. Neither do you. Be patient.",
	"Every wave was once just a ripple that believed in itself.",
}

var emotionAdvice = map[Emotion][]string{
	Joy: {
		"Bottle this feeling! Actually, don't. Messages in bottles are unreliable.",
		"Your happiness makes my tentacles tingle!",
		"Joy is just the ocean sparkling. You're the ocean.",
		"You're glowing brighter than a bioluminescent jellyfish right now.",
	},
	Sadness: {
		"Salt water heals everything: tears, sweat, and the sea.",
		"Even pearls start as an irritation. You're becoming a pearl.",
		"I'd pat your back but I'd use all eight arms and it'd be weird.",
		"The whale sings the saddest songs and it's still majestic. So are you.",
	},
	Anger: {
		"Channel that energy into something useful. Like opening a very stubborn jar.",
		"The pufferfish inflates when angry. Don't be a pufferfish.",
		"Breathe like the waves. In... and out... and in...",
		"Even the mighty orca takes a moment before it strikes. Patience.",
	},
	Fear: {
		"The anglerfish looks terrifying but is mostly just vibes. So is whatever scares you.",
		"I have no skeleton and I'm doing great. Structure is overrated.",
		"Behind every scary shadow is probably just a confused fish.",
		"The deep sea is dark but full of wonders. So is the unknown.",
	},
	Curiosity: {
		"The ocean is 80% unexplored. So is your potential. Probably.",
		"Keep asking questions. That's how you find the good coral.",
		"Curiosity is just your brain doing a little submarine dive. Go deeper.",
	},
	Sleepy: {
		"Sea otters hold hands while sleeping so they don't drift apart. Find your otter.",
		"Rest is not laziness. Even the tides take breaks.",
		"Whales sleep with half their brain. You deserve to use your whole bed.",
	},
	Silly: {
		"A group of squid is called a squad. You're part of my squad now.",
		"Did you know octopuses can fit through any hole larger than their beak? Unrelated. Just cool.",
		"The shrimp's heart is in its head. Thinking with your heart is literally their anatomy.",
	},
	Love: {
		"Octopuses are solitary creatures. But for you, I make an exception.",
		"My three hearts all beat for you. That's not a metaphor, it's biology.",
		"Love is the current that connects every ocean. You are that current.",
	},
}

// Advisor tracks input cadence and dispenses advice at intervals.
type Advisor struct {
	inputCount   int
	nextAdviceAt int
	rng          *rand.Rand
}

// NewAdvisor returns an Advisor ready to dispense wisdom.
func NewAdvisor() *Advisor {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Advisor{
		inputCount:   0,
		nextAdviceAt: 2 + r.Intn(2),
		rng:          r,
	}
}

// ShouldGiveAdvice increments the input counter and returns true
// when it's time to dispense advice (every 2–3 inputs).
func (a *Advisor) ShouldGiveAdvice() bool {
	a.inputCount++
	if a.inputCount >= a.nextAdviceAt {
		a.inputCount = 0
		a.nextAdviceAt = 2 + a.rng.Intn(2)
		return true
	}
	return false
}

// GetAdvice picks a random piece of advice for the given emotion.
func (a *Advisor) GetAdvice(emotion Emotion) string {
	pool := make([]string, 0, len(generalAdvice))
	pool = append(pool, generalAdvice...)
	if specific, ok := emotionAdvice[emotion]; ok {
		pool = append(pool, specific...)
	}
	return pool[a.rng.Intn(len(pool))]
}
