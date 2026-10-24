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
	return project.Task{ID: "T-1", Title: "Ship UI", Assignee: "Ana", Due: "Friday", Labels: []string{"UI"}}
}

func demoProjectBoard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	lanes := ui.Local(root, "lanes", func() []project.Lane {
		return []project.Lane{{Title: "Todo", Cards: []string{"Task"}}, {Title: "Doing"}}
	})
	root.Children(func() {
		if project.Board(c, lanes, project.BoardOptions{}).Changed() {
			log(fmt.Sprint(*lanes))
		}
	})
}

func demoTaskCard(c *ui.Context, log func(string)) {
	project.TaskCard(c, demoProjectTask())
}

func demoTaskRow(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	done := boolState(root, "done", false)
	root.Children(func() {
		if project.TaskRow(c, demoProjectTask(), done).Changed() {
			log("done " + fmt.Sprint(*done))
		}
	})
}

func demoSubtasks(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	items := ui.Local(root, "items", func() []project.Subtask { return []project.Subtask{{Title: "One"}} })
	root.Children(func() {
		if project.Subtasks(c, items).Changed() {
			log(fmt.Sprint(*items))
		}
	})
}

func demoPriority(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	level := ui.Local(root, "level", func() project.PriorityLevel { return project.Medium })
	root.Children(func() {
		if project.Priority(c, level, project.PriorityOptions{}).Changed() {
			log("level " + fmt.Sprint(*level))
		}
	})
}

func demoStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	status := ui.Local(root, "status", func() project.TaskStatus { return project.Todo })
	root.Children(func() {
		if project.Status(c, status, project.StatusOptions{}).Changed() {
			log("status " + fmt.Sprint(*status))
		}
	})
}

func demoAssignee(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	root.Children(func() {
		if project.Assignee(c, []string{"Ana", "Kai"}, selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoLabels(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := ui.Local(root, "selected", func() []bool { return []bool{false} })
	root.Children(func() {
		if project.Labels(c, []string{"UI"}, selected).Changed() {
			log(fmt.Sprint(*selected))
		}
	})
}

func demoEstimate(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	points := intState(root, "points", 3)
	root.Children(func() {
		if project.Estimate(c, points, project.EstimateOptions{}).Changed() {
			log("points " + fmt.Sprint(*points))
		}
	})
}

func demoSprintHeader(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	start := timeState(root, "start", c)
	end := ui.Local(root, "end", func() time.Time { return c.Now().Add(10 * 24 * time.Hour) })
	root.Children(func() {
		project.SprintHeader(c, "Sprint 1", *start, *end, project.SprintHeaderOptions{})
	})
}

func demoBurndown(c *ui.Context, log func(string)) {
	points := []project.BurndownPoint{{Day: "Mon", Remaining: 10, Ideal: 10}, {Day: "Tue", Remaining: 5, Ideal: 5}}
	project.Burndown(c, points)
}

func demoVelocity(c *ui.Context, log func(string)) {
	sprints := []project.Sprint{{Name: "S1", Points: 20}}
	project.Velocity(c, sprints)
}

func demoProjectTimeline(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	tasks := ui.Local(root, "tasks", func() []project.TimelineTask {
		return []project.TimelineTask{{Title: "Task", Start: 0.1, End: 0.5}}
	})
	root.Children(func() {
		if project.Timeline(c, tasks).Changed() {
			log("end " + fmt.Sprint((*tasks)[0].End))
		}
	})
}

func demoRoadmap(c *ui.Context, log func(string)) {
	quarters := []project.Quarter{{Label: "Q1", Items: []string{"UI refresh"}}}
	project.Roadmap(c, quarters)
}

func demoMilestone(c *ui.Context, log func(string)) {
	project.Milestone(c, "Launch", "Nov 1", 0.5, project.MilestoneOptions{})
}

func demoIssueCard(c *ui.Context, log func(string)) {
	issue := project.Issue{Key: "UI-1", Title: "Ship UI", Type: "Story", Priority: project.High}
	project.IssueCard(c, issue, project.IssueCardOptions{})
}

func demoIssueHeader(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	action := ui.Local(root, "action", func() project.IssueAction { return project.IssueNone })
	issue := project.Issue{Key: "UI-1", Title: "Ship UI", Type: "Story", Priority: project.High}
	root.Children(func() {
		if project.IssueHeader(c, issue, project.Todo, action, project.IssueHeaderOptions{}).Changed() {
			log("action " + fmt.Sprint(*action))
		}
	})
}

func demoComments(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	comments := ui.Local(root, "comments", func() []project.Comment { return nil })
	draft := textState(root, "draft", "")
	root.Children(func() {
		root, send := project.Comments(c, comments, draft, project.CommentsOptions{})
		if root.Changed() {
			log("draft " + fmt.Sprint(utf8.RuneCountInString(*draft)) + " runes")
		}
		if send.Changed() {
			log(fmt.Sprint(len(*comments)) + " comments")
		}
	})
}

func demoTimeTracker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	running := boolState(root, "running", false)
	total := ui.Local(root, "total", func() time.Duration { return 0 })
	root.Children(func() {
		if project.TimeTracker(c, running, total, project.TimeTrackerOptions{}).Changed() {
			log("running " + fmt.Sprint(*running))
		}
	})
}

func demoBacklog(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	issues := listState(root, "issues", "Alpha", "Beta")
	root.Children(func() {
		if project.Backlog(c, issues).Changed() {
			log(fmt.Sprint(*issues))
		}
	})
}
