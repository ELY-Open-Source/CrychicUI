package desktop

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TrayItem is a tray command.
type TrayItem struct{ Label, Icon, Shortcut string }

// TrayMenuOptions name the tray icon.
type TrayMenuOptions struct{ Label, Icon string }

// TrayMenu is a tray command menu; opening reports; cue Open.
func TrayMenu(c *ui.Context, items []TrayItem, chosen *string, o TrayMenuOptions) ui.Element {
	def(&o.Label, "Tray")
	def(&o.Icon, "app-window")
	labels := make([]string, len(items))
	for i, it := range items {
		labels[i] = it.Label
	}
	unique("TrayMenu", labels)
	if len(items) == 0 || chosen == nil {
		panic(fmt.Sprintf("desktop: TrayMenu %q needs items and a chosen command", o.Label))
	}
	t := c.Theme()
	root := ui.Column(c).AlignItems(ui.Start).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(o.Label).Role(ui.RoleMenuButton).Expanded(*open).Padding(t.Space(1.5)).Radius(t.Radius))
		trigger.Background(trigger.AnimateColor("bg", t.Border.Alpha(0.6*max(unit(trigger.Hovered()), unit(*open))), ui.Snappy))
		trigger.Children(func() { glyph(c, o.Icon, 16) })
		if trigger.Clicked() {
			*open = !*open
			if *open {
				root.SetChanged()
				sound.Play(sound.Open)
				kit.Log("TrayMenu", "opened")
			}
		}
		ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(56)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(0.5)).Role(ui.RoleMenu).Children(func() {
				for _, it := range items {
					c.Key(strconv.Quote(it.Label))
					item := ui.ButtonBase(c).Label(it.Label).Role(ui.RoleMenuItem).FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)
					item.Background(item.AnimateColor("bg", t.Accent.Alpha(unit(item.Hovered())), ui.Snappy))
					fg := t.Text.Mix(t.AccentText, unit(item.Hovered()))
					item.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							if it.Icon != "" {
								glyph(c, it.Icon, 14).TextColor(fg)
							}
							ui.Text(c, it.Label).FontSize(t.Rem(0.8)).TextColor(fg).Grow(1).MinWidth(0).SingleLine()
							if it.Shortcut != "" {
								ui.Text(c, it.Shortcut).FontSize(t.Rem(0.74)).TextColor(fg.Alpha(0.6)).Shrink(0).MaxWidth(t.Space(24)).SingleLine()
							}
						})
					})
					if item.Clicked() {
						*open = false
						*chosen = it.Label
						root.SetChanged()
						sound.Play(sound.Tap)
						kit.Log("TrayMenu", "chose")
					}
				}
			})
		})
	})
	return root
}
