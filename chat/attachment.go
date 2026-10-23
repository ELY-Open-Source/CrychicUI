package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// Attachment is a file chip showing its size.
func Attachment(c *ui.Context, name string, bytes int64) ui.Element {
	if name == "" || bytes < 0 {
		panic(fmt.Sprintf("chat: Attachment %q needs a name and a size", name))
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(name).Row().AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).
		Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Border(1, t.Border))
	root.Background(root.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(root.Hovered())), ui.Snappy))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).Background(t.Accent.Alpha(0.12)).
			TextColor(t.Accent).Children(func() { ui.Icon(c, icon.Must("file")).Size(16, 16) })
		ui.Column(c).Children(func() {
			ui.Text(c, name).Bold().SingleLine()
			typography.Bytes(c, bytes).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		})
	})
	if kit.Clicked(root, sound.Tap) {
		kit.Log("Attachment", "opened", "bytes", bytes)
	}
	return root
}
