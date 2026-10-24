package agent

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RunState is where an agent's run stands.
type RunState uint8

// Run states.
const (
	Idle RunState = iota
	Working
	Blocked
	Done
	runStateCount
)

// AgentStatusOptions name the states.
type AgentStatusOptions struct {
	Idle, Working, Blocked, Done string
}

// AgentStatus shows a state dot and text.
func AgentStatus(c *ui.Context, state RunState, o AgentStatusOptions) ui.Element {
	if state >= runStateCount {
		panic(fmt.Sprintf("agent: AgentStatus state %d unknown", state))
	}
	def(&o.Idle, "Idle")
	def(&o.Working, "Working")
	def(&o.Blocked, "Blocked")
	def(&o.Done, "Done")
	t := c.Theme()
	names := [4]string{o.Idle, o.Working, o.Blocked, o.Done}
	inks := [4]ui.Color{t.TextMuted, t.Accent, t.Warning, t.Success}
	name := names[state]
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).
		Radius(999).Role(ui.RoleStatus).Label(name)
	ink := root.AnimateColor("ink", inks[state], ui.Snappy)
	root.Background(ink.Alpha(0.12)).TextColor(ink)
	halo := float32(0)
	if state == Working {
		halo = float32(root.Loop("pulse", 1400*time.Millisecond, ui.Linear))
	}
	return root.Children(func() {
		dot := ui.Box(c).Size(t.Space(4), t.Space(4)).Shrink(0).Margin(-t.Space(1))
		dot.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/4
			if halo > 0 {
				grow := rad * (1 + float32(math.Sqrt(float64(halo))))
				p.Fill(ui.Rect{X: cx - grow, Y: cy - grow, W: 2 * grow, H: 2 * grow}, ink.Alpha(0.4*(1-halo)), 999)
			}
			p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: 2 * rad, H: 2 * rad}, ink, 999)
		}))
		ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().ContentTransition(ui.ContentTransition{})
	})
}
