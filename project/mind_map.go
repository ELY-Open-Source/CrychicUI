package project

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MindMapOptions word the add button.
type MindMapOptions struct{ Add string }

// MindMap shows nodes branching from a center; cue Pop.
func MindMap(c *ui.Context, center string, branches []string, added *int, o MindMapOptions) ui.Element {
	if center == "" || added == nil || *added < -1 || *added > len(branches) || len(branches) > 8 {
		panic(fmt.Sprintf("project: MindMap %q needs state and at most 8 branches", center))
	}
	unique("MindMap", branches)
	def(&o.Add, "Add branch")
	t := c.Theme()
	nodeH, pad, gap, centerW, centerH := t.Space(9), t.Space(3), t.Space(2), t.Space(26), t.Space(11)
	root := card(c)
	root.Children(func() {
		canvas := ui.Box(c).FillWidth().Height(t.Space(80))
		w, h := canvas.LayoutBounds().W, t.Space(80)
		kit.Settle(c, canvas, []float32{w}, func() []float32 { return []float32{canvas.LayoutBounds().W} })
		cx, cy := w/2, h/2
		// Side nodes fit between the center and the edge.
		nodeW := min(t.Space(24), max(0, cx-centerW/2-gap-pad))
		rx, ry := cx-nodeW/2-pad, cy-nodeH/2-pad
		at := func(k int) (float32, float32) {
			a := -math.Pi/2 + 2*math.Pi*float64(k)/float64(len(branches))
			return cx + rx*float32(math.Cos(a)), cy + ry*float32(math.Sin(a))
		}
		grow := make([]float32, len(branches))
		for k, branch := range branches {
			grow[k] = canvas.Keyframes("grow."+branch, []ui.Keyframe{{Value: 0}, {At: 400 * time.Millisecond, Value: 1, Ease: ui.EaseOut}},
				ui.KeyframeOptions{Delay: time.Duration(k) * 50 * time.Millisecond})
		}
		canvas.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			for k, branch := range branches {
				x, y := at(k)
				var path ui.Path
				mx := r.X + cx + (x-cx)/2
				path.MoveTo(r.X+cx, r.Y+cy).CubeTo(mx, r.Y+cy, mx, r.Y+y, r.X+x, r.Y+y)
				p.StrokePath(path.Trim(0, grow[k]), 1.5, labelInk(t, branch).Alpha(0.55))
			}
		}))
		canvas.Children(func() {
			ui.Box(c).Absolute().Left(cx-centerW/2).Top(cy-centerH/2).Size(centerW, centerH).Label(center).Center().Padding(0, t.Space(3)).
				Radius(999).Background(t.Accent).Shadow(0, 2, 6, 0, kit.Shade(0.12)).Children(func() {
				ui.Text(c, center).FontSize(t.Rem(0.85)).Bold().TextColor(t.AccentText).MinWidth(0).SingleLine()
			})
			for k, branch := range branches {
				c.Key(branch)
				x, y := at(k)
				ink := labelInk(t, branch)
				node := ui.Box(c).Absolute().Left(x-nodeW/2).Top(y-nodeH/2).Size(nodeW, nodeH).Label(branch).Center().Padding(0, t.Space(2)).
					Radius(999).Background(t.Background.Mix(ink, 0.14)).Border(1, ink.Alpha(0.45))
				kit.Arrive(node, ui.Pose{Scale: 0.5, Fade: 1}, time.Duration(k)*50*time.Millisecond).Children(func() {
					ui.Text(c, branch).FontSize(t.Rem(0.75)).Bold().TextColor(ink).MinWidth(0).SingleLine()
				})
			}
		})
		full := len(branches) == 8
		ui.Row(c).FillWidth().Padding(t.Space(2), t.Space(2.5)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			b := kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(1.5), t.Space(2.5)).Radius(999).Disabled(full))
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
			b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(b.Hovered())), ui.Snappy))
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "plus", 14)
					ui.Text(c, o.Add).FontSize(t.Rem(0.78)).Shrink(0)
				})
			})
			if b.Clicked() && !full {
				*added = len(branches)
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("MindMap", "added", "count", len(branches)+1)
			}
		})
	})
	return root
}
