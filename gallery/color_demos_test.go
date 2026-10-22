package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestColorDemosLog(t *testing.T) {
	click := func(labels ...string) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, l := range labels {
				if err := tt.Click(l); err != nil {
					return err
				}
				tt.Advance(300 * time.Millisecond)
			}
			return nil
		}
	}
	keys := func(ks ...ui.Key) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, k := range ks {
				tt.Key(0, k)
			}
			return nil
		}
	}
	bump := func(field string) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			if err := click(field)(tt); err != nil {
				return err
			}
			tt.Key(0, ui.KeyUp)
			return nil
		}
	}
	pick := func(open string) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			if err := click(open, "Hex")(tt); err != nil {
				return err
			}
			tt.Command("selectAll")
			tt.Type("00FF00")
			tt.Key(0, ui.KeyEnter)
			return nil
		}
	}
	acts := map[string]func(*ui.Tester) error{
		"Picker": pick("Brand"),
		"Well":   pick("Accent"),
		"HexField": func(tt *ui.Tester) error {
			tt.Key(0, ui.KeyTab)
			tt.Command("selectAll")
			tt.Type("#00FF00")
			return nil
		},
		"RGBFields":      bump("R"),
		"HSLFields":      bump("H"),
		"OklchFields":    bump("L"),
		"HueRing":        keys(ui.KeyTab, ui.KeyRight),
		"ShadePad":       keys(ui.KeyTab, ui.KeyRight),
		"AlphaSlider":    keys(ui.KeyTab, ui.KeyLeft),
		"SwatchGrid":     click("#22C55E"),
		"Contrast":       click("Darker text"),
		"GradientEditor": click("Fill"),
		"RecentColors":   click("Choose another"),
		"ColorName":      click("Next color"),
		"Harmony":        keys(ui.KeyTab, ui.KeySpace),
		"ThemePreview":   click("Switch theme"),
		"Eyedropper":     click("Strip"),
	}
	for _, e := range Entries() {
		act, ok := acts[e.Name]
		if !ok || e.Category != "color" {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 600)
		tt.Advance(time.Second)
		if err := act(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if len(logs) == 0 {
			t.Errorf("%s: the action logged nothing", e.Name)
		}
	}
}
