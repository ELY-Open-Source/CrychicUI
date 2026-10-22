// Package overlay holds dialogs, sheets, popovers and menus.
package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// moved reports open turning on or off.
func moved(root ui.Element, open bool, component string) (opened, closed bool) {
	was := ui.Local(root, "open", func() bool { return false })
	opened, closed = open && !*was, !open && *was
	*was = open
	if opened || closed {
		kit.Log(component, "open", open)
	}
	return opened, closed
}

// dialog shows a styled panel over a dim while open.
func dialog(c *ui.Context, open *bool, width float32, body func()) ui.Element {
	t := c.Theme()
	return ui.DialogBase(c, open, func(backdrop, panel ui.Element) {
		kit.Arrive(backdrop.Background(kit.Shade(0.4)), ui.Pose{Fade: 1}, 0)
		panel.FillWidth().MaxWidth(width).Padding(t.Space(5)).Radius(t.Radius*1.5).Background(t.Background).Shadow(0, 12, 32, 0, kit.Shade(0.24))
		kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Children(body)
	})
}

// heading is a title row with a close button.
func heading(c *ui.Context, title, closeLabel string) (shut ui.Element) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		ui.Text(c, title).Bold().FontSize(t.Rem(1.1)).Grow(1).SingleLine()
		shut = kit.Press(ui.ButtonBase(c).Label(closeLabel).Tooltip(closeLabel).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius))
		shut.Background(shut.AnimateColor("bg", t.SurfaceHover.Alpha(unit(shut.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("x")).Size(16, 16)
		})
	})
	return shut
}

// primary is a filled button in tone; its caller cues.
func primary(c *ui.Context, label string, tone, ink ui.Color) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(2), t.Space(4)).Radius(t.Radius))
	return b.Background(b.AnimateColor("bg", tone.Mix(ink, 0.12*unit(b.Hovered())), ui.Snappy)).TextColor(ink).
		Children(func() { ui.Text(c, label).Bold().SingleLine() })
}

// resting reports e still and unheld after a drag.
func resting(e ui.Element) bool {
	_, _, held := e.Dragged()
	return !held && !e.Dragging() && !e.Transitioning()
}

// side runs an overlay panel with a backdrop while open.
func side(c *ui.Context, open *bool, label string, place func(back ui.Element) ui.Element, from ui.Pose, body func()) (closed bool) {
	if !*open {
		return false
	}
	ui.Overlay(c, func() {
		back := ui.Row(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Modal()
		kit.Arrive(back.Background(kit.Shade(0.3)), ui.Pose{Fade: 1}, 0)
		back.Children(func() {
			panel := place(back).Background(c.Theme().Background).Role(ui.RoleDialog).Label(label).Shadow(0, 12, 32, 0, kit.Shade(0.24))
			kit.Arrive(panel, from, 0).Children(body)
			if panel.PressedOutside() || back.OverlayShortcut(0, ui.KeyEscape) {
				*open, closed = false, true
			}
		})
	})
	return closed
}

// quiet is a ghost button; its caller cues.
func quiet(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(2), t.Space(3)).Radius(t.Radius))
	return b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).
		Children(func() { ui.Text(c, label).SingleLine() })
}

// entries lists menu items; pick runs for a click.
func entries(c *ui.Context, items []MenuItem, on int, pick func(i int)) {
	t := c.Theme()
	ui.Column(c).Width(t.Space(60)).Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Children(func() {
		for i, it := range items {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(it.Label).FillWidth().Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius).Checked(i == on))
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if i == on {
				bg = t.Accent.Alpha(0.12)
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					if it.Icon != "" {
						ui.Icon(c, icon.Must(it.Icon)).Size(16, 16)
					}
					ui.Text(c, it.Label).Grow(1).SingleLine()
					if it.Keys != "" {
						kit.Muted(c, it.Keys).FontSize(t.Rem(0.8))
					}
				})
			})
			if b.Clicked() {
				pick(i)
			}
		}
	})
}

// items checks menu items and a chosen index or -1.
func items(component, label string, list []MenuItem, chosen *int) {
	if label == "" || len(list) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(list) {
		panic(fmt.Sprintf("overlay: %s %q needs items and a choice or -1", component, label))
	}
	for _, it := range list {
		if it.Label == "" {
			panic(fmt.Sprintf("overlay: %s %q items need labels", component, label))
		}
		if it.Icon != "" {
			icon.Must(it.Icon)
		}
	}
}

// lighter returns the lighter of two theme colors.
func lighter(a, b ui.Color) ui.Color {
	if int(a.R)+int(a.G)+int(a.B) > int(b.R)+int(b.G)+int(b.B) {
		return a
	}
	return b
}

// def sets an empty copy field to its English default.
func def(field *string, english string) {
	if *field == "" {
		*field = english
	}
}

// valid reports whether i indexes n items.
func valid(i, n int) bool { return i >= 0 && i < n }

func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}
