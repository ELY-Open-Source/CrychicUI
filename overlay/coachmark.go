package overlay

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CoachmarkOptions tune a Coachmark.
type CoachmarkOptions struct {
	// GotIt names the dismiss button.
	GotIt string
}

// Coachmark is one tip at a pulsing anchor.
func Coachmark(c *ui.Context, anchor ui.Element, title, text string, open *bool, o CoachmarkOptions) ui.Element {
	if title == "" || open == nil {
		panic(fmt.Sprintf("overlay: Coachmark %q needs open", title))
	}
	def(&o.GotIt, "Got it")
	t := c.Theme()
	root := ui.Box(c)
	if *open {
		ui.Overlay(c, func() {
			dot := ui.Box(c).AttachTo(anchor, ui.AnchorTopRight, ui.AnchorCenter).Size(t.Space(3), t.Space(3)).Radius(999).
				Background(t.Accent).PassThrough()
			k := dot.Loop("pulse", 1400*time.Millisecond, ui.EaseOut)
			dot.Border(t.Space(0.5), t.Accent.Alpha(0.5*(1-k))).Scale(1+k, 1+k)
			card := ui.Column(c).AttachTo(anchor, ui.AnchorBottom, ui.AnchorTop).Margin(t.Space(3), 0, 0, 0).Width(t.Space(68)).
				Padding(t.Space(4)).Gap(t.Space(2)).Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText).Role(ui.RoleDialog).Label(title)
			kit.Arrive(card, ui.Pose{Y: -t.Space(2), Fade: 1}, 0).Children(func() {
				ui.Text(c, title).Bold()
				ui.Text(c, text).FontSize(t.Rem(0.9))
				ui.Row(c).FillWidth().Justify(ui.End).Children(func() {
					b := kit.Press(ui.ButtonBase(c).Label(o.GotIt).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius))
					b.Background(b.AnimateColor("bg", t.AccentText.Alpha(0.16+0.1*unit(b.Hovered())), ui.Snappy)).
						Children(func() { ui.Text(c, o.GotIt).Bold() })
					if b.Clicked() {
						*open = false
						root.SetChanged()
					}
				})
			})
			if card.OverlayShortcut(0, ui.KeyEscape) {
				*open = false
				root.SetChanged()
			}
		})
	}
	if _, closed := moved(root, *open, "Coachmark"); closed {
		sound.Play(sound.Close)
	}
	return root
}
