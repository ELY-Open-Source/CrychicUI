package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Caption is text shown from From until To.
type Caption struct {
	From, To time.Duration
	Text     string
}

// CaptionsOptions tune a Captions.
type CaptionsOptions struct{}

// Captions shows the caption at pos, fading between lines.
func Captions(c *ui.Context, label string, lines []Caption, pos time.Duration, o CaptionsOptions) ui.Element {
	if label == "" || len(lines) == 0 || pos < 0 {
		panic(fmt.Sprintf("media: Captions %q needs a label, lines and a pos not negative", label))
	}
	var end time.Duration
	for _, l := range lines {
		if l.Text == "" || l.From < end || l.To <= l.From {
			panic(fmt.Sprintf("media: Captions %q lines need text and must run in order without overlap", label))
		}
		end = l.To
	}
	t := c.Theme()
	at := -1
	for i, l := range lines {
		if pos >= l.From && pos < l.To {
			at = i
		}
	}
	ink := t.Background
	if t.Dark {
		ink = t.Text
	}
	root := ui.Grid(c).Columns(1).FillWidth().MinHeight(t.Space(10)).Label(label)
	return root.Children(func() {
		if at < 0 {
			return
		}
		c.Key(at)
		row := ui.Row(c).ColumnStart(1).RowStart(1).FillWidth().Justify(ui.Center).AlignSelf(ui.End).
			Transition(ui.ElementTransition{Duration: 200 * time.Millisecond, Exit: &ui.Motion{}})
		kit.Arrive(row, ui.Pose{Fade: 1, Y: t.Space(1)}, 0).Children(func() {
			ui.Text(c, lines[at].Text).TextAlign(ui.Center).Padding(t.Space(1), t.Space(2.5)).Radius(t.Radius).
				Background(kit.Shade(0.72)).TextColor(ink)
		})
	})
}
