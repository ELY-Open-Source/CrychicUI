package feedback

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestAlertAndFriends(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			Alert(c, ToneWarning, "Disk almost full", "Free some space.")
			InlineMessage(c, ToneSuccess, "Saved")
			Notice(c, "inbox", "All caught up", "No new mail.")
			Callout(c, ToneInfo, "Tip", func() { ui.Text(c, "Press K to search.") })
		})
	}, 500, 500)
	for _, s := range []string{"Disk almost full", "Free some space.", "Saved", "All caught up", "Tip", "Press K to search."} {
		find(t, tt, s)
	}
	settles(t, tt, "arrivals")
	panics(t, "unknown tone", func() { render(func(c *ui.Context) { Alert(c, Tone(9), "x", "") }) })
	panics(t, "empty alert", func() { render(func(c *ui.Context) { Alert(c, ToneInfo, "", "") }) })
	panics(t, "nil callout", func() { render(func(c *ui.Context) { Callout(c, ToneInfo, "x", nil) }) })
}

func TestDismissibleAlert(t *testing.T) {
	open, changes := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if DismissibleAlert(c, ToneDanger, "Payment failed", "Card declined.", &open, DismissibleAlertOptions{}).Changed() {
			changes++
		}
	}, 400, 200)
	tt.Advance(time.Second)
	before := baseline()
	click(t, tt, "Dismiss")
	if open || changes != 1 {
		t.Fatalf("open %v changes %d", open, changes)
	}
	cues(t, before, sound.Close)
	settles(t, tt, "collapse")
	if tt.HasText("Payment failed") {
		t.Error("alert still shown")
	}
}

func TestBanner(t *testing.T) {
	open, acts := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		_, act := Banner(c, "A new version is out.", "Update", &open, BannerOptions{Close: "Hide"})
		if act.Clicked() {
			acts++
		}
	}, 500, 200)
	tt.Advance(time.Second)
	before := baseline()
	click(t, tt, "Update")
	click(t, tt, "Hide")
	if open || acts != 1 {
		t.Fatalf("open %v acts %d", open, acts)
	}
	cues(t, before, sound.Tap, sound.Close)
	settles(t, tt, "banner collapse")
}

func TestToaster(t *testing.T) {
	acted := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ui.Button(c, "Add").Clicked() {
			c.AddToast(ui.Toast{ID: "up", Title: "Uploaded", Description: "report.pdf", Type: "success", Timeout: -1,
				Action: "Open", OnAction: func() { acted++ }})
		}
		Toaster(c, ToasterOptions{})
	}, 600, 400)
	before := baseline()
	click(t, tt, "Add")
	find(t, tt, "report.pdf")
	tt.Frame()
	cues(t, before, sound.Notify)
	settles(t, tt, "toast entry")
	click(t, tt, "Open")
	if acted != 1 {
		t.Fatalf("action ran %d times", acted)
	}
	click(t, tt, "Add")
	click(t, tt, "Close")
	tt.Advance(2 * time.Second)
	if tt.HasText("report.pdf") {
		t.Error("toast not closed")
	}
	cues(t, before, sound.Notify, sound.Tap, sound.Notify, sound.Close)
	panics(t, "unknown type", func() {
		added := false
		ui.NewTester(func(c *ui.Context) {
			if !added {
				added = true
				c.AddToast(ui.Toast{Title: "x", Type: "odd"})
			}
			Toaster(c, ToasterOptions{})
		}, 200, 200).Frame()
	})
}

func TestToasterEscapeAndFirstBuild(t *testing.T) {
	added := false
	tt := ui.NewTester(func(c *ui.Context) {
		if !added {
			added = true
			c.AddToast(ui.Toast{ID: "a", Title: "Hello", Timeout: -1, Action: "Go"})
		}
		Toaster(c, ToasterOptions{})
	}, 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("first toast did not slide in")
	}
	tt.Advance(time.Second)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if tt.HasText("Hello") {
		t.Fatal("Escape left the toast")
	}
	cues(t, before, sound.Close)
}

func TestToasterHiddenToasts(t *testing.T) {
	var add []string
	closeID := ""
	tt := ui.NewTester(func(c *ui.Context) {
		for _, id := range add {
			c.AddToast(ui.Toast{ID: id, Title: "Toast " + id, Timeout: -1})
		}
		add = nil
		if closeID != "" {
			c.CloseToast(closeID)
			closeID = ""
		}
		Toaster(c, ToasterOptions{})
	}, 600, 500)
	before := baseline()
	add = []string{"A", "B", "C"}
	tt.Frame()
	add = []string{"D"}
	tt.Frame()
	if tt.HasText("Toast A") {
		t.Fatal("oldest toast not hidden")
	}
	closeID = "D"
	tt.Frame()
	tt.Advance(time.Second)
	if !tt.HasText("Toast A") {
		t.Fatal("hidden toast did not come back")
	}
	cues(t, before, sound.Notify, sound.Notify, sound.Notify, sound.Notify)
	closeID = "A"
	tt.Frame()
	add = []string{"A"}
	tt.Frame()
	cues(t, before, sound.Notify, sound.Notify, sound.Notify, sound.Notify, sound.Notify)
}

func TestToastCard(t *testing.T) {
	acts, closes := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		_, act, dismiss := ToastCard(c, ToneInfo, "Invite sent", "To ada@example.com", ToastCardOptions{Action: "Copy link"})
		if act.Clicked() {
			acts++
		}
		if dismiss.Clicked() {
			closes++
		}
	}, 500, 200)
	before := baseline()
	click(t, tt, "Copy link")
	click(t, tt, "Close")
	if acts != 1 || closes != 1 {
		t.Fatalf("acts %d closes %d", acts, closes)
	}
	cues(t, before, sound.Tap, sound.Close)
	settles(t, tt, "toast card")
}

func TestSnackbarInk(t *testing.T) {
	open := true
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		th.Inverse = th.Accent
		c.SetTheme(&th)
		Snackbar(c, "MMMM", "", &open)
	}, 400, 200)
	tt.Advance(time.Second)
	r := find(t, tt, "MMMM")
	img, inked := tt.Image(), false
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			p := img.RGBAAt(x, y)
			inked = inked || p.R > 240 && p.G > 240 && p.B > 240
		}
	}
	if !inked {
		t.Error("message not drawn in Background over Inverse")
	}
}

func TestSnackbar(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) { Snackbar(c, "Message archived", "Undo", &open) }, 500, 200)
	before := baseline()
	open = true
	tt.Frame()
	tt.Frame()
	cues(t, before, sound.Notify)
	settles(t, tt, "snackbar rise")
	click(t, tt, "Undo")
	if open {
		t.Fatal("action did not close the bar")
	}
	cues(t, before, sound.Notify, sound.Tap)
}
