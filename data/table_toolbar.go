package data

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TableToolbarOptions tune a TableToolbar.
type TableToolbarOptions struct {
	Search, Filters, Columns, Density string
	// Densities label compact, default and roomy.
	Densities [3]string
}

// TableToolbar holds search, filters, columns and density.
func TableToolbar(c *ui.Context, label string, query *string, d *Density, o TableToolbarOptions) (root, filters, columns ui.Element) {
	if label == "" || query == nil || d == nil {
		panic("data: TableToolbar needs a label, a query and a density")
	}
	def(&o.Search, "Search")
	def(&o.Filters, "Filters")
	def(&o.Columns, "Columns")
	def(&o.Density, "Row density")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleToolbar).Label(label)
	tool := func(name, glyph string) ui.Element {
		b := kit.Press(ui.ButtonBase(c.Key(name)).Label(name).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3)).
			Radius(t.Radius).Border(1, t.Border))
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(16, 16)
			ui.Text(c, name)
		})
		if b.Clicked() {
			sound.Play(sound.Tap)
			kit.Log("TableToolbar", "tool", "name", name)
		}
		return b
	}
	root.Children(func() {
		if ui.SearchField(c, query).Label(o.Search).Grow(1).Changed() {
			root.SetChanged()
			kit.Log("TableToolbar", "searched", "length", len([]rune(*query)))
		}
		filters = tool(o.Filters, "list-filter")
		columns = tool(o.Columns, "columns-3")
		if DensitySwitch(c, o.Density, d, DensitySwitchOptions{Names: o.Densities}).Changed() {
			root.SetChanged()
		}
	})
	return root, filters, columns
}
