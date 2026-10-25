package mail

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SignatureOptions word the card, the field and the preview.
type SignatureOptions struct{ Title, Placeholder, Preview, Empty string }

// Signature edits a signature with live preview.
func Signature(c *ui.Context, text *string, o SignatureOptions) ui.Element {
	if text == nil {
		panic("mail: Signature needs text")
	}
	def(&o.Title, "Email signature")
	def(&o.Placeholder, "Signature")
	def(&o.Preview, "Preview")
	def(&o.Empty, "No signature")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "pen-line", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Children(func() {
			field := ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).Border(1, t.Border)
			field.Children(func() {
				in := ui.TextAreaBase(c, text).Label(o.Placeholder).Placeholder(o.Placeholder).FontSize(t.Rem(0.85)).FillWidth().Lines(3, 6)
				field.BorderColor(field.AnimateColor("line", t.Border.Mix(t.Accent, unit(in.Focused())), ui.Snappy))
				if in.Changed() {
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("Signature", "edited", "runes", utf8.RuneCountInString(*text))
				}
			})
			kit.Muted(c, o.Preview).FontSize(t.Rem(0.7)).Bold()
			preview := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).
				Background(t.SurfaceHover.Alpha(0.6)).BorderWidth(0, 0, 0, t.Space(0.5)).BorderColor(t.Accent.Alpha(0.5))
			preview.Children(func() {
				if strings.TrimSpace(*text) == "" {
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.8))
				}
				for i, line := range strings.Split(strings.TrimSpace(*text), "\n") {
					if line == "" {
						continue
					}
					c.Key(i)
					l := ui.Text(c, line).FontSize(t.Rem(0.8)).Wrap().FillWidth()
					if i == 0 {
						l.Bold()
					} else {
						l.TextColor(t.TextMuted)
					}
				}
			})
		})
	})
	return root
}
