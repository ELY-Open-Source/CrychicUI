package agent

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Citation is a numbered chip opening a source.
func Citation(c *ui.Context, n int, url string) ui.Element {
	if n < 1 || url == "" {
		panic(fmt.Sprintf("agent: Citation %d needs a number and a URL", n))
	}
	t := c.Theme()
	name := strconv.Itoa(n)
	root := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(url).AlignSelf(ui.Start).Center().
		Padding(t.Space(0.5), t.Space(1.5)).Radius(999).Border(1, t.Border))
	root.Background(root.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(root.Hovered())), ui.Snappy)).TextColor(t.Accent)
	kit.Lift(root)
	root.Children(func() { ui.Text(c, name).FontSize(t.Rem(0.75)).Bold().SingleLine() })
	if root.Clicked() {
		if ui.SafeToOpen(url) {
			c.OpenURL(url)
		}
		root.SetChanged()
		sound.Play(sound.Open)
		kit.Log("Citation", "opened", "index", n)
	}
	return root
}
