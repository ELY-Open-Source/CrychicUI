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

// CallBannerOptions name the join and leave buttons.
type CallBannerOptions struct {
	Join, Leave string
}

// CallBanner is a "call in progress" bar with join.
func CallBanner(c *ui.Context, title string, joined *bool, o CallBannerOptions) ui.Element {
	if title == "" || joined == nil {
		panic(fmt.Sprintf("chat: CallBanner %q needs a title and joined", title))
	}
	def(&o.Join, "Join")
	def(&o.Leave, "Leave")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
		Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Role(ui.RoleGroup).Label(title)
	var btn ui.Element
	root.Children(func() {
		beacon := ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).
			Background(t.Success.Alpha(0.15)).TextColor(t.Success)
		ph := float64(beacon.Loop("pulse", 1500*time.Millisecond, ui.Linear))
		beacon.Opacity(0.6 + 0.4*float32(math.Sin(2*math.Pi*ph)))
		beacon.Children(func() { ui.Icon(c, icon.Must("phone")).Size(16, 16) })
		ui.Text(c, title).Bold().Grow(1).SingleLine()
		label := o.Join
		if *joined {
			label = o.Leave
		}
		btn = kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(3)).Radius(999).Checked(*joined))
		bg, fg := t.Accent.Mix(t.AccentHover, unit(btn.Hovered())), t.AccentText
		if *joined {
			bg, fg = t.SurfaceHover, t.Text
		}
		btn.Background(btn.AnimateColor("bg", bg, ui.Snappy)).TextColor(btn.AnimateColor("fg", fg, ui.Snappy)).
			Children(func() { ui.Text(c, label).Bold().SingleLine().ContentTransition(ui.ContentTransition{}) })
	})
	if btn.Clicked() {
		*joined = !*joined
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("CallBanner", "joined", *joined)
	}
	return root
}
