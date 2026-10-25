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
	children := []project.Task{
		{ID: "ATL-121", Title: "Cache trips in SQLite", Assignee: "Maya Chen", Done: true},
		{ID: "ATL-122", Title: "Queue edits made offline", Assignee: "Daniel Okafor", Done: true},
		{ID: "ATL-123", Title: "Resolve conflicts on reconnect", Assignee: "Priya Nair"},
		{ID: "ATL-124", Title: "Show a sync status badge", Assignee: "Lucas Martin"},
		{ID: "ATL-125", Title: "Background sync on iOS", Assignee: "Sofia Rossi"},
	}
	ui.Column(c).Width(440).Children(func() {
		project.Epic(c, "Offline mode for saved trips", children, project.EpicOptions{Key: "ATL-120 · Epic"})
	})
}

func demoDependencies(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	opened := textState(root, "opened", "")
	blockedBy := []project.Link{
		{Key: "ATL-151", Title: "Crash when opening a shared map link", Status: project.Doing},
		{Key: "ATL-138", Title: "Migrate push tokens", Status: project.Done},
	}
	blocks := []project.Link{
		{Key: "ATL-160", Title: "Public beta release", Status: project.Todo},
		{Key: "ATL-163", Title: "App Store screenshots", Status: project.Todo},
	}
	root.Children(func() {
		if project.Dependencies(c, blockedBy, blocks, opened, project.DependenciesOptions{}).Changed() {
			log("opened " + *opened)
		}
	})
}

func demoWatchers(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	watching := boolState(root, "watching", false)
	root.Children(func() {
		if project.Watchers(c, []string{"Maya Chen", "Daniel Okafor", "Priya Nair"}, watching, project.WatchersOptions{}).Changed() {
			log("watching " + fmt.Sprint(*watching))
		}
	})
}

func demoWorkload(c *ui.Context, log func(string)) {
	people := []project.WorkloadPerson{
		{Name: "Maya Chen", Role: "iOS", Load: 0.72},
		{Name: "Daniel Okafor", Role: "Backend", Load: 1.15},
		{Name: "Priya Nair", Role: "Design", Load: 0.9},
		{Name: "Lucas Martin", Role: "Android", Load: 0.55},
		{Name: "Sofia Rossi", Role: "QA", Load: 0.8},
	}
	ui.Column(c).Width(420).Children(func() { project.Workload(c, people, project.WorkloadOptions{}) })
}

func demoProjectCard(c *ui.Context, log func(string)) {
	members := []string{"Maya Chen", "Daniel Okafor", "Priya Nair", "Lucas Martin", "Sofia Rossi"}
	ui.Column(c).Width(360).Children(func() {
		project.ProjectCard(c, "Atlas mobile 2.0", members, 0.64, project.ProjectCardOptions{
			Description: "Offline trips, faster maps and a new onboarding.", Due: "Nov 12"})
	})
}

func demoProjectHeader(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	tab := textState(root, "tab", "Board")
	tabs := []string{"Overview", "Board", "Timeline", "Docs", "Releases"}
	root.Children(func() {
		o := project.ProjectHeaderOptions{Subtitle: "Travel planner for iOS and Android"}
		if project.ProjectHeader(c, "Atlas mobile 2.0", tab, tabs, o).Changed() {
			log("tab " + *tab)
		}
	})
}

func demoProjectChecklist(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	items := ui.Local(root, "items", func() []project.ChecklistItem {
		return []project.ChecklistItem{
			{ID: "strings", Title: "Freeze strings for translation", Done: true},
			{ID: "regression", Title: "Run the full regression suite", Done: true},
			{ID: "screens", Title: "Update App Store screenshots"},
			{ID: "notes", Title: "Write release notes"},
			{ID: "support", Title: "Brief the support team"},
		}
	})
	root.Children(func() {
		if project.Checklist(c, items, project.ChecklistOptions{Title: "Launch checklist"}).Changed() {
			done := 0
			for _, item := range *items {
				if item.Done {
					done++
				}
			}
			log(fmt.Sprintf("items %d, %d done", len(*items), done))
		}
	})
}

