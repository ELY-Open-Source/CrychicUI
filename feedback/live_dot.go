package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// LiveDotOptions tune a LiveDot.
type LiveDotOptions struct {
	// Text follows the dot; "LIVE" when empty.
	Text string
}

// LiveDot is a pulsing red dot beside "LIVE".
func LiveDot(c *ui.Context, o LiveDotOptions) ui.Element {
	if o.Text == "" {
		o.Text = "LIVE"
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(1.5)).AlignSelf(ui.Start).Padding(t.Space(0.5), t.Space(2)).Radius(999).
		Background(t.Danger.Alpha(0.12)).Role(ui.RoleStatus).Label(o.Text)
	return root.Children(func() {
		d := t.Space(2)
		ui.Box(c).Size(d, d).Children(func() {
			halo := ui.Box(c).Absolute().Top(0).Left(0).Right(0).Bottom(0).Radius(999)
			k := halo.Loop("pulse", 1200*time.Millisecond, ui.EaseOut)
			halo.Background(t.Danger.Alpha(0.5*(1-k))).Scale(1+1.2*k, 1+1.2*k)
			ui.Box(c).Fill().Radius(999).Background(t.Danger)
		})
		ui.Text(c, o.Text).Bold().FontSize(t.Rem(0.75)).LetterSpacing(t.Rem(0.04)).TextColor(t.Danger)
	})
}
