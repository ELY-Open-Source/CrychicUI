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
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleStatus).Label(from + " → " + to)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	glyph := "arrow-right"
	if c.Direction() == ui.RightToLeft {
		glyph = "arrow-left"
	}
	side := func(name string) {
		ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Width(t.Space(24)).Shrink(0).Children(func() {
			display.Avatar(c, name, nil)
			ui.Text(c, name).FontSize(t.Rem(0.85)).Bold().SingleLine()
		})
	}
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			side(from)
			track := ui.Box(c).Grow(1).Height(t.Space(6)).Center()
			x := track.Keyframes("crychic.pass", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			track.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				y := r.Y + r.H/2
				p.Fill(ui.Rect{X: r.X, Y: y - 1, W: r.W, H: 2}, t.Border, 1)
				p.Fill(ui.Rect{X: r.X, Y: y - 1, W: r.W * x, H: 2}, t.Accent, 1)
			}))
			track.Children(func() {
				ui.Box(c).Size(t.Space(6), t.Space(6)).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText).
					Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16) })
			})
			side(to)
		})
		ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Surface).Children(func() {
			if note != "" {
				ui.Text(c, note).FontSize(t.Rem(0.9)).Wrap()
			}
			kit.Muted(c, i18n.FormatDate(c.Region().FormatLocale(c.Locale()), at, i18n.DateOptions{Time: i18n.ShortStyle})).FontSize(t.Rem(0.75))
		})
	})
}
