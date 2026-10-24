package kit

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// Stagger arrives a row, staggered while its list is young.
func Stagger(c *ui.Context, list, row ui.Element, i int) ui.Element {
	born := ui.Local(list, "crychic.born", func() time.Time { return c.Now() })
	delay := ui.Local(row, "crychic.delay", func() time.Duration {
		if c.Now().Sub(*born) < time.Second {
			return time.Duration(min(i, 8)) * 40 * time.Millisecond
		}
		return 0
	})
	return Arrive(row, ui.Pose{Y: c.Theme().Space(1), Fade: 1}, *delay)
}

// CardHead is a card's title row: icon, title, then trailing.
func CardHead(c *ui.Context, glyph, title string, trailing func()) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).
		BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
		ui.Text(c, title).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
		if trailing != nil {
			trailing()
		}
	})
}

// Chip is a small tinted pill; new text cross-fades.
func Chip(c *ui.Context, text string, ink ui.Color) ui.Element {
	t := c.Theme()
	pill := ui.Row(c).Shrink(0).Padding(t.Space(0.25), t.Space(2)).Radius(999).Background(ink.Alpha(0.12))
	return pill.Children(func() {
		ui.Text(c, text).FontSize(t.Rem(0.7)).Bold().TextColor(ink).Shrink(0).ContentTransition(ui.ContentTransition{})
	})
}
