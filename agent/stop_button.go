package agent

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StopButtonOptions tune a StopButton.
type StopButtonOptions struct {
	Label string
}

// StopButton is a square stop that pulses while running.
func StopButton(c *ui.Context, running bool, o StopButtonOptions) ui.Element {
	def(&o.Label, "Stop")
	t := c.Theme()
	side := t.Space(10)
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Tooltip(o.Label).Size(side, side).Shrink(0).Center().Radius(999).
		Disabled(!running))
	on := root.AnimateTo("on", unit(running), ui.Snappy)
	beat := float32(0)
	if running {
		beat = float32(math.Sin(2 * math.Pi * float64(root.Loop("pulse", 1200*time.Millisecond, ui.Linear))))
	}
	root.Background(root.AnimateColor("bg", t.SurfaceHover.Mix(t.Text.Mix(t.Danger, unit(root.Hovered())), on), ui.Snappy))
	root.Children(func() {
		sq := t.Space(3.5)
		s := 1 + 0.12*beat
		ui.Box(c).Size(sq, sq).Radius(t.Space(0.75)).Background(t.Background.Mix(t.TextMuted, 1-on)).Scale(s, s)
	})
	if running && root.Clicked() {
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("StopButton", "stopped")
	}
	return root
}
