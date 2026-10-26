package desktop

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LaunchApp is an app in the launcher grid.
type LaunchApp struct{ Name, Icon, Category string }

// LauncherOptions word the search.
type LauncherOptions struct{ Placeholder, Empty string }

// Launcher is an app grid filtered by query; cue Select.
func Launcher(c *ui.Context, apps []LaunchApp, query, picked *string, o LauncherOptions) ui.Element {
	def(&o.Placeholder, "Search apps")
	def(&o.Empty, "No apps match")
	names := make([]string, len(apps))
	for i, a := range apps {
		if a.Icon == "" {
			panic(fmt.Sprintf("desktop: Launcher app %q needs an icon", a.Name))
		}
		names[i] = a.Name
	}
	unique("Launcher", names)
	if query == nil || picked == nil {
		panic("desktop: Launcher needs a query and a picked app")
	}
	t := c.Theme()
	root := window(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		frame := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(999).Background(t.SurfaceHover)
		frame.Border(1, frame.AnimateColor("ring", t.SurfaceHover.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
		frame.Children(func() {
			glyph(c, "search", 16).TextColor(t.TextMuted)
			in := ui.TextInputBase(c, query).Label(o.Placeholder).Placeholder(o.Placeholder).FontSize(t.Rem(0.86)).Grow(1).MinWidth(0)
			if in.Changed() {
				root.SetChanged()
				kit.Log("Launcher", "query", "runes", utf8.RuneCountInString(*query))
			}
		})
		q := strings.TrimSpace(*query)
		var shown []LaunchApp
		for _, a := range apps {
			lo, _ := kit.IndexFold(a.Name, q)
			lc, _ := kit.IndexFold(a.Category, q)
			if q == "" || lo >= 0 || lc >= 0 {
				shown = append(shown, a)
			}
		}
		if len(shown) == 0 {
			ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(8), 0).Children(func() {
				kit.Arrive(glyph(c, "search-x", 24).TextColor(t.TextMuted), ui.Pose{Scale: 0.7, Fade: 1}, 0)
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.84)).FillWidth().TextAlign(ui.Center).Wrap()
			})
			return
		}
		grid := ui.Grid(c).Columns(4).FillWidth().Gap(t.Space(2))
		grid.Children(func() {
			for i, a := range shown {
				c.Key(strconv.Quote(a.Name))
				on := a.Name == *picked
				b := kit.Press(ui.ButtonBase(c).Label(a.Name).Checked(on).Padding(t.Space(2), t.Space(1)).Radius(t.Radius * 1.5).AlignItems(ui.Center))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				kit.Stagger(c, grid, b, i).Children(func() {
					ui.Column(c).Gap(t.Space(1.5)).AlignItems(ui.Center).MinWidth(0).Children(func() {
						tile(c, a.Icon, t.Accent, t.Space(12)).Radius(t.Radius * 2)
						ui.Text(c, a.Name).FontSize(t.Rem(0.74)).MinWidth(0).MaxWidth(t.Space(22)).TextAlign(ui.Center).SingleLine()
					})
				})
				if b.Clicked() && !on {
					*picked = a.Name
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Launcher", "picked")
				}
			}
		})
	})
	return root
}
