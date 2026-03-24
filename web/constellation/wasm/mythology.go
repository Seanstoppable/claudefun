package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
)

// Mythology holds the procedurally generated backstory for a constellation.
type Mythology struct {
	Name        string // constellation name (passed in)
	Culture     string // the fictional culture that named it
	Story       string // the myth itself (2-3 paragraphs)
	Moral       string // the lesson/wisdom
	BestViewing string // fictional viewing conditions
}

// MythologyGenerator builds deterministic mythology from constellation metadata.
type MythologyGenerator struct {
	rng *rand.Rand
}

// NewMythologyGenerator creates a generator seeded from the input text for reproducibility.
func NewMythologyGenerator(seed string) *MythologyGenerator {
	h := sha256.Sum256([]byte(seed))
	s := int64(binary.BigEndian.Uint64(h[:8]))
	return &MythologyGenerator{rng: rand.New(rand.NewSource(s))}
}

// Generate produces a complete Mythology for the given constellation.
func (m *MythologyGenerator) Generate(constellationName string, starCount int) Mythology {
	culture := m.pickCulture()
	cultureName := extractCultureName(culture)

	story := m.buildStory(constellationName, starCount, culture, cultureName)
	moral := m.pickMoral()
	viewing := m.pickViewing(cultureName)

	return Mythology{
		Name:        constellationName,
		Culture:     culture,
		Story:       story,
		Moral:       moral,
		BestViewing: viewing,
	}
}

// ── pools ──────────────────────────────────────────────────────────────

var cultures = []string{
	"the ancient Selenari, moon-worshippers of the Silver Coast",
	"the nomadic Aethervolk, who navigated by starlight alone",
	"the deep-sea Thalassians, who believed the sky was an inverted ocean",
	"the mountain-dwelling Caelidrae, who carved star maps into glaciers",
	"the desert Ignari, whose astronomers worked only during sandstorms",
	"the forest-dwelling Arborites, who read constellations in leaf shadows",
	"the island Pelagians, who sang to the stars every equinox",
	"the subterranean Umbrani, who had never seen the sky but dreamed of it",
	"the cloud-city Stratosi, who lived so high they felt they could touch the stars",
	"the wandering Errantines, who claimed no homeland but the night sky itself",
}

var characters = []string{
	"a shepherd who could hear starlight",
	"a weaver who stitched the sky each night",
	"a thief who stole fire from the moon",
	"a dancer whose steps left trails of light",
	"a philosopher who argued with the void",
	"twin rivers that flowed upward into the heavens",
	"a whale made entirely of song",
}

var conflicts = []string{
	"grew too bright and threatened to outshine the sun",
	"fell in love with a mortal shadow",
	"was challenged by the darkness between stars",
	"forgot their purpose and began to wander",
	"discovered that the sky was not infinite after all",
	"tried to count every grain of sand to earn immortality",
}

var morals = []string{
	"Even the smallest light can guide a traveler home.",
	"The distance between stars is nothing compared to the distance between hearts.",
	"To name something is to give it power. Choose your words carefully.",
	"The sky remembers what the earth forgets.",
	"Not all who wander among the stars are lost — some are simply choosing their orbit.",
	"Brilliance is not the absence of darkness, but the courage to shine within it.",
	"Every constellation is a story waiting to be read by the right eyes.",
	"The stars do not compete with one another. They simply shine.",
}

var poeticConditions = []string{
	"the moon still whispered secrets to the tide",
	"the horizon was stitched with threads of gold",
	"silence itself had a sound and that sound was beautiful",
	"the sky wore its stars like an empress wears jewels",
	"every river ran with liquid moonlight",
}

var abstractThings = []string{
	"clocks", "written language", "sorrow", "borders",
	"the naming of seasons", "doubt", "gravity",
}

var events = []string{
	"the ocean tried to swallow the moon",
	"a child asked the sky its true name and received an answer",
	"every bird in the world fell silent for one perfect hour",
	"a single tear fell from the eye of the oldest mountain",
	"the wind itself forgot which direction to blow",
}

var adjectives = []string{
	"restless", "luminous", "proud", "melancholy", "fearless",
	"generous", "quarrelsome", "gentle", "ancient", "mischievous",
}

var verbs = []string{
	"name the seasons", "weep", "build walls", "tell lies",
	"count the days", "forget their dreams",
}

var origins = []string{
	"the shattered crown of a forgotten god",
	"the last breath of a dying comet",
	"a song that was too beautiful to remain mere sound",
	"a promise made between the earth and the void",
	"the sparks struck when two celestial swords clashed",
}

var wisdoms = []string{
	"beauty persists even when no one is watching",
	"nothing truly bright can be owned",
	"every ending is just a beginning viewed from the wrong direction",
	"even the cosmos needs witnesses",
	"what is lost in one world may be found in another",
}

var qualities = []string{
	"the stubbornness of hope",
	"the absurdity of devotion",
	"the quiet power of things that endure",
	"the strange mercy of the universe",
	"the elegance of unfinished things",
}

var latitudes = []string{
	"37°N to 42°N", "the southern tropics", "anywhere north of sadness",
	"the exact latitude of your happiest memory", "60°S, give or take a sigh",
}

var months = []string{
	"the month the ancients called 'Second Silence'",
	"late autumn", "midsummer", "the week before the equinox",
	"any month that contains a Tuesday",
}

