package main

// Emotion represents one of the octopus's emotional states.
type Emotion int

const (
	Joy Emotion = iota
	Sadness
	Anger
	Fear
	Curiosity
	Sleepy
	Silly
	Love
)

// EmotionInfo holds display metadata and trigger keywords for an emotion.
type EmotionInfo struct {
	Name     string
	Emoji    string
	Color    string
	Keywords []string
}

var emotionRegistry = map[Emotion]EmotionInfo{
	Joy: {
		Name:  "Joy",
		Emoji: "😊",
		Color: "#FFD700",
		Keywords: []string{
			"happy", "great", "awesome", "love", "wonderful",
			"amazing", "fantastic", "sunshine", "celebrate", "excited",
			"yay", "brilliant", "delighted", "cheerful", "ecstatic",
		},
	},
	Sadness: {
		Name:  "Sadness",
		Emoji: "😢",
		Color: "#4682B4",
		Keywords: []string{
			"sad", "cry", "miss", "lonely", "rain",
			"depressed", "gloomy", "heartbreak", "sorry", "tears",
			"melancholy", "blue", "down", "grief", "sorrow",
		},
	},
	Anger: {
		Name:  "Anger",
		Emoji: "😡",
		Color: "#FF4500",
		Keywords: []string{
			"angry", "hate", "furious", "rage", "ugh",
			"mad", "annoyed", "frustrated", "livid", "irritated",
			"outraged", "seething", "hostile", "bitter", "grr",
		},
	},
	Fear: {
		Name:  "Fear",
		Emoji: "😱",
		Color: "#9B59B6",
		Keywords: []string{
			"scared", "afraid", "nervous", "yikes", "help",
			"terrified", "anxious", "dread", "panic", "worry",
			"horror", "creepy", "spooky", "alarming", "eek",
		},
	},
	Curiosity: {
		Name:  "Curiosity",
		Emoji: "🤔",
		Color: "#1ABC9C",
		Keywords: []string{
			"why", "how", "wonder", "interesting", "hmm",
			"curious", "question", "puzzle", "mystery", "explore",
			"what", "think", "ponder", "fascinating", "intriguing",
		},
	},
	Sleepy: {
		Name:  "Sleepy",
		Emoji: "😴",
		Color: "#7B68EE",
		Keywords: []string{
			"tired", "exhausted", "sleep", "yawn", "zzz",
			"drowsy", "nap", "bed", "snooze", "fatigue",
			"weary", "rest", "dreamy", "lethargic", "nodding",
		},
	},
	Silly: {
		Name:  "Silly",
		Emoji: "🤪",
		Color: "#FF69B4",
		Keywords: []string{
			"lol", "haha", "bonkers", "weird", "bruh",
			"goofy", "ridiculous", "absurd", "wacky", "banana",
			"nonsense", "clown", "giggles", "derp", "shenanigans",
		},
	},
	Love: {
		Name:  "Love",
		Emoji: "💕",
		Color: "#FF1493",
		Keywords: []string{
			"love", "heart", "kiss", "hug", "darling",
			"sweetheart", "adore", "cherish", "smitten", "crush",
			"beloved", "romance", "tender", "embrace", "affection",
		},
	},
}

// Info returns the display metadata for the emotion.
func (e Emotion) Info() EmotionInfo {
	if info, ok := emotionRegistry[e]; ok {
		return info
	}
	return emotionRegistry[Curiosity]
}

// String returns the display name of the emotion.
func (e Emotion) String() string {
	return e.Info().Name
}

// AllEmotions returns a slice of all eight emotions.
func AllEmotions() []Emotion {
	return []Emotion{Joy, Sadness, Anger, Fear, Curiosity, Sleepy, Silly, Love}
}

// GetEyes returns just the eye expression string for an emotion.
func GetEyes(e Emotion) string {
	switch e {
	case Joy:
		return "◕ ‿ ◕"
	case Sadness:
		return "◕ ︵ ◕"
	case Anger:
		return "◣ ᨎ ◢"
	case Fear:
		return "◉ _ ◉"
	case Curiosity:
		return "◕ ‿ ◔"
	case Sleepy:
		return "◡ _ ◡"
	case Silly:
		return "◔ ‿ ◕"
	case Love:
		return "♥ ‿ ♥"
	default:
		return "◕ ‿ ◕"
	}
}
