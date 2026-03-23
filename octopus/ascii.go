package octopus



// OctopusFrame holds a single frame of ASCII art.
type OctopusFrame struct {
	Art string
}

const (
	frameWidth  = 28
	frameHeight = 16
)

// FrameWidth returns the max width of any frame.
func FrameWidth() int { return frameWidth }

// FrameHeight returns the max height of any frame.
func FrameHeight() int { return frameHeight }

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

// GetIdleFrames returns 2 gentle "breathing" animation frames.
func GetIdleFrames() []string {
	return []string{
		// Idle frame 1 — relaxed
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◕ ‿ ◕ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`     /  / | \  \        ` + "\n" +
			`    /  /  |  \  \       ` + "\n" +
			`   |  |   |   |  |      ` + "\n" +
			`   |  |   |   |  |      ` + "\n" +
			`    \  \  |  /  /       ` + "\n" +
			`     \  \ | /  /        ` + "\n" +
			`      \  \|/  /         ` + "\n" +
			`       \_/|\_/          ` + "\n" +
			`          |             ` + "\n" +
			`         ~~~            `,

		// Idle frame 2 — slight inhale
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◕ ‿ ◕ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`      / / | \ \         ` + "\n" +
			`     / /  |  \ \        ` + "\n" +
			`    |  |  |  |  |       ` + "\n" +
			`    |  |  |  |  |       ` + "\n" +
			`     \  \ | /  /        ` + "\n" +
			`      \  \|/  /         ` + "\n" +
			`       \ /|\ /          ` + "\n" +
			`        \_|_/           ` + "\n" +
			`          |             ` + "\n" +
			`         ~~~            `,
	}
}

// GetFrames returns 2 animation frames for the given emotion.
func GetFrames(emotion Emotion) []string {
	switch emotion {
	case Joy:
		return joyFrames()
	case Sadness:
		return sadnessFrames()
	case Anger:
		return angerFrames()
	case Fear:
		return fearFrames()
	case Curiosity:
		return curiosityFrames()
	case Sleepy:
		return sleepyFrames()
	case Silly:
		return sillyFrames()
	case Love:
		return loveFrames()
	default:
		return GetIdleFrames()
	}
}

func joyFrames() []string {
	return []string{
		// Joy frame 1 — arms wiggling outward
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◕ ‿ ◕ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`   ~ / / ||| \ \ ~     ` + "\n" +
			`    ~/  / ||| \  \~    ` + "\n" +
			`    |  /  |||  \  |     ` + "\n" +
			`    | |   |||   | |     ` + "\n" +
			`     \ \  |||  / /      ` + "\n" +
			`      \ \ ||| / /       ` + "\n" +
			`       \_\|||/_/        ` + "\n" +
			`          |||           ` + "\n" +
			`          |||           ` + "\n" +
			`           V            `,

		// Joy frame 2 — arms wiggling the other way
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◕ ‿ ◕ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`    ~ / / || \ \ ~      ` + "\n" +
			`     ~/  /|||\  \~     ` + "\n" +
			`      | / ||| \ |       ` + "\n" +
			`      ||  |||  ||       ` + "\n" +
			`      | \ ||| / |       ` + "\n" +
			`     / /  ||| \  \      ` + "\n" +
			`    /_/ \_|||_/ \_\     ` + "\n" +
			`          |||           ` + "\n" +
			`          |||           ` + "\n" +
			`           V            `,
	}
}

func sadnessFrames() []string {
	return []string{
		// Sadness frame 1 — arms hanging limp
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◕ ︵ ◕ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`       |||||||          ` + "\n" +
			`       |||||||          ` + "\n" +
			`       |||||||          ` + "\n" +
			`       ||| |||          ` + "\n" +
			`       ||| |||          ` + "\n" +
			`       ||   ||          ` + "\n" +
			`       ||   ||          ` + "\n" +
			`       |     |          ` + "\n" +
			`       |     |          ` + "\n" +
			`                        `,

		// Sadness frame 2 — subtle droop shift
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◕ ︵ ◕ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`       |||||||          ` + "\n" +
			`       |||||||          ` + "\n" +
			`       |||||||          ` + "\n" +
			`       ||| |||          ` + "\n" +
			`       ||| |||          ` + "\n" +
			`       ||   ||          ` + "\n" +
			`       |     |          ` + "\n" +
			`       |     |          ` + "\n" +
			`      .       .         ` + "\n" +
			`                        `,
	}
}

