package agent

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SourceCard shows a title, host and snippet; lifts on hover.
func SourceCard(c *ui.Context, n int, src Source) ui.Element {
	if n < 1 {
		panic(fmt.Sprintf("agent: SourceCard %d needs a number from 1", n))
	}
	validSource("SourceCard", src)
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(src.Title).Role(ui.RoleLink).Column().FillWidth().MinWidth(0).AlignItems(ui.Stretch).Gap(t.Space(1.5)).
		Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Shadow(0, 1, 3, 0, kit.Shade(0.06)))
	kit.Lift(root)
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(root.Hovered())), ui.Snappy))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			hostTile(c, src.Host)
			kit.Muted(c, src.Host).FontSize(t.Rem(0.8)).Grow(1).SingleLine()
			kit.Numeric(c, strconv.Itoa(n)).FontSize(t.Rem(0.7)).Bold().TextColor(t.Accent).
				Padding(0, t.Space(1.5)).Radius(999).Background(t.Accent.Alpha(0.12))
		})
		ui.Text(c, src.Title).Bold().Wrap()
		if src.Snippet != "" {
			kit.Muted(c, src.Snippet).FontSize(t.Rem(0.85)).Wrap()
		}
		ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.Accent).Children(func() {
			ui.Icon(c, icon.Must("external-link")).Size(16, 16)
			ui.Text(c, src.URL).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
		})
	})
	if root.Clicked() {
		if ui.SafeToOpen(src.URL) {
			c.OpenURL(src.URL)
		}
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("SourceCard", "opened", "index", n)
	}
	return root
}
