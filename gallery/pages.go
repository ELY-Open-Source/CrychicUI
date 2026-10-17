package gallery

import (
	"fmt"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

func (a *App) pages(c *ui.Context) {
	a.router.View(c, func(r *ui.Route) {
		switch {
		case r.Match("/"):
			r.Title("CrychicUI")
			a.home(c)
		case r.Match("/c/{cat}"):
			k, ok := category(r.Param("cat"))
			if !ok {
				notFound(c)
				return
			}
			r.Title(k.Name)
			a.categoryPage(c, k)
		case r.Match("/c/{cat}/{slug}"):
			for _, e := range inCategory(r.Param("cat")) {
				if e.Slug() == r.Param("slug") {
					r.Title(e.Name)
					a.entryPage(c, e)
					return
				}
			}
			notFound(c)
		default:
			notFound(c)
		}
	})
}

func page(c *ui.Context, fn func()) {
	t := c.Theme()
	ui.Scroll(c).Grow(1).Children(func() {
		ui.Column(c).FillWidth().MaxWidth(t.Rem(64)).Gap(t.Space(4)).Padding(t.Space(8)).Children(fn)
	})
}

func notFound(c *ui.Context) {
	page(c, func() {
		ui.Text(c, "Not found").Bold().FontSize(c.Theme().Rem(1.5))
		kit.Muted(c, "No component lives at this address.")
	})
}

func (a *App) home(c *ui.Context) {
	t := c.Theme()
	page(c, func() {
		ui.Text(c, "CrychicUI").Bold().FontSize(t.Rem(2))
		kit.Muted(c, "Components for MyGo native UI, in MyGo's own look, with motion and sound.")
		ui.Row(c).Gap(t.Space(4)).Children(func() {
			stat(c, "Components", len(entries))
			stat(c, "Categories", len(Categories))
			stat(c, "Sound cues", len(sound.Cues()))
		})
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).Children(func() {
			for _, k := range Categories {
				a.categoryCard(c, k)
			}
		})
	})
}

func stat(c *ui.Context, label string, n int) {
	t := c.Theme()
	kit.Card(c).Padding(t.Space(3), t.Space(4)).Gap(t.Space(1)).Children(func() {
		kit.Muted(c, label)
		ui.Text(c, fmt.Sprint(n)).Bold().FontSize(t.Rem(1.5))
	})
}

func (a *App) categoryCard(c *ui.Context, k Category) {
	t := c.Theme()
	card := kit.Press(kit.Lift(kit.Card(c).Key(k.Key).Width(t.Rem(13)).Padding(t.Space(3)).Gap(t.Space(1)).
		Label(k.Name)))
	card.Children(func() {
		ui.Icon(c, icon.Must(k.Icon)).Size(16, 16)
		ui.Text(c, k.Name).Bold()
		kit.Muted(c, fmt.Sprintf("%d components", len(inCategory(k.Key))))
	})
	if card.Clicked() {
		a.Go("/c/" + k.Key)
	}
}

func (a *App) categoryPage(c *ui.Context, k Category) {
	page(c, func() {
		ui.Text(c, k.Name).Bold().FontSize(c.Theme().Rem(1.5))
		list := inCategory(k.Key)
		kit.Muted(c, fmt.Sprintf("%d components", len(list)))
		for _, e := range list {
			a.entryRow(c, e)
		}
	})
}

func (a *App) entryPage(c *ui.Context, e Entry) {
	t := c.Theme()
	page(c, func() {
		kit.Muted(c, fmt.Sprintf("%04d · %s", e.ID, e.Category))
		ui.Text(c, e.Name).Bold().FontSize(t.Rem(1.5))
		kit.Muted(c, e.Summary)
		kit.Card(c).Key(e.ID).FillWidth().Padding(t.Space(6)).Children(func() {
			e.Demo(c, a.logger(e))
		})
		ui.Text(c, "Code").Bold()
		editor.View(c, a.code(e)).FillWidth().Height(t.Rem(22))
	})
}
