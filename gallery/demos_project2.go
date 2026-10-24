package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/project"
	"github.com/egoist/mygo/ui"
)

var projectEntries2 = []Entry{
	{ID: 1141, Name: "Epic", Category: "project", Summary: "Epic with child progress.", Demo: demoEpic},
	{ID: 1142, Name: "Dependencies", Category: "project", Summary: "Blocked by and blocks links.", Demo: demoDependencies},
	{ID: 1143, Name: "Watchers", Category: "project", Summary: "Watch toggle and avatars; cue Toggle.", Demo: demoWatchers},
	{ID: 1144, Name: "Workload", Category: "project", Summary: "People with load bars.", Demo: demoWorkload},
	{ID: 1145, Name: "ProjectCard", Category: "project", Summary: "Name, members, progress.", Demo: demoProjectCard},
	{ID: 1146, Name: "ProjectHeader", Category: "project", Summary: "Icon, name, tabs; cue Select.", Demo: demoProjectHeader},
	{ID: 1147, Name: "Checklist", Category: "project", Summary: "Editable checklist; cue Toggle, Pop.", Demo: demoProjectChecklist},
	{ID: 1148, Name: "Retro", Category: "project", Summary: "Went well, improve, actions columns; cue Pop.", Demo: demoRetro},
	{ID: 1149, Name: "Vote", Category: "project", Summary: "Dot votes on cards; cue Pop.", Demo: demoVote},
	{ID: 1150, Name: "StandUp", Category: "project", Summary: "Yesterday, today, blockers per person.", Demo: demoStandUp},
	{ID: 1151, Name: "DocCard", Category: "project", Summary: "Doc with icon and last edit.", Demo: demoDocCard},
	{ID: 1152, Name: "Whiteboard", Category: "project", Summary: "Sticky notes dragged on a canvas; cue Drop.", Demo: demoWhiteboard},
	{ID: 1153, Name: "MindMap", Category: "project", Summary: "Nodes branching from a center; cue Pop.", Demo: demoMindMap},
	{ID: 1154, Name: "OKR", Category: "project", Summary: "Objective with key results bars.", Demo: demoOKR},
	{ID: 1155, Name: "Risks", Category: "project", Summary: "Likelihood by impact grid.", Demo: demoRisks},
	{ID: 1156, Name: "Changes", Category: "project", Summary: "Changelog of a task.", Demo: demoChanges},
	{ID: 1157, Name: "Release", Category: "project", Summary: "Version with items and ship; cue Success.", Demo: demoRelease},
	{ID: 1158, Name: "TeamLoad", Category: "project", Summary: "Weeks by person shaded.", Demo: demoTeamLoad},
	{ID: 1159, Name: "Templates", Category: "project", Summary: "Template cards; cue Select.", Demo: demoTemplates},
	{ID: 1160, Name: "QuickAdd", Category: "project", Summary: "Enter adds a task with parsed due; cue Pop.", Demo: demoQuickAdd},
}

func demoEpic(c *ui.Context, log func(string)) {
	children := []project.Task{{ID: "a", Title: "A", Due: "done"}, {ID: "b", Title: "B", Due: "later"}}
	project.Epic(c, "Epic", children, project.EpicOptions{})
}

func demoDependencies(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := textState(root, "opened", "")
	root.Children(func() {
		if project.Dependencies(c, []string{"A"}, []string{"B"}, opened, project.DependenciesOptions{}).Changed() {
			log("opened " + *opened)
		}
	})
}

func demoWatchers(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	watching := boolState(root, "watching", false)
	root.Children(func() {
		if project.Watchers(c, []string{"Ana"}, watching, project.WatchersOptions{}).Changed() {
			log("watching " + fmt.Sprint(*watching))
		}
	})
}

func demoWorkload(c *ui.Context, log func(string)) {
	people := []project.WorkloadPerson{{Name: "Ana", Load: 0.8}}
	project.Workload(c, people)
}

func demoProjectCard(c *ui.Context, log func(string)) {
	project.ProjectCard(c, "Crychic", []string{"Ana"}, 0.5, project.ProjectCardOptions{})
}

