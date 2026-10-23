package agent

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PromptBoxOptions tune a PromptBox.
type PromptBoxOptions struct {
	Placeholder, Attach, Send string
}

// PromptBox is a large prompt input with model and attach.
func PromptBox(c *ui.Context, draft *string, model string, onSend func(text string), o PromptBoxOptions) (root, send, attach ui.Element) {
	if draft == nil || model == "" || onSend == nil {
		panic("agent: PromptBox needs a draft, a model and a send func")
	}
	def(&o.Placeholder, "Ask anything")
	def(&o.Attach, "Attach")
	def(&o.Send, "Send")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius*2).
		Background(t.Background).Border(1, t.Border)
	pending := ui.Local(root, "pending", func() string { return "" })
	blank := strings.TrimSpace(*draft) == ""
	root.Children(func() {
		in := ui.TextArea(c, draft).FillWidth().Lines(3, 8).Placeholder(o.Placeholder).Label(o.Placeholder).OnTextKey(func(k ui.TextKey) bool {
			if k.Key != ui.KeyEnter || k.Mods != 0 {
				return false
			}
			if text := strings.TrimSpace(k.Text()); text != "" {
				*pending = text
				k.SetText("")
			}
			return true
		})
		if in.Changed() {
			root.SetChanged()
		}
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			tag(c, model)
			ui.Box(c).Grow(1)
			attach = round(c, o.Attach, "paperclip")
			send = kit.Press(ui.ButtonBase(c).Label(o.Send).Tooltip(o.Send).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999)).
				Disabled(blank)
			send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy)).
				TextColor(t.AccentText).Children(func() { ui.Icon(c, icon.Must("send")).Size(16, 16) })
		})
	})
	if kit.Clicked(attach, sound.Tap) {
		kit.Log("PromptBox", "attach")
	}
	text := *pending
	*pending = ""
	if text == "" && !blank && send.Clicked() {
		text = strings.TrimSpace(*draft)
		*draft = ""
	}
	if text != "" {
		onSend(text)
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("PromptBox", "sent", "runes", utf8.RuneCountInString(text))
	}
	return root, send, attach
}
