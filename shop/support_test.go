package shop

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestReturns(t *testing.T) {
	selected := -1
	reason := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Returns(c, []Product{testProduct()}, "USD", []string{"Damaged", "Wrong size"}, &selected, &reason, ReturnsOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Desk lamp")
	click(t, tt, "Damaged")
	if selected != 0 || reason != "Damaged" {
		t.Errorf("selected %d reason %q", selected, reason)
	}
	cues(t, before, sound.Select, sound.Select)
}

func TestSeller(t *testing.T) {
	followed := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Seller(c, "芽衣商店", 4.8, &followed, SellerOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Follow")
	if !followed {
		t.Error("seller was not followed")
	}
	cues(t, before, sound.Toggle)
}

func TestQandA(t *testing.T) {
	questions := []Question{}
	draft := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		QandA(c, &questions, &draft, QandAOptions{})
	}), 500, 300)
	settles(t, tt)
	click(t, tt, "Ask a question")
	tt.Type("Is it dimmable?")
	before := baseline()
	click(t, tt, "Ask")
	if len(questions) != 1 || draft != "" {
		t.Errorf("questions %v draft %q", questions, draft)
	}
	cues(t, before, sound.Send)
}

func TestQandAFocus(t *testing.T) {
	questions := []Question{}
	draft := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		QandA(c, &questions, &draft, QandAOptions{})
	}), 500, 350)
	settles(t, tt)
	click(t, tt, "Ask a question")
	tt.Type("First")
	tt.Key(0, ui.KeyEnter)
	tt.Type("Second")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(questions) != 2 || questions[1].Text != "Second" {
		t.Errorf("questions %v", questions)
	}
}

func TestSizeGuide(t *testing.T) {
	open := false
	rows := []SizeRow{{Size: "M", Chest: "96 cm", Waist: "80 cm"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { SizeGuide(c, rows, &open, SizeGuideOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Size guide")
	if !open || !tt.HasText("96 cm") {
		t.Errorf("open %v texts %q", open, tt.Texts())
	}
	cues(t, before, sound.Open)
}

func TestSizeGuidePopover(t *testing.T) {
	open := false
	rows := []SizeRow{{Size: "S", Chest: "96 cm", Waist: "80 cm"}, {Size: "XXXL", Chest: "120 cm", Waist: "100 cm"}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SizeGuide(c, rows, &open, SizeGuideOptions{})
		ui.Text(c, "after")
	}), 500, 400)
	settles(t, tt)
	before, _ := tt.Find("after")
	click(t, tt, "Size guide")
	after, _ := tt.Find("after")
	if after.Y != before.Y {
		t.Errorf("popover moved content: before %v after %v", before.Y, after.Y)
	}
	head, _ := tt.Find("Chest")
	chest, _ := tt.Find("120 cm")
	waist, _ := tt.Find("100 cm")
	waistHead, _ := tt.Find("Waist")
	if chest.X != head.X || waist.X != waistHead.X {
		t.Errorf("columns not aligned: %v %v %v %v", head.X, chest.X, waistHead.X, waist.X)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open {
		t.Error("Escape left size guide open")
	}
}

func TestSpecs(t *testing.T) {
	groups := []SpecGroup{{Name: "Dimensions", Rows: []SpecRow{{Name: "Height", Value: "40 cm"}}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Specs(c, groups) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Dimensions") || !tt.HasText("40 cm") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSpecsMotion(t *testing.T) {
	groups := []SpecGroup{{Name: "Dimensions", Rows: []SpecRow{{Name: "Height", Value: "40 cm"}}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Specs(c, groups) }), 500, 180)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("specs have no arrival motion")
	}
}

func TestStoreHours(t *testing.T) {
	hours := []DayHours{
		{Day: "Sun"}, {Day: "Mon"}, {Day: "Tue"}, {Day: "Wed"}, {Day: "Thu"},
		{Day: "Fri", Open: "09:00", Close: "17:00"}, {Day: "Sat"},
	}
	now := time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC)
	tt := ui.NewTester(padded(func(c *ui.Context) { StoreHours(c, now, hours, StoreHoursOptions{}) }), 500, 350)
	settles(t, tt)
	if !tt.HasText("Open now") || !tt.HasText("09:00–17:00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestStoreHoursOvernight(t *testing.T) {
	hours := []DayHours{
		{Day: "Sun"}, {Day: "Mon"}, {Day: "Tue"}, {Day: "Wed"}, {Day: "Thu"},
		{Day: "Fri", Open: "22:00", Close: "02:00"}, {Day: "Sat"},
	}
	for _, now := range []time.Time{
		time.Date(2026, 10, 23, 23, 0, 0, 0, time.UTC),
		time.Date(2026, 10, 24, 1, 0, 0, 0, time.UTC),
	} {
		tt := ui.NewTester(padded(func(c *ui.Context) { StoreHours(c, now, hours, StoreHoursOptions{}) }), 500, 350)
		tt.Frame()
		if !tt.HasText("Open now") {
			t.Errorf("%v texts %q", now, tt.Texts())
		}
		if !tt.FramePending() {
			t.Error("store hours have no arrival motion")
		}
	}
}

func TestThankYou(t *testing.T) {
	before := baseline()
	tt := ui.NewTester(padded(func(c *ui.Context) { ThankYou(c, "1042", ThankYouOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Order 1042") {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Success)
}
