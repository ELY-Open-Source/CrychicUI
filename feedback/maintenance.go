package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// MaintenanceOptions tune a Maintenance; empty copy reads English.
type MaintenanceOptions struct {
	Title, Text string
}

// Maintenance is a rocking tools icon and a message.
func Maintenance(c *ui.Context, o MaintenanceOptions) ui.Element {
	if o.Title == "" {
		o.Title = "Down for maintenance"
	}
	if o.Text == "" {
		o.Text = "We will be back shortly."
	}
	root, _ := state(c, icon.Must("wrench"), c.Theme().Warning, o.Title, o.Text, "", false, func(i ui.Element) {
		k := i.Loop("rock", 1400*time.Millisecond, ui.Bounce(ui.EaseInOut))
		i.Rotate(-14 + 28*k)
	})
	return root
}
