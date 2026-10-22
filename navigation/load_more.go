package navigation

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LoadMoreOptions tune a LoadMore.
type LoadMoreOptions struct {
	// Step is shown per click; 10 when zero.
	Step int
	// More is the button's copy.
	More string
}

// LoadMore shows items and a button for more.
func LoadMore(c *ui.Context, label string, shown *int, total int, item func(i int), o LoadMoreOptions) ui.Element {
	if o.Step == 0 {
		o.Step = 10
	}
	if label == "" || shown == nil || item == nil || total < 0 || *shown < 0 || *shown > total || o.Step < 1 {
		panic(fmt.Sprintf("navigation: LoadMore %q needs 0 <= shown <= total, an item builder and Step >= 1", label))
	}
	def(&o.More, "Show more")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Label(label)
	start := ui.Local(root, "start", func() int { return 0 })
	return root.Children(func() {
		for i := range *shown {
			c.Key(i)
			kit.Arrive(ui.Column(c).FillWidth(), ui.Pose{Y: 6, Fade: 1}, time.Duration(min(max(i-*start, 0), 10))*40*time.Millisecond).
				Children(func() { item(i) })
		}
		if *shown < total {
			c.Key("more")
			if button.Outline(c, o.More).AlignSelf(ui.Center).Transition(kit.Enter).Clicked() {
				*start, *shown = *shown, min(total, *shown+o.Step)
				root.SetChanged()
				kit.Log("LoadMore", "more", "shown", *shown)
			}
		}
	})
}
