package chat

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChatEmpty is an empty chat with a wave.
func ChatEmpty(c *ui.Context, title, text, action string) (root, act ui.Element) {
	if title == "" || action == "" {
		panic(fmt.Sprintf("chat: ChatEmpty %q needs a title and an action", title))
	}
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(8)).Role(ui.RoleGroup).Label(title)
	root.Children(func() {
		hand := ui.Icon(c, icon.Must("hand")).Size(t.Space(10), t.Space(10)).TextColor(t.Accent)
		ph := float64(hand.Loop("wave", 1600*time.Millisecond, ui.Linear))
		hand.Rotate(14 * float32(math.Sin(2*math.Pi*ph)))
		ui.Text(c, title).Bold().FontSize(t.Rem(1.1))
		if text != "" {
			kit.Muted(c, text).FontSize(t.Rem(0.9))
		}
		act = kit.Press(ui.ButtonBase(c).Label(action).Padding(t.Space(1.5), t.Space(4)).Radius(999))
		act.Background(act.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(act.Hovered())), ui.Snappy)).
			TextColor(t.AccentText).Margin(t.Space(2), 0, 0, 0).
			Children(func() { ui.Text(c, action).Bold().SingleLine() })
	})
	if act.Clicked() {
		sound.Play(sound.Tap)
		kit.Log("ChatEmpty", "acted")
	}
	return root, act
}
