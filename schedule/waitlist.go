package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WaitlistOptions word the card; Total is the line's length.
type WaitlistOptions struct {
	Label, Place, Of, Ahead, Next string
	Total                         int
}

// Waitlist shows a position with rolling digits.
func Waitlist(c *ui.Context, position int, o WaitlistOptions) ui.Element {
	if position < 1 || o.Total < 0 || (o.Total > 0 && position > o.Total) {
		panic(fmt.Sprintf("schedule: Waitlist position %d of %d is invalid", position, o.Total))
	}
	def(&o.Label, "Your place in line")
	def(&o.Place, "#%d")
	def(&o.Of, "of %d")
	def(&o.Ahead, "%d ahead of you")
	def(&o.Next, "You're next")
	t := c.Theme()
	ahead := position - 1
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "users", t.Accent)
			ui.Text(c, o.Label).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0)
		})
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Place, position)).FontSize(t.Rem(2.4)).Bold().TextColor(t.Accent)
			if o.Total > 0 {
				kit.Muted(c, fmt.Sprintf(o.Of, o.Total)).FontSize(t.Rem(0.9)).Padding(0, 0, t.Space(1.5), 0)
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			shown := min(ahead, 9)
			for i := range shown {
				c.Key(i)
				dot := ui.Box(c).Size(t.Space(5), t.Space(5)).Radius(999).Background(t.TextMuted.Alpha(0.25))
				kit.Arrive(dot, ui.Pose{Scale: 0.4, Fade: 1}, time.Duration(i)*30*time.Millisecond)
			}
			c.Key("you")
			ui.Box(c).Size(t.Space(7), t.Space(7)).Radius(999).Background(t.Accent).Border(t.Space(0.75), t.Accent.Alpha(0.3))
		})
		text := fmt.Sprintf(o.Ahead, ahead)
		if ahead == 0 {
			text = o.Next
		}
		kit.Muted(c, text).FontSize(t.Rem(0.8)).ContentTransition(ui.ContentTransition{})
	})
	return root
}
