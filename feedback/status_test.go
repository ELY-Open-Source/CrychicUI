package feedback

import (
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestStatusLines(t *testing.T) {
	saving, sync, conn := true, Syncing, Reconnecting
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			SaveStatus(c, saving, SaveStatusOptions{})
			SyncStatus(c, sync, SyncStatusOptions{})
			ConnectionStatus(c, conn, ConnectionStatusOptions{})
			LiveDot(c, LiveDotOptions{})
		})
	}, 300, 200)
	for _, s := range []string{"Saving…", "Syncing", "Reconnecting", "LIVE"} {
		find(t, tt, s)
	}
	saving, sync, conn = false, SyncFailed, Online
	tt.Frame()
	for _, s := range []string{"Saved", "Sync failed", "Online"} {
		find(t, tt, s)
	}
	panics(t, "bad sync", func() { render(func(c *ui.Context) { SyncStatus(c, 9, SyncStatusOptions{}) }) })
	panics(t, "bad connection", func() { render(func(c *ui.Context) { ConnectionStatus(c, 9, ConnectionStatusOptions{}) }) })
}

func TestQuotaAndStorage(t *testing.T) {
	used := 50.0
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		ui.Column(c).Gap(8).Width(300).Children(func() {
			QuotaBar(c, "Seats", used, 100, QuotaBarOptions{})
			StorageBar(c, "Disk", []StoragePart{{"Photos", 40, th.Accent}, {"Apps", 20, th.Success}}, 128, StorageBarOptions{})
		})
	}, 300, 200)
	find(t, tt, "50 of 100")
	find(t, tt, "40 GB")
	settles(t, tt, "bars grow")
	used = 90
	tt.Frame()
	tt.Advance(2 * time.Second)
	r := find(t, tt, "Seats")
	p := tt.Image().RGBAAt(5, int(r.Y+r.H)+8)
	if p.R != th.Warning.R || p.G != th.Warning.G {
		t.Errorf("quota at 90%% is %v, want warning %v", p, th.Warning)
	}
	panics(t, "parts over total", func() {
		render(func(c *ui.Context) { StorageBar(c, "x", []StoragePart{{"A", 10, th.Accent}}, 5, StorageBarOptions{}) })
	})
}

func TestMeters(t *testing.T) {
	score := 1
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Width(300).Children(func() {
			Strength(c, "Password", score, StrengthOptions{})
			Battery(c, "Battery", 0.15, false)
		})
	}, 300, 300)
	find(t, tt, "Weak")
	score = 4
	tt.Frame()
	find(t, tt, "Strong")
	settles(t, tt, "meters")
	panics(t, "score 5", func() { render(func(c *ui.Context) { Strength(c, "x", 5, StrengthOptions{}) }) })
	tt2 := ui.NewTester(func(c *ui.Context) { Battery(c, "Battery", 0.5, true) }, 200, 100)
	loops(t, tt2, "charging bolt")
}

func TestLevelMeter(t *testing.T) {
	level := 0.9
	tt := ui.NewTester(func(c *ui.Context) { LevelMeter(c, "Mic", level) }, 200, 60)
	th := ui.LightTheme()
	r := find(t, tt, "Mic")
	lit := func(i int, k ui.Color) bool { return tinted(tt, int(r.X)+i*8+3, int(r.Y+r.H)-2, k) }
	tt.Advance(time.Second)
	if !lit(9, th.Warning) || lit(10, th.Danger) {
		t.Errorf("level 0.9: bar 9 %v, bar 10 %v", lit(9, th.Warning), lit(10, th.Danger))
	}
	level = 0.1
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	if !lit(0, th.Success) || lit(5, th.Success) || !lit(9, th.Warning) {
		t.Errorf("drop: bar 0 %v, bar 5 %v, peak %v", lit(0, th.Success), lit(5, th.Success), lit(9, th.Warning))
	}
	tt.Advance(2 * time.Second)
	if tt.FramePending() || lit(9, th.Warning) {
		t.Errorf("peak rest: pending %v, held %v", tt.FramePending(), lit(9, th.Warning))
	}
	panics(t, "NaN level", func() { render(func(c *ui.Context) { LevelMeter(c, "x", math.NaN()) }) })
}

