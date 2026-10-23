package media

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// opener is a pill whose chevron turns.
func opener(c *ui.Context, name, text string, open bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(name+": "+text).Tooltip(name).AlignSelf(ui.Start).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).
		Padding(t.Space(1.5), t.Space(3)).Radius(999).Border(1, t.Border))
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered() || open)), ui.Snappy))
	return b.Children(func() {
		kit.Numeric(c, text).ContentTransition(ui.ContentTransition{})
		ui.Icon(c, icon.Must("chevron-down")).Size(16, 16).AnimatePose(ui.Pose{Rotate: 180 * unit(open)}, ui.Snappy)
	})
}

// menu lists names under trigger, checking one.
func menu(c *ui.Context, trigger ui.Element, component, label string, names []string, at *int, open *bool) {
	t := c.Theme()
	if trigger.Clicked() {
		*open = !*open
		trigger.SetChanged()
		if *open {
			sound.Play(sound.Open)
		} else {
			sound.Play(sound.Close)
		}
		kit.Log(component, "toggled", "open", *open)
	}
	shown := *at
	if ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
		panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(36)).Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).
			Border(1, t.Border).Shadow(0, t.Space(2), t.Space(6), 0, kit.Shade(0.16)).TransformOrigin(0, 0).Role(ui.RoleMenu).
			Label(label)
		kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
		ui.Column(c).FillWidth().FocusGroup(ui.Vertical).Children(func() {
			for i, name := range names {
				c.Key(i)
				on := i == shown
				item := kit.Press(ui.ButtonBase(c).Label(name).Checked(on).FillWidth().Row().Gap(t.Space(2)).
					AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius))
				if on {
					item.AutoFocus()
				}
				item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy)).Children(func() {
					ui.Box(c).Size(16, 16).Shrink(0).Children(func() {
						if on {
							kit.Arrive(ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Accent), ui.Pose{Scale: 0.4, Fade: 1}, 0)
						}
					})
					ui.Text(c, name).SingleLine()
				})
				if item.Clicked() {
					*at, *open = i, false
					trigger.SetChanged()
					sound.Play(sound.Select)
					kit.Log(component, "chose", "index", i)
				}
			}
		})
	}).Dismissed() {
		trigger.SetChanged()
		kit.Log(component, "dismissed")
	}
}
