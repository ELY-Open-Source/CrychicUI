package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConfirmButtonOptions tune a ConfirmButton.
type ConfirmButtonOptions struct {
	// Prompt labels the armed state.
	Prompt string
}

// ConfirmButton arms once, then confirms within three seconds.
func ConfirmButton(c *ui.Context, label string, o ConfirmButtonOptions) ui.Element {
	if o.Prompt == "" {
		o.Prompt = "Sure?"
	}
	t := c.Theme()
	b := textBase(c, label).Border(1, t.Border).Transition(ui.ElementTransition{Size: true})
	armed := ui.Local(b, "armed", func() bool { return false })
	until := ui.Local(b, "until", func() time.Time { return time.Time{} })
	if *armed {
		if !c.Now().Before(*until) {
			*armed = false
			kit.Log("ConfirmButton", "disarmed")
		} else {
			c.After(until.Sub(c.Now()))
		}
	}
	if b.Clicked() {
		if *armed {
			*armed = false
			sound.Play(sound.Delete)
			b.SetChanged()
			kit.Log("ConfirmButton", "confirmed")
		} else {
			*armed = true
			*until = c.Now().Add(3 * time.Second)
			c.After(3 * time.Second)
			sound.Play(sound.Warning)
			kit.Log("ConfirmButton", "armed")
		}
	}
	bg, fg, text := t.Surface, t.Text, label
	if *armed {
		bg, fg, text = t.Danger, t.AccentText, o.Prompt
	}
	b.Label(text).Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(fg).Children(func() {
		ui.Text(c, text).SingleLine().ContentTransition(ui.ContentTransition{})
	})
	return b
}
