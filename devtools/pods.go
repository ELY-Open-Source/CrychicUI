package devtools

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PodStatus is a container's state.
type PodStatus int

const (
	PodRunning PodStatus = iota
	PodPending
	PodFailed
	PodTerminating
)

// Pod is one container; CPU and Memory are 0..1.
type Pod struct {
	Name, Age   string
	Status      PodStatus
	Restarts    int
	CPU, Memory float64
}

// PodsOptions word the panel; Ready takes ready and total.
type PodsOptions struct {
	Title, Ready, Restarts, CPU, Memory string
	Statuses                            [4]string
}

// Pods lists containers with status and load.
func Pods(c *ui.Context, pods []Pod, o PodsOptions) ui.Element {
	if len(pods) == 0 {
		panic("devtools: Pods needs pods")
	}
	names := make([]string, len(pods))
	ready := 0
	for i, p := range pods {
		if p.Age == "" || p.Status < PodRunning || p.Status > PodTerminating || p.Restarts < 0 ||
			!finite(p.CPU) || !finite(p.Memory) || p.CPU < 0 || p.CPU > 1 || p.Memory < 0 || p.Memory > 1 {
			panic(fmt.Sprintf("devtools: Pods %q needs an age, a status, restarts of at least 0 and loads in 0..1", p.Name))
		}
		names[i] = p.Name
		if p.Status == PodRunning {
			ready++
		}
	}
	unique("Pods", names)
	def(&o.Title, "Pods")
	def(&o.Ready, "%d of %d ready")
	def(&o.Restarts, "%d restarts")
	def(&o.CPU, "CPU")
	def(&o.Memory, "Mem")
	for i, name := range [...]string{"Running", "Pending", "Failed", "Terminating"} {
		def(&o.Statuses[i], name)
	}
	t := c.Theme()
	inks := [...]ui.Color{t.Success, t.Warning, t.Danger, t.TextMuted}
	load := func(v float64) ui.Color {
		switch {
		case v >= 0.9:
			return t.Danger
		case v >= 0.7:
			return t.Warning
		}
		return t.Accent
	}
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		tone := t.Success
		if ready < len(pods) {
			tone = t.Warning
		}
		kit.CardHead(c, "box", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Ready, ready, len(pods)), tone) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, p := range pods {
				ink := inks[p.Status]
				c.Key(p.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					d := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999)
					d.Background(d.AnimateColor("ink", ink, ui.Snappy))
					if p.Status == PodPending {
						d.Opacity(0.4 + 0.6*float32(d.Loop("pulse", 1200*time.Millisecond, ui.EaseInOut)))
					}
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, p.Name).Font("monospace").FontSize(t.Rem(0.78)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, o.Statuses[p.Status]).FontSize(t.Rem(0.7)).Bold().TextColor(ink).Shrink(0)
							dot(c)
							kit.Muted(c, p.Age).FontSize(t.Rem(0.7)).Shrink(0)
							if p.Restarts > 0 {
								dot(c)
								ui.Text(c, fmt.Sprintf(o.Restarts, p.Restarts)).FontSize(t.Rem(0.7)).TextColor(t.Warning).MinWidth(0).SingleLine()
							}
						})
					})
					for _, m := range []struct {
						name string
						v    float64
					}{{o.CPU, p.CPU}, {o.Memory, p.Memory}} {
						ui.Column(c).Width(t.Space(16)).Shrink(0).Gap(t.Space(0.75)).Children(func() {
							ui.Row(c).FillWidth().Children(func() {
								kit.Muted(c, m.name).FontSize(t.Rem(0.62)).Grow(1).MinWidth(0).SingleLine()
								kit.Numeric(c, fmt.Sprintf("%.0f%%", m.v*100)).FontSize(t.Rem(0.62)).Shrink(0)
							})
							meter(c, float32(m.v), load(m.v)).Height(t.Space(1))
						})
					}
				})
			}
		})
	})
	return root
}
