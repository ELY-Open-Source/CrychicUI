package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Action is an action's label and Lucide icon.
type Action struct{ Label, Icon string }

// BulkBarOptions tune a BulkBar.
type BulkBarOptions struct {
	// Selected takes the count.
	Selected, Clear string
}

// BulkBar slides up actions while count > 0.
func BulkBar(c *ui.Context, count int, actions []Action, chosen *int, o BulkBarOptions) (root, clear ui.Element) {
	actionsValid("BulkBar", actions)
	if count < 0 || chosen == nil || *chosen < -1 || *chosen >= len(actions) {
		panic("data: BulkBar needs a count and a choice")
	}
	def(&o.Selected, "%d selected")
	def(&o.Clear, "Clear selection")
	t := c.Theme()
	fill, ink := t.Inverse, t.InverseText
	if fill.A == 0 {
		fill = t.Text
	}
	if ink.A == 0 {
		ink = t.Background
	}
	root = ui.Box(c).FillWidth()
	root.Children(func() {
		if count == 0 {
			return
		}
		bar := ui.Row(c).Key("bar").FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius*1.5).
			Background(fill).TextColor(ink).Shadow(0, t.Space(2), t.Space(6), 0, kit.Shade(0.2)).
			Transition(ui.ElementTransition{Animation: &ui.Snappy, Exit: &ui.Motion{Y: t.Space(6)}})
		kit.Arrive(bar, ui.Pose{Y: t.Space(6), Fade: 1}, 0).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Selected, count)).Bold().Grow(1)
			for i, a := range actions {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(a.Label).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", ink.Alpha(0.12*unit(b.Hovered())), ui.Snappy)).Children(func() {
					ui.Icon(c, icon.Must(a.Icon)).Size(16, 16)
					ui.Text(c, a.Label)
				})
				if b.Clicked() {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("BulkBar", "chose", "index", i, "count", count)
				}
			}
			clear = kit.Press(ui.ButtonBase(c).Label(o.Clear).Tooltip(o.Clear).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius))
			clear.Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
			if clear.Clicked() {
				sound.Play(sound.Tap)
				kit.Log("BulkBar", "cleared")
			}
		})
	})
	return root, clear
}
