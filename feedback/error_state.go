package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ErrorStateOptions tune an ErrorState.
type ErrorStateOptions struct {
	// Title heads it; "Something went wrong" when empty.
	Title string
	// Retry labels the button; "Retry" when empty.
	Retry string
}

var shake = []ui.Keyframe{
	{Value: 0}, {At: 70 * time.Millisecond, Value: -6}, {At: 140 * time.Millisecond, Value: 6},
	{At: 210 * time.Millisecond, Value: -4}, {At: 280 * time.Millisecond, Value: 4}, {At: 350 * time.Millisecond, Value: 0},
}

// ErrorState is a shaking error icon, message and Retry.
func ErrorState(c *ui.Context, message string, o ErrorStateOptions) (root, retry ui.Element) {
	if o.Title == "" {
		o.Title = "Something went wrong"
	}
	if o.Retry == "" {
		o.Retry = "Retry"
	}
	root, retry = state(c, icon.Must("circle-alert"), c.Theme().Danger, o.Title, message, o.Retry, false, func(i ui.Element) {
		i.Translate(i.Keyframes("shake", shake), 0)
	})
	if once(root) {
		sound.Play(sound.Error)
		kit.Log("ErrorState", "shown")
	}
	if kit.Clicked(retry, sound.Tap) {
		kit.Log("ErrorState", "retry")
	}
	return root, retry
}
