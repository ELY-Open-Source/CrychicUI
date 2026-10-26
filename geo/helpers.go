// Package geo holds map, place and location components.
package geo

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// metersPerDegree is one degree of latitude.
const metersPerDegree = 111320.0

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

// glyph is a Lucide icon size DIP square.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	return kit.Arrive(ui.Column(c).FillWidth().Radius(t.Radius*1.5).Border(1, t.Border).Background(t.Background).
		Shadow(0, 1, 3, 0, t.Text.Alpha(0.06)), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// tile is a tinted icon square.
func tile(c *ui.Context, mark string, ink ui.Color, size float32) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink)
	return box.Children(func() { glyph(c, mark, size*0.5) })
}

// unique panics when items repeat or are blank.
func unique(component string, items []string) {
	seen := map[string]bool{}
	for _, s := range items {
		if strings.TrimSpace(s) == "" || seen[s] {
			panic(fmt.Sprintf("geo: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// valid panics unless at is a real place.
func valid(component, label string, at Point) {
	if !(at.Lat >= -85 && at.Lat <= 85 && at.Lng >= -180 && at.Lng <= 180) {
		panic(fmt.Sprintf("geo: %s %q needs latitude ±85 and longitude ±180, got %v", component, label, at))
	}
}

// scale is pixels per degree of latitude at zoom.
func scale(zoom float64) float64 { return 256 * math.Exp2(zoom) / 360 }

// wrap brings a longitude into -180 to 180.
func wrap(lng float64) float64 {
	l := math.Mod(lng+180, 360)
	if l < 0 {
		l += 360
	}
	return l - 180
}

// project maps at into r, centered on center at zoom.
func project(center, at Point, zoom float64, r ui.Rect) (float32, float32) {
	k := scale(zoom)
	cos := math.Cos(center.Lat * math.Pi / 180)
	return r.X + r.W/2 + float32(wrap(at.Lng-center.Lng)*k*cos), r.Y + r.H/2 - float32((at.Lat-center.Lat)*k)
}

// terrain paints land, parks, a river and streets.
func terrain(p *ui.Painter, r ui.Rect, t *ui.Theme, center Point, zoom float64) {
	p.Fill(r, t.SurfaceHover.Mix(t.Success, 0.08), 0)
	k := scale(zoom)
	cos := math.Cos(center.Lat * math.Pi / 180)
	lat0, lat1 := center.Lat-float64(r.H/2)/k, center.Lat+float64(r.H/2)/k
	lng0, lng1 := center.Lng-float64(r.W/2)/(k*cos), center.Lng+float64(r.W/2)/(k*cos)
	const block = 0.02
	// Parks and rivers below a few pixels are not drawn.
	for i := math.Floor(lat0 / block); k*block >= 8 && i*block < lat1; i++ {
		for j := math.Floor(lng0 / block); j*block < lng1; j++ {
			// 18,000 blocks circle the globe.
			if (int64(i)*73856093^int64(math.Mod(math.Mod(j, 18000)+18000, 18000))*19349663)%7 != 0 {
				continue
			}
			x0, y0 := project(center, Point{(i + 0.85) * block, (j + 0.15) * block}, zoom, r)
			x1, y1 := project(center, Point{(i + 0.15) * block, (j + 0.85) * block}, zoom, r)
			p.Fill(ui.Rect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}, t.Success.Alpha(0.16), float32(k*block*0.06))
		}
	}
	for band := math.Floor(lat0/0.05) - 1; k*0.05 >= 30 && band*0.05 < lat1+0.05; band++ {
		// Whole waves per turn keep the river seamless.
		var river ui.Path
		for n := 0; n <= 48; n++ {
			lng := lng0 + (lng1-lng0)*float64(n)/48
			x, y := project(center, Point{band*0.05 + 0.017 + 0.004*math.Sin(lng*math.Pi/180*10313), lng}, zoom, r)
			if n == 0 {
				river.MoveTo(x, y)
			} else {
				river.LineTo(x, y)
			}
		}
		p.StrokePath(&river, max(4, float32(k*0.0012)), t.Accent.Mix(t.Background, 0.7))
	}
	road := func(step float64, width float32, ink ui.Color) {
		if k*step < 18 {
			return
		}
		for lat := math.Floor(lat0/step) * step; lat < lat1; lat += step {
			_, y := project(center, Point{lat, center.Lng}, zoom, r)
			p.Fill(ui.Rect{X: r.X, Y: y - width/2, W: r.W, H: width}, ink, 0)
		}
		for lng := math.Floor(lng0/step) * step; lng < lng1; lng += step {
			x, _ := project(center, Point{center.Lat, lng}, zoom, r)
			p.Fill(ui.Rect{X: x - width/2, Y: r.Y, W: width, H: r.H}, ink, 0)
		}
	}
	road(0.004, 2, t.Background)
	road(0.02, 5, t.Background.Mix(t.Warning, 0.25))
}

// pin paints a map pin tipped at x, y.
func pin(p *ui.Painter, x, y, size float32, ink, hole ui.Color) {
	head := y - size
	var tip ui.Path
	tip.MoveTo(x-size*0.5, head+size*0.3).LineTo(x+size*0.5, head+size*0.3).LineTo(x, y).Close()
	p.FillPath(&tip, ink)
	var disc ui.Path
	p.FillPath(disc.Circle(x, head, size*0.6), ink)
	var dot ui.Path
	p.FillPath(dot.Circle(x, head, size*0.24), hole)
}

// Units word distances; each format takes a number string.
type Units struct{ Meters, Kilometers, Feet, Miles string }

func (u *Units) fill() {
	def(&u.Meters, "%s m")
	def(&u.Kilometers, "%s km")
	def(&u.Feet, "%s ft")
	def(&u.Miles, "%s mi")
}

// meters panics unless m is a distance on Earth.
func meters(component, label string, m float64) {
	if !(m >= 0 && m <= 4.1e7) {
		panic(fmt.Sprintf("geo: %s %q needs 0 to 41,000 km, got %v m", component, label, m))
	}
}

// span words m meters, metric or imperial.
func span(m float64, imperial bool, u Units) string {
	num := func(v float64, prec int) string { return strconv.FormatFloat(v, 'f', prec, 64) }
	if imperial {
		if ft := math.Round(m * 3.28084); ft < 1000 {
			return fmt.Sprintf(u.Feet, num(ft, 0))
		}
		mi := m / 1609.344
		if math.Round(mi*10) < 100 {
			return fmt.Sprintf(u.Miles, num(mi, 1))
		}
		return fmt.Sprintf(u.Miles, num(mi, 0))
	}
	if r := math.Round(m); r < 1000 {
		return fmt.Sprintf(u.Meters, num(r, 0))
	}
	km := m / 1000
	if math.Round(km*10) < 100 {
		return fmt.Sprintf(u.Kilometers, num(km, 1))
	}
	return fmt.Sprintf(u.Kilometers, num(km, 0))
}

// duration words minutes, switching to hours past one.
func duration(minutes int, short, long string) string {
	if minutes < 60 {
		return fmt.Sprintf(short, minutes)
	}
	return fmt.Sprintf(long, minutes/60, minutes%60)
}

// haversine is the great-circle meters from a to b.
func haversine(a, b Point) float64 {
	rad := math.Pi / 180
	dLat, dLng := (b.Lat-a.Lat)*rad, (b.Lng-a.Lng)*rad
	h := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(a.Lat*rad)*math.Cos(b.Lat*rad)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * 6371000 * math.Asin(math.Min(1, math.Sqrt(h)))
}

// bounds are the corners of pts, unwrapped from the first.
func bounds(pts []Point) (lo, hi Point) {
	lo, hi = pts[0], pts[0]
	for _, q := range pts {
		lng := pts[0].Lng + wrap(q.Lng-pts[0].Lng)
		lo = Point{min(lo.Lat, q.Lat), min(lo.Lng, lng)}
		hi = Point{max(hi.Lat, q.Lat), max(hi.Lng, lng)}
	}
	return lo, hi
}

// middle is the wrapped center of pts.
func middle(pts []Point) Point {
	lo, hi := bounds(pts)
	return Point{(lo.Lat + hi.Lat) / 2, wrap((lo.Lng + hi.Lng) / 2)}
}

// fit is the center and zoom showing pts in r.
func fit(pts []Point, r ui.Rect) (Point, float64) {
	lo, hi := bounds(pts)
	mid := middle(pts)
	cos := math.Cos(mid.Lat * math.Pi / 180)
	dLat, dLng := max(hi.Lat-lo.Lat, 0.002), max((hi.Lng-lo.Lng)*cos, 0.002)
	k := min(float64(r.H)*0.7/dLat, float64(r.W)*0.8/dLng)
	return mid, math.Min(18, math.Log2(k*360/256))
}

// number groups n's digits for the locale.
func number(c *ui.Context, n any) string {
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), n, i18n.NumberOptions{Integer: true})
}

// star is a filled five-point star inside r.
func star(r ui.Rect) *ui.Path {
	var p ui.Path
	cx, cy := r.X+r.W/2, r.Y+r.H*0.54
	outer, inner := min(r.W, r.H)/2, min(r.W, r.H)/5
	for i := range 10 {
		rad := outer
		if i%2 == 1 {
			rad = inner
		}
		th := float64(i) * math.Pi / 5
		x, y := cx+rad*float32(math.Sin(th)), cy-rad*float32(math.Cos(th))
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	return p.Close()
}

// rating draws five stars filled to score of 5.
func rating(c *ui.Context, score float64, size float32) {
	t := c.Theme()
	ui.Row(c).Gap(size * 0.15).Shrink(0).Children(func() {
		for i := range 5 {
			fill := float32(max(0, min(1, score-float64(i))))
			ui.Box(c).Size(size, size).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillPath(star(r), t.TextMuted.Alpha(0.25))
				if fill > 0 {
					p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, 0, func() { p.FillPath(star(r), t.Warning) })
				}
			})
		}
	})
}
