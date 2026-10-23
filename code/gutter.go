package code

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Gutter shows line numbers with breakpoints toggled.
func Gutter(c *ui.Context, lines int, breakpoints *[]int) ui.Element {
	if lines < 1 || breakpoints == nil {
		panic(fmt.Sprintf("code: Gutter needs lines and breakpoints, got %d", lines))
	}
	for _, b := range *breakpoints {
		if b < 1 || b > lines {
			panic(fmt.Sprintf("code: Gutter breakpoints must be lines from 1 to %d", lines))
		}
	}
	t := c.Theme()
	width := t.Space(float32(len(strconv.Itoa(lines))*2 + 5))
	root := ui.Column(c).AlignSelf(ui.Start).Width(width).Role(ui.RoleGroup).Label("Gutter")
	set := slices.Clone(*breakpoints)
	slices.Sort(set)
	return root.Children(func() {
		for i := 1; i <= lines; i++ {
			c.Key(i)
			on := slices.Contains(set, i)
			b := kit.Press(ui.ButtonBase(c).Label(strconv.Itoa(i)).Row().FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).
				Padding(0, t.Space(1)).Checked(on))
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Danger.Alpha(0.12)
			}
			b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
				dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999)
				dot.Background(dot.AnimateColor(fmt.Sprint("dot.", i), t.Danger.Alpha(unit(on)), ui.Snappy))
				kit.Muted(c, strconv.Itoa(i)).Font("monospace").FontSize(t.Rem(0.8)).Grow(1)
			})
			if b.Clicked() {
				if on {
					*breakpoints = slices.DeleteFunc(*breakpoints, func(v int) bool { return v == i })
				} else {
					*breakpoints = append(*breakpoints, i)
				}
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Gutter", "breakpoints", "count", len(*breakpoints))
			}
		}
	})
}
