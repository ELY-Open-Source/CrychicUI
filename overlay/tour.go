package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TourStep is one stop of a Tour.
type TourStep struct{ Title, Text string }

// TourOptions tune a Tour.
type TourOptions struct {
	Back, Next, Done string
}

// Tour spotlights anchors in turn, a card under each.
func Tour(c *ui.Context, label string, anchors []ui.Element, steps []TourStep, step *int, open *bool, o TourOptions) ui.Element {
	if label == "" || len(steps) == 0 || len(anchors) != len(steps) || step == nil || open == nil || !valid(*step, len(steps)) {
		panic(fmt.Sprintf("overlay: Tour %q needs an anchor per step, a step and open", label))
	}
	for _, s := range steps {
		if s.Title == "" {
			panic(fmt.Sprintf("overlay: Tour %q steps need titles", label))
		}
	}
	def(&o.Back, "Back")
	def(&o.Next, "Next")
	def(&o.Done, "Done")
	t := c.Theme()
	root := ui.Box(c)
	shown, last := *step, len(steps)-1
	move := func(to int) {
		*step = to
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Tour", "step", "index", to)
	}
	shut := func() {
		*open = false
		root.SetChanged()
	}
	if *open {
		target := anchors[shown]
		r, pad := target.Bounds(), t.Space(2)
		ui.Overlay(c, func() {
			scrim := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Modal().Role(ui.RoleDialog).Label(label)
			x, y := scrim.AnimateTo("x", r.X-pad, ui.Snappy), scrim.AnimateTo("y", r.Y-pad, ui.Snappy)
			w, h := scrim.AnimateTo("w", r.W+2*pad, ui.Snappy), scrim.AnimateTo("h", r.H+2*pad, ui.Snappy)
			dim := kit.Shade(0.5)
			scrim.Draw(func(p *ui.Painter, b ui.Rect) {
				p.Fill(ui.Rect{X: b.X, Y: b.Y, W: b.W, H: y - b.Y}, dim, 0)
				p.Fill(ui.Rect{X: b.X, Y: y + h, W: b.W, H: b.Y + b.H - y - h}, dim, 0)
				p.Fill(ui.Rect{X: b.X, Y: y, W: x - b.X, H: h}, dim, 0)
				p.Fill(ui.Rect{X: x + w, Y: y, W: b.X + b.W - x - w, H: h}, dim, 0)
			})
			kit.Settle(c, scrim, []float32{r.X, r.Y, r.W, r.H}, func() []float32 {
				b := target.Bounds()
				return []float32{b.X, b.Y, b.W, b.H}
			})
			kit.Arrive(scrim, ui.Pose{Fade: 1}, 0).Children(func() {
				// The card hangs off this, inside the modal's focus.
				hole := ui.Box(c).Absolute().Left(r.X).Top(r.Y).Width(r.W).Height(r.H)
				c.Key(shown)
				card := ui.Column(c).AttachTo(hole, ui.AnchorBottom, ui.AnchorTop).Margin(t.Space(3), 0, 0, 0).Width(t.Space(72)).
					Padding(t.Space(4)).Gap(t.Space(2)).Radius(t.Radius).Background(t.Background).Shadow(0, 8, 24, 0, kit.Shade(0.2))
				kit.Arrive(card, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() {
					kit.Muted(c, fmt.Sprintf("%d / %d", shown+1, len(steps))).FontSize(t.Rem(0.8))
					ui.Text(c, steps[shown].Title).Bold()
					kit.Muted(c, steps[shown].Text)
					ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
						if quiet(c.Key("back"), o.Back).Disabled(shown == 0).Clicked() {
							move(shown - 1)
						}
						if shown < last {
							if primary(c.Key("next"), o.Next, t.Accent, t.AccentText).Clicked() {
								move(shown + 1)
							}
							return
						}
						if primary(c.Key("done"), o.Done, t.Accent, t.AccentText).Clicked() {
							sound.Play(sound.Tick)
							shut()
						}
					})
				})
			})
			if scrim.OverlayShortcut(0, ui.KeyEscape) {
				shut()
			}
		})
	}
	moved(root, *open, "Tour")
	return root
}
