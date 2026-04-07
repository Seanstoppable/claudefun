package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ── colour palette ──────────────────────────────────────────────────────────

var (
	purple    = lipgloss.Color("#9B59B6")
	teal      = lipgloss.Color("#1ABC9C")
	gold      = lipgloss.Color("#F1C40F")
	secondary = lipgloss.Color("#95A5A6")
	dimColor  = lipgloss.Color("#7F8C8D")
	white     = lipgloss.Color("#ECF0F1")
)

// ── reusable styles ─────────────────────────────────────────────────────────

var (
	headerBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(0, 4).
			Align(lipgloss.Center).
			Foreground(teal)

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(white)
	artistStyle = lipgloss.NewStyle().Foreground(secondary)
	genreStyle  = lipgloss.NewStyle().Foreground(teal)
	labelStyle  = lipgloss.NewStyle().Foreground(dimColor)
	starsStyle  = lipgloss.NewStyle().Foreground(gold)

	trackNumStyle  = lipgloss.NewStyle().Foreground(dimColor).Width(4).Align(lipgloss.Right)
	trackNameStyle = lipgloss.NewStyle().Foreground(white)
	trackMetaStyle = lipgloss.NewStyle().Foreground(dimColor)
	featStyle      = lipgloss.NewStyle().Italic(true).Foreground(secondary)

	reviewerStyle = lipgloss.NewStyle().Foreground(teal).Bold(true)
	reviewText    = lipgloss.NewStyle().Foreground(secondary)
	helpfulStyle  = lipgloss.NewStyle().Foreground(dimColor)

	sectionIcon = lipgloss.NewStyle().Foreground(teal)
	divider     = lipgloss.NewStyle().Foreground(dimColor).Render(strings.Repeat("─", 55))

	artStyle = lipgloss.NewStyle().Foreground(purple)
)

func main() {
	mood := flag.String("mood", "", "mood / feeling for the album")
	tracks := flag.Int("tracks", 8, "number of tracks")
	discography := flag.Bool("discography", false, "generate 3 albums as a discography")
	seed := flag.Int64("seed", 0, "seed for reproducible output")
	flag.Parse()

	// Accept first positional arg as mood if -mood wasn't set.
	if *mood == "" && flag.NArg() > 0 {
		*mood = strings.Join(flag.Args(), " ")
	}

	if *mood == "" {
		fmt.Println("\n  🎵 What mood are you in? Describe it, and I'll compose the perfect elevator music.")
		fmt.Println("  Usage: elevator-music <mood> [-tracks N] [-discography] [-seed N]")
		fmt.Println("  Example: elevator-music \"melancholy\"")
		fmt.Println()
		os.Exit(0)
	}

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}

	// Print header.
	printHeader()

	if *discography {
		rng := rand.New(rand.NewSource(*seed))
		// Generate a shared artist.
		ag := NewAlbumGenerator(rng.Int63())
		firstAlbum := ag.Generate(*mood)
		sharedArtist := firstAlbum.Artist

		for i := 0; i < 3; i++ {
			albumGen := NewAlbumGenerator(*seed + int64(i)*1337)
			album := albumGen.Generate(*mood)
			album.Artist = sharedArtist
			album.Tracks = NewTrackGenerator(*seed + int64(i)*42).GenerateTracks(*mood, album.Genre, *tracks)
			renderAlbum(album)
			if i < 2 {
				fmt.Println()
				fmt.Println(lipgloss.NewStyle().Foreground(purple).Render(
					"  " + strings.Repeat("═", 55)))
				fmt.Println()
			}
		}
	} else {
		ag := NewAlbumGenerator(*seed)
		album := ag.Generate(*mood)
		album.Tracks = NewTrackGenerator(*seed + 42).GenerateTracks(*mood, album.Genre, *tracks)
		renderAlbum(album)
	}

	fmt.Println()
}

// ── header ──────────────────────────────────────────────────────────────────

func printHeader() {
	title := lipgloss.NewStyle().Bold(true).Foreground(purple).Render("🎵  ELEVATOR MUSIC COMPOSER  🎵")
	subtitle := lipgloss.NewStyle().Italic(true).Foreground(secondary).Render(`"Now playing: your feelings"`)

	box := headerBox.Render(title + "\n" + subtitle)
	fmt.Println()
	// Indent the box.
	for _, line := range strings.Split(box, "\n") {
		fmt.Println("  " + line)
	}
	fmt.Println()
}

// ── album rendering ─────────────────────────────────────────────────────────

func renderAlbum(album Album) {
	artGen := NewArtGenerator(album.CoverSeed)
	art := artGen.GenerateArt(album.Genre, album.Mood)
	renderAlbumInfo(album, art)
	renderTrackListing(album.Tracks)

	rg := NewReviewGenerator(album.CoverSeed + 7)
	stats := rg.GenerateStats(album.Title, album.Artist.Name, album.Genre)
	renderReviews(stats)
	renderFooter(stats)
}

