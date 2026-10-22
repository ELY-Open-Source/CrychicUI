package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AnnouncementBarOptions tune an AnnouncementBar.
type AnnouncementBarOptions struct {
	// Every is how long each message shows; 4 s.
	Every time.Duration
}

// AnnouncementBar slides messages up in turn; Changed marks rotations.
func AnnouncementBar(c *ui.Context, messages []string, o AnnouncementBarOptions) ui.Element {
	if len(messages) == 0 {
		panic("feedback: AnnouncementBar needs messages")
	}
	for _, m := range messages {
		if m == "" {
			panic("feedback: AnnouncementBar messages need text")
		}
	}
	if o.Every == 0 {
		o.Every = 4 * time.Second
	}
	if o.Every < 0 {
		panic("feedback: AnnouncementBar needs a positive Every")
	}
	t := c.Theme()
	h := t.Space(9)
	root := ui.Box(c).FillWidth().Height(h).Background(t.Accent).TextColor(t.AccentText).ClipY().Role(ui.RoleStatus)
	start := ui.Local(root, "start", func() time.Time { return c.Now() })
	n := int(c.Now().Sub(*start) / o.Every)
	if len(messages) > 1 {
		c.After(start.Add(time.Duration(n+1) * o.Every).Sub(c.Now()))
	}
	i := n % len(messages)
	last := ui.Local(root, "last", func() int { return n })
	if n != *last {
		*last = n
		kit.Log("AnnouncementBar", "rotated", "index", i)
		root.SetChanged()
	}
	slide := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{Y: h}, Exit: &ui.Motion{Y: -h}}
	return root.Label(messages[i]).Children(func() {
		c.Key(n)
		ui.Text(c, messages[i]).Absolute().Top(0).Bottom(0).Left(0).Right(0).TextAlign(ui.Center).Padding(t.Space(2.5), t.Space(4)).
			Bold().FontSize(t.Rem(0.9)).Transition(slide)
	})
}
