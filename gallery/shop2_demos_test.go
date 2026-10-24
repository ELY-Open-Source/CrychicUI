package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestShop2DemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"OrderCard": {press("View order"), "opened order"},
		"Delivery":  {press("Express"), "selected 1"},
		"Upsell":    {press("Add"), "added 0"},
		"Bundle":    {press("Arc desk lamp"), "false"},
		"MiniCart":  {press("Cart"), "open true"},
		"Returns":   {press("Damaged"), "reason Damaged"},
		"Seller":    {press("Follow"), "followed true"},
		"QandA": {func(tt *ui.Tester) error {
			if err := tt.Click("Ask a question"); err != nil {
				return err
			}
			tt.Type("Dimmable?")
			return tt.Click("Ask")
		}, "1 questions"},
		"SizeGuide": {press("Size guide"), "open true"},
	}
	quiet := map[string]bool{"OrderStatus": true, "Tracking": true, "FreeShipping": true, "Recently": true,
		"FlashSale": true, "GiftCard": true, "Loyalty": true, "Receipt": true, "Specs": true, "StoreHours": true, "ThankYou": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "shop" || e.ID < 996 {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 20 {
		t.Errorf("checked %d second shop demos", seen)
	}
}
