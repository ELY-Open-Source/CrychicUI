package code

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DebugBarOptions name the five actions and two states.
type DebugBarOptions struct {
	Continue, Over, Into, Out, Stop string
	Paused, Running                 string
}

// DebugBar has a status pill and five debugger actions.
func DebugBar(c *ui.Context, running bool, o DebugBarOptions) (root, cont, over, into, out, stop ui.Element) {
	def(&o.Continue, "Continue")
	def(&o.Over, "Step over")
	def(&o.Into, "Step in")
	def(&o.Out, "Step out")
	def(&o.Stop, "Stop")
	def(&o.Paused, "Paused")
	def(&o.Running, "Running")
	t := c.Theme()
	root = ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1.5)).Radius(999).
		Background(t.Background).Border(1, t.Border).Shadow(0, 1, 3, 0, kit.Shade(0.06)).FocusGroup(ui.Horizontal).
		Role(ui.RoleToolbar).Label("Debug")
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Icon(c, icon.Must("grip-vertical")).Size(14, 14).TextColor(t.TextMuted.Alpha(0.6))
		status, ink := o.Paused, t.Warning
		if running {
			status, ink = o.Running, t.Success
		}
		pill := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2.5)).Radius(999)
		pill.Background(pill.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).TextColor(pill.AnimateColor("ink", ink, ui.Snappy))
		pill.Children(func() {
			dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
			if running {
				ph := dot.Loop("pulse", 1200*time.Millisecond, ui.EaseInOut)
				dot.Opacity(0.4 + 0.6*float32(1-2*abs(ph-0.5)))
			}
			ui.Text(c, status).FontSize(t.Rem(0.75)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
		})
		ui.Box(c).Size(1, t.Space(5)).Margin(0, t.Space(1)).Background(t.Border)
		cont = round(c, o.Continue, "play").Disabled(running).TextColor(t.Success)
		over = round(c, o.Over, "redo-dot").Disabled(running)
		into = round(c, o.Into, "arrow-down-to-dot").Disabled(running)
		out = round(c, o.Out, "arrow-up-from-dot").Disabled(running)
		stop = round(c, o.Stop, "square").TextColor(t.Danger)
	})
	act := func(btn ui.Element, event string, enabled bool) {
		if enabled && btn.Clicked() {
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("DebugBar", event)
		}
	}
	act(cont, "continue", !running)
	act(over, "over", !running)
	act(into, "into", !running)
	act(out, "out", !running)
	act(stop, "stop", true)
	return root, cont, over, into, out, stop
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}
