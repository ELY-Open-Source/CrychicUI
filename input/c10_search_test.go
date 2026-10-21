package input

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestFilterSearch(t *testing.T) {
	query, selected := "", []string{}
	tt := ui.NewTester(func(c *ui.Context) {
		FilterSearch(c, &query, []string{"UI", "Docs"}, &selected, FilterSearchOptions{})
	}, 400, 240)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("motion")
	if query != "motion" || !appended(before, sound.Tick) {
		t.Fatalf("query %q", query)
	}
	if err := tt.Click("UI"); err != nil {
		t.Fatal(err)
	}
	if len(selected) != 1 || selected[0] != "UI" || !appended(before, sound.Toggle) {
		t.Errorf("selected %q", selected)
	}
}

func TestSearchHistory(t *testing.T) {
	query, history := "", []string{}
	tt := ui.NewTester(func(c *ui.Context) { SearchHistory(c, &query, &history, SearchHistoryOptions{}) }, 400, 300)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("first")
	tt.Key(0, ui.KeyEnter)
	if len(history) != 1 || history[0] != "first" {
		t.Fatalf("history %q", history)
	}
	query = ""
	tt.Frame()
	if err := tt.Click("Use first"); err != nil {
		t.Fatal(err)
	}
	if query != "first" || !appended(before, sound.Select) {
		t.Errorf("query %q", query)
	}
	if err := tt.Click("Remove first"); err != nil {
		t.Fatal(err)
	}
	if len(history) != 0 || !appended(before, sound.Delete) {
		t.Errorf("history %q", history)
	}
}

func TestDebouncedSearch(t *testing.T) {
	query, reported, changed := "", "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if DebouncedSearch(c, "Search", &query, &reported, DebouncedSearchOptions{Delay: 300 * time.Millisecond}).Changed() {
			changed++
		}
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("query")
	if reported != "" || !tt.HasText("Searching") {
		t.Fatal("reported before the pause")
	}
	tt.Advance(300 * time.Millisecond)
	if reported != "query" || changed == 0 || tt.HasText("Searching") {
		t.Errorf("reported %q, changed %d, texts %q", reported, changed, tt.Texts())
	}
}

func TestVerifyField(t *testing.T) {
	code := ""
	tt := ui.NewTester(func(c *ui.Context) { VerifyField(c, "Code", &code, VerifyFieldOptions{}) }, 400, 240)
	before := baseline()
	if err := tt.Click("Send code"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Resend in") || !appended(before, sound.Send) {
		t.Fatal("send did not count down")
	}
	tt.Advance(3 * time.Second)
	if !tt.HasText("Send code") {
		t.Error("resend did not re-enable")
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("1234")
	if code != "1234" {
		t.Errorf("code %q", code)
	}
}

func TestSignaturePad(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SignaturePad(c, SignaturePadOptions{}).Changed() {
			changed++
		}
	}, 400, 300)
	before := baseline()
	tt.Press(20, 20)
	tt.Move(80, 60)
	tt.Release(80, 60)
	if changed != 1 || !appended(before, sound.Tick) {
		t.Fatalf("changed %d, cues %v", changed, sound.Recent())
	}
	if err := tt.Click("Clear signature"); err != nil {
		t.Fatal(err)
	}
	if changed != 2 || !appended(before, sound.Close) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
}
