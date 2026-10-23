package agent

import (
	"fmt"
	"math"
	"time"

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
	name, ink := names[state], inks[state]
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(name)
	glow := float32(1)
	if state == Working {
		ph := float64(root.Loop("pulse", 1400*time.Millisecond, ui.Linear))
		glow = 0.5 + 0.5*float32(math.Abs(math.Sin(math.Pi*ph)))
	}
	return root.Children(func() {
		dot := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Background(ink.Alpha(glow))
		dot.Background(dot.AnimateColor("dot", ink, ui.Snappy)).Opacity(dot.AnimateTo("glow", glow, ui.Snappy))
		ui.Text(c, name).FontSize(t.Rem(0.9)).ContentTransition(ui.ContentTransition{})
	})
}
