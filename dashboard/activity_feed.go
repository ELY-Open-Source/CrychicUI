package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Activity is one event; Actor shows an avatar, else Icon.
type Activity struct {
	ID, Actor, Text, Time, Icon string
}

// ActivityFeedOptions word the card.
type ActivityFeedOptions struct{ Title, Empty string }

// ActivityFeed lists recent events; new rows slide in.
func ActivityFeed(c *ui.Context, events []Activity, o ActivityFeedOptions) ui.Element {
	ids := make([]string, len(events))
	for i, event := range events {
		if event.Text == "" || event.Time == "" {
			panic(fmt.Sprintf("dashboard: ActivityFeed event %q needs text and a time", event.ID))
		}
		ids[i] = event.ID
	}
	unique("ActivityFeed", ids)
	def(&o.Title, "Activity")
	def(&o.Empty, "No activity yet")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "activity", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3))
		list.Children(func() {
			if len(events) == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(5)).Transition(kit.Enter).
					Children(func() {
						ui.Icon(c, icon.Must("inbox")).Size(20, 20).TextColor(t.TextMuted)
						kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
					})
			}
			for i, event := range events {
				c.Key(event.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(events)-1)))
				kit.Stagger(c, list, row, i).Children(func() {
					if event.Actor != "" {
						display.Avatar(c, event.Actor, nil)
					} else {
						glyph := event.Icon
						if glyph == "" {
							glyph = "zap"
						}
						tile(c, glyph, palette(t, i))
					}
					spans := []ui.Span{{Text: event.Text}}
					if event.Actor != "" {
						spans = []ui.Span{{Text: event.Actor + " ", Weight: 700}, {Text: event.Text}}
					}
					ui.RichText(c, spans...).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
					kit.Muted(c, event.Time).FontSize(t.Rem(0.75)).Shrink(0)
				})
			}
		})
	})
	return root
}
