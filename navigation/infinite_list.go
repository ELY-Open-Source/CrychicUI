package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// InfiniteListOptions tune an InfiniteList.
type InfiniteListOptions struct {
	// Loading names the spinner row.
	Loading string
}

// InfiniteList asks for more near its end.
func InfiniteList(c *ui.Context, label string, count int, loading bool, item func(i int), more func(), o InfiniteListOptions) ui.Element {
	if label == "" || count < 0 || item == nil || more == nil {
		panic(fmt.Sprintf("navigation: InfiniteList %q needs a count, an item builder and more", label))
	}
	def(&o.Loading, "Loading")
	t := c.Theme()
	root := ui.Scroll(c).Label(label)
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	asked := ui.Local(root, "asked", func() int { return -1 })
	laid := ui.Local(root, "laid", func() int { return -1 })
	root.TrackScroll(st)
	var body ui.Element
	root.Children(func() {
		body = ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2))
		body.Children(func() {
			for i := range count {
				c.Key(i)
				ui.Column(c).FillWidth().Transition(kit.Enter).Children(func() { item(i) })
			}
			if loading {
				c.Key("loading")
				ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2)).AlignItems(ui.Center).Transition(kit.Enter).Children(func() {
					feedback.Spin(c, o.Loading, feedback.SpinOptions{})
					kit.Muted(c, o.Loading)
				})
			}
		})
	})
	// Paint follows the layout of the final build.
	body.DrawOver(func(*ui.Painter, ui.Rect) {
		if *laid != count {
			*laid = count
			c.Invalidate()
		}
	})
	measured := body.LayoutBounds().H
	if measured > 0 && *laid == count && !loading && count != *asked && st.MaxY-st.Y < t.Space(30) {
		*asked = count
		kit.Log("InfiniteList", "more", "count", count)
		more()
		c.Invalidate()
	}
	kit.Settle(c, root, []float32{st.Y, st.MaxY, measured}, func() []float32 { return []float32{st.Y, st.MaxY, body.LayoutBounds().H} })
	return root
}
