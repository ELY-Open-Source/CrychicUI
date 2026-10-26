package geo

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Address is a street line and its city.
type Address struct{ Line, City string }

// AddressLookupOptions word the lookup; Join formats line and city.
type AddressLookupOptions struct{ Placeholder, Empty, Join string }

// AddressLookup suggests addresses as you type; picking cues Select.
func AddressLookup(c *ui.Context, query *string, addresses []Address, chosen *string, o AddressLookupOptions) ui.Element {
	def(&o.Placeholder, "Start typing an address")
	def(&o.Empty, "No matching addresses")
	def(&o.Join, "%s, %s")
	full := make([]string, len(addresses))
	for i, a := range addresses {
		if strings.TrimSpace(a.Line) == "" || strings.TrimSpace(a.City) == "" {
			panic(fmt.Sprintf("geo: AddressLookup address %q needs a line and a city", a.Line))
		}
		full[i] = fmt.Sprintf(o.Join, a.Line, a.City)
	}
	unique("AddressLookup", full)
	if query == nil || chosen == nil || (*chosen != "" && !slices.Contains(full, *chosen)) {
		panic(fmt.Sprintf("geo: AddressLookup %q needs a query and a listed chosen address", o.Placeholder))
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(2)).Padding(t.Space(3)).Role(ui.RoleGroup)
	return root.Children(func() {
		c.Key("field")
		frame := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background)
		frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
		frame.Children(func() {
			glyph(c, "house", 16).TextColor(t.TextMuted)
			in := ui.TextInputBase(c, query).Label(o.Placeholder).Placeholder(o.Placeholder).FontSize(t.Rem(0.86)).Grow(1).MinWidth(0)
			if in.Changed() {
				// Editing away from the pick clears it.
				if *chosen != "" && *query != *chosen {
					*chosen = ""
				}
				root.SetChanged()
				kit.Log("AddressLookup", "query", "runes", utf8.RuneCountInString(*query))
			}
		})
		q := strings.TrimSpace(*query)
		c.Key("below")
		below := ui.Column(c).FillWidth().Gap(t.Space(0.5))
		below.Children(func() {
			if *chosen != "" {
				kit.Arrive(ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Background(t.Success.Alpha(0.1)),
					ui.Pose{Y: -t.Space(1), Fade: 1}, 0).Children(func() {
					glyph(c, "circle-check", 16).TextColor(t.Success)
					ui.Text(c, *chosen).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).Wrap()
				})
				return
			}
			if q == "" {
				return
			}
			n := 0
			for i, a := range addresses {
				lo, hi := kit.IndexFold(a.Line, q)
				if lo < 0 {
					if l, _ := kit.IndexFold(a.City, q); l < 0 {
						continue
					}
				}
				if n == 5 {
					break
				}
				c.Key(strconv.Quote(full[i]))
				b := kit.Press(ui.ButtonBase(c).Label(full[i]).FillWidth().Padding(t.Space(2)).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
				kit.Stagger(c, below, b, n).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						glyph(c, "map-pin", 16).TextColor(t.Accent)
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
							ui.RichText(c, kit.Marked(t, a.Line, lo, hi)...).FontSize(t.Rem(0.84)).FillWidth().MinWidth(0).SingleLine()
							kit.Muted(c, a.City).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
						})
					})
				})
				n++
				if b.Clicked() {
					*chosen, *query = full[i], full[i]
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("AddressLookup", "chose", "index", i)
				}
			}
			if n == 0 {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(2)).Children(func() {
					glyph(c, "search-x", 16).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
				})
			}
		})
	})
}
