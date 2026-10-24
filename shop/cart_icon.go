package shop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CartIconOptions name the cart; Items words the count aloud.
type CartIconOptions struct{ Label, Items string }

// bumpTime is how long a new count enlarges the badge.
const bumpTime = 120 * time.Millisecond

// CartIcon is a bag whose count badge rolls and bumps.
func CartIcon(c *ui.Context, count int, o CartIconOptions) ui.Element {
	if count < 0 {
		panic(fmt.Sprintf("shop: CartIcon count %d must not be negative", count))
	}
	def(&o.Label, "Cart")
	def(&o.Items, "%s, %d items")
	t := c.Theme()
	root := ui.Grid(c).Columns(1).AlignSelf(ui.Start).Size(t.Space(11), t.Space(11)).Role(ui.RoleGroup).Label(fmt.Sprintf(o.Items, o.Label, count))
	last := ui.Local(root, "last", func() int { return count })
	bumped := ui.Local(root, "bumped", func() time.Time { return time.Time{} })
	if count != *last {
		*last = count
		*bumped = c.Now()
	}
	big := c.Now().Sub(*bumped) < bumpTime
	if big {
		c.After(bumpTime - c.Now().Sub(*bumped))
	}
	root.Children(func() {
		ui.Box(c).ColumnStart(1).RowStart(1).Size(t.Space(10), t.Space(10)).AlignSelf(ui.End).Center().Radius(999).Background(t.Surface).
			Children(func() { ui.Icon(c, icon.Must("shopping-bag")).Size(20, 20) })
		if count == 0 {
			return
		}
		c.Key("badge")
		badge := ui.Box(c).ColumnStart(1).RowStart(1).AlignSelf(ui.Start).JustifySelf(ui.End).MinWidth(t.Space(5)).Height(t.Space(5)).
			Padding(0, t.Space(1.25)).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText).Border(t.Space(0.5), t.Background).
			Transition(kit.Pop)
		s := badge.AnimateTo("bump", 1+0.35*unit(big), ui.Bouncy)
		kit.Arrive(badge.Scale(s, s), ui.Pose{Scale: 0.6, Fade: 1}, 0).Children(func() {
			label := strconv.Itoa(count)
			if count > 99 {
				label = "99+"
			}
			ui.Text(c, label).FontSize(t.Rem(0.65)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root
}
