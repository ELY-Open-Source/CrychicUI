package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MemoriesOptions tune a Memories list.
type MemoriesOptions struct {
	Empty, Remove string
}

// Memories lists remembered facts with delete.
func Memories(c *ui.Context, label string, facts *[]string, o MemoriesOptions) ui.Element {
	if label == "" || facts == nil {
		panic(fmt.Sprintf("agent: Memories %q needs a label and facts", label))
	}
	def(&o.Empty, "Nothing remembered yet")
	def(&o.Remove, "Forget")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleGroup).Label(label)
	delAt := -1
	root.Children(func() {
		if len(*facts) == 0 {
			kit.Arrive(kit.Muted(c, o.Empty), ui.Pose{Fade: 1}, 0)
			return
		}
		for i, f := range *facts {
			if f == "" {
				panic(fmt.Sprintf("agent: Memories %q facts must not be empty", label))
			}
			c.Key(f)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			var del ui.Element
			row.Children(func() {
				ui.Icon(c, icon.Must("brain")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, f).Grow(1).Wrap()
				del = kit.Press(ui.ButtonBase(c).Label(o.Remove).Tooltip(o.Remove).Size(t.Space(6), t.Space(6)).Center().Radius(999))
				del.Background(del.AnimateColor("bg", t.Danger.Alpha(0.12*unit(del.Hovered())), ui.Snappy)).
					TextColor(t.TextMuted.Mix(t.Danger, unit(del.Hovered()))).
					Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
			})
			if del.Clicked() {
				delAt = i
			}
		}
	})
	if delAt >= 0 {
		*facts = append((*facts)[:delAt], (*facts)[delAt+1:]...)
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("Memories", "forgot", "left", len(*facts))
	}
	return root
}
