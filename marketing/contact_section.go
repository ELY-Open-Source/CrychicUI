package marketing

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ContactInfo is one way to reach the team.
type ContactInfo struct{ Icon, Label, Value string }

// ContactForm is a visitor's message.
type ContactForm struct{ Name, Email, Text string }

// ContactSectionOptions word the section.
type ContactSectionOptions struct{ Title, Text, Name, Email, Message, Send, Sent string }

// ContactSection is contact info beside a form; cue Success.
func ContactSection(c *ui.Context, info []ContactInfo, form *ContactForm, sent *bool, o ContactSectionOptions) (root, send ui.Element) {
	def(&o.Title, "Get in touch")
	def(&o.Name, "Your name")
	def(&o.Email, "Email address")
	def(&o.Message, "How can we help?")
	def(&o.Send, "Send message")
	def(&o.Sent, "Thanks, we will reply within a day")
	if form == nil || sent == nil {
		panic(fmt.Sprintf("marketing: ContactSection %q needs a form and a sent state", o.Title))
	}
	labels := make([]string, len(info))
	for i, in := range info {
		if strings.TrimSpace(in.Icon) == "" || strings.TrimSpace(in.Value) == "" {
			panic(fmt.Sprintf("marketing: ContactSection %q item %q needs an icon and a value", o.Title, in.Label))
		}
		labels[i] = in.Label
	}
	unique("ContactSection", labels)
	t := c.Theme()
	root = card(c).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(6)).Children(func() {
			ui.Column(c).Grow(1).Basis(t.Space(44)).MinWidth(0).Gap(t.Space(3)).Children(func() {
				ui.Text(c, o.Title).FontSize(t.Rem(1.2)).Bold().FillWidth().Wrap()
				if o.Text != "" {
					kit.Muted(c, o.Text).FontSize(t.Rem(0.82)).FillWidth().Wrap()
				}
				list := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(1), 0, 0, 0)
				list.Children(func() {
					for i, in := range info {
						r := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center)
						kit.Stagger(c, list, r, i).Children(func() {
							tile(c, in.Icon, t.Accent, t.Space(9))
							ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
								kit.Muted(c, in.Label).FontSize(t.Rem(0.7)).FillWidth().MinWidth(0).SingleLine()
								ui.Text(c, in.Value).FontSize(t.Rem(0.84)).Bold().FillWidth().Wrap()
							})
						})
					}
				})
			})
			col := ui.Column(c).Grow(1).Basis(t.Space(56)).MinWidth(0).Gap(t.Space(2.5))
			w := col.LayoutBounds().W
			kit.Settle(c, col, []float32{w}, func() []float32 { return []float32{col.LayoutBounds().W} })
			room := t.Space(60)
			if w > 0 {
				room = max(t.Space(10), min(room, w-2*t.Space(4)-16-t.Space(1.5)))
			}
			col.Children(func() {
				name := fieldIn(c, ui.Row(c).FillWidth().MinWidth(0), o.Name, "user", &form.Name)
				email := fieldIn(c, ui.Row(c).FillWidth().MinWidth(0), o.Email, "mail", &form.Email)
				frame := ui.Column(c).FillWidth().Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).Background(t.Background)
				frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
				var text ui.Element
				frame.Children(func() {
					text = ui.TextAreaBase(c, &form.Text).Label(o.Message).Placeholder(o.Message).FontSize(t.Rem(0.86)).FillWidth().Lines(3, 6)
				})
				if name.Changed() || email.Changed() || text.Changed() {
					*sent = false
					root.SetChanged()
					kit.Log("ContactSection", "edited", "runes", utf8.RuneCountInString(form.Name+form.Email+form.Text))
				}
				ok := strings.TrimSpace(form.Name) != "" && validEmail(form.Email) && strings.TrimSpace(form.Text) != "" && !*sent
				c.Key("send")
				ui.Row(c).FillWidth().Justify(ui.End).Children(func() {
					send = action(c, o.Send, "send", true, room).Disabled(!ok)
				})
				if (send.Clicked() || name.Submitted() || email.Submitted()) && ok {
					*sent = true
					root.SetChanged()
					sound.Play(sound.Success)
					kit.Log("ContactSection", "sent")
				}
				ui.Column(c).FillWidth().Children(func() {
					if *sent {
						c.Key("sent")
						note := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).TextColor(t.Success)
						kit.Arrive(note, ui.Pose{Y: t.Space(1.5), Fade: 1}, 0).Children(func() {
							glyph(c, "circle-check", 16)
							ui.Text(c, o.Sent).FontSize(t.Rem(0.82)).Bold().Grow(1).MinWidth(0).Wrap()
						})
					}
				})
			})
		})
	})
	return root, send
}
