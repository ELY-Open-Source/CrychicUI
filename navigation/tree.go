package navigation

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Node is a labeled destination; nil Children makes a leaf.
type Node struct {
	Label    string
	Children []Node
}

// nodes checks labels: present, slash-free, unique among siblings.
func nodes(component string, list []Node) {
	seen := map[string]bool{}
	for _, n := range list {
		if n.Label == "" || strings.Contains(n.Label, "/") || seen[n.Label] {
			panic(fmt.Sprintf("navigation: %s labels must be present, slash-free and unique among siblings: %q", component, n.Label))
		}
		seen[n.Label] = true
		nodes(component, n.Children)
	}
}

// branch is one tree build: shared state for every row.
type branch struct {
	component string
	changed   *bool
	open      map[string]bool
	chosen    *string
	glyph     func(n Node, open bool) string
	cue       sound.Cue
}

// rows builds list as indented rows under parent.
func (b branch) rows(c *ui.Context, list []Node, parent string, depth int) {
	t := c.Theme()
	back, forth := ui.KeyLeft, ui.KeyRight
	turn := float32(90)
	if c.Direction() == ui.RightToLeft {
		back, forth, turn = forth, back, -90
	}
	for _, n := range list {
		path := n.Label
		if parent != "" {
			path = parent + "/" + n.Label
		}
		c.Key(path)
		folder, open := n.Children != nil, b.open[path]
		r := kit.Press(ui.ButtonBase(c).Label(n.Label).FillWidth().Radius(t.Radius)).Transition(kit.Enter)
		if folder {
			r.Expanded(open)
		} else {
			r.Checked(path == *b.chosen)
		}
		bg := t.SurfaceHover.Alpha(0)
		switch {
		case path == *b.chosen:
			bg = t.Accent.Alpha(0.12)
		case r.Hovered():
			bg = t.SurfaceHover
		}
		toggle := func(to bool) {
			if folder && to != open {
				b.open[path] = to
				*b.changed = true
				sound.Play(sound.Toggle)
				kit.Log(b.component, "toggled", "depth", depth, "open", to)
			}
		}
		if r.Shortcut(0, forth) {
			toggle(true)
		}
		if r.Shortcut(0, back) {
			toggle(false)
		}
		r.Background(r.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2), t.Space(1.5), t.Space(2)+float32(depth)*t.Space(4)).
				Children(func() {
					arrow := ui.Box(c).Size(16, 16).Shrink(0)
					if folder {
						arrow.Rotate(turn * arrow.AnimateTo("turn", unit(open), ui.Snappy)).Children(func() {
							ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
						})
					}
					if b.glyph != nil {
						ui.Icon(c, icon.Must(b.glyph(n, open))).Size(16, 16).ContentTransition(ui.ContentTransition{})
					}
					ui.Text(c, n.Label).SingleLine()
				})
		})
		if r.Clicked() {
			if folder {
				toggle(!open)
			} else if path != *b.chosen {
				*b.chosen = path
				*b.changed = true
				sound.Play(b.cue)
				kit.Log(b.component, "chose", "depth", depth)
			}
		}
		if folder && b.open[path] {
			b.rows(c, n.Children, path, depth+1)
		}
	}
}

// tree lays nodes out as indented rows.
func tree(c *ui.Context, component, label string, list []Node, chosen *string, open map[string]bool, glyph func(Node, bool) string, cue sound.Cue) ui.Element {
	if label == "" || len(list) == 0 || chosen == nil || open == nil {
		panic(fmt.Sprintf("navigation: %s needs a label, nodes, a choice and an open set", component))
	}
	nodes(component, list)
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(64)).Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Label(label)
	changed := false
	root.Children(func() {
		branch{component, &changed, open, chosen, glyph, cue}.rows(c, list, "", 0)
	})
	if changed {
		root.SetChanged()
	}
	return root
}
