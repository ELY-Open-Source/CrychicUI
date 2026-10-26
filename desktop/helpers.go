// Package desktop holds desktop shell and window components.
package desktop

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

// glyph is a Lucide icon size DIP square.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	return kit.Arrive(ui.Column(c).FillWidth().Radius(t.Radius*1.5).Border(1, t.Border).Background(t.Background).
		Shadow(0, 1, 3, 0, t.Text.Alpha(0.06)), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// window is a floating window surface that settles in.
func window(c *ui.Context) ui.Element {
	t := c.Theme()
	return kit.Arrive(ui.Column(c).FillWidth().Clip().Radius(t.Radius*2).Border(1, t.Border).Background(t.Background).
		Shadow(0, 8, 24, 0, kit.Shade(0.16)), ui.Pose{Scale: 0.97, Fade: 1}, 0)
}

// titleBar is a window bar: lights, title, tail.
func titleBar(c *ui.Context, title string, active bool, tail func()) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).
		Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
		lights(c, active)
		ink := t.TextMuted
		if active {
			ink = t.Text
		}
		ui.Text(c, title).FontSize(t.Rem(0.8)).Bold().TextColor(ink).Grow(1).MinWidth(0).TextAlign(ui.Center).SingleLine()
		ui.Row(c).Width(t.Space(13)).Shrink(0).Justify(ui.End).Children(func() {
			if tail != nil {
				tail()
			}
		})
	})
}

// lights are a window's three dots, grey when inactive.
func lights(c *ui.Context, active bool) {
	t := c.Theme()
	ui.Row(c).Width(t.Space(13)).Shrink(0).Gap(t.Space(1.5)).Children(func() {
		for _, ink := range [3]ui.Color{t.Danger, t.Warning, t.Success} {
			if !active {
				ink = t.Border
			}
			ui.Box(c).Size(t.Space(3), t.Space(3)).Radius(999).Background(ink)
		}
	})
}

// tile is a tinted icon square.
func tile(c *ui.Context, mark string, ink ui.Color, size float32) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink)
	return box.Children(func() { glyph(c, mark, size*0.5) })
}

// unique panics when items repeat or are blank.
func unique(component string, items []string) {
	seen := map[string]bool{}
	for _, s := range items {
		if strings.TrimSpace(s) == "" || seen[s] {
			panic(fmt.Sprintf("desktop: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// avatar is a photo or initials circle units wide.
func avatar(c *ui.Context, name string, photo *ui.Bitmap, units float32) ui.Element {
	t := c.Theme()
	return display.Avatar(c, name, photo).Size(t.Space(units), t.Space(units)).FontSize(t.Rem(units * 0.08)).Shrink(0).AlignSelf(ui.Center)
}

// share checks a 0 to 1 level.
func share(component, label string, v float32) {
	if math.IsNaN(float64(v)) || v < 0 || v > 1 {
		panic(fmt.Sprintf("desktop: %s %q needs a level 0 to 1, got %v", component, label, v))
	}
}

// minute wakes the view at the next whole minute.
func minute(c *ui.Context) time.Time {
	now := c.Now()
	c.After(now.Truncate(time.Minute).Add(time.Minute).Sub(now))
	return now
}
