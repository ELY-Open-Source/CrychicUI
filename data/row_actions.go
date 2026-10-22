package data

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RowActions is a "…" button opening actions.
func RowActions(c *ui.Context, label string, actions []Action, chosen *int, open *bool) ui.Element {
	actionsValid("RowActions", actions)
	if label == "" || chosen == nil || open == nil || *chosen < -1 || *chosen >= len(actions) {
		panic("data: RowActions needs a label, a choice and open")
	}
	t := c.Theme()
	root := ui.Box(c).Shrink(0)
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius)).Expanded(*open)
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || b.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("ellipsis")).Size(16, 16)
		})
		if b.Clicked() {
			*open = !*open
			root.SetChanged()
			if *open {
				sound.Play(sound.Open)
			}
			kit.Log("RowActions", "toggled", "open", *open)
		}
		if ui.PopoverBase(c, b, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).Width(t.Space(48)).Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).
				Border(1, t.Border).Shadow(0, t.Space(2), t.Space(6), 0, kit.Shade(0.16)).TransformOrigin(1, 0).Role(ui.RoleMenu).Label(label)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			ui.Column(c).FillWidth().FocusGroup(ui.Vertical).Children(func() {
				for i, a := range actions {
					c.Key(i)
					item := kit.Press(ui.ButtonBase(c).Label(a.Label).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
						Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius))
					item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy)).Children(func() {
						ui.Icon(c, icon.Must(a.Icon)).Size(16, 16)
						ui.Text(c, a.Label)
					})
					if item.Clicked() {
						*chosen, *open = i, false
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("RowActions", "chose", "index", i)
					}
				}
			})
		}).Dismissed() {
			root.SetChanged()
			kit.Log("RowActions", "dismissed")
		}
	})
	return root
}
