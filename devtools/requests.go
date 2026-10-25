package devtools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// Request is one HTTP request in a network log.
type Request struct {
	ID, Method, Path string
	Status           int
	Time             time.Duration
	Size             int64
}

// RequestsOptions word the panel; Label takes method and path.
type RequestsOptions struct{ Title, Empty, Label string }

// Requests lists requests with method, status and time; cue Select.
func Requests(c *ui.Context, requests []Request, selected *string, o RequestsOptions) ui.Element {
	if selected == nil {
		panic("devtools: Requests needs a selection")
	}
	ids := make([]string, len(requests))
	var slowest time.Duration
	for i, r := range requests {
		if r.Method == "" || r.Path == "" || r.Status < 100 || r.Status > 599 || r.Time < 0 || r.Size < 0 {
			panic(fmt.Sprintf("devtools: Requests %q needs a method, path, status 100-599 and non-negative time and size", r.ID))
		}
		ids[i] = r.ID
		slowest = max(slowest, r.Time)
	}
	if len(requests) > 0 {
		unique("Requests", ids)
	}
	def(&o.Title, "Network")
	def(&o.Empty, "No requests yet")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "network", o.Title, func() { kit.Chip(c, fmt.Sprint(len(requests)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(requests) == 0 {
				empty(c, "network", o.Empty)
			}
			for i, r := range requests {
				on := *selected == r.ID
				c.Key(r.ID)
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, r.Method, r.Path)).FillWidth().
					Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						m := methodInk(t, r.Method)
						ui.Box(c).Width(t.Space(14)).Shrink(0).Center().Padding(t.Space(0.5), 0).Radius(t.Radius * 0.6).
							Background(m.Alpha(0.12)).TextColor(m).Children(func() {
							ui.Text(c, r.Method).Font("monospace").FontSize(t.Rem(0.65)).Bold().Shrink(0)
						})
						ui.Text(c, r.Path).Font("monospace").FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
						ui.Text(c, strconv.Itoa(r.Status)).Font("monospace").FontSize(t.Rem(0.72)).Bold().
							TextColor(statusInk(t, r.Status)).Shrink(0)
						typography.Bytes(c, r.Size).FontSize(t.Rem(0.7)).TextColor(t.TextMuted).
							Width(t.Space(14)).Shrink(0).TextAlign(ui.End)
						ui.Column(c).Width(t.Space(16)).Shrink(0).Gap(t.Space(0.5)).AlignItems(ui.End).Children(func() {
							kit.Numeric(c, took(r.Time)).Font("monospace").FontSize(t.Rem(0.7)).Shrink(0)
							share := float32(0)
							if slowest > 0 {
								share = float32(float64(r.Time) / float64(slowest))
							}
							meter(c, share, statusInk(t, r.Status)).Height(t.Space(0.75))
						})
					})
				})
				if b.Clicked() && !on {
					*selected = r.ID
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Requests", "selected", "index", i)
				}
			}
		})
	})
	return root
}
