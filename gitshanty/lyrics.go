package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// ShantyType describes the musical mood of the generated shanty.
type ShantyType int

const (
	CelebrationJig ShantyType = iota // features, releases
	MournfulBallad                   // bad merges, reverts, bugs
	EpicSaga                         // big refactors, large changes
	MutinyAnthem                     // force pushes, delete sprees
	WorkSong                         // regular commits, steady work
	LullabyCalm                      // docs, tiny changes, tests
)

// Shanty is the generated sea shanty with all its parts.
type Shanty struct {
	Title    string
	Type     ShantyType
	Tempo    string
	Verses   []string
	Chorus   string
	Bridge   string // optional
	Coda     string
	BandName string
}

// LyricsGenerator builds sea shanties from commit data.
type LyricsGenerator struct {
	rng *rand.Rand
}

// NewLyricsGenerator creates a generator with a deterministic seed.
func NewLyricsGenerator(seed int64) *LyricsGenerator {
	return &LyricsGenerator{rng: rand.New(rand.NewSource(seed))}
}

func (l *LyricsGenerator) pick(items []string) string {
	return items[l.rng.Intn(len(items))]
}

// ---------- band names ----------

var bandNames = []string{
	"The Salty Semicolons",
	"The Dangling Pointers",
	"Davy Jones' Code Review",
	"The Shipwrecked Sprints",
	"Blackbeard's Binary",
	"The Mutinous Maintainers",
	"Port & Starboard (The DevOps Duo)",
	"The Kraken Commits",
	"Dead Man's Deploy",
	"The Rusty Anchors of CI",
	"The Phantom Linters",
	"Yo Ho YAML",
}

func (l *LyricsGenerator) bandNameFor(author string) string {
	personalBands := []string{
		fmt.Sprintf("Captain %s & The Merge Conflicts", author),
		fmt.Sprintf("%s and the Loose Threads", author),
		fmt.Sprintf("The %s Experience (Unplugged)", author),
	}
	all := append([]string{}, bandNames...)
	all = append(all, personalBands...)
	return l.pick(all)
}

// ---------- titles ----------

func (l *LyricsGenerator) generateTitle(c Commit) string {
	adjectives := []string{"Last", "Great", "Terrible", "Magnificent", "Doomed", "Fearsome", "Legendary", "Forgotten"}
	nouns := []string{"Merge", "Commit", "Refactor", "Push", "Deploy", "Rebase", "Rollback"}
	drunkNouns := []string{"Drunken Merge", "Broken Pipeline", "Missing Semicolon", "Rogue Dependency", "Flaky Test", "Stale Branch"}
	events := []string{"Sinking", "Mutiny", "Plundering", "Maiden Voyage", "Last Stand", "Great Rewrite"}

	templates := []func() string{
		func() string { return fmt.Sprintf("The Ballad of %s", c.ShortHash) },
		func() string {
			return fmt.Sprintf("%s's %s %s", c.Author, l.pick(adjectives), l.pick(nouns))
		},
		func() string { return "Fifteen Files on a Dead Man's Branch" },
		func() string { return fmt.Sprintf("The %s of the Repo", l.pick(events)) },
		func() string { return fmt.Sprintf("What Shall We Do with a %s?", l.pick(drunkNouns)) },
		func() string { return fmt.Sprintf("A Shanty for %s", c.ShortHash) },
		func() string {
			return fmt.Sprintf("%d Lines to Davy Jones' Locker", c.Insertions+c.Deletions)
		},
		func() string { return fmt.Sprintf("The Night %s Pushed to Main", c.Author) },
	}
	return templates[l.rng.Intn(len(templates))]()
}

// ---------- tempo ----------

