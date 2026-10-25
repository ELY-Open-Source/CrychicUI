package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReleaseKind sorts a release note.
type ReleaseKind uint8

const (
	ReleaseAdded ReleaseKind = iota
	ReleaseFixed
	ReleaseChanged
)

// ReleaseItem is one release note.
type ReleaseItem struct {
	Kind ReleaseKind
	Text string
}

// ReleaseOptions word the card; Date is optional.
type ReleaseOptions struct{ Ship, Shipped, Draft, Date, Added, Fixed, Changed string }

// Release shows version notes and ship; cue Success.
func Release(c *ui.Context, version string, items []ReleaseItem, shipped *bool, o ReleaseOptions) ui.Element {
	if version == "" || len(items) == 0 || shipped == nil {
		panic(fmt.Sprintf("project: Release %q needs notes and state", version))
	}
	for i, item := range items {
		if item.Text == "" || item.Kind > ReleaseChanged {
			panic(fmt.Sprintf("project: Release %q note %d is invalid", version, i))
		}
	}
	def(&o.Ship, "Ship")
	def(&o.Shipped, "Shipped")
	def(&o.Draft, "Draft")
	def(&o.Added, "Added")
	def(&o.Fixed, "Fixed")
	def(&o.Changed, "Changed")
	t := c.Theme()
	looks := []struct {
		name, glyph string
		ink         ui.Color
	}{{o.Added, "circle-plus", t.Success}, {o.Fixed, "bug", t.Danger}, {o.Changed, "wrench", t.Accent}}
	root := card(c)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ink := t.Accent
			if *shipped {
				ink = t.Success
			}
			tile(c, "rocket", ink)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, version).FontSize(t.Rem(1)).Bold().MinWidth(0).SingleLine()
				if o.Date != "" {
					kit.Muted(c, o.Date).FontSize(t.Rem(0.72)).SingleLine()
				}
			})
			c.Key(*shipped)
			if *shipped {
				kit.Arrive(kit.Chip(c, o.Shipped, t.Success), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			} else {
				kit.Chip(c, o.Draft, t.Warning)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		list.Children(func() {
			for i, item := range items {
				c.Key(i)
				look := looks[item.Kind]
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).Width(t.Space(18)).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(look.ink).Children(func() {
						glyph(c, look.glyph, 14)
						ui.Text(c, look.name).FontSize(t.Rem(0.7)).Bold().MinWidth(0).SingleLine()
					})
					ui.Text(c, item.Text).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
				})
			}
		})
		ui.Box(c).FillWidth().Padding(0, t.Space(3), t.Space(3), t.Space(3)).Children(func() {
			name, mark := o.Ship, "rocket"
			if *shipped {
				name, mark = o.Shipped, "circle-check"
			}
			b := kit.Press(ui.ButtonBase(c).Label(name).FillWidth().Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Disabled(*shipped))
			b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())).Mix(t.Success, unit(*shipped)), ui.Snappy))
			b.TextColor(t.AccentText).Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					c.Key(mark)
					kit.Arrive(glyph(c, mark, 16), ui.Pose{Y: t.Space(1), Fade: 1}, 0)
					ui.Text(c, name).FontSize(t.Rem(0.85)).Bold().Shrink(0)
				})
			})
			if b.Clicked() && !*shipped {
				*shipped = true
				root.SetChanged()
				sound.Play(sound.Success)
				kit.Log("Release", "shipped", "notes", len(items))
			}
		})
	})
	return root
}