func demoProjectHeader(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	tab := textState(root, "tab", "Tasks")
	root.Children(func() {
		if project.ProjectHeader(c, "Crychic", tab, []string{"Tasks", "Docs"}, project.ProjectHeaderOptions{}).Changed() {
			log("tab " + *tab)
		}
	})
}

func demoProjectChecklist(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	items := ui.Local(root, "items", func() []project.ChecklistItem { return []project.ChecklistItem{{Title: "One"}} })
	root.Children(func() {
		if project.Checklist(c, items, project.ChecklistOptions{}).Changed() {
			log("items " + fmt.Sprint(len(*items)))
		}
	})
}

func demoRetro(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	columns := ui.Local(root, "columns", func() []project.RetroColumn { return []project.RetroColumn{{Title: "Went well"}} })
	root.Children(func() {
		if project.Retro(c, columns, project.RetroOptions{}).Changed() {
			log(fmt.Sprint(*columns))
		}
	})
}

func demoVote(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	votes := ui.Local(root, "votes", func() []int { return []int{0} })
	root.Children(func() {
		if project.Vote(c, []string{"Idea"}, votes, project.VoteOptions{}).Changed() {
			log(fmt.Sprint(*votes))
		}
	})
}

func demoStandUp(c *ui.Context, log func(string)) {
	people := []project.StandUpPerson{{Name: "Ana", Yesterday: "Review", Today: "Build", Blockers: "None"}}
	project.StandUp(c, people, project.StandUpOptions{})
}

func demoDocCard(c *ui.Context, log func(string)) {
	project.DocCard(c, "Spec", "2m ago", project.DocCardOptions{})
}

func demoWhiteboard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	notes := ui.Local(root, "notes", func() []project.StickyNote { return []project.StickyNote{{ID: "n1", Text: "Idea", X: 0.1, Y: 0.1}} })
	root.Children(func() {
		if project.Whiteboard(c, notes).Changed() {
			log(fmt.Sprintf("moved %.2f %.2f", (*notes)[0].X, (*notes)[0].Y))
		}
	})
}

func demoMindMap(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	branches := ui.Local(root, "branches", func() []string { return []string{"Branch"} })
	added := intState(root, "added", -1)
	root.Children(func() {
		if project.MindMap(c, "Center", *branches, added, project.MindMapOptions{}).Changed() {
			*branches = append(*branches, fmt.Sprintf("Branch %d", len(*branches)+1))
			log("branches " + fmt.Sprint(len(*branches)))
			*added = -1
		}
	})
}

func demoOKR(c *ui.Context, log func(string)) {
	results := []project.KeyResult{{Title: "Grow", Current: 50, Target: 100}}
	project.OKR(c, "Objective", results)
}

func demoRisks(c *ui.Context, log func(string)) {
	risks := []project.Risk{{Title: "Delay", Likelihood: 2, Impact: 3}}
	project.Risks(c, risks)
}

func demoChanges(c *ui.Context, log func(string)) {
	entries := []project.Change{{Text: "Updated", Time: "10:00"}}
	project.Changes(c, entries)
}

func demoRelease(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	shipped := boolState(root, "shipped", false)
	root.Children(func() {
		if project.Release(c, "v1.0", []string{"Feature"}, shipped, project.ReleaseOptions{}).Changed() {
			log("shipped")
		}
	})
}

func demoTeamLoad(c *ui.Context, log func(string)) {
	weeks := []project.TeamWeek{{Person: "Ana", Loads: []float64{0.5, 0.6, 0.7, 0.8}}}
	project.TeamLoad(c, weeks)
}

func demoTemplates(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	templates := []project.Template{{Title: "Kanban", Description: "Board"}}
	root.Children(func() {
		if project.Templates(c, templates, selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoQuickAdd(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	title, created := textState(root, "title", "tomorrow"), boolState(root, "created", false)
	due := ui.Local(root, "due", func() time.Time { return time.Time{} })
	root.Children(func() {
		root, add := project.QuickAdd(c, title, due, created, project.QuickAddOptions{})
		if root.Changed() {
			log("draft " + fmt.Sprint(utf8.RuneCountInString(*title)) + " runes")
		}
		if add.Changed() {
			log("created " + due.Format("2006-01-02"))
			*created = false
		}
	})
}