var tempos = map[ShantyType][]string{
	CelebrationJig: {"♩ = 140, With great merriment!", "♩ = 138, Joyously!", "♩ = 144, With tankards raised!"},
	MournfulBallad: {"♩ = 60, With deep sorrow...", "♩ = 56, Weeping softly...", "♩ = 64, Like a funeral at sea..."},
	EpicSaga:       {"♩ = 100, Grand and sweeping!", "♩ = 96, With orchestral grandeur!", "♩ = 104, Heroically!"},
	MutinyAnthem:   {"♩ = 160, FURIOUSLY!", "♩ = 168, WITH RECKLESS ABANDON!", "♩ = 155, Like there's no tomorrow!"},
	WorkSong:       {"♩ = 110, Steady as she goes", "♩ = 108, With quiet determination", "♩ = 112, One commit at a time"},
	LullabyCalm:    {"♩ = 70, Gently, gently...", "♩ = 66, Whispered into the void...", "♩ = 72, Like a linting breeze..."},
}

// ---------- ShantyTypeForCommit ----------

// ShantyTypeForCommit picks the shanty mood from a set of commit types.
func ShantyTypeForCommit(types []CommitType) ShantyType {
	has := func(ct CommitType) bool {
		for _, t := range types {
			if t == ct {
				return true
			}
		}
		return false
	}

	switch {
	case has(ForceEvent) || has(DeleteSpree):
		return MutinyAnthem
	case has(Feature) || has(InitialCommit):
		return CelebrationJig
	case has(MergeCommit) && has(Bugfix), has(Revert), has(HotfixUrgent):
		return MournfulBallad
	case has(Bugfix):
		return MournfulBallad
	case has(Refactor) || has(BigChange):
		return EpicSaga
	case has(Documentation) || has(TinyChange) || has(TestCommit):
		return LullabyCalm
	default:
		return WorkSong
	}
}

// ---------- verse / chorus builders ----------

func (l *LyricsGenerator) celebrationChorus(c Commit) string {
	templates := []string{
		fmt.Sprintf(
			"Heave ho, the feature's done!\n%s coded through the night!\n%d files changed 'neath the midnight sun,\nAnd the tests are running right!",
			c.Author, c.FilesChanged),
		fmt.Sprintf(
			"Raise the flag, the build is GREEN!\n%s sailed her into port!\nThe finest feature ever seen,\nA %d-file sort!",
			c.Author, c.FilesChanged),
	}
	return l.pick(templates)
}

func (l *LyricsGenerator) celebrationVerses(c Commit) []string {
	pool := []string{
		fmt.Sprintf(
			"'Twas on the branch of main we sailed,\nWhen %s raised the flag,\n\"%s\" the captain wailed,\nAnd not a soul did lag!",
			c.Author, truncMsg(c.Message, 40)),
		fmt.Sprintf(
			"They wrote the code with steady hand,\n%d lines of gold,\nThe finest feature in the land,\nA story to be told!",
			c.Insertions),
		fmt.Sprintf(
			"The PR was approved at dawn,\n%s merged with glee,\nThe old code dead, the new code born,\nAs green as any sea!",
			c.Author),
		fmt.Sprintf(
			"With %d insertions, bold and bright,\nAnd %d deletions cast aside,\nThe codebase gleamed in morning light,\nA developer's pride!",
			c.Insertions, c.Deletions),
	}
	return l.pickN(pool, 2)
}

func (l *LyricsGenerator) mournfulChorus(c Commit) string {
	templates := []string{
		fmt.Sprintf(
			"Oh, the merge went wrong, the merge went wrong,\nThe conflicts piled up high,\n%s sang a mournful song,\nAnd let a tear roll by...",
			c.Author),
		fmt.Sprintf(
			"Oh, the build is RED, the build is RED,\n%s hangs their head,\n%d lines of code are surely dead,\nThe pipeline filled with dread...",
			c.Author, c.Deletions),
	}
	return l.pick(templates)
}

