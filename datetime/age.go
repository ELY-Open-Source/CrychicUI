package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/egoist/mygo/ui"
)

// AgeOptions tune an Age.
type AgeOptions struct {
	// Unit follows the number; "years" when empty.
	Unit string
}

// Age counts whole years since birth; digits roll.
func Age(c *ui.Context, birth time.Time, o AgeOptions) ui.Element {
	now := c.Now()
	if birth.IsZero() || birth.After(now) {
		panic(fmt.Sprintf("datetime: Age needs a birth date before now, not %v", birth))
	}
	def(&o.Unit, "years")
	years := now.Year() - birth.Year()
	if now.Month() < birth.Month() || now.Month() == birth.Month() && now.Day() < birth.Day() {
		years--
	}
	next := time.Date(birth.Year()+years+1, birth.Month(), birth.Day(), 0, 0, 0, 0, now.Location())
	c.After(next.Sub(now))
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(2)).Role(ui.RoleText).
		Label(fmt.Sprintf("%d %s", years, o.Unit))
	return root.Children(func() {
		motion.RollingDigits(c, int64(years))
		kit.Muted(c, o.Unit)
	})
}
