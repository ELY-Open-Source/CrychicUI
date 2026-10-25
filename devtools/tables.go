package devtools

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// Table is one database table; Size is in bytes.
type Table struct {
	Name       string
	Rows, Size int64
}

// TablesOptions word the panel; Row and Rows take a count.
type TablesOptions struct{ Title, Empty, Row, Rows string }

// Tables lists database tables with rows and size; cue Select.
func Tables(c *ui.Context, tables []Table, selected *string, o TablesOptions) ui.Element {
	def(&o.Title, "Tables")
	if selected == nil {
		panic(fmt.Sprintf("devtools: Tables %q needs a selection", o.Title))
	}
	names := make([]string, len(tables))
	var largest int64
	for i, tb := range tables {
		if tb.Rows < 0 || tb.Size < 0 {
			panic(fmt.Sprintf("devtools: Tables %q needs non-negative rows and size", tb.Name))
		}
		names[i] = tb.Name
		largest = max(largest, tb.Size)
	}
	if len(tables) > 0 {
		unique("Tables", names)
	}
	def(&o.Empty, "No tables")
	def(&o.Row, "%s row")
	def(&o.Rows, "%s rows")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "database", o.Title, func() { kit.Chip(c, fmt.Sprint(len(tables)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(tables) == 0 {
				empty(c, "database", o.Empty)
			}
			for i, tb := range tables {
				on := *selected == tb.Name
				c.Key(tb.Name)
				b := kit.Press(ui.ButtonBase(c).Label(tb.Name).FillWidth().Padding(t.Space(2)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						ink := t.TextMuted
						if on {
							ink = t.Accent
						}
						tile(c, "table-2", ink, t.Space(8))
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							ui.Text(c, tb.Name).Font("monospace").FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).SingleLine()
							format := o.Rows
							if tb.Rows == 1 {
								format = o.Row
							}
							kit.Muted(c, fmt.Sprintf(format, number(c, tb.Rows))).FontSize(t.Rem(0.7)).FillWidth().MinWidth(0).SingleLine()
						})
						ui.Column(c).Width(t.Space(18)).Shrink(0).Gap(t.Space(0.75)).AlignItems(ui.End).Children(func() {
							typography.Bytes(c, tb.Size).FontSize(t.Rem(0.72)).Shrink(0)
							share := float32(0)
							if largest > 0 {
								share = float32(float64(tb.Size) / float64(largest))
							}
							meter(c, share, ink).Height(t.Space(1))
						})
					})
				})
				if b.Clicked() && !on {
					*selected = tb.Name
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Tables", "selected", "index", i)
				}
			}
		})
	})
	return root
}
