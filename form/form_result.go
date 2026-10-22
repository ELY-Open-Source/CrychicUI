package form

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/egoist/mygo/ui"
)

// Outcome is where a submitted form stands.
type Outcome int

// Outcomes of a submission.
const (
	Pending Outcome = iota
	Succeeded
	Failed
)

// FormResultOptions tune a FormResult.
type FormResultOptions struct {
	// Copy for the success and error panels.
	SuccessTitle, SuccessText, Done, ErrorTitle, ErrorText, Retry string
}

// FormResult swaps the form for a success or error panel.
func FormResult(c *ui.Context, outcome Outcome, o FormResultOptions, content func()) (root, action ui.Element) {
	if content == nil || outcome < Pending || outcome > Failed {
		panic(fmt.Sprintf("form: FormResult needs content and a known outcome, not %d", outcome))
	}
	def(&o.SuccessTitle, "Sent")
	def(&o.SuccessText, "We got it. Thank you.")
	def(&o.Done, "Done")
	def(&o.ErrorText, "We could not send it. Try again.")
	swap := ui.ElementTransition{Exit: &ui.Motion{Opacity: 0}}
	root = ui.Box(c).FillWidth()
	root.Children(func() {
		panel := ui.Column(c).Key(outcome).FillWidth().Transition(swap)
		panel.Children(func() {
			switch outcome {
			case Succeeded:
				_, action = feedback.SuccessState(c, o.SuccessTitle, o.SuccessText, o.Done)
			case Failed:
				_, action = feedback.ErrorState(c, o.ErrorText, feedback.ErrorStateOptions{Title: o.ErrorTitle, Retry: o.Retry})
			default:
				content()
			}
		})
	})
	return root, action
}
