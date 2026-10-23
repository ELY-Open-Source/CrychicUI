package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BranchesOptions word the position and buttons.
type BranchesOptions struct {
	Format      string
	Back, Forth string
}

// Branches switches between answers with a position.
func Branches(c *ui.Context, current *int, total int, o BranchesOptions) (root, back, forth ui.Element) {
	if current == nil || total < 1 || *current < 1 || *current > total {
		panic(fmt.Sprintf("agent: Branches needs a position from 1 to %d", total))
	}
	def(&o.Format, "%d / %d")
	def(&o.Back, "Previous")
	def(&o.Forth, "Next")
	t := c.Theme()
	backGlyph, forthGlyph := ends(c)
	root = ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup).
		Label(fmt.Sprintf(o.Format, *current, total))
	root.Children(func() {
		back = round(c, o.Back, backGlyph).Disabled(*current == 1)
		kit.Numeric(c, fmt.Sprintf(o.Format, *current, total)).FontSize(t.Rem(0.85)).TextColor(t.TextMuted)
		forth = round(c, o.Forth, forthGlyph).Disabled(*current == total)
	})
	turn := func(d int) {
		*current += d
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Branches", "at", "index", *current)
	}
	if back.Clicked() && *current > 1 {
		turn(-1)
	}
	if forth.Clicked() && *current < total {
		turn(1)
	}
	return root, back, forth
}
