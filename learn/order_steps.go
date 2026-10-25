package learn

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// OrderStepsOptions word the list; Count takes placed and total.
type OrderStepsOptions struct{ Title, Count, Done string }

// OrderSteps drags steps into order; cue Drop.
func OrderSteps(c *ui.Context, steps *[]string, correct []string, o OrderStepsOptions) ui.Element {
	def(&o.Title, "Put the steps in order")
	if steps == nil || len(correct) < 2 || len(*steps) != len(correct) {
		panic(fmt.Sprintf("learn: OrderSteps %q needs at least 2 steps, as many as the answer", o.Title))
	}
	unique("OrderSteps", correct)
	if !slices.Equal(slices.Sorted(slices.Values(*steps)), slices.Sorted(slices.Values(correct))) {
		panic(fmt.Sprintf("learn: OrderSteps %q steps are not the answer's steps", o.Title))
	}
	def(&o.Count, "%d of %d in place")
	def(&o.Done, "In order")
	t := c.Theme()
	placed := 0
	for i, s := range *steps {
		if s == correct[i] {
			placed++
		}
	}
	all := placed == len(correct)
	root := card(c).Role(ui.RoleList)
	id := ui.Local(root, "id", func() *int { return new(int) })
	from, to := -1, -1
	root.Children(func() {
		kit.CardHead(c, "list-ordered", o.Title, func() {
			c.Key(all)
			if all {
				kit.Arrive(kit.Chip(c, o.Done, t.Success), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			} else {
				kit.Chip(c, fmt.Sprintf(o.Count, placed, len(correct)), t.Accent)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(3))
		list.Children(func() {
			for i, s := range *steps {
				good := s == correct[i]
				c.Key(strconv.Quote(s))
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).
					Label(s).Focusable().Cursor(ui.CursorGrab).Drag(stepDrag{list: *id, from: i}).
					Transition(ui.ElementTransition{Animation: &ui.Snappy, Position: true})
				_, over := ui.DragOver[stepDrag](row)
				ink := t.Border
				if good {
					ink = t.Success
				}
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())).Mix(t.Accent.Alpha(0.12), unit(over)), ui.Snappy))
				row.Border(1, row.AnimateColor("line", ink.Mix(t.Accent, unit(over)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					glyph(c, "grip-vertical", 16).TextColor(t.TextMuted)
					badge := ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999)
					badge.Background(badge.AnimateColor("bg", t.SurfaceHover.Mix(t.Success, unit(good)), ui.Snappy))
					badge.TextColor(badge.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(good)), ui.Snappy))
					badge.Children(func() { ui.Text(c, fmt.Sprint(i+1)).FontSize(t.Rem(0.72)).Bold().Shrink(0) })
					ui.Text(c, s).FontSize(t.Rem(0.84)).Grow(1).MinWidth(0).Wrap()
				})
				if d, ok := ui.Drop[stepDrag](row); ok && d.list == *id && d.from != i {
					from, to = d.from, i
				}
			}
		})
	})
	if from >= 0 {
		s := (*steps)[from]
		*steps = slices.Insert(slices.Delete(slices.Clone(*steps), from, from+1), to, s)
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("OrderSteps", "moved", "from", from, "to", to)
	}
	return root
}

type stepDrag struct {
	list *int
	from int
}
