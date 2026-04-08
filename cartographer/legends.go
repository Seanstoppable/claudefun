package main

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strings"
)

// Legend represents a single piece of folklore for a place.
type Legend struct {
	Title    string
	Category string // "Creation Myth", "Local Monster", "Cursed Location", "Legendary Hero", "Ancient Mystery"
	Story    string // 2-4 sentences
}

// Lore holds the complete folklore for a named place.
type Lore struct {
	PlaceName    string
	Motto        string
	CreationMyth string
	Legends      []Legend
	Warning      string
}

// LoreGenerator produces deterministic folklore from place names.
type LoreGenerator struct {
	rng *rand.Rand
}

// NewLoreGenerator creates a LoreGenerator seeded from the given string.
func NewLoreGenerator(seed string) *LoreGenerator {
	h := fnv.New64a()
	h.Write([]byte(seed))
	return &LoreGenerator{rng: rand.New(rand.NewSource(int64(h.Sum64())))}
}

// pick returns a random element from a string slice.
func (lg *LoreGenerator) pick(pool []string) string {
	return pool[lg.rng.Intn(len(pool))]
}

// pickInt returns a random int in [lo, hi].
func (lg *LoreGenerator) pickInt(lo, hi int) int {
	return lo + lg.rng.Intn(hi-lo+1)
}

// inventName generates a whimsical proper name.
func (lg *LoreGenerator) inventName() string {
	prefixes := []string{
		"Grom", "Eld", "Thun", "Mor", "Fen", "Bram", "Gor", "Ald",
		"Vex", "Dun", "Cal", "Wynn", "Tor", "Ash", "Lyn", "Bel",
	}
	suffixes := []string{
		"wick", "dan", "thur", "bella", "gast", "worth", "mund", "ora",
		"iel", "ston", "dale", "mere", "bane", "fell", "shaw", "holm",
	}
	return lg.pick(prefixes) + lg.pick(suffixes)
}

// ---------------------------------------------------------------------------
// Content pools
// ---------------------------------------------------------------------------

var mottoTemplates = []string{
	"Where the map ends, {place} begins",
	"Founded by accident, sustained by stubbornness",
	"Population: optimistic",
	"Twinned with nowhere in particular",
	"{place}: We're on no map, and we like it that way",
	"Come for the mystery, stay because you're lost",
	"Established before records, remembered after everything",
	"Where even the compass points shrug",
	"{place}: probably real, definitely here",
	"Abandon certainty, all ye who enter",
	"Here be {place} — you have been mildly warned",
	"Not the destination, but certainly a location",
}

var creationMythTemplates = []string{
	"Legend holds that {place} was sneezed into existence by a sleeping giant named {name}. The hills are said to be the folds of their blanket, and the river their drool. The giant has never been found, which is considered a blessing.",
	"The founding of {place} is attributed to a cartographer who drew a map of a place that didn't exist — and then it did. Whether the map created the land or the land was always waiting to be drawn remains a matter of heated tavern debate.",
	"{place} emerged from the sea one Tuesday afternoon, fully formed and slightly confused. The fish who lived there before are still filing complaints.",
	"According to oral tradition, {place} was the result of two arguing gods who each tried to create their own paradise. The compromise is... {place}.",
	"The elders say {place} was once a word spoken by the wind, and the land grew to match the sound. This explains why {feature} looks the way it does.",
	"They say {place} fell out of a giant's pocket while they were searching for spare change. The giant never came back, presumably finding better currency elsewhere.",
	"Scholars believe {place} was originally a footnote in a divine manuscript that got promoted to a full paragraph. The formatting errors are still visible in the landscape.",
	"{place} was founded when {name} tripped over a rock and declared 'I'm not walking any further.' The rock is now a national monument.",
}

var monsterCreatures = []string{
	"Whelk", "Badger", "Heron", "Salamander", "Moth",
	"Toad", "Catfish", "Owl", "Hedgehog", "Fox",
	"Stoat", "Newt", "Crow", "Snail", "Goose",
}

var monsterAdjectives = []string{
	"Spectral", "Enormous", "Philosophical", "Melodramatic", "Tax-Collecting",
	"Invisible", "Argumentative", "Well-Dressed", "Melancholy", "Polite",
	"Bureaucratic", "Translucent", "Yodeling", "Unusually Punctual", "Three-Headed",
}

