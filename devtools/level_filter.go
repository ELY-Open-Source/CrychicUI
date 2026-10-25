package devtools

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LevelFilterOptions name Debug to Error.
type LevelFilterOptions struct{ Levels [4]string }

// LevelFilter toggles Debug to Error chips with counts; cue Toggle.
func LevelFilter(c *ui.Context, on *[4]bool, counts [4]int, o LevelFilterOptions) ui.Element {
	if on == nil {
		panic("devtools: LevelFilter needs state")
	}
	for i, n := range counts {
		if n < 0 {
			panic(fmt.Sprintf("devtools: LevelFilter count %d for level %d is negative", n, i))
		}
	}
	for i, name := range [...]string{"Debug", "Info", "Warn", "Error"} {
		def(&o.Levels[i], name)
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		for i := range 4 {
			ink := levelInk(t, Level(i))
			sel := on[i]
			b := kit.Press(ui.ButtonBase(c).Label(o.Levels[i]).Padding(t.Space(1.25), t.Space(2.5)).Radius(999).Checked(sel))
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.14), unit(sel)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(sel)), ui.Snappy))
			b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(ink, unit(sel)), ui.Snappy))
			kit.Stagger(c, root, b, i).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
					k := dot.AnimateTo("dot", 0.6+0.4*unit(sel), ui.Bouncy)
					dot.Scale(k, k)
					ui.Text(c, o.Levels[i]).FontSize(t.Rem(0.8)).Bold().Shrink(0)
					kit.Numeric(c, fmt.Sprint(counts[i])).FontSize(t.Rem(0.72)).TextColor(t.TextMuted).Shrink(0).
						ContentTransition(ui.ContentTransition{Numeric: true})
				})
			})
			if b.Clicked() {
				on[i] = !sel
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("LevelFilter", "toggled", "level", i, "on", !sel)
			}
		}
	})
	return root
}
