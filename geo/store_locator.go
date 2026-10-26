package geo

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Store is a shop shown on the locator.
type Store struct {
	Name, Address string
	At            Point
	Open          bool
}

// StoreLocatorOptions word the list; From adds distances.
type StoreLocatorOptions struct {
	Label, Open, Closed string
	From                *Point
	Imperial            bool
	Units               Units
}

// StoreLocator lists stores over a centering map; cue Select.
func StoreLocator(c *ui.Context, stores []Store, picked *string, o StoreLocatorOptions) ui.Element {
	def(&o.Label, "Stores")
	def(&o.Open, "Open")
	def(&o.Closed, "Closed")
	o.Units.fill()
	if len(stores) == 0 || picked == nil {
		panic(fmt.Sprintf("geo: StoreLocator %q needs stores and a picked store", o.Label))
	}
	names := make([]string, len(stores))
	pts := make([]Point, len(stores))
	for i, s := range stores {
		valid("StoreLocator", s.Name, s.At)
		names[i], pts[i] = s.Name, s.At
	}
	unique("StoreLocator", names)
	if *picked != "" && !slices.Contains(names, *picked) {
		panic(fmt.Sprintf("geo: StoreLocator %q has no store %q", o.Label, *picked))
	}
	if o.From != nil {
		valid("StoreLocator", o.Label, *o.From)
		pts = append(pts, *o.From)
	}
	t := c.Theme()
	root := card(c).Clip().Role(ui.RoleGroup).Label(o.Label)
	return root.Children(func() {
		aim := middle(pts)
		if i := slices.Index(names, *picked); i >= 0 {
			aim = stores[i].At
		}
		c.Key("map")
		m := ui.Box(c).FillWidth().Height(t.Space(44))
		// shown is unwrapped, so the map pans the short way.
		shown := ui.Local(m, "lng", func() float64 { return aim.Lng })
		*shown += wrap(aim.Lng - *shown)
		lat, lng := m.AnimateTo("lat", float32(aim.Lat), ui.Smooth), m.AnimateTo("lng", float32(*shown), ui.Smooth)
		m.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			// One step out keeps every store in view.
			_, zoom := fit(pts, r)
			zoom--
			center := Point{float64(lat), wrap(float64(lng))}
			terrain(p, r, t, center, zoom)
			if o.From != nil {
				x, y := project(center, *o.From, zoom, r)
				var halo, dot ui.Path
				p.FillPath(halo.Circle(x, y, t.Space(3)), t.Accent.Alpha(0.2))
				p.FillPath(dot.Circle(x, y, t.Space(1.5)), t.Accent)
			}
			for _, s := range stores {
				x, y := project(center, s.At, zoom, r)
				if s.Name == *picked {
					continue
				}
				pin(p, x, y, t.Space(3), t.TextMuted, t.Background)
			}
			if i := slices.Index(names, *picked); i >= 0 {
				x, y := project(center, stores[i].At, zoom, r)
				pin(p, x, y, t.Space(4.5), t.Danger, t.Background)
			}
		}))
		c.Key("list")
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(2)).Role(ui.RoleList)
		list.Children(func() {
			for i, s := range stores {
				c.Key(strconv.Quote(s.Name))
				on := s.Name == *picked
				b := kit.Press(ui.ButtonBase(c).Label(s.Name).FillWidth().Padding(t.Space(2)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						tile(c, "store", t.Accent.Mix(t.Danger, unit(on)), t.Space(8))
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
							ui.Text(c, s.Name).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).SingleLine()
							if s.Address != "" {
								kit.Muted(c, s.Address).FontSize(t.Rem(0.74)).FillWidth().MinWidth(0).SingleLine()
							}
						})
						ui.Column(c).Shrink(0).Gap(t.Space(0.25)).AlignItems(ui.End).Children(func() {
							if o.From != nil {
								ui.Text(c, span(haversine(*o.From, s.At), o.Imperial, o.Units)).FontSize(t.Rem(0.74)).Bold().Shrink(0).MaxWidth(t.Space(24)).SingleLine()
							}
							state, ink := o.Closed, t.Danger
							if s.Open {
								state, ink = o.Open, t.Success
							}
							ui.Text(c, state).FontSize(t.Rem(0.7)).Bold().TextColor(ink).Shrink(0).MaxWidth(t.Space(24)).SingleLine()
						})
					})
				})
				if b.Clicked() && !on {
					*picked = s.Name
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("StoreLocator", "picked", "index", i)
				}
			}
		})
	})
}
