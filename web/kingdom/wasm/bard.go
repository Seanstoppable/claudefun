package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// Bard narrates every event in rhyming verse — dramatic, medieval, and slightly ridiculous.
type Bard struct {
	rng         *rand.Rand
	kingdomName string
}

// NewBard creates a new Bard with the given random seed and kingdom name.
func NewBard(seed int64, kingdomName string) *Bard {
	return &Bard{
		rng:         rand.New(rand.NewSource(seed)),
		kingdomName: kingdomName,
	}
}

// pick returns a random element from a string slice.
func (b *Bard) pick(options []string) string {
	return options[b.rng.Intn(len(options))]
}

// fillKingdom replaces {kingdom} placeholders with the kingdom name.
func (b *Bard) fillKingdom(s string) string {
	return strings.ReplaceAll(s, "{kingdom}", b.kingdomName)
}

// NarratePolicy narrates a policy choice in rhyming couplets.
func (b *Bard) NarratePolicy(question string, chosen string, flavor string) string {
	reactions := []string{"cheer", "gasp", "weep", "shrug", "faint", "swoon", "jeer", "pray", "toast", "wince", "groan", "applaud"}
	consequences := []string{"hope", "dread", "change", "doubt", "fear", "pride", "rage", "peace", "shock", "glee", "grief", "awe"}
	nouns := []string{"tide", "wave", "storm", "drum", "flame", "bell", "horn", "wheel", "wind", "plague", "song", "march"}

	templates := []string{
		"The ruler pondered, stroked their chin,\n" +
			"And chose: '%s' — let it begin!\n" +
			"%s\n" +
			"The bard records this fateful call,\n" +
			"May it not lead to kingdom's fall!",

		"A choice was made in %s's hall,\n" +
			"'%s' echoed wall to wall!\n" +
			"The citizens did %s in turn,\n" +
			"As %s like a %s did churn.",

		"'Twas on this day the decree was set,\n" +
			"A choice the kingdom won't forget!\n" +
			"%s\n" +
			"The scrolls shall note what came to pass,\n" +
			"Whether this was wise... or crass.",

		"The question posed: '%s'\n" +
			"The ruler's answer: '%s' — oh well!\n" +
			"%s\n" +
			"And so the wheel of fate does turn,\n" +
			"Will %s prosper now, or burn?",

		"With quill in hand the scribe did write,\n" +
			"The ruler chose '%s' tonight!\n" +
			"%s\n" +
			"In %s the deed is done,\n" +
			"Another chapter has begun!",

		"The council asked: '%s'\n" +
			"The ruler said: '%s' — and off they went!\n" +
			"The people did %s without a care,\n" +
			"As %s billowed through the air.",

		"By royal decree, let it be known,\n" +
			"'%s' is the seed that's sown!\n" +
			"%s\n" +
			"The bard takes note, the ink still wet,\n" +
			"A choice the realm won't soon forget!",

		"The herald cried across the square,\n" +
			"'%s!' rang through %s's air!\n" +
			"The citizens did %s and sigh,\n" +
			"As %s rose up to the sky.",

		"Beneath the banners, tall and grand,\n" +
			"The ruler gave their firm command:\n" +
			"'%s' — so let it be!\n" +
			"%s\n" +
			"May history judge us graciously!",

		"The court assembled, hushed and still,\n" +
			"To hear the sovereign's iron will.\n" +
			"'%s' — the answer came,\n" +
			"%s\n" +
			"And %s shall never be the same!",
	}

	idx := b.rng.Intn(len(templates))
	switch idx {
	case 0:
		return fmt.Sprintf(templates[0], chosen, flavor)
	case 1:
		return fmt.Sprintf(templates[1],
			b.kingdomName, chosen,
			b.pick(reactions), b.pick(consequences), b.pick(nouns))
	case 2:
		return fmt.Sprintf(templates[2], flavor)
	case 3:
		return fmt.Sprintf(templates[3], question, chosen, flavor, b.kingdomName)
	case 4:
		return fmt.Sprintf(templates[4], chosen, flavor, b.kingdomName)
	case 5:
		return fmt.Sprintf(templates[5], question, chosen,
			b.pick(reactions), b.pick(consequences))
	case 6:
		return fmt.Sprintf(templates[6], chosen, flavor)
	case 7:
		return fmt.Sprintf(templates[7], chosen, b.kingdomName,
			b.pick(reactions), b.pick(nouns))
	case 8:
		return fmt.Sprintf(templates[8], chosen, flavor)
	case 9:
		return fmt.Sprintf(templates[9], chosen, flavor, b.kingdomName)
	}
	return ""
}

