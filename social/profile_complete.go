package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ProfileStep is one profile completion step.
type ProfileStep struct {
	ID, Title string
	Done      bool
}

// ProfileCompleteOptions word the card; Complete takes done and total.
type ProfileCompleteOptions struct{ Title, Complete, Done string }

// ProfileComplete shows steps needed to finish a profile.
func ProfileComplete(c *ui.Context, steps []ProfileStep, o ProfileCompleteOptions) ui.Element {
	if len(steps) == 0 {
		panic("social: ProfileComplete needs steps")
	}
	ids := make([]string, len(steps))
	done := 0
	for i, step := range steps {
		if step.Title == "" {
			panic(fmt.Sprintf("social: ProfileComplete step %q needs a title", step.ID))
		}
		ids[i] = step.ID
		if step.Done {
			done++
		}
	}
	unique("ProfileComplete", ids)
	def(&o.Title, "Complete your profile")
	def(&o.Complete, "%d/%d complete")
	def(&o.Done, "All set")
	t := c.Theme()
	share := float32(done) / float32(len(steps))
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3.5)).Children(func() {
			meter := ui.Box(c).Size(t.Space(14), t.Space(14)).Shrink(0).Center()
			fill := meter.AnimateTo("fill", share, ui.Smooth)
			meter.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(0.75)
				p.StrokePath(ring(cx, cy, rad), t.Space(1.25), t.SurfaceHover)
				if fill > 0 {
					p.StrokePath(ring(cx, cy, rad).Trim(0, fill), t.Space(1.25), t.Success)
				}
			}))
			meter.Children(func() {
				kit.Numeric(c, fmt.Sprintf("%d%%", done*100/len(steps))).FontSize(t.Rem(0.82)).Bold().Shrink(0)
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, o.Title).FontSize(t.Rem(0.95)).Bold().Wrap().FillWidth()
				kit.Numeric(c, fmt.Sprintf(o.Complete, done, len(steps))).FontSize(t.Rem(0.78)).TextColor(t.TextMuted)
			})
			c.Key(done == len(steps))
			if done == len(steps) {
				kit.Arrive(kit.Chip(c, o.Done, t.Success), ui.Pose{Scale: 0.6, Fade: 1}, 0)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(0, t.Space(1.5), t.Space(1.5), t.Space(1.5))
		list.Children(func() {
			for i, step := range steps {
				c.Key(step.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered() && !step.Done)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					mark, ink := "circle", t.TextMuted
					if step.Done {
						mark, ink = "circle-check-big", t.Success
					}
					glyph(c, mark, 18).TextColor(ink)
					title := ui.Text(c, step.Title).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					if step.Done {
						title.TextColor(t.TextMuted).Strikethrough()
					} else {
						glyph(c, "chevron-right", 16).TextColor(t.TextMuted)
					}
				})
			}
		})
	})
	return root
}
