package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SubmitState is a SubmitButton's current phase.
type SubmitState int

// Submit phases.
const (
	SubmitIdle SubmitState = iota
	SubmitLoading
	SubmitSuccess
	SubmitError
)

// SubmitButton morphs through idle, loading, success and error.
func SubmitButton(c *ui.Context, label string, state SubmitState) ui.Element {
	if state < SubmitIdle || state > SubmitError {
		panic("button: SubmitButton state out of range")
	}
	t := c.Theme()
	root := textBase(c, label).Disabled(state == SubmitLoading).Border(1, t.Border).Background(t.Surface)
	previous := ui.Local(root, "previous", func() SubmitState { return state })
	if state != *previous {
		if state == SubmitSuccess {
			sound.Play(sound.Success)
		} else if state == SubmitError {
			sound.Play(sound.Error)
		}
		*previous = state
	}
	root.Children(func() {
		name := "arrow-up"
		switch state {
		case SubmitLoading:
			name = "loader-circle"
		case SubmitSuccess:
			name = "check"
		case SubmitError:
			name = "x"
		}
		glyph := ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
		if state == SubmitLoading {
			glyph.Rotate(glyph.Loop("spin", time.Second, ui.Linear) * 360)
		}
		ui.Text(c, label).SingleLine()
	})
	return root
}
