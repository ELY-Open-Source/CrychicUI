package form

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ErrorSummaryOptions tune an ErrorSummary.
type ErrorSummaryOptions struct {
	// Title formats the error count in.
	Title string
}

// ErrorSummary counts the form's errors; each focuses its field.
func ErrorSummary(c *ui.Context, errors []feedback.FormError, o ErrorSummaryOptions) ui.Element {
	def(&o.Title, "Fix %d problems")
	root := ui.Column(c).FillWidth()
	shown := ui.Local(root, "count", func() int { return 0 })
	if n := len(errors); n > *shown {
		sound.Play(sound.Error)
		kit.Log("ErrorSummary", "errors", "count", n)
	}
	*shown = len(errors)
	return root.Children(func() {
		if len(errors) > 0 {
			ui.Box(c).FillWidth().Transition(ui.ElementTransition{Exit: &ui.Motion{}}).Children(func() {
				feedback.ErrorList(c, errors, feedback.ErrorListOptions{Title: fmt.Sprintf(o.Title, len(errors))})
			})
		}
	})
}
