package overlay

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MenuItem is a label, Lucide icon and shortcut text.
type MenuItem struct{ Label, Icon, Keys string }

// Menu is a button opening a styled list of items.
func Menu(c *ui.Context, label string, list []MenuItem, chosen *int, open *bool) ui.Element {
	items("Menu", label, list, chosen)
	if open == nil {
		panic("overlay: Menu needs open")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	shown := *chosen
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(label).Gap(t.Space(1.5)).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)).Expanded(*open)
		trigger.Background(trigger.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || trigger.Hovered())), ui.Snappy)).Children(func() {
			ui.Text(c, label).SingleLine()
			chev := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16)
			chev.Rotate(180 * chev.AnimateTo("turn", unit(*open), ui.Snappy))
		})
		if trigger.Clicked() {
			*open = !*open
			root.SetChanged()
		}
		if ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border).
				Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			entries(c, list, shown, func(i int) {
				*chosen, *open = i, false
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Menu", "chose", "index", i)
			})
		}).Dismissed() {
			root.SetChanged()
		}
	})
	if opened, _ := moved(root, *open, "Menu"); opened {
		sound.Play(sound.Open)
	}
	return root
}
