package desktop

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MenuItem is a command with an optional shortcut.
type MenuItem struct{ Label, Shortcut string }

// Menu is a titled list of commands.
type Menu struct {
	Title string
	Items []MenuItem
}

// MenuBar is an app menu bar; opening reports; cue Open.
func MenuBar(c *ui.Context, app string, menus []Menu, chosen *string) ui.Element {
	titles, labels := make([]string, len(menus)), []string{}
	for i, m := range menus {
		if len(m.Items) == 0 {
			panic(fmt.Sprintf("desktop: MenuBar menu %q needs items", m.Title))
		}
		for _, it := range m.Items {
			labels = append(labels, it.Label)
		}
		titles[i] = m.Title
	}
	unique("MenuBar", titles)
	// chosen names a command, so labels are unique across menus.
	unique("MenuBar", labels)
	if strings.TrimSpace(app) == "" || chosen == nil {
		panic(fmt.Sprintf("desktop: MenuBar %q needs an app name and a chosen command", app))
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(0.5)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).
		Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Role(ui.RoleMenuBar)
	kit.Arrive(root, ui.Pose{Y: -t.Space(1), Fade: 1}, 0)
	opens := ui.Local(root, "opens", func() map[string]*bool { return map[string]*bool{} })
	root.Children(func() {
		ui.Text(c, app).FontSize(t.Rem(0.8)).Bold().Padding(t.Space(1), t.Space(2)).Shrink(0).MaxWidth(t.Space(36)).SingleLine()
		for _, m := range menus {
			c.Key(strconv.Quote(m.Title))
			open := (*opens)[m.Title]
			if open == nil {
				open = new(bool)
				(*opens)[m.Title] = open
			}
			trigger := ui.ButtonBase(c).Label(m.Title).Role(ui.RoleMenuButton).Expanded(*open).Padding(t.Space(1), t.Space(2)).Radius(t.Radius)
			trigger.Background(trigger.AnimateColor("bg", t.Border.Alpha(0.6*max(unit(trigger.Hovered()), unit(*open))), ui.Snappy))
			trigger.Children(func() {
				ui.Text(c, m.Title).FontSize(t.Rem(0.8)).Shrink(0).MaxWidth(t.Space(30)).SingleLine()
			})
			if trigger.Clicked() {
				*open = !*open
				if *open {
					root.SetChanged()
					sound.Play(sound.Open)
					kit.Log("MenuBar", "opened")
				}
			}
			ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
				panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(52)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
					Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0, 0)
				kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
				ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(0.5)).Role(ui.RoleMenu).Children(func() {
					for _, it := range m.Items {
						c.Key(strconv.Quote(it.Label))
						item := ui.ButtonBase(c).Label(it.Label).Role(ui.RoleMenuItem).FillWidth().Padding(t.Space(1.25), t.Space(2)).Radius(t.Radius)
						item.Background(item.AnimateColor("bg", t.Accent.Alpha(unit(item.Hovered())), ui.Snappy))
						fg := t.Text.Mix(t.AccentText, unit(item.Hovered()))
						item.Children(func() {
							ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Children(func() {
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
							kit.Log("MenuBar", "chose")
						}
					}
				})
			})
		}
	})
	return root
}
