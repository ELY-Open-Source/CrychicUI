package display

import (
	"math"
	"strings"
	"unicode"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// finite reports whether v is finite.
func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// unit converts on to an animation target.
func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// formatNumber writes a finite value in the viewer's region.
func formatNumber(c *ui.Context, value float64) string {
	if !finite(value) {
		panic("display: number must be finite")
	}
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), value, i18n.NumberOptions{})
}

// formatMoney writes an amount with two decimals.
func formatMoney(c *ui.Context, value float64) string {
	if !finite(value) {
		panic("display: money must be finite")
	}
	o := i18n.NumberOptions{MinFractionDigits: 2, MaxFractionDigits: 2}
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), value, o)
}

// initials returns up to two uppercase name initials.
func initials(name string) string {
	words := strings.Fields(name)
	if len(words) == 0 {
		panic("display: initials need a name")
	}
	out := make([]rune, 0, 2)
	for _, word := range words {
		for _, r := range word {
			out = append(out, unicode.ToUpper(r))
			break
		}
		if len(out) == 2 {
			break
		}
	}
	return string(out)
}

// nameColor gives a name one deterministic theme color.
func nameColor(c *ui.Context, name string) ui.Color {
	t := c.Theme()
	palette := [...]ui.Color{t.Accent, t.Success, t.Warning}
	var hash uint32
	for _, r := range name {
		hash = hash*31 + uint32(r)
	}
	return palette[hash%uint32(len(palette))]
}

// avatarCircle draws an image or initials in a circle.
func avatarCircle(c *ui.Context, name string, image *ui.Bitmap, size float32) ui.Element {
	if !finite(float64(size)) || size <= 0 {
		panic("display: avatar size must be finite and positive")
	}
	root := ui.Column(c).AlignSelf(ui.Start).Size(size, size).Center().Radius(999).Clip()
	if image != nil {
		root.Children(func() {
			c.Key("image")
			img := ui.Image(c, image).Fill().Fit(ui.Cover).ContentTransition(ui.ContentTransition{})
			kit.Arrive(img, ui.Pose{Fade: 1}, 0)
		})
		return kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	}
	color := nameColor(c, name)
	root.Background(color.Alpha(0.16)).TextColor(color).Bold().Children(func() {
		c.Key("initials")
		kit.Arrive(ui.Text(c, initials(name)), ui.Pose{Fade: 1}, 0)
	})
	return kit.Arrive(root, ui.Pose{Fade: 1}, 0)
}
