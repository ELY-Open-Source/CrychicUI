package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/egoist/mygo/ui"
)

// TopBarOptions tune a TopBar.
type TopBarOptions struct {
	// Back names the back button.
	Back string
}

// TopBar is a title between back and actions; titles cross-fade.
func TopBar(c *ui.Context, title string, o TopBarOptions, actions func()) (root, back ui.Element) {
	if title == "" {
		panic("navigation: TopBar needs a title")
	}
	def(&o.Back, "Back")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Height(t.Space(14)).AlignItems(ui.Center).Gap(t.Space(2)).PaddingX(t.Space(2)).
		Background(t.Surface).Border(1, t.Border).Radius(t.Radius)
	root.Children(func() {
		back = button.IconButton(c, o.Back, "chevron-left")
		ui.Text(c, title).Bold().Grow(1).SingleLine().ContentTransition(ui.ContentTransition{})
		if actions != nil {
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(actions)
		}
	})
	return root, back
}
