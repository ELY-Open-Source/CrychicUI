package agent

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ToolCall shows a tool name, arguments and status.
func ToolCall(c *ui.Context, name, args string, running bool) ui.Element {
	if name == "" {
		panic("agent: ToolCall needs a name")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).
		Radius(t.Radius).Background(t.Surface).Role(ui.RoleGroup).Label(name)
	return root.Children(func() {
		ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Children(func() {
			if running {
				ic := ui.Icon(c, icon.Must("rotate-cw")).Size(16, 16).TextColor(t.Accent)
				ic.Rotate(360 * float32(ic.Loop("spin", time.Second, ui.Linear)))
			} else {
				done := ui.Icon(c, icon.Must("circle-check")).Size(16, 16).TextColor(t.Success)
				kit.Arrive(done, ui.Pose{Scale: 0.5, Fade: 1}, 0)
			}
		})
		ui.Icon(c, icon.Must("terminal")).Size(16, 16).TextColor(t.TextMuted)
		ui.Text(c, name).Bold().SingleLine()
		if args != "" {
			kit.Muted(c, args).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).SingleLine()
		}
	})
}
