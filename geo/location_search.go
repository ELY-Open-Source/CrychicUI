package geo

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LocationSearchOptions word the search; From adds distances.
type LocationSearchOptions struct {
	Placeholder, Empty string
	From               *Point
	Imperial           bool
	Units              Units
}

// LocationSearch filters places by name or address; cue Select.
func LocationSearch(c *ui.Context, query *string, places []Place, picked *string, o LocationSearchOptions) ui.Element {
	def(&o.Placeholder, "Search places")
	def(&o.Empty, "No places found")
	o.Units.fill()
	names := make([]string, len(places))
	for i, p := range places {
		checkPlace("LocationSearch", p)
		names[i] = p.Name
	}
	unique("LocationSearch", names)
	if query == nil || picked == nil {
		panic("geo: LocationSearch needs a query and a picked place")
	}
	if o.From != nil {
		valid("LocationSearch", o.Placeholder, *o.From)
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup)
	return root.Children(func() {
		c.Key("field")
		frame := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(999).Background(t.Background)
		frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
		frame.Children(func() {
			glyph(c, "search", 16).TextColor(t.TextMuted)
			in := ui.TextInputBase(c, query).Label(o.Placeholder).Placeholder(o.Placeholder).FontSize(t.Rem(0.86)).Grow(1).MinWidth(0)
			if in.Changed() {
				root.SetChanged()
				kit.Log("LocationSearch", "query", "runes", utf8.RuneCountInString(*query))
			}
		})
		q := strings.TrimSpace(*query)
		var hits []Place
		for _, p := range places {
			if lo, _ := kit.IndexFold(p.Name, q); lo >= 0 {
				hits = append(hits, p)
			} else if lo, _ := kit.IndexFold(p.Address, q); lo >= 0 {
				hits = append(hits, p)
			}
		}
		c.Key("list")
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleList)
		list.Children(func() {
			if len(hits) == 0 {
				ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(5), 0).Children(func() {
					kit.Arrive(glyph(c, "search-x", 24).TextColor(t.TextMuted), ui.Pose{Scale: 0.7, Fade: 1}, 0)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.84)).FillWidth().TextAlign(ui.Center).Wrap()
				})
				return
			}
			for i, p := range hits {
				def(&p.Icon, "map-pin")
				c.Key(strconv.Quote(p.Name))
				on := p.Name == *picked
				b := kit.Press(ui.ButtonBase(c).Label(p.Name).FillWidth().Padding(t.Space(2)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						tile(c, p.Icon, t.Accent, t.Space(8))
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
							lo, hi := kit.IndexFold(p.Name, q)
							ui.RichText(c, kit.Marked(t, p.Name, lo, hi)...).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).SingleLine()
							if p.Address != "" {
								lo, hi := kit.IndexFold(p.Address, q)
								ui.RichText(c, kit.Marked(t, p.Address, lo, hi)...).FontSize(t.Rem(0.74)).TextColor(t.TextMuted).
									FillWidth().MinWidth(0).SingleLine()
							}
						})
						if o.From != nil {
							kit.Muted(c, span(haversine(*o.From, p.At), o.Imperial, o.Units)).FontSize(t.Rem(0.74)).Shrink(0).MaxWidth(t.Space(20)).SingleLine()
						}
					})
				})
				if b.Clicked() && !on {
					*picked = p.Name
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("LocationSearch", "picked")
				}
			}
		})
	})
}
