package chat

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Conversation is one row of a ConversationList.
type Conversation struct {
	Name, Last string
	At         time.Time
	Unread     int
	Avatar     *ui.Bitmap
}

// ConversationList shows avatar, name, last line, time and unread.
func ConversationList(c *ui.Context, label string, items []Conversation, selected *int) ui.Element {
	if label == "" || len(items) == 0 || selected == nil || *selected < -1 || *selected >= len(items) {
		panic(fmt.Sprintf("chat: ConversationList %q needs a label, items and a choice or -1", label))
	}
	for _, k := range items {
		if k.Name == "" || k.At.IsZero() || k.Unread < 0 {
			panic(fmt.Sprintf("chat: ConversationList %q items need a name, a time and a count", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, k := range items {
			c.Key(i)
			on := i == *selected
			b := kit.Press(ui.ButtonBase(c).Label(k.Name).Row().FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).
				Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Checked(on))
			kit.Arrive(b, ui.Pose{Y: t.Space(1.5), Fade: 1}, time.Duration(i)*30*time.Millisecond)
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Accent.Alpha(0.1)
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
				display.Avatar(c, k.Name, k.Avatar)
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, k.Name).Bold().SingleLine()
					if k.Last != "" {
						kit.Muted(c, k.Last).SingleLine().FontSize(t.Rem(0.85))
					}
				})
				ui.Column(c).AlignItems(ui.End).Gap(t.Space(1)).Shrink(0).Children(func() {
					kit.Muted(c, when(c, k.At)).FontSize(t.Rem(0.75))
					if k.Unread > 0 {
						pill := ui.Box(c).Padding(0, t.Space(1.5)).MinHeight(t.Space(4.5)).Center().Radius(999).
							Background(t.Accent).TextColor(t.AccentText)
						kit.Arrive(pill, ui.Pose{Scale: 0.6, Fade: 1}, 0).Children(func() {
							kit.Numeric(c, strconv.Itoa(k.Unread)).FontSize(t.Rem(0.75)).Bold()
						})
					}
				})
			})
			if b.Clicked() && !on {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ConversationList", "selected", "index", i)
			}
		}
	})
}
