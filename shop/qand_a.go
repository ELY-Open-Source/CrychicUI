package shop

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Question is one customer question and answer.
type Question struct {
	Text, Answer string
}

// QandAOptions are the question labels.
type QandAOptions struct{ Placeholder, Ask, Empty string }

// QandA shows questions and asks new ones; cue Send.
func QandA(c *ui.Context, questions *[]Question, draft *string, o QandAOptions) (root, ask ui.Element) {
	if questions == nil || draft == nil {
		panic("shop: QandA needs questions and draft")
	}
	def(&o.Placeholder, "Ask a question")
	def(&o.Ask, "Ask")
	def(&o.Empty, "No answers yet")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	pending := ui.Local(root, "pending", func() string { return "" })
	var in ui.Element
	root.Children(func() {
		for i, question := range *questions {
			if question.Text == "" {
				panic("shop: QandA question must not be empty")
			}
			c.Key(i)
			ui.Column(c).Gap(t.Space(1)).Children(func() {
				ui.Text(c, question.Text).Bold().Wrap()
				if question.Answer == "" {
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.8))
				} else {
					ui.Text(c, question.Answer).FontSize(t.Rem(0.85)).Wrap()
				}
			})
		}
		c.Key("ask-form")
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			in = ui.TextInput(c, draft).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1)
			in.OnTextKey(func(k ui.TextKey) bool {
				if k.Key != ui.KeyEnter || strings.TrimSpace(k.Text()) == "" {
					return false
				}
				*pending = strings.TrimSpace(k.Text())
				k.SetText("")
				return true
			})
			ask = addButton(c, o.Ask).Disabled(strings.TrimSpace(*draft) == "")
			ask.Children(func() { ui.Text(c, o.Ask).SingleLine() })
		})
	})
	if in.Changed() {
		root.SetChanged()
	}
	if *pending != "" || ask.Clicked() {
		text := *pending
		if text == "" {
			text = strings.TrimSpace(*draft)
		}
		if text != "" {
			*questions = append(*questions, Question{Text: text})
			*draft = ""
			ask.SetChanged()
			root.SetChanged()
			sound.Play(sound.Send)
			kit.Log("QandA", "asked", "runes", utf8.RuneCountInString(text))
		}
		*pending = ""
	}
	return root, ask
}
