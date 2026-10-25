// Package account holds account and billing components.
package account

import (
	"fmt"
	"math"
	"net/mail"
	"strings"
	"time"

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

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func validEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

// glyph is a Lucide icon size DIP square.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius*1.5).Border(1, t.Border).Background(t.Background).
		Shadow(0, 1, 3, 0, t.Text.Alpha(0.06))
	return kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// tile is a tinted round icon square.
func tile(c *ui.Context, mark string, ink ui.Color, size float32) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink)
	return box.Children(func() { glyph(c, mark, size*0.5) })
}

// head is a form header: icon tile, title, optional subtitle.
func head(c *ui.Context, mark, title, sub string) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
		kit.Arrive(tile(c, mark, t.Accent, t.Space(10)), ui.Pose{Scale: 0.6, Fade: 1}, 0)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, title).FontSize(t.Rem(1.05)).Bold().Wrap().FillWidth()
			if sub != "" {
				kit.Muted(c, sub).FontSize(t.Rem(0.8)).Wrap().FillWidth()
			}
		})
	})
}

// field is a framed input built before its caption.
func field(c *ui.Context, caption, mark string, value *string, password bool) ui.Element {
	t := c.Theme()
	var in ui.Element
	ui.Grid(c).Columns(1).FillWidth().Gap(t.Space(1)).Children(func() {
		frame := ui.Row(c).RowStart(2).ColumnStart(1).FillWidth().MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).
			Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background)
		frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
		frame.Children(func() {
			glyph(c, mark, 16).TextColor(t.TextMuted)
			in = ui.TextInputBase(c, value).Label(caption).Placeholder(caption).FontSize(t.Rem(0.88)).Grow(1).MinWidth(0)
			if password {
				in.Password()
			}
		})
		ui.Text(c, caption).RowStart(1).ColumnStart(1).FontSize(t.Rem(0.78)).Bold().TextColor(t.TextMuted)
	})
	return in
}

// primary is the accent action with an optional icon.
func primary(c *ui.Context, label, mark string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(2), t.Space(3.5)).Radius(t.Radius))
	b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
	return b.Children(func() { content(c, label, mark) })
}

// outline is the bordered secondary action.
func outline(c *ui.Context, label, text, mark string, ink ui.Color) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Shrink(0))
	b.Background(b.AnimateColor("bg", t.Background.Mix(ink.Alpha(0.1), unit(b.Hovered())), ui.Snappy))
	b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(b.Hovered())*0.6), ui.Snappy)).TextColor(ink)
	return b.Children(func() { content(c, text, mark) })
}

func content(c *ui.Context, text, mark string) {
	t := c.Theme()
	ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
		if mark != "" {
			glyph(c, mark, 16)
		}
		ui.Text(c, text).FontSize(t.Rem(0.85)).Bold().Shrink(0)
	})
}

// note is a tinted status line that pops in.
func note(c *ui.Context, mark, text string, ink ui.Color) ui.Element {
	t := c.Theme()
	row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).
		Background(ink.Alpha(0.1)).TextColor(ink)
	return kit.Arrive(row, ui.Pose{Y: -t.Space(1), Fade: 1}, 0).Children(func() {
		glyph(c, mark, 16)
		ui.Text(c, text).FontSize(t.Rem(0.82)).Bold().Grow(1).MinWidth(0).Wrap()
	})
}

// empty is a centered empty-list state.
func empty(c *ui.Context, mark, text string) {
	t := c.Theme()
	col := ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(6), t.Space(3))
	kit.Arrive(col, ui.Pose{Scale: 0.9, Fade: 1}, 0).Children(func() {
		kit.Arrive(tile(c, mark, t.TextMuted, t.Space(10)), ui.Pose{Scale: 0.5, Fade: 1}, 120*time.Millisecond)
		kit.Muted(c, text).FontSize(t.Rem(0.85))
	})
}

// meter is a bar filled to share in ink.
func meter(c *ui.Context, share float32, ink ui.Color) ui.Element {
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).Background(t.SurfaceHover)
	grow := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	w := grow * bar.AnimateTo("share", max(0, min(1, share)), ui.Smooth)
	col := bar.AnimateColor("ink", ink, ui.Snappy)
	return bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, col, 999)
	}))
}

// dot joins two meta texts with a drawn dot.
func dot(c *ui.Context) {
	t := c.Theme()
	ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.TextMuted.Alpha(0.6))
}

// unique panics when items repeat or are blank.
func unique(component string, items []string) {
	seen := map[string]bool{}
	for _, s := range items {
		if strings.TrimSpace(s) == "" || seen[s] {
			panic(fmt.Sprintf("account: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// rowExit fades a removed row out.
var rowExit = ui.ElementTransition{Duration: 160 * time.Millisecond, Exit: &ui.Motion{}}
