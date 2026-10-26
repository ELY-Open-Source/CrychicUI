// Package marketing holds landing page components.
package marketing

import (
	"fmt"
	"math"
	"net/mail"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

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
	return kit.Arrive(surface(c), ui.Pose{Y: c.Theme().Space(2), Fade: 1}, 0)
}

// surface is a card that leaves arrival to its caller.
func surface(c *ui.Context) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Radius(t.Radius*1.5).Border(1, t.Border).Background(t.Background).Shadow(0, 1, 3, 0, t.Text.Alpha(0.06))
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
			panic(fmt.Sprintf("marketing: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// number groups an integer for the current locale.
func number(c *ui.Context, n any) string {
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), n, i18n.NumberOptions{Integer: true})
}

// avatar is a photo or initials circle units wide.
func avatar(c *ui.Context, name string, photo *ui.Bitmap, units float32) ui.Element {
	t := c.Theme()
	return display.Avatar(c, name, photo).Size(t.Space(units), t.Space(units)).FontSize(t.Rem(units * 0.08)).Shrink(0).AlignSelf(ui.Center)
}

// action is a pill button; room caps its label.
func action(c *ui.Context, label, mark string, solid bool, room float32) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Shrink(0).Padding(t.Space(2.25), t.Space(4)).Radius(999))
	hover := unit(b.Hovered())
	if solid {
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.Text, 0.12*hover), ui.Snappy)).TextColor(t.AccentText)
	} else {
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, hover), ui.Snappy)).Border(1, t.Border)
	}
	shift := b.AnimateTo("nudge", hover, ui.Snappy)
	return b.Children(func() {
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, label).FontSize(t.Rem(0.84)).Bold().Shrink(0).MaxWidth(room).SingleLine()
			if mark != "" {
				glyph(c, mark, 16).Translate(t.Space(0.75)*shift, 0)
			}
		})
	})
}

func validEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

// emailField is a bordered email input; it returns the input.
func emailField(c *ui.Context, label string, email *string) ui.Element {
	return fieldIn(c, ui.Row(c).Grow(1).Basis(c.Theme().Space(40)).MinWidth(0), label, "mail", email)
}

// fieldIn fills frame with a mark and an input.
func fieldIn(c *ui.Context, frame ui.Element, label, mark string, value *string) ui.Element {
	t := c.Theme()
	var in ui.Element
	frame.Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(999).Background(t.Background)
	frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
	frame.Children(func() {
		glyph(c, mark, 16).TextColor(t.TextMuted)
		in = ui.TextInputBase(c, value).Label(label).Placeholder(label).FontSize(t.Rem(0.86)).Grow(1).MinWidth(0)
	})
	return in
}

// heading is a centered section title and line.
func heading(c *ui.Context, title, text string) {
	t := c.Theme()
	if title == "" && text == "" {
		return
	}
	ui.Column(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
		if title != "" {
			kit.Arrive(ui.Text(c, title).FontSize(t.Rem(1.3)).Bold().FillWidth().TextAlign(ui.Center).Wrap(), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
		}
		if text != "" {
			kit.Muted(c, text).FontSize(t.Rem(0.86)).FillWidth().MaxWidth(t.Space(110)).TextAlign(ui.Center).Wrap()
		}
	})
}

// decimal formats v with digits places for the current locale.
func decimal(c *ui.Context, v float64, digits int) string {
	if digits == 0 {
		return number(c, v)
	}
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), v, i18n.NumberOptions{MinFractionDigits: digits, MaxFractionDigits: digits})
}

// price shows whole amounts bare and others with cents.
func price(c *ui.Context, v float64) string {
	if v == math.Trunc(v) {
		return decimal(c, v, 0)
	}
	return decimal(c, v, 2)
}

// Plan is a pricing plan; Yearly is billed yearly.
type Plan struct {
	Name, Text      string
	Monthly, Yearly float64
	Features        []string
	Highlight       bool
	Badge           string
}

