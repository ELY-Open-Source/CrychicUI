package project

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// panics wants a panic whose message holds want.
func panics(t *testing.T, want string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want %q", r, want)
		}
	}()
	ui.Render(view, 600, 500, 1)
}

// swapUnderPress presses name, reorders the data, then releases there.
func swapUnderPress(t *testing.T, tt *ui.Tester, name string, reorder func()) {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("%q missing", name)
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	reorder()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func drag(tt *ui.Tester, from, to ui.Rect) {
	tt.Press(from.X+from.W/2, from.Y+from.H/2)
	tt.Move(to.X+to.W/2, to.Y+to.H/2)
	tt.Release(to.X+to.W/2, to.Y+to.H/2)
	tt.Frame()
}

func TestBoardDropsBeforeACard(t *testing.T) {
	lanes := []Lane{{Title: "To do", Cards: []Task{{ID: "a", Title: "Alpha"}}}, {Title: "Doing", Cards: []Task{{ID: "b", Title: "Beta"}}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Board(c, &lanes, BoardOptions{}) }), 700, 500)
	settles(t, tt)
	a, _ := tt.Find("Alpha")
	b, _ := tt.Find("Beta")
	before := baseline()
	drag(tt, a, b)
	if len(lanes[0].Cards) != 0 || len(lanes[1].Cards) != 2 || lanes[1].Cards[0].ID != "a" {
		t.Errorf("lanes %v", lanes)
	}
	cues(t, before, sound.Drop)
	if !tt.HasText("Drop cards here") {
		t.Errorf("an empty lane shows no hint: %q", tt.Texts())
	}
}

func TestBoardChecks(t *testing.T) {
	lanes := []Lane{{Title: "To do", Cards: []Task{{ID: "a", Title: "A"}, {ID: "a", Title: "B"}}}}
	panics(t, `repeated "a"`, func(c *ui.Context) { Board(c, &lanes, BoardOptions{}) })
	lanes = []Lane{{Title: "To do"}, {Title: "To do"}}
	panics(t, `repeated "To do"`, func(c *ui.Context) { Board(c, &lanes, BoardOptions{}) })
	lanes = []Lane{{Title: "To do", Cards: []Task{{ID: "new-1", Title: "Taken"}}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Board(c, &lanes, BoardOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Add card")
	if got := lanes[0].Cards[1].ID; got == "new-1" {
		t.Errorf("added card reused ID %q", got)
	}
}

func TestPickersKeepIdentityAcrossReorder(t *testing.T) {
	people := []string{"Maya Chen", "Kai Moreno"}
	selected := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Assignee(c, people, &selected, AssigneeOptions{}) }), 500, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Maya Chen", func() { people[0], people[1] = people[1], people[0] })
	if selected != -1 {
		t.Errorf("Assignee picked %d after a reorder", selected)
	}
	labels, on := []string{"Mobile", "Backend"}, []bool{false, false}
	tt = ui.NewTester(padded(func(c *ui.Context) { Labels(c, labels, &on) }), 500, 200)
	settles(t, tt)
	swapUnderPress(t, tt, "Mobile", func() { labels[0], labels[1] = labels[1], labels[0] })
	if slices.Contains(on, true) {
		t.Errorf("Labels toggled a moved label: %v", on)
	}
	items := []Subtask{{Title: "Cache"}, {Title: "Queue"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Subtasks(c, &items, SubtasksOptions{}) }), 500, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Cache", func() { items[0], items[1] = items[1], items[0] })
	if items[0].Done || items[1].Done {
		t.Errorf("Subtasks toggled a moved item: %v", items)
	}
}

func TestPickerChecks(t *testing.T) {
	selected := -1
	panics(t, `repeated "Ana"`, func(c *ui.Context) { Assignee(c, []string{"Ana", "Ana"}, &selected, AssigneeOptions{}) })
	points := 0
	panics(t, "option 3", func(c *ui.Context) { Estimate(c, &points, EstimateOptions{Points: []int{3, 3}}) })
	items := []Subtask{{Title: "A"}, {Title: "A"}}
	panics(t, `repeated "A"`, func(c *ui.Context) { Subtasks(c, &items, SubtasksOptions{}) })
}

func TestSprintHeaderStates(t *testing.T) {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	far := time.Date(2500, 1, 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)
	cases := []struct {
		start, end time.Time
		want       string
	}{
		{now.Add(-9 * 24 * time.Hour), midnight.Add(-time.Second), "Ends today"},
		{now.Add(-9 * 24 * time.Hour), midnight.Add(time.Hour), "1 day left"},
		{now.Add(-10 * 24 * time.Hour), now.Add(-time.Hour), "Ended"},
		{now, far, fmt.Sprintf("%d days left", (far.Unix()-today.Unix())/86400)},
	}
	for _, tc := range cases {
		tt := ui.NewTester(padded(func(c *ui.Context) { SprintHeader(c, "Sprint 14", tc.start, tc.end, SprintHeaderOptions{}) }), 500, 260)
		tt.Frame()
		if !tt.HasText(tc.want) {
			t.Errorf("texts %q, want %q", tt.Texts(), tc.want)
		}
	}
}

func TestChartTotals(t *testing.T) {
	points := []BurndownPoint{{Day: "D1", Remaining: 10, Ideal: 10}, {Day: "D2", Remaining: 2.5, Ideal: 5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Burndown(c, points, BurndownOptions{}) }), 500, 400)
	tt.Frame()
	if !tt.HasText("2.5 pts left") {
		t.Errorf("texts %q", tt.Texts())
	}
	sprints := []Sprint{{Name: "S1", Points: 30}, {Name: "S2", Points: 41}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Velocity(c, sprints, VelocityOptions{}) }), 500, 400)
	tt.Frame()
	if !tt.HasText("Avg 36 pts") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTimelineKeepsIdentityAcrossReorder(t *testing.T) {
	tasks := []TimelineTask{{Title: "Research", Start: 0, End: 0.2}, {Title: "Design", Start: 0.2, End: 0.4}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Timeline(c, &tasks, TimelineOptions{}) }), 700, 300)
	settles(t, tt)
	r, _ := tt.Find("Research")
	x, y := r.X+r.W+50, r.Y+r.H/2
	tt.Press(x, y)
	tasks[0], tasks[1] = tasks[1], tasks[0]
	tt.Frame()
	tt.Move(x+100, y)
	tt.Release(x+100, y)
	tt.Frame()
	if tasks[0].End != 0.4 || tasks[1].End <= 0.2 {
		t.Errorf("drag resized the wrong task: %v", tasks)
	}
	panics(t, `repeated "Beta"`, func(c *ui.Context) {
		dup := []TimelineTask{{Title: "Beta", Start: 0, End: 0.5}, {Title: "Beta", Start: 0.5, End: 1}}
		Timeline(c, &dup, TimelineOptions{})
	})
}