var monsterBehaviors = []string{
	"collect debts from the unwary",
	"sing off-key at midnight",
	"rearrange signposts",
	"critique travelers' fashion choices",
	"demand riddles (but accept knock-knock jokes)",
	"steal left shoes exclusively",
	"file noise complaints on behalf of the trees",
	"leave strongly worded letters under doors",
	"photobomb landscape paintings",
	"charge tolls in biscuits",
}

var monsterAdvice = []string{
	"running in a zigzag pattern",
	"offering a compliment and walking briskly away",
	"carrying a spare sandwich at all times",
	"pretending you didn't see it",
	"asking it about its feelings",
	"humming loudly and avoiding eye contact",
	"offering it a cup of tea",
	"backing away slowly while applauding",
	"reciting poetry in a calm but authoritative voice",
	"throwing a decoy hat in the opposite direction",
	"complimenting its posture and excusing yourself",
	"standing very still and thinking about cheese",
	"apologising profusely in a foreign accent",
	"challenging it to a game of chess (carry a travel set)",
	"reading aloud from the nearest Terms and Conditions",
	"waving a white handkerchief while moonwalking",
}

var cursedFeatures = []string{
	"Well", "Bridge", "Crossroads", "Clearing", "Stone Circle",
	"Tower", "Pond", "Tree", "Doorway", "Cellar",
	"Fountain", "Staircase", "Gazebo", "Bench", "Sundial",
}

var cursedAdjectives = []string{
	"Whispering", "Bottomless", "Spinning", "Upside-Down", "Forgetting",
	"Shrinking", "Multiplying", "Backwards", "Hiccupping", "Migrating",
}

var cursedEffects = []string{
	"an irresistible urge to yodel",
	"temporarily forgetting their own name",
	"their socks switching feet",
	"speaking in rhyming couplets for 24 hours",
	"an inexplicable craving for turnips",
	"hearing faint accordion music from no discernible source",
	"a persistent belief that it is always Wednesday",
	"walking slightly to the left for the rest of the day",
	"suddenly knowing the complete history of spoons",
	"an unshakeable conviction that they are being narrated",
	"their shadow arriving three seconds late",
	"involuntarily winking at strangers",
	"perceiving all doors as slightly too small",
	"developing strong opinions about fonts",
	"their applause becoming inexplicably slow",
	"briefly forgetting what chairs are for",
}

var cursedSigns = []string{
	"YOU HAVE BEEN WARNED",
	"Closed for Cursing — Back at 3",
	"Not Our Problem",
	"Enter At Your Own Existential Risk",
	"The Management Accepts No Liability For Temporal Anomalies",
	"If You Can Read This, It's Already Too Late",
	"Please Curse Responsibly",
	"No Refunds On Your Sense Of Direction",
	"Haunted By Appointment Only",
	"Trespassers Will Be Mildly Inconvenienced",
	"Do Not Tap The Glass (There Is No Glass)",
	"This Area Provided As-Is, Without Warranty",
	"Caution: Cursing In Progress — Hard Hats Optional",
	"Beyond This Point, All Complaints Are Alphabetised And Ignored",
	"Mind The Existential Gap",
	"Abandon All Hats, Ye Who Enter Here",
}

var cursedDays = []string{
	"Tuesday", "Thursday", "every other Wednesday", "full moon",
	"days ending in 'y'", "bank holiday", "the third Sunday of months with an R in them",
	"the solstice", "new moon", "leap day",
	"any day the wind blows from the north", "Mondays (obviously)",
	"the second breakfast hour", "days when the fog is particularly judgmental",
	"whenever the church bell rings thirteen times", "alternate Saturdays",
	"the anniversary of something nobody can remember",
}

var heroNames = []string{
	"Grendis", "Morrigan", "Thistlewick", "Bronwyn", "Oakshield",
	"Fernsby", "Caldwell", "Brambleheart", "Ashwick", "Thornley",
	"Cobblesworth", "Kettleburn", "Millhaven", "Crumpet", "Bixby",
}

var heroTitles = []string{
	"Adequate", "Mostly Brave", "Slightly Lost", "Reluctant",
	"Accidentally Famous", "Moderately Heroic", "Surprisingly Effective",
	"Well-Intentioned", "Arguably Competent", "Bafflingly Lucky",
}

