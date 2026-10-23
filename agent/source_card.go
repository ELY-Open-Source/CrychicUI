package agent

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SourceCard shows a title, host and snippet; lifts on hover.
func SourceCard(c *ui.Context, title, host, snippet string) ui.Element {
	if title == "" || host == "" {
		panic("agent: SourceCard needs a title and a host")
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(title).Column().AlignSelf(ui.Start).Width(t.Space(76)).Gap(t.Space(1)).
		Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Shadow(0, 1, 3, 0, kit.Shade(0.06)))
	kit.Lift(root)
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(root.Hovered())), ui.Snappy))
	root.Children(func() {
		ui.Text(c, title).Bold().Wrap()
		kit.Muted(c, host).FontSize(t.Rem(0.8)).SingleLine()
		if snippet != "" {
			kit.Muted(c, snippet).FontSize(t.Rem(0.9)).Wrap()
		}
	})
	if root.Clicked() {
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("SourceCard", "opened")
	}
	return root
}
