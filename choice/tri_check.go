package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TriState is off, on or mixed.
type TriState uint8

// The three states of a TriCheck.
const (
	TriOff TriState = iota
	TriOn
	TriMixed
)

// TriCheck is a check whose mark morphs to a dash.
func TriCheck(c *ui.Context, label string, state *TriState) ui.Element {
	if label == "" || state == nil || *state > TriMixed {
		panic("choice: TriCheck needs a label and a valid state")
	}
	t := c.Theme()
	on := *state == TriOn
	root := ui.CheckboxBase(c, &on).Label(label).Gap(t.Space(2))
	if *state == TriMixed {
		root.Mixed()
	}
	shown := *state
	root.Children(func() {
		box(c, root, shown == TriOn, shown == TriMixed)
		ui.Text(c, label)
	})
	if kit.Changed(root, sound.Toggle) {
		*state = TriOff
		if on {
			*state = TriOn
		}
		kit.Log("TriCheck", "toggled", "state", *state)
	}
	return root
}
