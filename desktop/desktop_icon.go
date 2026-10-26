package desktop

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DesktopIcon selects on click; opening bumps opened; cue Open.
func DesktopIcon(c *ui.Context, name, mark string, selected *bool, opened *int) ui.Element {
	if strings.TrimSpace(name) == "" || mark == "" || selected == nil || opened == nil || *opened < 0 {
		panic(fmt.Sprintf("desktop: DesktopIcon %q needs an icon, a selection and an open count of 0 or more", name))
	}
	t := c.Theme()
	b := ui.ButtonBase(c).Label(name).Width(t.Space(22)).Padding(t.Space(1.5)).Radius(t.Radius).AlignItems(ui.Center)
	kit.Arrive(b, ui.Pose{Scale: 0.85, Fade: 1}, 0)
	on := *selected
	b.Background(b.AnimateColor("bg", t.Border.Alpha(0.4*unit(b.Hovered() && !on)), ui.Snappy))
	b.Children(func() {
		ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			ic := tile(c, mark, t.Accent, t.Space(12)).Radius(t.Radius * 1.5)
			ic.Border(2, ic.AnimateColor("ring", t.Accent.Alpha(unit(on)), ui.Snappy))
			label := ui.Text(c, name).FontSize(t.Rem(0.74)).TextAlign(ui.Center).MaxWidth(t.Space(19)).MaxLines(2).Padding(0, t.Space(1)).Radius(t.Radius * 0.5)
			label.Background(label.AnimateColor("pill", t.Accent.Alpha(unit(on)), ui.Snappy)).TextColor(label.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
		})
	})
	open := b.DoubleClicked() || (on && b.Shortcut(0, ui.KeyEnter))
	switch {
	case open && *opened < math.MaxInt:
		*selected = true
		*opened++
		b.SetChanged()
		sound.Play(sound.Open)
		kit.Log("DesktopIcon", "opened")
	case b.Clicked() && !on:
		*selected = true
		b.SetChanged()
		sound.Play(sound.Select)
		kit.Log("DesktopIcon", "selected")
	}
	return b
}
