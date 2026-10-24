package project

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestBoardSameColumnAndDuplicates(t *testing.T) {
	lanes := []Lane{{Title: "Todo", Cards: []string{"Fix"}}, {Title: "Doing", Cards: []string{"Fix"}}, {Title: "Review"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Board(c, &lanes, BoardOptions{}) }), 700, 400)
	settles(t, tt)
	click(t, tt, "Add card")
	click(t, tt, "Add card")
	r, _ := tt.Find("Fix")
	zone, _ := tt.Find("Todo")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(zone.X+zone.W/2, zone.Y+zone.H+20)
	tt.Release(zone.X+zone.W/2, zone.Y+zone.H+20)
	tt.Frame()
	if len(lanes[0].Cards) != 3 || lanes[1].Cards[0] != "Fix" {
		t.Errorf("same-column state %v", lanes)
	}
}

func TestBoardSourceIdentity(t *testing.T) {
	lanes := []Lane{{Title: "A", Cards: []string{"Fix"}}, {Title: "B", Cards: []string{"Fix"}}, {Title: "C"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Board(c, &lanes, BoardOptions{}) }), 700, 400)
	settles(t, tt)
	r, _ := tt.Find("Fix")
	zone, _ := tt.Find("C")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(zone.X+zone.W/2, zone.Y+zone.H+20)
	tt.Release(zone.X+zone.W/2, zone.Y+zone.H+20)
	tt.Frame()
	if len(lanes[0].Cards) != 0 || len(lanes[1].Cards) != 1 || len(lanes[2].Cards) != 1 {
		t.Errorf("wrong source moved: %v", lanes)
	}
}

func TestTimelineClickNoChange(t *testing.T) {
	tasks := []TimelineTask{{Title: "Task", Start: 0.1, End: 0.11}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Timeline(c, &tasks) }), 700, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Task")
	if tasks[0].End != 0.11 {
		t.Errorf("click resized to %v", tasks[0].End)
	}
	cues(t, before)
}

func TestTimeTrackerInitialRunning(t *testing.T) {
	running, total := true, time.Duration(0)
	tt := ui.NewTester(padded(func(c *ui.Context) { TimeTracker(c, &running, &total, TimeTrackerOptions{}) }), 500, 180)
	settles(t, tt)
	click(t, tt, "Stop")
	if total > 3*time.Second {
		t.Errorf("initial running total %v", total)
	}
}

func TestTimeTrackerOverflow(t *testing.T) {
	running, total := false, time.Duration(math.MaxInt64-1)
	tt := ui.NewTester(padded(func(c *ui.Context) { TimeTracker(c, &running, &total, TimeTrackerOptions{}) }), 500, 180)
	settles(t, tt)
	click(t, tt, "Start")
	tt.Advance(2 * time.Millisecond)
	defer func() {
		if recover() == nil {
			t.Error("expected overflow panic")
		}
	}()
	click(t, tt, "Stop")
}

func TestAssigneeUnicode(t *testing.T) {
	selected := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Assignee(c, []string{"张三"}, &selected) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("张") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestIssueCardPriority(t *testing.T) {
	issue := Issue{Key: "UI-1", Title: "Ship UI", Type: "Story", Priority: Urgent}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		IssueCard(c, issue, IssueCardOptions{Urgent: "Critical"})
	}), 500, 180)
	settles(t, tt)
	if !tt.HasText("Critical") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestIssueHeaderStatus(t *testing.T) {
	action := IssueNone
	issue := Issue{Key: "UI-1", Title: "Ship UI", Type: "Story", Priority: High}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		IssueHeader(c, issue, Doing, &action, IssueHeaderOptions{Doing: "In progress"})
	}), 500, 180)
	settles(t, tt)
	if !tt.HasText("In progress") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTaskRowUnchecked(t *testing.T) {
	done := false
	tt := ui.NewTester(padded(func(c *ui.Context) { TaskRow(c, testTask(), &done) }), 500, 180)
	settles(t, tt)
	if tt.HasText("check") {
		t.Errorf("unchecked row has checkmark: %q", tt.Texts())
	}
}

func TestBurndownRoadmapMotion(t *testing.T) {
	points := []BurndownPoint{{Day: "Mon", Remaining: 10, Ideal: 10}, {Day: "Tue", Remaining: 5, Ideal: 5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Burndown(c, points) }), 500, 250)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("burndown has no motion")
	}
	quarters := []Quarter{{Label: "Q1", Items: []string{"UI refresh"}}}
	other := ui.NewTester(padded(func(c *ui.Context) { Roadmap(c, quarters) }), 500, 220)
	other.Frame()
	if !other.FramePending() {
		t.Error("roadmap has no motion")
	}
}

func TestCopyOverrides(t *testing.T) {
	comments := []Comment{}
	draft := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Comments(c, &comments, &draft, CommentsOptions{Author: "Reviewer"})
	}), 500, 300)
	settles(t, tt)
	click(t, tt, "Write a comment")
	tt.Type("Looks good")
	click(t, tt, "Send")
	if comments[0].Author != "Reviewer" {
		t.Errorf("author %q", comments[0].Author)
	}
	start := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 10)
	header := ui.NewTester(padded(func(c *ui.Context) {
		SprintHeader(c, "Sprint 1", start, end, SprintHeaderOptions{DateFormat: "2006-01-02"})
	}), 500, 180)
	if !header.HasText("2026-10-01") || !header.HasText("2026-10-11") {
		t.Errorf("date format %q", header.Texts())
	}
}

func TestEstimateValidation(t *testing.T) {
	points := -3
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	ui.Render(func(c *ui.Context) { Estimate(c, &points, EstimateOptions{}) }, 500, 180, 1)
}
