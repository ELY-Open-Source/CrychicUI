package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ScrubberOptions tune a Scrubber.
type ScrubberOptions struct {
	// Step is a key's move; zero gives 1 s.
	Step time.Duration
}

// Scrubber scrubs pos; a preview frame rides over the track.
func Scrubber(c *ui.Context, label string, frames []*ui.Bitmap, pos *time.Duration, total time.Duration, o ScrubberOptions) ui.Element {
	pictures("Scrubber", label, frames)
	if pos == nil || total <= 0 || *pos < 0 || *pos > total || o.Step < 0 {
		panic(fmt.Sprintf("media: Scrubber %q needs a total above 0, a pos within it and a step not negative", label))
	}
	if o.Step == 0 {
		o.Step = time.Second
	}
	t := c.Theme()
	v := pos.Seconds()
	pw, ph := t.Space(28), t.Space(16)
	kw, h, side := t.Space(3.5), t.Space(1), t.Space(5)
	rtl := c.Direction() == ui.RightToLeft
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		lane := ui.Box(c).FillWidth().Height(ph + t.Space(6))
		s := ui.SliderBase(c, &v, 0, total.Seconds()).Step(o.Step.Seconds()).Label(label).FillWidth().Height(side).PaddingX(kw / 2)
		if s.Changed() {
			*pos = time.Duration(v * float64(time.Second))
			root.SetChanged()
		}
		b := s.LayoutBounds()
		kit.Settle(c, s, []float32{b.W}, func() []float32 { return []float32{s.LayoutBounds().W} })
		played := float32(v / total.Seconds())
		_, _, held := s.Dragged()
		x, _, over := s.PointerPosition()
		f, shown := played, held || s.Focused()
		if over && s.Hovered() && !held && b.W > kw {
			f, shown = max(0, min(1, (x-kw/2)/(b.W-kw))), true
			if rtl {
				f = 1 - f
			}
		}
		s.Children(func() {
			ui.Box(c).Fill().Children(func() {
				y := (side - h) / 2
				ui.Box(c).Absolute().Start(0).End(0).Top(y).Height(h).Radius(999).Background(t.Border)
				ui.Box(c).Absolute().Start(0).Top(y).Height(h).Radius(999).Background(t.Accent).WidthPercent(100 * played)
				dx := -kw / 2
				if rtl {
					dx = kw / 2
				}
				pose := ui.Pose{}
				if held || s.Hovered() || s.Focused() {
					pose.Scale = 1.3
				}
				ui.Box(c).Absolute().Size(kw, kw).Radius(999).Background(t.Accent).StartPercent(100*played).
					Top((side-kw)/2).Translate(dx, 0).AnimatePose(pose, ui.Snappy)
			})
		})
		released(s, "Scrubber", v)
		mid := kw/2 + f*(b.W-kw)
		if rtl {
			mid = b.W - mid
		}
		left := max(0, min(b.W-pw, mid-pw/2))
		frame := frames[min(len(frames)-1, int(f*float32(len(frames))))]
		lane.Children(func() {
			card := ui.Column(c).Absolute().Left(left).Top(0).Width(pw).Gap(t.Space(1)).AlignItems(ui.Center).
				TransformOrigin(0.5, 1).Opacity(lane.AnimateTo("seen", unit(shown), ui.Snappy))
			card.AnimatePose(ui.Pose{Scale: 1 - 0.06*(1-unit(shown))}, ui.Snappy).Children(func() {
				ui.Box(c).Size(pw, ph).Radius(t.Radius).Clip().Border(1, t.Border).Background(t.Surface).Children(func() {
					ui.Image(c, frame).Fill().Fit(ui.Cover)
				})
				kit.Numeric(c, clock(time.Duration(float64(f)*float64(total)))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			})
		})
		ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
			kit.Numeric(c, clock(*pos)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			kit.Numeric(c, clock(total)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		})
	})
}
