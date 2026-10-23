package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// TokenCountOptions tune a TokenCount.
type TokenCountOptions struct {
	Label string
}

// TokenCount shows tokens used of a limit; digits roll.
func TokenCount(c *ui.Context, used, limit int, o TokenCountOptions) ui.Element {
	if used < 0 || limit <= 0 {
		panic(fmt.Sprintf("agent: TokenCount needs a count and a limit, got %d of %d", used, limit))
	}
	def(&o.Label, "tokens")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleStatus)
	return root.Children(func() {
		n := typography.Number(c, float64(used))
		if used > limit {
			n.TextColor(t.Danger)
		}
		kit.Muted(c, "/").FontSize(t.Rem(0.85))
		typography.Number(c, float64(limit)).TextColor(t.TextMuted).FontSize(t.Rem(0.85))
		kit.Muted(c, o.Label).FontSize(t.Rem(0.85))
	})
}
