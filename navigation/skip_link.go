package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SkipLink opens while focused and moves focus to target.
func SkipLink(c *ui.Context, label string, target *ui.Handle) (link ui.Element) {
	if label == "" || target == nil {
		panic(fmt.Sprintf("navigation: SkipLink %q needs a label and a target", label))
	}
	t := c.Theme()
	// In flow: Tab reaches absolute elements last.
	bar := ui.Box(c).FillWidth().Clip()
	bar.Children(func() {
		link = kit.Press(ui.ButtonBase(c).Label(label).Height(t.Space(9)).Margin(t.Space(1)).Padding(0, t.Space(3)).Center().
			Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText))
		link.Children(func() { ui.Text(c, label).SingleLine() })
		if kit.Clicked(link, sound.Tap) {
			target.Focus(c)
			kit.Log("SkipLink", "jumped")
		}
	})
	k := link.AnimateTo("open", unit(link.Focused()), ui.Snappy)
	link.Opacity(k)
	bar.Height(k * t.Space(11))
	return link
}
