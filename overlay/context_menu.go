package overlay

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ContextMenu opens a styled menu where target is right-clicked.
func ContextMenu(c *ui.Context, label string, target ui.Element, list []MenuItem, chosen *int, open *bool) ui.Element {
	items("ContextMenu", label, list, chosen)
	if open == nil {
		panic("overlay: ContextMenu needs open")
	}
	t := c.Theme()
	root := ui.Box(c)
	at := ui.Local(root, "at", func() ui.Point { return ui.Point{} })
	// System menus cannot be styled; this reads the click itself.
	again := target.RightClicked()
	if again {
		x, y, _ := target.PointerPosition()
		*at, *open = ui.Point{X: x, Y: y}, true
		root.SetChanged()
	}
	if target.Shortcut(ui.Shift, ui.KeyF10) {
		*at, *open = ui.Point{Y: target.Bounds().H}, true
		root.SetChanged()
	}
	shown := *chosen
	if *open {
		ui.Overlay(c, func() {
			// The panel flips off this point to fit.
			point := ui.Box(c).AttachTo(target, ui.AnchorTopLeft, ui.AnchorTopLeft).Left(at.X).Top(at.Y).Width(0).Height(0)
			panel := ui.Box(c).AttachTo(point, ui.AnchorBottomRight, ui.AnchorTopLeft).Padding(t.Space(1)).
				Radius(t.Radius).Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).
				Role(ui.RoleMenu).Label(label).TransformOrigin(0, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0).Children(func() {
				entries(c, list, shown, func(i int) {
					*chosen, *open = i, false
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("ContextMenu", "chose", "index", i)
				})
			})
			if !again && panel.PressedOutside() || panel.OverlayShortcut(0, ui.KeyEscape) {
				*open = false
				root.SetChanged()
			}
		})
	}
	if opened, _ := moved(root, *open, "ContextMenu"); opened {
		sound.Play(sound.Open)
	}
	return root
}
