package gallery

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestFormDemosLogEdits(t *testing.T) {
	for _, e := range Entries() {
		switch e.Name {
		case "FormCard", "InlineForm", "ContactForm", "LoginBox", "NewsletterForm", "ProfileForm", "PaymentForm":
		default:
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 600, 700)
		tt.Advance(0)
		tt.Key(0, ui.KeyTab)
		tt.Type("4")
		tt.Frame()
		if len(logs) == 0 {
			t.Errorf("%s: typing logged nothing", e.Name)
		}
	}
}

func TestSurveyDemoLogsBack(t *testing.T) {
	for _, e := range Entries() {
		if e.Name != "Survey" {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 600, 500)
		if err := tt.Click("Daily"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(time.Second)
		if err := tt.Click("Back"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		if !slices.Contains(logs, "back") {
			t.Errorf("logs %q", logs)
		}
	}
}
