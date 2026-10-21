package display

import (
	"math"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestInfoCard(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		InfoCard(c, "Title", "Body text", func() { ui.Text(c, "Footer") })
	}, 400, 200)
	r := find(t, tt, "Title")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "InfoCard hover")
	find(t, tt, "Body text")
	find(t, tt, "Footer")
}

func TestMediaCard(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		MediaCard(c, "Photo", "A useful caption", testBitmap())
	}, 400, 250)
	r := find(t, tt, "Photo")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "MediaCard zoom")
	find(t, tt, "A useful caption")
	panics(t, "MediaCard image", func() { ui.Render(func(c *ui.Context) { MediaCard(c, "Photo", "Body", nil) }, 10, 10, 1) })
}

func TestStatCard(t *testing.T) {
	value := 1234.5
	tt := ui.NewTester(func(c *ui.Context) { StatCard(c, "Revenue", value, 4.5) }, 400, 200)
	value = 2345.5
	tt.Frame()
	settles(t, tt, "StatCard value")
	find(t, tt, "2,345.5")
	find(t, tt, "+4.5%")
}

func TestStat(t *testing.T) {
	value := 12.0
	tt := ui.NewTester(func(c *ui.Context) { Stat(c, "Users", value) }, 300, 200)
	value = 42
	tt.Frame()
	settles(t, tt, "Stat value")
	find(t, tt, "42")
}

func TestKPI(t *testing.T) {
	value := 12000.0
	tt := ui.NewTester(func(c *ui.Context) { KPI(c, "Growth", value, 2.5) }, 400, 200)
	settles(t, tt, "KPI arrow")
	find(t, tt, "12,000")
	panics(t, "KPI delta", func() { ui.Render(func(c *ui.Context) { KPI(c, "Growth", 1, math.NaN()) }, 10, 10, 1) })
}

func TestDelta(t *testing.T) {
	value := 1.5
	tt := ui.NewTester(func(c *ui.Context) { Delta(c, value) }, 300, 200)
	value = -2
	tt.Frame()
	settles(t, tt, "Delta value")
	find(t, tt, "-2.0")
}

func TestListItem(t *testing.T) {
	t.Run("clickable", func(t *testing.T) {
		clicks := 0
		tt := ui.NewTester(func(c *ui.Context) {
			if ListItem(c, "Inbox", "Three messages", ListItemOptions{Icon: "inbox", Clickable: true}, func() {
				ui.Text(c, "Now")
			}).Clicked() {
				clicks++
			}
		}, 500, 200)
		r := find(t, tt, "Inbox")
		tt.Move(r.X+r.W/2, r.Y+r.H/2)
		settles(t, tt, "ListItem hover")
		before := baseline()
		if err := tt.Click("Inbox"); err != nil {
			t.Fatal(err)
		}
		if clicks != 1 || !appended(before, sound.Tap) {
			t.Errorf("clicks %d, cues %v", clicks, sound.Recent())
		}
	})
	t.Run("static", func(t *testing.T) {
		tt := ui.NewTester(func(c *ui.Context) {
			ListItem(c, "Inbox", "Three messages", ListItemOptions{}, nil)
		}, 500, 200)
		r := find(t, tt, "Inbox")
		tt.Move(r.X+r.W/2, r.Y+r.H/2)
		settles(t, tt, "ListItem hover")
		before := baseline()
		if err := tt.Click("Inbox"); err != nil {
			t.Fatal(err)
		}
		if !appended(before) {
			t.Error("static row played a cue")
		}
	})
}

func TestDescriptionList(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		DescriptionList(c, []KeyValuePair{{Key: "Name", Value: "Crychic"}, {Key: "Kind", Value: "UI"}})
	}, 500, 200)
	settles(t, tt, "DescriptionList")
	find(t, tt, "Name")
	find(t, tt, "Crychic")
}

func TestKeyValue(t *testing.T) {
	value := "Old"
	tt := ui.NewTester(func(c *ui.Context) { KeyValue(c, "Status", value) }, 300, 200)
	value = "New"
	tt.Frame()
	settles(t, tt, "KeyValue")
	find(t, tt, "New")
}

func TestTimeline(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Timeline(c, []string{"Created", "Shipped"}) }, 400, 200)
	settles(t, tt, "Timeline")
	find(t, tt, "Created")
	find(t, tt, "Shipped")
}

func TestActivity(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Activity(c, "Ada", "edited", "Spec", "now") }, 500, 200)
	settles(t, tt, "Activity")
	find(t, tt, "Ada")
	find(t, tt, "edited")
	find(t, tt, "Spec")
}

func TestRibbon(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Ribbon(c, "New", func() { ui.Box(c).Size(120, 80) })
	}, 300, 200)
	settles(t, tt, "Ribbon")
	find(t, tt, "New")
}

func TestIndicator(t *testing.T) {
	var accent, success ui.Color
	color := ui.Color{}
	tt := ui.NewTester(func(c *ui.Context) {
		accent, success = c.Theme().Accent, c.Theme().Success
		if color.A == 0 {
			color = accent
		}
		Indicator(c, color)
	}, 100, 100)
	tt.Frame()
	color = success
	tt.Frame()
	settles(t, tt, "Indicator color")
}

func TestIconView(t *testing.T) {
	name := "house"
	tt := ui.NewTester(func(c *ui.Context) { IconView(c, name, IconViewOptions{}) }, 100, 100)
	name = "settings"
	tt.Frame()
	settles(t, tt, "IconView swap")
	panics(t, "IconView size", func() { ui.Render(func(c *ui.Context) { IconView(c, "house", IconViewOptions{Size: -1}) }, 10, 10, 1) })
}

func TestIconBox(t *testing.T) {
	var accent, success ui.Color
	color := ui.Color{}
	tt := ui.NewTester(func(c *ui.Context) {
		accent, success = c.Theme().Accent, c.Theme().Success
		if color.A == 0 {
			color = accent
		}
		IconBox(c, "bell", IconBoxOptions{Color: color})
	}, 100, 100)
	tt.Frame()
	color = success
	tt.Frame()
	settles(t, tt, "IconBox tint")
}

func TestLogo(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Logo(c, "sparkles", "CrychicUI") }, 300, 200)
	settles(t, tt, "Logo turn")
	find(t, tt, "CrychicUI")
}

func TestPlaceholder(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Placeholder(c, "Drop files", 160, 80) }, 300, 200)
	r := find(t, tt, "Drop files")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "Placeholder hover")
	panics(t, "Placeholder size", func() { ui.Render(func(c *ui.Context) { Placeholder(c, "Drop", 0, 10) }, 10, 10, 1) })
}
