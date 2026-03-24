package octopus

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	normalTickInterval     = 200 * time.Millisecond
	transitionTickInterval = 100 * time.Millisecond
	squishTicks            = 2
)

// AnimTickMsg is the message type sent on each animation tick.
type AnimTickMsg time.Time

// AnimTick returns a bubbletea Cmd that fires an AnimTickMsg after the
// standard animation interval.
func AnimTick() tea.Cmd {
	return tea.Tick(normalTickInterval, func(t time.Time) tea.Msg {
		return AnimTickMsg(t)
	})
}

// AnimTickFast returns a bubbletea Cmd with the faster transition interval.
func AnimTickFast() tea.Cmd {
	return tea.Tick(transitionTickInterval, func(t time.Time) tea.Msg {
		return AnimTickMsg(t)
	})
}

// AnimationState drives frame-based animation with smooth transitions
// between emotions. It integrates with bubbletea's tick-based update loop.
type AnimationState struct {
	current        Emotion
	target         Emotion
	frame          int
	transitioning  bool
	transitionTick int    // counts ticks within a transition
	lastFrame      string // last frame produced by Tick()
}

// NewAnimationState creates an AnimationState starting in the idle/Curiosity
// pose so the octopus is gently moving from the moment it appears.
func NewAnimationState() *AnimationState {
	return &AnimationState{
		current: Curiosity,
		target:  Curiosity,
	}
}

// SetEmotion triggers a transition to the given emotion. If the octopus is
// already displaying that emotion, this is a no-op.
func (a *AnimationState) SetEmotion(e Emotion) {
	if e == a.current && !a.transitioning {
		return
	}
	a.target = e
	a.transitioning = true
	a.transitionTick = 0
}

// IsTransitioning reports whether the animation is mid-transition.
func (a *AnimationState) IsTransitioning() bool {
	return a.transitioning
}

// TickInterval returns the Duration that should be used for the next tick
// command — faster during transitions for a snappy feel.
func (a *AnimationState) TickInterval() time.Duration {
	if a.transitioning {
		return transitionTickInterval
	}
	return normalTickInterval
}

// CurrentEmotion returns the emotion currently being displayed.
func (a *AnimationState) CurrentEmotion() Emotion {
	return a.current
}

// CurrentFrame returns the ASCII frame that was last produced by Tick,
// without advancing animation state. Safe to call from View().
func (a *AnimationState) CurrentFrame() string {
	if a.lastFrame == "" {
		return getFrameForEmotion(a.current, 0)
	}
	return a.lastFrame
}

// Tick advances the animation by one frame and returns the ASCII art to
// render. Call this from your bubbletea Update when you receive an
// AnimTickMsg.
//
// Normal mode: alternates between frame 0 and frame 1 of the current
// emotion (or the idle animation when no emotion is active).
//
// Transition mode: plays a brief "squish" effect (compressed frame for
// squishTicks), then expands into frame 0 of the target emotion.
func (a *AnimationState) Tick() string {
	if a.transitioning {
		a.transitionTick++

		if a.transitionTick <= squishTicks {
			// Squish phase — show a vertically compressed frame.
			a.lastFrame = squishFrame(a.current)
			return a.lastFrame
		}

		// Expand into the new emotion on the first post-squish tick.
		a.current = a.target
		a.transitioning = false
		a.transitionTick = 0
		a.frame = 0
		a.lastFrame = getFrameForEmotion(a.current, 0)
		return a.lastFrame
	}

	// Normal animation: alternate between the two frames.
	a.lastFrame = getFrameForEmotion(a.current, a.frame)
	a.frame = (a.frame + 1) % 2
	return a.lastFrame
}

// getFrameForEmotion returns the requested frame for an emotion, falling
// back to the idle animation for unknown values.
func getFrameForEmotion(e Emotion, idx int) string {
	frames := GetFrames(e)
	if len(frames) == 0 {
		frames = GetIdleFrames()
	}
	if idx < 0 || idx >= len(frames) {
		idx = 0
	}
	return frames[idx]
}

// squishFrame produces a vertically compressed version of frame 0 for the
// given emotion, keeping the head and collapsing the body so the octopus
// looks like it's "squishing" before springing into a new pose.
func squishFrame(e Emotion) string {
	frames := GetFrames(e)
	if len(frames) == 0 {
		frames = GetIdleFrames()
	}
	lines := strings.Split(frames[0], "\n")

	// We keep the top portion (head, ~6 lines), skip the middle, and keep
	// a compressed bottom. The goal is roughly half the original height
	// padded back to FrameHeight() so the layout doesn't jump.
	height := FrameHeight()

	headEnd := 6
	if headEnd > len(lines) {
		headEnd = len(lines)
	}
	head := lines[:headEnd]

	// Take last 2 body lines as a compressed "foot".
	var foot []string
	if len(lines) > headEnd+2 {
		foot = lines[len(lines)-2:]
	} else if len(lines) > headEnd {
		foot = lines[headEnd:]
	}

	// Build the squished frame: head + foot + blank padding to full height.
	var b strings.Builder
	written := 0
	for _, l := range head {
		if written > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(l)
		written++
	}
	for _, l := range foot {
		b.WriteByte('\n')
		b.WriteString(l)
		written++
	}
	// Pad remaining lines so the frame height stays constant.
	blank := strings.Repeat(" ", FrameWidth())
	for written < height {
		b.WriteByte('\n')
		b.WriteString(blank)
		written++
	}
	return b.String()
}
