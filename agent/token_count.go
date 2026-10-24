package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// exact is the largest count a float64 holds exactly.
const exact = 1 << 53

// TokenCountOptions tune a TokenCount.
type TokenCountOptions struct {
	Label, Left, Over string
}

// TokenCount shows tokens used of a limit; digits roll.
func TokenCount(c *ui.Context, used, limit int, o TokenCountOptions) ui.Element {
	if used < 0 || limit <= 0 || used > exact || limit > exact {
		panic(fmt.Sprintf("agent: TokenCount needs counts from 0 to 2^53, got %d of %d", used, limit))
	}
	def(&o.Label, "Tokens")
	def(&o.Left, "%d left")
	def(&o.Over, "%d over")
	t := c.Theme()
	share := float64(used) / float64(limit)
	ink := meterInk(c, share)
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleMeter).Label(o.Label).Range(0, float64(limit), float64(min(used, limit)))
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
			kit.Muted(c, o.Label).FontSize(t.Rem(0.8)).Grow(1)
			n := typography.Number(c, float64(used)).Bold()
			n.TextColor(n.AnimateColor("ink", t.Text.Mix(t.Danger, unit(used > limit)), ui.Snappy))
			kit.Muted(c, "/").FontSize(t.Rem(0.8))
			typography.Number(c, float64(limit)).TextColor(t.TextMuted).FontSize(t.Rem(0.8))
		})
		bar(c, "used", share, ink)
		note := fmt.Sprintf(o.Left, limit-used)
		if used > limit {
			note = fmt.Sprintf(o.Over, used-limit)
		}
		kit.Numeric(c, note).FontSize(t.Rem(0.75)).TextColor(ink).ContentTransition(ui.ContentTransition{})
	})
}
