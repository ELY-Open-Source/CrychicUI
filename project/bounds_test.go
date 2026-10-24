package project

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// narrow lays view out 320 DIP wide, padded 10.
func narrow(view func(c *ui.Context)) func(c *ui.Context) {
	return padded(func(c *ui.Context) { ui.Column(c).Width(320).Children(func() { view(c) }) })
}

func TestTimeTrackerCarriesFractions(t *testing.T) {
	running, total := false, 1500*time.Millisecond
	tt := ui.NewTester(padded(func(c *ui.Context) { TimeTracker(c, &running, &total, TimeTrackerOptions{}) }), 500, 200)
	settles(t, tt)
	click(t, tt, "Start")
	tt.Advance(1500 * time.Millisecond)
	if !tt.HasText("0:00:03") {
		t.Errorf("1.5 s stored plus 1.5 s run: texts %q", tt.Texts())
	}
}

func TestTimelineShrinkNeverGrows(t *testing.T) {
	tasks := []TimelineTask{{Title: "Spike", Start: 0.1, End: 0.11}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Timeline(c, &tasks, TimelineOptions{}) }), 700, 220)
	settles(t, tt)
	r, _ := tt.Find("Spike")
	x, y := r.X+r.W+60, r.Y+r.H/2
	tt.Press(x, y)
	tt.Move(x-20, y)
	tt.Release(x-20, y)
	tt.Frame()
	if tasks[0].End > 0.11 {
		t.Errorf("a shrink grew the span to %v", tasks[0].End)
	}
}

func TestBoardNoOpDropReportsNothing(t *testing.T) {
	lanes := []Lane{{Title: "To do", Cards: []Task{{ID: "a", Title: "Alpha"}, {ID: "b", Title: "Beta"}}}}
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Board(c, &lanes, BoardOptions{}).Changed() {
			changes++
		}
	}), 500, 400)
	settles(t, tt)
	a, _ := tt.Find("Alpha")
	b, _ := tt.Find("Beta")
	before := baseline()
	drag(tt, a, b)
	if changes != 0 || lanes[0].Cards[0].ID != "a" {
		t.Errorf("%d changes, lanes %v", changes, lanes)
	}
	cues(t, before)
}

func TestLongLabelsStayInside(t *testing.T) {
	long := "Enterprise customer onboarding and account provisioning improvements"
	task := Task{ID: "a", Title: "Provisioning", Labels: []string{long}}
	tt := ui.NewTester(narrow(func(c *ui.Context) { TaskCard(c, task) }), 500, 300)
	settles(t, tt)
	if r, _ := tt.Find(long); r.X+r.W > 330 {
		t.Errorf("TaskCard label ends at %v, past 330", r.X+r.W)
	}
	on := []bool{false}
	tt = ui.NewTester(narrow(func(c *ui.Context) { Labels(c, []string{long}, &on) }), 500, 300)
	settles(t, tt)
	if r, _ := tt.Find(long); r.X+r.W > 330 {
		t.Errorf("Labels pill ends at %v, past 330", r.X+r.W)
	}
}

func TestBurndownAxisFits(t *testing.T) {
	var points []BurndownPoint
	for i := range 14 {
		points = append(points, BurndownPoint{Day: fmt.Sprintf("Oct %02d", i+1), Remaining: float64(40 - 2*i), Ideal: float64(40 - 3*i)})
	}
	tt := ui.NewTester(narrow(func(c *ui.Context) { Burndown(c, points, BurndownOptions{}) }), 500, 500)
	settles(t, tt)
	for _, p := range points {
		if r, ok := tt.Find(p.Day); ok && r.W > 0 && (r.X < 10 || r.X+r.W > 330) {
			t.Errorf("%s at %v leaves the card", p.Day, r)
		}
	}
	if !tt.HasText("Oct 14") || !tt.HasText("Oct 01") {
		t.Errorf("ends of the axis missing: %q", tt.Texts())
	}
}

func TestRoadmapCurrentAnimates(t *testing.T) {
	quarters := []Quarter{{Label: "Q3", Items: []string{"Trips"}}, {Label: "Q4", Items: []string{"Sync"}}}
	current := "Q3"
	tt := ui.NewTester(padded(func(c *ui.Context) { Roadmap(c, quarters, RoadmapOptions{Current: current}) }), 600, 300)
	settles(t, tt)
	r, _ := tt.Find("Q4")
	x, y := int(r.X)-4, int(r.Y+r.H/2)
	current = "Q4"
	tt.Frame()
	tt.Advance(16 * time.Millisecond)
	mid := tt.Image().RGBAAt(x, y)
	settles(t, tt)
	if end := tt.Image().RGBAAt(x, y); mid == end {
		t.Errorf("the current quarter's fill jumped to %v", end)
	}
}

func TestBoardLongLabelStaysInCard(t *testing.T) {
	long := "Enterprise customer onboarding and account provisioning improvements"
	lanes := []Lane{
		{Title: "To do", Cards: []Task{{ID: "a", Title: "Provisioning", Labels: []string{long}}}},
		{Title: "Doing"}, {Title: "Done"},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(520).Children(func() { Board(c, &lanes, BoardOptions{}) })
	}), 700, 400)
	settles(t, tt)
	label, _ := tt.Find(long)
	title, _ := tt.Find("Provisioning")
	if label.W == 0 || label.X+label.W > title.X+title.W+1 {
		t.Errorf("label %v passes its card's text edge %v", label, title.X+title.W)
	}
}

func TestBurndownLongLabelsStayApart(t *testing.T) {
	var points []BurndownPoint
	for i := range 14 {
		points = append(points, BurndownPoint{Day: fmt.Sprintf("Wednesday, Oct %02d", i+1), Remaining: float64(40 - 2*i), Ideal: float64(40 - 3*i)})
	}
	tt := ui.NewTester(narrow(func(c *ui.Context) { Burndown(c, points, BurndownOptions{}) }), 500, 500)
	settles(t, tt)
	end := 0.0
	for _, p := range points {
		r, ok := tt.Find(p.Day)
		if !ok || r.W == 0 {
			continue
		}
		if r.X < 10 || r.X+r.W > 330 || float64(r.X) < end {
			t.Errorf("%s at %v leaves the card or overlaps a label ending at %v", p.Day, r, end)
		}
		end = float64(r.X + r.W)
	}
}

func TestSprintHeaderCountsFromDayOne(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for _, tc := range []struct {
		start time.Time
		want  string
	}{
		{today, "Day 1 of 10"},
		{today.AddDate(0, 0, -1), "Day 2 of 10"},
		{today.AddDate(0, 0, -9), "Day 10 of 10"},
	} {
		end := tc.start.AddDate(0, 0, 10)
		tt := ui.NewTester(padded(func(c *ui.Context) { SprintHeader(c, "Sprint 14", tc.start, end, SprintHeaderOptions{}) }), 500, 260)
		tt.Frame()
		if !tt.HasText(tc.want) {
			t.Errorf("texts %q, want %q", tt.Texts(), tc.want)
		}
	}
}
