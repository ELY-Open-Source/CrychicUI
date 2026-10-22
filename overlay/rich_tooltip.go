package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RichTooltip is a fading tooltip with title and text.
func RichTooltip(c *ui.Context, anchor ui.Element, title, text string) ui.Element {
	if title == "" || text == "" {
		panic(fmt.Sprintf("overlay: RichTooltip %q needs a title and text", title))
	}
	t := c.Theme()
	anchor.Description(text)
	ink, paper := t.Inverse, t.InverseText
	if ink == (ui.Color{}) {
		ink, paper = t.Text, t.Background
	}
	return ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.Margin(0, 0, t.Space(1.5), 0).Width(t.Space(60)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(ink).TextColor(paper)
		kit.Arrive(tip, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
			ui.Column(c).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, title).Bold()
				ui.Text(c, text).FontSize(t.Rem(0.85)).TextColor(paper.Alpha(0.8))
			})
		})
	})
}
