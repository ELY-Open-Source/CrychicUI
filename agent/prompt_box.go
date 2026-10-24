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
	Placeholder, Attach, Send, Hint string
}

// PromptBox is a large prompt input with model and attach.
func PromptBox(c *ui.Context, draft *string, model string, onSend func(text string), o PromptBoxOptions) (root, send, attach ui.Element) {
	if draft == nil || model == "" || onSend == nil {
		panic("agent: PromptBox needs a draft, a model and a send func")
	}
	def(&o.Placeholder, "Ask anything")
	def(&o.Attach, "Attach")
	def(&o.Send, "Send")
	def(&o.Hint, "Enter to send · Shift+Enter for a new line")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius*2).Background(t.Background).
		Shadow(0, 1, 3, 0, kit.Shade(0.06))
	root.Border(1, root.AnimateColor("ring", t.Border.Mix(t.Accent, unit(root.FocusWithin())), ui.Snappy))
	pending := ui.Local(root, "pending", func() string { return "" })
	blank := strings.TrimSpace(*draft) == ""
	root.Children(func() {
		in := ui.TextAreaBase(c, draft).FillWidth().Lines(2, 8).Placeholder(o.Placeholder).Label(o.Placeholder).OnTextKey(func(k ui.TextKey) bool {
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
			attach = round(c, o.Attach, "paperclip")
			ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Radius(999).Background(t.Surface).
				Children(func() {
					ui.Icon(c, icon.Must("sparkles")).Size(16, 16).TextColor(t.Accent)
					ui.Text(c, model).FontSize(t.Rem(0.8)).Shrink(0)
				})
			kit.Muted(c, o.Hint).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
			send = kit.Press(ui.ButtonBase(c).Label(o.Send).Tooltip(o.Send).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999)).
				Disabled(blank)
			bg := t.SurfaceHover
			if !blank {
				bg = t.Accent.Mix(t.AccentHover, unit(send.Hovered()))
			}
			send.Background(send.AnimateColor("bg", bg, ui.Snappy)).TextColor(send.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(!blank)), ui.Snappy)).
				Children(func() { ui.Icon(c, icon.Must("arrow-up")).Size(16, 16) })
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
