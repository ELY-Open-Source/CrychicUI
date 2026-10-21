package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CheckGroupOptions tune a CheckGroup.
type CheckGroupOptions struct {
	// AllLabel names the parent check; "Select all" when empty.
	AllLabel string
}

// CheckGroup is checks under a parent that goes mixed.
func CheckGroup(c *ui.Context, labels []string, selected *[]string, o CheckGroupOptions) ui.Element {
	distinct("CheckGroup", labels)
	subset("CheckGroup", labels, selected)
	if o.AllLabel == "" {
		o.AllLabel = "Select all"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		n := len(*selected)
		all := n == len(labels)
		c.Key("all")
		parent := ui.CheckboxBase(c, &all).Label(o.AllLabel).Gap(t.Space(2))
		if n > 0 && !all {
			parent.Mixed()
		}
		parent.Children(func() {
			box(c, parent, n == len(labels), n > 0 && n < len(labels))
			ui.Text(c, o.AllLabel).Bold()
		})
		if kit.Changed(parent, sound.Toggle) {
			*selected = nil
			if all {
				*selected = slices.Clone(labels)
			}
			root.SetChanged()
			kit.Log("CheckGroup", "all", "count", len(*selected))
		}
		ui.Column(c).PaddingStart(t.Space(6)).Gap(t.Space(2)).Children(func() {
			for _, label := range labels {
				on := slices.Contains(*selected, label)
				c.Key(label)
				if Check(c, label, &on).Changed() {
					*selected = flip(*selected, label)
					root.SetChanged()
					kit.Log("CheckGroup", "toggled", "count", len(*selected))
				}
			}
		})
	})
}