var heroDeeds = []string{
	"single-handedly defeating a dragon (which turned out to be a large iguana)",
	"discovering the cure for hiccups (since lost)",
	"negotiating peace between two warring bakeries",
	"mapping the unmappable swamp (and getting only slightly lost)",
	"inventing the umbrella 200 years after everyone else",
	"convincing a mountain to move three feet to the right",
	"reading an entire library and summarizing it as 'meh'",
	"accidentally founding the postal service while trying to return a borrowed book",
	"winning a staring contest with a statue",
}

var heroFestivals = []string{
	"Thistlewick Day",
	"The Feast of Mild Heroism",
	"The Annual Parade of Questionable Bravery",
	"The Festival of Found Things",
	"The Grand Celebration of Adequate Achievement",
	"The Procession of Polite Applause",
	"The Week of Lukewarm Remembrance",
	"The Symposium of Accidental Victories",
	"The Evening of Reluctant Toasts",
	"The March of the Vaguely Competent",
	"The Biannual Shrug of Gratitude",
	"The Afternoon of Gentle Exaggeration",
	"The Gala of Participation Trophies",
	"The Quiet Acknowledging Nod Festival",
	"The Day of Surprisingly Low Expectations",
	"The Pageant of Almost Getting It Right",
}

var mysteryAdjectives = []string{
	"Humming", "Singing", "Glowing", "Wandering", "Counting",
	"Whispering", "Rotating", "Debating", "Sulking", "Dancing",
}

var mysteryFeatures = []string{
	"Path", "River", "Lake", "Statue", "Garden",
	"Door", "Clock", "Library", "Map", "Bell",
}

var mysteryTimePeriods = []string{
	"full moon", "leap year", "third Thursday of the month",
	"whenever it feels like it", "every seventeen minutes",
	"the autumn equinox", "Tuesdays (but only in winter)",
	"precisely noon on days that don't exist",
	"the vernal equinox (give or take)",
	"every time someone nearby sneezes",
	"at the exact moment tea is served",
	"during particularly dramatic sunsets",
	"the stroke of midnight on a Wednesday",
	"once per fortnight, weather permitting",
	"when the barometric pressure is 'vibes'",
}

var warnings = []string{
	"Do not feed the wildlife. The wildlife does not need encouragement.",
	"Maps of {place} are approximate. Reality here is also approximate.",
	"Time moves differently in {place}. Pack accordingly.",
	"Visitors are advised to bring comfortable shoes and low expectations.",
	"If the fog asks you a question, you are under no obligation to answer.",
	"The locals are friendly. The geography is less so.",
	"Do not accept directions from anyone who glows faintly.",
	"All exits are also entrances. Choose wisely.",
	"The weather in {place} is best described as 'optional'.",
	"If the ground hums beneath your feet, it is not your imagination. It is, however, not your concern.",
}

var landscapeFeatures = []string{
	"the crooked hill", "the eastern bluff", "the old millpond",
	"the twisted oak grove", "the northern ridge", "the sunken meadow",
	"the chalk cliffs", "the mossy ravine", "the wobbling spire",
	"the whispering hollow", "the lopsided tor", "the boggy ditch",
	"the peculiar outcrop", "the leaning pines", "the forgotten quarry",
	"the perpetually damp gully",
}

var landmarks = []string{
	"the old bridge", "the town square", "the abandoned mill",
	"the northern gate", "the market fountain", "the bell tower",
	"the docks at dawn", "the crossroads shrine", "the graveyard fence",
	"the council steps", "the baker's lane", "the tilted signpost",
	"the parish noticeboard", "the second-best well", "the lamplighter's corner",
	"the overgrown bandstand",
}

// ---------------------------------------------------------------------------
// Legend generators (one per category)
// ---------------------------------------------------------------------------

func (lg *LoreGenerator) generateMonster(placeName string) Legend {
	adj := lg.pick(monsterAdjectives)
	creature := lg.pick(monsterCreatures)
	title := fmt.Sprintf("The %s %s of %s", adj, creature, placeName)
	size := lg.pickInt(7, 30)
	landmark := lg.pick(landmarks)
	behavior := lg.pick(monsterBehaviors)
	advice := lg.pick(monsterAdvice)

	story := fmt.Sprintf(
		"Reportedly %d feet tall, the %s has been spotted near %s. It is known to %s. Experts recommend %s.",
		size, creature, landmark, behavior, advice,
	)
	return Legend{Title: title, Category: "Local Monster", Story: story}
}

