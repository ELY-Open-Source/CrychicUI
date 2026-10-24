package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EpicOptions are the epic labels.
type EpicOptions struct{ Done string }

// Epic shows an epic with child progress.
func Epic(c *ui.Context, title string, children []Task, o EpicOptions) ui.Element {
	if title == "" || len(children) == 0 {
		panic("project: Epic needs a title and children")
	}
	def(&o.Done, "%d/%d done")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	done := 0
	for _, child := range children {
		if child.ID == "" || child.Title == "" {
			panic("project: Epic child must not be empty")
		}
		if child.Done {
			done++
		}
	}
	root.Children(func() {
		ui.Text(c, title).Bold()
		kit.Numeric(c, fmt.Sprintf(o.Done, done, len(children))).TextColor(t.TextMuted)
		progress(c, float64(done)/float64(len(children)), t.Success)
	})
	return root
}
