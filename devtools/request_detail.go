package devtools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Header is one HTTP header.
type Header struct {
	Name, Value string
}

// RequestInfo is one request's detail; Body is optional.
type RequestInfo struct {
	Method, URL string
	Status      int
	Time        time.Duration
	Headers     []Header
	Body        string
}

// RequestDetailOptions word the tabs; Headers takes the count.
type RequestDetailOptions struct{ Headers, Body, NoBody string }

// RequestDetail shows headers and body tabs; cue Select.
func RequestDetail(c *ui.Context, req RequestInfo, tab *int, o RequestDetailOptions) ui.Element {
	if req.Method == "" || req.URL == "" || req.Status < 100 || req.Status > 599 || req.Time < 0 || tab == nil || *tab < 0 || *tab > 1 {
		panic(fmt.Sprintf("devtools: RequestDetail %q needs a method, URL, status 100-599 and tab 0 or 1", req.URL))
	}
	def(&o.Headers, "Headers %d")
	def(&o.Body, "Body")
	def(&o.NoBody, "No body")
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(4)).Children(func() {
			m := methodInk(t, req.Method)
			ui.Text(c, req.Method).Font("monospace").FontSize(t.Rem(0.75)).Bold().TextColor(m).Shrink(0)
			ui.Text(c, req.URL).Font("monospace").FontSize(t.Rem(0.78)).Grow(1).MinWidth(0).SingleLine()
			kit.Chip(c, strconv.Itoa(req.Status), statusInk(t, req.Status))
			kit.Muted(c, took(req.Time)).Font("monospace").FontSize(t.Rem(0.72)).Shrink(0)
		})
		names := []string{fmt.Sprintf(o.Headers, len(req.Headers)), o.Body}
		bar := ui.Row(c).FillWidth().Padding(0, t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Role(ui.RoleTabList)
		bar.Children(func() {
			for i, name := range names {
				on := *tab == i
				b := kit.Press(ui.ButtonBase(c).Label(name).Column().Gap(t.Space(1.5)).Padding(t.Space(2), t.Space(3), 0, t.Space(3)).
					Role(ui.RoleTab).Checked(on))
				b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, max(unit(on), unit(b.Hovered()))), ui.Snappy))
				b.Children(func() {
					ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().Shrink(0)
					line := ui.Box(c).FillWidth().Height(t.Space(0.5)).Radius(999)
					line.Background(line.AnimateColor("line", t.Accent.Alpha(0).Mix(t.Accent, unit(on)), ui.Snappy))
				})
				if b.Clicked() && !on {
					*tab = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("RequestDetail", "tab", "index", i)
				}
			}
		})
		c.Key(*tab)
		body := ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(4), t.Space(4), t.Space(4))
		kit.Arrive(body, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
			if *tab == 0 {
				for i, h := range req.Headers {
					c.Key(i)
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start).Padding(t.Space(1), 0).Children(func() {
						kit.Muted(c, h.Name).Font("monospace").FontSize(t.Rem(0.72)).Width(t.Space(36)).Shrink(0).Wrap()
						ui.Text(c, h.Value).Font("monospace").FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).Wrap()
					})
				}
				return
			}
			if req.Body == "" {
				empty(c, "file-x", o.NoBody)
				return
			}
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(60)).Padding(t.Space(2.5)).Radius(t.Radius).Background(t.Surface).Children(func() {
				ui.Text(c, req.Body).Font("monospace").FontSize(t.Rem(0.75)).FillWidth().Wrap()
			})
		})
	})
	return root
}