func demoRetro(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	columns := ui.Local(root, "columns", func() []project.RetroColumn {
		return []project.RetroColumn{
			{Title: "Went well", Cards: []string{"Offline sync landed a week early", "Design reviews took a day, not three"}},
			{Title: "To improve", Cards: []string{"Flaky UI tests blocked merges twice"}},
			{Title: "Action items", Cards: []string{"Quarantine flaky tests by Friday"}},
		}
	})
	root.Children(func() {
		if project.Retro(c, columns, project.RetroOptions{}).Changed() {
			counts := make([]int, len(*columns))
			for i, column := range *columns {
				counts[i] = len(column.Cards)
			}
			log(fmt.Sprintf("cards %v", counts))
		}
	})
}

func demoVote(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	cards := []string{"Fix the flaky UI tests", "Ship dark mode first", "Weekly design office hours", "Move stand-up to 9:45"}
	votes := ui.Local(root, "votes", func() []int { return []int{4, 2, 1, 0} })
	root.Children(func() {
		if project.Vote(c, cards, votes, project.VoteOptions{Title: "Next sprint focus"}).Changed() {
			log(fmt.Sprintf("votes %v", *votes))
		}
	})
}

func demoStandUp(c *ui.Context, log func(string)) {
	people := []project.StandUpPerson{
		{Name: "Maya Chen", Yesterday: "Finished the SQLite trip cache", Today: "Pair with Priya on the sync badge"},
		{Name: "Daniel Okafor", Yesterday: "Reproduced the shared link crash", Today: "Guard the deep link router",
			Blockers: "Waiting on staging push certificates"},
		{Name: "Priya Nair", Yesterday: "Reviewed onboarding copy", Today: "Final sync badge states"},
	}
	ui.Column(c).Width(440).Children(func() { project.StandUp(c, people, project.StandUpOptions{}) })
}

func demoDocCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Gap(8).Children(func() {
		project.DocCard(c, "Offline sync design", "2h ago", project.DocCardOptions{Author: "Priya Nair",
			Preview: "How trips, edits and conflicts move between the device and the server."})
		project.DocCard(c, "Release 2.0 runbook", "yesterday", project.DocCardOptions{Author: "Daniel Okafor", Icon: "notebook-pen"})
	})
}

func demoWhiteboard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	notes := ui.Local(root, "notes", func() []project.StickyNote {
		return []project.StickyNote{
			{ID: "n1", Text: "Map loads slowly on 3G", X: 0.04, Y: 0.06},
			{ID: "n2", Text: "Show trip size before download", X: 0.38, Y: 0.1},
			{ID: "n3", Text: "Sync badge confuses testers", X: 0.68, Y: 0.05},
			{ID: "n4", Text: "Offer a Wi-Fi only option", X: 0.12, Y: 0.55},
			{ID: "n5", Text: "Retry quietly, then tell me", X: 0.5, Y: 0.6},
		}
	})
	root.Children(func() {
		if project.Whiteboard(c, notes).Changed() {
			log("notes moved")
		}
	})
}

func demoMindMap(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	branches := ui.Local(root, "branches", func() []string { return []string{"Offline", "Maps", "Onboarding", "Sharing", "Widgets"} })
	added := intState(root, "added", -1)
	more := []string{"Dark mode", "Payments", "Search"}
	root.Children(func() {
		if project.MindMap(c, "Atlas 2.0", *branches, added, project.MindMapOptions{}).Changed() {
			*branches = append(*branches, more[len(*branches)-5])
			log("branches " + fmt.Sprint(len(*branches)))
			*added = -1
		}
	})
}

func demoOKR(c *ui.Context, log func(string)) {
	results := []project.KeyResult{
		{Title: "Ship offline mode to every user", Current: 64, Target: 100},
		{Title: "Reach a 4.7 App Store rating", Current: 4.5, Target: 4.7},
		{Title: "Grow weekly active planners to 40,000", Current: 26500, Target: 40000},
	}
	ui.Column(c).Width(440).Children(func() {
		project.OKR(c, "Make Atlas the fastest way to plan a trip", results, project.OKROptions{Owner: "Maya Chen", Period: "Q4 2026"})
	})
}

