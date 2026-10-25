package devtools

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StateNode is one key of app state.
type StateNode struct {
	Key, Value string
	Children   []StateNode
}

// StateTreeOptions word the panel; Count takes a child count.
type StateTreeOptions struct{ Title, Count string }

// StateTree shows app state; changed values flash; cue Toggle.
func StateTree(c *ui.Context, root StateNode, o StateTreeOptions) ui.Element {
	def(&o.Title, "State")
	def(&o.Count, "{%d}")
	t := c.Theme()
	box := card(c).Role(ui.RoleTree)
	open := ui.Local(box, "open", func() map[string]bool { return map[string]bool{strconv.Quote(root.Key): true} })
	seen := ui.Local(box, "seen", func() map[string]string { return map[string]string{} })
	flash := ui.Local(box, "flash", func() map[string]time.Time { return map[string]time.Time{} })
	now := c.Now()
	box.Children(func() {
		kit.CardHead(c, "braces", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(2))
		list.Children(func() {
			var walk func(n StateNode, path string, depth int)
			walk = func(n StateNode, path string, depth int) {
				if n.Key == "" {
					panic(fmt.Sprintf("devtools: StateTree node under %q needs a key", path))
				}
				keys := make([]string, len(n.Children))
				for i, k := range n.Children {
					keys[i] = k.Key
				}
				if len(keys) > 0 {
					unique("StateTree", keys)
				}
				branch := len(n.Children) > 0
				if !branch {
					if last, ok := (*seen)[path]; ok && last != n.Value {
						(*flash)[path] = now
					}
					(*seen)[path] = n.Value
				}
				lit := now.Sub((*flash)[path]) < 150*time.Millisecond
				if lit {
					c.After(150*time.Millisecond - now.Sub((*flash)[path]))
				}
				c.Key(path)
				row := kit.Press(ui.ButtonBase(c).Label(n.Key).FillWidth().Padding(t.Space(1), t.Space(2), t.Space(1),
					t.Space(2)+t.Space(4)*float32(depth)).Radius(t.Radius).Expanded(branch && (*open)[path]))
				glow := row.AnimateTo("flash", unit(lit), ui.Smooth)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy).Mix(t.Warning.Alpha(0.3), glow))
				row.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						chev := glyph(c, "chevron-right", 14).TextColor(t.TextMuted).Opacity(unit(branch))
						chev.Rotate(90 * chev.AnimateTo("turn", unit(branch && (*open)[path]), ui.Snappy))
						ui.Text(c, n.Key).Font("monospace").FontSize(t.Rem(0.78)).TextColor(t.Accent).Shrink(0).MaxWidth(t.Space(36)).SingleLine()
						if branch {
							kit.Muted(c, fmt.Sprintf(o.Count, len(n.Children))).Font("monospace").FontSize(t.Rem(0.72)).Shrink(0)
							return
						}
						v := ui.Text(c, n.Value).Font("monospace").FontSize(t.Rem(0.78)).Grow(1).MinWidth(0).SingleLine()
						v.TextColor(valueInk(t, n.Value))
					})
				})
				if branch && row.Clicked() {
					(*open)[path] = !(*open)[path]
					box.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("StateTree", "toggled", "depth", depth, "open", (*open)[path])
				}
				if branch && (*open)[path] {
					for _, k := range n.Children {
						walk(k, path+"/"+strconv.Quote(k.Key), depth+1)
					}
				}
			}
			walk(root, strconv.Quote(root.Key), 0)
		})
	})
	return box
}

// valueInk colors a value by its look.
func valueInk(t *ui.Theme, v string) ui.Color {
	switch {
	case strings.HasPrefix(v, `"`):
		return t.Success
	case v == "true" || v == "false":
		return t.Warning
	case v == "null" || v == "":
		return t.TextMuted
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return t.Accent
	}
	return t.Text
}
