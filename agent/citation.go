package agent

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Citation is a numbered chip opening a source.
func Citation(c *ui.Context, n int, src Source) ui.Element {
	if n < 1 {
		panic(fmt.Sprintf("agent: Citation %d needs a number from 1", n))
	}
	validSource("Citation", src)
	t := c.Theme()
	name := strconv.Itoa(n)
	root := kit.Press(ui.ButtonBase(c).Label(name+" "+src.Title).Tooltip(src.Title+" · "+src.Host).AlignSelf(ui.Center).
		Center().MinWidth(t.Space(5)).Padding(t.Space(0.25), t.Space(1.25)).Radius(999))
	root.Background(root.AnimateColor("bg", t.Accent.Alpha(0.12).Mix(t.Accent, unit(root.Hovered())), ui.Snappy))
	root.TextColor(root.AnimateColor("ink", t.Accent.Mix(t.AccentText, unit(root.Hovered())), ui.Snappy))
	kit.Lift(root)
	root.Children(func() { ui.Text(c, name).FontSize(t.Rem(0.7)).Bold().SingleLine() })
	if root.Clicked() {
		if ui.SafeToOpen(src.URL) {
			c.OpenURL(src.URL)
		}
		root.SetChanged()
		sound.Play(sound.Open)
		kit.Log("Citation", "opened", "index", n)
	}
	return root
}
