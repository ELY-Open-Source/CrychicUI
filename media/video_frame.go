package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VideoFrameOptions tune a VideoFrame.
type VideoFrameOptions struct {
	Play, Seek string
}

// VideoFrame is a poster with play and controls.
func VideoFrame(c *ui.Context, label string, poster *ui.Bitmap, pos *time.Duration, total time.Duration, playing *bool, o VideoFrameOptions) ui.Element {
	pictures("VideoFrame", label, []*ui.Bitmap{poster})
	if pos == nil || playing == nil || total <= 0 || *pos < 0 || *pos > total {
		panic(fmt.Sprintf("media: VideoFrame %q needs a total above 0, a pos within it and state", label))
	}
	def(&o.Play, "Play")
	def(&o.Seek, "Position")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Label(label)
	now := c.Now()
	advance(root, now, pos, total, playing, "VideoFrame")
	ink := t.Background
	if t.Dark {
		ink = t.Text
	}
	return root.Children(func() {
		ui.Grid(c).Columns(1).FillWidth().AspectRatio(16.0 / 9).Radius(t.Radius).Clip().Background(kit.Shade(0.9)).
			Children(func() {
				shot := ui.Image(c, poster).ColumnStart(1).RowStart(1).Fill().Fit(ui.Cover)
				shot.AnimatePose(ui.Pose{Scale: 1 + 0.03*unit(*playing)}, ui.Smooth)
				if *playing {
					return
				}
				c.Key("overlay")
				wrap := ui.Box(c).ColumnStart(1).RowStart(1).Fill().Center().
					Transition(ui.ElementTransition{Duration: 200 * time.Millisecond, Exit: &ui.Motion{Scale: 1.2}})
				kit.Arrive(wrap, ui.Pose{Fade: 1, Scale: 0.8}, 0).Children(func() {
					b := round(c, o.Play, "play", t.Space(16))
					b.Background(b.AnimateColor("plate", kit.Shade(0.45+0.15*unit(b.Hovered())), ui.Snappy)).TextColor(ink)
					if b.Clicked() {
						*playing = true
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("VideoFrame", "played")
						advance(root, now, pos, total, playing, "VideoFrame")
					}
				})
			})
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			if button.PlayButton(c, playing, button.PlayButtonOptions{Label: o.Play}).Changed() {
				root.SetChanged()
				advance(root, now, pos, total, playing, "VideoFrame")
			}
			if SeekBar(c, o.Seek, pos, total, total, SeekBarOptions{}).Grow(1).Changed() {
				root.SetChanged()
			}
		})
		if *playing {
			c.AnimationFrame()
		}
	})
}
