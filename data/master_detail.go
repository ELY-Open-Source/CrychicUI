package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MasterDetailOptions tune a MasterDetail.
type MasterDetailOptions struct {
	// Empty shows while nothing is chosen.
	Empty string
}

// MasterDetail lists items beside the chosen one's detail.
func MasterDetail(c *ui.Context, label string, items []string, chosen *int, detail func(i int), o MasterDetailOptions) ui.Element {
	if label == "" || len(items) == 0 || chosen == nil || detail == nil || *chosen < -1 || *chosen >= len(items) {
		panic(fmt.Sprintf("data: MasterDetail %q needs items, a choice or -1, and a detail", label))
	}
	unique("MasterDetail", label, items)
	def(&o.Empty, "Choose an item")
	t := c.Theme()
	root := ui.Row(c).FillWidth().AlignItems(ui.Stretch).Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup).Label(label)
	shown := *chosen
	fade := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	return root.Children(func() {
		ui.Column(c).Width(t.Space(56)).Shrink(0).Padding(t.Space(1)).Gap(t.Space(0.5)).Background(t.Surface).
			Role(ui.RoleList).FocusGroup(ui.Vertical).Children(func() {
			for i, it := range items {
				c.Key(it)
				on := i == shown
				b := kit.Press(ui.ButtonBase(c).Label(it).FillWidth().Justify(ui.Start).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.14), unit(on)), ui.Snappy)).
					Children(func() { ui.Text(c, it).SingleLine() })
				if b.Clicked() && !on {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("MasterDetail", "chose", "index", i)
				}
			}
		})
		ui.Box(c).Grow(1).Padding(t.Space(4)).Children(func() {
			c.Key(shown)
			ui.Column(c).FillWidth().Gap(t.Space(2)).Transition(fade).Children(func() {
				if shown < 0 {
					kit.Muted(c, o.Empty)
					return
				}
				detail(shown)
			})
		})
	})
}
