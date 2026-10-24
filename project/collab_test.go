package project

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestEpic(t *testing.T) {
	children := []Task{{ID: "a", Title: "A", Done: true}, {ID: "b", Title: "B"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Epic(c, "Epic", children, EpicOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("1/2 done") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDependencies(t *testing.T) {
	opened := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Dependencies(c, []string{"A"}, []string{"B"}, &opened, DependenciesOptions{})
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "A")
	if opened != "A" {
		t.Errorf("opened %q", opened)
	}
	cues(t, before, sound.Open)
}

func TestWatchers(t *testing.T) {
	watching := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Watchers(c, []string{"Ana"}, &watching, WatchersOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Watch")
	if !watching {
		t.Error("not watching")
	}
	cues(t, before, sound.Toggle)
}

func TestWorkload(t *testing.T) {
	people := []WorkloadPerson{{Name: "Ana", Load: 0.8}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Workload(c, people) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Ana") || !tt.HasText("80%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestProjectCard(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ProjectCard(c, "Crychic", []string{"Ana"}, 0.5, ProjectCardOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Crychic") || !tt.HasText("1 members") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestProjectHeader(t *testing.T) {
	tab := "Tasks"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ProjectHeader(c, "Crychic", &tab, []string{"Tasks", "Docs"}, ProjectHeaderOptions{})
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Docs")
	if tab != "Docs" {
		t.Errorf("tab %q", tab)
	}
	cues(t, before, sound.Select)
}

func TestChecklist(t *testing.T) {
	items := []ChecklistItem{{Title: "One"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Checklist(c, &items, ChecklistOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Toggle")
	click(t, tt, "Add item")
	if !items[0].Done || len(items) != 2 {
		t.Errorf("items %v", items)
	}
	cues(t, before, sound.Toggle, sound.Pop)
}

func TestRetro(t *testing.T) {
	columns := []RetroColumn{{Title: "Went well"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Retro(c, &columns, RetroOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add card")
	if len(columns[0].Cards) != 1 {
		t.Errorf("cards %v", columns[0].Cards)
	}
	cues(t, before, sound.Pop)
}

func TestVote(t *testing.T) {
	votes := []int{0}
	tt := ui.NewTester(padded(func(c *ui.Context) { Vote(c, []string{"Idea"}, &votes, VoteOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Vote")
	if votes[0] != 1 {
		t.Errorf("votes %v", votes)
	}
	cues(t, before, sound.Pop)
}

func TestStandUp(t *testing.T) {
	people := []StandUpPerson{{Name: "Ana", Yesterday: "Review", Today: "Build", Blockers: "None"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { StandUp(c, people, StandUpOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Ana") || !tt.HasText("Today: Build") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDocCard(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { DocCard(c, "Spec", "2m ago", DocCardOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Spec") || !tt.HasText("Edited 2m ago") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestWhiteboard(t *testing.T) {
	notes := []StickyNote{{ID: "n1", Text: "Idea", X: 0.1, Y: 0.1}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Whiteboard(c, &notes) }), 600, 400)
	settles(t, tt)
	r, _ := tt.Find("Idea")
	before := baseline()
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	t.Logf("rect %v focused %v", r, tt.Focused("Idea"))
	tt.Move(x+40, y+20)
	tt.Release(x+40, y+20)
	tt.Frame()
	if notes[0].X <= 0.1 || notes[0].Y <= 0.1 {
		t.Errorf("note did not move: %+v", notes[0])
	}
	cues(t, before, sound.Drop)
}

func TestMindMap(t *testing.T) {
	added := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { MindMap(c, "Center", []string{"Branch"}, &added, MindMapOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add branch")
	if added != 1 {
		t.Errorf("added %d", added)
	}
	cues(t, before, sound.Pop)
}

func TestOKR(t *testing.T) {
	results := []KeyResult{{Title: "Grow", Current: 50, Target: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { OKR(c, "Objective", results) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Objective") || !tt.HasText("50/100") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestRisks(t *testing.T) {
	risks := []Risk{{Title: "Delay", Likelihood: 2, Impact: 3}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Risks(c, risks) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Delay") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestChanges(t *testing.T) {
	entries := []Change{{Text: "Updated", Time: "10:00"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Changes(c, entries) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Updated") || !tt.HasText("10:00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestRelease(t *testing.T) {
	shipped := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Release(c, "v1.0", []string{"Feature"}, &shipped, ReleaseOptions{}) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Ship")
	if !shipped {
		t.Error("release not shipped")
	}
	cues(t, before, sound.Success)
}

func TestTeamLoad(t *testing.T) {
	weeks := []TeamWeek{{Person: "Ana", Loads: []float64{0.5, 0.6, 0.7, 0.8}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { TeamLoad(c, weeks) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Ana") || !tt.HasText("80%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTemplates(t *testing.T) {
	selected := -1
	templates := []Template{{Title: "Kanban", Description: "Board"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Templates(c, templates, &selected) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Kanban")
	if selected != 0 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}

func TestQuickAdd(t *testing.T) {
	title, created := "tomorrow", false
	var due time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		QuickAdd(c, &title, &due, &created, QuickAddOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add")
	if !created || due.IsZero() || due.Day() != time.Now().AddDate(0, 0, 1).Day() {
		t.Errorf("created %v due %v", created, due)
	}
	cues(t, before, sound.Pop)
}
