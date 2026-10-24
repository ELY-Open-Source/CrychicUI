package code

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Watch is an expression and its value.
type Watch struct {
	Expr, Value string
}

// WatchesOptions word a Watches list.
type WatchesOptions struct {
	Placeholder, Add, Remove, Missing, Empty string
}

var watchExit = ui.ElementTransition{Animation: &ui.Snappy, Exit: &ui.Motion{Scale: 0.9}}

// Watches shows expressions and values; add and remove.
func Watches(c *ui.Context, label string, watches *[]Watch, o WatchesOptions) ui.Element {
	if label == "" || watches == nil {
		panic(fmt.Sprintf("code: Watches %q needs a label and watches", label))
	}
	seen := map[string]bool{}
	for _, w := range *watches {
		if w.Expr == "" || seen[w.Expr] {
			panic(fmt.Sprintf("code: Watches %q watches need unique expressions", label))
		}
		seen[w.Expr] = true
	}
	def(&o.Placeholder, "Watch an expression")
	def(&o.Add, "Add")
	def(&o.Remove, "Remove")
	def(&o.Missing, "not available")
	def(&o.Empty, "No expressions watched")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	draft := ui.Local(root, "draft", func() string { return "" })
	pending := ui.Local(root, "pending", func() string { return "" })
	delAt := -1
	var add ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Icon(c, icon.Must("eye")).Size(14, 14).TextColor(t.TextMuted)
			ui.Text(c, label).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Numeric(c, strconv.Itoa(len(*watches))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		ui.Column(c).FillWidth().Padding(t.Space(1), 0).Children(func() {
			if len(*watches) == 0 {
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.8)).Padding(t.Space(2), t.Space(3)).Transition(kit.Enter)
			}
			for i, w := range *watches {
				c.Key(w.Expr)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).
					Transition(watchExit).Role(ui.RoleListItem).Label(w.Expr)
				stagger(c, root, row, i)
				var del ui.Element
				row.Children(func() {
					ui.Text(c, w.Expr).Font("monospace").FontSize(t.Rem(0.8)).Shrink(0).MaxWidthPercent(55).SingleLine()
					kit.Muted(c, "=").Font("monospace").FontSize(t.Rem(0.8)).Shrink(0)
					value, ink := w.Value, valueInk(t, w.Value)
					if value == "" {
						value, ink = o.Missing, t.TextMuted
					}
					ui.Text(c, value).Font("monospace").FontSize(t.Rem(0.8)).TextColor(ink).Grow(1).MinWidth(0).SingleLine()
					del = kit.Press(ui.ButtonBase(c).Label(o.Remove+" "+w.Expr).Tooltip(o.Remove).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999))
					del.Background(del.AnimateColor("bg", t.Danger.Alpha(0.12*unit(del.Hovered())), ui.Snappy)).
						TextColor(del.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(del.Hovered())), ui.Snappy)).
						Children(func() { ui.Icon(c, icon.Must("x")).Size(14, 14) })
				})
				if del.Clicked() {
					delAt = i
				}
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Icon(c, icon.Must("plus")).Size(14, 14).TextColor(t.TextMuted)
			in := ui.TextInput(c, draft).Grow(1).Placeholder(o.Placeholder).Label(o.Placeholder).Font("monospace").FontSize(t.Rem(0.8))
			in.OnTextKey(func(k ui.TextKey) bool {
				if k.Key != ui.KeyEnter || k.Mods != 0 {
					return false
				}
				if expr := strings.TrimSpace(k.Text()); expr != "" {
					*pending = expr
					k.SetText("")
				}
				return true
			})
			add = kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(1), t.Space(3)).Radius(999)).Disabled(strings.TrimSpace(*draft) == "")
			add.Background(add.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(add.Hovered())), ui.Snappy)).TextColor(t.AccentText).
				Children(func() { ui.Text(c, o.Add).FontSize(t.Rem(0.8)).Bold().Shrink(0) })
		})
	})
	if add.Clicked() && strings.TrimSpace(*draft) != "" {
		if addWatch(root, watches, strings.TrimSpace(*draft)) {
			*draft = ""
		}
	}
	if expr := *pending; expr != "" {
		*pending = ""
		addWatch(root, watches, expr)
	}
	if delAt >= 0 {
		*watches = slices.Delete(slices.Clone(*watches), delAt, delAt+1)
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Watches", "removed", "count", len(*watches))
	}
	return root
}

// addWatch appends expr unless present; reports which way it went.
func addWatch(root ui.Element, watches *[]Watch, expr string) (added bool) {
	if slices.ContainsFunc(*watches, func(w Watch) bool { return w.Expr == expr }) {
		kit.Log("Watches", "duplicate")
		return false
	}
	*watches = append(slices.Clone(*watches), Watch{Expr: expr})
	root.SetChanged()
	sound.Play(sound.Pop)
	kit.Log("Watches", "added", "count", len(*watches))
	return true
}
