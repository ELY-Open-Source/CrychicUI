package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// XPBarOptions word the bar; XP takes earned and needed.
type XPBarOptions struct{ Level, XP, Up string }

// XPBar shows level progress; a new level plays Success.
func XPBar(c *ui.Context, level int, xp, need int64, o XPBarOptions) ui.Element {
	def(&o.Level, "Level %d")
	if level < 1 || need <= 0 || xp < 0 || xp >= need {
		panic(fmt.Sprintf("game: XPBar %q needs a level of 1 or more and 0 <= xp < need", fmt.Sprintf(o.Level, level)))
	}
	def(&o.XP, "%s / %s XP")
	def(&o.Up, "Level up!")
	t := c.Theme()
	root := card(c).Row().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(4)).Role(ui.RoleMeter).
		Label(fmt.Sprintf(o.Level, level))
	last := ui.Local(root, "last", func() int { return level })
	upAt := ui.Local(root, "up", func() time.Time { return time.Time{} })
	now := c.Now()
	if level > *last {
		*upAt = now
		sound.Play(sound.Success)
		kit.Log("XPBar", "levelled", "level", level)
	}
	*last = level
	fresh := !upAt.IsZero() && now.Sub(*upAt) < 2*time.Second
	if fresh {
		c.After(2*time.Second - now.Sub(*upAt))
	}
	root.Children(func() {
		c.Key(level)
		badge := ui.Box(c).Size(t.Space(11), t.Space(11)).Shrink(0).Center().Radius(999).Background(t.Warning.Alpha(0.15)).
			Border(t.Space(0.5), t.Warning).TextColor(t.Warning.Mix(t.Text, 0.3))
		kit.Arrive(badge, ui.Pose{Scale: 0.5, Rotate: -30, Fade: 1}, 0).Children(func() {
			kit.Numeric(c, fmt.Sprint(level)).FontSize(t.Rem(1)).Bold().MaxWidth(t.Space(9)).SingleLine()
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf(o.Level, level)).FontSize(t.Rem(0.82)).Bold().Shrink(0).MaxWidth(t.Space(24)).SingleLine()
				c.Key("up")
				up := kit.Chip(c, o.Up, t.Warning)
				up.Opacity(up.AnimateTo("shown", unit(fresh), ui.Smooth))
				kit.Muted(c, fmt.Sprintf(o.XP, number(c, xp), number(c, need))).Font("monospace").FontSize(t.Rem(0.7)).
					Grow(1).MinWidth(0).TextAlign(ui.End).SingleLine()
			})
			// Each level gets a fresh bar that grows from empty.
			c.Key(level)
			meter(c, float32(float64(xp)/float64(need)), t.Warning).Height(t.Space(2.5))
		})
	})
	return root
}
