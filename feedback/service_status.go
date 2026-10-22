package feedback

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ServiceState is operational, degraded or down.
type ServiceState uint8

// The three service states.
const (
	Operational ServiceState = iota
	Degraded
	Down
)

// Service is a name, state and uptime, 0 to 1.
type Service struct {
	Name   string
	State  ServiceState
	Uptime float64
}

// ServiceStatusOptions tune a ServiceStatus.
type ServiceStatusOptions struct {
	// States name the three states; English when zero.
	States [3]string
	// Uptime writes the uptime percent; "%.2f%%" when empty.
	Uptime string
}

// ServiceStatus is service rows with status dots and uptime.
func ServiceStatus(c *ui.Context, services []Service, o ServiceStatusOptions) ui.Element {
	names := make([]string, len(services))
	for i, s := range services {
		fraction("ServiceStatus", s.Uptime)
		if s.Name == "" || slices.Contains(names[:i], s.Name) || s.State > Down {
			panic("feedback: ServiceStatus needs distinct names and known states")
		}
		names[i] = s.Name
	}
	if len(services) == 0 {
		panic("feedback: ServiceStatus needs services")
	}
	if o.States == ([3]string{}) {
		o.States = [3]string{"Operational", "Degraded", "Down"}
	}
	if o.Uptime == "" {
		o.Uptime = "%.2f%%"
	}
	t := c.Theme()
	tints := [3]ui.Color{t.Success, t.Warning, t.Danger}
	root := kit.Card(c).FillWidth().Dividers(1, t.Border).Gap(1)
	return root.Children(func() {
		for i, s := range services {
			c.Key(s.Name)
			row := ui.Row(c).Gap(t.Space(3)).Padding(t.Space(2.5), t.Space(3)).Role(ui.RoleStatus).Label(s.Name)
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999)
				dot.Background(dot.AnimateColor("tint", tints[s.State], quick))
				ui.Text(c, s.Name).Grow(1)
				kit.Muted(c, o.States[s.State]).FontSize(t.Rem(0.85)).ContentTransition(ui.ContentTransition{})
				kit.Numeric(c, fmt.Sprintf(o.Uptime, s.Uptime*100)).FontSize(t.Rem(0.85)).MinWidth(t.Space(16)).TextAlign(ui.End)
			})
		}
	})
}