// PlanWords word a plan card; Choose takes the name.
type PlanWords struct{ Currency, Per, Free, Choose, Start, Current string }

func (w *PlanWords) fill() {
	def(&w.Currency, "$")
	def(&w.Per, "/ month")
	def(&w.Free, "Free")
	def(&w.Choose, "Choose %s")
	def(&w.Start, "Get started")
	def(&w.Current, "Current plan")
}

func validPlan(component string, p Plan) {
	if strings.TrimSpace(p.Name) == "" || !(p.Monthly >= 0 && p.Monthly <= 1e9) || !(p.Yearly >= 0 && p.Yearly <= 1e9) {
		panic(fmt.Sprintf("marketing: %s plan %q needs a name and prices of 0 to 1e9", component, p.Name))
	}
	unique(component, p.Features)
}

// planCard draws one plan; current disables its button.
func planCard(c *ui.Context, p Plan, yearly, current bool, w PlanWords) (root, button ui.Element) {
	t := c.Theme()
	root = card(c).Gap(t.Space(3)).Padding(t.Space(5))
	if p.Highlight {
		root.Border(t.Space(0.5), t.Accent).Background(t.Background.Mix(t.Accent, 0.04))
	}
	root.Children(func() {
		if p.Badge != "" {
			ui.Row(c).Children(func() {
				kit.Arrive(kit.Chip(c, p.Badge, t.Accent), ui.Pose{Scale: 0.8, Fade: 1}, 200*time.Millisecond)
			})
		}
		ui.Text(c, p.Name).FontSize(t.Rem(1)).Bold().FillWidth().MinWidth(0).SingleLine()
		if p.Text != "" {
			kit.Muted(c, p.Text).FontSize(t.Rem(0.8)).FillWidth().Wrap()
		}
		amount := p.Monthly
		if yearly {
			amount = p.Yearly
		}
		ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
			if amount == 0 {
				ui.Text(c, w.Free).FontSize(t.Rem(1.8)).Bold().Shrink(0)
				return
			}
			kit.Numeric(c, w.Currency+price(c, amount)).FontSize(t.Rem(1.8)).Bold().MinWidth(0).SingleLine().
				ContentTransition(ui.ContentTransition{Numeric: true})
			kit.Muted(c, w.Per).FontSize(t.Rem(0.78)).Shrink(0).Padding(0, 0, t.Space(1), 0)
		})
		label, shown := fmt.Sprintf(w.Choose, p.Name), w.Start
		if current {
			label, shown = w.Current, w.Current
		}
		c.Key("choose")
		button = kit.Press(ui.ButtonBase(c).Label(label).FillWidth().Padding(t.Space(2.25), t.Space(3)).Radius(999).Disabled(current))
		if p.Highlight {
			button.Background(button.AnimateColor("bg", t.Accent.Mix(t.Text, 0.12*unit(button.Hovered())), ui.Snappy)).TextColor(t.AccentText)
		} else {
			button.Background(button.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(button.Hovered())), ui.Snappy)).Border(1, t.Border)
		}
		button.Children(func() {
			ui.Text(c, shown).FontSize(t.Rem(0.84)).Bold().MinWidth(0).SingleLine()
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(1.75)).Padding(t.Space(1), 0, 0, 0)
		list.Children(func() {
			for i, f := range p.Features {
				r := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
				kit.Stagger(c, list, r, i).Children(func() {
					glyph(c, "check", 14).TextColor(t.Success)
					ui.Text(c, f).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
				})
			}
		})
	})
	return root, button
}

// arrow is a round outlined icon button.
func arrow(c *ui.Context, label, mark string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(9), t.Space(9)).Shrink(0).Center().Radius(999).Border(1, t.Border))
	b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
	return b.Children(func() { glyph(c, mark, 16) })
}

