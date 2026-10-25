package learn

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HintRevealOptions word the card and its button.
type HintRevealOptions struct{ Title, Ask, Note string }

// HintReveal keeps a hint hidden until asked; cue Open.
func HintReveal(c *ui.Context, hint string, shown *bool, o HintRevealOptions) ui.Element {
	if shown == nil || strings.TrimSpace(hint) == "" {
		panic(fmt.Sprintf("learn: HintReveal %q needs a hint and a state", hint))
	}
	def(&o.Title, "Hint")
	def(&o.Ask, "Show hint")
	def(&o.Note, "Try once more before you look.")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			bulb := tile(c, "lightbulb", t.Warning, t.Space(9))
			bulb.Background(bulb.AnimateColor("bg", t.Warning.Alpha(0.12).Mix(t.Warning.Alpha(0.28), unit(*shown)), ui.Snappy))
			ui.Text(c, o.Title).FontSize(t.Rem(0.9)).Bold().Grow(1).MinWidth(0).SingleLine()
		})
		c.Key(*shown)
		if *shown {
			body := ui.Text(c, hint).FontSize(t.Rem(0.88)).FillWidth().Wrap()
			kit.Arrive(body, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			return
		}
		// Hidden lines stand in for the hint.
		lines := ui.Column(c).FillWidth().Gap(t.Space(2))
		lines.Children(func() {
			for i, w := range []float32{1, 0.82, 0.55} {
				row := ui.Row(c).FillWidth()
				kit.Stagger(c, lines, row, i).Children(func() {
					ui.Box(c).Grow(w).Basis(0).Height(t.Space(2.5)).Radius(999).Background(t.SurfaceHover.Mix(t.Border, 0.3))
					ui.Box(c).Grow(1 - w).Basis(0)
				})
			}
		})
		kit.Muted(c, o.Note).FontSize(t.Rem(0.76)).FillWidth().Wrap()
		b := kit.Press(ui.ButtonBase(c).Label(o.Ask).Padding(t.Space(2), t.Space(4)).Radius(999).Border(1, t.Border).Shrink(0))
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "eye", 14)
				ui.Text(c, o.Ask).FontSize(t.Rem(0.82)).Bold().Shrink(0)
			})
		})
		if b.Clicked() && !*shown {
			*shown = true
			root.SetChanged()
			sound.Play(sound.Open)
			kit.Log("HintReveal", "revealed")
		}
	})
	return root
}
