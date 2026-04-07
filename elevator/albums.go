package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// Genre represents a sub-genre of elevator music.
type Genre string

const (
	AmbientLobby      Genre = "Ambient Lobby"
	CorporateZen      Genre = "Corporate Zen"
	WaitingRoomJazz   Genre = "Waiting Room Jazz"
	HoldMusicDeluxe   Genre = "Hold Music Deluxe"
	DentistOfficeCore Genre = "Dentist Office Core"
	ParkingGarageWave Genre = "Parking Garage Wave"
	SoftRockPurgatory Genre = "Soft Rock Purgatory"
	AcousticBeige     Genre = "Acoustic Beige"
	LoFiElevator      Genre = "Lo-Fi Elevator"
	SmoothBureaucracy Genre = "Smooth Bureaucracy"
	MidtempoMalaise   Genre = "Midtempo Malaise"
	SuburbanAmbient   Genre = "Suburban Ambient"
)

// AllGenres returns every available elevator music genre.
func AllGenres() []Genre {
	return []Genre{
		AmbientLobby, CorporateZen, WaitingRoomJazz, HoldMusicDeluxe,
		DentistOfficeCore, ParkingGarageWave, SoftRockPurgatory, AcousticBeige,
		LoFiElevator, SmoothBureaucracy, MidtempoMalaise, SuburbanAmbient,
	}
}

// Artist represents the fictional performer of an elevator music album.
type Artist struct {
	Name string
	Bio  string
}

// Album is a procedurally generated elevator music album.
type Album struct {
	Title     string
	Artist    Artist
	Genre     Genre
	Label     string
	Year      int
	Tracks    []Track
	CoverSeed int64
	Mood      string
}

// AlbumGenerator produces fictional elevator music albums.
type AlbumGenerator struct {
	rng *rand.Rand
}

// NewAlbumGenerator creates a generator seeded with the given value.
func NewAlbumGenerator(seed int64) *AlbumGenerator {
	return &AlbumGenerator{rng: rand.New(rand.NewSource(seed))}
}

// ---------------------------------------------------------------------------
// Name pools
// ---------------------------------------------------------------------------

var artistFirstNames = []string{
	"Gary", "Linda", "Soft", "Gentle", "Ambient", "Smooth", "Quiet", "Mild",
	"Tepid", "Room Temperature", "Beige", "Kenneth", "Patricia", "The Honorable", "DJ",
}

var artistLastNames = []string{
	"Feelings", "Elevator", "Saxworth", "Smoothington", "Blandford", "Mellow",
	"Neutralson", "Waiting", "Background", "Holdmusic", "Lobbyton", "Ambiance",
}

var bandNames = []string{
	"The Comfortable Silences",
	"Muzak & The Holding Patterns",
	"The Inoffensive",
	"Passive Aggressive Jazz Quartet",
	"The Background Noise Collective",
	"Elevator & The Shafts",
	"The Lukewarm Takes",
	"Lobby Lizards",
}

var artistBios = []string{
	"Has been playing in elevators since before elevators had music.",
	"Voted 'Most Likely to Be Ignored' three years running.",
	"Their music has been described as 'the sound of waiting, but polished.'",
	"Pioneered the 'is this even playing?' movement in 2003.",
	"Former dentist who found their true calling in background noise.",
	"Once held a concert where nobody noticed it was happening.",
	"Believes every room deserves a soundtrack no one asked for.",
	"Has a PhD in Ambient Nothingness from the University of Beige.",
	"Their debut album was accidentally used as hold music for the IRS.",
	"Known for performing exclusively in buildings with at least three floors.",
}

// ---------------------------------------------------------------------------
// Album titles, labels, adjectives
// ---------------------------------------------------------------------------

var albumTitleTemplates = []string{
	"%s Feelings Vol. %d",
	"Songs for %s Elevators",
	"The %s Side of %s",
	"%s: A Lobby Suite",
	"Music for People Who Are %s",
	"Waiting (%s Edition)",
	"%s in the Key of Beige",
	"The Essential %s Collection",
	"Hold Please: %s Remixes",
	"%s & Other Mild Sensations",
}

var titleAdjectives = []string{
	"Softer", "Smoother", "Quieter", "Milder", "Gentler",
	"Beiger", "Lukewarm", "Pleasantly Dull", "Inoffensive", "Tepid",
}

