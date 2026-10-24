package dashboard

import (
	"math"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestWidgetFrameMenuSurvivesNote(t *testing.T) {
	note, changes := "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if WidgetFrame(c, "Signups", WidgetFrameOptions{Note: note}, func() { ui.Text(c, "1,246") }).Changed() {
			changes++
		}
	}), 500, 300)
	settles(t, tt)
	click(t, tt, "Widget actions")
	note = "Updated 2m ago"
	tt.Frame()
	if err := tt.ChooseMenuItem("Refresh"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if changes != 1 {
		t.Errorf("refresh after a note: %d changes, want 1", changes)
	}
}

func TestMRRNetIsExact(t *testing.T) {
	for _, tc := range []struct {
		new, expansion, churn float64
		want                  string
	}{
		{1 << 53, 1, 1 << 53, "Net new +USD 1.00"},
		{1e308, 1e308, 1e308, "Net new +USD"},
	} {
		tt := ui.NewTester(padded(func(c *ui.Context) { MRR(c, 1000, tc.new, tc.expansion, tc.churn, "USD", MRROptions{}) }), 500, 400)
		tt.Frame()
		found := false
		for _, s := range tt.Texts() {
			if strings.HasPrefix(s, tc.want) {
				found = true
			}
		}
		if !found {
			t.Errorf("MRR %g+%g-%g: texts %q, want %q", tc.new, tc.expansion, tc.churn, tt.Texts(), tc.want)
		}
	}
	panicsSaying(t, `"MRR" net overflows`, func(c *ui.Context) { MRR(c, 0, math.MaxFloat64, math.MaxFloat64, 0, "USD", MRROptions{}) })
}

func TestWidgetFramePanicNamesTitle(t *testing.T) {
	panicsSaying(t, `"Revenue"`, func(c *ui.Context) { WidgetFrame(c, "Revenue", WidgetFrameOptions{}, nil) })
}
