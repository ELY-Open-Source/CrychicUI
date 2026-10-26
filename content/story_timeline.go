package content

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Beat is a dated moment in a story.
type Beat struct {
	Date              time.Time
	Title, Text, Icon string
}

// StoryTimelineOptions set the date layout.
type StoryTimelineOptions struct{ Date string }

// StoryTimeline is dated beats along a rail, oldest first.
func StoryTimeline(c *ui.Context, beats []Beat, o StoryTimelineOptions) ui.Element {
	def(&o.Date, "Jan 2006")
	if len(beats) == 0 {
		panic("content: StoryTimeline needs beats")
	}
	for i, b := range beats {
		if strings.TrimSpace(b.Title) == "" || b.Icon == "" || b.Date.IsZero() {
			panic(fmt.Sprintf("content: StoryTimeline beat %q needs a title, an icon and a date", b.Title))
		}
		if i > 0 && b.Date.Before(beats[i-1].Date) {
			panic(fmt.Sprintf("content: StoryTimeline beat %q comes before the one above it", b.Title))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleList)
	root.Children(func() {
		for i, b := range beats {
			c.Key(i)
			last := i == len(beats)-1
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch).Role(ui.RoleListItem)
			kit.Stagger(c, root, row, i).Children(func() {
				ui.Column(c).Shrink(0).AlignItems(ui.Center).Children(func() {
					ink := t.TextMuted
					if last {
						ink = t.Accent
					}
					kit.Arrive(tile(c, b.Icon, ink, t.Space(8)).Radius(999), ui.Pose{Scale: 0.6, Fade: 1}, 0)
					if !last {
						ui.Box(c).Width(1).Grow(1).MinHeight(t.Space(4)).Background(t.Border)
					}
				})
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Padding(t.Space(0.5), 0, t.Space(5), 0).Children(func() {
					ui.Text(c, b.Date.Format(o.Date)).FontSize(t.Rem(0.7)).Bold().TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine()
					ui.Text(c, b.Title).FontSize(t.Rem(0.9)).Bold().FillWidth().Wrap()
					if b.Text != "" {
						kit.Muted(c, b.Text).FontSize(t.Rem(0.8)).LineHeight(1.5).FillWidth().Wrap()
					}
				})
			})
		}
	})
	return root
}
