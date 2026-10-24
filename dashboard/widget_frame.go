package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WidgetFrameOptions word the frame; Note sits by the title.
type WidgetFrameOptions struct{ Menu, Refresh, Icon, Note string }

// WidgetFrame frames a widget; menu refresh reloads it; cue Tap.
func WidgetFrame(c *ui.Context, title string, o WidgetFrameOptions, content func()) ui.Element {
	if title == "" || content == nil {
		panic(fmt.Sprintf("dashboard: WidgetFrame %q needs a title and content", title))
	}
	def(&o.Menu, "Widget actions")
	def(&o.Refresh, "Refresh")
	def(&o.Icon, "layout-dashboard")
	t := c.Theme()
	root := card(c)
	loads := ui.Local(root, "loads", func() int { return 0 })
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(3)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Icon(c, icon.Must(o.Icon)).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
			ui.Text(c, title).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
			if o.Note != "" {
				kit.Muted(c, o.Note).FontSize(t.Rem(0.7)).Shrink(0)
			}
			c.Key("menu")
			menu := kit.Press(ui.ButtonBase(c).Label(o.Menu).Tooltip(o.Menu).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999))
			menu.Background(menu.AnimateColor("bg", t.SurfaceHover.Alpha(unit(menu.Hovered())), ui.Snappy))
			menu.Children(func() { ui.Icon(c, icon.Must("ellipsis")).Size(16, 16) })
			menu.Menu(func(m *ui.Menu) {
				if m.Item(o.Refresh).Chosen() {
					*loads++
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("WidgetFrame", "refreshed", "loads", *loads)
				}
			})
		})
		c.Key(*loads)
		body := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		kit.Arrive(body, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(content)
	})
	return root
}
