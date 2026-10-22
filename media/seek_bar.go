package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SeekBarOptions tune a SeekBar.
type SeekBarOptions struct {
	// Step is a key's seek; zero gives 5 s.
	Step time.Duration
}

// SeekBar seeks pos, showing buffered and played parts.
func SeekBar(c *ui.Context, label string, pos *time.Duration, buffered, total time.Duration, o SeekBarOptions) ui.Element {
	if label == "" || pos == nil || total <= 0 || *pos < 0 || *pos > total || buffered < 0 || buffered > total || o.Step < 0 {
		panic(fmt.Sprintf("media: SeekBar %q needs a total above 0 with pos and buffered within it", label))
	}
	if o.Step == 0 {
		o.Step = 5 * time.Second
	}
	t := c.Theme()
	v := pos.Seconds()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5))
	return root.Children(func() {
		kw, h, side := t.Space(3.5), t.Space(1), t.Space(5)
		s := ui.SliderBase(c, &v, 0, total.Seconds()).Step(o.Step.Seconds()).Label(label).FillWidth().Height(side).PaddingX(kw / 2)
		if s.Changed() {
			*pos = time.Duration(v * float64(time.Second))
			root.SetChanged()
		}
		_, _, held := s.Dragged()
		played := float32(v / total.Seconds())
		s.Children(func() {
			ui.Box(c).Fill().Children(func() {
				y := (side - h) / 2
				ui.Box(c).Absolute().Start(0).End(0).Top(y).Height(h).Radius(999).Background(t.Border)
				b := ui.Box(c).Absolute().Start(0).Top(y).Height(h).Radius(999).Background(t.TextMuted.Alpha(0.35))
				b.WidthPercent(100 * b.AnimateTo("buffered", float32(buffered.Seconds()/total.Seconds()), ui.Snappy))
				ui.Box(c).Absolute().Start(0).Top(y).Height(h).Radius(999).Background(t.Accent).WidthPercent(100 * played)
				dx := -kw / 2
				if c.Direction() == ui.RightToLeft {
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
		released(s, "SeekBar", v)
		ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
			kit.Numeric(c, clock(*pos)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			kit.Numeric(c, clock(total)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		})
	})
}
