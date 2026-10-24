package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ServerLoadOptions word the card; Host names the machine.
type ServerLoadOptions struct{ Title, CPU, Memory, Host string }

// ServerLoad shows CPU and memory bars that warm with load.
func ServerLoad(c *ui.Context, cpu, memory float64, o ServerLoadOptions) ui.Element {
	if cpu < 0 || cpu > 100 || memory < 0 || memory > 100 || !finite(cpu) || !finite(memory) {
		panic("dashboard: ServerLoad needs CPU and memory in [0,100]")
	}
	def(&o.Title, "Server load")
	def(&o.CPU, "CPU")
	def(&o.Memory, "Memory")
	t := c.Theme()
	heat := func(v float64) ui.Color {
		switch {
		case v >= 90:
			return t.Danger
		case v >= 70:
			return t.Warning
		}
		return t.Accent
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "server", o.Title, func() {
			if o.Host != "" {
				kit.Chip(c, o.Host, t.TextMuted)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3))
		list.Children(func() {
			for i, row := range []struct {
				name, glyph string
				value       float64
			}{{o.CPU, "cpu", cpu}, {o.Memory, "memory-stick", memory}} {
				c.Key(i)
				ink := heat(row.value)
				line := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
				kit.Stagger(c, list, line, i).Children(func() {
					tile(c, row.glyph, ink)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
							ui.Text(c, row.name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							kit.Numeric(c, fmt.Sprintf("%.0f%%", row.value)).Bold().FontSize(t.Rem(0.85)).TextColor(ink).Shrink(0).
								ContentTransition(ui.ContentTransition{Numeric: true})
						})
						progress(c, row.value/100, ink)
					})
				})
			}
		})
	})
	return root
}
