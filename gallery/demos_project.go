package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/project"
	"github.com/egoist/mygo/ui"
)

var projectEntries = []Entry{
	{ID: 1121, Name: "Board", Category: "project", Summary: "Kanban columns with counts and add; cards drag; cue Drop.", Demo: demoProjectBoard},
	{ID: 1122, Name: "TaskCard", Category: "project", Summary: "Title, labels, assignee, due; lifts.", Demo: demoTaskCard},
	{ID: 1123, Name: "TaskRow", Category: "project", Summary: "Check, title, due, assignee; check strikes title; cue Toggle.", Demo: demoTaskRow},
	{ID: 1124, Name: "Subtasks", Category: "project", Summary: "Checklist with done count bar; cue Toggle.", Demo: demoSubtasks},
	{ID: 1125, Name: "Priority", Category: "project", Summary: "Urgent to low picker with icons; cue Select.", Demo: demoPriority},
	{ID: 1126, Name: "Status", Category: "project", Summary: "Status picker with colored dots; cue Select.", Demo: demoStatus},
	{ID: 1127, Name: "Assignee", Category: "project", Summary: "Person picker with avatars; cue Select.", Demo: demoAssignee},
	{ID: 1128, Name: "Labels", Category: "project", Summary: "Colored label picker; cue Toggle.", Demo: demoLabels},
	{ID: 1129, Name: "Estimate", Category: "project", Summary: "Story points chips; cue Select.", Demo: demoEstimate},
	{ID: 1130, Name: "SprintHeader", Category: "project", Summary: "Sprint name, dates, days left.", Demo: demoSprintHeader},
	{ID: 1131, Name: "Burndown", Category: "project", Summary: "Remaining against ideal line.", Demo: demoBurndown},
	{ID: 1132, Name: "Velocity", Category: "project", Summary: "Points per sprint bars.", Demo: demoVelocity},
	{ID: 1133, Name: "Timeline", Category: "project", Summary: "Gantt bars by task; drag resizes; cue Tick.", Demo: demoProjectTimeline},
	{ID: 1134, Name: "Roadmap", Category: "project", Summary: "Quarters with initiatives.", Demo: demoRoadmap},
	{ID: 1135, Name: "Milestone", Category: "project", Summary: "Diamond with date and progress.", Demo: demoMilestone},
	{ID: 1136, Name: "IssueCard", Category: "project", Summary: "Key, title, type icon, priority.", Demo: demoIssueCard},
	{ID: 1137, Name: "IssueHeader", Category: "project", Summary: "Key, title, status and actions.", Demo: demoIssueHeader},
	{ID: 1138, Name: "Comments", Category: "project", Summary: "Thread with composer; cue Send.", Demo: demoComments},
	{ID: 1139, Name: "TimeTracker", Category: "project", Summary: "Start and stop with total; cue Toggle.", Demo: demoTimeTracker},
	{ID: 1140, Name: "Backlog", Category: "project", Summary: "Ordered issues dragged; cue Drop.", Demo: demoBacklog},
}

func demoProjectTask() project.Task {
	return project.Task{ID: "ATL-142", Title: "Offline sync for saved trips", Assignee: "Maya Chen", Due: "Oct 30", Labels: []string{"Mobile", "Sync"}}
}

func demoProjectIssues() []project.Issue {
	return []project.Issue{
		{Key: "ATL-142", Title: "Offline sync for saved trips", Type: "Story", Priority: project.High, Points: 8, Assignee: "Maya Chen"},
		{Key: "ATL-151", Title: "Crash when opening a shared map link", Type: "Bug", Priority: project.Urgent, Points: 3, Assignee: "Daniel Okafor"},
		{Key: "ATL-147", Title: "Dark mode for the trip planner", Type: "Story", Priority: project.Medium, Points: 5, Assignee: "Priya Nair"},
		{Key: "ATL-139", Title: "Migrate push tokens to the new service", Type: "Task", Priority: project.Medium, Points: 3},
		{Key: "ATL-160", Title: "Trip sharing and invites", Type: "Epic", Priority: project.Low, Points: 13, Assignee: "Lucas Martin"},
	}
}

