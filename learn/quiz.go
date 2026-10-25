package learn

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Question is a quiz question; Answer indexes Choices.
type Question struct {
	Text    string
	Choices []string
	Answer  int
	Explain string
}

func validQuestion(component string, q Question) {
	if strings.TrimSpace(q.Text) == "" || len(q.Choices) < 2 || len(q.Choices) > 6 || q.Answer < 0 || q.Answer >= len(q.Choices) {
		panic(fmt.Sprintf("learn: %s %q needs 2 to 6 choices and an answer among them", component, q.Text))
	}
	unique(component, q.Choices)
}

// QuizOptions word the check button and verdicts.
type QuizOptions struct{ Check, Right, Wrong string }

// Quiz picks a choice, then checks it; cues Success, Error.
func Quiz(c *ui.Context, q Question, picked *int, checked *bool, o QuizOptions) ui.Element {
	validQuestion("Quiz", q)
	if picked == nil || checked == nil || *picked < -1 || *picked >= len(q.Choices) || *checked && *picked < 0 {
		panic(fmt.Sprintf("learn: Quiz %q needs a pick in range, and one before checking", q.Text))
	}
	def(&o.Check, "Check")
	def(&o.Right, "Correct")
	def(&o.Wrong, "Not quite")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, q.Text).FontSize(t.Rem(1)).Bold().FillWidth().Wrap()
		list := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleRadioGroup)
		list.Children(func() {
			for i, ch := range q.Choices {
				on := *picked == i
				right, wrong := *checked && i == q.Answer, *checked && on && i != q.Answer
				ink := t.Accent
				switch {
				case right:
					ink = t.Success
				case wrong:
					ink = t.Danger
				}
				c.Key(strconv.Quote(ch))
				b := kit.Press(ui.ButtonBase(c).Label(ch).Role(ui.RoleRadio).Checked(on).FillWidth().Padding(t.Space(2.5), t.Space(3)).
					Radius(t.Radius).Disabled(*checked))
				lit := on || right
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.1), unit(lit)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(lit)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						badge := ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999)
						badge.Background(badge.AnimateColor("bg", t.SurfaceHover.Mix(ink, unit(lit)), ui.Snappy))
						badge.TextColor(badge.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(lit)), ui.Snappy))
						badge.Children(func() {
							c.Key([2]bool{right, wrong})
							switch {
							case right:
								kit.Arrive(glyph(c, "check", 14), ui.Pose{Scale: 0.4, Fade: 1}, 0)
							case wrong:
								kit.Arrive(glyph(c, "x", 14), ui.Pose{Scale: 0.4, Fade: 1}, 0)
							default:
								ui.Text(c, string(rune('A'+i))).FontSize(t.Rem(0.72)).Bold().Shrink(0)
							}
						})
						ui.Text(c, ch).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
					})
				})
				if b.Clicked() && !*checked && !on {
					*picked = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Quiz", "picked", "choice", i)
				}
			}
		})
		if *checked {
			good := *picked == q.Answer
			ink, mark, word := t.Danger, "circle-x", o.Wrong
			if good {
				ink, mark, word = t.Success, "circle-check", o.Right
			}
			c.Key("verdict")
			note := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).Background(ink.Alpha(0.1))
			kit.Arrive(note, ui.Pose{Y: t.Space(1.5), Fade: 1}, 0).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(ink).Children(func() {
					glyph(c, mark, 16)
					ui.Text(c, word).FontSize(t.Rem(0.85)).Bold().MinWidth(0).SingleLine()
				})
				if q.Explain != "" {
					ui.Text(c, q.Explain).FontSize(t.Rem(0.78)).TextColor(t.Text.Mix(t.TextMuted, 0.3)).FillWidth().Wrap()
				}
			})
		}
		c.Key("check")
		b := kit.Press(ui.ButtonBase(c).Label(o.Check).FillWidth().Padding(t.Space(2.5), 0).Radius(999).Disabled(*picked < 0 || *checked))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
		b.Children(func() { ui.Text(c, o.Check).FontSize(t.Rem(0.85)).Bold().Shrink(0) })
		if b.Clicked() && *picked >= 0 && !*checked {
			*checked = true
			root.SetChanged()
			if *picked == q.Answer {
				sound.Play(sound.Success)
			} else {
				sound.Play(sound.Error)
			}
			kit.Log("Quiz", "checked", "right", *picked == q.Answer)
		}
	})
	return root
}
