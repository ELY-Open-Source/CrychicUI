package devtools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Check is one endpoint; History is recent passes, oldest first.
type Check struct {
	Endpoint string
	State    HealthState
	Latency  time.Duration
	History  []bool
}

// HealthOptions word the panel; Count takes healthy and total.
type HealthOptions struct {
	Title, Healthy, Count, Passed string
	States                        [3]string
}

// Health lists endpoints with their recent checks.
func Health(c *ui.Context, checks []Check, o HealthOptions) ui.Element {
	def(&o.Title, "Health checks")
	if len(checks) == 0 {
		panic(fmt.Sprintf("devtools: Health %q needs checks", o.Title))
	}
	names := make([]string, len(checks))
	up, down := 0, 0
	for i, k := range checks {
		validHealth("Health", k.Endpoint, k.State)
		if k.Latency < 0 || len(k.History) > 60 {
			panic(fmt.Sprintf("devtools: Health %q needs a non-negative latency and at most 60 checks", k.Endpoint))
		}
		names[i] = k.Endpoint
		switch k.State {
		case HealthUp:
			up++
		case HealthDown:
			down++
		}
	}
	unique("Health", names)
	def(&o.Healthy, "All healthy")
	def(&o.Count, "%d of %d healthy")
	def(&o.Passed, "%s%% passed")
	for i, name := range [...]string{"Healthy", "Degraded", "Down"} {
		def(&o.States[i], name)
	}
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "heart-pulse", o.Title, func() {
			switch {
			case up == len(checks):
				kit.Chip(c, o.Healthy, t.Success)
			case down > 0:
				kit.Chip(c, fmt.Sprintf(o.Count, up, len(checks)), t.Danger)
			default:
				kit.Chip(c, fmt.Sprintf(o.Count, up, len(checks)), t.Warning)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, k := range checks {
				ink := healthInk(t, k.State)
				c.Key(k.Endpoint)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					d := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999)
					d.Background(d.AnimateColor("ink", ink, ui.Snappy))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, k.Endpoint).Font("monospace").FontSize(t.Rem(0.78)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, o.States[k.State]).FontSize(t.Rem(0.7)).Bold().TextColor(ink).Shrink(0)
							dot(c)
							kit.Numeric(c, took(k.Latency)).Font("monospace").FontSize(t.Rem(0.7)).TextColor(t.TextMuted).Shrink(0)
						})
					})
					if len(k.History) > 0 {
						history(c, k.History, o.Passed)
					}
				})
			}
		})
	})
	return root
}

// history is pass and fail ticks with a pass rate.
func history(c *ui.Context, runs []bool, format string) {
	t := c.Theme()
	passed := 0
	for _, ok := range runs {
		if ok {
			passed++
		}
	}
	// At most 60 runs, so one fail shows below 100%.
	rate := float64(passed) / float64(len(runs)) * 100
	ui.Column(c).Width(t.Space(30)).Shrink(0).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
		strip := ui.Box(c).FillWidth().Height(t.Space(4))
		reveal := strip.Keyframes("reveal", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		strip.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			w := r.W / float32(len(runs))
			gap := min(t.Space(0.5), w/3)
			for i, ok := range runs {
				if float32(i+1) > reveal*float32(len(runs))+0.999 {
					break
				}
				ink := t.Success
				if !ok {
					ink = t.Danger
				}
				p.Fill(ui.Rect{X: r.X + float32(i)*w, Y: r.Y, W: w - gap, H: r.H}, ink.Alpha(0.8), t.Space(0.25))
			}
		}))
		kit.Muted(c, fmt.Sprintf(format, strconv.FormatFloat(rate, 'f', 1, 64))).FontSize(t.Rem(0.65)).Shrink(0)
	})
}
