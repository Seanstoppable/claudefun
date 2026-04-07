package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// Track represents a single track on an elevator music album.
type Track struct {
	Number   int
	Title    string
	Duration string
	BPM      int
	Key      string
	Featured string
}

// TrackGenerator produces procedural track listings.
type TrackGenerator struct {
	rng *rand.Rand
}

// NewTrackGenerator creates a track generator seeded with the given value.
func NewTrackGenerator(seed int64) *TrackGenerator {
	return &TrackGenerator{rng: rand.New(rand.NewSource(seed))}
}

// ---------------------------------------------------------------------------
// Track title templates — some use the mood, most are standalone
// ---------------------------------------------------------------------------

var staticTrackTitles = []string{
	"Waiting for the Elevator to Arrive (But Not Urgently)",
	"Mild Contentment in a Fluorescent-Lit Room",
	"The Feeling of Being on Hold But It's Fine",
	"Ambient Awareness of a Water Cooler",
	"Soft Realization That It's Only Tuesday",
	"Gentle Acceptance of Moderate Traffic",
	"The Sound of Someone Else's Printer",
	"Contemplating the Snack Machine (Interlude)",
	"Existing in a Temperature-Controlled Environment",
	"Background Process",
	"Lobby of the Mind",
	"Smooth Transition to Nothing in Particular",
	"The Comfort of Predictable Weather",
	"Untroubled by Specific Events",
	"Saxophone Solo for No Occasion",
	"Piano Thoughts (Lightly Seasoned)",
	"Hold Music for the Soul",
	"Interlude: Brief Pause for Reflection (or Not)",
	"Quiet Enthusiasm for Carpet Samples",
	"The Hum of Central Air",
	"Nodding Politely at a Stranger",
	"Third Floor: Miscellaneous Feelings",
	"Soft Jazz for Quarterly Reviews",
	"Elevator Doors Closing (Take 7)",
	"Patiently Awaiting the Microwave",
	"The Gentle Art of Small Talk",
	"Afternoon Slump in D Minor",
	"Ambient Spreadsheet Energy",
	"Saxophone in the Stairwell",
	"Politely Declining a Meeting Invite",
	"The Emotional Range of a Thermostat",
	"Vibes (Room Temperature)",
	"Soothing Sounds of a Parking Structure",
	"Contemplative Muzak for the Modern Lobby",
	"The Space Between Floors",
	"Calm Resignation at the DMV",
	"Acoustic Wallpaper",
	"Light Drizzle on a Skylight (Extended Mix)",
	"The Inoffensive Hour",
	"Beige Sunset",
}

var moodTrackTemplates = []string{
	"%s But Make It Elevator",
	"Muzak for %s Moments",
	"%s in the Waiting Room",
	"A %s Feeling (Lobby Remix)",
	"The %s Elevator Suite",
	"%s: Floor 7",
	"Gently %s (Instrumental)",
	"%s Tones for Conference Room B",
	"On Hold with %s Feelings",
	"The %s Side of the Mezzanine",
}

// ---------------------------------------------------------------------------
// Musical keys & featured artists
// ---------------------------------------------------------------------------

var musicalKeys = []string{
	"C Major", "C Minor", "C♯ Minor",
	"D Major", "D Minor",
	"E♭ Major", "E Major", "E Minor",
	"F Major", "F Minor", "F♯ Minor",
	"G Major", "G Minor",
	"A♭ Major", "A Major", "A Minor",
	"B♭ Major", "B♭ Minor", "B Major", "B Minor",
}

var featuredArtists = []string{
	"feat. The Lobby Strings",
	"feat. DJ Smooth Transition",
	"feat. The Saxophone of Ambiguity",
	"feat. Kenneth (from accounting)",
	"feat. Gentle Gary on Keys",
	"feat. The Waiting Room Choir",
	"feat. Anonymous Flautist",
	"feat. Patricia & The Hold Tones",
}

// ---------------------------------------------------------------------------
// Genre-influenced bonus titles
// ---------------------------------------------------------------------------

var genreBonusTitles = map[Genre][]string{
	AmbientLobby:      {"Lobby Drift", "Atrium Echo", "Revolving Door Reverie"},
	CorporateZen:      {"Synergy in E♭", "Team-Building Ballad", "Fiscal Quarter Lullaby"},
	WaitingRoomJazz:   {"Magazine Table Blues", "Old Highlights Reel", "Waiting Room Waltz"},
	HoldMusicDeluxe:   {"Please Continue to Hold", "Your Call Is Important", "Estimated Wait: Forever"},
	DentistOfficeCore:  {"Open Wide (Softly)", "Novocaine Dreams", "Drill-Free Zone"},
	ParkingGarageWave:  {"Level P3 Ambient", "Concrete Echo", "Ticket Lost, Vibes Found"},
	SoftRockPurgatory:  {"Hotel Breakfast Anthem", "Endless Checkout Lane", "FM Dial Purgatory"},
	AcousticBeige:      {"Off-White Serenade", "Khaki Afternoon", "Eggshell Emotions"},
	LoFiElevator:       {"Crackle & Ascend", "Dusty Lobby Tape", "VHS Lobby Footage"},
	SmoothBureaucracy:  {"Form B-7 (Smooth Mix)", "Permit Pending Groove", "Filing Cabinet Funk"},
	MidtempoMalaise:    {"Not Sad, Just Indoors", "Moderate Disappointment", "Grey Sky Groove"},
	SuburbanAmbient:    {"Cul-de-Sac Calm", "Sprinkler System ASMR", "HOA-Approved Vibes"},
}

// ---------------------------------------------------------------------------
// Generation
// ---------------------------------------------------------------------------

// GenerateTracks produces a track listing for the given mood and genre.
func (tg *TrackGenerator) GenerateTracks(mood string, genre Genre, count int) []Track {
	titleMood := strings.Title(mood) //nolint:staticcheck

	// Build a pool of candidate titles.
	pool := make([]string, 0, len(staticTrackTitles)+len(moodTrackTemplates)+3)
	pool = append(pool, staticTrackTitles...)

	for _, tmpl := range moodTrackTemplates {
		pool = append(pool, fmt.Sprintf(tmpl, titleMood))
	}

	if bonus, ok := genreBonusTitles[genre]; ok {
		pool = append(pool, bonus...)
	}

	// Shuffle and pick unique titles.
	tg.rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	tracks := make([]Track, count)
	for i := 0; i < count; i++ {
		title := pool[i%len(pool)]

		// Duration: 2:30 – 5:00  (150–300 seconds)
		totalSec := 150 + tg.rng.Intn(151)
		minutes := totalSec / 60
		seconds := totalSec % 60
		duration := fmt.Sprintf("%d:%02d", minutes, seconds)

		bpm := 70 + tg.rng.Intn(51) // 70–120
		key := musicalKeys[tg.rng.Intn(len(musicalKeys))]

		var featured string
		if tg.rng.Intn(5) == 0 { // 20% chance
			featured = featuredArtists[tg.rng.Intn(len(featuredArtists))]
		}

		tracks[i] = Track{
			Number:   i + 1,
			Title:    title,
			Duration: duration,
			BPM:      bpm,
			Key:      key,
			Featured: featured,
		}
	}

	return tracks
}