func TestRoadmapAndMilestone(t *testing.T) {
	quarters := []Quarter{{Label: "Q3", Items: []string{"Trips"}}, {Label: "Q4", Items: []string{"Sync"}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Roadmap(c, quarters, RoadmapOptions{Current: "Q4"}) }), 600, 300)
	settles(t, tt)
	if !tt.HasText("Now") {
		t.Errorf("texts %q", tt.Texts())
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { Milestone(c, "Beta", "Nov 12", 0.999, MilestoneOptions{}) }), 500, 200)
	settles(t, tt)
	if !tt.HasText("99%") || tt.HasText("Reached") {
		t.Errorf("texts %q", tt.Texts())
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { Milestone(c, "Beta", "Nov 12", 1, MilestoneOptions{}) }), 500, 200)
	settles(t, tt)
	if !tt.HasText("Reached") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestIssueParts(t *testing.T) {
	bad := Issue{Key: "A-1", Title: "Crash", Type: "Bug", Points: -1}
	panics(t, `invalid issue "A-1"`, func(c *ui.Context) { IssueCard(c, bad, IssueCardOptions{}) })
	action := IssueNone
	issue := Issue{Key: "A-1", Title: "Crash", Type: "Bug", Priority: Urgent}
	tt := ui.NewTester(padded(func(c *ui.Context) { IssueHeader(c, issue, Doing, &action, IssueHeaderOptions{}) }), 600, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Edit")
	if action != IssueEdit || !tt.HasText("Urgent") || !tt.HasText("Doing") {
		t.Errorf("action %v texts %q", action, tt.Texts())
	}
	cues(t, before, sound.Tap)
}

func TestCommentsEmptyAndEnter(t *testing.T) {
	comments, draft := []Comment{}, ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Comments(c, &comments, &draft, CommentsOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("No comments yet") {
		t.Errorf("texts %q", tt.Texts())
	}
	click(t, tt, "Write a comment")
	tt.Type("Ship it")
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	if len(comments) != 1 || comments[0].Text != "Ship it" || tt.HasText("No comments yet") {
		t.Errorf("comments %v texts %q", comments, tt.Texts())
	}
}

func TestTimeTrackerShowsLiveTotal(t *testing.T) {
	running, total := false, 2*time.Hour+14*time.Minute+9*time.Second
	tt := ui.NewTester(padded(func(c *ui.Context) { TimeTracker(c, &running, &total, TimeTrackerOptions{}) }), 500, 200)
	settles(t, tt)
	click(t, tt, "Start")
	tt.Advance(65 * time.Second)
	if !tt.HasText("2:15:14") || !tt.HasText("Running") {
		t.Errorf("texts %q", tt.Texts())
	}
	big := time.Duration(math.MaxInt64 - int64(time.Second))
	running = true
	tt = ui.NewTester(padded(func(c *ui.Context) { TimeTracker(c, &running, &big, TimeTrackerOptions{}) }), 500, 200)
	tt.Advance(3 * time.Second)
	if tt.HasText("-") {
		t.Errorf("clock wrapped: %q", tt.Texts())
	}
}

func TestBacklogChecks(t *testing.T) {
	issues := []Issue{{Key: "A-1", Title: "One", Type: "Task", Points: math.MaxInt}, {Key: "A-2", Title: "Two", Type: "Task", Points: 1}}
	panics(t, "points overflow", func(c *ui.Context) { Backlog(c, &issues, BacklogOptions{}) })
	issues = []Issue{{Key: "A-1", Title: "One", Type: "Task"}, {Key: "A-1", Title: "Two", Type: "Task"}}
	panics(t, `repeated "A-1"`, func(c *ui.Context) { Backlog(c, &issues, BacklogOptions{}) })
	issues = []Issue{{Key: "A-1", Title: "One", Type: "Task", Points: 3}, {Key: "A-2", Title: "Two", Type: "Bug", Points: 5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Backlog(c, &issues, BacklogOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("2 issues · 8 pts") {
		t.Errorf("texts %q", tt.Texts())
	}
}
