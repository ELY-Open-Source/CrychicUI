package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Insight is one generated insight card.
type Insight struct {
	Title, Body string
}

// Insights shows generated insight cards.
func Insights(c *ui.Context, cards []Insight) ui.Element {
	if len(cards) == 0 {
		panic("dashboard: Insights needs cards")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, card := range cards {
			if card.Title == "" || card.Body == "" {
				panic("dashboard: Insights has an invalid card")
			}
			c.Key(i)
			panel := ui.Column(c).Gap(t.Space(1)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
			kit.Arrive(panel, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			panel.Children(func() {
				ui.Text(c, card.Title).Bold()
				ui.Text(c, card.Body).FontSize(t.Rem(0.85)).Wrap()
			})
		}
	})
	return root
}
