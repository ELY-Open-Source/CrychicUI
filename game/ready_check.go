package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReadyCheckOptions word the check; Count takes ready and total.
type ReadyCheckOptions struct{ Title, Ready, NotReady, Count, All string }

// ReadyCheck toggles your ready mark beside the others'; cue Select.
func ReadyCheck(c *ui.Context, others []Seat, ready *bool, o ReadyCheckOptions) ui.Element {
	def(&o.Title, "Ready check")
	if ready == nil || len(others) > 15 {
		panic(fmt.Sprintf("game: ReadyCheck %q needs your state and at most 15 others", o.Title))
	}
	names := make([]string, len(others))
	count := 0
	for i, p := range others {
		names[i] = p.Name
		if p.Ready {
			count++
		}
	}
	if len(names) > 0 {
		unique("ReadyCheck", names)
	}
	if *ready {
		count++
	}
	total := len(others) + 1
	def(&o.Ready, "Ready")
	def(&o.NotReady, "Not ready")
	def(&o.Count, "%d of %d ready")
	def(&o.All, "Everyone is ready")
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			glyph(c, "users", 16).TextColor(t.Accent)
			ui.Text(c, o.Title).FontSize(t.Rem(0.9)).Bold().Grow(1).MinWidth(0).SingleLine()
			c.Key(count == total)
			if count == total {
				kit.Arrive(kit.Chip(c, o.All, t.Success), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			} else {
				kit.Chip(c, fmt.Sprintf(o.Count, count, total), t.Warning)
			}
		})
		meter(c, float32(count)/float32(total), t.Warning.Mix(t.Success, unit(count == total))).Height(t.Space(1.5))
		seats := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2))
		seats.Children(func() {
			for i, p := range others {
				c.Key(p.Name)
				chip := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2.5), t.Space(1), t.Space(1)).Radius(999)
				chip.Background(chip.AnimateColor("bg", t.SurfaceHover.Mix(t.Success.Alpha(0.14), unit(p.Ready)), ui.Snappy))
				kit.Stagger(c, seats, chip, i).Children(func() {
					display.Avatar(c, p.Name, nil).Size(t.Space(6), t.Space(6)).FontSize(t.Rem(0.55))
					ui.Text(c, p.Name).FontSize(t.Rem(0.75)).Bold().Shrink(0).MaxWidth(t.Space(28)).SingleLine()
					c.Key(p.Ready)
					mark, ink := "hourglass", t.TextMuted
					if p.Ready {
						mark, ink = "circle-check", t.Success
					}
					kit.Arrive(glyph(c, mark, 14).TextColor(ink), ui.Pose{Scale: 0.5, Fade: 1}, 0)
				})
			}
		})
		label := o.Ready
		if *ready {
			label = o.NotReady
		}
		b := kit.Press(ui.ButtonBase(c).Label(label).FillWidth().Padding(t.Space(2.5)).Radius(t.Radius).Checked(*ready))
		b.Background(b.AnimateColor("bg", t.Success.Mix(t.Success.Mix(t.Text, 0.15), unit(b.Hovered())).Mix(t.SurfaceHover, unit(*ready)), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.Background.Mix(t.Text, unit(*ready)), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				c.Key(*ready)
				mark := "circle-check"
				if *ready {
					mark = "x"
				}
				kit.Arrive(glyph(c, mark, 16), ui.Pose{Scale: 0.5, Fade: 1}, 0)
				ui.Text(c, label).FontSize(t.Rem(0.9)).Bold().Shrink(0)
			})
		})
		if b.Clicked() {
			*ready = !*ready
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("ReadyCheck", "ready", *ready)
		}
	})
	return root
}
