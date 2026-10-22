package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AppSwitcher opens a grid of apps from an icon button.
func AppSwitcher(c *ui.Context, label string, apps []NavItem, picked *int, open *bool) ui.Element {
	if label == "" || len(apps) == 0 || picked == nil || open == nil || *picked < -1 || *picked >= len(apps) {
		panic(fmt.Sprintf("navigation: AppSwitcher %q needs apps, a pick or -1, and open", label))
	}
	for _, a := range apps {
		if a.Label == "" {
			panic(fmt.Sprintf("navigation: AppSwitcher %q apps need labels", label))
		}
		icon.Must(a.Icon)
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	was, shown := *open, *picked
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Size(t.Space(9), t.Space(9)).Center().Radius(t.Radius)).Expanded(*open)
		trigger.Background(trigger.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || trigger.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("grid-3x3")).Size(16, 16)
		})
		popup(c, trigger, open, func() {
			ui.Grid(c).Columns(3).Gap(t.Space(1)).Children(func() {
				for i, a := range apps {
					c.Key(i)
					tile := kit.Press(ui.ButtonBase(c).Label(a.Label).Width(t.Space(22)).Padding(t.Space(2)).Radius(t.Radius).Checked(i == shown))
					bg := t.SurfaceHover.Alpha(unit(tile.Hovered()))
					if i == shown {
						bg = t.Accent.Alpha(0.12)
					}
					tile.Background(tile.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
						ui.Column(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Icon(c, icon.Must(a.Icon)).Size(16, 16)
							ui.Text(c, a.Label).FontSize(t.Rem(0.8)).SingleLine()
						})
					})
					if kit.Clicked(tile, sound.Select) {
						*picked, *open = i, false
						root.SetChanged()
						kit.Log("AppSwitcher", "picked", "index", i)
					}
				}
			})
		})
	})
	reopened(root, was, *open, "AppSwitcher")
	return root
}
