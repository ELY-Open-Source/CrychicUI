package account

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Device is one remembered device.
type Device struct {
	ID, Name, Last string
}

// DevicesListOptions are the device labels.
type DevicesListOptions struct{ Remove, Empty string }

// DevicesList removes devices; cue Delete.
func DevicesList(c *ui.Context, devices *[]Device, o DevicesListOptions) ui.Element {
	if devices == nil {
		panic("account: DevicesList needs devices")
	}
	def(&o.Remove, "Remove")
	def(&o.Empty, "No devices")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	remove := -1
	root.Children(func() {
		if len(*devices) == 0 {
			kit.Muted(c, o.Empty)
		}
		for i, device := range *devices {
			if device.ID == "" || device.Name == "" || device.Last == "" {
				panic("account: invalid device")
			}
			c.Key(device.ID)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, device.Name)
					kit.Muted(c, device.Last)
				})
				b := button(c, device.ID+" "+o.Remove)
				b.Children(func() { ui.Text(c, o.Remove) })
				if b.Clicked() {
					remove = i
				}
			})
		}
	})
	if remove >= 0 {
		*devices = slices.Delete(*devices, remove, remove+1)
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("DevicesList", "removed", "count", len(*devices))
	}
	return root
}
