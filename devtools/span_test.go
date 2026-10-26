package devtools

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestSpanDetail(t *testing.T) {
	span := Span{ID: "s2", Name: "POST /orders", Service: "kiln-api", Duration: 128 * time.Millisecond, Error: true}
	attrs := []Attr{{"http.status_code", "502"}, {"db.system", "postgres"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { SpanDetail(c, span, attrs, SpanDetailOptions{}) }), 480, 400)
	if !tt.FramePending() {
		t.Error("attribute rows did not arrive")
	}
	settles(t, tt)
	for _, s := range []string{"POST /orders", "kiln-api", "128 ms", "Error", "http.status_code", "502", "postgres"} {
		if !exact(tt, s) {
			t.Errorf("%q missing from %q", s, tt.Texts())
		}
	}
	if exact(tt, "OK") || exact(tt, "No attributes") {
		t.Errorf("a failed span with attributes shows %q", tt.Texts())
	}
	span.Error = false
	tt = ui.NewTester(padded(func(c *ui.Context) { SpanDetail(c, span, nil, SpanDetailOptions{}) }), 480, 400)
	settles(t, tt)
	if !exact(tt, "OK") || !exact(tt, "No attributes") {
		t.Errorf("an ok span without attributes shows %q", tt.Texts())
	}
	panics(t, "SpanDetail item", func(c *ui.Context) {
		SpanDetail(c, span, []Attr{{"db.system", "a"}, {"db.system", "b"}}, SpanDetailOptions{})
	})
	panics(t, "needs a name", func(c *ui.Context) { SpanDetail(c, Span{ID: "x", Service: "api"}, nil, SpanDetailOptions{}) })
}