func (l *LyricsGenerator) mournfulVerses(c Commit) []string {
	pool := []string{
		fmt.Sprintf(
			"In the depths of branch main,\nWhere merge conflicts dwell,\n%s tried to make a stand,\nBut the diff was straight from hell.",
			c.Author),
		fmt.Sprintf(
			"%d lines were cast away,\nLike sailors lost at sea,\n\"Revert, revert!\" we heard them say,\nBut it was not to be...",
			c.Deletions),
		fmt.Sprintf(
			"The CI pipeline screamed in pain,\nAs %s hit 'confirm',\nThe staging server, once so sane,\nNow writhing like a worm.",
			c.Author),
		fmt.Sprintf(
			"\"%s\" — the commit that broke the dam,\nA %d-file catastrophe,\nNow every dev must give a damn,\nAbout this tragedy.",
			truncMsg(c.Message, 30), c.FilesChanged),
	}
	return l.pickN(pool, 2)
}

func (l *LyricsGenerator) epicChorus(c Commit) string {
	templates := []string{
		fmt.Sprintf(
			"Refactor! Refactor! Across the codebase wide!\n%d files fell before the tide!\n%s stood with pride inside,\nThe greatest diff you ever spied!",
			c.FilesChanged, c.Author),
		fmt.Sprintf(
			"Onward! Through the legacy code!\n%s carries the load!\n%d insertions, a mighty ode,\nTo the refactoring road!",
			c.Author, c.Insertions),
	}
	return l.pick(templates)
}

func (l *LyricsGenerator) epicVerses(c Commit) []string {
	pool := []string{
		fmt.Sprintf(
			"In days of old, the code was cruel,\nSpaghetti ruled the land,\nThen %s picked up the tool,\nAnd refactored, grand!",
			c.Author),
		fmt.Sprintf(
			"%d files changed in a single sweep,\nThe architects looked on in awe,\nThrough layers tangled, dark, and deep,\nThe new design held no flaw.",
			c.FilesChanged),
		fmt.Sprintf(
			"They say that %s worked for days,\nRewriting line by line,\n%d insertions in a blaze,\nThe codebase, now divine!",
			c.Author, c.Insertions),
		fmt.Sprintf(
			"The old abstractions crumbled down,\nLike walls before the sea,\n%d deletions wore the crown,\nOf what once used to be.",
			c.Deletions),
	}
	return l.pickN(pool, 2)
}

func (l *LyricsGenerator) mutinyChorus(c Commit) string {
	templates := []string{
		fmt.Sprintf(
			"FORCE PUSH TO MAIN! FORCE PUSH TO MAIN!\nThe history is gone, boys, never seen again!\n%s lit the match and laughed through the pain,\nFORCE PUSH TO MAIN! FORCE PUSH TO MAIN!",
			c.Author),
		fmt.Sprintf(
			"DELETE IT ALL! DELETE IT ALL!\n%s answered the call!\n%d lines fell, both big and small,\nDELETE IT ALL! DELETE IT ALL!",
			c.Author, c.Deletions),
	}
	return l.pick(templates)
}

func (l *LyricsGenerator) mutinyVerses(c Commit) []string {
	pool := []string{
		fmt.Sprintf(
			"They said \"don't force push,\" we said \"WATCH US NOW!\"\n%d lines deleted with a solemn vow,\nThe git log wept, the CI server bowed,\n%s stood defiant in the crowd!",
			c.Deletions, c.Author),
		fmt.Sprintf(
			"-f was the flag that sealed our fate,\n%s is now but a memory, mate,\nThe reviewers cried \"TOO LATE! TOO LATE!\"\nBut the push had sailed through the gate!",
			c.ShortHash),
		fmt.Sprintf(
			"No branch is safe, no tag survives,\nWhen %s takes the wheel,\nThe history rewrites our lives,\nWith --force and nerves of steel!",
			c.Author),
		fmt.Sprintf(
			"The junior devs all ran and hid,\nWhen %s raised the flag,\n\"I'm cleaning up what others did!\"\n(The Slack channel began to lag.)",
			c.Author),
	}
	return l.pickN(pool, 2)
}

