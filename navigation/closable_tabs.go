package navigation

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ClosableTabsOptions tune ClosableTabs.
type ClosableTabsOptions struct {
	// Add names the add button; Close and Untitled format names.
	Add, Close, Untitled string
}

// tabKey tells apart tabs of one name by their order.
type tabKey struct {
	name string
	n    int
}

// ClosableTabs are tabs that close and add.
func ClosableTabs(c *ui.Context, label string, names *[]string, active *int, o ClosableTabsOptions) ui.Element {
	if label == "" || names == nil || active == nil || len(*names) == 0 && *active != -1 || len(*names) > 0 && !valid(*active, len(*names)) {
		panic("navigation: ClosableTabs needs a label, tabs and an active index, -1 when empty")
	}
	def(&o.Add, "New tab")
	def(&o.Close, "Close %s")
	def(&o.Untitled, "Untitled %d")
	t := c.Theme()
	root := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
	shown, closing, adding := *active, -1, false
	pop := ui.ElementTransition{Enter: &ui.Motion{Scale: 0.8}, Exit: &ui.Motion{Scale: 0.8}}
	root.Children(func() {
		if len(*names) > 0 {
			parts := ui.TabsBase(c, active, len(*names))
			list := parts.List.Label(label).Gap(t.Space(1))
			list.Children(func() {
				for i, name := range *names {
					c.Key(tabKey{name, i - slices.Index(*names, name)})
					tab := parts.Tab(i).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).Radius(t.Radius).Transition(pop)
					bg := t.Surface.Alpha(0)
					if i == shown {
						bg = t.Surface
					} else if tab.Hovered() {
						bg = t.SurfaceHover
					}
					tab.Background(tab.AnimateColor("bg", bg, ui.Snappy)).Border(1, tab.AnimateColor("edge", t.Border.Alpha(unit(i == shown)), ui.Snappy))
					tab.Children(func() {
						ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, name).SingleLine()
							x := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Close, name)).Size(t.Space(6), t.Space(6)).Center().Radius(t.Radius))
							x.Background(x.AnimateColor("bg", t.SurfaceHover.Alpha(unit(x.Hovered())), ui.Snappy)).Children(func() {
								ui.Icon(c, icon.Must("x")).Size(16, 16)
							})
							if x.Clicked() {
								closing = i
							}
						})
					})
				}
			})
			if list.Changed() && closing < 0 {
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ClosableTabs", "picked", "index", *active)
			}
		}
		add := kit.Press(ui.ButtonBase(c).Label(o.Add).Tooltip(o.Add).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius))
		add.Background(add.AnimateColor("bg", t.SurfaceHover.Alpha(unit(add.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("plus")).Size(16, 16)
		})
		adding = add.Clicked()
	})
	switch {
	case closing >= 0:
		*names = slices.Delete(*names, closing, closing+1)
		if closing < *active || *active >= len(*names) {
			*active--
		}
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("ClosableTabs", "closed", "tabs", len(*names))
	case adding:
		*names = append(*names, fmt.Sprintf(o.Untitled, len(*names)+1))
		*active = len(*names) - 1
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("ClosableTabs", "added", "tabs", len(*names))
	}
	return root
}
