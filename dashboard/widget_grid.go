package dashboard

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WidgetGridOptions word the grid; Hint sits under each title.
type WidgetGridOptions struct{ Hint string }

// WidgetGrid lays widgets in two columns, dragged; cue Drop.
func WidgetGrid(c *ui.Context, titles *[]string, o WidgetGridOptions) ui.Element {
	if titles == nil || len(*titles) == 0 {
		panic("dashboard: WidgetGrid needs widgets")
	}
	unique("WidgetGrid", *titles)
	def(&o.Hint, "Drag to reorder")
	t := c.Theme()
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	gridID := ui.Local(root, "gridID", func() *int { return new(int) })
	moveTitle, moveTo := "", -1
	root.Children(func() {
		for i, title := range *titles {
			c.Key(title)
			tile := kit.Card(c).Row().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2.5)).MinWidth(0).
				Focusable().Label(title).Cursor(ui.CursorGrab).Drag(widgetDrag{id: *gridID, title: title}).
				Transition(ui.ElementTransition{Animation: &ui.Snappy, Position: true})
			over := false
			if d, ok := ui.DragOver[widgetDrag](tile); ok && d.id == *gridID && d.title != title {
				over = true
			}
			lift := unit(tile.Hovered()) * 0.5
			tile.Background(tile.AnimateColor("bg", t.Background.Mix(t.Surface, lift).Mix(t.Accent.Alpha(0.1), unit(over)).
				Mix(t.Surface, unit(tile.Dragging())), ui.Snappy)).
				BorderColor(tile.AnimateColor("edge", t.Border.Mix(t.Accent, unit(over)), ui.Snappy))
			kit.Stagger(c, root, tile, i).Children(func() {
				ink := palette(t, i)
				ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink).
					Children(func() { ui.Icon(c, icon.Must("layout-dashboard")).Size(14, 14) })
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
					ui.Text(c, title).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
					kit.Muted(c, o.Hint).FontSize(t.Rem(0.7)).FillWidth().SingleLine()
				})
				ui.Icon(c, icon.Must("grip-vertical")).Size(16, 16).Shrink(0).TextColor(t.TextMuted)
			})
			if d, ok := ui.Drop[widgetDrag](tile); ok && d.id == *gridID {
				moveTitle, moveTo = d.title, i
			}
		}
	})
	if moveTitle != "" {
		from := slices.Index(*titles, moveTitle)
		if from >= 0 && from != moveTo {
			*titles = slices.Delete(*titles, from, from+1)
			*titles = slices.Insert(*titles, moveTo, moveTitle)
			root.SetChanged()
			sound.Play(sound.Drop)
			kit.Log("WidgetGrid", "moved", "from", from, "to", moveTo)
		}
	}
	return root
}

type widgetDrag struct {
	id    *int
	title string
}