func (l *LyricsGenerator) workChorus(c Commit) string {
	templates := []string{
		fmt.Sprintf(
			"Commit by commit, line by line,\nWe sail the codebase, rain or shine,\n%s's at the helm, and the code is fine,\nAnother day upon the brine!",
			c.Author),
		fmt.Sprintf(
			"Push and pull, pull and push,\nThrough the branches and the bush,\n%s keeps it steady, no need to rush,\n%d files in the morning hush!",
			c.Author, c.FilesChanged),
	}
	return l.pick(templates)
}

func (l *LyricsGenerator) workVerses(c Commit) []string {
	pool := []string{
		fmt.Sprintf(
			"Another morning, another commit,\n%s clocks in, bit by bit,\n%d lines changed — just enough, legit,\nThe daily grind, we must admit.",
			c.Author, c.Insertions+c.Deletions),
		fmt.Sprintf(
			"The stand-up said, \"what did you do?\"\n%s replied, \"a thing or two,\"\n\"%s\" — the ticket, right on cue,\nAnother day, another review.",
			c.Author, truncMsg(c.Message, 30)),
		fmt.Sprintf(
			"No fanfare here, no trumpets sound,\nJust %s on solid ground,\n%d files changed, the work is bound,\nTo keep this ship from running aground.",
			c.Author, c.FilesChanged),
	}
	return l.pickN(pool, 2)
}

func (l *LyricsGenerator) lullabyChorus(c Commit) string {
	templates := []string{
		fmt.Sprintf(
			"Hush now, the README's been updated,\n%s typed so soft and slow,\nA comma here, a typo abated,\nThe quietest commit you'll know...",
			c.Author),
		fmt.Sprintf(
			"Shh... the tests are sleeping sound,\n%s tiptoed through the code,\nNot a single error found,\nOn this gentle, moonlit road...",
			c.Author),
	}
	return l.pick(templates)
}

func (l *LyricsGenerator) lullabyVerses(c Commit) []string {
	pool := []string{
		fmt.Sprintf(
			"A single line was all it took,\n%s fixed the typo in the book,\nThe diff so small you'd need to look,\nTwice, in every nook.",
			c.Author),
		fmt.Sprintf(
			"The changelog whispers, soft and thin,\n\"%s\" — barely worth a grin,\nBut %s checked it in,\nWith a peaceful, sleepy spin.",
			truncMsg(c.Message, 25), c.Author),
		fmt.Sprintf(
			"No alarms, no Slack pings ring,\n%d files, a gentle thing,\n%s did a tiny fling,\nBarely worth mentioning...",
			c.FilesChanged, c.Author),
	}
	return l.pickN(pool, 2)
}

// ---------- bridges ----------

func (l *LyricsGenerator) generateBridge(c Commit, st ShantyType) string {
	// Not every shanty gets a bridge
	if l.rng.Intn(3) == 0 {
		return ""
	}

	bridges := map[ShantyType][]string{
		CelebrationJig: {
			fmt.Sprintf("[Spoken, dramatically]\nAnd lo, on the %s of %s, the build turned GREEN...",
				c.Date.Format("2nd"), c.Date.Format("January 2006")),
			"[Accordion solo — 16 bars of pure joy]",
		},
		MournfulBallad: {
			fmt.Sprintf("[Whispered]\n%s stared at the terminal... the cursor blinked... and blinked...",
				c.Author),
			"[A lone fiddle plays as the CI pipeline times out]",
		},
		EpicSaga: {
			"[Orchestral interlude — the sound of a thousand keyboards]",
			fmt.Sprintf("[Narrator]\nAnd on that day, %s lines were born anew...", fmt.Sprint(c.Insertions)),
		},
		MutinyAnthem: {
			"[GUITAR SOLO — played on a banjo, aggressively]",
			fmt.Sprintf("[Chanted]\n%s! %s! %s!\n(The hash echoes into oblivion)", c.ShortHash, c.ShortHash, c.ShortHash),
		},
		WorkSong: {
			"[Harmonica break — 8 bars, steady tempo]",
			"[The crew hums along, tapping their keyboards in rhythm]",
		},
		LullabyCalm: {
			"[A music box plays softly as the diff loads...]",
			"[Gentle humming... the linter found nothing... all is well...]",
		},
	}

	if options, ok := bridges[st]; ok {
		return l.pick(options)
	}
	return ""
}

