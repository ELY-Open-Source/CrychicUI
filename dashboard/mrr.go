package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// MRROptions are the MRR labels.
type MRROptions struct{ New, Expansion, Churn string }

// MRR shows MRR with new, expansion, and churn.
func MRR(c *ui.Context, current, new, expansion, churn float64, currency string, o MRROptions) ui.Element {
	if current < 0 || new < 0 || expansion < 0 || churn < 0 || !finite(current) || !finite(new) || !finite(expansion) || !finite(churn) || currency == "" {
		panic("dashboard: MRR has invalid values")
	}
	def(&o.New, "New")
	def(&o.Expansion, "Expansion")
	def(&o.Churn, "Churn")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		typography.Money(c, current, currency).Bold().FontSize(t.Rem(1.4))
		for _, row := range []struct {
			name  string
			value float64
			ink   ui.Color
		}{{o.New, new, t.Success}, {o.Expansion, expansion, t.Accent}, {o.Churn, churn, t.Danger}} {
			ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, row.name).Grow(1)
				typography.Money(c, row.value, currency).TextColor(row.ink)
			})
		}
	})
	return root
}
