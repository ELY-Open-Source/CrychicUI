package learn

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Said is a true-or-false answer.
type Said int

const (
	Unsaid Said = iota
	SaidTrue
	SaidFalse
)

// TrueFalseOptions word the two buttons, verdicts and why.
type TrueFalseOptions struct{ True, False, Right, Wrong, Explain string }

// TrueFalse answers a statement once; cues Success, Error.
func TrueFalse(c *ui.Context, statement string, truth bool, said *Said, o TrueFalseOptions) ui.Element {
	if said == nil || strings.TrimSpace(statement) == "" || *said < Unsaid || *said > SaidFalse {
		panic(fmt.Sprintf("learn: TrueFalse %q needs a statement and an answer state", statement))
	}
	def(&o.True, "True")
	def(&o.False, "False")
	def(&o.Right, "Correct")
	def(&o.Wrong, "Not quite")
	t := c.Theme()
	answer := SaidFalse
	if truth {
		answer = SaidTrue
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
			tile(c, "message-circle-question", t.Accent, t.Space(9))
			ui.Text(c, statement).FontSize(t.Rem(0.95)).Bold().Grow(1).MinWidth(0).Wrap()
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for _, s := range []Said{SaidTrue, SaidFalse} {
				label, mark := o.True, "check"
				if s == SaidFalse {
					label, mark = o.False, "x"
				}
				on := *said == s
				ink := t.Accent
				switch {
				case *said != Unsaid && s == answer:
					ink = t.Success
				case on:
					ink = t.Danger
				}
				lit := *said != Unsaid && (on || s == answer)
				c.Key(s)
				b := kit.Press(ui.ButtonBase(c).Label(label).Grow(1).Basis(0).MinWidth(0).Height(t.Space(16)).Column().Center().
					Gap(t.Space(1)).Radius(t.Radius * 1.2).Checked(on).Disabled(*said != Unsaid))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(0.6).Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.14), unit(lit)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(lit)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.Text.Mix(ink, unit(lit)), ui.Snappy))
				b.Children(func() {
					g := glyph(c, mark, 22)
					k := g.AnimateTo("pop", 1+0.2*unit(on), ui.Bouncy)
					g.Scale(k, k)
					ui.Text(c, label).FontSize(t.Rem(0.88)).Bold().MaxWidth(t.Space(30)).SingleLine()
				})
				if b.Clicked() && *said == Unsaid {
					*said = s
					root.SetChanged()
					if s == answer {
						sound.Play(sound.Success)
					} else {
						sound.Play(sound.Error)
					}
					kit.Log("TrueFalse", "answered", "right", s == answer)
				}
			}
		})
		if *said != Unsaid {
			ink, mark, word := t.Danger, "circle-x", o.Wrong
			if *said == answer {
				ink, mark, word = t.Success, "circle-check", o.Right
			}
			c.Key("verdict")
			note := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).Background(ink.Alpha(0.1))
			kit.Arrive(note, ui.Pose{Y: t.Space(1.5), Fade: 1}, 0).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(ink).Children(func() {
					glyph(c, mark, 16)
					ui.Text(c, word).FontSize(t.Rem(0.85)).Bold().MinWidth(0).SingleLine()
				})
				if o.Explain != "" {
					ui.Text(c, o.Explain).FontSize(t.Rem(0.78)).TextColor(t.Text.Mix(t.TextMuted, 0.3)).FillWidth().Wrap()
				}
			})
		}
	})
	return root
}