// ---------- codas ----------

func (l *LyricsGenerator) generateCoda(c Commit, st ShantyType) string {
	codas := map[ShantyType][]string{
		CelebrationJig: {
			fmt.Sprintf("And that, dear friends, is how %s shipped it! 🎉", c.Author),
			"Three cheers for the merge! Hip hip — HOORAY! 🍺",
			"May your builds be green and your deploys be clean! ⚓",
		},
		MournfulBallad: {
			fmt.Sprintf("Rest in peace, %s... you served us well. 🪦", c.ShortHash),
			"And the revert was merged at dawn... 😢",
			"Some say the conflict markers are still there to this day...",
		},
		EpicSaga: {
			fmt.Sprintf("Thus ends the Great Refactoring of %s! ⚔️", c.Date.Format("2006")),
			"The codebase was never the same again. 🏰",
			"And the tech debt was vanquished... for now. 🐉",
		},
		MutinyAnthem: {
			"git reflog remembers. git reflog ALWAYS remembers. 💀",
			fmt.Sprintf("Legend says %s's force push can still be heard on quiet nights... 👻", c.Author),
			"--force: because sometimes democracy is overrated. ☠️",
		},
		WorkSong: {
			"Same time tomorrow, lads. Same time tomorrow. ⚓",
			fmt.Sprintf("And %s pushed, and it was Tuesday. The end.", c.Author),
			"Another commit, another day closer to retirement. 🚢",
		},
		LullabyCalm: {
			"...and the CI ran green, and all was quiet. 🌙",
			"Good night, sweet codebase. Good night. 💤",
			fmt.Sprintf("Fin. (%s whispers: \"I just fixed a typo.\")", c.Author),
		},
	}

	if options, ok := codas[st]; ok {
		return l.pick(options)
	}
	return "And so the tale is told. ⚓"
}

// ---------- main generators ----------

// GenerateShanty creates a complete sea shanty for a single commit.
func (l *LyricsGenerator) GenerateShanty(commit Commit) Shanty {
	st := ShantyTypeForCommit(commit.Types)
	tempoOptions := tempos[st]

	var chorus string
	var verses []string

	switch st {
	case CelebrationJig:
		chorus = l.celebrationChorus(commit)
		verses = l.celebrationVerses(commit)
	case MournfulBallad:
		chorus = l.mournfulChorus(commit)
		verses = l.mournfulVerses(commit)
	case EpicSaga:
		chorus = l.epicChorus(commit)
		verses = l.epicVerses(commit)
	case MutinyAnthem:
		chorus = l.mutinyChorus(commit)
		verses = l.mutinyVerses(commit)
	case WorkSong:
		chorus = l.workChorus(commit)
		verses = l.workVerses(commit)
	case LullabyCalm:
		chorus = l.lullabyChorus(commit)
		verses = l.lullabyVerses(commit)
	}

	return Shanty{
		Title:    l.generateTitle(commit),
		Type:     st,
		Tempo:    l.pick(tempoOptions),
		Verses:   verses,
		Chorus:   chorus,
		Bridge:   l.generateBridge(commit, st),
		Coda:     l.generateCoda(commit, st),
		BandName: l.bandNameFor(commit.Author),
	}
}

