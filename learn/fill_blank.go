package learn

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FillBlankOptions word the gap, the check and the verdicts.
type FillBlankOptions struct{ Gap, Check, Right, Wrong string }

// FillBlank checks a typed missing word; cues Success, Error.
func FillBlank(c *ui.Context, before, after, answer string, text *string, checked *bool, o FillBlankOptions) ui.Element {
	if text == nil || checked == nil || strings.TrimSpace(answer) == "" || strings.TrimSpace(before+after) == "" {
		panic(fmt.Sprintf("learn: FillBlank %q needs an answer, words around it and a text", answer))
	}
	def(&o.Gap, "Missing word")
	def(&o.Check, "Check")
	def(&o.Right, "Correct")
	def(&o.Wrong, "The answer is %s")
	t := c.Theme()
	good := strings.EqualFold(strings.TrimSpace(*text), strings.TrimSpace(answer))
	ink := t.Accent
	if *checked {
		ink = t.Danger
		if good {
			ink = t.Success
		}
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	check := func() {
		good := strings.EqualFold(strings.TrimSpace(*text), strings.TrimSpace(answer))
		*checked = true
		root.SetChanged()
		if good {
			sound.Play(sound.Success)
		} else {
			sound.Play(sound.Error)
		}
		kit.Log("FillBlank", "checked", "right", good, "runes", utf8.RuneCountInString(*text))
	}
	root.Children(func() {
		line := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1)).AlignItems(ui.Center)
		line.Children(func() {
			for _, w := range strings.Fields(before) {
				ui.Text(c, w).FontSize(t.Rem(0.95)).Shrink(0)
			}
			c.Key("gap")
			gap := ui.Box(c).Width(t.Space(30)).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
				BorderWidth(0, 0, t.Space(0.5), 0)
			gap.BorderColor(gap.AnimateColor("line", t.Border.Mix(ink, unit(*checked || gap.FocusWithin())), ui.Snappy))
			gap.Background(gap.AnimateColor("bg", t.SurfaceHover.Alpha(0.6).Mix(ink.Alpha(0.1), unit(*checked)), ui.Snappy))
			gap.Children(func() {
				in := ui.TextInputBase(c, text).Label(o.Gap).Placeholder("…").FontSize(t.Rem(0.95)).FillWidth().MinWidth(0).
					TextColor(t.Text.Mix(ink, unit(*checked)))
				if in.Changed() {
					*checked = false
					root.SetChanged()
				}
				if in.Submitted() && strings.TrimSpace(*text) != "" && !*checked {
					check()
				}
			})
			for _, w := range strings.Fields(after) {
				ui.Text(c, w).FontSize(t.Rem(0.95)).Shrink(0)
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			c.Key(*checked)
			if *checked {
				mark, word := "circle-x", fmt.Sprintf(o.Wrong, answer)
				if good {
					mark, word = "circle-check", o.Right
				}
				note := ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(ink)
				kit.Arrive(note, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
					glyph(c, mark, 16)
					ui.Text(c, word).FontSize(t.Rem(0.82)).Bold().Grow(1).MinWidth(0).Wrap()
				})
			} else {
				ui.Box(c).Grow(1)
			}
			c.Key("check")
			b := kit.Press(ui.ButtonBase(c).Label(o.Check).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0).
				Disabled(*checked || strings.TrimSpace(*text) == ""))
			b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
			b.Children(func() { ui.Text(c, o.Check).FontSize(t.Rem(0.82)).Bold().Shrink(0) })
			if b.Clicked() && !*checked && strings.TrimSpace(*text) != "" {
				check()
			}
		})
	})
	return root
}
