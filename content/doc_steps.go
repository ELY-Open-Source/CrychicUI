package content

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DocStep is a numbered step's title and text.
type DocStep struct{ Title, Text string }

// DocSteps are numbered steps joined by a line.
func DocSteps(c *ui.Context, steps []DocStep) ui.Element {
	if len(steps) == 0 || len(steps) > 20 {
		panic(fmt.Sprintf("content: DocSteps needs 1 to 20 steps, got %d", len(steps)))
	}
	for _, s := range steps {
		if strings.TrimSpace(s.Title) == "" {
			panic("content: DocSteps step needs a title")
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleList)
	root.Children(func() {
		for i, s := range steps {
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
			kit.Stagger(c, root, row, i).Children(func() {
				ui.Column(c).Shrink(0).AlignItems(ui.Center).Children(func() {
					n := ui.Box(c).Size(t.Space(7), t.Space(7)).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText)
					kit.Arrive(n, ui.Pose{Scale: 0.6, Fade: 1}, 0).Children(func() {
						ui.Text(c, strconv.Itoa(i+1)).FontSize(t.Rem(0.78)).Bold()
					})
					if i < len(steps)-1 {
						ui.Box(c).Width(1).Grow(1).MinHeight(t.Space(4)).Background(t.Border)
					}
				})
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Padding(t.Space(1), 0, t.Space(5), 0).Children(func() {
					ui.Text(c, s.Title).FontSize(t.Rem(0.9)).Bold().FillWidth().Wrap()
					if s.Text != "" {
						kit.Muted(c, s.Text).FontSize(t.Rem(0.82)).LineHeight(1.55).FillWidth().Wrap()
					}
				})
			})
		}
	})
	return root
}
