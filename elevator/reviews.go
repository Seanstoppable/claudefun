package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// Review represents a single user review of an album.
type Review struct {
	Reviewer  string
	Rating    float64
	Text      string
	PlayCount int
	Helpful   int
}

// AlbumStats holds aggregate statistics and reviews for an album.
type AlbumStats struct {
	AvgRating        float64
	TotalPlays       int
	Reviews          []Review
	SimilarTo        []string
	TopPlaylist      string
	MonthlyListeners int
}

// ReviewGenerator produces procedural album reviews and stats.
type ReviewGenerator struct {
	rng *rand.Rand
}

// NewReviewGenerator creates a ReviewGenerator with the given seed.
func NewReviewGenerator(seed int64) *ReviewGenerator {
	return &ReviewGenerator{rng: rand.New(rand.NewSource(seed))}
}

var reviewerNames = []string{
	"lobby_listener_42",
	"elevator_enthusiast",
	"smooth_jazz_steve",
	"background_patricia",
	"hold_music_connoisseur",
	"cubicle_carl",
	"waiting_room_warrior",
	"muzak_maven_99",
	"the_real_kenny_g",
	"ambient_alice",
	"beige_beats_fan",
	"soft_rock_sandra",
	"tepid_tony",
	"mild_mike_2003",
}

// templates keyed by star tier: 5, 3-4, 1-2
var fiveStarTemplates = []string{
	"I put this on in the elevator at work and three people complimented the ambiance. Career highlight.",
	"Finally, music that matches the exact energy of standing in line at the DMV. Transcendent.",
	"This album got me through a 47-minute hold with my insurance company. I owe {artist} my sanity.",
	"I fell asleep to this and woke up a better person. 10/10, would gently doze again.",
}

var midStarTemplates = []string{
	"Perfectly adequate. Would listen to again if it happened to be playing somewhere I was.",
	"Track 4 made me feel something, which isn't really what I come to elevator music for, but I'll allow it.",
	"Good background music for existing. Not great, not terrible. Like a C+ in music form.",
	"My therapist recommended this album. I'm not sure if that's an endorsement.",
}

var lowStarTemplates = []string{
	"Too engaging. I caught myself nodding along during a meeting. Unprofessional.",
	"This has a saxophone solo that demands attention. That defeats the entire purpose.",
	"I could identify individual notes. Elevator music should be more... vague.",
	"Made my coworker ask 'what are we listening to?' — the cardinal sin of background music.",
}

var similarArtists = []string{
	"Muzak & The Holding Patterns",
	"The Lobby Experience",
	"Smooth Jazz for Nobody in Particular",
	"Gary Elevator",
	"The Beige Album Collective",
	"DJ Hold Please",
	"The Saxophones of Ambiguity",
	"Patricia Smooth",
	"The Waiting Room All-Stars",
	"Background Noise Orchestra",
}

var playlistNames = []string{
	"Songs to Wait By",
	"Lobby Essentials",
	"Corporate Chill",
	"The Inoffensive Mix",
	"Background Music for Background Living",
	"Elevator Hits Vol. 47",
	"Smooth Operations",
	"Dental Waiting Room Vibes",
	"Music for Pretending to Read Magazines",
	"The 'Is This Even Playing?' Playlist",
}

// GenerateStats produces a full AlbumStats for the given album.
func (rg *ReviewGenerator) GenerateStats(albumTitle string, artistName string, genre Genre) AlbumStats {
	numReviews := 3 + rg.rng.Intn(3) // 3-5 reviews

	// Pick unique reviewers
	perm := rg.rng.Perm(len(reviewerNames))
	reviews := make([]Review, numReviews)
	var ratingSum float64

	for i := 0; i < numReviews; i++ {
		reviewer := reviewerNames[perm[i%len(reviewerNames)]]
		rating, text := rg.pickReview(artistName)
		playCount := 1000 + rg.rng.Intn(499001) // 1K-500K
		helpful := rg.rng.Intn(120)

		reviews[i] = Review{
			Reviewer:  reviewer,
			Rating:    rating,
			Text:      text,
			PlayCount: playCount,
			Helpful:   helpful,
		}
		ratingSum += rating
	}

	avgRating := ratingSum / float64(numReviews)
	// Round to nearest half star
	avgRating = float64(int(avgRating*2+0.5)) / 2.0

	totalPlays := 0
	for _, r := range reviews {
		totalPlays += r.PlayCount
	}

	// Pick 2-4 similar artists
	numSimilar := 2 + rg.rng.Intn(3)
	similarPerm := rg.rng.Perm(len(similarArtists))
	similar := make([]string, numSimilar)
	for i := 0; i < numSimilar; i++ {
		similar[i] = similarArtists[similarPerm[i]]
	}

	playlist := playlistNames[rg.rng.Intn(len(playlistNames))]
	monthlyListeners := 500 + rg.rng.Intn(49501) // 500-50K

	return AlbumStats{
		AvgRating:        avgRating,
		TotalPlays:       totalPlays,
		Reviews:          reviews,
		SimilarTo:        similar,
		TopPlaylist:      playlist,
		MonthlyListeners: monthlyListeners,
	}
}

func (rg *ReviewGenerator) pickReview(artistName string) (float64, string) {
	tier := rg.rng.Intn(10)
	var rating float64
	var text string

	switch {
	case tier < 4: // 40% chance 5 stars
		rating = 5.0
		text = fiveStarTemplates[rg.rng.Intn(len(fiveStarTemplates))]
	case tier < 8: // 40% chance 3-4 stars
		rating = 3.0 + float64(rg.rng.Intn(3))*0.5 // 3.0, 3.5, 4.0
		text = midStarTemplates[rg.rng.Intn(len(midStarTemplates))]
	default: // 20% chance 1-2 stars
		rating = 1.0 + float64(rg.rng.Intn(3))*0.5 // 1.0, 1.5, 2.0
		text = lowStarTemplates[rg.rng.Intn(len(lowStarTemplates))]
	}

	text = strings.ReplaceAll(text, "{artist}", artistName)
	return rating, text
}

// FormatStars returns a star rating string like "★★★★☆" or "★★★½☆".
func FormatStars(rating float64) string {
	full := int(rating)
	half := (rating - float64(full)) >= 0.5
	empty := 5 - full
	if half {
		empty--
	}

	var b strings.Builder
	for i := 0; i < full; i++ {
		b.WriteString("★")
	}
	if half {
		b.WriteString("½")
	}
	for i := 0; i < empty; i++ {
		b.WriteString("☆")
	}
	return b.String()
}

// FormatPlayCount formats a number with K/M suffix.
func FormatPlayCount(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}