func demoRisks(c *ui.Context, log func(string)) {
	risks := []project.Risk{
		{Title: "App Store review delay", Likelihood: 3, Impact: 5},
		{Title: "Push certificates expire", Likelihood: 2, Impact: 4},
		{Title: "Map tile costs spike", Likelihood: 3, Impact: 3},
		{Title: "Designer on leave", Likelihood: 4, Impact: 2},
		{Title: "Offline data loss", Likelihood: 1, Impact: 5},
	}
	ui.Column(c).Width(440).Children(func() { project.Risks(c, risks, project.RisksOptions{}) })
}

func demoChanges(c *ui.Context, log func(string)) {
	entries := []project.Change{
		{Who: "Daniel Okafor", Text: "created the issue", Time: "Oct 19"},
		{Who: "Maya Chen", Text: "set priority", Time: "Oct 19", From: "Medium", To: "Urgent"},
		{Who: "Daniel Okafor", Text: "moved it", Time: "Oct 21", From: "To do", To: "In progress"},
		{Who: "Priya Nair", Text: "attached a screen recording", Time: "Oct 22"},
		{Text: "Linked to release v2.0.0", Time: "Oct 23"},
	}
	ui.Column(c).Width(420).Children(func() { project.Changes(c, entries, project.ChangesOptions{}) })
}

func demoRelease(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	shipped := boolState(root, "shipped", false)
	items := []project.ReleaseItem{
		{Kind: project.ReleaseAdded, Text: "Offline mode for saved trips"},
		{Kind: project.ReleaseAdded, Text: "Dark mode across every screen"},
		{Kind: project.ReleaseFixed, Text: "Crash when opening a shared map link"},
		{Kind: project.ReleaseChanged, Text: "Faster map tiles on slow networks"},
	}
	root.Children(func() {
		if project.Release(c, "v2.0.0", items, shipped, project.ReleaseOptions{Date: "Planned for Nov 12"}).Changed() {
			log("shipped")
		}
	})
}

func demoTeamLoad(c *ui.Context, log func(string)) {
	weeks := []project.TeamWeek{
		{Person: "Maya Chen", Loads: []float64{0.7, 0.85, 0.9, 0.6}},
		{Person: "Daniel Okafor", Loads: []float64{1.1, 1.2, 0.95, 0.8}},
		{Person: "Priya Nair", Loads: []float64{0.9, 0.6, 0.5, 0.4}},
		{Person: "Lucas Martin", Loads: []float64{0.5, 0.7, 1.05, 0.9}},
		{Person: "Sofia Rossi", Loads: []float64{0.6, 0.8, 0.9, 1}},
	}
	ui.Column(c).Width(480).Children(func() {
		project.TeamLoad(c, weeks, project.TeamLoadOptions{Weeks: []string{"Oct 26", "Nov 2", "Nov 9", "Nov 16"}})
	})
}

func demoTemplates(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	selected := intState(root, "selected", 0)
	templates := []project.Template{
		{Title: "Kanban board", Description: "To do, doing, done lanes for steady work.", Icon: "kanban"},
		{Title: "Sprint planning", Description: "Two-week sprints with a backlog and burndown.", Icon: "zap"},
		{Title: "Product launch", Description: "Checklist, risks and a release plan.", Icon: "rocket"},
		{Title: "Bug triage", Description: "Sort incoming bugs by priority and owner.", Icon: "bug"},
	}
	root.Children(func() {
		if project.Templates(c, templates, selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoQuickAdd(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	title, created := textState(root, "title", "Write release notes tomorrow"), boolState(root, "created", false)
	due := ui.Local(root, "due", func() time.Time { return time.Time{} })
	root.Children(func() {
		root, add := project.QuickAdd(c, title, due, created, project.QuickAddOptions{})
		if root.Changed() {
			log("draft " + fmt.Sprint(utf8.RuneCountInString(*title)) + " runes")
		}
		if add.Changed() {
			log("created " + due.Format("2006-01-02"))
			*title, *created = "", false
		}
	})
}
