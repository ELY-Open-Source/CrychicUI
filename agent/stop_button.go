package agent

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
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
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Row().AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).
		Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Disabled(!running))
	glow := float32(0.75)
	if running {
		ph := float64(root.Loop("pulse", 1400*time.Millisecond, ui.Linear))
		glow = 0.75 + 0.25*float32(math.Abs(math.Sin(math.Pi*ph)))
	}
	root.Background(root.AnimateColor("bg", t.Danger.Mix(t.Danger.Alpha(0.85), unit(root.Hovered())), ui.Snappy)).
		TextColor(t.AccentText).Opacity(root.AnimateTo("glow", glow, ui.Snappy))
	root.Children(func() {
		ui.Icon(c, icon.Must("square")).Size(16, 16)
		ui.Text(c, o.Label).Bold().SingleLine()
	})
	if running && root.Clicked() {
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("StopButton", "stopped")
	}
	return root
}
