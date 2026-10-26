package desktop

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Spaces are workspace dots; the current one widens; cue Select.
func Spaces(c *ui.Context, names []string, current *int) ui.Element {
	unique("Spaces", names)
	if len(names) == 0 || current == nil || *current < 0 || *current >= len(names) {
		panic(fmt.Sprintf("desktop: Spaces needs workspaces and a current one, got %d", len(names)))
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2.5)).Radius(999).
		Background(t.SurfaceHover).Border(1, t.Border).Role(ui.RoleTabList)
	kit.Arrive(root, ui.Pose{Scale: 0.9, Fade: 1}, 0)
	return root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			for i, name := range names {
				c.Key(strconv.Quote(name))
				on := i == *current
				b := ui.ButtonBase(c).Label(name).Role(ui.RoleTab).Checked(on).Padding(t.Space(1), 0).Center()
				w := b.AnimateTo("width", t.Space(2.5)+t.Space(5.5)*unit(on), ui.Snappy)
				b.Children(func() {
					d := ui.Box(c).Size(w, t.Space(2.5)).Radius(999)
					d.Background(d.AnimateColor("dot", t.TextMuted.Alpha(0.4+0.4*unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
				})
				if b.Clicked() && !on {
					*current = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Spaces", "current", "index", i)
				}
			}
		})
		c.Key(strconv.Quote(names[*current]))
		kit.Arrive(ui.Text(c, names[*current]).FontSize(t.Rem(0.74)).Bold().Shrink(0).MaxWidth(t.Space(36)).SingleLine(), ui.Pose{X: t.Space(1), Fade: 1}, 0)
	})
}
