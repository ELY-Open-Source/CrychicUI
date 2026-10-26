package choice

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestOption(t *testing.T) {
	plan, changes := "free", 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			for _, v := range []string{"free", "pro"} {
				label := map[string]string{"free": "Free plan", "pro": "Pro plan"}[v]
				if Option(c, label, v, &plan).Changed() {
					changes++
				}
			}
		})
	}, 300, 200)
	before := baseline()
	click(t, tt, "Pro plan")
	if plan != "pro" || changes != 1 {
		t.Errorf("pick gave %q after %d changes", plan, changes)
	}
	cues(t, before, sound.Select)
	settles(t, tt, "Option dot")
	before = baseline()
	click(t, tt, "Pro plan")
	if plan != "pro" || changes != 1 {
		t.Errorf("picking the picked option gave %q after %d changes", plan, changes)
	}
	cues(t, before)
	func() {
		defer func() {
			if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), `"Free plan"`) {
				t.Errorf("nil selection panic %v, want the label", r)
			}
		}()
		ui.Render(func(c *ui.Context) { Option(c, "Free plan", "free", nil) }, 100, 100, 1)
	}()
	panics(t, "blank label", func() { ui.Render(func(c *ui.Context) { Option(c, "", "free", &plan) }, 100, 100, 1) })
}
