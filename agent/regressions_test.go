package agent

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestAgentBoardColumnsAndFlight(t *testing.T) {
	for _, settle := range []time.Duration{800 * time.Millisecond, 1500 * time.Millisecond} {
		agents := []BoardAgent{{Name: "Alpha", Status: Working}, {Name: "Beta", Status: Blocked}}
		tt := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).Width(500).Padding(10).Children(func() { AgentBoard(c, "Team", agents, AgentBoardOptions{}) })
		}, 600, 500)
		tt.Advance(settle)
		a, _ := tt.Find("Alpha")
		b, _ := tt.Find("Beta")
		if a.W > 250 || a.X+a.W > b.X {
			t.Fatalf("columns overlap: %v %v", a, b)
		}
		agents[0].Status = Blocked
		tt.Frame()
		tt.Advance(16 * time.Millisecond)
		if n := dark(tt, a); n < 20 {
			t.Errorf("after %v: card left its column at once: %d dark pixels where it was", settle, n)
		}
		tt.Advance(time.Second)
		if n := dark(tt, a); n > 0 {
			t.Errorf("after %v: card stayed: %d dark pixels", settle, n)
		}
	}
}

func TestSearchStepsLabelsLikeKeys(t *testing.T) {
	steps := []SearchStep{{Label: "sep.Reading"}, {Label: "Reading", State: StepRunning}}
	render(func(c *ui.Context) { SearchSteps(c, "Search", steps) })
	defer func() {
		if recover() == nil {
			t.Error("duplicate labels did not panic")
		}
	}()
	SearchSteps(nil, "Search", []SearchStep{{Label: "Reading"}, {Label: "Reading"}})
}

// dark counts near-black text pixels inside r.
func dark(tt *ui.Tester, r ui.Rect) int {
	img, n := tt.Image(), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); p.R < 80 && p.G < 80 && p.B < 80 {
				n++
			}
		}
	}
	return n
}

func TestContextFilesLongNameAndEmpty(t *testing.T) {
	files := []ContextFile{{Name: strings.Repeat("quarterly-revenue-", 6) + ".csv", Size: 10}}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(320).Padding(10).Children(func() { ContextFiles(c, "Files", &files, ContextFilesOptions{}) })
	}, 500, 300)
	tt.Advance(time.Second)
	r, ok := tt.Find("Remove " + files[0].Name)
	if !ok || r.X+r.W > 330 {
		t.Fatalf("remove pushed out: %v %v", r, ok)
	}
	click(t, tt, "Remove "+files[0].Name)
	tt.Advance(time.Second)
	if len(files) != 0 || !tt.HasText("No files attached") {
		t.Errorf("files %v texts %q", files, tt.Texts())
	}
}

func TestRateAnswerLongReasonFits(t *testing.T) {
	rating, reason, sent := -1, -1, false
	long := []string{strings.Repeat("The answer ignored the attached spreadsheet ", 3)}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(320).Padding(10).Children(func() { RateAnswer(c, &rating, long, &reason, &sent, RateAnswerOptions{}) })
	}, 500, 400)
	tt.Advance(time.Second)
	if r := find(t, tt, long[0]); r.X+r.W > 330 {
		t.Errorf("reason overflows to %v", r.X+r.W)
	}
}

func TestDiffKinds(t *testing.T) {
	lines := []string{"--- a/x.go", "+++ b/x.go", "@@ -1 +1 @@", "-----", "++counter;", " same"}
	if got := string(diffKinds(lines)); got != "   -+ " {
		t.Errorf("kinds %q", got)
	}
}

func TestAgentWelcomeSharedLabel(t *testing.T) {
	chosen := -1
	render(func(c *ui.Context) {
		AgentWelcome(c, "Hi", []Suggestion{{Text: "Research", Icon: "search"}}, []Capability{{Icon: "search", Title: "Research"}}, &chosen, AgentWelcomeOptions{})
	})
}