func demoProjectBoard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	lanes := ui.Local(root, "lanes", func() []project.Lane {
		return []project.Lane{
			{Title: "To do", Cards: []project.Task{
				{ID: "ATL-147", Title: "Dark mode for the trip planner", Assignee: "Priya Nair", Labels: []string{"Design"}},
				{ID: "ATL-139", Title: "Migrate push tokens", Due: "Nov 2", Labels: []string{"Backend"}},
			}},
			{Title: "In progress", Cards: []project.Task{demoProjectTask()}},
			{Title: "Done", Cards: []project.Task{{ID: "ATL-128", Title: "New onboarding screens", Assignee: "Sofia Rossi", Labels: []string{"Design"}, Done: true}}},
		}
	})
	root.Children(func() {
		if project.Board(c, lanes, project.BoardOptions{}).Changed() {
			counts := []int{}
			for _, lane := range *lanes {
				counts = append(counts, len(lane.Cards))
			}
			log(fmt.Sprint("cards ", counts))
		}
	})
}

func demoTaskCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Children(func() { project.TaskCard(c, demoProjectTask()) })
}

func demoTaskRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440).Gap(4)
	done, shipped := boolState(root, "done", false), boolState(root, "shipped", true)
	root.Children(func() {
		if project.TaskRow(c, demoProjectTask(), done).Changed() {
			log(fmt.Sprint("done ", *done))
		}
		if project.TaskRow(c, project.Task{ID: "ATL-128", Title: "New onboarding screens", Assignee: "Sofia Rossi", Due: "Oct 20"}, shipped).Changed() {
			log(fmt.Sprint("shipped ", *shipped))
		}
	})
}

func demoSubtasks(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	items := ui.Local(root, "items", func() []project.Subtask {
		return []project.Subtask{{Title: "Cache trips in SQLite", Done: true}, {Title: "Queue edits made offline", Done: true},
			{Title: "Resolve conflicts on reconnect"}, {Title: "Show a sync status badge"}}
	})
	root.Children(func() {
		if project.Subtasks(c, items, project.SubtasksOptions{}).Changed() {
			n := 0
			for _, item := range *items {
				if item.Done {
					n++
				}
			}
			log(fmt.Sprintf("%d done", n))
		}
	})
}

func demoPriority(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	level := ui.Local(root, "level", func() project.PriorityLevel { return project.Medium })
	root.Children(func() {
		if project.Priority(c, level, project.PriorityOptions{}).Changed() {
			log(fmt.Sprint("level ", *level))
		}
	})
}

func demoStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	status := ui.Local(root, "status", func() project.TaskStatus { return project.Doing })
	root.Children(func() {
		if project.Status(c, status, project.StatusOptions{Doing: "In progress"}).Changed() {
			log(fmt.Sprint("status ", *status))
		}
	})
}

func demoAssignee(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	selected := intState(root, "selected", 0)
	people := []string{"Maya Chen", "Daniel Okafor", "Priya Nair", "Lucas Martin"}
	root.Children(func() {
		if project.Assignee(c, people, selected, project.AssigneeOptions{}).Changed() {
			log(fmt.Sprint("selected ", *selected))
		}
	})
}

func demoLabels(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	selected := ui.Local(root, "selected", func() []bool { return []bool{true, false, true, false, false} })
	labels := []string{"Mobile", "Backend", "Sync", "Design", "Performance"}
	root.Children(func() {
		if project.Labels(c, labels, selected).Changed() {
			log(fmt.Sprint(*selected))
		}
	})
}

func demoEstimate(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	points := intState(root, "points", 3)
	root.Children(func() {
		if project.Estimate(c, points, project.EstimateOptions{}).Changed() {
			log(fmt.Sprint("points ", *points))
		}
	})
}

func demoSprintHeader(c *ui.Context, log func(string)) {
	day := c.Now().Truncate(24 * time.Hour)
	start, end := day.AddDate(0, 0, -3), day.AddDate(0, 0, 7)
	o := project.SprintHeaderOptions{Goal: "Ship offline trips to the beta group"}
	ui.Column(c).Width(400).Children(func() { project.SprintHeader(c, "Sprint 14 · Atlas mobile", start, end, o) })
}

func demoBurndown(c *ui.Context, log func(string)) {
	remaining := []float64{42, 40, 37, 37, 31, 26, 24, 19}
	points := make([]project.BurndownPoint, len(remaining))
	for i, v := range remaining {
		points[i] = project.BurndownPoint{Day: fmt.Sprint("D", i+1), Remaining: v, Ideal: 42 - 42*float64(i)/9}
	}
	ui.Column(c).Width(440).Children(func() { project.Burndown(c, points, project.BurndownOptions{Title: "Sprint 14 burndown"}) })
}

