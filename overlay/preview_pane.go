package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PreviewPaneOptions tune a PreviewPane.
type PreviewPaneOptions struct {
	// Empty shows while nothing is selected.
	Empty string
}

// PreviewPane shows the selected item; content cross-fades by key.
func PreviewPane(c *ui.Context, label string, key int, o PreviewPaneOptions, content func()) ui.Element {
	if label == "" || content == nil || key < -1 {
		panic(fmt.Sprintf("overlay: PreviewPane %q needs content and a key or -1", label))
	}
	def(&o.Empty, "Select an item to preview")
	t := c.Theme()
	root := ui.Grid(c).Columns(1).FillWidth().Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border).Label(label)
	fade := ui.ElementTransition{Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	return root.Children(func() {
		c.Key(key)
		page := ui.Column(c).ColumnStart(1).RowStart(1).Gap(t.Space(2)).Transition(fade)
		if key < 0 {
			page.Children(func() { kit.Muted(c, o.Empty) })
			return
		}
		page.Children(content)
	})
}
