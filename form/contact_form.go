package form

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Contact is the value a ContactForm edits.
type Contact struct{ Name, Email, Message string }

// ContactFormOptions tune a ContactForm.
type ContactFormOptions struct {
	// Labels first, then the error messages.
	Name, Email, Message, Send          string
	NameError, EmailError, MessageError string
}

// ContactForm takes name, email and message; Send enables when valid.
func ContactForm(c *ui.Context, value *Contact, o ContactFormOptions) (root, send ui.Element) {
	if value == nil {
		panic("form: ContactForm needs a value")
	}
	def(&o.Name, "Name")
	def(&o.Email, "Email")
	def(&o.Message, "Message")
	def(&o.Send, "Send")
	def(&o.NameError, "Enter your name.")
	def(&o.EmailError, "Enter a valid email address.")
	def(&o.MessageError, "Write a message.")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(4))
	root.Children(func() {
		entry(c, o.Name, when(blank(value.Name), o.NameError), func() ui.Element {
			in := ui.TextInput(c, &value.Name).FillWidth()
			edited(root, in)
			return in
		})
		entry(c, o.Email, when(!emailPattern.MatchString(value.Email), o.EmailError), func() ui.Element {
			in := ui.TextInput(c, &value.Email).FillWidth()
			edited(root, in)
			return in
		})
		entry(c, o.Message, when(blank(value.Message), o.MessageError), func() ui.Element {
			in := ui.TextArea(c, &value.Message).FillWidth().Lines(4, 8)
			edited(root, in)
			return in
		})
		ok := !blank(value.Name) && emailPattern.MatchString(value.Email) && !blank(value.Message)
		send = primary(c, o.Send).Disabled(!ok)
	})
	if kit.Clicked(send, sound.Send) {
		kit.Log("ContactForm", "send", "message", len(value.Message))
	}
	return root, send
}
