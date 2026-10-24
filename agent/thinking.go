package agent

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/fx"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ThinkingOptions tune a Thinking label.
type ThinkingOptions struct {
	Label  string
	Format string
}

// Thinking is a shimmering label with the elapsed time.
func Thinking(c *ui.Context, since time.Time, o ThinkingOptions) ui.Element {
	if since.IsZero() {
		panic("agent: Thinking needs a start time")
	}
	def(&o.Label, "Thinking")
	def(&o.Format, "%ds")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3)).
		Radius(999).Background(t.Surface).Role(ui.RoleStatus).Label(o.Label)
	took := max(0, c.Now().Sub(since))
	c.After(time.Second - took%time.Second)
	return root.Children(func() {
		brain := ui.Icon(c, icon.Must("brain")).Size(16, 16).TextColor(t.Accent)
		ph := float64(brain.Loop("pulse", 1400*time.Millisecond, ui.Linear))
		s := float32(1 + 0.08*math.Sin(2*math.Pi*ph))
		brain.Scale(s, s)
		fx.ShinyText(c, o.Label)
		kit.Numeric(c, fmt.Sprintf(o.Format, int64(took/time.Second))).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
	})
}
