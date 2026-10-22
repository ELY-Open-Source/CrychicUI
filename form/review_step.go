package form

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ReviewSection is a titled group of answers to review.
type ReviewSection struct {
	Title string
	Items []display.KeyValuePair
}

// ReviewStepOptions tune a ReviewStep.
type ReviewStepOptions struct {
	// Edit is the link text; EditFor formats its accessible name.
	Edit, EditFor string
}

// ReviewStep summarizes sections; Edit sets *edit to that section.
func ReviewStep(c *ui.Context, sections []ReviewSection, edit *int, o ReviewStepOptions) ui.Element {
	if len(sections) == 0 || edit == nil {
		panic("form: ReviewStep needs sections and an edit index")
	}
	def(&o.Edit, "Edit")
	def(&o.EditFor, "Edit %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(5))
	return root.Children(func() {
		for i, s := range sections {
			ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
					ui.Text(c, s.Title).FontWeight(600).Grow(1)
					if button.Ghost(c, o.Edit).Label(fmt.Sprintf(o.EditFor, s.Title)).Clicked() {
						*edit = i
						root.SetChanged()
						kit.Log("ReviewStep", "edit", "section", i)
					}
				})
				display.DescriptionList(c, s.Items)
			})
		}
	})
}
