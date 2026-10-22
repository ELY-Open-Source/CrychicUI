package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// DigitalClockOptions tune a DigitalClock.
type DigitalClockOptions struct {
	// Location is the zone shown; local time when nil.
	Location *time.Location
	// Seconds adds a seconds field.
	Seconds bool
}

// DigitalClock rolls its digits; the colon blinks each second.
func DigitalClock(c *ui.Context, o DigitalClockOptions) ui.Element {
	if o.Location == nil {
		o.Location = time.Local
	}
	t := c.Theme()
	now := c.Now().In(o.Location)
	c.After(time.Second - time.Duration(now.Nanosecond()))
	skeleton := "jm"
	if o.Seconds {
		skeleton = "jms"
	}
	label := i18n.FormatDate(locale(c), now, i18n.DateOptions{Skeleton: skeleton})
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).FontSize(t.Rem(2)).FontWeight(600).
		Role(ui.RoleText).Label(label)
	lit := unit(c.Preferences().ReduceMotion || now.Second()%2 == 0)
	return root.Children(func() {
		hour := now.Hour()
		h12 := twelve(c)
		if h12 {
			hour = (hour+11)%12 + 1
		}
		kit.Numeric(c, fmt.Sprintf("%02d", hour))
		colon := ui.Text(c, ":")
		colon.Opacity(colon.AnimateTo("blink", lit, ui.Smooth))
		kit.Numeric(c, fmt.Sprintf("%02d", now.Minute()))
		if o.Seconds {
			ui.Text(c, ":")
			kit.Numeric(c, fmt.Sprintf("%02d", now.Second()))
		}
		if h12 {
			am, pm := i18n.DayPeriods(locale(c), i18n.Abbreviated)
			period := am
			if now.Hour() >= 12 {
				period = pm
			}
			ui.Text(c, period).FontSize(t.Rem(1)).TextColor(t.TextMuted).Margin(0, 0, 0, t.Space(2))
		}
	})
}
