package agent

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Memory is one remembered fact.
type Memory struct {
	Fact, Saved string
}

// MemoriesOptions tune a Memories list.
type MemoriesOptions struct {
	Count, Empty, Remove string
}

// Memories lists remembered facts with delete.
func Memories(c *ui.Context, label string, facts *[]Memory, o MemoriesOptions) ui.Element {
	if label == "" || facts == nil {
		panic(fmt.Sprintf("agent: Memories %q needs a label and facts", label))
	}
	seen := map[string]bool{}
	for _, m := range *facts {
		if m.Fact == "" || seen[m.Fact] {
			panic(fmt.Sprintf("agent: Memories %q facts must be unique and not empty", label))
		}
		seen[m.Fact] = true
	}
	def(&o.Count, "%d remembered")
	def(&o.Empty, "Nothing remembered yet")
	def(&o.Remove, "Forget")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Role(ui.RoleList).Label(label)
	delAt := -1
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, 0, t.Space(1), 0).Children(func() {
			tile(c, "brain", t.Accent, t.Space(7))
			ui.Text(c, label).Bold().Grow(1).SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Count, len(*facts))).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).
				ContentTransition(ui.ContentTransition{})
		})
		if len(*facts) == 0 {
			empty := ui.Column(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(5))
			kit.Arrive(empty, ui.Pose{Fade: 1}, 0).Children(func() {
				ui.Icon(c, icon.Must("brain")).Size(16, 16).TextColor(t.TextMuted)
				kit.Muted(c, o.Empty)
			})
			return
		}
		for i, m := range *facts {
			c.Key(m.Fact)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
				Transition(ui.ElementTransition{Animation: &ui.Snappy, Exit: &ui.Motion{Collapse: true}})
			kit.Arrive(row, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(i)*30*time.Millisecond)
			row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())), ui.Snappy))
			var del ui.Element
			row.Children(func() {
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, m.Fact).Wrap()
					if m.Saved != "" {
						kit.Muted(c, m.Saved).FontSize(t.Rem(0.75))
					}
				})
				del = kit.Press(ui.ButtonBase(c).Label(o.Remove+" "+m.Fact).Tooltip(o.Remove).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999))
				del.Opacity(del.AnimateTo("show", max(unit(row.Hovered()), unit(del.Focused()), 0.45), ui.Snappy))
				del.Background(del.AnimateColor("bg", t.Danger.Alpha(0.12*unit(del.Hovered())), ui.Snappy)).
					TextColor(del.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(del.Hovered())), ui.Snappy)).
					Children(func() { ui.Icon(c, icon.Must("trash-2")).Size(16, 16) })
			})
			if del.Clicked() {
				delAt = i
			}
		}
	})
	if delAt >= 0 {
		*facts = slices.Delete(slices.Clone(*facts), delAt, delAt+1)
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("Memories", "forgot", "left", len(*facts))
	}
	return root
}
