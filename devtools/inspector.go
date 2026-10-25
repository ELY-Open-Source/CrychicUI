package devtools

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Node is one element in an inspected tree.
type Node struct {
	ID, Tag  string
	Props    []Attr
	Children []Node
}

// InspectorOptions word the panel; Label takes tag and ID.
type InspectorOptions struct{ Title, Label, Props, Pick, None string }

// Inspector shows an element tree and props; cue Select.
func Inspector(c *ui.Context, root Node, selected *string, o InspectorOptions) ui.Element {
	def(&o.Title, "Inspector")
	if selected == nil {
		panic(fmt.Sprintf("devtools: Inspector %q needs a selection", o.Title))
	}
	type row struct {
		n     Node
		depth int
	}
	var rows []row
	var walk func(n Node, depth int)
	walk = func(n Node, depth int) {
		if strings.TrimSpace(n.Tag) == "" {
			panic(fmt.Sprintf("devtools: Inspector %q node %q needs a tag", o.Title, n.ID))
		}
		rows = append(rows, row{n, depth})
		for _, k := range n.Children {
			walk(k, depth+1)
		}
	}
	walk(root, 0)
	ids := make([]string, len(rows))
	for i, r := range rows {
		ids[i] = r.n.ID
	}
	unique("Inspector", ids)
	def(&o.Label, "%s %s")
	def(&o.Props, "Props")
	def(&o.Pick, "Pick an element to see its props")
	def(&o.None, "No props")
	t := c.Theme()
	box := card(c).Role(ui.RoleTree)
	box.Children(func() {
		kit.CardHead(c, "code-xml", o.Title, func() { kit.Chip(c, fmt.Sprint(len(rows)), t.Accent) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5))
		list.Children(func() {
			for i, r := range rows {
				on := *selected == r.n.ID
				c.Key(r.n.ID)
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, r.n.Tag, r.n.ID)).FillWidth().
					Padding(t.Space(1), t.Space(2), t.Space(1), t.Space(2)+t.Space(4)*float32(r.depth)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.RichText(c, ui.Span{Text: "<", Color: t.TextMuted}, ui.Span{Text: r.n.Tag, Color: t.Accent, Weight: 700},
							ui.Span{Text: ">", Color: t.TextMuted}).Font("monospace").FontSize(t.Rem(0.78)).Shrink(0).MaxWidth(t.Space(36)).SingleLine()
						if len(r.n.Props) > 0 {
							p := r.n.Props[0]
							kit.Muted(c, p.Key+"="+p.Value).Font("monospace").FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).SingleLine()
						}
					})
				})
				if b.Clicked() && !on {
					*selected = r.n.ID
					box.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Inspector", "selected", "index", i, "depth", r.depth)
				}
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2.5), t.Space(3)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			kit.Muted(c, o.Props).FontSize(t.Rem(0.7)).Bold()
			picked := -1
			for i, r := range rows {
				if r.n.ID == *selected {
					picked = i
				}
			}
			if picked < 0 {
				kit.Muted(c, o.Pick).FontSize(t.Rem(0.75)).Wrap()
				return
			}
			n := rows[picked].n
			c.Key(n.ID)
			panel := ui.Column(c).FillWidth().Gap(t.Space(1))
			kit.Arrive(panel, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
				if len(n.Props) == 0 {
					kit.Muted(c, o.None).FontSize(t.Rem(0.75))
				}
				for _, p := range n.Props {
					ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
						kit.Muted(c, p.Key).Font("monospace").FontSize(t.Rem(0.72)).Width(t.Space(24)).Shrink(0).SingleLine()
						ui.Text(c, p.Value).Font("monospace").FontSize(t.Rem(0.72)).TextColor(valueInk(t, p.Value)).Grow(1).MinWidth(0).Wrap()
					})
				}
			})
		})
	})
	return box
}
