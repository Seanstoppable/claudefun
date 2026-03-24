package main

import (
	"math/rand"
)

// Faction represents a group within the kingdom whose approval matters.
type Faction string

const (
	Farmers   Faction = "Farmers"
	Merchants Faction = "Merchants"
	Nobles    Faction = "Nobles"
	Scholars  Faction = "Scholars"
	Jesters   Faction = "Jesters"
)

// Effect describes the consequences of a policy decision on the kingdom.
type Effect struct {
	Treasury       int
	Population     int
	Happiness      int
	Military       int
	Culture        int
	Food           int
	Reputation     int
	FactionEffects map[Faction]int
}

// Policy represents an absurd dilemma the ruler must resolve.
type Policy struct {
	Question string // the dilemma presented to the ruler
	OptionA  string // first choice
	OptionB  string // second choice
	EffectA  Effect // consequences of choosing A
	EffectB  Effect // consequences of choosing B
	FlavorA  string // funny description of what happens
	FlavorB  string // funny description of what happens
}

// PolicyGenerator hands out random, non-repeating policies from the pool.
type PolicyGenerator struct {
	rng      *rand.Rand
	used     map[int]bool
	policies []Policy
}

// NewPolicyGenerator creates a generator seeded for reproducible absurdity.
func NewPolicyGenerator(seed int64) *PolicyGenerator {
	return &PolicyGenerator{
		rng:      rand.New(rand.NewSource(seed)),
		used:     make(map[int]bool),
		policies: allPolicies(),
	}
}

// NextPolicy returns a random unused policy. Panics if none remain.
func (pg *PolicyGenerator) NextPolicy() Policy {
	if !pg.HasMore() {
		panic("no more policies — the kingdom has exhausted all absurdity")
	}

	for {
		idx := pg.rng.Intn(len(pg.policies))
		if !pg.used[idx] {
			pg.used[idx] = true
			return pg.policies[idx]
		}
	}
}

// HasMore reports whether unused policies remain.
func (pg *PolicyGenerator) HasMore() bool {
	return len(pg.used) < len(pg.policies)
}

func fe(effects map[Faction]int) map[Faction]int { return effects }