// circle is a round path centered on cx, cy.
func circle(cx, cy, r float32) *ui.Path {
	const k = 0.5522847498
	var p ui.Path
	p.MoveTo(cx, cy-r)
	p.CubeTo(cx+r*k, cy-r, cx+r, cy-r*k, cx+r, cy)
	p.CubeTo(cx+r, cy+r*k, cx+r*k, cy+r, cx, cy+r)
	p.CubeTo(cx-r*k, cy+r, cx-r, cy+r*k, cx-r, cy)
	p.CubeTo(cx-r, cy-r*k, cx-r*k, cy-r, cx, cy-r)
	return p.Close()
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

// countdown is whole seconds left; it wakes each second.
func countdown(c *ui.Context, until time.Time) int64 {
	now := c.Now()
	// Unix seconds: Sub saturates past 292 years.
	secs := until.Unix() - now.Unix()
	if until.Nanosecond() > now.Nanosecond() {
		secs++
	}
	if secs <= 0 {
		return 0
	}
	frac := time.Duration((until.Nanosecond() - now.Nanosecond() + 1e9) % 1e9)
	if frac == 0 {
		frac = time.Second
	}
	c.After(frac + time.Millisecond)
	return secs
}

// counted picks one or many by n.
func counted(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf(one, n)
	}
	return fmt.Sprintf(many, n)
}

// dot joins two meta texts with a drawn dot.
func dot(c *ui.Context) {
	t := c.Theme()
	ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.TextMuted.Alpha(0.6))
}

// track is a switch drawing with a sliding knob.
func track(c *ui.Context, on, hover bool) {
	t := c.Theme()
	b := ui.Box(c).Size(t.Space(9), t.Space(5)).Shrink(0).Radius(999)
	b.Background(b.AnimateColor("track", t.SurfaceHover.Mix(t.Border, 0.6*unit(hover)).Mix(t.Accent, unit(on)), ui.Snappy))
	x := b.AnimateTo("knob", unit(on), ui.Snappy)
	pad := t.Space(0.5)
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		d := r.H - 2*pad
		p.Fill(ui.Rect{X: r.X + pad + x*(r.W-r.H), Y: r.Y + pad, W: d, H: d}, t.Background, d/2)
	})
}

// Post is a blog post; zero Minutes hides read time.
type Post struct {
	Title, Excerpt, Tag, Author string
	Photo, Cover                *ui.Bitmap
	Date                        time.Time
	Minutes                     int
}

// PostWords word a post; Date is a layout, Read minutes.
type PostWords struct{ Date, Read string }

func (w *PostWords) fill() {
	def(&w.Date, "Jan 2, 2006")
	def(&w.Read, "%d min read")
}

func validPost(component string, p Post) {
	if strings.TrimSpace(p.Title) == "" || strings.TrimSpace(p.Author) == "" || p.Date.IsZero() || p.Minutes < 0 || p.Minutes > 600 {
		panic(fmt.Sprintf("marketing: %s post %q needs an author, a date and 0 to 600 minutes", component, p.Title))
	}
}

// meta is read time then the date, which truncates.
func meta(c *ui.Context, p Post, w PostWords) {
	t := c.Theme()
	row := ui.Row(c).FillWidth().MinWidth(0).Gap(t.Space(1.5)).AlignItems(ui.Center)
	seen := row.LayoutBounds().W
	kit.Settle(c, row, []float32{seen}, func() []float32 { return []float32{row.LayoutBounds().W} })
	// Read time gets at most half the line.
	room := t.Space(20)
	if seen > 0 {
		room = max(t.Space(6), seen/2)
	}
	row.Children(func() {
		if p.Minutes > 0 {
			kit.Muted(c, fmt.Sprintf(w.Read, p.Minutes)).FontSize(t.Rem(0.72)).Shrink(0).MaxWidth(room).SingleLine()
			dot(c)
		}
		kit.Muted(c, p.Date.Format(w.Date)).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
	})
}
