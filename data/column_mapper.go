package data

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ColumnMapperOptions tune a ColumnMapper.
type ColumnMapperOptions struct {
	// Skip names leaving a column unmapped.
	Skip string
	// Map names a column's menu; it takes the column.
	Map string
}

// ColumnMapper maps source columns to fields; "" skips one.
func ColumnMapper(c *ui.Context, label string, sources, fields []string, mapping *[]string, o ColumnMapperOptions) ui.Element {
	if label == "" || len(sources) == 0 || len(fields) == 0 || mapping == nil || len(*mapping) != len(sources) {
		panic(fmt.Sprintf("data: ColumnMapper %q needs sources, fields and a mapping per source", label))
	}
	unique("ColumnMapper", label, sources)
	unique("ColumnMapper", label, fields)
	def(&o.Skip, "Skip")
	def(&o.Map, "Map %s")
	if slices.Contains(fields, o.Skip) {
		panic(fmt.Sprintf("data: ColumnMapper %q fields cannot include %q", label, o.Skip))
	}
	for _, m := range *mapping {
		if m != "" && !slices.Contains(fields, m) {
			panic(fmt.Sprintf("data: ColumnMapper %q maps to %q, not a field", label, m))
		}
	}
	t := c.Theme()
	options := append([]string{o.Skip}, fields...)
	root := frame(c, label).Dividers(1, t.Border)
	return root.Children(func() {
		for i, src := range sources {
			c.Key(src)
			mapped := (*mapping)[i] != ""
			ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
				ui.Text(c, src).Grow(1).SingleLine()
				arrow := ui.Icon(c, icon.Must("arrow-right")).Size(16, 16)
				arrow.TextColor(arrow.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(mapped)), ui.Snappy))
				pick := (*mapping)[i]
				if pick == "" {
					pick = o.Skip
				}
				if choice.Dropdown(c, fmt.Sprintf(o.Map, src), options, &pick).Width(t.Space(44)).Changed() {
					(*mapping)[i] = pick
					if pick == o.Skip {
						(*mapping)[i] = ""
					}
					root.SetChanged()
					n := 0
					for _, m := range *mapping {
						if m != "" {
							n++
						}
					}
					kit.Log("ColumnMapper", "mapped", "index", i, "mapped", n)
				}
			})
		}
	})
}
