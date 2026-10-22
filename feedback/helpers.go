// Package feedback holds alerts, progress, loaders and states.
package feedback

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var quick = ui.Animation{Duration: 120 * time.Millisecond}

// look is a tone's color and Lucide icon.
func look(t *ui.Theme, tone Tone) (ui.Color, string) {
	switch tone {
	case ToneInfo:
		return t.Accent, "info"
	case ToneSuccess:
		return t.Success, "circle-check"
	case ToneWarning:
		return t.Warning, "triangle-alert"
	case ToneDanger:
		return t.Danger, "circle-x"
	}
	panic("feedback: unknown tone")
}

// unit converts on to an animation target.
func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// fraction panics unless v is within 0 and 1.
func fraction(component string, v float64) float32 {
	if !(v >= 0 && v <= 1) {
		panic("feedback: " + component + " needs a value from 0 to 1")
	}
	return float32(v)
}

// button is a quiet action button that presses in.
func button(c *ui.Context, label string, primary bool) ui.Element {
	t := c.Theme()
	b := ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)
	bg, fg := t.Surface, t.Text
	switch {
	case primary && b.Pressed():
		bg, fg = t.AccentPressed, t.AccentText
	case primary && b.Hovered():
		bg, fg = t.AccentHover, t.AccentText
	case primary:
		bg, fg = t.Accent, t.AccentText
	case b.Pressed():
		bg = t.SurfacePressed
	case b.Hovered():
		bg = t.SurfaceHover
	}
	b.Background(b.AnimateColor("bg", bg, quick)).TextColor(fg)
	if !primary {
		b.Border(1, t.Border)
	}
	return kit.Press(b).Children(func() { ui.Text(c, label).SingleLine() })
}

// closer is a small cross button.
func closer(c *ui.Context, label string) ui.Element { return cross(c, ui.ButtonBase(c).Label(label)) }

// cross styles b as a small cross button.
func cross(c *ui.Context, b ui.Element) ui.Element {
	t := c.Theme()
	b.Size(t.Space(7), t.Space(7)).Radius(t.Radius).TextColor(t.TextMuted)
	bg := ui.Color{}
	if b.Hovered() {
		bg = t.Text.Alpha(0.08)
	}
	b.Background(b.AnimateColor("bg", bg, quick))
	return kit.Press(b).Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
}

// once reports true on e's first build only.
func once(e ui.Element) bool {
	seen := ui.Local(e, "crychic.seen", func() bool { return false })
	first := !*seen
	*seen = true
	return first
}

// sheen paints a skeleton block with a band sweeping by.
func sheen(c *ui.Context, e ui.Element, radius float32) ui.Element {
	t := c.Theme()
	base := t.Text.Alpha(0.07)
	shine := t.Text.Alpha(0.13)
	return e.Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		w, view := t.Space(60), p.Visible()
		phase := float32(p.Now().UnixMilli()%1600) / 1600
		x := view.X - w + (view.W+2*w)*phase
		p.Fill(r, base, radius)
		p.Clip(r, radius, func() {
			p.FillGradient(ui.Rect{X: x, Y: r.Y, W: w, H: r.H},
				ui.LinearGradient{Angle: 90, Stops: []ui.GradientStop{{Color: base, At: 0}, {Color: shine, At: 0.5}, {Color: base, At: 1}}}, 0)
		})
	})
}

// polar is the point at angle degrees clockwise from right.
func polar(cx, cy, r float32, deg float64) (float32, float32) {
	a := deg * math.Pi / 180
	return cx + r*float32(math.Cos(a)), cy + r*float32(math.Sin(a))
}
