package learn

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QuizResultOptions word the result; Pass is the passing share.
type QuizResultOptions struct {
	Title, Score, Passed, Failed, Answer, Chose, Skipped string
	Pass                                                 float64
}

// QuizResult shows a score ring and each answer reviewed.
func QuizResult(c *ui.Context, qs []Question, picks []int, o QuizResultOptions) ui.Element {
	def(&o.Title, "Quiz results")
	if len(qs) == 0 || len(picks) != len(qs) || o.Pass < 0 || o.Pass > 1 || o.Pass != o.Pass {
		panic(fmt.Sprintf("learn: QuizResult %q needs questions, a pick each and a pass share in 0..1", o.Title))
	}
	texts := make([]string, len(qs))
	right := 0
	for i, q := range qs {
		validQuestion("QuizResult", q)
		if picks[i] < -1 || picks[i] >= len(q.Choices) {
			panic(fmt.Sprintf("learn: QuizResult %q pick %d is out of range", q.Text, picks[i]))
		}
		texts[i] = q.Text
		if picks[i] == q.Answer {
			right++
		}
	}
	unique("QuizResult", texts)
	def(&o.Score, "%d of %d correct")
	def(&o.Passed, "Passed")
	def(&o.Failed, "Keep practicing")
	def(&o.Answer, "Answer: %s")
	def(&o.Chose, "You chose %s")
	def(&o.Skipped, "Skipped")
	if o.Pass == 0 {
		o.Pass = 0.7
	}
	t := c.Theme()
	share := float64(right) / float64(len(qs))
	pass := share >= o.Pass
	ink, verdict := t.Warning, o.Failed
	if pass {
		ink, verdict = t.Success, o.Passed
	}
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).Children(func() {
			dial(c, t.Space(24), float32(share), ink, func() {
				ui.Text(c, fmt.Sprintf("%d%%", percent(right, len(qs)))).FontSize(t.Rem(1.2)).Bold().Shrink(0)
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
				ui.Text(c, o.Title).FontSize(t.Rem(1)).Bold().FillWidth().MinWidth(0).SingleLine()
				kit.Muted(c, fmt.Sprintf(o.Score, right, len(qs))).FontSize(t.Rem(0.8)).FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).Children(func() {
					kit.Arrive(kit.Chip(c, verdict, ink), ui.Pose{Scale: 0.7, Fade: 1}, 500*time.Millisecond)
				})
			})
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(0, t.Space(3), t.Space(3), t.Space(3))
		list.Children(func() {
			for i, q := range qs {
				good := picks[i] == q.Answer
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start).Padding(t.Space(2.5)).Radius(t.Radius).
					BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
				kit.Stagger(c, list, row, i).Children(func() {
					mark, mk := "circle-x", t.Danger
					if good {
						mark, mk = "circle-check", t.Success
					}
					glyph(c, mark, 18).TextColor(mk)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Text(c, q.Text).FontSize(t.Rem(0.82)).Bold().FillWidth().Wrap()
						if good {
							return
						}
						chose := o.Skipped
						if picks[i] >= 0 {
							chose = fmt.Sprintf(o.Chose, q.Choices[picks[i]])
						}
						ui.Text(c, chose).FontSize(t.Rem(0.74)).TextColor(t.Danger).FillWidth().Wrap()
						ui.Text(c, fmt.Sprintf(o.Answer, q.Choices[q.Answer])).FontSize(t.Rem(0.74)).Bold().TextColor(t.Success).FillWidth().Wrap()
					})
				})
			}
		})
	})
	return root
}