// GenerateEpic creates a multi-commit saga shanty summarising a series of commits.
func (l *LyricsGenerator) GenerateEpic(commits []Commit) Shanty {
	if len(commits) == 0 {
		return Shanty{
			Title:    "The Shanty of Nothing",
			Type:     LullabyCalm,
			Tempo:    "♩ = 0, Silence.",
			Chorus:   "There were no commits, not a one,\nThe repo sat beneath the sun,\nNo code was written, nothing done,\nThe greatest voyage... never begun.",
			Coda:     "git log returned empty. We wept.",
			BandName: l.pick(bandNames),
		}
	}

	// Aggregate stats
	totalFiles, totalIns, totalDel := 0, 0, 0
	authors := map[string]int{}
	var allTypes []CommitType
	for _, c := range commits {
		totalFiles += c.FilesChanged
		totalIns += c.Insertions
		totalDel += c.Deletions
		authors[c.Author]++
		allTypes = append(allTypes, c.Types...)
	}

	topAuthor := ""
	topCount := 0
	for a, n := range authors {
		if n > topCount {
			topAuthor = a
			topCount = n
		}
	}

	st := ShantyTypeForCommit(allTypes)
	if len(commits) >= 10 {
		st = EpicSaga
	}

	// Opening verse
	opening := fmt.Sprintf(
		"Gather 'round, ye coders bold,\nA tale of %d commits be told,\nFrom %s to %s, through nights so cold,\n%s led the crew into the fold!",
		len(commits),
		commits[len(commits)-1].Date.Format("Jan 2"),
		commits[0].Date.Format("Jan 2"),
		topAuthor,
	)

	// Per-commit summary lines (up to 8)
	var summaryLines []string
	limit := len(commits)
	if limit > 8 {
		limit = 8
	}
	for i := 0; i < limit; i++ {
		c := commits[i]
		line := fmt.Sprintf("  ♪ %s — \"%s\" (%+d/-%d)",
			c.ShortHash, truncMsg(c.Message, 35), c.Insertions, c.Deletions)
		summaryLines = append(summaryLines, line)
	}
	if len(commits) > 8 {
		summaryLines = append(summaryLines, fmt.Sprintf("  ♪ ...and %d more commits lost to the mist...", len(commits)-8))
	}
	commitVerse := strings.Join(summaryLines, "\n")

	// Grand chorus
	chorus := fmt.Sprintf(
		"Sail on! Sail on! Through %d files changed!\n%d lines born, %d lines estranged!\nThe codebase grew, the codebase ranged,\nAnd nothing was the same — all rearranged!",
		totalFiles, totalIns, totalDel,
	)

	// Closing verse
	closing := fmt.Sprintf(
		"And so we reached the harbour's end,\n%d commits from foe and friend,\nThe branch is merged, the wounds will mend,\n'Til the next sprint starts again!",
		len(commits),
	)

	// Epic title
	epicTitles := []string{
		fmt.Sprintf("The %d-Commit Odyssey", len(commits)),
		fmt.Sprintf("%s's Grand Voyage: A %d-Commit Saga", topAuthor, len(commits)),
		"An Epic of Diffs and Deploys",
		fmt.Sprintf("The Voyage of the %d Files", totalFiles),
		"From First Push to Final Merge: An Epic",
	}

	return Shanty{
		Title:    l.pick(epicTitles),
		Type:     st,
		Tempo:    l.pick(tempos[st]),
		Verses:   []string{opening, commitVerse, closing},
		Chorus:   chorus,
		Bridge:   "[All instruments join — the sound of a hundred keyboards clacking in unison]",
		Coda:     fmt.Sprintf("Performed by %s — %d commits, %d authors, 0 regrets. ⚓", l.bandNameFor(topAuthor), len(commits), len(authors)),
		BandName: l.bandNameFor(topAuthor),
	}
}

// ---------- helpers ----------

func (l *LyricsGenerator) pickN(pool []string, n int) []string {
	if n >= len(pool) {
		return pool
	}
	// Fisher-Yates partial shuffle
	picked := make([]string, len(pool))
	copy(picked, pool)
	for i := 0; i < n; i++ {
		j := i + l.rng.Intn(len(picked)-i)
		picked[i], picked[j] = picked[j], picked[i]
	}
	return picked[:n]
}

func truncMsg(msg string, maxLen int) string {
	msg = strings.SplitN(msg, "\n", 2)[0] // first line only
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}
