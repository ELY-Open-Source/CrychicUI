package dashboard

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WidgetGrid reorders widgets by drag; cue Drop.
func WidgetGrid(c *ui.Context, titles *[]string) ui.Element {
	if titles == nil || len(*titles) == 0 {
		panic("dashboard: WidgetGrid needs widgets")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	gridID := ui.Local(root, "gridID", func() *int { return new(int) })
	moveTitle, moveTo := "", -1
	root.Children(func() {
		for i, title := range *titles {
			if title == "" {
				panic("dashboard: WidgetGrid title must not be empty")
			}
			c.Key(title)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Focusable().Label(title).Cursor(ui.CursorGrab).Drag(widgetDrag{id: *gridID, title: title})
			if row.Dragging() {
				row.Opacity(0.4)
			}
			row.Children(func() {
				ui.Icon(c, icon.Must("grip-vertical")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, title).Grow(1)
			})
			if d, ok := ui.Drop[widgetDrag](row); ok && d.id == *gridID {
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
