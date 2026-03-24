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
	reactions := []string{"cheer", "gasp", "weep", "shrug", "faint"}
	consequences := []string{"hope", "dread", "change", "doubt"}
	nouns := []string{"tide", "wave", "storm", "drum"}

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
	}
	mid := []string{
		"The treasury sits at middling height,\nNot lavish, but we'll be alright.",
		"Some coins remain within the chest,\nNot great, not dire — call it 'blessed.'",
		"The budget's tight but holding fast,\nLet's hope our spending doesn't blast.",
		"A modest sum of gold we keep,\nEnough to eat, but not to leap.",
	}
	low := []string{
		"The treasury echoes, bare and cold,\nWhere once sat silver, now grows mold.",
		"The royal purse is looking thin,\nWe've barely got a coin within.",
		"Our gold has fled like morning dew,\nThe taxman weeps — there's naught to do.",
		"The coffers hold but dust and air,\nThe kingdom's wealth? Best not compare.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateHappiness(val int) string {
	high := []string{
		"The people dance, the people sing,\nJoy resounds through everything!",
		"Such happiness! Such merry cheer!\nThe taverns overflow with beer!",
		"The citizens are smiling wide,\nContent and bursting full of pride!",
		"Laughter echoes lane to lane,\nNot a single soul would complain!",
	}
	mid := []string{
		"The folk seem fine — not great, not grim,\nTheir happiness is on the brim.",
		"The people shrug and carry on,\nNot quite upset, not quite withdrawn.",
		"A tepid mood hangs in the air,\nThe people manage, somewhat fair.",
		"Content enough, though some complain,\nThe mood is... middling, let's be plain.",
	}
	low := []string{
		"The people grumble, spirits sink,\nThe kingdom teeters on the brink...",
		"The streets are dark, the mood is foul,\nYou'd think the moon replaced the owl.",
		"Despair creeps through the kingdom's heart,\nThe people's patience falls apart.",
		"No songs are sung, no dances done,\nThe people's joy has come undone.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateMilitary(val int) string {
	high := []string{
		"The army stands both proud and tall,\nNo foe shall breach this mighty wall!",
		"Our soldiers march in perfect rows,\nStrike fear into our kingdom's foes!",
		"A thousand swords gleam in the light,\nOur enemies retreat in fright!",
		"The military's strength is vast,\nNo siege shall hold, no war shall last!",
	}
	mid := []string{
		"The guards patrol with modest zeal,\nThey've got some rust upon their steel.",
		"Our army's... present, one might say,\nThey mostly show up every day.",
		"The soldiers train — when they remember,\nMorale is warm, like dying ember.",
		"A decent force defends our gate,\nNot fearsome, but adequate.",
	}
	low := []string{
		"The guards are few, their swords are dull,\nDefense is... let's say, aspirational.",
		"Our 'army' is a tired old goat,\nAnd one brave guard in a torn coat.",
		"Invaders knock — we hide and pray,\nOur military ran away.",
		"The kingdom's might is... let me think...\nOne catapult. It's on the blink.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateCulture(val int) string {
	high := []string{
		"The arts do flourish, minds expand,\nThe finest culture in the land!",
		"Poets sing and painters thrive,\nNever has art been more alive!",
		"Our culture shines like morning dew,\nPhilosophers in every pew!",
		"The theaters are packed each night,\nOur bards perform to pure delight!",
	}
	mid := []string{
		"Some art exists — a play or two,\nThe culture's getting halfway through.",
		"A wandering minstrel plays a tune,\nThe people hum, then lose it soon.",
		"The library grows, book upon book,\n(Though most are just on how to cook.)",
		"Culture persists in modest form,\nNot quite a drought, not quite a storm.",
	}
	low := []string{
		"The library holds but one sad book,\n(And even that, the rats have took.)",
		"No art, no song, no tale is told,\nThe kingdom's culture? Stone age old.",
		"The finest art in all the realm?\nA stick-figure drawing at the helm.",
		"Culture here has hit rock bottom,\nAll the poets? Long forgotten.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateFood(val int) string {
	high := []string{
		"The granaries burst, the tables groan,\nNo belly empty, no hungry moan!",
		"A feast for all! The harvest grand!\nThe finest produce in the land!",
		"So much food the silos crack,\nThe royal chef has lost the track!",
		"Apples, bread, and roasted boar,\nThe people cry: 'We need no more!'",
	}
	mid := []string{
		"The food supply holds steady still,\nEnough to eat, but not to spill.",
		"The people dine on bread and stew,\nNot gourmet fare, but it'll do.",
		"The harvest came — a modest yield,\nNot empty barns, but not a filled field.",
		"We've food enough to see us through,\nThough seconds are for precious few.",
	}
	low := []string{
		"The pantry's bare, the soup is thin,\nThe kingdom's belt pulled to its skin.",
		"The rats eat better than the court,\nOur food supplies are running short.",
		"One turnip left — we share it round,\nThe saddest meal the bard has found.",
		"Starvation knocks upon the door,\nThe kingdom's cupboards hold no more.",
	}
	return b.pickByThreshold(val, high, mid, low)
}

func (b *Bard) narrateReputation(val int) string {
	high := []string{
		"Our name rings out through distant shores,\nAll kingdoms envy what is ours!",
		"Renowned! Respected! Known afar!\nOur reputation: shining star!",
		"The world looks on with awe and praise,\nOur kingdom's glory sets ablaze!",
		"From east to west, the legends spread,\nOf our great realm, much has been said!",
	}
	mid := []string{
		"Our reputation's 'fine, I guess,'—\nNeighbors know us, more or less.",
		"We're known enough in nearby towns,\nWe get some smiles, we get some frowns.",
		"The kingdom's name does gently float,\nOn lips of merchants, barely of note.",
		"Some know our name, some scratch their head,\n'Oh that place? Ah yes,' they said.",
	}
	low := []string{
		"Our reputation? Best not ask,\nDisguising shame's a daily task.",
		"The neighbors laugh, the traders scoff,\nOur kingdom's name? They've written it off.",
		"If reputations could be sold,\nOurs wouldn't fetch a coin of gold.",
		"Unknown, unloved, unsung, unmissed,\nWe're barely on the map — if that exists.",
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
	}

	idx := b.rng.Intn(len(templates))
	switch idx {
	case 2:
		return fmt.Sprintf(templates[2], eventDesc, b.kingdomName)
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
	}

	idx := b.rng.Intn(len(templates))
	switch idx {
	case 0, 3, 4, 5:
		return fmt.Sprintf(templates[idx], turn, b.kingdomName)
	default:
		return fmt.Sprintf(templates[idx], turn)
	}
}
