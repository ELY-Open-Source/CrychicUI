package color

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ContrastOptions tune a Contrast.
type ContrastOptions struct {
	// Sample is shown text; AA and AAA name marks.
	Sample, AA, AAA string
	// Pass and Fail tell assistive technology each result.
	Pass, Fail string
}

// Contrast shows fg on bg, its ratio and marks.
func Contrast(c *ui.Context, fg, bg ui.Color, o ContrastOptions) ui.Element {
	def(&o.Sample, "Aa")
	def(&o.AA, "AA")
	def(&o.AAA, "AAA")
	def(&o.Pass, "passes")
	def(&o.Fail, "fails")
	t := c.Theme()
	r := ratio(fg.Over(bg), bg)
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		card := ui.Box(c).Size(t.Space(20), t.Space(14)).Center().Radius(t.Radius).Border(1, t.Border)
		card.Background(card.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
			s := ui.Text(c, o.Sample).FontSize(t.Rem(1.75)).Bold()
			s.TextColor(s.AnimateColor("fg", fg, ui.Snappy))
		})
		ui.Column(c).Gap(t.Space(1.5)).Children(func() {
			kit.Numeric(c, fmt.Sprintf("%.2f:1", r)).FontSize(t.Rem(1.25)).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
			ui.Row(c).Gap(t.Space(2)).Children(func() {
				mark(c, o.AA, r >= 4.5, o)
				mark(c, o.AAA, r >= 7, o)
			})
		})
	})
}

// mark is a pass or fail pill.
func mark(c *ui.Context, name string, pass bool, o ContrastOptions) {
	t := c.Theme()
	tint, glyph, word := t.Danger, "x", o.Fail
	if pass {
		tint, glyph, word = t.Success, "check", o.Pass
	}
	m := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).TextColor(t.AccentText).
		Role(ui.RoleText).Label(name + " " + word)
	m.Background(m.AnimateColor("tint", tint, ui.Snappy)).Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
		ui.Text(c, name).FontSize(t.Rem(0.85)).Bold()
	})
}