func allPolicies() []Policy {
	return []Policy{
		// 1
		{
			Question: "Tax all left-handed citizens?",
			OptionA:  "Yes, they're suspiciously dexterous",
			OptionB:  "No, that's absurd",
			EffectA:  Effect{Treasury: 30, Happiness: -15, Reputation: -10},
			EffectB:  Effect{Happiness: 5},
			FlavorA:  "The left-handed underground forms almost immediately. They're surprisingly well-organized.",
			FlavorB:  "Left-handed citizens breathe a sigh of relief and continue being mildly inconvenienced by scissors.",
		},
		// 2
		{
			Question: "Make cheese the national currency?",
			OptionA:  "Gouda idea!",
			OptionB:  "Keep using boring gold",
			EffectA:  Effect{Treasury: -20, Culture: 15, Happiness: 10, Food: -10},
			EffectB:  Effect{Treasury: 5},
			FlavorA:  "Banks now require climate control. The exchange rate is measured in wedges.",
			FlavorB:  "Gold remains king. The dairy farmers shrug and go back to making lunch.",
		},
		// 3
		{
			Question: "The court jester demands a seat on the council",
			OptionA:  "Grant it. What's the worst that could happen?",
			OptionB:  "Deny it. We have standards.",
			EffectA:  Effect{Culture: 20, Reputation: -10, FactionEffects: fe(map[Faction]int{Jesters: 30})},
			EffectB:  Effect{FactionEffects: fe(map[Faction]int{Jesters: -20, Nobles: 10})},
			FlavorA:  "Council meetings are now 40% more entertaining and 60% less productive.",
			FlavorB:  "The jester juggles angrily outside the council chamber for three straight hours.",
		},
		// 4
		{
			Question: "Scholars want to build a university. It's expensive.",
			OptionA:  "Fund it!",
			OptionB:  "We need soldiers, not scholars",
			EffectA:  Effect{Treasury: -40, Culture: 25, Population: 10, FactionEffects: fe(map[Faction]int{Scholars: 30})},
			EffectB:  Effect{Military: 15, FactionEffects: fe(map[Faction]int{Scholars: -20})},
			FlavorA:  "The University of Applied Nonsense opens its doors. First course: Advanced Hat Theory.",
			FlavorB:  "The scholars mutter darkly about 'anti-intellectualism' while the army gets shinier helmets.",
		},
		// 5
		{
			Question: "A dragon has been spotted near the kingdom. What do we do?",
			OptionA:  "Send the army!",
			OptionB:  "Offer it a diplomatic position",
			EffectA:  Effect{Military: -20, Reputation: 25},
			EffectB:  Effect{Culture: 10, Reputation: 15, Happiness: 5},
			FlavorA:  "The army charges valiantly. The dragon seems mostly confused.",
			FlavorB:  "Ambassador Flame-Breath attends their first trade meeting. Nobody argues with their proposals.",
		},
		// 6
		{
			Question: "The farmers want to unionize",
			OptionA:  "Allow it. Fair wages for all!",
			OptionB:  "Unions? In MY kingdom?",
			EffectA:  Effect{Food: 15, Treasury: -15, FactionEffects: fe(map[Faction]int{Farmers: 25, Nobles: -15})},
			EffectB:  Effect{Food: -5, FactionEffects: fe(map[Faction]int{Farmers: -30, Nobles: 10})},
			FlavorA:  "The Turnip Workers Local 47 holds its first meeting. They demand dental coverage.",
			FlavorB:  "The farmers grumble. Crop yields drop as 'accidents' befall the noble's vegetable gardens.",
		},
		// 7
		{
			Question: "A traveling merchant offers a 'definitely not cursed' crown",
			OptionA:  "Buy it! (50 gold)",
			OptionB:  "Decline suspiciously",
			EffectA:  Effect{Treasury: -50, Culture: 10, Happiness: -10},
			EffectB:  Effect{Reputation: 5},
			FlavorA:  "The crown whispers unsettling tax advice at 3am. It was definitely cursed.",
			FlavorB:  "The merchant disappears in a puff of suspiciously glittery smoke. Wise choice.",
		},
		// 8
		{
			Question: "The people demand a national holiday celebrating turnips",
			OptionA:  "Declare Turnip Day!",
			OptionB:  "Turnips aren't that special",
			EffectA:  Effect{Happiness: 20, Food: -10, Treasury: -10, FactionEffects: fe(map[Faction]int{Farmers: 15})},
			EffectB:  Effect{Happiness: -5, FactionEffects: fe(map[Faction]int{Farmers: -15})},
			FlavorA:  "The parade features a giant turnip float. Children weep with joy. Adults weep for other reasons.",
			FlavorB:  "The farmers begin a quiet but effective turnip embargo. Soups get much worse.",
		},
		// 9
		{
			Question: "Your advisor suggests building a moat filled with pudding",
			OptionA:  "Brilliant defense strategy!",
			OptionB:  "Use water like a normal kingdom",
			EffectA:  Effect{Military: 5, Food: -20, Happiness: 15, Reputation: -10},
			EffectB:  Effect{Military: 10},
			FlavorA:  "Invaders are confused and slightly delicious. Children keep sneaking snacks from the fortifications.",
			FlavorB:  "The moat is boring but effective. Fish move in. Someone starts charging for fishing rights.",
		},
		// 10
		{
			Question: "Neighboring kingdom challenges you to a poetry contest",
			OptionA:  "Accept! Our verses are superior!",
			OptionB:  "We fight with swords, not words",
			EffectA:  Effect{Culture: 20, Reputation: 15},
			EffectB:  Effect{Military: 10, Culture: -10},
			FlavorA:  "Your kingdom's entry: a 47-stanza epic about soup. The judges are moved to tears.",
			FlavorB:  "The neighboring kingdom writes a devastating limerick about your cowardice. It rhymes perfectly.",
		},
		// 11
		{
			Question: "The treasury is haunted. The accountant quit.",
			OptionA:  "Hire a ghost accountant",
			OptionB:  "Perform an exorcism",
			EffectA:  Effect{Treasury: 10, Culture: 5, Happiness: -5},
			EffectB:  Effect{Treasury: -15, Happiness: 10},
			FlavorA:  "The ghost accountant works 24/7 and never asks for a raise. Morale is... haunted.",
			FlavorB:  "The priest charges an exorcism fee. The ghost leaves a strongly-worded Yelp review.",
		},
		// 12
		{
			Question: "Citizens petition for mandatory nap time",
			OptionA:  "2pm naps for everyone!",
			OptionB:  "This is a kingdom, not a daycare",
			EffectA:  Effect{Happiness: 25, Military: -10, Food: -5},
			EffectB:  Effect{Happiness: -10, Military: 5},
			FlavorA:  "Productivity drops 30% but workplace satisfaction hits an all-time high. Snoring echoes through the streets.",
			FlavorB:  "The citizens grumble through their afternoon slump. Coffee bean imports triple.",
		},
		// 13
		{
			Question: "A group of bards wants to start a music festival",
			OptionA:  "Fund the 'Bardbonanza'!",
			OptionB:  "The budget doesn't have room for festivities",
			EffectA:  Effect{Treasury: -25, Culture: 30, Happiness: 20, FactionEffects: fe(map[Faction]int{Jesters: 20})},
			EffectB:  Effect{Culture: -5, FactionEffects: fe(map[Faction]int{Jesters: -10})},
			FlavorA:  "Three days of lute solos, mead, and questionable life choices. Best festival ever.",
			FlavorB:  "The bards play sad songs about budget cuts. Ironically, it's their best work.",
		},
		// 14
		{
			Question: "The royal chef demands exotic ingredients from overseas",
			OptionA:  "Spare no expense for the royal palate!",
			OptionB:  "Porridge was good enough for my ancestors",
			EffectA:  Effect{Treasury: -30, Happiness: 10, Reputation: 15},
			EffectB:  Effect{Happiness: -5, Treasury: 5},
			FlavorA:  "Dinner now includes 'deconstructed phoenix egg' and 'artisanal swamp foam.' The nobles applaud.",
			FlavorB:  "The chef quits in a huff. The new chef's specialty is 'bread, but sad.'",
		},
		// 15
		{
			Question: "Should we teach the army to dance?",
			OptionA:  "Dancing soldiers = intimidating soldiers!",
			OptionB:  "Train them to fight, not pirouette",
			EffectA:  Effect{Military: 5, Culture: 15, Happiness: 10},
			EffectB:  Effect{Military: 15},
			FlavorA:  "The synchronized battle waltz terrifies neighboring kingdoms. Morale through the roof.",
			FlavorB:  "The army drills in grim silence. They're effective but deeply boring at parties.",
		},
		// 16
		{
			Question: "Merchants propose opening trade with the Goblin Markets",
			OptionA:  "Profit is profit!",
			OptionB:  "Goblins can't be trusted",
			EffectA:  Effect{Treasury: 40, Reputation: -15, FactionEffects: fe(map[Faction]int{Merchants: 20})},
			EffectB:  Effect{Reputation: 5, FactionEffects: fe(map[Faction]int{Merchants: -10})},
			FlavorA:  "Gold flows in. Nobody asks where the 'mystery meat' shipments come from.",
			FlavorB:  "The merchants sulk. The goblins trade with your rivals instead. They seem suspiciously prosperous.",
		},
		// 17
		{
			Question: "The nobles demand that peasants bow three times instead of two",
			OptionA:  "Three bows it is!",
			OptionB:  "One bow is enough for anyone",
			EffectA:  Effect{Happiness: -15, FactionEffects: fe(map[Faction]int{Nobles: 25, Farmers: -20})},
			EffectB:  Effect{Happiness: 10, FactionEffects: fe(map[Faction]int{Nobles: -15, Farmers: 15})},
			FlavorA:  "The peasants develop chronic back problems. The nobles have never been more pleased.",
			FlavorB:  "The nobles gasp. The peasants stand a little taller. Someone writes a folk song about it.",
		},
		// 18
		{
			Question: "Someone has been putting googly eyes on all the statues",
			OptionA:  "It's art. Leave them.",
			OptionB:  "Find the culprit!",
			EffectA:  Effect{Culture: 10, Happiness: 15, FactionEffects: fe(map[Faction]int{Nobles: -10})},
			EffectB:  Effect{Military: -5, Happiness: -5},
			FlavorA:  "Tourism increases 200%. The statue of your grandfather looks perpetually surprised.",
			FlavorB:  "After months of investigation, the culprit is revealed: everyone. It was everyone.",
		},
		// 19
		{
			Question: "The kingdom's cats are organizing. They seem... strategic.",
			OptionA:  "Appoint a Cat Chancellor",
			OptionB:  "Increase the dog population as a counterbalance",
			EffectA:  Effect{Happiness: 20, Reputation: -15, Culture: 10},
			EffectB:  Effect{Food: -10, Military: 5},
			FlavorA:  "Chancellor Whiskers III passes sweeping nap legislation. Approval ratings soar.",
			FlavorB:  "The dogs are enthusiastic but disorganized. The cats seem amused. This concerns you.",
		},
		// 20
		{
			Question: "A prophet claims the world ends next Tuesday",
			OptionA:  "Throw a massive end-of-world party!",
			OptionB:  "Ignore them. Prophets are 0 for 347.",
			EffectA:  Effect{Treasury: -35, Happiness: 30, Food: -15},
			EffectB:  Effect{FactionEffects: fe(map[Faction]int{Scholars: 10})},
			FlavorA:  "Best. Party. Ever. Wednesday arrives and everyone has a hangover and existential clarity.",
			FlavorB:  "Tuesday passes without incident. The prophet mumbles about 'calendar discrepancies.'",
		},
		// 21
		{
			Question: "The sewers need renovation. It's expensive but... smelly.",
			OptionA:  "Fix the sewers!",
			OptionB:  "The people have noses. They'll adapt.",
			EffectA:  Effect{Treasury: -30, Happiness: 15, Population: 5},
			EffectB:  Effect{Happiness: -20, Population: -5},
			FlavorA:  "Clean water flows again. The sewer workers demand a parade. You grant it reluctantly.",
			FlavorB:  "Citizens develop a unique 'kingdom cologne.' Visitors from abroad are... not impressed.",
		},
		// 22
		{
			Question: "Pirates offer an alliance",
			OptionA:  "Pirates make great allies!",
			OptionB:  "We're a respectable kingdom",
			EffectA:  Effect{Military: 20, Reputation: -25, Treasury: 20},
			EffectB:  Effect{Reputation: 10, Military: -5},
			FlavorA:  "Your navy now says 'arrr' unironically. Trade routes are surprisingly secure though.",
			FlavorB:  "The pirates shrug and raid your coast anyway. At least you have the moral high ground.",
		},
		// 23
		{
			Question: "The kingdom's only bridge is trolled by an actual troll",
			OptionA:  "Pay the troll toll",
			OptionB:  "Build a second bridge",
			EffectA:  Effect{Treasury: -15, FactionEffects: fe(map[Faction]int{Merchants: 10})},
			EffectB:  Effect{Treasury: -40, Population: 10, FactionEffects: fe(map[Faction]int{Merchants: 15})},
			FlavorA:  "The troll is surprisingly professional about it. Issues receipts and everything.",
			FlavorB:  "The new bridge is lovely. The troll applies for unemployment benefits.",
		},
		// 24
		{
			Question: "Scholars discovered that the earth might be round",
			OptionA:  "Fascinating! Fund more research!",
			OptionB:  "Nonsense. Jail the scholars.",
			EffectA:  Effect{Treasury: -20, Culture: 25, FactionEffects: fe(map[Faction]int{Scholars: 20})},
			EffectB:  Effect{Culture: -15, Happiness: -5, FactionEffects: fe(map[Faction]int{Scholars: -30})},
			FlavorA:  "The scholars build a telescope. They immediately use it to spy on the neighboring kingdom.",
			FlavorB:  "The scholars revolt. They barricade themselves in the library and refuse to return overdue books.",
		},
		// 25
		{
			Question: "The army wants to replace swords with very sharp breadsticks",
			OptionA:  "Dual-purpose weapons! Genius!",
			OptionB:  "Keep the swords. Obviously.",
			EffectA:  Effect{Military: -10, Food: 5, Happiness: 15, Culture: 5},
			EffectB:  Effect{Military: 10},
			FlavorA:  "Soldiers can now fight AND have lunch. Enemy morale drops when hit with seasoned garlic bread.",
			FlavorB:  "The army grumbles. They were really looking forward to the garlic butter upgrades.",
		},
		// 26
		{
			Question: "A wizard offers to make it rain gold for one day",
			OptionA:  "What's the catch? ...Do it anyway!",
			OptionB:  "Wizards are trouble. Decline.",
			EffectA:  Effect{Treasury: 60, Happiness: -20, Culture: -10},
			EffectB:  Effect{FactionEffects: fe(map[Faction]int{Scholars: 5})},
			FlavorA:  "Gold rains down! Then the wizard's invoice arrives. Also everything smells like sulfur now.",
			FlavorB:  "The wizard turns your castle slightly pink out of spite. It fades in a week.",
		},
		// 27
		{
			Question: "The royal portrait painter made you look... unflattering",
			OptionA:  "Hang it with pride! Honesty matters.",
			OptionB:  "Off with their brush!",
			EffectA:  Effect{Culture: 15, Happiness: 5, Reputation: 10},
			EffectB:  Effect{Culture: -15, FactionEffects: fe(map[Faction]int{Jesters: 10})},
			FlavorA:  "The painting becomes the kingdom's most popular attraction. Gift shop sales are through the roof.",
			FlavorB:  "The painter's brushes are ceremonially snapped. They switch to abstract pottery out of spite.",
		},
		// 28
		{
			Question: "Neighboring kingdom asks to merge. They're terrible at everything.",
			OptionA:  "More territory! Absorb them!",
			OptionB:  "We have enough problems",
			EffectA:  Effect{Population: 30, Treasury: -20, Happiness: -10, Food: -15},
			EffectB:  Effect{Reputation: 10},
			FlavorA:  "Your kingdom doubles in size and problems. Their national sport was 'competitive complaining.'",
			FlavorB:  "They merge with the kingdom next door instead. That kingdom now has twice the incompetence.",
		},
		// 29
		{
			Question: "The royal food taster is on strike",
			OptionA:  "Taste it yourself. Assert dominance.",
			OptionB:  "Triple their salary",
			EffectA:  Effect{Happiness: 5, Reputation: 10, Military: -5},
			EffectB:  Effect{Treasury: -15, Happiness: 5},
			FlavorA:  "You survive the tasting. The food was fine. You feel invincible and mildly nauseous.",
			FlavorB:  "The food taster returns wearing a monocle. They now insist on being called 'Flavor Consultant.'",
		},
		// 30
		{
			Question: "Citizens request a public swimming pool",
			OptionA:  "The Royal Splash Zone!",
			OptionB:  "There's a perfectly good river right there",
			EffectA:  Effect{Treasury: -25, Happiness: 25, Culture: 5},
			EffectB:  Effect{Happiness: -10},
			FlavorA:  "The pool includes a waterslide shaped like you. Cannonballs are the kingdom's fastest growing sport.",
			FlavorB:  "Citizens swim in the river. Three lose their shoes to the current. One finds a sword.",
		},
		// 31
		{
			Question: "Wandering monks seek shelter in the kingdom",
			OptionA:  "Welcome them!",
			OptionB:  "We're full, try next kingdom",
			EffectA:  Effect{Culture: 15, Food: -10, Population: 5, FactionEffects: fe(map[Faction]int{Scholars: 10})},
			EffectB:  Effect{Reputation: -10},
			FlavorA:  "The monks share ancient wisdom. Also a really good recipe for lentil soup.",
			FlavorB:  "The monks leave, cursing softly. Your crops taste slightly worse for a season.",
		},
		// 32
		{
			Question: "The clocktower is 3 hours wrong. No one agrees which direction.",
			OptionA:  "Set it forward",
			OptionB:  "Set it backward",
			EffectA:  Effect{Happiness: -5, FactionEffects: fe(map[Faction]int{Merchants: 5, Farmers: -5})},
			EffectB:  Effect{Happiness: -5, FactionEffects: fe(map[Faction]int{Farmers: 5, Merchants: -5})},
			FlavorA:  "Merchants love the early start. Farmers miss breakfast. Nobody knows what time it actually is.",
			FlavorB:  "Farmers enjoy sleeping in. Merchants miss their appointments. The clockmaker drinks heavily.",
		},
		// 33
		{
			Question: "Should the kingdom anthem have more cowbell?",
			OptionA:  "Always more cowbell!",
			OptionB:  "The anthem is dignified as-is",
			EffectA:  Effect{Happiness: 15, Culture: 10, FactionEffects: fe(map[Faction]int{Nobles: -10})},
			EffectB:  Effect{FactionEffects: fe(map[Faction]int{Nobles: 5, Jesters: -5})},
			FlavorA:  "The new anthem slaps. Foreign dignitaries can't stop tapping their feet. The nobles wince rhythmically.",
			FlavorB:  "The anthem remains a 12-minute dirge about honor. The cowbell players form a support group.",
		},
		// 34
		{
			Question: "A sentient hedge maze has appeared outside the castle",
			OptionA:  "Charge admission!",
			OptionB:  "Burn it before it grows",
			EffectA:  Effect{Treasury: 20, Happiness: 10, Reputation: 5},
			EffectB:  Effect{Military: 5, Culture: -10},
			FlavorA:  "Tourists flock in. Only 60% find their way out. The rest seem happy in there though.",
			FlavorB:  "The hedge maze screams as it burns. Nobody sleeps well for a week.",
		},
		// 35
		{
			Question: "The royal treasurer has been replaced by three raccoons in a coat",
			OptionA:  "They're doing a great job. Don't question it.",
			OptionB:  "Hire an actual person",
			EffectA:  Effect{Treasury: 5, Happiness: 10},
			EffectB:  Effect{Treasury: -10},
			FlavorA:  "The raccoons implement a surprisingly effective savings plan. They insist on being paid in shiny things.",
			FlavorB:  "The new treasurer is competent but boring. The raccoons get jobs at the rival kingdom.",
		},
		// 36
		{
			Question: "Philosophers debate: is a hot dog a sandwich?",
			OptionA:  "Yes, obviously",
			OptionB:  "No, and anyone who says otherwise is banished",
			EffectA:  Effect{Culture: 5, FactionEffects: fe(map[Faction]int{Scholars: -10, Farmers: 5})},
			EffectB:  Effect{Population: -5, FactionEffects: fe(map[Faction]int{Scholars: 10})},
			FlavorA:  "The 'sandwich truth' faction gains power. All bread-wrapped foods are reclassified. Chaos in bakeries.",
			FlavorB:  "A small but vocal group is exiled. They found a settlement called 'Sandwich Freedom.'",
		},
		// 37
		{
			Question: "A bard wrote a song mocking you. It's very catchy.",
			OptionA:  "If you can't beat 'em, hum along",
			OptionB:  "Ban all music",
			EffectA:  Effect{Culture: 15, Happiness: 10, Reputation: -5},
			EffectB:  Effect{Culture: -25, Happiness: -20, FactionEffects: fe(map[Faction]int{Jesters: -30})},
			FlavorA:  "You're caught humming it in court. The bard gets a record deal. Your approval rating somehow rises.",
			FlavorB:  "The kingdom falls silent. People communicate through aggressive interpretive dance instead.",
		},
		// 38
		{
			Question: "Invent a new sport: competitive sheep stacking",
			OptionA:  "May the best stacker win!",
			OptionB:  "That's dangerous for the sheep",
			EffectA:  Effect{Happiness: 20, Food: -5, Culture: 10, FactionEffects: fe(map[Faction]int{Farmers: 10})},
			EffectB:  Effect{FactionEffects: fe(map[Faction]int{Farmers: 5})},
			FlavorA:  "The Sheep Stacking Championship draws crowds from across the land. The sheep seem indifferent.",
			FlavorB:  "The sheep appreciate your concern. They show this by not eating your garden. Briefly.",
		},
		// 39
		{
			Question: "The kingdom's flag is too boring",
			OptionA:  "Add a dragon riding a unicorn!",
			OptionB:  "Plain flags build character",
			EffectA:  Effect{Culture: 10, Military: 5, Reputation: 10, Treasury: -10},
			EffectB:  Effect{Culture: -5},
			FlavorA:  "The new flag is magnificent. Enemy armies pause to admire it mid-charge. Merch sales explode.",
			FlavorB:  "The flag remains a beige rectangle. Neighboring kingdoms struggle to identify it at summits.",
		},
		// 40
		{
			Question: "Tax season approaches. The people look nervous.",
			OptionA:  "Heavy taxes! The coffers must overflow!",
			OptionB:  "Tax holiday this year!",
			EffectA:  Effect{Treasury: 40, Happiness: -25, FactionEffects: fe(map[Faction]int{Farmers: -15, Merchants: -15})},
			EffectB:  Effect{Treasury: -30, Happiness: 30, FactionEffects: fe(map[Faction]int{Farmers: 15, Merchants: 15})},
			FlavorA:  "Gold pours in. Smiles pour out. The tax collectors need bodyguards to do grocery shopping.",
			FlavorB:  "The citizens celebrate wildly. The treasury echoes when you shout into it.",
		},
		// 41
		{
			Question: "A flock of geese has occupied the throne room",
			OptionA:  "Negotiate terms with the Goose Council",
			OptionB:  "Call in the army to reclaim the throne",
			EffectA:  Effect{Culture: 10, Happiness: 15, Reputation: -10, FactionEffects: fe(map[Faction]int{Nobles: -15})},
			EffectB:  Effect{Military: -10, Happiness: -5, Reputation: 5},
			FlavorA:  "The geese agree to a timeshare arrangement. Tuesdays and Thursdays are 'goose court.'",
			FlavorB:  "Three soldiers are hospitalized. Geese are vicious. The throne has bite marks now.",
		},
		// 42
		{
			Question: "A mysterious fog rolls in every night and rearranges the furniture",
			OptionA:  "Charge for 'immersive interior design experiences'",
			OptionB:  "Hire a wizard to investigate",
			EffectA:  Effect{Treasury: 15, Happiness: 10, Culture: 10},
			EffectB:  Effect{Treasury: -20, Culture: 5, FactionEffects: fe(map[Faction]int{Scholars: 10})},
			FlavorA:  "Nobles pay top coin for the 'Fog Feng Shui.' One wakes up with their bed on the roof.",
			FlavorB:  "The wizard concludes it's 'vibes.' They charge full price. The fog continues unimpressed.",
		},
		// 43
		{
			Question: "The kingdom's bees have learned to spell words with their formations",
			OptionA:  "Establish a Bee Postal Service",
			OptionB:  "This is concerning. Contain them.",
			EffectA:  Effect{Culture: 15, Treasury: 10, Happiness: 10, FactionEffects: fe(map[Faction]int{Merchants: 10})},
			EffectB:  Effect{Food: -10, Happiness: -5, Military: 5},
			FlavorA:  "Bee-mail is surprisingly reliable. Delivery stings are considered acceptable postage.",
			FlavorB:  "The containment fails. The bees spell 'we remember this' above the castle. Ominous.",
		},
		// 44
		{
			Question: "An underground fighting ring for vegetables has been discovered",
			OptionA:  "Legalize and regulate Veggie Fight Club",
			OptionB:  "Shut it down. Vegetables have rights.",
			EffectA:  Effect{Treasury: 20, Happiness: 15, Food: -10, FactionEffects: fe(map[Faction]int{Farmers: -10})},
			EffectB:  Effect{Food: 5, Culture: 5, FactionEffects: fe(map[Faction]int{Farmers: 10})},
			FlavorA:  "The Turnip vs. Carrot championship draws record crowds. The eggplant remains undefeated.",
			FlavorB:  "The vegetables are liberated. Farmers report they seem 'calmer' now.",
		},
		// 45
		{
			Question: "Your crown has started giving unsolicited advice",
			OptionA:  "Listen to the wise crown",
			OptionB:  "Wear a different hat",
			EffectA:  Effect{Culture: 5, Happiness: -10, FactionEffects: fe(map[Faction]int{Scholars: 15, Nobles: -10})},
			EffectB:  Effect{Reputation: -5, Happiness: 5},
			FlavorA:  "The crown's tax policy suggestions are actually brilliant. Its dating advice is not.",
			FlavorB:  "You switch to a floppy wizard hat. The crown sulks on its pillow, muttering about interest rates.",
		},
		// 46
		{
			Question: "A visiting dignitary accidentally insulted the kingdom's soup",
			OptionA:  "Declare a Soup War",
			OptionB:  "Accept the criticism and improve the recipe",
			EffectA:  Effect{Military: -5, Reputation: -10, Happiness: 15, Culture: 10},
			EffectB:  Effect{Food: 10, Reputation: 10, Culture: 5},
			FlavorA:  "Both kingdoms mobilize their finest chefs. Casualties are measured in burnt tongues.",
			FlavorB:  "The improved soup wins international acclaim. The dignitary sends a formal apology and a ladle.",
		},
	}
}
