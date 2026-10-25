package devtools

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Delivery is one webhook send; Status 0 means no response.
type Delivery struct {
	ID, Event, URL, When string
	Status, Attempts     int
}

// WebhooksOptions word the panel; Retry takes the event name.
type WebhooksOptions struct{ Title, Retry, Again, Attempt, Attempts, Timeout, Failing, Delivered string }

// Webhooks lists deliveries; Retry sets retry; cue Send.
func Webhooks(c *ui.Context, deliveries []Delivery, retry *string, o WebhooksOptions) ui.Element {
	def(&o.Title, "Webhooks")
	if retry == nil || len(deliveries) == 0 {
		panic(fmt.Sprintf("devtools: Webhooks %q needs deliveries and a retry target", o.Title))
	}
	ids := make([]string, len(deliveries))
	failing := 0
	for i, d := range deliveries {
		if d.Event == "" || d.URL == "" || d.Attempts < 1 || d.Status != 0 && (d.Status < 100 || d.Status > 599) {
			panic(fmt.Sprintf("devtools: Webhooks %q needs an event, URL, one attempt and a status of 0 or 100-599", d.ID))
		}
		ids[i] = d.ID
		if failed(d) {
			failing++
		}
	}
	unique("Webhooks", ids)
	def(&o.Retry, "Retry %s")
	def(&o.Again, "Retry")
	def(&o.Attempt, "%d attempt")
	def(&o.Attempts, "%d attempts")
	def(&o.Timeout, "Timeout")
	def(&o.Failing, "%d failing")
	def(&o.Delivered, "All delivered")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "webhook", o.Title, func() {
			if failing > 0 {
				kit.Chip(c, fmt.Sprintf(o.Failing, failing), t.Danger)
				return
			}
			kit.Chip(c, o.Delivered, t.Success)
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, d := range deliveries {
				bad := failed(d)
				c.Key(d.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())).Mix(t.Danger.Alpha(0.05), unit(bad)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.75)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Box(c).Padding(t.Space(0.5), t.Space(1.5)).Radius(t.Radius * 0.6).Background(t.Accent.Alpha(0.1)).
								TextColor(t.Accent).Shrink(0).Children(func() {
								ui.Text(c, d.Event).Font("monospace").FontSize(t.Rem(0.7)).Bold().Shrink(0).MaxWidth(t.Space(32)).SingleLine()
							})
							if d.When != "" {
								kit.Muted(c, d.When).FontSize(t.Rem(0.68)).MinWidth(0).SingleLine()
							}
						})
						kit.Muted(c, d.URL).Font("monospace").FontSize(t.Rem(0.7)).FillWidth().MinWidth(0).SingleLine()
					})
					ui.Column(c).Shrink(0).Gap(t.Space(0.75)).AlignItems(ui.End).Children(func() {
						code, ink := strconv.Itoa(d.Status), t.Success
						if d.Status == 0 {
							code = o.Timeout
						}
						if bad {
							ink = t.Danger
						}
						kit.Chip(c, code, ink)
						format := o.Attempts
						if d.Attempts == 1 {
							format = o.Attempt
						}
						kit.Muted(c, fmt.Sprintf(format, d.Attempts)).FontSize(t.Rem(0.65)).Shrink(0)
					})
					if !bad {
						return
					}
					b := outline(c, fmt.Sprintf(o.Retry, d.Event), o.Again, "rotate-cw", t.Accent)
					if b.Clicked() {
						*retry = d.ID
						root.SetChanged()
						sound.Play(sound.Send)
						kit.Log("Webhooks", "retry", "index", i, "attempts", d.Attempts)
					}
				})
			}
		})
	})
	return root
}

// failed says a delivery got no 2xx reply.
func failed(d Delivery) bool { return d.Status < 200 || d.Status > 299 }
