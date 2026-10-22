package form

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Question is one survey question and its choices.
type Question struct {
	Text    string
	Options []string
}

// SurveyOptions tune a Survey.
type SurveyOptions struct {
	// Progress formats "%d of %d"; Back and Done are copy.
	Progress, Back, Done string
}

// Survey asks one question at a time; -1 means unanswered.
func Survey(c *ui.Context, questions []Question, answers *[]int, o SurveyOptions) (root, back ui.Element) {
	if len(questions) == 0 || answers == nil || len(*answers) != len(questions) {
		panic("form: Survey needs questions and one answer slot each")
	}
	for i, q := range questions {
		if q.Text == "" || len(q.Options) == 0 || (*answers)[i] < -1 || (*answers)[i] >= len(q.Options) {
			panic(fmt.Sprintf("form: Survey question %d is malformed", i))
		}
	}
	def(&o.Progress, "%d of %d")
	def(&o.Back, "Back")
	def(&o.Done, "Thanks for your answers.")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(3))
	at := ui.Local(root, "at", func() int { return 0 })
	dir := ui.Local(root, "dir", func() float32 { return 1 })
	root.Children(func() {
		i, shift := *at, *dir*t.Space(8)
		if i == len(questions) {
			kit.Arrive(ui.Text(c, o.Done).FontWeight(600), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			return
		}
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			kit.Muted(c, fmt.Sprintf(o.Progress, i+1, len(questions))).Grow(1)
			if i == 0 {
				return
			}
			if back = button.Ghost(c, o.Back); back.Clicked() {
				*at, *dir = i-1, -1
				kit.Log("Survey", "back", "question", i-1)
			}
		})
		ui.Box(c).FillWidth().Clip().Children(func() {
			q := ui.Column(c).Key(i).FillWidth().Gap(t.Space(2)).
				Transition(ui.ElementTransition{Enter: &ui.Motion{X: shift}, Exit: &ui.Motion{X: -shift}})
			q.Children(func() {
				ui.Text(c, questions[i].Text).FontWeight(600)
				for j, opt := range questions[i].Options {
					if answer(c, opt, (*answers)[i] == j).Clicked() {
						(*answers)[i], *at, *dir = j, i+1, 1
						sound.Play(sound.Select)
						root.SetChanged()
						kit.Log("Survey", "answer", "question", i, "option", j)
					}
				}
			})
		})
	})
	return root, back
}

// answer is a survey option row tinted when chosen.
func answer(c *ui.Context, label string, on bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).FillWidth().Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius))
	border, bg := t.Border, t.Background
	if on {
		border, bg = t.Accent, t.Accent.Alpha(0.08)
	} else if b.Hovered() {
		bg = t.Surface
	}
	return b.Checked(on).Border(1, b.AnimateColor("border", border, ui.Snappy)).
		Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() { ui.Text(c, label) })
}