func (lg *LoreGenerator) generateCursedLocation(placeName string) Legend {
	adj := lg.pick(cursedAdjectives)
	feature := lg.pick(cursedFeatures)
	title := fmt.Sprintf("The %s %s", adj, feature)
	effect := lg.pick(cursedEffects)
	day := lg.pick(cursedDays)
	sign := lg.pick(cursedSigns)

	story := fmt.Sprintf(
		"Those who visit the %s %s report %s. The locals avoid it on %ss. A sign posted nearby reads '%s'.",
		adj, feature, effect, day, sign,
	)
	return Legend{Title: title, Category: "Cursed Location", Story: story}
}

func (lg *LoreGenerator) generateHero(_ string) Legend {
	name := lg.pick(heroNames)
	heroTitle := lg.pick(heroTitles)
	title := fmt.Sprintf("%s the %s", name, heroTitle)
	deed := lg.pick(heroDeeds)
	festival := lg.pick(heroFestivals)

	story := fmt.Sprintf(
		"Known for %s. Their statue stands in the town square, though it bears no resemblance to them whatsoever. The annual %s is held in their honor.",
		deed, festival,
	)
	return Legend{Title: title, Category: "Legendary Hero", Story: story}
}

func (lg *LoreGenerator) generateMystery(placeName string) Legend {
	// Two mystery templates
	switch lg.rng.Intn(2) {
	case 0:
		adj := lg.pick(mysteryAdjectives)
		n := lg.pickInt(5, 13)
		title := fmt.Sprintf("The %s Stones of %s", adj, placeName)
		story := fmt.Sprintf(
			"A circle of %d stones that hum on the solstice. Nobody knows who built them, but a plaque claims credit on behalf of the local gardening club.",
			n,
		)
		return Legend{Title: title, Category: "Ancient Mystery", Story: story}
	default:
		feature := lg.pick(mysteryFeatures)
		timePeriod := lg.pick(mysteryTimePeriods)
		distance := lg.pickInt(2, 40)
		papers := lg.pickInt(12, 347)
		title := fmt.Sprintf("The Disappearing %s", feature)
		story := fmt.Sprintf(
			"Every %s, the %s vanishes completely, only to reappear %d feet to the left. Scholars have published %d papers on this phenomenon. None agree.",
			timePeriod, feature, distance, papers,
		)
		return Legend{Title: title, Category: "Ancient Mystery", Story: story}
	}
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// Generate produces a complete Lore for the given place name.
func (lg *LoreGenerator) Generate(placeName string) *Lore {
	// Re-seed from the place name for per-place determinism.
	h := fnv.New64a()
	h.Write([]byte(placeName))
	lg.rng = rand.New(rand.NewSource(int64(h.Sum64())))

	lore := &Lore{PlaceName: placeName}

	// Motto
	motto := lg.pick(mottoTemplates)
	lore.Motto = strings.ReplaceAll(motto, "{place}", placeName)

	// Creation myth
	myth := lg.pick(creationMythTemplates)
	myth = strings.ReplaceAll(myth, "{place}", placeName)
	myth = strings.ReplaceAll(myth, "{name}", lg.inventName())
	myth = strings.ReplaceAll(myth, "{feature}", lg.pick(landscapeFeatures))
	lore.CreationMyth = myth

	// Legends: 3-5 per place, mixing categories.
	legendCount := lg.pickInt(3, 5)
	categories := []func(string) Legend{
		lg.generateMonster,
		lg.generateCursedLocation,
		lg.generateHero,
		lg.generateMystery,
	}

	// Always include at least one monster and one hero; fill the rest randomly.
	lore.Legends = append(lore.Legends, lg.generateMonster(placeName))
	lore.Legends = append(lore.Legends, lg.generateHero(placeName))
	for i := 2; i < legendCount; i++ {
		gen := categories[lg.rng.Intn(len(categories))]
		lore.Legends = append(lore.Legends, gen(placeName))
	}

	// Warning
	warning := lg.pick(warnings)
	lore.Warning = strings.ReplaceAll(warning, "{place}", placeName)

	return lore
}
