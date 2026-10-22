package motion

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SequenceOptions tune a Sequence.
type SequenceOptions struct {
	// Action labels the button; "Continue" when empty.
	Action string
}

// Sequence brings in title, body, action in turn; run replays.
func Sequence(c *ui.Context, run uint64, title, body string, o SequenceOptions) (root, action ui.Element) {
	if title == "" || body == "" {
		panic("motion: Sequence needs a title and a body")
	}
	if o.Action == "" {
		o.Action = "Continue"
	}
	t := c.Theme()
	root = kit.Card(c).Padding(t.Space(4)).Gap(t.Space(2)).AlignItems(ui.Start)
	tl := root.Timeline("crychic.sequence", ui.TimelineOptions{Run: run})
	in := ui.Tween{From: 0, To: 1, Duration: 240 * time.Millisecond}
	a := tl.Tween(in)
	b := tl.Tween(in, ui.AfterPrevious(-120*time.Millisecond))
	e := tl.Tween(in, ui.AfterPrevious(-120*time.Millisecond))
	rise := t.Space(3)
	root.Children(func() {
		ui.Text(c, title).Bold().FontSize(t.Rem(1.25)).Opacity(a).Translate(0, (1-a)*rise)
		kit.Muted(c, body).Opacity(b).Translate(0, (1-b)*rise)
		action = ui.Button(c, o.Action).Opacity(e).Translate(0, (1-e)*rise)
		if kit.Clicked(action, sound.Tap) {
			kit.Log("Sequence", "action")
		}
	})
	return root, action
}
