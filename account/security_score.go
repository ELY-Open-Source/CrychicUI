package account

import (
	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SecurityStep is one account security recommendation.
type SecurityStep struct {
	Title string
	Done  bool
}

// SecurityScoreOptions are the security score labels.
type SecurityScoreOptions struct{ Label string }

// SecurityScore shows a score ring with steps.
func SecurityScore(c *ui.Context, score int, steps []SecurityStep, o SecurityScoreOptions) ui.Element {
	if score < 0 || score > 100 || len(steps) == 0 {
		panic("account: SecurityScore needs score and steps")
	}
	def(&o.Label, "Security score")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		feedback.ProgressRing(c, o.Label, float64(score)/100, feedback.ProgressRingOptions{Format: "%.0f"})
		ui.Column(c).Grow(1).Gap(t.Space(1)).Children(func() {
			for i, step := range steps {
				if step.Title == "" {
					panic("account: security step must not be empty")
				}
				c.Key(i)
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					glyph, color := "circle", t.TextMuted
					if step.Done {
						glyph, color = "check-circle-2", t.Success
					}
					ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(color)
					ui.Text(c, step.Title).Wrap()
				})
			}
		})
	})
	return root
}
