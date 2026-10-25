package devtools

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Column is one table column; Default and Ref are optional.
type Column struct {
	Name, Type, Default, Ref string
	Nullable, Primary        bool
}

// SchemaOptions word the card; Count takes the column count.
type SchemaOptions struct{ Count, Null, Default string }

// Schema shows a table's columns and types.
func Schema(c *ui.Context, table string, columns []Column, o SchemaOptions) ui.Element {
	if strings.TrimSpace(table) == "" || len(columns) == 0 {
		panic(fmt.Sprintf("devtools: Schema %q needs a table name and columns", table))
	}
	names := make([]string, len(columns))
	for i, col := range columns {
		if strings.TrimSpace(col.Type) == "" {
			panic(fmt.Sprintf("devtools: Schema %q column %q needs a type", table, col.Name))
		}
		names[i] = col.Name
	}
	unique("Schema", names)
	def(&o.Count, "%d columns")
	def(&o.Null, "NULL")
	def(&o.Default, "default %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "table-2", table, func() { kit.Chip(c, fmt.Sprintf(o.Count, len(columns)), t.Accent) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5))
		list.Children(func() {
			for i, col := range columns {
				c.Key(col.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Box(c).Size(16, 16).Shrink(0).Children(func() {
						switch {
						case col.Primary:
							glyph(c, "key-round", 16).TextColor(t.Warning)
						case col.Ref != "":
							glyph(c, "link", 16).TextColor(t.Accent)
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, col.Name).Font("monospace").FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).SingleLine()
						if col.Ref != "" || col.Default != "" {
							ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
								if col.Ref != "" {
									glyph(c, "arrow-right", 12).TextColor(t.Accent)
									ui.Text(c, col.Ref).Font("monospace").FontSize(t.Rem(0.7)).TextColor(t.Accent).MinWidth(0).SingleLine()
								}
								if col.Default != "" {
									kit.Muted(c, fmt.Sprintf(o.Default, col.Default)).Font("monospace").FontSize(t.Rem(0.7)).MinWidth(0).SingleLine()
								}
							})
						}
					})
					if col.Nullable {
						ui.Box(c).Padding(t.Space(0.25), t.Space(1.25)).Radius(t.Radius*0.6).Border(1, t.Border).Shrink(0).Children(func() {
							kit.Muted(c, o.Null).Font("monospace").FontSize(t.Rem(0.62)).Shrink(0)
						})
					}
					ui.Box(c).Padding(t.Space(0.5), t.Space(1.5)).Radius(t.Radius * 0.6).Background(t.Accent.Alpha(0.1)).TextColor(t.Accent).
						Shrink(0).Children(func() {
						ui.Text(c, col.Type).Font("monospace").FontSize(t.Rem(0.7)).Bold().Shrink(0).MaxWidth(t.Space(30)).SingleLine()
					})
				})
			}
		})
	})
	return root
}
