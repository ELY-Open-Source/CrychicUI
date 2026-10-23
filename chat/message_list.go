package chat

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MessageListOptions tune a MessageList.
type MessageListOptions struct {
	// Height of the list; zero gives the default.
	Height float32
}

// MessageList groups bubbles by sender; new ones slide up.
func MessageList(c *ui.Context, label string, list []Message, o MessageListOptions) ui.Element {
	if label == "" {
		panic("chat: MessageList needs a label")
	}
	messages("MessageList", list)
	t := c.Theme()
	size("MessageList", &o.Height, t.Space(90))
	root := ui.Column(c).FillWidth().Height(o.Height).Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup).Label(label)
	state := ui.Local(root, "list", func() ui.ListState { return ui.ListState{FollowEnd: true, FollowAnimation: &ui.Smooth} })
	seen := ui.Local(root, "seen", func() int { return len(list) })
	from := *seen
	*seen = len(list)
	root.Children(func() {
		ui.List(c, state, len(list), func(i int) {
			m := list[i]
			first := i == 0 || !grouped(list[i-1], m)
			last := i == len(list)-1 || !grouped(m, list[i+1])
			row := ui.Column(c).FillWidth().Padding(t.Space(0.5), t.Space(3)).Gap(t.Space(1))
			if first {
				row.Margin(t.Space(2), 0, 0, 0)
			}
			if *ui.Local(row, "fresh", func() bool { return i >= from }) {
				kit.Arrive(row, ui.Pose{Y: t.Space(4), Fade: 1}, 0)
			}
			row.Children(func() {
				if first && !m.Mine {
					kit.Muted(c, m.Sender).FontSize(t.Rem(0.8)).Bold()
				}
				shell(c, m.Mine).Children(func() { ui.Text(c, m.Text).Wrap() })
				if last {
					stamp := kit.Muted(c, when(c, m.At)).FontSize(t.Rem(0.75))
					if m.Mine {
						stamp.AlignSelf(ui.End)
					}
				}
			})
		}).Grow(1).Justify(ui.End).Padding(0, 0, t.Space(2), 0)
	})
	if len(list) > from {
		theirs := 0
		for _, m := range list[from:] {
			if !m.Mine {
				theirs++
			}
		}
		kit.Log("MessageList", "appended", "count", len(list)-from, "theirs", theirs)
	}
	return root
}

// grouped reports b following a closely from one sender.
func grouped(a, b Message) bool {
	return a.Sender == b.Sender && a.Mine == b.Mine && !b.At.Before(a.At) && b.At.Sub(a.At) < 5*time.Minute
}
