package datetime

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// zones are common IANA zones, the default choices.
var zones = []string{"UTC", "Pacific/Honolulu", "America/Anchorage", "America/Los_Angeles", "America/Denver",
	"America/Chicago", "America/New_York", "America/Sao_Paulo", "Europe/London", "Europe/Paris", "Europe/Berlin",
	"Europe/Moscow", "Africa/Cairo", "Africa/Johannesburg", "Asia/Dubai", "Asia/Kolkata", "Asia/Bangkok",
	"Asia/Shanghai", "Asia/Singapore", "Asia/Tokyo", "Australia/Sydney", "Pacific/Auckland"}

// ZoneSelectOptions tune a ZoneSelect.
type ZoneSelectOptions struct {
	// Zones are the choices; common IANA zones when nil.
	Zones []string
	// Search tunes the filter's copy.
	Search choice.SearchDropdownOptions
}

// ZoneSelect picks a time zone from a filtered list.
func ZoneSelect(c *ui.Context, label string, zone *string, o ZoneSelectOptions) ui.Element {
	if o.Zones == nil {
		o.Zones = zones
	}
	root := kit.Arrive(ui.Column(c).AlignSelf(ui.Start), ui.Pose{Y: c.Theme().Space(2), Fade: 1}, 0)
	return root.Children(func() {
		if choice.SearchDropdown(c, label, o.Zones, zone, o.Search).Changed() {
			root.SetChanged()
			kit.Log("ZoneSelect", "picked")
		}
	})
}
