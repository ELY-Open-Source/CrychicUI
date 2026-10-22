package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// StrengthOptions tune a Strength.
type StrengthOptions struct {
	// Words name scores 1 to 4; Weak to Strong.
	Words [4]string
}

// Strength is password strength 0 to 4: bars, word.
func Strength(c *ui.Context, label string, score int, o StrengthOptions) ui.Element {
	if label == "" || score < 0 || score > 4 {
		panic("feedback: Strength needs a label and a score from 0 to 4")
	}
	if o.Words == ([4]string{}) {
		o.Words = [4]string{"Weak", "Fair", "Good", "Strong"}
	}
	t := c.Theme()
	tints := [4]ui.Color{t.Danger, t.Warning, t.Accent, t.Success}
	word, tint := "", t.Border
	if score > 0 {
		word, tint = o.Words[score-1], tints[score-1]
	}
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleMeter).Label(label).Value(word).Range(0, 4, float64(score))
	return root.Children(func() {
		ui.Row(c).Grow(1).Gap(t.Space(1)).Children(func() {
			for i := range 4 {
				c.Key(i)
				a := ui.Animation{Duration: 160 * time.Millisecond, Delay: time.Duration(i) * 40 * time.Millisecond}
				bar := ui.Box(c).Grow(1).Basis(0).Height(t.Space(1)).Radius(t.Space(0.5))
				col := t.Border
				if i < score {
					col = tint
				}
				bar.Background(bar.AnimateColor("tint", col, a))
			}
		})
		ui.Text(c, word).FontSize(t.Rem(0.85)).Width(t.Space(14)).TextColor(tint).ContentTransition(ui.ContentTransition{})
	})
}
