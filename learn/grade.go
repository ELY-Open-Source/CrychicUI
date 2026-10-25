package learn

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GradeOptions word the badge; Letters run A to F.
type GradeOptions struct {
	Label, Score string
	Letters      [5]string
}

// Grade shows a letter grade for a 0–100 score.
func Grade(c *ui.Context, score float64, o GradeOptions) ui.Element {
	if !(score >= 0 && score <= 100) {
		panic(fmt.Sprintf("learn: Grade %q needs a score of 0 to 100, got %v", o.Label, score))
	}
	def(&o.Score, "%s%%")
	for i, l := range [...]string{"A", "B", "C", "D", "F"} {
		def(&o.Letters[i], l)
	}
	t := c.Theme()
	band := 4
	for i, floor := range []float64{90, 80, 70, 60} {
		if score >= floor {
			band = i
			break
		}
	}
	ink := [...]ui.Color{t.Success, t.Accent, t.Warning, t.Warning.Mix(t.Danger, 0.5), t.Danger}[band]
	shown := math.Floor(score*10) / 10
	text := fmt.Sprintf(o.Score, fmt.Sprint(shown))
	root := ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(o.Letters[band])
	root.Children(func() {
		c.Key(band)
		badge := ui.Box(c).Size(t.Space(16), t.Space(16)).Shrink(0).Center().Radius(t.Radius*1.5).
			Background(ink.Alpha(0.14)).Border(t.Space(0.5), ink.Alpha(0.5)).TextColor(ink)
		kit.Arrive(badge, ui.Pose{Scale: 0.5, Rotate: -12, Fade: 1}, 0).Children(func() {
			ui.Text(c, o.Letters[band]).FontSize(t.Rem(1.9)).Bold().MaxWidth(t.Space(14)).SingleLine()
		})
		ui.Column(c).Gap(t.Space(0.5)).MinWidth(0).Children(func() {
			kit.Numeric(c, text).FontSize(t.Rem(1.1)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			if o.Label != "" {
				kit.Muted(c, o.Label).FontSize(t.Rem(0.74)).MaxWidth(t.Space(50)).SingleLine()
			}
		})
	})
	return root
}
