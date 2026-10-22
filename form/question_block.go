package form

import (
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QuestionBlockOptions tune a QuestionBlock.
type QuestionBlockOptions struct {
	// Hint and Placeholder are copy; Max is 280 when zero.
	Hint, Placeholder string
	Max               int
}

// QuestionBlock asks a prompt with a hint and counted answer.
func QuestionBlock(c *ui.Context, prompt string, answer *string, o QuestionBlockOptions) ui.Element {
	if prompt == "" || answer == nil || o.Max < 0 {
		panic("form: QuestionBlock needs a prompt, an answer and a limit of 0 or more")
	}
	if o.Max == 0 {
		o.Max = 280
	}
	t := c.Theme()
	root := kit.Arrive(ui.Column(c).FillWidth().Gap(t.Space(1)), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		edited(root, input.CountedField(c, prompt, answer, input.CountedFieldOptions{Placeholder: o.Placeholder, Max: o.Max}))
		if o.Hint != "" {
			kit.Muted(c, o.Hint).FontSize(t.Rem(0.85))
		}
	})
}
