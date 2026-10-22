package choice

import (
	"slices"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCheck(t *testing.T) {
	on, changes := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Check(c, "Sync", &on).Changed() {
			changes++
		}
	}, 200, 80)
	before := baseline()
	click(t, tt, "Sync")
	if !on || changes != 1 {
		t.Fatalf("on %v changes %d", on, changes)
	}
	cues(t, before, sound.Toggle)
	settles(t, tt, "Check mark")
	panics(t, "nil state", func() { ui.Render(func(c *ui.Context) { Check(c, "x", nil) }, 50, 50, 1) })
}

func TestCheckCard(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) { CheckCard(c, "Backups", "Nightly copies.", &on) }, 320, 120)
	find(t, tt, "Nightly copies.")
	before := baseline()
	click(t, tt, "Backups")
	if !on {
		t.Fatal("card did not check")
	}
	cues(t, before, sound.Toggle)
	settles(t, tt, "CheckCard border")
}

func TestCheckGroup(t *testing.T) {
	sel := []string{"Mail"}
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if CheckGroup(c, []string{"Mail", "Calendar", "Messages"}, &sel, CheckGroupOptions{AllLabel: "Everything"}).Changed() {
			changes++
		}
	}, 300, 200)
	before := baseline()
	click(t, tt, "Everything")
	if len(sel) != 3 || changes != 1 {
		t.Fatalf("selection %v changes %d", sel, changes)
	}
	click(t, tt, "Calendar")
	if slices.Contains(sel, "Calendar") || changes != 2 {
		t.Fatalf("selection %v changes %d", sel, changes)
	}
	click(t, tt, "Everything")
	if len(sel) != 3 {
		t.Errorf("mixed parent did not select all: %v", sel)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle)
	settles(t, tt, "CheckGroup marks")
	panics(t, "unknown selection", func() {
		bad := []string{"Nope"}
		ui.Render(func(c *ui.Context) { CheckGroup(c, []string{"A"}, &bad, CheckGroupOptions{}) }, 50, 50, 1)
	})
	panics(t, "duplicate labels", func() {
		ui.Render(func(c *ui.Context) { CheckGroup(c, []string{"A", "A"}, &sel, CheckGroupOptions{}) }, 50, 50, 1)
	})
}

func TestToggle(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) { Toggle(c, "Wi-Fi", &on) }, 200, 80)
	before := baseline()
	click(t, tt, "Wi-Fi")
	if !on {
		t.Fatal("switch did not turn on")
	}
	cues(t, before, sound.Toggle)
	settles(t, tt, "Toggle knob")
}

func TestToggleRow(t *testing.T) {
	on, changes := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ToggleRow(c, "Sounds", "Play cues on actions.", &on).Changed() {
			changes++
		}
	}, 360, 100)
	before := baseline()
	click(t, tt, "Play cues on actions.")
	if on || changes != 1 {
		t.Fatalf("on %v changes %d", on, changes)
	}
	cues(t, before, sound.Toggle)
	settles(t, tt, "ToggleRow knob")
}

func TestOptions(t *testing.T) {
	plan, seats, changes := "Free", "One", 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			if OptionGroup(c, "Plan", []string{"Free", "Pro"}, &plan, OptionGroupOptions{Row: true}).Changed() {
				changes++
			}
			OptionCard(c, "Team", "For five or more.", "Team", &seats)
		})
	}, 400, 200)
	before := baseline()
	click(t, tt, "Pro")
	if plan != "Pro" || changes != 1 {
		t.Fatalf("plan %q changes %d", plan, changes)
	}
	tt.Key(0, ui.KeyLeft)
	if plan != "Free" || changes != 2 {
		t.Fatalf("arrow left gave %q, changes %d", plan, changes)
	}
	click(t, tt, "Team")
	if seats != "Team" {
		t.Fatalf("card gave %q", seats)
	}
	cues(t, before, sound.Select, sound.Select, sound.Select)
	settles(t, tt, "Option dots")
}

func TestSegments(t *testing.T) {
	view, icons := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			Segments(c, "View", []string{"Day", "Week", "Month"}, &view)
			IconSegments(c, "Layout", []SegmentItem{{"List", "list"}, {"Grid", "layout-grid"}}, &icons)
		})
	}, 400, 160)
	before := baseline()
	click(t, tt, "Month")
	if view != 2 {
		t.Fatalf("view %d", view)
	}
	settles(t, tt, "Segments thumb")
	click(t, tt, "Grid")
	if icons != 1 {
		t.Fatalf("icons %d", icons)
	}
	cues(t, before, sound.Select, sound.Select)
	settles(t, tt, "IconSegments thumb")
}

func TestChips(t *testing.T) {
	size := "M"
	filters := []string{}
	toggled := []string{}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			ChipChoice(c, "Size", []string{"S", "M", "L"}, &size)
			FilterChips(c, "Filters", []string{"Open", "Mine"}, &filters)
			ToggleSet(c, "Style", []string{"Bold", "Italic"}, &toggled)
		})
	}, 400, 200)
	before := baseline()
	click(t, tt, "L")
	click(t, tt, "Mine")
	click(t, tt, "Italic")
	click(t, tt, "Bold")
	if size != "L" || !slices.Equal(filters, []string{"Mine"}) || !slices.Equal(toggled, []string{"Italic", "Bold"}) {
		t.Fatalf("size %q filters %v toggled %v", size, filters, toggled)
	}
	cues(t, before, sound.Select, sound.Toggle, sound.Toggle, sound.Toggle)
	settles(t, tt, "chips")
	if !tt.HasText("1") {
		t.Errorf("filter count missing: %q", tt.Texts())
	}
}

func TestSegmentsLabelOverThumb(t *testing.T) {
	view := 0
	tt := ui.NewTester(func(c *ui.Context) { Segments(c, "View", []string{"Day", "Week", "Month"}, &view) }, 400, 80)
	click(t, tt, "Week")
	settles(t, tt, "Segments thumb")
	r, img, dark := find(t, tt, "Week"), tt.Image(), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); int(p.R)+int(p.G)+int(p.B) < 300 {
				dark++
			}
		}
	}
	if dark < 10 {
		t.Errorf("the chosen label is hidden: %d text pixels", dark)
	}
}
