package marketing

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QA is one question and its answer.
type QA struct{ Question, Answer string }

// FAQOptions word the list.
type FAQOptions struct{ Title string }

// FAQ opens one answer at a time; cue Toggle.
func FAQ(c *ui.Context, items []QA, open *string, o FAQOptions) ui.Element {
	def(&o.Title, "Frequently asked questions")
	questions := make([]string, len(items))
	for i, q := range items {
		if q.Answer == "" {
			panic(fmt.Sprintf("marketing: FAQ question %q needs an answer", q.Question))
		}
		questions[i] = q.Question
	}
	unique("FAQ", questions)
	if len(items) == 0 || open == nil || *open != "" && !slices.Contains(questions, *open) {
		panic(fmt.Sprintf("marketing: FAQ %q needs questions and an open one among them", o.Title))
	}
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "circle-help", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(2))
		list.Children(func() {
			for i, q := range items {
				on := q.Question == *open
				c.Key(strconv.Quote(q.Question))
				item := ui.Column(c).FillWidth().Role(ui.RoleListItem)
				if i > 0 {
					item.BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
				}
				kit.Stagger(c, list, item, i).Children(func() {
					b := kit.Press(ui.ButtonBase(c).Label(q.Question).Checked(on).FillWidth().AlignItems(ui.Stretch).
						Padding(t.Space(3), t.Space(2)).Radius(t.Radius))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
					b.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, q.Question).FontSize(t.Rem(0.88)).Bold().TextColor(t.Text.Mix(t.Accent, unit(on))).Grow(1).MinWidth(0).Wrap()
							chev := glyph(c, "chevron-down", 16).TextColor(t.TextMuted)
							chev.Rotate(180 * chev.AnimateTo("turn", unit(on), ui.Snappy))
						})
					})
					if b.Clicked() {
						next := q.Question
						if on {
							next = ""
						}
						*open = next
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("FAQ", "toggled", "index", i, "open", !on)
					}
					layout.Collapse(c, on, func() {
						kit.Muted(c, q.Answer).FontSize(t.Rem(0.84)).FillWidth().Padding(0, t.Space(2), t.Space(3), t.Space(2)).Wrap()
					})
				})
			}
		})
	})
	return root
}
