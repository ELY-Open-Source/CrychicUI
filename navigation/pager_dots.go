package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PagerDotsOptions tune PagerDots.
type PagerDotsOptions struct {
	// Page names a dot; %d is its number.
	Page string
}

// PagerDots are dots for pages; the active one widens.
func PagerDots(c *ui.Context, label string, page *int, n int, o PagerDotsOptions) ui.Element {
	if label == "" || n < 2 || page == nil || !valid(*page, n) {
		panic(fmt.Sprintf("navigation: PagerDots %q needs two pages and a current one", label))
	}
	def(&o.Page, "Page %d")
	t := c.Theme()
	shown := *page
	parts := ui.TabsBase(c, page, n)
	root := parts.List.Label(label).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center)
	root.Children(func() {
		for i := range n {
			c.Key(i)
			tab := parts.Tab(i).Label(fmt.Sprintf(o.Page, i+1)).Padding(t.Space(1))
			tab.Children(func() { pip(c, ui.Box(c), i == shown) })
		}
	})
	if kit.Changed(root, sound.Tick) {
		kit.Log("PagerDots", "picked", "index", *page)
	}
	return root
}
