package shop

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Question is one customer question and its answer, if any.
type Question struct {
	Text, Answer string
}

// QandAOptions word the list; Q and A mark each pair.
type QandAOptions struct{ Placeholder, Ask, Empty, None, Q, A string }

// QandA lists answered questions and asks new ones; cue Send.
func QandA(c *ui.Context, questions *[]Question, draft *string, o QandAOptions) (root, ask ui.Element) {
	if questions == nil || draft == nil {
		panic("shop: QandA needs questions and a draft")
	}
	def(&o.Placeholder, "Ask a question")
	def(&o.Ask, "Ask")
	def(&o.Empty, "No answers yet")
	def(&o.None, "No questions yet. Be the first to ask.")
	def(&o.Q, "Q")
	def(&o.A, "A")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	pending := ui.Local(root, "pending", func() string { return "" })
	mark := func(letter string, ink ui.Color) {
		ui.Text(c, letter).FontSize(t.Rem(0.7)).Bold().TextColor(ink).TextAlign(ui.Center).Width(t.Space(5)).Shrink(0).
			Padding(t.Space(0.5), 0).Radius(t.Space(1)).Background(ink.Alpha(0.12))
	}
	var in ui.Element
	root.Children(func() {
		list := ui.Column(c).Key("list").FillWidth()
		list.Children(func() {
			if len(*questions) == 0 {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3), 0).TextColor(t.TextMuted).Children(func() {
					ui.Icon(c, icon.Must("message-square-text")).Size(16, 16)
					ui.Text(c, o.None).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
				})
			}
			seen := map[string]int{}
			for i, q := range *questions {
				if strings.TrimSpace(q.Text) == "" {
					panic("shop: QandA questions need text")
				}
				c.Key(fmt.Sprintf("%q#%d", q.Text, seen[q.Text]))
				seen[q.Text]++
				row := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2.5), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(*questions)-1)))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
						mark(o.Q, t.Accent)
						ui.Text(c, q.Text).Bold().FontSize(t.Rem(0.9)).Grow(1).MinWidth(0).Wrap()
					})
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
						mark(o.A, t.TextMuted)
						if q.Answer == "" {
							kit.Muted(c, o.Empty).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0)
						} else {
							ui.Text(c, q.Answer).FontSize(t.Rem(0.85)).LineHeight(1.5).Grow(1).MinWidth(0).Wrap()
						}
					})
				})
			}
		})
		c.Key("ask-form")
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
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
			ask.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must("send-horizontal")).Size(14, 14)
					ui.Text(c, o.Ask).Bold().Shrink(0)
				})
			})
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
