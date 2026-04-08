package octopus

import (
	"math/rand"
	"strings"
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
		"You're sparkling like sun on the waves. Suspicious, but beautiful.",
		"Quick, do a barrel roll! That's what dolphins do when they're happy.",
		"Your joy is making the plankton bloom. In a good way, not the toxic way.",
		"Even the grumpy hagfish is smiling right now because of you.",
		"If happiness were water, you'd be the entire Pacific right now.",
		"The seahorses are dancing in your honor. They do that, you know.",
	},
	Sadness: {
		"Salt water heals everything: tears, sweat, and the sea.",
		"Even pearls start as an irritation. You're becoming a pearl.",
		"I'd pat your back but I'd use all eight arms and it'd be weird.",
		"The whale sings the saddest songs and it's still majestic. So are you.",
		"The deep sea is ninety percent darkness and one hundred percent beautiful. Like you right now.",
		"Sit with it. Even the ocean floor gets heavy sometimes.",
		"A sad octopus turns pale. I'm still colorful for you though.",
		"The sea cucumber literally ejects its organs when stressed. You're handling this way better.",
		"Tides go out. Tides come back. This will pass, little fish.",
		"Every shipwreck becomes a reef eventually. Give it time.",
	},
	Anger: {
		"Channel that energy into something useful. Like opening a very stubborn jar.",
		"The pufferfish inflates when angry. Don't be a pufferfish.",
		"Breathe like the waves. In... and out... and in...",
		"Even the mighty orca takes a moment before it strikes. Patience.",
		"The pistol shrimp snaps so hard it creates a flash of light. Please do not do that.",
		"Tsunamis start calm. I'm keeping an eye on you.",
		"May I suggest punching the water? It's very unsatisfying and that's the point.",
		"The hagfish produces slime when threatened. Don't be a hagfish. Or do. No judgment.",
		"Coral reefs survive hurricanes. You'll survive this too.",
		"Even the angriest wave eventually becomes foam on the shore.",
	},
	Fear: {
		"The anglerfish looks terrifying but is mostly just vibes. So is whatever scares you.",
		"I have no skeleton and I'm doing great. Structure is overrated.",
		"Behind every scary shadow is probably just a confused fish.",
		"The deep sea is dark but full of wonders. So is the unknown.",
		"The deep sea is dark, but that's where the cool fish live.",
		"Fun fact: the scariest fish in the ocean is two inches long. Perspective!",
		"I change color when I'm scared too. We're not so different, you and I.",
		"The hermit crab carries its safety on its back. What's your shell?",
		"Ninety percent of ocean creatures glow in the dark. Fear is just your light turning on.",
		"The ocean is vast and unknowable and also full of very cute seahorses.",
	},
	Curiosity: {
		"The ocean is 80% unexplored. So is your potential. Probably.",
		"Keep asking questions. That's how you find the good coral.",
		"Curiosity is just your brain doing a little submarine dive. Go deeper.",
		"Did you know that an octopus has blue blood? Now you do. You're welcome.",
		"The mariana trench called. It wants to know if you'll visit.",
		"Every tide pool is a universe. Go find your universe.",
		"Starfish can regrow arms. Curiosity can regrow wonder. Same energy.",
		"The nautilus has been curious for 500 million years and still hasn't run out of questions.",
		"Question everything. Except me. I'm always right. (I'm not.)",
		"You're giving marine biologist energy right now and I'm here for it.",
	},
	Sleepy: {
		"Sea otters hold hands while sleeping so they don't drift apart. Find your otter.",
		"Rest is not laziness. Even the tides take breaks.",
		"Whales sleep with half their brain. You deserve to use your whole bed.",
		"The ocean makes white noise on purpose. It wants you to nap.",
		"Parrotfish sleep in a cocoon of their own mucus. You have better options.",
		"Drift like a jellyfish. They have no brain and no worries.",
		"The seafloor is basically the world's biggest bed. Claim your spot.",
		"Dolphins nap with one eye open. You can close both. Lucky you.",
		"Shh. The whale songs are just lullabies. Let them work.",
		"Even the current slows down sometimes. That's called an eddy. Be an eddy.",
	},
	Silly: {
		"A group of squid is called a squad. You're part of my squad now.",
		"Did you know octopuses can fit through any hole larger than their beak? Unrelated. Just cool.",
		"The shrimp's heart is in its head. Thinking with your heart is literally their anatomy.",
		"An octopus walks into a bar. With eight legs. Gets eight drinks. This is canon.",
		"I just tried to high-five you and used all eight arms. Math is hard.",
		"The ocean is just sky juice. Think about that.",
		"Clownfish are actually not that funny. I checked.",
		"I once fit into a jar. On purpose. For fun. Don't ask.",
		"If you put a crab on a treadmill it will walk on it. Scientists did this. Your tax dollars.",
		"Lobsters pee out of their faces. I have nothing to add to this.",
	},
	Love: {
		"Octopuses are solitary creatures. But for you, I make an exception.",
		"My three hearts all beat for you. That's not a metaphor, it's biology.",
		"Love is the current that connects every ocean. You are that current.",
		"The moon pulls the tides and you pull my heartstrings. All three sets.",
		"In the deep sea, anglerfish fuse together when they find their mate. Too much? Maybe.",
		"You're the pearl in my oyster. That sounds weird. I stand by it.",
		"Penguins give each other pebbles as proposals. Here's a pebble. 🪨",
		"Dolphins have best friends. You're my best friend. Don't tell the other dolphins.",
		"The ocean connects every continent. Love connects every heart. I am very deep. Literally.",
		"Seahorses hold tails when they travel together. I'd hold your tail if you had one.",
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
		nextAdviceAt: 2 + r.Intn(2), // 2 or 3
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

// GetAdvice picks a random piece of advice (emotion-specific or general)
// and returns it wrapped in a speech bubble.
func (a *Advisor) GetAdvice(emotion Emotion) string {
	pool := make([]string, 0, len(generalAdvice))
	pool = append(pool, generalAdvice...)
	if specific, ok := emotionAdvice[emotion]; ok {
		pool = append(pool, specific...)
	}
	text := pool[a.rng.Intn(len(pool))]
	return a.FormatBubble(text)
}

// FormatBubble wraps text in a cute ASCII speech bubble, word-wrapping
// to fit roughly 35 characters wide inside the border.
func (a *Advisor) FormatBubble(text string) string {
	const maxWidth = 35

	lines := wordWrap(text, maxWidth)

	// Find the widest line to size the box.
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}

	border := strings.Repeat("─", width+2)
	var b strings.Builder
	b.WriteString("  ╭" + border + "╮\n")
	for _, l := range lines {
		padding := strings.Repeat(" ", width-len(l))
		b.WriteString("  │ " + l + padding + " │\n")
	}
	b.WriteString("  ╰" + border + "╯")
	return b.String()
}

// wordWrap splits text into lines no longer than maxWidth, breaking on spaces.
func wordWrap(text string, maxWidth int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	current := words[0]
	for _, w := range words[1:] {
		if len(current)+1+len(w) > maxWidth {
			lines = append(lines, current)
			current = w
		} else {
			current += " " + w
		}
	}
	lines = append(lines, current)
	return lines
}
