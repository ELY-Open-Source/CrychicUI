package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestOverlayDemosLog(t *testing.T) {
	type step func(*ui.Tester) error
	click := func(l string) step { return func(tt *ui.Tester) error { return tt.Click(l) } }
	key := func(k ui.Key) step { return func(tt *ui.Tester) error { tt.Key(0, k); return nil } }
	right := func(l string) step { return func(tt *ui.Tester) error { return tt.RightClick(l) } }
	typ := func(s string) step { return func(tt *ui.Tester) error { tt.Type(s); return nil } }
	drag := func(l string) step {
		return func(tt *ui.Tester) error {
			r, _ := tt.Find(l)
			tt.Press(r.X+r.W/2, r.Y+r.H/2)
			for i := 1; i <= 8; i++ {
				tt.Move(r.X+r.W/2+float32(i)*10, r.Y+r.H/2+float32(i)*5)
				tt.Advance(16 * time.Millisecond)
			}
			tt.Release(r.X+r.W/2+80, r.Y+r.H/2+40)
			return nil
		}
	}
	runs := map[string]struct {
		steps []step
		want  string
	}{
		"Modal":             {[]step{click("Open settings"), key(ui.KeyEscape)}, "closed"},
		"ConfirmDialog":     {[]step{click("Publish"), click("Confirm")}, "confirmed"},
		"DangerDialog":      {[]step{click("Delete project"), click("Delete")}, "deleted"},
		"PromptDialog":      {[]step{click("Rename"), typ(" 2"), key(ui.KeyEnter)}, "renamed, 10 characters"},
		"Drawer":            {[]step{click("Filters"), key(ui.KeyEscape)}, "closed"},
		"BottomSheet":       {[]step{click("Details"), key(ui.KeyUp)}, "open true, height 1"},
		"TopSheet":          {[]step{click("Announcements"), key(ui.KeyEscape)}, "closed"},
		"PopoverCard":       {[]step{click("About this plan"), key(ui.KeyEscape)}, "dismissed"},
		"RichTooltip":       {[]step{click("Share")}, "shared"},
		"ContextMenu":       {[]step{right("report.pdf"), click("Rename")}, "open false, chose 0"},
		"Menu":              {[]step{click("Actions"), click("Share")}, "open false, chose 2"},
		"CommandPalette":    {[]step{click("Commands"), key(ui.KeyEnter)}, "open false, chose 0"},
		"SpotlightSearch":   {[]step{click("Search"), key(ui.KeyEnter)}, "open false, chose [0 0]"},
		"Lightbox":          {[]step{click("View photos"), click("Next")}, "open true, photo 2"},
		"FullscreenModal":   {[]step{click("Open editor"), key(ui.KeyEscape)}, "closed"},
		"Inspector":         {[]step{click("Inspect layer"), click("Close")}, "closed"},
		"Tour":              {[]step{click("Start tour"), click("Next")}, "open true, step 2"},
		"Coachmark":         {[]step{click("Show tip"), click("Got it")}, "dismissed"},
		"WhatsNew":          {[]step{click("What's new"), click("Next")}, "open true, page 2"},
		"NotificationPanel": {[]step{click("Notifications")}, "open true"},
		"QuickView":         {[]step{click("Invoice 1042"), click("Close")}, "closed"},
		"PreviewPane":       {[]step{click("Budget")}, "selected Budget"},
		"EmojiPicker":       {[]step{click("dog")}, "chose 🐶"},
		"BusyVeil":          {[]step{click("Start saving")}, "busy true"},
		"LockOverlay":       {[]step{click("Unlock")}, "unlocked"},
		"Paywall":           {[]step{click("Upgrade")}, "upgrade"},
		"ShareSheet":        {[]step{click("Share"), click("Mail")}, "target 0, open false"},
		"ActionSheet":       {[]step{click("Message actions"), click("Share")}, "chose 2, open false"},
		"MiniPlayer":        {[]step{click("Play")}, "playing true"},
		"FloatingPanel":     {[]step{drag("Drag me anywhere.")}, "moved to"},
		"PiPBox":            {[]step{key(ui.KeyTab), key(ui.KeyLeft)}, "corner 2"},
		"CalloutPopover":    {[]step{click("Before you publish"), key(ui.KeyEscape)}, "dismissed"},
		"InfoPopover":       {[]step{click("About pricing")}, "open true"},
		"ConfirmPopover":    {[]step{click("Delete draft"), click("Yes")}, "deleted"},
		"SelectMenu":        {[]step{click("Fruit"), key(ui.KeyEnter)}, "chose 0, open false"},
		"MentionMenu":       {[]step{click("Message"), typ("hi @gr"), key(ui.KeyEnter)}, "mentioned"},
		"SlashMenu":         {[]step{click("Note"), typ("/ta"), key(ui.KeyEnter)}, "inserted Table"},
		"ShortcutSheet":     {[]step{click("Keyboard shortcuts"), key(ui.KeyEscape)}, "closed"},
		"StepDialog":        {[]step{click("Start setup"), click("Next")}, "step 1, open true"},
		"RenameDialog":      {[]step{click("Rename"), typ("notes"), key(ui.KeyEnter)}, "renamed, 9 characters"},
		"SessionTimeout":    {[]step{click("Simulate idle"), click("Stay signed in")}, "stayed"},
		"FeedbackWidget":    {[]step{click("Feedback"), typ("Nice"), click("Send")}, "sent 4 characters"},
		"SheetStack":        {[]step{click("Open settings"), click("Open Account"), key(ui.KeyEscape)}, "depth 1"},
	}
	for _, e := range Entries() {
		run, ok := runs[e.Name]
		if !ok || e.Category != "overlay" {
			continue
		}
		delete(runs, e.Name)
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		for _, s := range run.steps {
			if err := s(tt); err != nil {
				t.Fatalf("%s: %v", e.Name, err)
			}
			tt.Advance(500 * time.Millisecond)
		}
		if len(logs) == 0 || !strings.HasPrefix(logs[len(logs)-1], run.want) {
			t.Errorf("%s logged %q, want %q last", e.Name, logs, run.want)
		}
	}
	if len(runs) > 0 {
		t.Errorf("no overlay demo for %d runs", len(runs))
	}
}

func TestOverlayDemosLogMore(t *testing.T) {
	runs := map[string]struct {
		steps []string
		want  string
	}{
		"CalloutPopover": {[]string{"Before you publish"}, "open true"},
		"ConfirmPopover": {[]string{"Delete draft"}, "open true"},
		"RenameDialog":   {[]string{"Rename", "Cancel"}, "closed"},
		"PromptDialog":   {[]string{"Rename", "Cancel"}, "closed"},
	}
	for _, e := range Entries() {
		run, ok := runs[e.Name]
		if !ok || e.Category != "overlay" {
			continue
		}
		delete(runs, e.Name)
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		for _, s := range run.steps {
			if err := tt.Click(s); err != nil {
				t.Fatalf("%s: %v", e.Name, err)
			}
			tt.Advance(500 * time.Millisecond)
		}
		if len(logs) == 0 || logs[len(logs)-1] != run.want {
			t.Errorf("%s logged %q, want %q last", e.Name, logs, run.want)
		}
	}
	if len(runs) > 0 {
		t.Errorf("no overlay demo for %d runs", len(runs))
	}
}