func angerFrames() []string {
	return []string{
		// Anger frame 1 — arms coiled and thrust outward
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◣ ᨎ ◢ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`  \\\  / ||| \  ///     ` + "\n" +
			`   \\\/ /|||\  \///    ` + "\n" +
			`    \\| / ||| \ |//     ` + "\n" +
			`     \||  |||  ||/      ` + "\n" +
			`      ||  |||  ||       ` + "\n" +
			`     /||  |||  ||\      ` + "\n" +
			`    //|\  |||  /|\\     ` + "\n" +
			`   ///  \_|||_/  \\\    ` + "\n" +
			`  ///     |||     \\\   ` + "\n" +
			`          V             `,

		// Anger frame 2 — arms flex tighter
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◣ ᨎ ◢ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`   \\\ / ||| \ ///     ` + "\n" +
			`    \\\/ ||| \///      ` + "\n" +
			`     \\| ||| |//       ` + "\n" +
			`      \| ||| |/        ` + "\n" +
			`       | ||| |          ` + "\n" +
			`      /| ||| |\         ` + "\n" +
			`     //| ||| |\\        ` + "\n" +
			`    /// \_|||_/ \\\     ` + "\n" +
			`   ///    |||    \\\    ` + "\n" +
			`          V             `,
	}
}

func fearFrames() []string {
	return []string{
		// Fear frame 1 — arms curled inward, hugging body
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◉ _ ◉ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`       /|||||\          ` + "\n" +
			`      | ||||| |         ` + "\n" +
			`      | ||||| |         ` + "\n" +
			`       \|||||/          ` + "\n" +
			`       /|||||\          ` + "\n" +
			`      | || || |         ` + "\n" +
			`       \|   |/          ` + "\n" +
			`        |   |           ` + "\n" +
			`         \_/            ` + "\n" +
			`                        `,

		// Fear frame 2 — trembling slightly
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◉ _ ◉ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`      /||||||\          ` + "\n" +
			`     | |||||| |         ` + "\n" +
			`     | |||||| |         ` + "\n" +
			`      \||||||/          ` + "\n" +
			`      /||||||\          ` + "\n" +
			`     | |||  | |         ` + "\n" +
			`      \|    |/          ` + "\n" +
			`       |    |           ` + "\n" +
			`        \__/            ` + "\n" +
			`                        `,
	}
}

func curiosityFrames() []string {
	return []string{
		// Curiosity frame 1 — one arm reaching up/out
		`         ___        ?   ` + "\n" +
			`       /     \     /   ` + "\n" +
			`      |       |   /    ` + "\n" +
			`      | ◕ ‿ ◔ |  /     ` + "\n" +
			`       \_____/ _/      ` + "\n" +
			`        |||||/         ` + "\n" +
			`     /  |||||          ` + "\n" +
			`    /  / |||           ` + "\n" +
			`   |  |  |||           ` + "\n" +
			`   |  |  |||           ` + "\n" +
			`    \  \ |||           ` + "\n" +
			`     \  \|||           ` + "\n" +
			`      \  |||           ` + "\n" +
			`       \_|||           ` + "\n" +
			`         |||           ` + "\n" +
			`          V            `,

		// Curiosity frame 2 — arm shifts slightly
		`         ___      ?     ` + "\n" +
			`       /     \    |    ` + "\n" +
			`      |       |   |    ` + "\n" +
			`      | ◕ ‿ ◔ |  /     ` + "\n" +
			`       \_____/_/       ` + "\n" +
			`        |||||/         ` + "\n" +
			`     /  |||||          ` + "\n" +
			`    /  / |||           ` + "\n" +
			`   |  /  |||           ` + "\n" +
			`   | |   |||           ` + "\n" +
			`    \ \  |||           ` + "\n" +
			`     \ \ |||           ` + "\n" +
			`      \ \|||           ` + "\n" +
			`       \_|||           ` + "\n" +
			`         |||           ` + "\n" +
			`          V            `,
	}
}

