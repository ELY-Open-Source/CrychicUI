package feedback

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestEmptyState(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if _, act := EmptyState(c, "inbox", "No messages", "Mail you get shows here.", "Compose"); act.Clicked() {
			n++
		}
	}, 400, 400)
	settles(t, tt, "EmptyState rise")
	before := baseline()
	click(t, tt, "Compose")
	if n != 1 {
		t.Fatalf("clicks %d", n)
	}
	cues(t, before, sound.Tap)
	ui.NewTester(func(c *ui.Context) {
		if _, act := EmptyState(c, "inbox", "Nothing", "", ""); act.Valid() {
			t.Error("no action, yet an action element")
		}
	}, 200, 200)
}

func TestErrorState(t *testing.T) {
	before := baseline()
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if _, retry := ErrorState(c, "The server did not answer.", ErrorStateOptions{Retry: "Again"}); retry.Clicked() {
			retries++
		}
	}, 400, 400)
	find(t, tt, "Something went wrong")
	settles(t, tt, "shake")
	click(t, tt, "Again")
	if retries != 1 {
		t.Fatalf("retries %d", retries)
	}
	cues(t, before, sound.Error, sound.Tap)
}

func TestSuccessState(t *testing.T) {
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { SuccessState(c, "Payment done", "A receipt is on its way.", "View order") }, 400, 400)
	settles(t, tt, "check")
	click(t, tt, "View order")
	cues(t, before, sound.Success, sound.Tap)
}

func TestNotFound(t *testing.T) {
	homes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if _, home := NotFound(c, NotFoundOptions{Home: "Back"}); home.Clicked() {
			homes++
		}
	}, 400, 400)
	find(t, tt, "404")
	find(t, tt, "Page not found")
	settles(t, tt, "NotFound rise")
	before := baseline()
	click(t, tt, "Back")
	if homes != 1 {
		t.Fatalf("homes %d", homes)
	}
	cues(t, before, sound.Tap)
}
