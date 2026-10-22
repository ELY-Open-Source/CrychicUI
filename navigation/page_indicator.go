package navigation

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PageIndicatorOptions tune a PageIndicator.
type PageIndicatorOptions struct {
	// Back and Forth name the arrows.
	Back, Forth string
}

// PageIndicator shows "3 / 10" between arrows; digits roll.
func PageIndicator(c *ui.Context, label string, page *int, n int, o PageIndicatorOptions) ui.Element {
	if label == "" || n < 1 || page == nil || !valid(*page, n) {
		panic(fmt.Sprintf("navigation: PageIndicator %q needs pages and a current one", label))
	}
	def(&o.Back, "Previous page")
	def(&o.Forth, "Next page")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Label(label)
	shown, back, forth := *page, "chevron-left", "chevron-right"
	if c.Direction() == ui.RightToLeft {
		back, forth = forth, back
	}
	return root.Children(func() {
		if arrowButton(c, o.Back, back, shown == 0, sound.Tick) {
			*page = shown - 1
			root.SetChanged()
			kit.Log("PageIndicator", "moved", "page", *page)
		}
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			motion.RollingDigits(c, int64(shown+1)).FontSize(t.Rem(1))
			kit.Muted(c, "/ "+strconv.Itoa(n))
		})
		if arrowButton(c, o.Forth, forth, shown == n-1, sound.Tick) {
			*page = shown + 1
			root.SetChanged()
			kit.Log("PageIndicator", "moved", "page", *page)
		}
	})
}
