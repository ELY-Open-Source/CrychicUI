package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WidgetFrameOptions are the widget labels.
type WidgetFrameOptions struct{ Menu, Refresh string }

// WidgetFrame frames a widget with menu and refresh; cue Tap.
func WidgetFrame(c *ui.Context, title string, o WidgetFrameOptions, content func()) ui.Element {
	if title == "" || content == nil {
		panic("dashboard: WidgetFrame needs title and content")
	}
	def(&o.Menu, "Widget actions")
	def(&o.Refresh, "Refresh")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, title).Bold().Grow(1)
			menu := kit.Press(ui.ButtonBase(c).Label(o.Menu).Tooltip(o.Menu).Size(t.Space(7), t.Space(7)).Center().Radius(999))
			menu.Children(func() { ui.Icon(c, icon.Must("more-horizontal")).Size(16, 16) })
			menu.Menu(func(m *ui.Menu) {
				if m.Item(o.Refresh).Chosen() {
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("WidgetFrame", "refreshed")
				}
			})
		})
		content()
	})
	return root
}