var recordLabels = []string{
	"Muzak Records",
	"Smooth Operations Inc.",
	"Lobby International",
	"Hold Music Heritage Foundation",
	"Beige Note Records",
	"ElevatorFi",
	"The Waiting Room Label",
	"Inoffensive Records",
	"Background Presence Music",
	"Tepid Grooves Publishing",
}

// ---------------------------------------------------------------------------
// Mood → Genre mapping
// ---------------------------------------------------------------------------

var moodGenreMap = map[string][]Genre{
	"calm":       {AmbientLobby, SuburbanAmbient},
	"peaceful":   {AmbientLobby, SuburbanAmbient},
	"relaxed":    {AmbientLobby, SuburbanAmbient},
	"corporate":  {CorporateZen, SmoothBureaucracy},
	"business":   {CorporateZen, SmoothBureaucracy},
	"work":       {CorporateZen, SmoothBureaucracy},
	"jazz":       {WaitingRoomJazz, LoFiElevator},
	"smooth":     {WaitingRoomJazz, LoFiElevator},
	"mellow":     {WaitingRoomJazz, LoFiElevator},
	"sad":        {MidtempoMalaise, AcousticBeige},
	"melancholy": {MidtempoMalaise, AcousticBeige},
	"blue":       {MidtempoMalaise, AcousticBeige},
	"happy":      {HoldMusicDeluxe, DentistOfficeCore},
	"cheerful":   {HoldMusicDeluxe, DentistOfficeCore},
	"anxious":    {ParkingGarageWave},
	"nervous":    {ParkingGarageWave},
	"bored":      {SoftRockPurgatory},
}

func (ag *AlbumGenerator) pickGenre(mood string) Genre {
	lower := strings.ToLower(mood)
	for keyword, genres := range moodGenreMap {
		if strings.Contains(lower, keyword) {
			return genres[ag.rng.Intn(len(genres))]
		}
	}
	all := AllGenres()
	return all[ag.rng.Intn(len(all))]
}

// ---------------------------------------------------------------------------
// Generation helpers
// ---------------------------------------------------------------------------

func (ag *AlbumGenerator) generateArtist() Artist {
	var name string
	if ag.rng.Intn(3) == 0 {
		name = bandNames[ag.rng.Intn(len(bandNames))]
	} else {
		first := artistFirstNames[ag.rng.Intn(len(artistFirstNames))]
		last := artistLastNames[ag.rng.Intn(len(artistLastNames))]
		name = first + " " + last
	}
	bio := artistBios[ag.rng.Intn(len(artistBios))]
	return Artist{Name: name, Bio: bio}
}

func (ag *AlbumGenerator) generateTitle(mood string) string {
	idx := ag.rng.Intn(len(albumTitleTemplates))
	tmpl := albumTitleTemplates[idx]
	titleMood := strings.Title(mood) //nolint:staticcheck

	switch idx {
	case 0: // "%s Feelings Vol. %d"
		return fmt.Sprintf(tmpl, titleMood, ag.rng.Intn(9)+1)
	case 2: // "The %s Side of %s"
		adj := titleAdjectives[ag.rng.Intn(len(titleAdjectives))]
		return fmt.Sprintf(tmpl, adj, titleMood)
	default:
		return fmt.Sprintf(tmpl, titleMood)
	}
}

// Generate produces a complete Album for the given mood.
func (ag *AlbumGenerator) Generate(mood string) Album {
	genre := ag.pickGenre(mood)
	artist := ag.generateArtist()
	title := ag.generateTitle(mood)
	label := recordLabels[ag.rng.Intn(len(recordLabels))]
	year := 1975 + ag.rng.Intn(50) // 1975–2024
	trackCount := 8 + ag.rng.Intn(6) // 8-13 tracks
	coverSeed := ag.rng.Int63()

	tg := NewTrackGenerator(ag.rng.Int63())
	tracks := tg.GenerateTracks(mood, genre, trackCount)

	return Album{
		Title:     title,
		Artist:    artist,
		Genre:     genre,
		Label:     label,
		Year:      year,
		Tracks:    tracks,
		CoverSeed: coverSeed,
		Mood:      mood,
	}
}
