package form

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Rule is one AdvancedSearch condition.
type Rule struct{ Field, Op, Value string }

// AdvancedSearchOptions tune an AdvancedSearch.
type AdvancedSearchOptions struct {
	// Operators default to contains, is, is not and starts with.
	Operators []string
	// The rest are labels and button copy.
	FieldLabel, OpLabel, ValueLabel, Remove, Add, Search string
}

// AdvancedSearch edits field, operator and value rules.
func AdvancedSearch(c *ui.Context, fields []string, rules *[]Rule, o AdvancedSearchOptions) (root, search ui.Element) {
	if len(fields) == 0 || rules == nil {
		panic("form: AdvancedSearch needs fields and rules")
	}
	if len(o.Operators) == 0 {
		o.Operators = []string{"contains", "is", "is not", "starts with"}
	}
	def(&o.FieldLabel, "Field")
	def(&o.OpLabel, "Operator")
	def(&o.ValueLabel, "Value")
	def(&o.Remove, "Remove rule")
	def(&o.Add, "Add rule")
	def(&o.Search, "Search")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(3))
	ids := ui.Local(root, "rows", func() *rows { return &rows{} })
	keys := (*ids).keys(len(*rules))
	root.Children(func() {
		gone := -1
		for i, key := range keys {
			r := &(*rules)[i]
			ui.Row(c).Key(key).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Transition(rowMotion(c)).Children(func() {
				edited(root, choice.Dropdown(c, o.FieldLabel, fields, &r.Field))
				edited(root, choice.Dropdown(c, o.OpLabel, o.Operators, &r.Op))
				edited(root, ui.TextInput(c, &r.Value).Label(o.ValueLabel).Grow(1))
				if button.IconButton(c, o.Remove, "x").Clicked() {
					gone = i
				}
			})
		}
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			if button.Ghost(c, o.Add).Clicked() {
				*rules = append(*rules, Rule{Field: fields[0], Op: o.Operators[0]})
				(*ids).add()
				root.SetChanged()
				kit.Log("AdvancedSearch", "add", "rules", len(*rules))
			}
			ui.Box(c).Grow(1)
			search = button.Solid(c, o.Search).Disabled(len(*rules) == 0)
		})
		if gone >= 0 {
			*rules = append((*rules)[:gone], (*rules)[gone+1:]...)
			(*ids).drop(gone)
			root.SetChanged()
			kit.Log("AdvancedSearch", "remove", "rules", len(*rules))
		}
	})
	return root, search
}
