package form

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConsentOptions tune a Consent.
type ConsentOptions struct {
	// Submit labels the button the check enables.
	Submit string
}

// Consent is a terms check that enables its submit button.
func Consent(c *ui.Context, terms string, agreed *bool, o ConsentOptions) (root, submit ui.Element) {
	if terms == "" || agreed == nil {
		panic("form: Consent needs terms and a value")
	}
	def(&o.Submit, "Continue")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(4))
	root.Children(func() {
		edited(root, choice.Check(c, terms, agreed))
		submit = primary(c, o.Submit).Disabled(!*agreed)
	})
	if kit.Clicked(submit, sound.Tap) {
		kit.Log("Consent", "submit")
	}
	return root, submit
}