func sleepyFrames() []string {
	return []string{
		// Sleepy frame 1 — arms drooping, z's floating
		`         ___    z       ` + "\n" +
			`       /     \   z     ` + "\n" +
			`      |       |   Z    ` + "\n" +
			`      | ◡ _ ◡ |        ` + "\n" +
			`       \_____/         ` + "\n" +
			`        |||||          ` + "\n" +
			`       |||||||         ` + "\n" +
			`      / ||||| \        ` + "\n" +
			`     / /||||| .\       ` + "\n" +
			`    . / |||||  .\      ` + "\n" +
			`   .  / || ||   .      ` + "\n" +
			`   . /  ||  |   .      ` + "\n" +
			`    ./  |    \  .      ` + "\n" +
			`    .   |     \.       ` + "\n" +
			`    .   |      .       ` + "\n" +
			`        |              `,

		// Sleepy frame 2 — z's drift up
		`         ___      z     ` + "\n" +
			`       /     \    z    ` + "\n" +
			`      |       |    Z   ` + "\n" +
			`      | ◡ _ ◡ |        ` + "\n" +
			`       \_____/         ` + "\n" +
			`        |||||          ` + "\n" +
			`       |||||||         ` + "\n" +
			`      / ||||| \        ` + "\n" +
			`     / /|||||. \       ` + "\n" +
			`    . / |||||  .\      ` + "\n" +
			`   .  / || ||   .      ` + "\n" +
			`   . /  ||  |   .      ` + "\n" +
			`    ./  |    \  .      ` + "\n" +
			`    .   |     \.       ` + "\n" +
			`    .   |      .       ` + "\n" +
			`        |              `,
	}
}

func sillyFrames() []string {
	return []string{
		// Silly frame 1 — arms tangled in goofy knots
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◔ ‿ ◕ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`    /\  |||||  /\       ` + "\n" +
			`   /  \/|||||\/ .\      ` + "\n" +
			`   \  /\|||||/\  /      ` + "\n" +
			`    \/  \|||/  \/       ` + "\n" +
			`    /\   |||   /\       ` + "\n" +
			`   /  \  |||  /  \      ` + "\n" +
			`   \   \_|||_/   /      ` + "\n" +
			`    \    |||    /       ` + "\n" +
			`     \   |||   /        ` + "\n" +
			`          V             `,

		// Silly frame 2 — arms tangled differently
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ◔ ‿ ◕ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`    \/  |||||  \/       ` + "\n" +
			`    /\  |||||  /\       ` + "\n" +
			`   /  \/|||||\/ .\      ` + "\n" +
			`   \  /\|||||/\  /      ` + "\n" +
			`    \/  \|||/  \/       ` + "\n" +
			`    /    |||    \       ` + "\n" +
			`   /   /_|||_\   \      ` + "\n" +
			`   \    /|||\    /      ` + "\n" +
			`    \   |||    /        ` + "\n" +
			`         V              `,
	}
}

func loveFrames() []string {
	return []string{
		// Love frame 1 — two arms forming a heart, others gentle
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ♥ ‿ ♥ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`     _  |||||  _        ` + "\n" +
			`    / \ ||||| / \       ` + "\n" +
			`   |   \|||||/   |      ` + "\n" +
			`    \   |||||   /       ` + "\n" +
			`     \  |||||  /        ` + "\n" +
			`      \ || || /         ` + "\n" +
			`       \|   |/          ` + "\n" +
			`        \   /           ` + "\n" +
			`         \ /            ` + "\n" +
			`          V             `,

		// Love frame 2 — heart pulses slightly bigger
		`         ___            ` + "\n" +
			`       /     \          ` + "\n" +
			`      |       |         ` + "\n" +
			`      | ♥ ‿ ♥ |         ` + "\n" +
			`       \_____/          ` + "\n" +
			`        |||||           ` + "\n" +
			`    _   |||||   _       ` + "\n" +
			`   / \  |||||  / \      ` + "\n" +
			`  |   \ ||||| /   |     ` + "\n" +
			`   \   \|||||/   /      ` + "\n" +
			`    \   |||||   /       ` + "\n" +
			`     \  || ||  /        ` + "\n" +
			`      \ |   | /         ` + "\n" +
			`       \|   |/          ` + "\n" +
			`        \   /           ` + "\n" +
			`         \_/            `,
	}
}