func renderAlbumInfo(album Album, art string) {
	artLines := strings.Split(strings.TrimRight(art, "\n"), "\n")
	styledArt := make([]string, len(artLines))
	for i, line := range artLines {
		styledArt[i] = artStyle.Render(line)
	}

	// Build info lines to sit beside the art.
	info := []string{
		titleStyle.Render(album.Title),
		artistStyle.Render("by " + album.Artist.Name),
		"",
		genreStyle.Render("Genre: " + string(album.Genre)),
		labelStyle.Render(fmt.Sprintf("Label: %s (%d)", album.Label, album.Year)),
	}

	// We add the rating line after generating stats for the header.
	rg := NewReviewGenerator(album.CoverSeed + 7)
	stats := rg.GenerateStats(album.Title, album.Artist.Name, album.Genre)
	ratingLine := starsStyle.Render(FormatStars(stats.AvgRating)) + " " +
		lipgloss.NewStyle().Foreground(dimColor).Render(
			fmt.Sprintf("%.1f", stats.AvgRating)) + " · " +
		lipgloss.NewStyle().Foreground(dimColor).Render(
			FormatPlayCount(stats.TotalPlays)+" plays")
	info = append(info, ratingLine)

	// Pad info to match art height.
	for len(info) < len(styledArt) {
		info = append(info, "")
	}
	for len(styledArt) < len(info) {
		styledArt = append(styledArt, strings.Repeat(" ", 20))
	}

	// Print side-by-side.
	for i := 0; i < len(styledArt); i++ {
		infoStr := ""
		if i < len(info) {
			infoStr = info[i]
		}
		fmt.Printf("  %s    %s\n", styledArt[i], infoStr)
	}
	fmt.Println()
}

func renderTrackListing(tracks []Track) {
	fmt.Println("  " + sectionIcon.Render("♫") + " " + titleStyle.Render("Track Listing"))
	fmt.Println("  " + divider)

	totalSeconds := 0
	for _, t := range tracks {
		// Parse duration for total.
		parts := strings.Split(t.Duration, ":")
		if len(parts) == 2 {
			m, _ := strconv.Atoi(parts[0])
			s, _ := strconv.Atoi(parts[1])
			totalSeconds += m*60 + s
		}

		num := trackNumStyle.Render(fmt.Sprintf("%d.", t.Number))
		name := trackNameStyle.Render(t.Title)
		dur := trackMetaStyle.Render(t.Duration)
		bpm := trackMetaStyle.Render(fmt.Sprintf("♩=%d", t.BPM))
		key := trackMetaStyle.Render(t.Key)

		fmt.Printf("  %s %s  %s  %s  %s\n", num, name, dur, bpm, key)

		if t.Featured != "" {
			pad := strings.Repeat(" ", 6)
			fmt.Printf("  %s%s\n", pad, featStyle.Render(t.Featured))
		}
	}

	fmt.Println("  " + divider)
	totalMin := totalSeconds / 60
	totalSec := totalSeconds % 60
	fmt.Printf("  Total: %s\n\n", trackMetaStyle.Render(fmt.Sprintf("%d:%02d", totalMin, totalSec)))
}

func renderReviews(stats AlbumStats) {
	fmt.Println("  " + sectionIcon.Render("⭐") + " " + titleStyle.Render("Reviews"))
	fmt.Println("  " + divider)

	maxReviews := 3
	if len(stats.Reviews) < maxReviews {
		maxReviews = len(stats.Reviews)
	}

	for i := 0; i < maxReviews; i++ {
		r := stats.Reviews[i]
		name := reviewerStyle.Render(r.Reviewer)
		stars := starsStyle.Render(FormatStars(r.Rating))
		fmt.Printf("  %s  %s\n", name, stars)

		// Word-wrap review text to ~52 chars with surrounding quotes.
		wrapped := wordWrap(r.Text, 52)
		for j, line := range wrapped {
			prefix := " "
			suffix := ""
			if j == 0 {
				prefix = `"`
			}
			if j == len(wrapped)-1 {
				suffix = `"`
			}
			fmt.Printf("  %s\n", reviewText.Render(prefix+line+suffix))
		}
		fmt.Printf("   %s\n", helpfulStyle.Render(fmt.Sprintf("↳ %d people found this helpful", r.Helpful)))
		fmt.Println()
	}

	fmt.Println("  " + divider)
}

func renderFooter(stats AlbumStats) {
	// Similar artists.
	similar := strings.Join(stats.SimilarTo, ", ")
	fmt.Printf("  %s Similar Artists: %s\n",
		sectionIcon.Render("📻"),
		lipgloss.NewStyle().Foreground(secondary).Render(similar))

	// Featured playlist.
	fmt.Printf("  %s Featured on: %s\n",
		sectionIcon.Render("📋"),
		lipgloss.NewStyle().Foreground(secondary).Render(`"`+stats.TopPlaylist+`"`))

	// Monthly listeners.
	fmt.Printf("  %s %s monthly listeners\n",
		sectionIcon.Render("👥"),
		lipgloss.NewStyle().Foreground(secondary).Render(FormatPlayCount(stats.MonthlyListeners)))

	fmt.Println("  " + divider)
}

// ── helpers ─────────────────────────────────────────────────────────────────

func wordWrap(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	current := words[0]
	for _, w := range words[1:] {
		if len(current)+1+len(w) > width {
			lines = append(lines, current)
			current = " " + w
		} else {
			current += " " + w
		}
	}
	lines = append(lines, current)
	return lines
}
