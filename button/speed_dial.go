package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SpeedDialAction names one fanned-out action.
type SpeedDialAction struct{ Label, Icon string }

// SpeedDialOptions tune a SpeedDial.
type SpeedDialOptions struct {
	// Label names the main action.
	Label string
}

// SpeedDial fans actions out with a stagger.
func SpeedDial(c *ui.Context, open *bool, actions []SpeedDialAction, selected *int, o SpeedDialOptions) ui.Element {
	if open == nil || selected == nil || len(actions) == 0 || *selected < -1 || *selected >= len(actions) {
		panic("button: SpeedDial needs state and actions")
	}
	if o.Label == "" {
		o.Label = "Actions"
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).AlignItems(ui.End).Gap(t.Space(2))
	return root.Children(func() {
		if *open {
			for i, action := range actions {
				c.Key(i)
				b := textBase(c, action.Label).Border(1, t.Border).Background(t.Surface)
				a := ui.Snappy
				a.Delay = time.Duration(i) * 30 * time.Millisecond
				b.Transition(ui.ElementTransition{Animation: &a, Enter: &ui.Motion{Y: 8}, Exit: &ui.Motion{Y: 8}}).
					Children(func() {
						ui.Icon(c, icon.Must(action.Icon)).Size(16, 16)
						ui.Text(c, action.Label).SingleLine()
					})
				if kit.Clicked(b, sound.Tap) {
					*selected = i
					*open = false
					root.SetChanged()
					kit.Log("SpeedDial", "selected", "index", i)
				}
			}
		}
		name := "plus"
		if *open {
			name = "x"
		}
		main := fabBase(c.Key("main"), o.Label, name)
		if main.Clicked() {
			*open = !*open
			root.SetChanged()
			if *open {
				sound.Play(sound.Open)
			} else {
				sound.Play(sound.Close)
			}
			kit.Log("SpeedDial", "toggled", "open", *open)
		}
	})
}
