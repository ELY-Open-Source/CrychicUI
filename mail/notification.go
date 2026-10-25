package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Notice is one notification row; Detail and Icon are optional.
type Notice struct {
	ID, Text, Detail, Time, Icon string
}

// Notification opens a notification row; cue Tap.
func Notification(c *ui.Context, notice Notice, opened *bool) ui.Element {
	if notice.ID == "" || notice.Text == "" || notice.Time == "" || opened == nil {
		panic(fmt.Sprintf("mail: Notification %q needs text, time, and state", notice.ID))
	}
	if notice.Icon == "" {
		notice.Icon = "bell"
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(notice.Text).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(2.5), t.Space(3)).
		Radius(t.Radius).Border(1, t.Border))
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(root.Hovered())).Mix(t.Accent.Alpha(0.05), unit(!*opened)), ui.Snappy))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start).Children(func() {
			tile(c, notice.Icon, t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				text := ui.Text(c, notice.Text).FontSize(t.Rem(0.85)).Wrap().FillWidth()
				if !*opened {
					text.Bold()
				}
				if notice.Detail != "" {
					kit.Muted(c, notice.Detail).FontSize(t.Rem(0.75)).Wrap().FillWidth()
				}
				kit.Muted(c, notice.Time).FontSize(t.Rem(0.7))
			})
			ui.Box(c).Shrink(0).Padding(t.Space(1), 0, 0, 0).Children(func() {
				dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999)
				dot.Background(dot.AnimateColor("dot", t.Accent.Alpha(unit(!*opened)), ui.Snappy))
			})
		})
	})
	if root.Clicked() {
		*opened = true
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("Notification", "opened")
	}
	return root
}
