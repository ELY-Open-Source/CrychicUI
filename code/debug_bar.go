package code

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DebugBarOptions name the five actions.
type DebugBarOptions struct {
	Continue, Over, Into, Out, Stop string
}

// DebugBar has continue, step over, in, out and stop.
func DebugBar(c *ui.Context, running bool, o DebugBarOptions) (root, cont, over, into, out, stop ui.Element) {
	def(&o.Continue, "Continue")
	def(&o.Over, "Step over")
	def(&o.Into, "Step in")
	def(&o.Out, "Step out")
	def(&o.Stop, "Stop")
	t := c.Theme()
	root = ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
		Background(t.Surface).Border(1, t.Border).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label("Debug")
	root.Children(func() {
		cont = round(c, o.Continue, "play").Disabled(running)
		over = round(c, o.Over, "redo").Disabled(running)
		into = round(c, o.Into, "arrow-down").Disabled(running)
		out = round(c, o.Out, "arrow-up").Disabled(running)
		stop = round(c, o.Stop, "square").Disabled(!running)
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
	act(stop, "stop", running)
	return root, cont, over, into, out, stop
}
