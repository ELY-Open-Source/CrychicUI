package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChatHeaderOptions name the call buttons.
type ChatHeaderOptions struct {
	Voice, Video string
}

// ChatHeader shows avatar, name, status and call buttons.
func ChatHeader(c *ui.Context, name, status string, image *ui.Bitmap, o ChatHeaderOptions) (root, voice, video ui.Element) {
	if name == "" || status == "" {
		panic(fmt.Sprintf("chat: ChatHeader %q needs a name and a status", name))
	}
	def(&o.Voice, "Voice call")
	def(&o.Video, "Video call")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Role(ui.RoleGroup).Label(name)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
			display.Avatar(c, name, image)
			ui.Column(c).Grow(1).Children(func() {
				ui.Text(c, name).Bold().SingleLine()
				kit.Muted(c, status).FontSize(t.Rem(0.8)).SingleLine().ContentTransition(ui.ContentTransition{})
			})
			voice = round(c, o.Voice, "phone")
			video = round(c, o.Video, "video")
		})
		ui.Box(c).FillWidth().Height(1).Background(t.Border)
	})
	if voice.Clicked() {
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("ChatHeader", "voice")
	}
	if video.Clicked() {
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("ChatHeader", "video")
	}
	return root, voice, video
}
