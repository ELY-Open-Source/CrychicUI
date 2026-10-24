package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QuestionData is one question with votes and answers.
type QuestionData struct {
	ID, Title, Body string
	Votes, Answers  int
}

// QuestionOptions are the question labels.
type QuestionOptions struct{ Answers, Up, Down string }

// Question shows a question with votes and answers; cue Tick.
func Question(c *ui.Context, question *QuestionData, vote *int, o QuestionOptions) ui.Element {
	if question == nil || question.ID == "" || question.Title == "" || question.Body == "" || question.Answers < 0 || vote == nil {
		panic("social: invalid question")
	}
	def(&o.Answers, "%d answers")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, question.Title).Bold()
		ui.Text(c, question.Body).Wrap()
		ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			votes := VoteArrowsOptions{Up: o.Up, Down: o.Down}
			if VoteArrows(c, &question.Votes, vote, votes).Changed() {
				root.SetChanged()
			}
			kit.Numeric(c, fmt.Sprintf(o.Answers, question.Answers)).TextColor(t.TextMuted)
		})
	})
	return root
}
