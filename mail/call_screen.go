package mail

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CallScreenOptions are the call screen labels.
type CallScreenOptions struct{ End string }

// CallScreen shows caller and timer; end cues Close.
func CallScreen(c *ui.Context, caller string, started time.Time, ended *bool, o CallScreenOptions) ui.Element {
	if caller == "" || started.IsZero() || started.After(c.Now()) || ended == nil {
		panic("mail: CallScreen needs caller, start, and state")
	}
	def(&o.End, "End")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border)
	frozen := ui.Local(root, "elapsed", func() time.Duration { return c.Now().Sub(started).Truncate(time.Second) })
	wasEnded := ui.Local(root, "wasEnded", func() bool { return *ended })
	elapsed := c.Now().Sub(started).Truncate(time.Second)
	if *ended && !*wasEnded {
		*frozen = elapsed
	}
	*wasEnded = *ended
	if *ended {
		elapsed = *frozen
	}
	if !*ended {
		c.After(time.Second - c.Now().Sub(started)%time.Second)
	}
	root.Children(func() {
		ui.Text(c, caller).Bold()
		kit.Numeric(c, fmt.Sprintf("%02d:%02d", int(elapsed.Minutes()), int(elapsed.Seconds())%60)).TextColor(t.TextMuted)
		end := button(c, o.End).Background(t.Danger).TextColor(t.AccentText).Disabled(*ended)
		end.Children(func() { ui.Text(c, o.End) })
		if end.Clicked() && !*ended {
			*ended = true
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("CallScreen", "ended")
		}
	})
	return root
}
