package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// AlarmRowOptions tune an AlarmRow.
type AlarmRowOptions struct {
	// Days names the day toggles.
	Days choice.DayTogglesOptions
}

// AlarmRow shows time, days and a switch; off dims it.
func AlarmRow(c *ui.Context, label string, at time.Time, days *[7]bool, on *bool, o AlarmRowOptions) ui.Element {
	if label == "" || days == nil || on == nil {
		panic("datetime: AlarmRow needs a label, days and a switch")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(4)).
		Radius(t.Radius).Border(1, t.Border)
	return root.Children(func() {
		info := ui.Column(c).Grow(1).Gap(t.Space(2))
		info.Opacity(info.AnimateTo("lit", 0.45+0.55*unit(*on), ui.Smooth)).Children(func() {
			ui.Column(c).Children(func() {
				ui.Text(c, i18n.FormatDate(locale(c), at, i18n.DateOptions{Skeleton: "jm"})).FontSize(t.Rem(1.75)).FontWeight(600)
				kit.Muted(c, label)
			})
			if choice.DayToggles(c, label, days, o.Days).Changed() {
				root.SetChanged()
				kit.Log("AlarmRow", "days")
			}
		})
		if kit.Changed(ui.Switch(c, on).Label(label), sound.Toggle) {
			root.SetChanged()
			kit.Log("AlarmRow", "switched", "on", *on)
		}
	})
}
