package project

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestBoard(t *testing.T) {
	lanes := []Lane{{Title: "Todo", Cards: []Task{{ID: "t1", Title: "Task"}}}, {Title: "Doing"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Board(c, &lanes, BoardOptions{}) }), 700, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add card")
	if len(lanes[0].Cards) != 2 {
		t.Errorf("cards %v", lanes[0].Cards)
	}
	r, _ := tt.Find("Task")
	zone, _ := tt.Find("Doing")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(zone.X+zone.W/2, zone.Y+zone.H+20)
	tt.Release(zone.X+zone.W/2, zone.Y+zone.H+20)
	tt.Frame()
	if len(lanes[0].Cards) != 1 || len(lanes[1].Cards) != 1 || lanes[1].Cards[0].ID != "t1" {
		t.Errorf("lanes %v", lanes)
	}
	cues(t, before, sound.Pop, sound.Drop)
}

func TestTaskCard(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { TaskCard(c, testTask()) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Ship UI") || !tt.HasText("Friday") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTaskRow(t *testing.T) {
	done := false
	tt := ui.NewTester(padded(func(c *ui.Context) { TaskRow(c, testTask(), &done) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Ship UI")
	if !done {
		t.Error("task not checked")
	}
	cues(t, before, sound.Toggle)
}

func TestSubtasks(t *testing.T) {
	items := []Subtask{{Title: "One"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Subtasks(c, &items, SubtasksOptions{}) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "One")
	if !items[0].Done || !tt.HasText("1 of 1") {
		t.Errorf("items %v texts %q", items, tt.Texts())
	}
	cues(t, before, sound.Toggle)
}

func TestPriority(t *testing.T) {
	level := Medium
	tt := ui.NewTester(padded(func(c *ui.Context) { Priority(c, &level, PriorityOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Urgent")
	if level != Urgent {
		t.Errorf("level %v", level)
	}
	cues(t, before, sound.Select)
}

func TestStatus(t *testing.T) {
	status := Todo
	tt := ui.NewTester(padded(func(c *ui.Context) { Status(c, &status, StatusOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Done")
	if status != Done {
		t.Errorf("status %v", status)
	}
	cues(t, before, sound.Select)
}

func TestAssignee(t *testing.T) {
	selected := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Assignee(c, []string{"Ana Lima", "Kai Moreno"}, &selected, AssigneeOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Kai Moreno")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}

func TestLabels(t *testing.T) {
	selected := []bool{false}
	tt := ui.NewTester(padded(func(c *ui.Context) { Labels(c, []string{"UI"}, &selected) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "UI")
	if !selected[0] {
		t.Errorf("selected %v", selected)
	}
	cues(t, before, sound.Toggle)
}

func TestEstimate(t *testing.T) {
	points := 3
	tt := ui.NewTester(padded(func(c *ui.Context) { Estimate(c, &points, EstimateOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "5")
	if points != 5 {
		t.Errorf("points %d", points)
	}
	cues(t, before, sound.Select)
}

func TestSprintHeader(t *testing.T) {
	start := time.Now()
	end := start.AddDate(0, 0, 10)
	tt := ui.NewTester(padded(func(c *ui.Context) { SprintHeader(c, "Sprint 1", start, end, SprintHeaderOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Sprint 1") || !tt.HasText("10 days left") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestBurndown(t *testing.T) {
	points := []BurndownPoint{{Day: "Mon", Remaining: 10, Ideal: 10}, {Day: "Tue", Remaining: 5, Ideal: 5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Burndown(c, points, BurndownOptions{}) }), 500, 250)
	settles(t, tt)
}

func TestVelocity(t *testing.T) {
	sprints := []Sprint{{Name: "S1", Points: 20}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Velocity(c, sprints, VelocityOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("S1") || !tt.HasText("20") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTimeline(t *testing.T) {
	tasks := []TimelineTask{{Title: "Task", Start: 0.1, End: 0.5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Timeline(c, &tasks, TimelineOptions{}) }), 700, 220)
	settles(t, tt)
	r, _ := tt.Find("Task")
	before := baseline()
	x, y := r.X+r.W+50, r.Y+r.H/2
	tt.Press(x, y)
	tt.Move(x+100, y)
	tt.Release(x+100, y)
	tt.Frame()
	if tasks[0].End <= 0.5 {
		t.Errorf("end %v", tasks[0].End)
	}
	cues(t, before, sound.Tick)
}

func TestRoadmap(t *testing.T) {
	quarters := []Quarter{{Label: "Q1", Items: []string{"UI refresh"}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Roadmap(c, quarters, RoadmapOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Q1") || !tt.HasText("UI refresh") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMilestone(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Milestone(c, "Launch", "Nov 1", 0.5, MilestoneOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Launch") || !tt.HasText("Due Nov 1") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestIssueCard(t *testing.T) {
	issue := Issue{Key: "UI-1", Title: "Ship UI", Type: "Story", Priority: High}
	tt := ui.NewTester(padded(func(c *ui.Context) { IssueCard(c, issue, IssueCardOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("UI-1") || !tt.HasText("Story") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestIssueHeader(t *testing.T) {
	action := IssueNone
	issue := Issue{Key: "UI-1", Title: "Ship UI", Type: "Story", Priority: High}
	tt := ui.NewTester(padded(func(c *ui.Context) { IssueHeader(c, issue, Todo, &action, IssueHeaderOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Close")
	if action != IssueClose {
		t.Errorf("action %v", action)
	}
	cues(t, before, sound.Tap)
}

func TestComments(t *testing.T) {
	comments := []Comment{}
	draft := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Comments(c, &comments, &draft, CommentsOptions{})
	}), 500, 300)
	settles(t, tt)
	click(t, tt, "Write a comment")
	tt.Type("Looks good")
	before := baseline()
	click(t, tt, "Send")
	if len(comments) != 1 || draft != "" {
		t.Errorf("comments %v draft %q", comments, draft)
	}
	cues(t, before, sound.Send)
}

func TestTimeTracker(t *testing.T) {
	running, total := false, time.Duration(0)
	tt := ui.NewTester(padded(func(c *ui.Context) { TimeTracker(c, &running, &total, TimeTrackerOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Start")
	tt.Advance(100 * time.Millisecond)
	click(t, tt, "Stop")
	if running || total <= 0 {
		t.Errorf("running %v total %v", running, total)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
}

func TestBacklog(t *testing.T) {
	issues := []Issue{{Key: "A-1", Title: "Alpha", Type: "Task"}, {Key: "A-2", Title: "Beta", Type: "Bug"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Backlog(c, &issues, BacklogOptions{}) }), 500, 300)
	settles(t, tt)
	a, _ := tt.Find("Alpha")
	b, _ := tt.Find("Beta")
	before := baseline()
	tt.Press(b.X+b.W/2, b.Y+b.H/2)
	tt.Move(a.X+a.W/2, a.Y+a.H/2)
	tt.Release(a.X+a.W/2, a.Y+a.H/2)
	tt.Frame()
	if issues[0].Title != "Beta" {
		t.Errorf("issues %v", issues)
	}
	cues(t, before, sound.Drop)
}
