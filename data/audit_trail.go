package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Audit is who did what, and when.
type Audit struct{ Who, What, When string }

// AuditTrail lists changes by person; rows stagger in.
func AuditTrail(c *ui.Context, label string, audits []Audit) ui.Element {
	if label == "" || len(audits) == 0 {
		panic(fmt.Sprintf("data: AuditTrail %q needs audits", label))
	}
	for _, a := range audits {
		if a.Who == "" || a.What == "" {
			panic(fmt.Sprintf("data: AuditTrail %q audits need who and what", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, a := range audits {
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				display.Avatar(c, a.Who, nil)
				ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, a.Who).Bold()
					ui.Text(c, a.What)
				})
				kit.Muted(c, a.When).Shrink(0)
			})
		}
	})
}
