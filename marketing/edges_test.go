package marketing

import (
	"fmt"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// focusHolds clicks name, shows optional copy, then presses Space.
func focusHolds(t *testing.T, what, name string, more *bool, hits *int, view func(c *ui.Context)) {
	t.Helper()
	tt := ui.NewTester(padded(view), 520, 700)
	settles(t, tt)
	click(t, tt, name)
	*more = true
	tt.Frame()
	settles(t, tt)
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if *hits != 2 {
		t.Errorf("%s: %d activations after copy appeared, want 2", what, *hits)
	}
}

func TestActionsKeepFocus(t *testing.T) {
	more, hits := false, 0
	word := func(s string) string {
		if more {
			return s
		}
		return ""
	}
	focusHolds(t, "Hero", "Start", &more, &hits, func(c *ui.Context) {
		_, first, _ := Hero(c, "Run your studio", word("Kilns, classes and members."), "Start", "Tour", HeroOptions{Eyebrow: word("New")})
		if first.Clicked() {
			hits++
		}
	})
	more, hits = false, 0
	focusHolds(t, "HeroSplit", "Start", &more, &hits, func(c *ui.Context) {
		_, b := HeroSplit(c, "Run your studio", "", "Start", photo(), HeroSplitOptions{Eyebrow: word("New")})
		if b.Clicked() {
			hits++
		}
	})
	more, hits = false, 0
	focusHolds(t, "CTA", "Start", &more, &hits, func(c *ui.Context) {
		_, b := CTA(c, "Ready?", word("Free for 14 days."), "Start", CTAOptions{})
		if b.Clicked() {
			hits++
		}
	})
	more, hits, at := false, 0, 0
	quotes := []Quote{{Text: "Fewer missed firings.", Name: "Maya Chen"}, {Text: "Classes fill.", Name: "Priya Nair"},
		{Text: "Billing runs itself.", Name: "Lucas Martin"}}
	focusHolds(t, "Testimonials", "Next quote", &more, &hits, func(c *ui.Context) {
		before := at
		Testimonials(c, quotes, &at, TestimonialsOptions{Title: word("Loved by studios")})
		if at != before {
			hits++
		}
	})
}

func TestBillingSaveExact(t *testing.T) {
	for _, c := range []struct {
		monthly, year float64
		want          string
	}{{10, 96, "Save 20%"}, {10, 108, "Save 10%"}, {24, 230.4, "Save 20%"}} {
		yearly := true
		tt := ui.NewTester(padded(func(ctx *ui.Context) { BillingSwitch(ctx, c.monthly, c.year, &yearly, BillingSwitchOptions{}) }), 480, 300)
		settles(t, tt)
		if !exact(tt, c.want) {
			t.Errorf("%v / %v: %q, want %q", c.monthly, c.year, tt.Texts(), c.want)
		}
	}
	yearly := false
	panics(t, "at least 0.01", func(c *ui.Context) { BillingSwitch(c, 0.001, 1, &yearly, BillingSwitchOptions{}) })
}

func TestSiteNavFoldsOnText(t *testing.T) {
	long := []string{"Kiln scheduling", "Member billing", "Class registration", "Studio reporting"}
	active := ""
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(520).Children(func() { SiteNav(c, "Glaze", long, &active, "Sign in", SiteNavOptions{}) })
	}, 520, 400)
	settles(t, tt)
	sign := rectOf(t, tt, "Sign in")
	for _, l := range long {
		if r, ok := tt.Find(l); ok && r.W > 0 && r.X+r.W > sign.X {
			t.Errorf("%q runs under Sign in: %v vs %v", l, r, sign)
		}
	}
	short := []string{"Features", "Pricing", "Customers", "Docs"}
	st := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(520).Children(func() { SiteNav(c, "Glaze", short, &active, "Sign in", SiteNavOptions{}) })
	}, 520, 400)
	settles(t, st)
	if _, ok := st.Find("Menu"); ok {
		t.Errorf("four short links folded at 520: %q", st.Texts())
	}
	nt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(320).Children(func() { SiteNav(c, "Glaze Studio Manager", short, &active, "Start free trial", SiteNavOptions{}) })
	}, 320, 400)
	settles(t, nt)
	if m := rectOf(t, nt, "Menu"); m.W == 0 || m.X+m.W > 320 {
		t.Errorf("Menu out of view at 320: %v", m)
	}
	click(t, nt, "Menu")
	settles(t, nt)
	if r, ok := nt.Find("Pricing"); !ok || r.H == 0 {
		t.Errorf("Menu opened nothing: %q", nt.Texts())
	}
}

func TestSiteNavMenuScrolls(t *testing.T) {
	var links []string
	for i := 1; i <= 40; i++ {
		links = append(links, fmt.Sprintf("Page %02d", i))
	}
	active := ""
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(320).Children(func() { SiteNav(c, "Glaze", links, &active, "Sign in", SiteNavOptions{}) })
	}, 320, 400)
	settles(t, tt)
	click(t, tt, "Menu")
	settles(t, tt)
	r := rectOf(t, tt, "Page 02")
	for range 8 {
		tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 400)
		tt.Advance(500 * time.Millisecond)
	}
	if l, ok := tt.Find("Page 40"); !ok || l.H == 0 || l.Y+l.H > 400 {
		t.Errorf("Page 40 out of reach: %v", l)
	}
}

func TestHeroSplitActionFits(t *testing.T) {
	act := "Explore all studio management features"
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(520).Children(func() { HeroSplit(c, "Every firing, on schedule", "", act, photo(), HeroSplitOptions{}) })
	}, 520, 600)
	settles(t, tt)
	title, b := rectOf(t, tt, "Every firing, on schedule"), rectOf(t, tt, act)
	if b.X+b.W > title.X+title.W+0.5 {
		t.Errorf("action ends at %v, its column at %v", b.X+b.W, title.X+title.W)
	}
}

func TestNewsletterEnterAndTitle(t *testing.T) {
	email, sent := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) { Newsletter(c, &email, &sent, NewsletterOptions{}) }), 480, 400)
	settles(t, tt)
	click(t, tt, "Email address")
	tt.Type("maya@kiln.co")
	tt.Frame()
	before := baseline()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	cues(t, before, sound.Success)
	if !sent {
		t.Error("Enter did not report the signup")
	}
	panics(t, "Studio letter", func(c *ui.Context) { Newsletter(c, nil, nil, NewsletterOptions{Title: "Studio letter"}) })
}
