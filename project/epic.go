package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EpicOptions word the epic; Key is an optional issue key.
type EpicOptions struct{ Done, Key string }

// Epic shows an epic with child progress.
func Epic(c *ui.Context, title string, children []Task, o EpicOptions) ui.Element {
	if title == "" || len(children) == 0 {
		panic(fmt.Sprintf("project: Epic %q needs a title and children", title))
	}
	ids := make([]string, len(children))
	done := 0
	for i, child := range children {
		validTask("Epic", child)
		ids[i] = child.ID
		if child.Done {
			done++
		}
	}
	unique("Epic", ids)
	def(&o.Done, "%d of %d done")
	t := c.Theme()
	share := float64(done) / float64(len(children))
	ink := t.Accent
	if done == len(children) {
		ink = t.Success
	}
	root := card(c)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
				tile(c, "zap", t.Warning)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					if o.Key != "" {
						kit.Muted(c, o.Key).FontSize(t.Rem(0.7)).Bold().SingleLine()
					}
					ui.Text(c, title).FontSize(t.Rem(0.9)).Bold().Wrap().FillWidth()
				})
				kit.Chip(c, fmt.Sprintf("%d%%", done*100/len(children)), ink)
			})
			progress(c, share, ink)
			kit.Numeric(c, fmt.Sprintf(o.Done, done, len(children))).FontSize(t.Rem(0.7)).TextColor(t.TextMuted)
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(1.5), t.Space(2), t.Space(1.5)).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
		list.Children(func() {
			for i, child := range children {
				c.Key(child.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5))
				kit.Stagger(c, list, row, i).Children(func() {
					status := Todo
					if child.Done {
						status = Done
					}
					name, mark := statusLook(t, status)
					glyph(c, name, 16).TextColor(mark)
					text := ui.Text(c, child.Title).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
					if child.Done {
						text.Strikethrough().TextColor(t.TextMuted)
					}
					if child.Assignee != "" {
						display.Avatar(c, child.Assignee, nil).Size(t.Space(6), t.Space(6)).FontSize(t.Rem(0.55))
					}
				})
			}
		})
	})
	return root
}
