package geo

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CoordinatesOptions word the fields; Sides are N, S, E, W.
type CoordinatesOptions struct {
	Lat, Lng, Invalid, Copy, Copied string
	Sides                           [4]string
}

// dms words v as degrees, minutes and seconds.
func dms(v float64, pos, neg string) string {
	side := pos
	if v < 0 {
		side = neg
	}
	secs := int(math.Round(math.Abs(v) * 3600))
	return fmt.Sprintf("%d°%02d′%02d″%s", secs/3600, secs/60%60, secs%60, side)
}

// Coordinates edits a point as numbers; Copy cues Pop.
func Coordinates(c *ui.Context, at *Point, o CoordinatesOptions) (root, copy ui.Element) {
	def(&o.Lat, "Latitude")
	def(&o.Lng, "Longitude")
	def(&o.Invalid, "Enter a number from %s to %s")
	def(&o.Copy, "Copy")
	def(&o.Copied, "Copied")
	for i, d := range [4]string{"N", "S", "E", "W"} {
		def(&o.Sides[i], d)
	}
	if at == nil {
		panic("geo: Coordinates needs a point")
	}
	valid("Coordinates", o.Lat, *at)
	t := c.Theme()
	text := func(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
	root = card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	texts := ui.Local(root, "texts", func() [2]string { return [2]string{text(at.Lat), text(at.Lng)} })
	bad := ui.Local(root, "bad", func() [2]bool { return [2]bool{} })
	// A point set from outside rewrites both fields.
	last := ui.Local(root, "last", func() Point { return *at })
	if *at != *last {
		*texts, *bad, *last = [2]string{text(at.Lat), text(at.Lng)}, [2]bool{}, *at
	}
	root.Children(func() {
		c.Key("fields")
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).Children(func() {
			for i, label := range [2]string{o.Lat, o.Lng} {
				limit := [2]float64{85, 180}[i]
				c.Key(i)
				ui.Column(c).Grow(1).Basis(t.Space(40)).MinWidth(0).Gap(t.Space(1)).Children(func() {
					kit.Muted(c, label).FontSize(t.Rem(0.72)).Bold().FillWidth().MinWidth(0).SingleLine()
					frame := ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background)
					frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())).Mix(t.Danger, unit(bad[i])), ui.Snappy))
					frame.Children(func() {
						in := ui.TextInputBase(c, &texts[i]).Label(label).Font("monospace").FontSize(t.Rem(0.84)).Grow(1).MinWidth(0)
						if in.Changed() {
							v, err := strconv.ParseFloat(strings.TrimSpace(texts[i]), 64)
							bad[i] = err != nil || !(v >= -limit && v <= limit)
							if !bad[i] {
								if i == 0 {
									at.Lat = v
								} else {
									at.Lng = v
								}
								*last = *at
							}
							root.SetChanged()
							kit.Log("Coordinates", "edited", "field", i, "runes", utf8.RuneCountInString(texts[i]), "valid", !bad[i])
						}
					})
					if bad[i] {
						c.Key("bad")
						ui.Text(c, fmt.Sprintf(o.Invalid, text(-limit), text(limit))).FontSize(t.Rem(0.7)).TextColor(t.Danger).FillWidth().Wrap()
					}
				})
			}
		})
		c.Key("foot")
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			glyph(c, "map-pin", 14).TextColor(t.Accent)
			ui.Text(c, dms(at.Lat, o.Sides[0], o.Sides[1])).Font("monospace").FontSize(t.Rem(0.76)).Shrink(0).SingleLine()
			ui.Text(c, dms(at.Lng, o.Sides[2], o.Sides[3])).Font("monospace").FontSize(t.Rem(0.76)).Grow(1).MinWidth(0).SingleLine()
			value := text(at.Lat) + "," + text(at.Lng)
			copy = kit.Press(ui.ButtonBase(c).Label(o.Copy).Row().Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).
				Padding(t.Space(1), t.Space(2.5)).Radius(999).Border(1, t.Border))
			copy.Background(copy.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(copy.Hovered())), ui.Snappy))
			done := ui.Local(copy, "until", func() time.Time { return time.Time{} })
			held := ui.Local(copy, "held", func() string { return "" })
			if copy.Clicked() {
				c.WriteClipboard(value)
				*done, *held = c.Now().Add(2*time.Second), value
				sound.Play(sound.Pop)
				kit.Log("Coordinates", "copied")
			}
			copied := c.Now().Before(*done) && *held == value
			if copied {
				c.After(done.Sub(c.Now()))
			}
			copy.Children(func() {
				shown, mark := o.Copy, "copy"
				if copied {
					shown, mark = o.Copied, "check"
				}
				glyph(c, mark, 12).ContentTransition(ui.ContentTransition{})
				ui.Text(c, shown).FontSize(t.Rem(0.74)).Bold().Shrink(0).MaxWidth(t.Space(24)).SingleLine().ContentTransition(ui.ContentTransition{})
			})
		})
	})
	return root, copy
}
