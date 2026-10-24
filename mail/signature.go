package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SignatureOptions are the signature labels.
type SignatureOptions struct{ Placeholder, Preview string }

// Signature edits a signature with live preview.
func Signature(c *ui.Context, text *string, o SignatureOptions) ui.Element {
	if text == nil {
		panic("mail: Signature needs text")
	}
	def(&o.Placeholder, "Signature")
	def(&o.Preview, "Preview")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		in := ui.TextArea(c, text).Label(o.Placeholder).Placeholder(o.Placeholder).FillWidth().Lines(3, 6)
		if in.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("Signature", "edited", "runes", len([]rune(*text)))
		}
		kit.Muted(c, o.Preview).FontSize(t.Rem(0.75))
		preview := ui.Text(c, *text).Wrap().FillWidth().Padding(t.Space(2)).Radius(t.Radius).Background(t.Surface)
		preview.Transition(kit.Enter)
	})
	return root
}