var foods = []string{
	"pomegranates", "salted bread", "something you baked yourself",
	"honeycomb", "cold soup", "anything fermented under moonlight",
}

// ── helpers ────────────────────────────────────────────────────────────

func (m *MythologyGenerator) pick(pool []string) string {
	return pool[m.rng.Intn(len(pool))]
}

func (m *MythologyGenerator) pickCulture() string  { return m.pick(cultures) }
func (m *MythologyGenerator) pickMoral() string     { return m.pick(morals) }

func (m *MythologyGenerator) pickViewing(cultureName string) string {
	switch m.rng.Intn(6) {
	case 0:
		return "Best observed on clear nights when you've forgotten something important."
	case 1:
		return fmt.Sprintf("Visible from %s during %s, preferably while eating %s.",
			m.pick(latitudes), m.pick(months), m.pick(foods))
	case 2:
		return "Most vivid at 3 AM, when the boundary between dream and waking thins."
	case 3:
		return "Can only be seen by those who look up at exactly the right moment — which is always now."
	case 4:
		return "Traditionally observed while standing in shallow water under a new moon."
	default:
		return fmt.Sprintf("The %s believed this constellation was clearest when viewed through tears of joy.", cultureName)
	}
}

// extractCultureName pulls the short culture name (e.g. "Selenari") from the full description.
func extractCultureName(culture string) string {
	// Format: "the [adj] NAME, who ..." — grab the capitalised word before the comma.
	parts := strings.SplitN(culture, ",", 2)
	words := strings.Fields(parts[0])
	if len(words) > 0 {
		return words[len(words)-1]
	}
	return "ancients"
}

// ── story construction ─────────────────────────────────────────────────

func (m *MythologyGenerator) buildOpening(name string, starCount int, cultureName string) string {
	switch m.rng.Intn(5) {
	case 0:
		return fmt.Sprintf("In the age before %s, when %s, there lived %s.",
			m.pick(abstractThings), m.pick(poeticConditions), m.pick(characters))
	case 1:
		return fmt.Sprintf("It is said that %s first appeared on the night that %s.",
			name, m.pick(events))
	case 2:
		return fmt.Sprintf("The elders speak of a time when %d stars gathered in council to decide the fate of %s.",
			starCount, m.pick(abstractThings))
	case 3:
		return fmt.Sprintf("Long before mortals learned to %s, the stars themselves were %s and %s.",
			m.pick(verbs), m.pick(adjectives), m.pick(adjectives))
	default:
		return fmt.Sprintf("This constellation was born from %s — or so the %s believed.",
			m.pick(origins), cultureName)
	}
}

func (m *MythologyGenerator) buildConflict(name string, starCount int) string {
	character := m.pick(characters)
	conflict := m.pick(conflicts)

	templates := []string{
		fmt.Sprintf("Among them was %s, who %s. For %d nights the struggle played out across the heavens, each star a witness, each shadow a co-conspirator.",
			character, conflict, starCount+m.rng.Intn(40)+10),
		fmt.Sprintf("But %s %s, and the celestial order trembled. The other constellations turned their faces away, all except %s, which burned only brighter.",
			character, conflict, name),
		fmt.Sprintf("In those days, %s %s. The %d stars of %s watched in silence, for stars are patient beyond all mortal understanding, and they knew that even this too would pass — though not without cost.",
			character, conflict, starCount, name),
	}
	return templates[m.rng.Intn(len(templates))]
}

func (m *MythologyGenerator) buildResolution(name string, starCount int) string {
	switch m.rng.Intn(5) {
	case 0:
		return fmt.Sprintf("And so they were placed among the stars as a reminder that %s. To this day, %s hangs in the firmament, unchanged and unchanging, a silent sermon written in light.",
			m.pick(wisdoms), name)
	case 1:
		return fmt.Sprintf("In the end, they shattered into %d fragments, each becoming a star in %s. Some say if you listen closely, you can still hear the pieces singing to one another across the dark.",
			starCount, name)
	case 2:
		return fmt.Sprintf("The gods, moved by this devotion, wove them into the firmament. %s was their answer — %d points of light arranged just so, a geometry that means 'you are forgiven' in the language of the cosmos.",
			name, starCount)
	case 3:
		return fmt.Sprintf("And so they remain, a testament to %s. Scholars have debated the precise arrangement of %s's %d stars for centuries, but the stars themselves have offered no comment.",
			m.pick(qualities), name, starCount)
	default:
		return fmt.Sprintf("They transformed into light itself, forever dancing across the sky. If %s seems to shimmer on warm nights, it is because the dance has never truly ended — it merely waits for the right eyes to notice.",
			name)
	}
}

func (m *MythologyGenerator) buildStory(name string, starCount int, culture, cultureName string) string {
	opening := m.buildOpening(name, starCount, cultureName)
	conflict := m.buildConflict(name, starCount)
	resolution := m.buildResolution(name, starCount)

	// Weave in a brief cultural attribution between conflict and resolution.
	bridge := fmt.Sprintf(
		"\n\nThe %s recorded this tale in their most sacred texts — though 'sacred' is perhaps too solemn a word for a people who understood that the universe, above all else, has a sense of humor.",
		cultureName,
	)

	return opening + " " + conflict + bridge + " " + resolution
}
