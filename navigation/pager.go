package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PagerOptions tune a Pager.
type PagerOptions struct {
	// Prev, Next name arrows; Format writes the page.
	Prev, Next, Format string
}

// Pager steps pages around "Page x of y".
func Pager(c *ui.Context, label string, page *int, pages int, o PagerOptions) ui.Element {
	if label == "" || page == nil || pages < 1 || *page < 1 || *page > pages {
		panic(fmt.Sprintf("navigation: Pager %q needs a page from 1 to %d", label, pages))
	}
	def(&o.Prev, "Previous")
	def(&o.Next, "Next")
	def(&o.Format, "Page %d of %d")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(2)).FocusGroup(ui.Horizontal).Label(label)
	down := ui.Local(root, "down", func() bool { return false })
	shown := *page
	to := func(p int) {
		*down, *page = p < shown, p
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Pager", "paged", "page", p)
	}
	return root.Children(func() {
		if pageButton(c, o.Prev, "chevron-left", false, shown == 1).Clicked() {
			to(shown - 1)
		}
		kit.Numeric(c, fmt.Sprintf(o.Format, shown, pages)).ContentTransition(ui.ContentTransition{Numeric: true, Down: *down})
		if pageButton(c, o.Next, "chevron-right", false, shown == pages).Clicked() {
			to(shown + 1)
		}
	})
}
