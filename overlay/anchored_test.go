package overlay

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// anchor is a toggle button for a test overlay.
func anchor(c *ui.Context, label string, open *bool) ui.Element {
	b := ui.ButtonBase(c).Label(label).Size(100, 32).Children(func() { ui.Text(c, label) })
	if b.Clicked() {
		*open = !*open
	}
	return b
}

func TestPopoverCardAndQuickView(t *testing.T) {
	open, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		a := anchor(c, "Info", &open)
		if PopoverCard(c, a, &open, func() { ui.Text(c, "Details") }).Changed() {
			changed++
		}
	}, 600, 400)
	before := baseline()
	click(t, tt, "Info")
	if !tt.HasText("Details") || !tt.FramePending() {
		t.Fatalf("did not scale in: %q", tt.Texts())
	}
	settles(t, tt, "PopoverCard")
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || changed != 1 {
		t.Errorf("Escape: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Open)

	open = false
	tq := ui.NewTester(func(c *ui.Context) {
		a := anchor(c, "Row 1", &open)
		QuickView(c, a, "Invoice", &open, QuickViewOptions{}, func() { ui.Text(c, "Total") })
	}, 600, 400)
	before = baseline()
	click(t, tq, "Row 1")
	if !tq.HasText("Total") {
		t.Fatalf("quick view %q", tq.Texts())
	}
	click(t, tq, "Close")
	if open {
		t.Error("quick view still open")
	}
	cues(t, before, sound.Open)
}

func TestHoverCardAndRichTooltip(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		a := ui.Box(c).Label("Ada").Size(100, 32).Children(func() { ui.Text(c, "Ada") })
		HoverCard(c, "Profile", a, func() { ui.Text(c, "Mathematician") })
	}, 600, 400)
	if tt.HasText("Mathematician") {
		t.Fatal("shown before hover")
	}
	r := find(t, tt, "Ada")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !tt.HasText("Mathematician") || !tt.FramePending() {
		t.Fatalf("hover did not fade it in: %q", tt.Texts())
	}
	settles(t, tt, "HoverCard")
	card := find(t, tt, "Profile")
	// Cross the gap in steps, a frame apart.
	for _, y := range []float32{r.Y + r.H + 2, r.Y + r.H + 5, card.Y + card.H/2} {
		tt.Move(card.X+card.W/2, y)
		tt.Advance(16 * time.Millisecond)
	}
	if !tt.HasText("Mathematician") {
		t.Fatal("left on its way into the card")
	}
	tt.Move(590, 390)
	tt.Frame()
	if !tt.HasText("Mathematician") {
		t.Fatal("went before its grace")
	}
	tt.Advance(hoverGrace + 50*time.Millisecond)
	if tt.HasText("Mathematician") {
		t.Error("stayed after the pointer left")
	}

	tr := ui.NewTester(func(c *ui.Context) {
		a := ui.ButtonBase(c).Label("Share").Size(100, 32)
		RichTooltip(c, a, "Share", "Send a link to others")
	}, 600, 400)
	r = find(t, tr, "Share")
	tr.Move(r.X+r.W/2, r.Y+r.H/2)
	tr.Advance(time.Second)
	if !tr.HasText("Send a link to others") {
		t.Fatalf("tip did not show: %q", tr.Texts())
	}
	tip, bg := find(t, tr, "Send a link to others"), ui.LightTheme().Background
	if p := tr.Image().RGBAAt(int(tip.X)+1, int(tip.Y)+1); p.R == bg.R && p.G == bg.G && p.B == bg.B {
		t.Error("the tip text is not on the tip")
	}
	panics(t, "no text", func() { render(func(c *ui.Context) { RichTooltip(c, ui.Box(c), "x", "") }) })
}

func TestCoachmark(t *testing.T) {
	open, changed := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		a := ui.ButtonBase(c).Label("Export").Size(100, 32)
		if Coachmark(c, a, "New: export", "Save as PDF here.", &open, CoachmarkOptions{}).Changed() {
			changed++
		}
	}, 600, 400)
	if !tt.HasText("Save as PDF here.") || !tt.FramePending() {
		t.Fatalf("tip %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Got it")
	if open || changed != 1 || tt.HasText("Save as PDF here.") {
		t.Errorf("got it: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Close)
}

func TestNotificationPanel(t *testing.T) {
	open := false
	groups := []NoticeGroup{{"Today", []Notice{{"Build passed", "main is green", true}, {"New comment", "On PR 12", true}}},
		{"Earlier", []Notice{{"Deployed", "v1.2", false}}}}
	tt := ui.NewTester(func(c *ui.Context) { NotificationPanel(c, "Notifications", groups, &open, NotificationPanelOptions{}) }, 600, 500)
	if !tt.HasText("2") {
		t.Fatalf("unread badge %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Notifications")
	if !open || !tt.HasText("Earlier") || !tt.HasText("On PR 12") {
		t.Fatalf("panel %q", tt.Texts())
	}
	cues(t, before, sound.Open)
	none := false
	shut := ui.NewTester(func(c *ui.Context) { NotificationPanel(c, "Bell", nil, &none, NotificationPanelOptions{}) }, 400, 300)
	click(t, shut, "Bell")
	if !shut.HasText("You're all caught up") {
		t.Errorf("empty panel %q", shut.Texts())
	}
}

func TestMenuAndContextMenu(t *testing.T) {
	list := []MenuItem{{"Rename", "pencil", "F2"}, {"Delete", "trash-2", ""}}
	chosen, open, changed := -1, false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Menu(c, "Actions", list, &chosen, &open).Changed() {
			changed++
		}
	}, 600, 400)
	before := baseline()
	click(t, tt, "Actions")
	if !open || !tt.HasText("F2") {
		t.Fatalf("menu %q", tt.Texts())
	}
	click(t, tt, "Delete")
	if chosen != 1 || open || changed != 2 {
		t.Errorf("menu pick: chosen %d, open %v, changed %d", chosen, open, changed)
	}
	cues(t, before, sound.Open, sound.Select)

	chosen, open = -1, false
	tc := ui.NewTester(func(c *ui.Context) {
		target := ui.ButtonBase(c).Label("File").Size(120, 40)
		ContextMenu(c, "File actions", target, list, &chosen, &open)
	}, 600, 400)
	before = baseline()
	if err := tc.RightClick("File"); err != nil {
		t.Fatal(err)
	}
	tc.Frame()
	if !open || !tc.HasText("Rename") {
		t.Fatalf("right click: open %v, %q", open, tc.Texts())
	}
	if tc.Menu() != nil {
		t.Errorf("a system menu opened too: %q", tc.Menu())
	}
	click(t, tc, "Rename")
	if chosen != 0 || open {
		t.Errorf("context pick: chosen %d, open %v", chosen, open)
	}
	tc.Key(0, ui.KeyTab)
	tc.Key(ui.Shift, ui.KeyF10)
	tc.Frame()
	if !open {
		t.Error("Shift+F10 did not open the menu")
	}
	cues(t, before, sound.Open, sound.Select, sound.Open)
	panics(t, "bad choice", func() { render(func(c *ui.Context) { ch, o := 2, false; Menu(c, "x", list, &ch, &o) }) })
}
