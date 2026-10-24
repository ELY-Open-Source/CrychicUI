package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SpamOptions are the spam banner labels.
type SpamOptions struct{ Warning, NotSpam string }

// Spam marks a message as not spam; cue Tap.
func Spam(c *ui.Context, spam *bool, o SpamOptions) ui.Element {
	if spam == nil {
		panic("mail: Spam needs state")
	}
	def(&o.Warning, "This message looks like spam")
	def(&o.NotSpam, "Not spam")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
		Background(t.Warning.Alpha(0.15)).Border(1, t.Warning)
	root.Opacity(root.AnimateTo("spam", float32(unit(*spam)), ui.Snappy))
	root.Children(func() {
		ui.Text(c, o.Warning).Grow(1).Wrap()
		b := button(c, o.NotSpam).Disabled(!*spam)
		b.Children(func() { ui.Text(c, o.NotSpam) })
		if b.Clicked() && *spam {
			*spam = false
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("Spam", "notSpam")
		}
	})
	return root
}
