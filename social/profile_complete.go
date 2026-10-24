package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ProfileStep is one profile completion step.
type ProfileStep struct {
	ID, Title string
	Done      bool
}

// ProfileCompleteOptions are the completion labels.
type ProfileCompleteOptions struct{ Complete string }

// ProfileComplete shows steps needed to finish a profile.
func ProfileComplete(c *ui.Context, steps []ProfileStep, o ProfileCompleteOptions) ui.Element {
	if len(steps) == 0 {
		panic("social: ProfileComplete needs steps")
	}
	def(&o.Complete, "%d/%d complete")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	done := 0
	for _, step := range steps {
		if step.ID == "" || step.Title == "" {
			panic("social: invalid profile step")
		}
		if step.Done {
			done++
		}
	}
	root.Children(func() {
		kit.Numeric(c, fmt.Sprintf(o.Complete, done, len(steps))).TextColor(t.TextMuted)
		bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
		w := bar.AnimateTo("complete", float32(done)/float32(len(steps)), ui.Smooth)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, t.Success, 999)
		}))
		for _, step := range steps {
			c.Key(step.ID)
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				glyph, color := "circle", t.TextMuted
				if step.Done {
					glyph, color = "check-circle-2", t.Success
				}
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(color)
				ui.Text(c, step.Title)
			})
		}
	})
	return root
}
