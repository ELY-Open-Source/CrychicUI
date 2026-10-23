package chat

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Channel is one row of a ChannelList.
type Channel struct {
	Name   string
	Unread int
}

// ChannelList shows "#" channels with unread dots.
func ChannelList(c *ui.Context, label string, channels []Channel, selected *int) ui.Element {
	if label == "" || len(channels) == 0 || selected == nil || *selected < -1 || *selected >= len(channels) {
		panic(fmt.Sprintf("chat: ChannelList %q needs a label, channels and a choice or -1", label))
	}
	for _, ch := range channels {
		if ch.Name == "" || ch.Unread < 0 {
			panic(fmt.Sprintf("chat: ChannelList %q channels need a name and a count", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, ch := range channels {
			c.Key(i)
			on := i == *selected
			b := kit.Press(ui.ButtonBase(c).Label(ch.Name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius).Checked(on))
			kit.Arrive(b, ui.Pose{Y: t.Space(1.5), Fade: 1}, time.Duration(i)*30*time.Millisecond)
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Accent.Alpha(0.1)
			}
			fg := t.Text
			if on {
				fg = t.Accent
			}
			b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).TextColor(b.AnimateColor(fmt.Sprint("fg.", i), fg, ui.Snappy)).
				Children(func() {
					ui.Icon(c, icon.Must("hash")).Size(16, 16).TextColor(t.TextMuted)
					ui.Text(c, ch.Name).Grow(1).SingleLine()
					if ch.Unread > 0 {
						dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(t.Accent)
						kit.Arrive(dot, ui.Pose{Scale: 0.5, Fade: 1}, 0)
					}
				})
			if b.Clicked() && !on {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ChannelList", "selected", "index", i)
			}
		}
	})
}
