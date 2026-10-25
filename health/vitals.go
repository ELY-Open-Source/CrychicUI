package health

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Range is where a vital reading sits.
type Range int

const (
	InRange Range = iota
	High
	Low
)

// Vital is one reading; Value is shown as given.
type Vital struct {
	Name, Value, Unit, Icon string
	Range                   Range
}

// VitalsOptions word the grid; Ranges are normal, high and low.
type VitalsOptions struct {
	Title  string
	Ranges [3]string
}

// Vitals lays readings out as tiles marked by range.
func Vitals(c *ui.Context, vitals []Vital, o VitalsOptions) ui.Element {
	def(&o.Title, "Vitals")
	if len(vitals) == 0 {
		panic(fmt.Sprintf("health: Vitals %q needs readings", o.Title))
	}
	names := make([]string, len(vitals))
	for i, v := range vitals {
		if strings.TrimSpace(v.Value) == "" || v.Icon == "" || v.Range < InRange || v.Range > Low {
			panic(fmt.Sprintf("health: Vitals reading %q needs a value, an icon and a range", v.Name))
		}
		names[i] = v.Name
	}
	unique("Vitals", names)
	for i, r := range [...]string{"Normal", "High", "Low"} {
		def(&o.Ranges[i], r)
	}
	t := c.Theme()
	inks := [3]ui.Color{t.Success, t.Danger, t.Warning}
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "activity", o.Title, nil)
		grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		grid.Children(func() {
			for i, v := range vitals {
				ink := inks[v.Range]
				cell := ui.Column(c).MinWidth(0).Gap(t.Space(1.5)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).
					Background(ink.Alpha(0.04))
				kit.Stagger(c, grid, cell, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						glyph(c, v.Icon, 16).TextColor(ink)
						kit.Muted(c, v.Name).FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).SingleLine()
					})
					ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
						kit.Numeric(c, v.Value).FontSize(t.Rem(1.2)).Bold().MinWidth(0).SingleLine()
						kit.Muted(c, v.Unit).FontSize(t.Rem(0.66)).Shrink(0).MaxWidth(t.Space(14)).SingleLine()
					})
					ui.Row(c).Children(func() { kit.Chip(c, o.Ranges[v.Range], ink) })
				})
			}
		})
	})
	return root
}
