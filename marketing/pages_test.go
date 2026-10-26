package marketing

import (
	"image"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func photo() *ui.Bitmap { return ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 8, 6))) }

func TestHeroActions(t *testing.T) {
	var picks []bool
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, a, b := Hero(c, "Run your studio", "One place.", "Start trial", "Book demo", HeroOptions{Eyebrow: "New"})
		if a.Clicked() {
			picks = append(picks, true)
		}
		if b.Clicked() {
			picks = append(picks, false)
		}
	}), 520, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Start trial")
	click(t, tt, "Book demo")
	cues(t, before, sound.Tap, sound.Tap)
	if len(picks) != 2 || !picks[0] || picks[1] {
		t.Errorf("picks %v", picks)
	}
	panics(t, "two different actions", func(c *ui.Context) { Hero(c, "T", "", "Go", "Go", HeroOptions{}) })
}

func TestHeroSplitAction(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		HeroSplit(c, "Every firing", "Text", "See how", photo(), HeroSplitOptions{Points: []string{"Loads", "Shelves"}})
	}), 520, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "See how")
	cues(t, before, sound.Tap)
	panics(t, "needs a title, an action and media", func(c *ui.Context) { HeroSplit(c, "T", "", "Go", nil, HeroSplitOptions{}) })
	panics(t, "blank or repeated", func(c *ui.Context) {
		HeroSplit(c, "T", "", "Go", photo(), HeroSplitOptions{Points: []string{"a", "a"}})
	})
}

func TestFeatureGridChecks(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		FeatureGrid(c, []Feature{{"flame", "Kiln", "Loads"}, {"palette", "Glazes", "Recipes"}}, FeatureGridOptions{Title: "Built for studios"})
	}), 520, 400)
	settles(t, tt)
	for _, s := range []string{"Built for studios", "Kiln", "Recipes"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	panics(t, "1 to 4 columns", func(c *ui.Context) { FeatureGrid(c, []Feature{{"flame", "A", ""}}, FeatureGridOptions{Columns: 5}) })
	panics(t, "needs an icon", func(c *ui.Context) { FeatureGrid(c, []Feature{{"", "A", ""}}, FeatureGridOptions{}) })
	panics(t, "blank or repeated", func(c *ui.Context) {
		FeatureGrid(c, []Feature{{"flame", "A", ""}, {"sun", "A", ""}}, FeatureGridOptions{})
	})
}

func TestFeatureRowAlternates(t *testing.T) {
	rows := []Showcase{{Title: "Fill every wheel", Text: "Classes", Media: photo()}, {Title: "Know who fires", Text: "Members", Media: photo()}}
	tt := ui.NewTester(padded(func(c *ui.Context) { FeatureRow(c, rows, FeatureRowOptions{}) }), 620, 700)
	settles(t, tt)
	first, second := rectOf(t, tt, "Fill every wheel"), rectOf(t, tt, "Know who fires")
	if !(second.X > first.X+100) {
		t.Errorf("second row's text did not move right: %.1f vs %.1f", second.X, first.X)
	}
	panics(t, "needs media", func(c *ui.Context) { FeatureRow(c, []Showcase{{Title: "A"}}, FeatureRowOptions{}) })
}

func TestBentoWideSpans(t *testing.T) {
	tiles := []Bento{{Icon: "flame", Title: "Wide A", Wide: true}, {Icon: "users", Title: "Left B"}, {Icon: "sun", Title: "Right C"},
		{Icon: "palette", Title: "Wide D", Wide: true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { BentoFeatures(c, tiles, BentoFeaturesOptions{}) }), 520, 900)
	settles(t, tt)
	a, b, cc, d := rectOf(t, tt, "Wide A"), rectOf(t, tt, "Left B"), rectOf(t, tt, "Right C"), rectOf(t, tt, "Wide D")
	if !(cc.X > b.X+100) || b.X != a.X || d.X != a.X || !(b.Y > a.Y) || !(d.Y > b.Y) {
		t.Errorf("tiles out of place: a %v b %v c %v d %v", a, b, cc, d)
	}
}

func TestLogoCloudLoops(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		LogoCloud(c, []Logo{{"Clayworks", "flame"}, {"Greenware", "leaf"}}, LogoCloudOptions{})
	}), 520, 300)
	settles(t, tt)
	if !tt.FramePending() || !exact(tt, "Trusted by teams at") {
		t.Errorf("marquee: %q", tt.Texts())
	}
	panics(t, "two logos or more", func(c *ui.Context) { LogoCloud(c, []Logo{{"A", "sun"}}, LogoCloudOptions{}) })
}

func TestStatsBandCounts(t *testing.T) {
	stats := []Stat{{Value: 1240, Suffix: "+", Label: "Studios"}, {Value: 3.2, Decimals: 1, Suffix: "M", Label: "Firings"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { StatsBand(c, stats, StatsBandOptions{}) }), 520, 400)
	tt.Advance(300 * time.Millisecond)
	if exact(tt, "1,240+") {
		t.Errorf("no count-up: %q", tt.Texts())
	}
	settles(t, tt)
	if !exact(tt, "1,240+") || !exact(tt, "3.2M") {
		t.Errorf("final: %q", tt.Texts())
	}
	panics(t, "0 to 3 places", func(c *ui.Context) { StatsBand(c, []Stat{{Value: 1, Decimals: 4, Label: "A"}}, StatsBandOptions{}) })
	panics(t, "within 1e15", func(c *ui.Context) { StatsBand(c, []Stat{{Value: 1e16, Label: "A"}}, StatsBandOptions{}) })
}

func TestTestimonialsPages(t *testing.T) {
	quotes := []Quote{{Text: "First", Name: "Priya Nair"}, {Text: "Second", Name: "Daniel Okafor", Rating: 5}, {Text: "Third", Name: "Sofia Rossi"}}
	at := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { Testimonials(c, quotes, &at, TestimonialsOptions{}) }), 480, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Previous quote")
	cues(t, before)
	before = baseline()
	click(t, tt, "Next quote")
	cues(t, before, sound.Swipe)
	settles(t, tt)
	if at != 1 || !exact(tt, "Second") {
		t.Errorf("next: %d %q", at, tt.Texts())
	}
	click(t, tt, "Quote 3")
	before = baseline()
	click(t, tt, "Next quote")
	cues(t, before)
	if at != 2 {
		t.Errorf("at %d", at)
	}
	panics(t, "position among them", func(c *ui.Context) { a := 3; Testimonials(c, quotes, &a, TestimonialsOptions{}) })
	panics(t, "rating of 0 to 5", func(c *ui.Context) {
		a := 0
		Testimonials(c, []Quote{{Text: "x", Name: "y", Rating: 6}}, &a, TestimonialsOptions{})
	})
}
