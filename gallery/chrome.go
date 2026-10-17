package gallery

import (
	"fmt"
	"log/slog"
	"strings"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var modes = []string{"System", "Light", "Dark"}

func (a *App) topBar(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(4)).AlignItems(ui.Center).
		Border(1, t.Border).Children(func() {
		ui.Text(c, "CrychicUI").Bold().FontSize(t.Rem(1.15))
		ui.Spacer(c)
		ui.SearchField(c, &a.query).Label("Search components").Width(240)
		if ui.Segmented(c, &a.mode, modes...).Label("Appearance").Changed() {
			slog.Debug("crychic: gallery appearance", "mode", modes[a.mode])
		}
		on, glyph := sound.Enabled(), "volume-x"
		if on {
			glyph = "volume-2"
		}
		ui.Icon(c, icon.Must(glyph)).Size(16, 16)
		if ui.Switch(c, &on).Label("Sound").Changed() {
			sound.SetEnabled(on)
			slog.Debug("crychic: gallery sound", "on", on)
		}
	})
}

func (a *App) sidebar(c *ui.Context) {
	page := "/"
	if parts := strings.Split(a.router.Path(), "/"); len(parts) > 2 && parts[1] == "c" {
		page = "/c/" + parts[2]
	}
	if ui.Sidebar(c, &page, func() {
		ui.SidebarItem(c, "/", icon.Must("house"), "Home")
		ui.SidebarSection(c, "Components", nil, func() {
			for _, k := range Categories {
				ui.SidebarItem(c, "/c/"+k.Key, icon.Must(k.Icon), k.Name).Children(func() {
					ui.Badge(c, fmt.Sprint(len(inCategory(k.Key))))
				})
			}
		})
	}).Width(220).Changed() {
		a.Go(page)
	}
}

func (a *App) statusBar(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(1.5), t.Space(4)).Border(1, t.Border).
		Children(func() {
			kit.Muted(c, "Log")
			ui.Text(c, strings.Join(a.events, "  ·  ")).MaxLines(1).Grow(1).Shrink(1)
			kit.Muted(c, "Cues")
			var cues []string
			for _, q := range recentCues(8) {
				cues = append(cues, q.String())
			}
			ui.Text(c, strings.Join(cues, " ")).MaxLines(1)
		})
}

func recentCues(n int) []sound.Cue {
	r := sound.Recent()
	return r[max(0, len(r)-n):]
}

func (a *App) results(c *ui.Context) {
	t := c.Theme()
	found := matches(a.query)
	ui.Scroll(c).Grow(1).Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(6)).Children(func() {
			ui.Text(c, fmt.Sprintf("%d results", len(found))).Bold()
			for _, e := range found {
				a.entryRow(c, e)
			}
		})
	})
}

func (a *App) entryRow(c *ui.Context, e Entry) {
	t := c.Theme()
	row := ui.Row(c).Key(e.Path()).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
		Label(e.Name)
	bg := t.Background
	if row.Hovered() {
		bg = t.SurfaceHover
	}
	row.Background(row.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
		kit.Muted(c, fmt.Sprintf("%04d", e.ID)).Width(t.Rem(3))
		ui.Text(c, e.Name).Bold().Width(t.Rem(12))
		kit.Muted(c, e.Summary).MaxLines(1).Grow(1).Shrink(1)
	})
	if row.Clicked() {
		a.Go(e.Path())
	}
}
