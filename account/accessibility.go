package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// A11y holds accessibility choices; TextSize runs 0 to 3.
type A11y struct {
	ReduceMotion, HighContrast bool
	TextSize                   int
}

// AccessibilityOptions word the card; Size takes the step number.
type AccessibilityOptions struct{ Title, Motion, MotionNote, Contrast, ContrastNote, Text, Size, Glyph, Sample string }

// Accessibility sets motion, contrast, and text size; cue Toggle.
func Accessibility(c *ui.Context, value *A11y, o AccessibilityOptions) ui.Element {
	if value == nil || value.TextSize < 0 || value.TextSize > 3 {
		panic(fmt.Sprintf("account: Accessibility %q needs a value with text size 0..3", o.Title))
	}
	def(&o.Title, "Accessibility")
	def(&o.Motion, "Reduce motion")
	def(&o.MotionNote, "Fewer animations across the app.")
	def(&o.Contrast, "High contrast")
	def(&o.ContrastNote, "Stronger borders and text.")
	def(&o.Text, "Text size")
	def(&o.Size, "Text size %d")
	def(&o.Glyph, "Aa")
	def(&o.Sample, "The quick brown fox jumps over the lazy dog.")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		head(c, "accessibility", o.Title, "")
		for k, row := range []struct {
			title, note, mark string
			on                *bool
		}{{o.Motion, o.MotionNote, "move", &value.ReduceMotion}, {o.Contrast, o.ContrastNote, "contrast", &value.HighContrast}} {
			c.Key(k)
			b := kit.Press(ui.ButtonBase(c).Label(row.title).FillWidth().Padding(t.Space(2), 0).Radius(t.Radius).
				Role(ui.RoleSwitch).Checked(*row.on))
			b.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
					ink := t.TextMuted
					if *row.on {
						ink = t.Accent
					}
					tile(c, row.mark, ink, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, row.title).FontSize(t.Rem(0.85)).Bold().Wrap().FillWidth()
						kit.Muted(c, row.note).FontSize(t.Rem(0.75)).Wrap().FillWidth()
					})
					track(c, *row.on, b.Hovered())
				})
			})
			if b.Clicked() {
				*row.on = !*row.on
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Accessibility", "toggled", "row", k, "on", *row.on)
			}
		}
		c.Key("size")
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3), 0, 0, 0).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Text(c, o.Text).FontSize(t.Rem(0.85)).Bold()
			ui.Row(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleRadioGroup).Children(func() {
				for i := range 4 {
					on := value.TextSize == i
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Size, i+1)).Grow(1).Basis(0).Height(t.Space(11)).
						Radius(t.Radius).Checked(on))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
					b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
					b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(on)), ui.Snappy))
					b.Children(func() { ui.Text(c, o.Glyph).FontSize(t.Rem(0.75 + 0.2*float32(i))).Bold().Shrink(0) })
					if b.Clicked() && !on {
						value.TextSize = i
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("Accessibility", "size", "step", i)
					}
				}
			})
			sample := ui.Text(c, o.Sample).Wrap().FillWidth()
			sample.FontSize(sample.AnimateTo("size", t.Rem(0.8+0.15*float32(value.TextSize)), ui.Smooth))
		})
	})
	return root
}