func demoVelocity(c *ui.Context, log func(string)) {
	sprints := []project.Sprint{{Name: "S9", Points: 31}, {Name: "S10", Points: 38}, {Name: "S11", Points: 29},
		{Name: "S12", Points: 41}, {Name: "S13", Points: 36}, {Name: "S14", Points: 44}}
	ui.Column(c).Width(420).Children(func() { project.Velocity(c, sprints, project.VelocityOptions{}) })
}

func demoProjectTimeline(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	tasks := ui.Local(root, "tasks", func() []project.TimelineTask {
		return []project.TimelineTask{{Title: "Research", Start: 0, End: 0.2}, {Title: "Design", Start: 0.15, End: 0.45},
			{Title: "Build", Start: 0.4, End: 0.8}, {Title: "Beta", Start: 0.75, End: 0.95}}
	})
	root.Children(func() {
		if project.Timeline(c, tasks, project.TimelineOptions{Ticks: []string{"Oct 5", "Oct 12", "Oct 19", "Oct 26", "Nov 2"}}).Changed() {
			ends := []string{}
			for _, task := range *tasks {
				ends = append(ends, fmt.Sprintf("%.2f", task.End))
			}
			log(fmt.Sprint("ends ", ends))
		}
	})
}

func demoRoadmap(c *ui.Context, log func(string)) {
	quarters := []project.Quarter{
		{Label: "Q3 2026", Items: []string{"Saved trips", "New onboarding"}},
		{Label: "Q4 2026", Items: []string{"Offline sync", "Dark mode", "Shared maps"}},
		{Label: "Q1 2027", Items: []string{"Trip invites", "Android widgets"}},
	}
	ui.Column(c).Width(520).Children(func() { project.Roadmap(c, quarters, project.RoadmapOptions{Current: "Q4 2026"}) })
}

func demoMilestone(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Gap(8).Children(func() {
		project.Milestone(c, "Public beta", "Nov 12", 0.64, project.MilestoneOptions{})
		project.Milestone(c, "Design freeze", "Oct 16", 1, project.MilestoneOptions{})
	})
}

func demoIssueCard(c *ui.Context, log func(string)) {
	issues := demoProjectIssues()
	ui.Column(c).Width(440).Gap(8).Children(func() {
		for i := range 3 {
			c.Key(i)
			project.IssueCard(c, issues[i], project.IssueCardOptions{})
		}
	})
}

func demoIssueHeader(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	action := ui.Local(root, "action", func() project.IssueAction { return project.IssueNone })
	root.Children(func() {
		if project.IssueHeader(c, demoProjectIssues()[1], project.Doing, action, project.IssueHeaderOptions{Doing: "In progress"}).Changed() {
			log(fmt.Sprint("action ", *action))
		}
	})
}

func demoComments(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	comments := ui.Local(root, "comments", func() []project.Comment {
		return []project.Comment{
			{Author: "Daniel Okafor", When: "2h ago", Text: "Repro: open a shared link while the app is cold. It crashes in MapRouter."},
			{Author: "Maya Chen", When: "1h ago", Text: "Looks like the deep link fires before the store loads. I'll guard it."},
			{Author: "Priya Nair", When: "20m ago", Text: "Can we also show a toast when the link has expired?"},
		}
	})
	draft := textState(root, "draft", "")
	root.Children(func() {
		thread, send := project.Comments(c, comments, draft, project.CommentsOptions{Author: "Lucas Martin"})
		if thread.Changed() {
			log(fmt.Sprintf("draft %d runes", utf8.RuneCountInString(*draft)))
		}
		if send.Changed() {
			log(fmt.Sprintf("%d comments", len(*comments)))
		}
	})
}

func demoTimeTracker(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	running := boolState(root, "running", false)
	total := ui.Local(root, "total", func() time.Duration { return 2*time.Hour + 14*time.Minute + 9*time.Second })
	root.Children(func() {
		if project.TimeTracker(c, running, total, project.TimeTrackerOptions{Task: "ATL-151 · Shared link crash"}).Changed() {
			log(fmt.Sprint("running ", *running))
		}
	})
}

func demoBacklog(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	issues := ui.Local(root, "issues", demoProjectIssues)
	root.Children(func() {
		if project.Backlog(c, issues, project.BacklogOptions{}).Changed() {
			keys := []string{}
			for _, issue := range *issues {
				keys = append(keys, issue.Key)
			}
			log(fmt.Sprint(keys))
		}
	})
}
