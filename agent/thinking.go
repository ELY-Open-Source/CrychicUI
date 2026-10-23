package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/fx"
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
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(o.Label)
	took := c.Now().Sub(since)
	c.After(time.Second - took%time.Second)
	return root.Children(func() {
		fx.ShinyText(c, o.Label)
		kit.Numeric(c, fmt.Sprintf(o.Format, int64(took/time.Second))).FontSize(t.Rem(0.85)).TextColor(t.TextMuted)
	})
}
