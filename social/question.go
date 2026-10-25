package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QuestionData is one question; Tags are optional.
type QuestionData struct {
	ID, Title, Body string
	Votes, Answers  int
	Tags            []string
}

// QuestionOptions word the vote buttons; Answer, Answers take the count.
type QuestionOptions struct{ Answer, Answers, Up, Down string }

// Question shows a question with votes and answers; cue Tick.
func Question(c *ui.Context, question *QuestionData, vote *int, o QuestionOptions) ui.Element {
	if question == nil {
		panic("social: Question needs a question")
	}
	if question.ID == "" || question.Title == "" || question.Body == "" || question.Answers < 0 {
		panic(fmt.Sprintf("social: Question %q needs id, body and answers", question.Title))
	}
	if vote == nil || *vote < -1 || *vote > 1 {
		panic(fmt.Sprintf("social: Question %q needs a vote of -1, 0 or 1", question.Title))
	}
	unique("Question", question.Tags)
	def(&o.Answer, "%s answer")
	def(&o.Answers, "%s answers")
	answers := o.Answers
	if question.Answers == 1 {
		answers = o.Answer
	}
	t := c.Theme()
	root := card(c).Row().Gap(t.Space(3)).AlignItems(ui.Start).Padding(t.Space(3))
	root.Children(func() {
		votes := VoteArrowsOptions{Up: o.Up, Down: o.Down}
		if VoteArrows(c, &question.Votes, vote, votes).Changed() {
			root.SetChanged()
		}
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
			ui.Text(c, question.Title).FontSize(t.Rem(0.95)).Bold().Wrap().FillWidth()
			kit.Muted(c, question.Body).FontSize(t.Rem(0.82)).Wrap().FillWidth()
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(0.5), 0, 0, 0).Children(func() {
				ink := t.TextMuted
				if question.Answers > 0 {
					ink = t.Success
				}
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).
					Background(ink.Alpha(0.12)).TextColor(ink).Children(func() {
					glyph(c, "message-square", 12)
					kit.Numeric(c, fmt.Sprintf(answers, number(c, question.Answers))).FontSize(t.Rem(0.72)).Bold().Shrink(0)
				})
				for _, tag := range question.Tags {
					c.Key(tag)
					kit.Chip(c, tag, t.Accent)
				}
			})
		})
	})
	return root
}
