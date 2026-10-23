package code

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TerminalOptions tune a Terminal.
type TerminalOptions struct {
	// Height of the view; zero gives the default.
	Height float32
	Prompt string
}

// Terminal shows lines of output and a blinking prompt.
func Terminal(c *ui.Context, label string, lines []string, o TerminalOptions) ui.Element {
	if label == "" {
		panic("code: Terminal needs a label")
	}
	def(&o.Prompt, "$")
	t := c.Theme()
	size("Terminal", &o.Height, t.Space(60))
	bg, fg := t.Inverse, t.InverseText
	if bg.A == 0 {
		bg = t.Text
	}
	if fg.A == 0 {
		fg = t.Background
	}
	root := ui.Column(c).FillWidth().Height(o.Height).Radius(t.Radius).Clip().Background(bg).TextColor(fg).
		Role(ui.RoleGroup).Label(label)
	state := ui.Local(root, "list", func() ui.ListState { return ui.ListState{FollowEnd: true, FollowAnimation: &ui.Smooth} })
	seen := ui.Local(root, "seen", func() int { return len(lines) })
	from := *seen
	*seen = len(lines)
	root.Children(func() {
		ui.List(c, state, len(lines)+1, func(i int) {
			row := ui.Row(c).FillWidth().Padding(0, t.Space(3)).Gap(t.Space(1))
			if *ui.Local(row, "fresh", func() bool { return i >= from }) {
				kit.Arrive(row, ui.Pose{Fade: 1}, 0)
			}
			row.Children(func() {
				if i < len(lines) {
					ui.Text(c, lines[i]).Font("monospace").FontSize(t.Rem(0.85)).Wrap()
					return
				}
				ui.Text(c, o.Prompt).Font("monospace").FontSize(t.Rem(0.85))
				caret := ui.Box(c).Width(t.Space(2)).Height(t.Rem(1)).Background(fg)
				ph := float64(caret.Loop("blink", time.Second, ui.Linear))
				caret.Opacity(unit(ph < 0.6))
			})
		}).Grow(1).Justify(ui.End).Padding(t.Space(2), 0)
	})
	if len(lines) > from {
		kit.Log("Terminal", "appended", "count", len(lines)-from)
	}
	return root
}
