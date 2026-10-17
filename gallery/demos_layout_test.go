package gallery

import (
	"testing"
	"time"
)

func TestScrollerDemoLogsArrows(t *testing.T) {
	e := Entries()[23]
	if e.Name != "Scroller" {
		t.Fatalf("entry 24 is %s", e.Name)
	}
	_, tt := open(e.Path())
	tt.Advance(time.Second)
	if err := tt.Click("Scroll forward"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Scroller: forward") {
		t.Errorf("log strip lacks the event: %q", tt.Texts())
	}
}
