package code

import (
	"fmt"
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

// WatchesOptions tune a Watches list.
type WatchesOptions struct {
	Placeholder, Add, Remove string
}

// Watches shows expressions and values; add and remove.
func Watches(c *ui.Context, label string, watches *[]Watch, o WatchesOptions) ui.Element {
	if label == "" || watches == nil {
		panic(fmt.Sprintf("code: Watches %q needs a label and watches", label))
	}
	def(&o.Placeholder, "Watch an expression")
	def(&o.Add, "Add")
	def(&o.Remove, "Remove")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	draft := ui.Local(root, "draft", func() string { return "" })
	pending := ui.Local(root, "pending", func() string { return "" })
	delAt := -1
	root.Children(func() {
		for i, w := range *watches {
			if w.Expr == "" {
				panic(fmt.Sprintf("code: Watches %q watches need expressions", label))
			}
			c.Key(w.Expr)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			var del ui.Element
			row.Children(func() {
				ui.Text(c, w.Expr).Font("monospace").FontSize(t.Rem(0.85)).Grow(1).SingleLine()
				value := w.Value
				if value == "" {
					value = "—"
				}
				ui.Text(c, value).Font("monospace").FontSize(t.Rem(0.85)).Bold().SingleLine()
				del = kit.Press(ui.ButtonBase(c).Label(o.Remove).Tooltip(o.Remove).Size(t.Space(5), t.Space(5)).Center().Radius(999))
				del.Background(del.AnimateColor("bg", t.Danger.Alpha(0.12*unit(del.Hovered())), ui.Snappy)).
					TextColor(t.TextMuted.Mix(t.Danger, unit(del.Hovered()))).
					Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
			})
			if del.Clicked() {
				delAt = i
			}
		}
		ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			in := ui.TextInput(c, draft).Grow(1).Placeholder(o.Placeholder).Label(o.Placeholder).Font("monospace").FontSize(t.Rem(0.85))
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
			add := kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(1), t.Space(2)).Radius(999).Border(1, t.Border)).
				Disabled(strings.TrimSpace(*draft) == "")
			add.Background(add.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(add.Hovered())), ui.Snappy)).Children(func() {
				ui.Text(c, o.Add).FontSize(t.Rem(0.85)).SingleLine()
			})
			if add.Clicked() && strings.TrimSpace(*draft) != "" {
				expr := strings.TrimSpace(*draft)
				if addWatch(c, root, watches, expr) {
					*draft = ""
				}
			}
		})
	})
	if expr := *pending; expr != "" {
		*pending = ""
		addWatch(c, root, watches, expr)
	}
	if delAt >= 0 {
		*watches = append((*watches)[:delAt], (*watches)[delAt+1:]...)
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Watches", "removed", "count", len(*watches))
	}
	return root
}

// addWatch appends expr unless present; reports which way it went.
func addWatch(c *ui.Context, root ui.Element, watches *[]Watch, expr string) (added bool) {
	for _, w := range *watches {
		if w.Expr == expr {
			kit.Log("Watches", "duplicate")
			return false
		}
	}
	*watches = append(*watches, Watch{Expr: expr})
	root.SetChanged()
	sound.Play(sound.Pop)
	kit.Log("Watches", "added", "count", len(*watches))
	return true
}
