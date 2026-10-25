package devtools

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Span is one timed operation; roots have no Parent.
type Span struct {
	ID, Parent, Name, Service string
	Start, Duration           time.Duration
	Error                     bool
}

// TracesOptions word the panel.
type TracesOptions struct{ Title string }

// Traces is a span waterfall that picks spans; cue Select.
func Traces(c *ui.Context, spans []Span, selected *string, o TracesOptions) ui.Element {
	if len(spans) == 0 || selected == nil {
		panic("devtools: Traces needs spans and a selection")
	}
	def(&o.Title, "Trace")
	ids := make([]string, len(spans))
	depth := map[string]int{}
	var total time.Duration
	for i, s := range spans {
		if s.Name == "" || s.Service == "" || s.Start < 0 || s.Duration < 0 {
			panic(fmt.Sprintf("devtools: Traces span %q needs a name, service and non-negative times", s.ID))
		}
		if s.Parent != "" {
			d, ok := depth[s.Parent]
			if !ok {
				panic(fmt.Sprintf("devtools: Traces span %q comes before its parent %q", s.ID, s.Parent))
			}
			depth[s.ID] = d + 1
		} else {
			depth[s.ID] = 0
		}
		if s.Start > math.MaxInt64-s.Duration {
			panic(fmt.Sprintf("devtools: Traces span %q ends past the largest duration", s.ID))
		}
		ids[i] = s.ID
		total = max(total, s.Start+s.Duration)
	}
	unique("Traces", ids)
	if total <= 0 {
		panic(fmt.Sprintf("devtools: Traces %q needs a span with time", o.Title))
	}
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger.Mix(t.Accent, 0.5)}
	services := map[string]ui.Color{}
	for _, s := range spans {
		if _, ok := services[s.Service]; !ok {
			services[s.Service] = inks[len(services)%len(inks)]
		}
	}
	name := t.Space(44)
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "waypoints", o.Title, func() { kit.Chip(c, took(total), t.Accent) })
		ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Box(c).Width(name).Shrink(0)
			ui.Row(c).Grow(1).MinWidth(0).Children(func() {
				// Each quarter starts with its tick and label.
				for k := range 4 {
					tick := total/4*time.Duration(k) + total%4*time.Duration(k)/4
					ui.Row(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(1, t.Space(2.5)).Shrink(0).Background(t.Border)
						kit.Muted(c, took(tick)).Font("monospace").FontSize(t.Rem(0.65)).MinWidth(0).SingleLine()
					})
				}
			})
			ui.Box(c).Width(t.Space(14)).Shrink(0)
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1))
		list.Children(func() {
			for i, s := range spans {
				on := *selected == s.ID
				ink := services[s.Service]
				if s.Error {
					ink = t.Danger
				}
				c.Key(s.ID)
				b := kit.Press(ui.ButtonBase(c).Label(s.Name).FillWidth().Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Row(c).Width(name).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).
							Padding(0, 0, 0, t.Space(3)*float32(depth[s.ID])).Children(func() {
							ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(ink)
							ui.Text(c, s.Name).Font("monospace").FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
						})
						lane := ui.Box(c).Grow(1).MinWidth(0).Height(t.Space(5))
						grow := lane.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 500 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
						x0 := float32(float64(s.Start) / float64(total))
						w := float32(float64(s.Duration) / float64(total))
						lane.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
							bw := max(r.W*w*grow, t.Space(0.5))
							p.Fill(ui.Rect{X: r.X + r.W*x0, Y: r.Y + r.H*0.2, W: bw, H: r.H * 0.6}, ink.Alpha(0.85), t.Radius*0.5)
						}))
						kit.Numeric(c, took(s.Duration)).Font("monospace").FontSize(t.Rem(0.68)).TextColor(t.TextMuted).
							Width(t.Space(14)).Shrink(0).TextAlign(ui.End)
					})
				})
				if b.Clicked() && !on {
					*selected = s.ID
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Traces", "selected", "index", i)
				}
			}
		})
	})
	return root
}