func TestSignal(t *testing.T) {
	strength := 2
	tt := ui.NewTester(func(c *ui.Context) { Signal(c, "Wi-Fi", strength) }, 200, 60)
	th := ui.LightTheme()
	r := find(t, tt, "Wi-Fi")
	lit := func(i int) bool { return tinted(tt, int(r.X)+i*6+2, int(r.Y+r.H)-2, th.Text) }
	settles(t, tt, "bars rise")
	if !lit(0) || !lit(1) || lit(2) || lit(3) {
		t.Errorf("strength 2: %v %v %v %v", lit(0), lit(1), lit(2), lit(3))
	}
	strength = 4
	tt.Frame()
	settles(t, tt, "bars fill")
	if !lit(3) {
		t.Error("strength 4: top bar dark")
	}
	panics(t, "signal 5", func() { render(func(c *ui.Context) { Signal(c, "x", 5) }) })
}

func TestServiceStatus(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ServiceStatus(c, []Service{{"API", Operational, 0.9998}, {"Search", Degraded, 0.97}}, ServiceStatusOptions{})
	}, 400, 200)
	find(t, tt, "99.98%")
	find(t, tt, "Degraded")
	settles(t, tt, "rows arrive")
	panics(t, "bad state", func() {
		render(func(c *ui.Context) { ServiceStatus(c, []Service{{"A", 7, 1}}, ServiceStatusOptions{}) })
	})
}

func TestAnnouncementBar(t *testing.T) {
	turns := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if AnnouncementBar(c, []string{"Free shipping", "New colors"}, AnnouncementBarOptions{Every: time.Second}).Changed() {
			turns++
		}
	}, 400, 80)
	find(t, tt, "Free shipping")
	tt.Advance(1100 * time.Millisecond)
	find(t, tt, "New colors")
	tt.Advance(time.Second)
	find(t, tt, "Free shipping")
	if turns != 2 {
		t.Errorf("rotations reported %d, want 2", turns)
	}
	if !tt.FramePending() {
		t.Error("rotation stopped")
	}
}

func TestUndoToast(t *testing.T) {
	open, undos := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if _, undo := UndoToast(c, "Message deleted", &open, UndoToastOptions{For: 2 * time.Second}); undo.Clicked() {
			undos++
		}
	}, 500, 200)
	before := baseline()
	click(t, tt, "Undo")
	if open || undos != 1 {
		t.Fatalf("open %v undos %d", open, undos)
	}
	open = true
	tt.Frame()
	tt.Advance(2500 * time.Millisecond)
	if open {
		t.Error("toast did not expire")
	}
	open = true
	tt.Frame()
	click(t, tt, "Close")
	cues(t, before, sound.Tap, sound.Close)
}

func TestUnlockToast(t *testing.T) {
	open := true
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { UnlockToast(c, "First streak", "Seven days in a row.", &open) }, 500, 200)
	find(t, tt, "First streak")
	settles(t, tt, "shine")
	tt.Frame()
	cues(t, before, sound.Success)
}

func TestBatteryBoltOverFill(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Battery(c, "Battery", 1, true) }, 200, 100)
	tt.Advance(2 * time.Second)
	r, ok := tt.Find("Battery")
	if !ok {
		t.Fatal("battery not found")
	}
	img, dark := tt.Image(), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); int(p.R)+int(p.G)+int(p.B) < 200 {
				dark++
			}
		}
	}
	if dark < 30 {
		t.Errorf("the bolt is hidden under a full fill: %d dark pixels", dark)
	}
	accent, spill := ui.LightTheme().Accent, 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		if y >= int(r.Y) && y < int(r.Y+r.H) {
			continue
		}
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if tinted(tt, x, y, accent) {
				spill++
			}
		}
	}
	if spill > 0 {
		t.Errorf("the fill spills past the battery: %d pixels", spill)
	}
}
