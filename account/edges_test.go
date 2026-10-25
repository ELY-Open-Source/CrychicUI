package account

import (
	"math"
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"
)

// narrow puts a view in a 320-DIP column.
func narrow(view func(c *ui.Context)) func(c *ui.Context) {
	return padded(func(c *ui.Context) { ui.Column(c).Width(320).Children(func() { view(c) }) })
}

func TestKeysStayInside(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("a key escaped its component: %v", r)
		}
	}()
	ui.Render(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			for i := range 2 {
				c.Key(i)
				Providers(c, []string{"Google"}, new(string), ProvidersOptions{})
				ui.Text(c, "Footer")
			}
			for i := range 2 {
				c.Key(10 + i)
				SignUp(c, &SignUpValue{}, SignUpOptions{})
				ui.Text(c, "Footer")
			}
		})
	}, 600, 1200, 1)
}

func TestUsageHugeValues(t *testing.T) {
	r := []Resource{{Name: "Storage", Unit: "GB", Used: math.MaxFloat64, Limit: math.MaxFloat64}}
	tt := ui.NewTester(narrow(func(c *ui.Context) { Usage(c, r, UsageOptions{}) }), 500, 400)
	settles(t, tt)
	if name, ok := tt.Find("Storage"); !ok || name.W < 40 {
		t.Errorf("the resource name shrank away: %+v", name)
	}
	big := strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64)
	v, ok := tt.Find(big + " of " + big + " GB")
	if !ok || v.X+v.W > 330.5 {
		t.Errorf("value %v spans to %.1f past the 330 card edge", ok, v.X+v.W)
	}
	panics(t, `Usage "Storage" percent overflows`, func(c *ui.Context) {
		Usage(c, []Resource{{Name: "Storage", Unit: "GB", Used: 1e308, Limit: 1}}, UsageOptions{})
	})
}

func TestLongMetaLeavesActions(t *testing.T) {
	long := "October 21, 2026 at 10:45 AM"
	sessions := []Session{{ID: "a", Device: "iPad Air", Location: "Seattle, WA", Last: long}}
	tt := ui.NewTester(narrow(func(c *ui.Context) { Sessions(c, &sessions, SessionsOptions{}) }), 500, 400)
	settles(t, tt)
	clear(t, tt, long, "Sign out iPad Air")

	created := "Created " + long
	keys := []APIKey{{ID: "s", Name: "Staging", Prefix: "ck_test_", Secret: "2", Created: created}}
	tt = ui.NewTester(narrow(func(c *ui.Context) { APIKeys(c, &keys, APIKeysOptions{}) }), 500, 400)
	settles(t, tt)
	clear(t, tt, created, "Copy Staging")
}

// clear fails when text reaches into the action's box.
func clear(t *testing.T, tt *ui.Tester, text, action string) {
	t.Helper()
	m, ok := tt.Find(text)
	b, found := tt.Find(action)
	if !ok || !found {
		t.Fatalf("%q or %q missing: %q", text, action, tt.Texts())
	}
	if m.X+m.W > b.X+0.5 {
		t.Errorf("%q ends at %.1f, over %q at %.1f", text, m.X+m.W, action, b.X)
	}
}

func TestPlanBorderAnimates(t *testing.T) {
	current, change := false, false
	plan := Plan{Name: "Studio", Price: "$24"}
	tt := ui.NewTester(padded(func(c *ui.Context) { PlanCard(c, plan, current, &change, PlanCardOptions{}) }), 500, 400)
	settles(t, tt)
	current = true
	tt.Frame()
	if !tt.FramePending() {
		t.Error("current border snapped instead of animating")
	}
}
