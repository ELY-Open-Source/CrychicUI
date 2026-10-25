package devtools

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PingSample is one round trip; Lost means no reply.
type PingSample struct {
	RTT  time.Duration
	Lost bool
}

// PingHost is one target; Seq counts every sample taken.
type PingHost struct {
	Name    string
	Seq     int
	Samples []PingSample
}

// PingOptions word the panel; Loss takes a percent.
type PingOptions struct{ Title, Avg, Loss, Down string }

// Ping shows recent latency as dots; new samples pop in.
func Ping(c *ui.Context, hosts []PingHost, o PingOptions) ui.Element {
	def(&o.Title, "Ping")
	if len(hosts) == 0 {
		panic(fmt.Sprintf("devtools: Ping %q needs hosts", o.Title))
	}
	names := make([]string, len(hosts))
	for i, h := range hosts {
		if h.Seq < len(h.Samples) || len(h.Samples) > 120 {
			panic(fmt.Sprintf("devtools: Ping %q needs at most 120 samples and a count at least that", h.Name))
		}
		for _, s := range h.Samples {
			if s.RTT < 0 {
				panic(fmt.Sprintf("devtools: Ping %q has a negative round trip", h.Name))
			}
		}
		names[i] = h.Name
	}
	unique("Ping", names)
	def(&o.Avg, "%s avg")
	def(&o.Loss, "%s%% loss")
	def(&o.Down, "No replies")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "activity", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2), t.Space(3))
		list.Children(func() {
			for i, h := range hosts {
				c.Key(h.Name)
				row := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(1.5), 0)
				kit.Stagger(c, list, row, i).Children(func() { pingRow(c, h, o) })
			}
		})
	})
	return root
}

// pingRow is one host's summary and dots.
func pingRow(c *ui.Context, h PingHost, o PingOptions) {
	t := c.Theme()
	var sum float64
	var slowest time.Duration
	got, lost := 0, 0
	for _, s := range h.Samples {
		if s.Lost {
			lost++
			continue
		}
		sum += float64(s.RTT)
		slowest = max(slowest, s.RTT)
		got++
	}
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		ui.Text(c, h.Name).Font("monospace").FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
		if len(h.Samples) > 0 && got == 0 {
			kit.Chip(c, o.Down, t.Danger)
			return
		}
		if got > 0 {
			// The float mean may round past the slowest sample.
			avg := slowest
			if mean := sum / float64(got); mean < float64(slowest) {
				avg = time.Duration(mean)
			}
			kit.Numeric(c, fmt.Sprintf(o.Avg, took(avg))).Font("monospace").FontSize(t.Rem(0.72)).TextColor(rttInk(t, avg)).Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
		}
		if lost > 0 {
			pct := math.Max(1, math.Round(float64(lost)/float64(len(h.Samples))*100))
			kit.Chip(c, fmt.Sprintf(o.Loss, strconv.FormatFloat(min(pct, 99), 'f', 0, 64)), t.Warning)
		}
	})
	size, gap := t.Space(2.5), t.Space(1)
	dots := ui.Row(c).FillWidth().Gap(gap).AlignItems(ui.Center).Height(size)
	w := dots.LayoutBounds().W
	kit.Settle(c, dots, []float32{w}, func() []float32 { return []float32{dots.LayoutBounds().W} })
	fit := min(len(h.Samples), int((w+gap)/(size+gap)))
	first := len(h.Samples) - fit
	dots.Children(func() {
		for i := first; i < len(h.Samples); i++ {
			s := h.Samples[i]
			c.Key(h.Seq - len(h.Samples) + i)
			d := ui.Box(c).Size(size, size).Shrink(0).Radius(999)
			if s.Lost {
				d.Border(t.Space(0.5), t.Danger)
			} else {
				d.Background(rttInk(t, s.RTT))
			}
			kit.Arrive(d, ui.Pose{Scale: 0.2, Fade: 1}, 0)
		}
	})
}

// rttInk colors a round trip by speed.
func rttInk(t *ui.Theme, d time.Duration) ui.Color {
	switch {
	case d < 80*time.Millisecond:
		return t.Success
	case d < 200*time.Millisecond:
		return t.Warning
	}
	return t.Danger
}
