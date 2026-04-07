package main

import (
	"encoding/json"
	"hash/fnv"
	"syscall/js"
)

func seedFromMood(mood string) int64 {
	h := fnv.New64a()
	h.Write([]byte(mood))
	return int64(h.Sum64())
}

func generateAlbum(this js.Value, args []js.Value) interface{} {
	mood := args[0].String()
	trackCount := 8
	if len(args) > 1 {
		trackCount = args[1].Int()
	}

	seed := seedFromMood(mood)
	albumGen := NewAlbumGenerator(seed)
	album := albumGen.Generate(mood)

	trackGen := NewTrackGenerator(seed + 1)
	album.Tracks = trackGen.GenerateTracks(mood, album.Genre, trackCount)

	artGen := NewArtGenerator(album.CoverSeed)
	coverArt := artGen.GenerateArt(album.Genre, mood)

	reviewGen := NewReviewGenerator(seed + 2)
	stats := reviewGen.GenerateStats(album.Title, album.Artist.Name, album.Genre)

	tracks := make([]map[string]interface{}, len(album.Tracks))
	for i, t := range album.Tracks {
		tracks[i] = map[string]interface{}{
			"number":   t.Number,
			"title":    t.Title,
			"duration": t.Duration,
			"bpm":      t.BPM,
			"key":      t.Key,
			"featured": t.Featured,
		}
	}

	reviews := make([]map[string]interface{}, len(stats.Reviews))
	for i, r := range stats.Reviews {
		reviews[i] = map[string]interface{}{
			"reviewer":  r.Reviewer,
			"rating":    r.Rating,
			"stars":     FormatStars(r.Rating),
			"text":      r.Text,
			"playCount": FormatPlayCount(r.PlayCount),
			"helpful":   r.Helpful,
		}
	}

	result := map[string]interface{}{
		"title":  album.Title,
		"artist": album.Artist.Name,
		"bio":    album.Artist.Bio,
		"genre":  string(album.Genre),
		"label":  album.Label,
		"year":   album.Year,
		"mood":   album.Mood,
		"tracks": tracks,
		"coverArt":         coverArt,
		"avgRating":        stats.AvgRating,
		"avgRatingStars":   FormatStars(stats.AvgRating),
		"totalPlays":       FormatPlayCount(stats.TotalPlays),
		"reviews":          reviews,
		"similarTo":        stats.SimilarTo,
		"topPlaylist":      stats.TopPlaylist,
		"monthlyListeners": FormatPlayCount(stats.MonthlyListeners),
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		js.Global().Get("console").Call("error", "JSON marshal error:", err.Error())
		return `{"error":"internal encoding error"}`
	}
	return string(jsonBytes)
}

func main() {
	js.Global().Set("generateAlbum", js.FuncOf(generateAlbum))
	select {}
}
