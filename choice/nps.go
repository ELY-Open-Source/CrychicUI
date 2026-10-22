package choice

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NPSOptions tune an NPS.
type NPSOptions struct {
	// Low and High caption the ends; "Not likely", "Very likely".
	Low, High string
}

// NPS is a 0 to 10 scale; -1 is unanswered.
func NPS(c *ui.Context, label string, score *int, o NPSOptions) ui.Element {
	if label == "" || score == nil || *score < -1 || *score > 10 {
		panic("choice: NPS needs a label and a score from -1 to 10")
	}
	if o.Low == "" {
		o.Low = "Not likely"
	}
	if o.High == "" {
		o.High = "Very likely"
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(1.5)).AlignSelf(ui.Start)
	return root.Children(func() {
		changed := false
		ui.RadioGroup(c, func() {
			for i := range 11 {
				band := t.Success
				switch {
				case i <= 6:
					band = t.Danger
				case i <= 8:
					band = t.Warning
				}
				c.Key(i)
				n := strconv.Itoa(i)
				cell := ui.RadioBase(c, score, i).Label(n).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius)
				bg, edge, fg := band.Alpha(0.08), band.Alpha(0.35), t.Text
				if *score == i {
					bg, edge, fg = band, band, t.AccentText
				} else if cell.Hovered() {
					bg = band.Alpha(0.18)
				}
				cell.Background(cell.AnimateColor("bg", bg, quick)).Border(1, cell.AnimateColor("edge", edge, quick)).
					TextColor(cell.AnimateColor("fg", fg, quick))
				kit.Press(cell).Children(func() { ui.Text(c, n) })
				changed = kit.Changed(cell, sound.Select) || changed
			}
		}).Label(label).Row().Gap(t.Space(1))
		if changed {
			root.SetChanged()
			kit.Log("NPS", "scored", "score", *score)
		}
		ui.Row(c).Children(func() {
			kit.Muted(c, o.Low).FontSize(t.Rem(0.8))
			ui.Spacer(c)
			kit.Muted(c, o.High).FontSize(t.Rem(0.8))
		})
	})
}
