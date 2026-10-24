package project

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestTimelineStationaryTrack(t *testing.T) {
	tasks := []TimelineTask{{Title: "Task", Start: 0.1, End: 0.11}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Timeline(c, &tasks) }), 700, 220)
	settles(t, tt)
	r, _ := tt.Find("Task")
	before := baseline()
	x, y := r.X+r.W+50, r.Y+r.H/2
	tt.Press(x, y)
	tt.Release(x, y)
	tt.Frame()
	if tasks[0].End != 0.11 {
		t.Errorf("stationary press resized to %v", tasks[0].End)
	}
	cues(t, before)
}

func TestBoardDuplicateSourceIndex(t *testing.T) {
	lanes := []Lane{{Title: "A", Cards: []string{"Fix", "Other", "Fix"}}, {Title: "B"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Board(c, &lanes, BoardOptions{}) }), 700, 400)
	settles(t, tt)
	var r ui.Rect
	for _, text := range tt.Texts() {
		_ = text
	}
	first, _ := tt.Find("Fix")
	other, _ := tt.Find("Other")
	r = other
	r.Y = other.Y + (other.Y - first.Y)
	zone, _ := tt.Find("B")
	t.Logf("first %v third %v zone %v", first, r, zone)
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(zone.X+zone.W/2, zone.Y+zone.H+20)
	tt.Release(zone.X+zone.W/2, zone.Y+zone.H+20)
	tt.Frame()
	if len(lanes[0].Cards) != 2 || lanes[0].Cards[0] != "Fix" || lanes[0].Cards[1] != "Other" {
		t.Errorf("wrong duplicate removed: %v", lanes[0].Cards)
	}
}

func TestSprintHeaderNextChange(t *testing.T) {
	start := time.Now()
	end := start.Add(48 * time.Hour)
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SprintHeader(c, "Sprint 1", start, end, SprintHeaderOptions{})
	}), 500, 180)
	settles(t, tt)
	tt.Advance(4 * time.Hour)
	if !tt.HasText("1 days left") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTimeTrackerResume(t *testing.T) {
	running, total := false, time.Duration(0)
	tt := ui.NewTester(padded(func(c *ui.Context) { TimeTracker(c, &running, &total, TimeTrackerOptions{}) }), 500, 180)
	settles(t, tt)
	click(t, tt, "Start")
	tt.Advance(100 * time.Millisecond)
	click(t, tt, "Stop")
	previous := total
	tt.Advance(time.Second)
	running = true
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	click(t, tt, "Stop")
	if total > previous+time.Second {
		t.Errorf("resumed total %v previous %v", total, previous)
	}
}
