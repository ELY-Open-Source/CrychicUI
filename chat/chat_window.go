package chat

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ChatWindowOptions tune a ChatWindow and its parts.
type ChatWindowOptions struct {
	// Height of the message list; zero gives the default.
	Height   float32
	Avatar   *ui.Bitmap
	Header   ChatHeaderOptions
	Composer ComposerOptions
}

// ChatWindow joins a header, a message list and a composer.
func ChatWindow(c *ui.Context, name, status string, list []Message, onSend func(text string), o ChatWindowOptions) (root, send, voice, video, attach ui.Element) {
	if name == "" || status == "" || onSend == nil {
		panic("chat: ChatWindow needs a name, a status and a send func")
	}
	messages("ChatWindow", list)
	size("ChatWindow", &o.Height, c.Theme().Space(90))
	t := c.Theme()
	root = ui.Column(c).FillWidth().Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup).Label(name)
	draft := ui.Local(root, "draft", func() string { return "" })
	root.Children(func() {
		var v, d, box ui.Element
		_, v, d = ChatHeader(c, name, status, o.Avatar, o.Header)
		voice, video = v, d
		MessageList(c, name, list, MessageListOptions{Height: o.Height}).Grow(1)
		box, send, attach = Composer(c, draft, onSend, o.Composer)
		if box.Changed() {
			root.SetChanged()
		}
	})
	if voice.Clicked() {
		kit.Log("ChatWindow", "voice")
	}
	if video.Clicked() {
		kit.Log("ChatWindow", "video")
	}
	if attach.Clicked() {
		kit.Log("ChatWindow", "attach")
	}
	return root, send, voice, video, attach
}
