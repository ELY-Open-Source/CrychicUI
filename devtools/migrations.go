package devtools

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Migration is one schema change; At is when applied.
type Migration struct {
	ID, Name, At string
	Applied      bool
}

// MigrationsOptions word the panel; Applied and Pending take counts.
type MigrationsOptions struct{ Title, Applied, Pending, Next string }

// Migrations shows applied and pending schema changes.
func Migrations(c *ui.Context, migrations []Migration, o MigrationsOptions) ui.Element {
	def(&o.Title, "Migrations")
	if len(migrations) == 0 {
		panic(fmt.Sprintf("devtools: Migrations %q needs migrations", o.Title))
	}
	ids := make([]string, len(migrations))
	applied := 0
	for i, m := range migrations {
		if m.Name == "" || m.Applied && m.At == "" {
			panic(fmt.Sprintf("devtools: Migrations %q needs a name, and a time once applied", m.ID))
		}
		if m.Applied && i > applied {
			panic(fmt.Sprintf("devtools: Migrations %q is applied after a pending one", m.ID))
		}
		if m.Applied {
			applied++
		}
		ids[i] = m.ID
	}
	unique("Migrations", ids)
	def(&o.Applied, "%d applied")
	def(&o.Pending, "%d pending")
	def(&o.Next, "Next")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "layers", o.Title, func() {
			kit.Chip(c, fmt.Sprintf(o.Applied, applied), t.Success)
			if applied < len(migrations) {
				kit.Chip(c, fmt.Sprintf(o.Pending, len(migrations)-applied), t.Warning)
			}
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(3))
		list.Children(func() {
			for i, m := range migrations {
				last := i == len(migrations)-1
				c.Key(m.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(3))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Width(t.Space(5)).Shrink(0).AlignItems(ui.Center).Children(func() {
						c.Key(m.Applied)
						if m.Applied {
							kit.Arrive(glyph(c, "circle-check", 18).TextColor(t.Success), ui.Pose{Scale: 0.5, Fade: 1},
								time.Duration(i)*60*time.Millisecond)
						} else {
							glyph(c, "circle-dashed", 18).TextColor(t.Warning)
						}
						if !last {
							line := t.Border
							if m.Applied && migrations[i+1].Applied {
								line = t.Success.Alpha(0.5)
							}
							ui.Box(c).Width(1).Grow(1).MinHeight(t.Space(3)).Background(line)
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Padding(0, 0, t.Space(3), 0).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, m.Name).FontSize(t.Rem(0.82)).Bold().Grow(1).MinWidth(0).SingleLine()
							if !m.Applied && i == applied {
								kit.Chip(c, o.Next, t.Warning)
							}
						})
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							kit.Muted(c, m.ID).Font("monospace").FontSize(t.Rem(0.68)).MinWidth(0).SingleLine()
							if m.Applied {
								dot(c)
								kit.Muted(c, m.At).FontSize(t.Rem(0.68)).Shrink(0)
							}
						})
					})
				})
			}
		})
	})
	return root
}
