package mood

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ssmith/mood-octopus/octopus"
)

const maxEntries = 100

// MoodEntry records a single mood observation.
type MoodEntry struct {
	Emotion   octopus.Emotion `json:"emotion"`
	Input     string          `json:"input"`
	Timestamp time.Time       `json:"timestamp"`
}

// History tracks mood entries and persists them to disk.
type History struct {
	entries []MoodEntry
	path    string
}

// NewHistory loads (or initialises) the mood history from ~/.mood-octopus/history.json.
func NewHistory() (*History, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".mood-octopus")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	h := &History{path: filepath.Join(dir, "history.json")}

	data, err := os.ReadFile(h.path)
	if err == nil {
		// Ignore malformed files; start fresh.
		_ = json.Unmarshal(data, &h.entries)
	}

	return h, nil
}

// Record appends a mood entry and persists the history to disk.
func (h *History) Record(emotion octopus.Emotion, input string) error {
	h.entries = append(h.entries, MoodEntry{
		Emotion:   emotion,
		Input:     input,
		Timestamp: time.Now(),
	})
	h.trim()
	return h.save()
}

// LastMood returns the most recent entry, or nil if there is no history.
func (h *History) LastMood() *MoodEntry {
	if len(h.entries) == 0 {
		return nil
	}
	e := h.entries[len(h.entries)-1]
	return &e
}

// Recent returns the last n entries (or fewer if the history is shorter).
func (h *History) Recent(n int) []MoodEntry {
	if n <= 0 {
		return nil
	}
	if n > len(h.entries) {
		n = len(h.entries)
	}
	out := make([]MoodEntry, n)
	copy(out, h.entries[len(h.entries)-n:])
	return out
}

// StartupGreeting returns a contextual greeting based on the last recorded mood.
func (h *History) StartupGreeting() string {
	last := h.LastMood()
	if last == nil {
		return "Hello! I'm your new octopus friend! 🐙 Type anything and watch me react!"
	}

	switch last.Emotion {
	case octopus.Joy:
		return "Welcome back! Last time you were happy — let's keep that going! 🌟"
	case octopus.Sadness:
		return "Hey, you were feeling down last time. I hope things are better now 💙"
	case octopus.Anger:
		return "You were pretty fired up last time. Deep breaths... 🌊"
	case octopus.Fear:
		return "Last time was a bit scary, huh? Don't worry, I'm here! 🤗"
	case octopus.Curiosity:
		return "Welcome back, curious one! What shall we explore today? 🔍"
	case octopus.Sleepy:
		return "Ah, you were sleepy last time. Feeling rested? ☕"
	case octopus.Silly:
		return "We had fun last time! Ready for more shenanigans? 🤪"
	case octopus.Love:
		return "Aww, last time was full of love. My hearts are still warm 💕"
	default:
		return "Welcome back! Let's see what mood we're in today 🐙"
	}
}

// MoodSparkline returns a string of emoji representing the last `width` moods.
func (h *History) MoodSparkline(width int) string {
	recent := h.Recent(width)
	if len(recent) == 0 {
		return ""
	}
	var b strings.Builder
	for _, e := range recent {
		b.WriteString(e.Emotion.Info().Emoji)
	}
	return b.String()
}

// trim keeps the history at most maxEntries long.
func (h *History) trim() {
	if len(h.entries) > maxEntries {
		h.entries = h.entries[len(h.entries)-maxEntries:]
	}
}

func (h *History) save() error {
	data, err := json.Marshal(h.entries)
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}
