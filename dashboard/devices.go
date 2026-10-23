package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DeviceShare is one device type share.
type DeviceShare struct {
	Kind  string
	Share float64
}

// Devices shows desktop, mobile, and tablet shares.
func Devices(c *ui.Context, shares []DeviceShare) ui.Element {
	if len(shares) == 0 {
		panic("dashboard: Devices needs shares")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, share := range shares {
			if share.Kind == "" || share.Share < 0 || share.Share > 1 || !finite(share.Share) {
				panic("dashboard: Devices has an invalid share")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, share.Kind).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%.0f%%", share.Share*100)).TextColor(t.TextMuted)
				})
				progress(c, share.Share, t.Accent)
			})
		}
	})
	return root
}
