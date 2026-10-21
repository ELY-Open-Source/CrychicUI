package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Prose spaces long-form paragraphs and fades in.
func Prose(c *ui.Context, paragraphs []string) ui.Element {
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Children(func() {
		for i, paragraph := range paragraphs {
			ui.Text(c, paragraph).Key(i).ContentTransition(ui.ContentTransition{})
		}
	})
	return kit.Arrive(root, ui.Pose{Fade: 1}, 0)
}
