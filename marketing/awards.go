package marketing

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Award is one badge; Icon defaults to an award ribbon.
type Award struct{ Title, Org, Year, Icon string }

// AwardsOptions word the row.
type AwardsOptions struct{ Title string }

// Awards is a row of award badges.
func Awards(c *ui.Context, awards []Award, o AwardsOptions) ui.Element {
	if len(awards) == 0 {
		panic(fmt.Sprintf("marketing: Awards %q needs awards", o.Title))
	}
	titles := make([]string, len(awards))
	for i, a := range awards {
		titles[i] = a.Title
	}
	unique("Awards", titles)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, "")
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(4)).Justify(ui.Center).Children(func() {
			for i, a := range awards {
				mark := a.Icon
				if mark == "" {
					mark = "award"
				}
				c.Key(strconv.Quote(a.Title))
				badge := ui.Column(c).Width(t.Space(32)).Gap(t.Space(1.5)).AlignItems(ui.Center)
				kit.Arrive(badge, ui.Pose{Scale: 0.6, Rotate: -12, Fade: 1}, time.Duration(i)*80*time.Millisecond).Children(func() {
					ui.Box(c).Size(t.Space(14), t.Space(14)).Shrink(0).Center().Radius(999).Background(t.Warning.Alpha(0.12)).
						Border(t.Space(0.5), t.Warning.Alpha(0.4)).TextColor(t.Warning).Children(func() { glyph(c, mark, 26) })
					ui.Text(c, a.Title).FontSize(t.Rem(0.8)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
					if a.Org != "" {
						kit.Muted(c, a.Org).FontSize(t.Rem(0.7)).FillWidth().TextAlign(ui.Center).Wrap()
					}
					if a.Year != "" {
						ui.Text(c, a.Year).FontSize(t.Rem(0.7)).Bold().TextColor(t.Warning).Shrink(0)
					}
				})
			}
		})
	})
	return root
}