// NarrateState narrates the current kingdom state poetically.
func (b *Bard) NarrateState(treasury, happiness, military, culture, food, reputation int) string {
	var lines []string

	lines = append(lines, b.narrateTreasury(treasury))
	lines = append(lines, b.narrateHappiness(happiness))
	lines = append(lines, b.narrateMilitary(military))
	lines = append(lines, b.narrateCulture(culture))
	lines = append(lines, b.narrateFood(food))
	lines = append(lines, b.narrateReputation(reputation))

	header := fmt.Sprintf("~ The State of %s ~\n", b.kingdomName)
	return header + strings.Join(lines, "\n")
}

func (b *Bard) narrateTreasury(val int) string {
	high := []string{
		"The coffers gleam with mountains gold,\nMore wealth than castle walls can hold!",
		"Gold coins spill out across the floor,\nThe treasurer cries: 'There is no more... room!'",
		"The vaults are bursting at the seams,\nOur wealth exceeds our wildest dreams!",
		"So rich the kingdom, legend states,\nEven the moat has golden gates!",
		"The gold piles high as castle spires,\nOur wealth is all the world desires!",
		"We swim in coin like ducks in ponds,\nThe treasury overflows with bonds!",
		"A dragon's hoard would blush with shame,\nOur riches put that beast to blame!",
		"The royal mint works overtime,\nWe've gold enough to pave each climb!",
		"So wealthy are we, truth be told,\nThe servants' mops are made of gold!",
		"The merchants come from lands afar,\nTo glimpse our treasure's shining star!",
		"Our coffers sing a golden tune,\nWe're richer than the harvest moon!",
		"The vaults require a second floor,\nFor gold keeps pouring more and more!",
	}
	mid := []string{
		"The treasury sits at middling height,\nNot lavish, but we'll be alright.",
		"Some coins remain within the chest,\nNot great, not dire — call it 'blessed.'",
		"The budget's tight but holding fast,\nLet's hope our spending doesn't blast.",
		"A modest sum of gold we keep,\nEnough to eat, but not to leap.",
		"The coins we have could fill a boot,\nNot quite a tree, but there's some fruit.",
		"We're not yet poor, we're not yet rich,\nThe treasury's stuck in between the ditch.",
		"The gold trickles in, a modest stream,\nNot quite a nightmare, not a dream.",
		"Our wealth is like a middling stew,\nIt's warm enough to see us through.",
		"The treasurer shrugs, 'Could be much worse,'\nAt least there's something in the purse.",
		"We count our coins with careful hands,\nNot lavish wealth, but steady lands.",
		"The budget bends but doesn't break,\nWe've just enough to bake a cake.",
		"A reasonable sum of silver bright,\nNot day, not dark — a fiscal twilight.",
	}
	low := []string{
		"The treasury echoes, bare and cold,\nWhere once sat silver, now grows mold.",
		"The royal purse is looking thin,\nWe've barely got a coin within.",
		"Our gold has fled like morning dew,\nThe taxman weeps — there's naught to do.",
		"The coffers hold but dust and air,\nThe kingdom's wealth? Best not compare.",
		"A spider spins where gold once sat,\nThe treasury's home to just a rat.",
		"We owe the baker, owe the smith,\nOur fortune's now a fairy myth.",
		"The royal vault? An empty room,\nIt echoes with financial doom.",
		"Our debts pile high, our coins run dry,\nThe treasurer has learned to cry.",
		"If poverty were made of gold,\nWe'd be the richest realm of old.",
		"The crown jewels? Already pawned,\nOur golden age has come and gone.",
		"We've scraped the barrel, checked the floor,\nThe kingdom's wealth exists no more.",
		"The moths have eaten through the purse,\nOur fiscal state could not be worse.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateHappiness(val int) string {
	high := []string{
		"The people dance, the people sing,\nJoy resounds through everything!",
		"Such happiness! Such merry cheer!\nThe taverns overflow with beer!",
		"The citizens are smiling wide,\nContent and bursting full of pride!",
		"Laughter echoes lane to lane,\nNot a single soul would complain!",
		"The children play, the elders grin,\nWhat a glorious state we're in!",
		"From every window, songs take flight,\nThe kingdom glows with pure delight!",
		"The baker sings, the blacksmith hums,\nHappiness beats louder than drums!",
		"The market rings with joyful cries,\nBeneath the bluest, brightest skies!",
		"No frown exists in all the land,\nLife is merry, life is grand!",
		"The festivals run day and night,\nEvery heart is shining bright!",
		"Even the grumpiest old knight,\nWas seen to smile — what a sight!",
		"The kingdom's joy could fill the sea,\nNever have folk been this carefree!",
	}
	mid := []string{
		"The folk seem fine — not great, not grim,\nTheir happiness is on the brim.",
		"The people shrug and carry on,\nNot quite upset, not quite withdrawn.",
		"A tepid mood hangs in the air,\nThe people manage, somewhat fair.",
		"Content enough, though some complain,\nThe mood is... middling, let's be plain.",
		"They neither dance nor sit and weep,\nTheir joy runs shallow, not too deep.",
		"A yawn, a stretch, a halfhearted smile,\nThe people's mood is 'worth the while.'",
		"Some whistle tunes, some drag their feet,\nThe mood is lukewarm on the street.",
		"The taverns serve, the people sip,\nNo revolution on the tip.",
		"Not quite content, not quite distressed,\nThe people call it 'second best.'",
		"A shrug, a nod, a 'could be worse,'\nThe people clutch a neutral purse.",
		"The mood is steady, calm, and flat,\nLike porridge — and about as bland as that.",
		"Half the folk say things are fine,\nThe other half just sip their wine.",
	}
	low := []string{
		"The people grumble, spirits sink,\nThe kingdom teeters on the brink...",
		"The streets are dark, the mood is foul,\nYou'd think the moon replaced the owl.",
		"Despair creeps through the kingdom's heart,\nThe people's patience falls apart.",
		"No songs are sung, no dances done,\nThe people's joy has come undone.",
		"The market square sits still and bleak,\nThe future here looks rather meek.",
		"A scowl adorns each passing face,\nThis kingdom is a joyless place.",
		"The children cry, the elders moan,\nThe throne might as well be made of stone.",
		"Revolt is whispered lane to lane,\nThe people seethe with bitter pain.",
		"The tavern's dry, the jokes have died,\nThe people's hope has stepped outside.",
		"Each dawn arrives with heavy sighs,\nBeneath a kingdom's weeping skies.",
		"The court fool quit — that says it all,\nEven humor's hit a wall.",
		"Misery spreads from shore to shore,\nThe people simply can't take more.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateMilitary(val int) string {
	high := []string{
		"The army stands both proud and tall,\nNo foe shall breach this mighty wall!",
		"Our soldiers march in perfect rows,\nStrike fear into our kingdom's foes!",
		"A thousand swords gleam in the light,\nOur enemies retreat in fright!",
		"The military's strength is vast,\nNo siege shall hold, no war shall last!",
		"The generals smile with battle plans,\nVictory rests within our hands!",
		"Our fortress walls are thick as oaks,\nWe've archers, knights, and cannon strokes!",
		"The cavalry rides like thunder's crack,\nNo army dares to launch attack!",
		"Our barracks burst with eager might,\nEach soldier trained to win the fight!",
		"The war machines stand row on row,\nA fearsome and imposing show!",
		"Our battle standards fly so high,\nThey nearly scrape the very sky!",
		"The scouts report all borders clear,\nNo enemy would venture near!",
		"The armory gleams with polished steel,\nOur strength is one no foe can steal!",
	}
	mid := []string{
		"The guards patrol with modest zeal,\nThey've got some rust upon their steel.",
		"Our army's... present, one might say,\nThey mostly show up every day.",
		"The soldiers train — when they remember,\nMorale is warm, like dying ember.",
		"A decent force defends our gate,\nNot fearsome, but adequate.",
		"The knights have armor, slightly dented,\nTheir bravery is... supplemented.",
		"We've got a wall, we've got a moat,\nAnd one suspicious-looking goat.",
		"The archers aim with middling skill,\nThey hit the target... sometimes still.",
		"Our defenses hold, for now at least,\nWe're safe from man, if not from beast.",
		"The army drills on Tuesdays sharp,\nTheir battle cry's more like a carp.",
		"Enough to guard, not quite to conquer,\nOur military's not much bonker.",
		"The watch keeps watch, the scouts keep scouting,\nThe captain's confidence is... doubting.",
		"Our forces hold a middle ground,\nNot lost, not won — just hanging round.",
	}
	low := []string{
		"The guards are few, their swords are dull,\nDefense is... let's say, aspirational.",
		"Our 'army' is a tired old goat,\nAnd one brave guard in a torn coat.",
		"Invaders knock — we hide and pray,\nOur military ran away.",
		"The kingdom's might is... let me think...\nOne catapult. It's on the blink.",
		"The only shield we have is hope,\nOur guards can barely swing a rope.",
		"Our 'fortress' is a picket fence,\nOur 'strategy' makes zero sense.",
		"The armor's gone, the swords are sold,\nOur bravest knight is nine years old.",
		"A scarecrow guards the outer wall,\nThat is our army — that is all.",
		"The war horn sounds — nobody comes,\nThe soldiers fled to beat the drums... elsewhere.",
		"One rusty blade, one broken shield,\nWe'd lose a battle in a field.",
		"The enemy could conquer us,\nWith nothing but a horse and bus.",
		"Our defenses? A locked front door,\nA 'KEEP OUT' sign — and nothing more.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateCulture(val int) string {
	high := []string{
		"The arts do flourish, minds expand,\nThe finest culture in the land!",
		"Poets sing and painters thrive,\nNever has art been more alive!",
		"Our culture shines like morning dew,\nPhilosophers in every pew!",
		"The theaters are packed each night,\nOur bards perform to pure delight!",
		"Great sculptures grace the palace lawn,\nA renaissance has truly dawned!",
		"The poets duel with sharpened verse,\nEach line more brilliant than the first!",
		"Our opera house is world-renowned,\nThe finest voices can be found!",
		"The galleries stretch mile on mile,\nEach masterpiece provokes a smile!",
		"Scholars gather, wisdom grows,\nKnowledge blooms like summer's rose!",
		"The printing press works day and night,\nSpreading learning, spreading light!",
		"From music halls to lecture stages,\nOur culture spans a thousand ages!",
		"The kingdom's art is so refined,\nIt elevates the common mind!",
	}
	mid := []string{
		"Some art exists — a play or two,\nThe culture's getting halfway through.",
		"A wandering minstrel plays a tune,\nThe people hum, then lose it soon.",
		"The library grows, book upon book,\n(Though most are just on how to cook.)",
		"Culture persists in modest form,\nNot quite a drought, not quite a storm.",
		"A painter works upon the square,\nHis audience? A pigeon pair.",
		"The bard attempts a sonnet grand,\nThe crowd responds with lukewarm hands.",
		"One theater stands, it shows one play,\nThe same one, every single day.",
		"The scholars debate by candlelight,\nThey mostly argue who is right.",
		"A statue stands in middling pose,\nIs it a king? Nobody knows.",
		"Some poetry floats through the air,\nIt rhymes, at least — beyond that, fair.",
		"The culture's not completely dead,\nIt's merely sleeping in its bed.",
		"A modest muse inspires the town,\nNot lifting up, not dragging down.",
	}
	low := []string{
		"The library holds but one sad book,\n(And even that, the rats have took.)",
		"No art, no song, no tale is told,\nThe kingdom's culture? Stone age old.",
		"The finest art in all the realm?\nA stick-figure drawing at the helm.",
		"Culture here has hit rock bottom,\nAll the poets? Long forgotten.",
		"The last musician broke his lute,\nThe kingdom's culture follows suit.",
		"The theater collapsed last week,\nThe arts outlook is rather bleak.",
		"We tried to host a poetry night,\nThe only guest was candlelight.",
		"The greatest work of art we own?\nA slightly interesting-looking stone.",
		"Our culture's like a barren field,\nNo harvest here, no art revealed.",
		"The muses fled, the ink has dried,\nCreativity has truly died.",
		"The bard's own quill has lost its flair,\nThere's simply nothing worth to share.",
		"The kingdom's taste in art? A riddle,\nSomewhere between 'dire' and 'little.'",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateFood(val int) string {
	high := []string{
		"The granaries burst, the tables groan,\nNo belly empty, no hungry moan!",
		"A feast for all! The harvest grand!\nThe finest produce in the land!",
		"So much food the silos crack,\nThe royal chef has lost the track!",
		"Apples, bread, and roasted boar,\nThe people cry: 'We need no more!'",
		"The orchards bloom with fruit so sweet,\nEvery citizen has more than enough to eat!",
		"The markets overflow with grain,\nAbundance falls like golden rain!",
		"The royal kitchen never rests,\nFeeding all our honored guests!",
		"From cheese to cake to honeyed wine,\nThe kingdom's menu? Simply divine!",
		"The harvest festival runs all week,\nThe tables bend — the floorboards creak!",
		"Our cows are fat, our fields are green,\nThe plumpest realm you've ever seen!",
		"The pantries burst from floor to shelf,\nEven the scarecrow's fed himself!",
		"A cornucopia, left and right,\nEvery meal a sheer delight!",
	}
	mid := []string{
		"The food supply holds steady still,\nEnough to eat, but not to spill.",
		"The people dine on bread and stew,\nNot gourmet fare, but it'll do.",
		"The harvest came — a modest yield,\nNot empty barns, but not a filled field.",
		"We've food enough to see us through,\nThough seconds are for precious few.",
		"The soup is warm, the bread is fresh,\nBut nothing here will quite impress.",
		"The farmers work from dawn to dusk,\nThe crop is wheat — not gold, not musk.",
		"Potatoes, turnips, daily fare,\nNot quite a feast, but food is there.",
		"The rations hold, the stores suffice,\nWe eat our porridge — sometimes rice.",
		"The kitchen keeps a steady pace,\nFeeding folk with modest grace.",
		"No one starves, but no one feasts,\nWe're somewhere between men and beasts.",
		"The cellar's stocked with middling stores,\nEnough for now — but nothing more.",
		"Our meals are plain but fill the gut,\nThe pantry door stays firmly shut.",
	}
	low := []string{
		"The pantry's bare, the soup is thin,\nThe kingdom's belt pulled to its skin.",
		"The rats eat better than the court,\nOur food supplies are running short.",
		"One turnip left — we share it round,\nThe saddest meal the bard has found.",
		"Starvation knocks upon the door,\nThe kingdom's cupboards hold no more.",
		"The kitchen fire has long gone cold,\nThere's nothing left to cook or hold.",
		"The baker quit — no flour remains,\nJust empty shelves and hunger pains.",
		"Our finest dish? A bowl of air,\nGarnished with a side of despair.",
		"The livestock's gone, the crops have failed,\nThe kingdom's stomach has been jailed.",
		"The royal feast? A single pea,\nShared by the court — including me.",
		"The chickens flew, the pigs have fled,\nThere's not a crumb of daily bread.",
		"The harvest yielded dust and sand,\nFamine stretches cross the land.",
		"The people boil their leather shoes,\nThat's front-page kingdom dinner news.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateReputation(val int) string {
	high := []string{
		"Our name rings out through distant shores,\nAll kingdoms envy what is ours!",
		"Renowned! Respected! Known afar!\nOur reputation: shining star!",
		"The world looks on with awe and praise,\nOur kingdom's glory sets ablaze!",
		"From east to west, the legends spread,\nOf our great realm, much has been said!",
		"Ambassadors queue at our gate,\nTo dine with rulers truly great!",
		"The bards abroad all sing our name,\nOur kingdom's wrapped in lasting fame!",
		"Our banners fly in foreign halls,\nRespect resounds beyond our walls!",
		"The maps all mark us bold and bright,\nA beacon of renown and might!",
		"Travelers speak with reverent tone,\nOf this great land upon its throne!",
		"Our treaties stack from floor to sky,\nAll want to be our ally!",
		"The scholars write of our great deeds,\nOur fame is planted like the seeds!",
		"Our kingdom's name is on each tongue,\nThe greatest song that's ever sung!",
	}
	mid := []string{
		"Our reputation's 'fine, I guess,'—\nNeighbors know us, more or less.",
		"We're known enough in nearby towns,\nWe get some smiles, we get some frowns.",
		"The kingdom's name does gently float,\nOn lips of merchants, barely of note.",
		"Some know our name, some scratch their head,\n'Oh that place? Ah yes,' they said.",
		"We're on the map, if barely so,\nA middling realm where middling winds blow.",
		"The merchants trade, the envoys call,\nWe're known enough — but not to all.",
		"Our reputation's holding steady,\nNot quite famous, not quite petty.",
		"The neighbors nod, the traders wave,\nWe're not a joke, but not yet brave.",
		"A moderate name, a humble brand,\nWe're somewhat known across the land.",
		"We ring a bell in foreign courts,\nA footnote in the traders' reports.",
		"Not infamous, not yet revered,\nOur reputation's... middling-geared.",
		"The world has heard of us, perhaps,\nWe're somewhere in between the maps.",
	}
	low := []string{
		"Our reputation? Best not ask,\nDisguising shame's a daily task.",
		"The neighbors laugh, the traders scoff,\nOur kingdom's name? They've written it off.",
		"If reputations could be sold,\nOurs wouldn't fetch a coin of gold.",
		"Unknown, unloved, unsung, unmissed,\nWe're barely on the map — if that exists.",
		"The world has heard of us — and winced,\nOur reputation's left them unconvinced.",
		"Our envoys travel, none return,\nThe bridges of goodwill all burn.",
		"The foreign courts refuse our mail,\nOur kingdom's name is beyond pale.",
		"We're known for all the wrong affairs,\nThe world collectively... just stares.",
		"Our reputation hit the floor,\nThen somehow managed to sink more.",
		"The merchants pass us on their route,\nOur kingdom's name is in dispute.",
		"If shame were fame, we'd top the chart,\nOur reputation's fallen apart.",
		"The pigeons won't deliver here,\nEven birds avoid us, that is clear.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) pickByThreshold(val int, high, mid, low []string) string {
	switch {
	case val >= 70:
		return b.pick(high)
	case val >= 35:
		return b.pick(mid)
	default:
		return b.pick(low)
	}
}

// NarrateEvent narrates a random event in dramatic verse.
func (b *Bard) NarrateEvent(eventDesc string) string {
	templates := []string{
		"Hear ye! Hear ye! News has come,\n" +
			"%s!\n" +
			"The bard strikes a chord upon his drum,\n" +
			"What fortune or disaster — the tale's just begun!",

		"A twist of fate! The heavens shake!\n" +
			"%s!\n" +
			"Is this a blessing or mistake?\n" +
			"The bard records it, for the kingdom's sake!",

		"Stop the presses! Hold the ale!\n" +
			"%s!\n" +
			"The bard must now update the tale,\n" +
			"Of %s — through triumph and travail!",

		"Lo! An event most unexpected!\n" +
			"%s!\n" +
			"The kingdom's course has been redirected,\n" +
			"Was this by fate or chance selected?",

		"Breaking news from %s's gate!\n" +
			"%s!\n" +
			"The bard scrambles to narrate,\n" +
			"This most peculiar twist of fate!",

		"The ravens bring a startling word,\n" +
			"%s!\n" +
			"The strangest thing the realm has heard,\n" +
			"The bard takes note — how quite absurd!",

		"A commotion stirs within the square!\n" +
			"%s!\n" +
			"The townsfolk gawk, the townsfolk stare,\n" +
			"What happens next? Does anyone care?",

		"Hark! A happening of note!\n" +
			"%s!\n" +
			"The scribe adjusts his well-worn coat,\n" +
			"And adds this chapter — freshly wrote!",

		"The crystal ball reveals the scene,\n" +
			"%s!\n" +
			"The most bizarre there's ever been,\n" +
			"In all of %s's halls so keen!",

		"Thunder rumbles! Lightning cracks!\n" +
			"%s!\n" +
			"The bard adjusts his almanacs,\n" +
			"And jots this down with quill and wax!",
	}

	idx := b.rng.Intn(len(templates))
	switch idx {
	case 2:
		return fmt.Sprintf(templates[2], eventDesc, b.kingdomName)
	case 4:
		return fmt.Sprintf(templates[4], b.kingdomName, eventDesc)
	case 8:
		return fmt.Sprintf(templates[8], eventDesc, b.kingdomName)
	default:
		return fmt.Sprintf(templates[idx], eventDesc)
	}
}

// NarrateGameOver delivers the final narration — victory or defeat.
func (b *Bard) NarrateGameOver(victory bool, reason string) string {
	if victory {
		return b.narrateVictory(reason)
	}
	return b.narrateDefeat(reason)
}

func (b *Bard) narrateVictory(reason string) string {
	templates := []string{
		"Rejoice! Rejoice! The tale is told,\n" +
			"Of %s, legendary and bold!\n" +
			"%s\n" +
			"Through chaos, laughs, and bitter strife,\n" +
			"The ruler earned eternal life!\n" +
			"(Not literally. That's a different game.)",

		"The trumpets sound! The banners wave!\n" +
			"%s stands triumphant, strong, and brave!\n" +
			"%s\n" +
			"Let all the land rejoice and cheer,\n" +
			"The greatest ruler of the year!",

		"And so concludes this glorious tale,\n" +
			"Of %s — which did not fail!\n" +
			"%s\n" +
			"The bard bows deep, the crowd goes wild,\n" +
			"Victory! (The narrator smiled.)",

		"Huzzah! Huzzah! The deed is done!\n" +
			"%s has gloriously won!\n" +
			"%s\n" +
			"The scrolls record this legendary reign,\n" +
			"May such a ruler rise again!",

		"The golden age of %s is here!\n" +
			"%s\n" +
			"The people raise a thundering cheer,\n" +
			"A triumph sung from ear to ear!",

		"What glory shines on %s tonight!\n" +
			"%s\n" +
			"The stars themselves burn twice as bright,\n" +
			"In honor of this sovereign's might!",

		"The chronicles of %s are penned,\n" +
			"A story with a glorious end!\n" +
			"%s\n" +
			"The bard retires, his duty done,\n" +
			"The greatest tale beneath the sun!",

		"A legend born in %s's halls,\n" +
			"%s\n" +
			"The victory echoes off the walls,\n" +
			"And into history's great annals falls!",

		"The feast is set, the wine is poured,\n" +
			"%s has triumphed! Praise the lord!\n" +
			"%s\n" +
			"This ruler's name shall echo long,\n" +
			"In every tale and every song!",

		"Three cheers for %s, mighty and true!\n" +
			"%s\n" +
			"The bard salutes with verses new,\n" +
			"A finer reign the world ne'er knew!",
	}

	t := b.pick(templates)
	return fmt.Sprintf(t, b.kingdomName, reason)
}

func (b *Bard) narrateDefeat(reason string) string {
	templates := []string{
		"Alas, alas, the kingdom fell,\n" +
			"%s\n" +
			"The bard packs up, with tales to tell,\n" +
			"Of %s's rise... and fare-thee-well.",

		"And so it ends — not with a cheer,\n" +
			"But with a sigh and silent tear.\n" +
			"%s\n" +
			"The tale of %s ends right here.",

		"The curtain falls, the torches dim,\n" +
			"The kingdom's fate was rather grim.\n" +
			"%s\n" +
			"Farewell, dear %s — on fortune's whim.",

		"A tragedy! A woeful plight!\n" +
			"%s brought %s's final night.\n" +
			"The bard weeps softly, quill in hand,\n" +
			"And writes 'THE END' across the land.",

		"The bells toll low, the flags hang slack,\n" +
			"%s\n" +
			"There is no road that leads us back,\n" +
			"To %s's former golden track.",

		"The story ends in ash and sorrow,\n" +
			"%s\n" +
			"There'll be no brighter dawn tomorrow,\n" +
			"For %s — no time left to borrow.",

		"And so the sun sets on this reign,\n" +
			"%s\n" +
			"The bard records the bitter pain,\n" +
			"Of %s — ne'er to rise again.",

		"The chronicles conclude with woe,\n" +
			"%s\n" +
			"Where once stood glory, now hangs low,\n" +
			"The tattered flag of %s's show.",

		"A somber tune, a final verse,\n" +
			"%s\n" +
			"Could things have gone from bad to worse?\n" +
			"For %s, it seems — a royal curse.",

		"The throne sits empty, cold, and bare,\n" +
			"%s\n" +
			"The people wander in despair,\n" +
			"Once %s stood proud — now, nothing's there.",
	}

	idx := b.rng.Intn(len(templates))
	switch idx {
	case 3:
		return fmt.Sprintf(templates[3], reason, b.kingdomName)
	default:
		return fmt.Sprintf(templates[idx], reason, b.kingdomName)
	}
}

// NarrateWelcome delivers the opening narration for a new game.
func (b *Bard) NarrateWelcome(kingdomName string) string {
	templates := []string{
		"Welcome, O Ruler, to %s fair!\n" +
			"(Well, 'fair' is generous, but we're getting there.)\n" +
			"Your kingdom awaits your wise decree,\n" +
			"Or terrible ones — we'll just have to see!\n\n" +
			"Treasury: modest. Army: questionable.\n" +
			"Happiness: fragile. Outlook: debatable.\n" +
			"But fear not! For with you at the helm,\n" +
			"Nothing can go wrong in this realm!\n" +
			"...Right?",

		"Hear ye! Hear ye! Gather 'round!\n" +
			"A brand new ruler has been found!\n" +
			"Welcome to %s, land of... potential,\n" +
			"Your leadership is now essential!\n\n" +
			"The people watch with bated breath,\n" +
			"(Half expecting certain death.)\n" +
			"But never mind their worried faces,\n" +
			"Let's rule this kingdom through its paces!",

		"The throne awaits in %s grand,\n" +
			"(Well, 'grand' might be a stretch, but understand—)\n" +
			"This humble realm is yours to lead,\n" +
			"Through every triumph, every deed!\n\n" +
			"The castle leaks, the army's small,\n" +
			"The treasury's got... well, not much at all.\n" +
			"But with your wisdom, sharp and keen,\n" +
			"This'll be the best reign ever seen!\n" +
			"...Probably.",

		"A new dawn breaks o'er %s's hills,\n" +
			"A ruler arrives! The kingdom thrills!\n" +
			"(Or trembles slightly — hard to say,\n" +
			"Your reputation's TBD today.)\n\n" +
			"The peasants wave, the nobles stare,\n" +
			"The court jester trips upon a chair.\n" +
			"But pay them no mind, O Sovereign bright,\n" +
			"Your kingdom needs you! ...Starting tonight.",

		"Behold! The gates of %s swing wide,\n" +
			"A brand new ruler steps inside!\n" +
			"The cobblestones are slightly cracked,\n" +
			"The royal welcome somewhat lacked—\n\n" +
			"But worry not! Great things await!\n" +
			"A kingdom poised to meet its fate!\n" +
			"With you in charge, the sky's the limit,\n" +
			"(Unless the budget's in it.)",

		"Greetings, Sovereign! Step right in,\n" +
			"To %s — where your tale begins!\n" +
			"The kingdom's seen much better days,\n" +
			"But your arrival earns some praise!\n\n" +
			"The cook is ready, the guard's awake,\n" +
			"The bard has tuned his lute for your sake.\n" +
			"So take the crown, assume the throne,\n" +
			"This realm is now your very own!",

		"The stars aligned, the prophecy told,\n" +
			"A ruler comes to %s bold!\n" +
			"(The prophecy was rather vague,\n" +
			"It also mentioned a minor plague.)\n\n" +
			"But here you are! The people cheer!\n" +
			"...Or at least they will, once they appear.\n" +
			"Take heart! Take charge! The realm's in need,\n" +
			"Of someone brave to plant the seed!",

		"All hail the ruler of %s!\n" +
			"(No pressure, but we're counting on your plans.)\n" +
			"The castle's drafty, the moat's a mess,\n" +
			"The kingdom's in a state of... progress?\n\n" +
			"But every legend starts somewhere small,\n" +
			"Perhaps a crumbling castle wall.\n" +
			"So raise your banner, claim your crown,\n" +
			"And try not to burn the kingdom down!",
	}

	return fmt.Sprintf(b.pick(templates), kingdomName)
}

// NarrateTurnStart announces the beginning of a new turn.
func (b *Bard) NarrateTurnStart(turn int) string {
	templates := []string{
		"Turn %d dawns upon the land,\n" +
			"The fate of %s in your hand!",

		"Day %d — the sun peeks through the clouds,\n" +
			"The people gather, forming crowds...",

		"Chapter %d of our kingdom's tale,\n" +
			"Will it triumph? Will it fail?",

		"And so begins turn number %d,\n" +
			"In %s, where legends never end!",

		"The rooster crows! 'Tis turn %d!\n" +
			"%s stirs — what will be done?",

		"Turn %d! The plot, it thickens so,\n" +
			"Where %s leads, we all shall go!",

		"The herald calls: 'Tis turn %d!\n" +
			"Another day in %s has begun!",

		"Page %d of the royal decree,\n" +
			"What shall the ruler's next move be?",

		"The clock strikes forward — turn %d rings,\n" +
			"The bard awaits what fortune brings!",

		"A new morn breaks — turn %d is nigh,\n" +
			"Beneath the %s banner high!",

		"The scroll unrolls to turn %d now,\n" +
			"The kingdom waits — but when and how?",

		"Turn %d arrives with trumpet's blare,\n" +
			"What fate awaits in %s fair?",
	}

	idx := b.rng.Intn(len(templates))
	switch idx {
	case 0, 3, 4, 5, 6, 9, 11:
		return fmt.Sprintf(templates[idx], turn, b.kingdomName)
	default:
		return fmt.Sprintf(templates[idx], turn)
	}
}
