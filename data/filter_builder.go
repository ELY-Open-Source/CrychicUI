package data

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Rule is a field, operator and value.
type Rule struct{ Field, Op, Value string }

// FilterBuilderOptions tune a FilterBuilder.
type FilterBuilderOptions struct {
	Add string
	// Rule part names; each takes the rule number.
	Field, Op, Value, Remove string
}

// FilterBuilder edits rules; rows enter and leave.
func FilterBuilder(c *ui.Context, label string, fields, ops []string, rs *[]Rule, o FilterBuilderOptions) ui.Element {
	if label == "" || rs == nil {
		panic(fmt.Sprintf("data: FilterBuilder %q needs rules", label))
	}
	rules("FilterBuilder", label, fields, ops, *rs)
	def(&o.Add, "Add filter")
	def(&o.Field, "Field %d")
	def(&o.Op, "Operator %d")
	def(&o.Value, "Value %d")
	def(&o.Remove, "Remove filter %d")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	ids := ui.Local(root, "ids", func() []int { return nil })
	next := ui.Local(root, "next", func() int { return 0 })
	if len(*ids) != len(*rs) {
		*ids = nil
		for range *rs {
			*ids = append(*ids, *next)
			*next++
		}
	}
	say := ruleCopy{o.Field, o.Op, o.Value, o.Remove}
	moves := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{Y: -t.Space(2)}, Exit: &ui.Motion{Collapse: true}}
	root.Children(func() {
		drop := -1
		for i := range *rs {
			c.Key((*ids)[i])
			ui.Column(c).FillWidth().Transition(moves).Children(func() {
				edited, remove := ruleRow(c, fields, ops, &(*rs)[i], i+1, say)
				if edited {
					root.SetChanged()
					kit.Log("FilterBuilder", "edited", "index", i)
				}
				if remove {
					drop = i
				}
			})
		}
		if drop >= 0 {
			*rs, *ids = slices.Delete(slices.Clone(*rs), drop, drop+1), slices.Delete(slices.Clone(*ids), drop, drop+1)
			root.SetChanged()
			sound.Play(sound.Delete)
			kit.Log("FilterBuilder", "removed", "index", drop)
		}
		add := kit.Press(ui.ButtonBase(c.Key("add")).Label(o.Add).AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius))
		add.Background(add.AnimateColor("bg", t.SurfaceHover.Alpha(unit(add.Hovered())), ui.Snappy)).TextColor(t.Accent).Children(func() {
			ui.Icon(c, icon.Must("plus")).Size(16, 16)
			ui.Text(c, o.Add)
		})
		if add.Clicked() {
			*rs = append(*rs, Rule{Field: fields[0], Op: ops[0]})
			*ids = append(*ids, *next)
			*next++
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("FilterBuilder", "added", "count", len(*rs))
		}
	})
	return root
}
