package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Node is a label with nested nodes.
type Node struct {
	Label    string
	Children []Node
}

// CheckTree cascades checks; checked holds leaf paths.
func CheckTree(c *ui.Context, label string, nodes []Node, checked *ui.Selection[string]) ui.Element {
	if label == "" || len(nodes) == 0 || checked == nil {
		panic(fmt.Sprintf("data: CheckTree %q needs nodes and checks", label))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleTree).Label(label)
	// leaves lists the leaf paths under n.
	var leaves func(n Node, path string) []string
	leaves = func(n Node, path string) []string {
		if len(n.Children) == 0 {
			return []string{path}
		}
		var out []string
		for _, k := range n.Children {
			out = append(out, leaves(k, path+"/"+k.Label)...)
		}
		return out
	}
	var level func(ns []Node, at string, depth int)
	level = func(ns []Node, at string, depth int) {
		var names []string
		for _, n := range ns {
			names = append(names, n.Label)
		}
		named("CheckTree", label, names)
		for _, n := range ns {
			path := at + n.Label
			all := leaves(n, path)
			got := 0
			for _, p := range all {
				if checked.Has(p) {
					got++
				}
			}
			on := got == len(all)
			c.Key(n.Label)
			ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), 0, t.Space(1), t.Space(5)*float32(depth)).Children(func() {
					if checkbox(c, n.Label, &on, got > 0 && got < len(all)).Changed() {
						for _, p := range all {
							if on {
								checked.Add(p)
							} else {
								checked.Remove(p)
							}
						}
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("CheckTree", "checked", "leaves", len(all), "on", on)
					}
					ui.Text(c, n.Label)
				})
				level(n.Children, path+"/", depth+1)
			})
		}
	}
	root.Children(func() { level(nodes, "", 0) })
	return root
}
