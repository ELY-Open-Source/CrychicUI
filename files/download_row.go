package files

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DownloadRowOptions word the download.
type DownloadRowOptions struct {
	Open, Of, Done string
	// Left words the time to go, "%s" taking FormatLeft.
	Left       string
	FormatLeft func(seconds float64) string
}

// leftText words seconds to wait: 45s, 12m 4s, 9d.
func leftText(secs float64) string {
	s := math.Floor(secs)
	switch {
	case s < 60:
		return fmt.Sprintf("%.0fs", s)
	case s < 3600:
		return fmt.Sprintf("%.0fm %.0fs", math.Floor(s/60), math.Mod(s, 60))
	case s < 48*3600:
		return fmt.Sprintf("%.0fh %.0fm", math.Floor(s/3600), math.Floor(math.Mod(s, 3600)/60))
	}
	return fmt.Sprintf("%.0fd", math.Floor(s/86400))
}

// DownloadRow shows progress, speed and time left; cue Success.
func DownloadRow(c *ui.Context, item FileItem, progress float64, speed int64, opened *bool, o DownloadRowOptions) (root, open ui.Element) {
	validFile("DownloadRow", item)
	if speed < 0 || opened == nil || !(progress >= 0 && progress <= 1) || math.IsInf(progress, 0) {
		panic(fmt.Sprintf("files: DownloadRow %q needs progress in [0,1], a speed and bound state", item.Name))
	}
	def(&o.Open, "Open")
	def(&o.Of, "of")
	def(&o.Done, "Done")
	def(&o.Left, "%s left")
	if o.FormatLeft == nil {
		o.FormatLeft = leftText
	}
	t := c.Theme()
	done := progress >= 1
	got := uploaded(UploadItem{Size: item.Size, Progress: progress})
	root = kit.Card(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(2.5), t.Space(3)).Role(ui.RoleGroup).Label(item.Name)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Row().AlignItems(ui.Center).Children(func() {
		fileTile(c, item, t.Space(9))
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
			ui.Text(c, item.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().SingleLine()
			ink := t.Accent
			if done {
				ink = t.Success
			}
			progressBar(c, progress, ink)
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				if done {
					ui.Icon(c, icon.Must("circle-check")).Size(12, 12).TextColor(t.Success)
					ui.Text(c, o.Done).FontSize(t.Rem(0.75)).TextColor(t.Success)
					kit.Muted(c, "·").FontSize(t.Rem(0.75))
					sizeText(c, item.Size).FontSize(t.Rem(0.75))
					return
				}
				sizeText(c, got).FontSize(t.Rem(0.75))
				kit.Muted(c, o.Of).FontSize(t.Rem(0.75))
				sizeText(c, item.Size).FontSize(t.Rem(0.75))
				kit.Muted(c, "·").FontSize(t.Rem(0.75))
				ui.Row(c).Children(func() {
					sizeText(c, speed).FontSize(t.Rem(0.75))
					kit.Muted(c, "/s").FontSize(t.Rem(0.75))
				})
				if speed > 0 {
					secs := float64(item.Size-got) / float64(speed)
					kit.Muted(c, "·").FontSize(t.Rem(0.75))
					kit.Muted(c, fmt.Sprintf(o.Left, o.FormatLeft(secs))).FontSize(t.Rem(0.75)).Shrink(0)
				}
			})
		})
		if done {
			c.Key("open")
			open = kit.Press(ui.ButtonBase(c).Label(o.Open).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)).Transition(kit.Pop)
			open.Background(open.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(open.Hovered())), ui.Snappy)).TextColor(t.AccentText).
				Children(func() { ui.Text(c, o.Open).Bold().Shrink(0) })
		}
	})
	if done && open.Clicked() && !*opened {
		*opened = true
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("DownloadRow", "opened")
	}
	return root, open
}
