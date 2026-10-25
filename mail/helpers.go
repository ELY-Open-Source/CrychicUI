// Package mail holds mail and communication components.
package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Mail is one message summary and body; Attached marks files.
type Mail struct {
	ID, Sender, Subject, Preview, Body, Time string
	Unread, Attached                         bool
}

// Contact is a person; Phone and Role are optional.
type Contact struct {
	ID, Name, Email, Phone, Role string
}

// Participant is one call participant.
type Participant struct {
	ID, Name        string
	Muted, Speaking bool
}

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

func validMail(component string, m Mail) {
	if m.ID == "" || m.Sender == "" || m.Subject == "" || m.Time == "" {
		panic(fmt.Sprintf("mail: %s message %q needs id, sender, subject, and time", component, m.ID))
	}
}

func validContact(component string, p Contact) {
	if p.ID == "" || p.Name == "" || p.Email == "" {
		panic(fmt.Sprintf("mail: %s contact %q needs id, name, and email", component, p.ID))
	}
}

// unique panics when names repeat or are empty.
func unique(component string, names []string) {
	seen := map[string]bool{}
	for _, name := range names {
		if name == "" || seen[name] {
			panic(fmt.Sprintf("mail: %s has an empty or repeated %q", component, name))
		}
		seen[name] = true
	}
}

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	return kit.Arrive(kit.Card(c).FillWidth().Role(ui.RoleGroup), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// glyph is a small icon in the current text color.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// tile is a tinted rounded icon square.
func tile(c *ui.Context, name string, ink ui.Color) ui.Element {
	t := c.Theme()
	return ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).
		TextColor(ink).Children(func() { glyph(c, name, 16) })
}

// avatar is a person's circle, units theme spaces wide.
func avatar(c *ui.Context, name string, units float32) ui.Element {
	t := c.Theme()
	return display.Avatar(c, name, nil).Size(t.Space(units), t.Space(units)).FontSize(t.Rem(units * 0.08)).Shrink(0).AlignSelf(ui.Center)
}

// rowInk is a row's animated hover and selected tint.
func rowInk(t *ui.Theme, e ui.Element, on bool) ui.Element {
	return e.Background(e.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(e.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
}

// row is a full-width pressable list row.
func row(c *ui.Context, label string, on bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(2), t.Space(2.5)).
		Radius(t.Radius).Checked(on))
	return rowInk(t, b, on)
}

// toggle is a labeled switch with a sliding knob.
func toggle(c *ui.Context, label string, on bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(9), t.Space(5)).Shrink(0).Radius(999).Checked(on))
	return track(c, b, on, b.Hovered())
}

// track paints e as a switch with a sliding knob.
func track(c *ui.Context, e ui.Element, on, hover bool) ui.Element {
	t := c.Theme()
	e.Background(e.AnimateColor("track", t.SurfaceHover.Mix(t.Border, 0.6*unit(hover)).Mix(t.Accent, unit(on)), ui.Snappy))
	x := e.AnimateTo("knob", unit(on), ui.Snappy)
	pad := t.Space(0.5)
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		d := r.H - 2*pad
		p.Fill(ui.Rect{X: r.X + pad + x*(r.W-r.H), Y: r.Y + pad, W: d, H: d}, t.Background, d/2)
	})
}

// roundAction is a round icon button captioned below.
func roundAction(c *ui.Context, label, mark string, bg, ink ui.Color) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Column().AlignItems(ui.Center).Gap(t.Space(1.5)).Shrink(0))
	hover := unit(b.Hovered())
	return b.Children(func() {
		circle := ui.Box(c).Size(t.Space(12), t.Space(12)).Center().Radius(999).TextColor(ink)
		circle.Background(circle.AnimateColor("bg", bg.Mix(t.Text, 0.1*hover), ui.Snappy)).Children(func() {
			c.Key(mark)
			kit.Arrive(glyph(c, mark, 18), ui.Pose{Scale: 0.6, Fade: 1}, 0)
		})
		ui.Text(c, label).FontSize(t.Rem(0.72)).Shrink(0)
	})
}

// timeChoice is a row offering one time.
func timeChoice(c *ui.Context, label, mark, when string, on, off bool) ui.Element {
	t := c.Theme()
	b := row(c, label, on).Disabled(off)
	b.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			glyph(c, mark, 16).TextColor(t.Accent)
			ui.Text(c, label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
			kit.Muted(c, when).FontSize(t.Rem(0.75)).Shrink(0)
			c.Key(on)
			if on {
				kit.Arrive(glyph(c, "check", 14).TextColor(t.Accent), ui.Pose{Scale: 0.5, Fade: 1}, 0)
			}
		})
	})
	return b
}

func bytes(n int64) string {
	if n < 0 {
		panic("mail: attachment size must not be negative")
	}
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	value := float64(n)
	for _, suffix := range []string{"KB", "MB", "GB", "TB"} {
		value /= unit
		if value < unit {
			return fmt.Sprintf("%.1f %s", value, suffix)
		}
	}
	return fmt.Sprintf("%.1f PB", value/unit)
}
