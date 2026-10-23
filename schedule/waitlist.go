package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WaitlistOptions are the waitlist labels.
type WaitlistOptions struct{ Label string }

// Waitlist shows a position with rolling digits.
func Waitlist(c *ui.Context, position int, o WaitlistOptions) ui.Element {
	if position < 1 {
		panic("schedule: Waitlist position must be positive")
	}
	def(&o.Label, "Position")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		kit.Muted(c, o.Label).FontSize(t.Rem(0.75)).Grow(1)
		kit.Numeric(c, fmt.Sprint(position)).Bold().FontSize(t.Rem(1.3))
	})
	return root
}
