package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// sectioned builds titled sections in a scroll that st tracks.
func sectioned(c *ui.Context, st *ui.ScrollState, titles []string, section func(i int)) (tops func() []float32) {
	t := c.Theme()
	var secs []ui.Element
	ui.Scroll(c).Grow(1).TrackScroll(st).Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(6)).Padding(t.Space(3)).Children(func() {
			for i, title := range titles {
				c.Key(i)
				s := ui.Column(c).FillWidth().Gap(t.Space(2))
				s.Children(func() {
					ui.Text(c, title).Bold().FontSize(t.Rem(1.1))
					section(i)
				})
				secs = append(secs, s)
			}
		})
	})
	return func() []float32 {
		out := make([]float32, len(secs))
		for i, s := range secs {
			out[i] = s.LayoutBounds().Y - secs[0].LayoutBounds().Y
		}
		return out
	}
}

// opener is a menu trigger whose chevron turns while open.
func opener(c *ui.Context, label string, open, link bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Gap(t.Space(1)).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)).Expanded(open)
	bg, fg := t.SurfaceHover.Alpha(0), t.Text
	switch {
	case link:
		fg = t.Accent
	case open || b.Hovered():
		bg = t.SurfaceHover
	}
	return b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(fg).Children(func() {
		ui.Text(c, label).SingleLine()
		chev := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16)
		chev.Rotate(180 * chev.AnimateTo("turn", unit(open), ui.Snappy))
	})
}

// toggled flips open on a click, cueing Open.
func toggled(b ui.Element, open *bool) {
	if b.Clicked() {
		if *open = !*open; *open {
			sound.Play(sound.Open)
		}
	}
}

// reopened marks root changed when the open menu moved.
func reopened[T comparable](root ui.Element, was, now T, component string) {
	if was != now {
		root.SetChanged()
		kit.Log(component, "open", now)
	}
}

// popup opens body in an arriving panel under trigger.
func popup(c *ui.Context, trigger ui.Element, open *bool, body func()) {
	t := c.Theme()
	toggled(trigger, open)
	ui.Popover(c, trigger, open, func() {
		kit.Arrive(ui.Column(c).Gap(t.Space(0.5)).Padding(t.Space(1)), ui.Pose{Scale: 0.96, Fade: 1}, 0).Children(body)
	})
}

// checked is a menu row with a popping check.
func checked(c *ui.Context, label string, on bool, lead func()) ui.Element {
	t := c.Theme()
	return rowBase(c, label).Checked(on).Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			if lead != nil {
				lead()
			}
			ui.Text(c, label).Grow(1).SingleLine()
			mark := ui.Box(c).Size(16, 16).Shrink(0)
			if on {
				mark.Children(func() {
					kit.Arrive(ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Accent), ui.Pose{Scale: 0.5, Fade: 1}, 0)
				})
			}
		})
	})
}
