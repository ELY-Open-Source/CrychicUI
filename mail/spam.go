package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SpamOptions word the banner and its button.
type SpamOptions struct{ Warning, Detail, NotSpam string }

// Spam marks a message as not spam; cue Tap.
func Spam(c *ui.Context, spam *bool, o SpamOptions) ui.Element {
	if spam == nil {
		panic("mail: Spam needs state")
	}
	def(&o.Warning, "This message looks like spam")
	def(&o.Detail, "Be careful with its links and attachments.")
	def(&o.NotSpam, "Not spam")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).
		Background(t.Warning.Alpha(0.1)).Border(1, t.Warning.Alpha(0.5)).Role(ui.RoleGroup)
	shown := root.AnimateTo("spam", unit(*spam), ui.Snappy)
	root.Opacity(shown).Translate(0, -t.Space(2)*(1-shown))
	root.Children(func() {
		tile(c, "shield-alert", t.Warning)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, o.Warning).FontSize(t.Rem(0.85)).Bold().Wrap().FillWidth()
			kit.Muted(c, o.Detail).FontSize(t.Rem(0.75)).Wrap().FillWidth()
		})
		b := kit.Press(ui.ButtonBase(c).Label(o.NotSpam).Padding(t.Space(1.5), t.Space(3)).Radius(999).Shrink(0).Disabled(!*spam))
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
		b.Border(1, t.Border).Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "check", 14).TextColor(t.Success)
				ui.Text(c, o.NotSpam).FontSize(t.Rem(0.8)).Bold().Shrink(0)
			})
		})
		if b.Clicked() && *spam {
			*spam = false
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("Spam", "notSpam")
		}
	})
	return root
}
