package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// AgentHandoff is an arrow between two agents with a note.
func AgentHandoff(c *ui.Context, from, to, note string, at time.Time) ui.Element {
	if from == "" || to == "" || at.IsZero() {
		panic(fmt.Sprintf("agent: AgentHandoff %q needs two agents and a time", from))
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.Center).Role(ui.RoleStatus)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	glyph := "arrow-right"
	if c.Direction() == ui.RightToLeft {
		glyph = "arrow-left"
	}
	return root.Children(func() {
		display.Avatar(c, from, nil)
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(t.TextMuted)
		display.Avatar(c, to, nil)
		ui.Column(c).Children(func() {
			if note != "" {
				ui.Text(c, note).FontSize(t.Rem(0.9))
			}
			kit.Muted(c, i18n.FormatDate(c.Region().FormatLocale(c.Locale()), at, i18n.DateOptions{Time: i18n.ShortStyle})).FontSize(t.Rem(0.75))
		})
	})
}
