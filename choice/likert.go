package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LikertOptions tune a Likert.
type LikertOptions struct {
	// Labels run disagree to agree; five English steps.
	Labels []string
}

// Likert is a disagree to agree scale; -1 is unanswered.
func Likert(c *ui.Context, statement string, value *int, o LikertOptions) ui.Element {
	if len(o.Labels) == 0 {
		o.Labels = []string{"Strongly disagree", "Disagree", "Neutral", "Agree", "Strongly agree"}
	}
	distinct("Likert", o.Labels)
	if statement == "" || value == nil || *value < -1 || *value >= len(o.Labels) {
		panic("choice: Likert needs a statement and a value from -1 to len(Labels)-1")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		ui.Text(c, statement).Bold()
		changed := false
		ui.RadioGroup(c, func() {
			for i, name := range o.Labels {
				c.Key(i)
				b := ui.RadioBase(c, value, i).Label(name).Grow(1).Basis(0).Column().Gap(t.Space(1.5)).
					Padding(t.Space(2), t.Space(1)).Radius(t.Radius).AlignItems(ui.Center)
				on := *value == i
				bg := ui.Color{}
				if on {
					bg = t.Accent.Alpha(0.08)
				} else if b.Hovered() {
					bg = t.SurfaceHover
				}
				b.Background(b.AnimateColor("bg", bg, quick)).Children(func() {
					dot(c, b, on)
					ui.Text(c, name).FontSize(t.Rem(0.85)).TextAlign(ui.Center)
				})
				changed = kit.Changed(b, sound.Select) || changed
			}
		}).Label(statement).Row().FillWidth().Gap(t.Space(1)).AlignItems(ui.Start)
		if changed {
			root.SetChanged()
			kit.Log("Likert", "answered", "value", *value)
		}
	})
}
