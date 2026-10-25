package devtools

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ErrorGroup is one grouped error; Trend is recent counts.
type ErrorGroup struct {
	ID, Type, Message, Last string
	Count, Users            int64
	Trend                   []float64
	Resolved                bool
}

// ErrorsOptions word the panel; Label takes type and message.
type ErrorsOptions struct{ Title, Label, Event, Events, User, Users, Resolved, Empty string }

// Errors lists grouped errors with counts; cue Select.
func Errors(c *ui.Context, groups []ErrorGroup, selected *string, o ErrorsOptions) ui.Element {
	def(&o.Title, "Errors")
	if selected == nil {
		panic(fmt.Sprintf("devtools: Errors %q needs a selection", o.Title))
	}
	ids := make([]string, len(groups))
	var total int64
	for i, g := range groups {
		if g.Type == "" || g.Count < 1 || g.Users < 0 {
			panic(fmt.Sprintf("devtools: Errors %q needs a type, a count of at least 1 and users of at least 0", g.ID))
		}
		if len(g.Trend) > 0 {
			if lo, _ := bounds("Errors "+g.ID, g.Trend); lo < 0 {
				panic(fmt.Sprintf("devtools: Errors %q trend has a negative count", g.ID))
			}
		}
		if g.Count > math.MaxInt64-total {
			panic(fmt.Sprintf("devtools: Errors %q total count overflows", o.Title))
		}
		total += g.Count
		ids[i] = g.ID
	}
	if len(groups) > 0 {
		unique("Errors", ids)
	}
	def(&o.Label, "%s: %s")
	def(&o.Event, "%s event")
	def(&o.Events, "%s events")
	def(&o.User, "%s user")
	def(&o.Users, "%s users")
	def(&o.Resolved, "Resolved")
	def(&o.Empty, "No errors")
	t := c.Theme()
	plural := func(n int64, one, many string) string {
		if n == 1 {
			return fmt.Sprintf(one, number(c, n))
		}
		return fmt.Sprintf(many, number(c, n))
	}
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "bug", o.Title, func() { kit.Chip(c, plural(total, o.Event, o.Events), t.Danger) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(groups) == 0 {
				empty(c, "shield-check", o.Empty)
			}
			for i, g := range groups {
				on := *selected == g.ID
				ink, mark := t.Danger, "bug"
				if g.Resolved {
					ink, mark = t.Success, "circle-check"
				}
				c.Key(g.ID)
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, g.Type, g.Message)).FillWidth().Padding(t.Space(2)).
					Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						tile(c, mark, ink, t.Space(8))
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							ui.Text(c, g.Type).Font("monospace").FontSize(t.Rem(0.78)).Bold().TextColor(ink).FillWidth().MinWidth(0).SingleLine()
							if g.Message != "" {
								ui.Text(c, g.Message).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
							}
							ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
								kit.Muted(c, plural(g.Users, o.User, o.Users)).FontSize(t.Rem(0.68)).Shrink(0)
								if g.Last != "" {
									dot(c)
									kit.Muted(c, g.Last).FontSize(t.Rem(0.68)).MinWidth(0).SingleLine()
								}
							})
						})
						ui.Column(c).Width(t.Space(20)).Shrink(0).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
							if g.Resolved {
								kit.Chip(c, o.Resolved, t.Success)
							} else {
								kit.Numeric(c, number(c, g.Count)).FontSize(t.Rem(0.9)).Bold().Shrink(0).MaxWidth(t.Space(20)).SingleLine()
							}
							if len(g.Trend) > 1 {
								trend(c, g.Trend, ink)
							}
						})
					})
				})
				if b.Clicked() && !on {
					*selected = g.ID
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Errors", "selected", "index", i)
				}
			}
		})
	})
	return root
}

// trend is a small line that draws in.
func trend(c *ui.Context, points []float64, ink ui.Color) {
	t := c.Theme()
	lo, hi := bounds("trend", points)
	box := ui.Box(c).FillWidth().Height(t.Space(5))
	draw := box.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		p.StrokePath(spark(r, points, lo, hi).Trim(0, draw), t.Space(0.375), ink)
	}))
}
