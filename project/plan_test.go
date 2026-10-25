package project

import (
	"math"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// exact reports a text shown verbatim, not as a substring.
func exact(tt *ui.Tester, s string) bool { return slices.Contains(tt.Texts(), s) }

func TestDependenciesKeepIdentity(t *testing.T) {
	blocked := []Link{{Key: "ATL-1", Title: "Auth"}, {Key: "ATL-2", Title: "Billing", Status: Done}}
	opened := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Dependencies(c, blocked, nil, &opened, DependenciesOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "1 blocking") || !exact(tt, "None") {
		t.Errorf("texts %q", tt.Texts())
	}
	swapUnderPress(t, tt, "ATL-1", func() { blocked[0], blocked[1] = blocked[1], blocked[0] })
	if opened == "ATL-2" {
		t.Error("release opened the link that slid in")
	}
	blocked[0].Status, blocked[1].Status = Done, Done
	tt.Frame()
	if !exact(tt, "Unblocked") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, `repeated "ATL-1"`, func(c *ui.Context) {
		Dependencies(c, []Link{{Key: "ATL-1", Title: "A"}, {Key: "ATL-1", Title: "B"}}, nil, &opened, DependenciesOptions{})
	})
}

func TestWatchersCountsYou(t *testing.T) {
	watching := false
	people := []string{"Maya Chen"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Watchers(c, people, &watching, WatchersOptions{}) }), 500, 200)
	settles(t, tt)
	if !exact(tt, "1 watcher") {
		t.Errorf("texts %q", tt.Texts())
	}
	click(t, tt, "Watch")
	if !exact(tt, "2 watchers") {
		t.Errorf("texts %q", tt.Texts())
	}
	people, watching = nil, false
	tt.Frame()
	if !exact(tt, "No watchers yet") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, `repeated "You"`, func(c *ui.Context) { Watchers(c, []string{"You"}, &watching, WatchersOptions{}) })
}

func TestWorkloadOverCapacity(t *testing.T) {
	people := []WorkloadPerson{{Name: "Daniel Okafor", Load: 1.15}, {Name: "Maya Chen", Load: 0.7}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Workload(c, people, WorkloadOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "1 over capacity") || !exact(tt, "115%") {
		t.Errorf("texts %q", tt.Texts())
	}
	people[0].Load = 0.9
	tt.Frame()
	if !exact(tt, "Balanced") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "load 2.5", func(c *ui.Context) { Workload(c, []WorkloadPerson{{Name: "A", Load: 2.5}}, WorkloadOptions{}) })
}

func TestProjectCardPercent(t *testing.T) {
	for _, tc := range []struct {
		value float64
		want  string
	}{{0.29, "29%"}, {0.999, "99%"}, {1, "100%"}} {
		tt := ui.NewTester(padded(func(c *ui.Context) {
			ProjectCard(c, "Atlas", []string{"Maya Chen"}, tc.value, ProjectCardOptions{})
		}), 500, 260)
		tt.Frame()
		if !exact(tt, tc.want) || !exact(tt, "1 member") {
			t.Errorf("value %v: texts %q, want %q", tc.value, tt.Texts(), tc.want)
		}
	}
}

func TestProjectHeaderTabsKeepIdentity(t *testing.T) {
	tabs, tab := []string{"Overview", "Board", "Docs"}, "Docs"
	tt := ui.NewTester(padded(func(c *ui.Context) { ProjectHeader(c, "Atlas", &tab, tabs, ProjectHeaderOptions{}) }), 500, 260)
	settles(t, tt)
	swapUnderPress(t, tt, "Overview", func() { tabs[0], tabs[1] = tabs[1], tabs[0] })
	if tab == "Board" {
		t.Error("release picked the tab that slid in")
	}
	panics(t, `no tab "Gone"`, func(c *ui.Context) {
		gone := "Gone"
		ProjectHeader(c, "Atlas", &gone, []string{"Board"}, ProjectHeaderOptions{})
	})
}

func TestChecklistEmptyAndAdd(t *testing.T) {
	var items []ChecklistItem
	tt := ui.NewTester(padded(func(c *ui.Context) { Checklist(c, &items, ChecklistOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "No items yet") || !exact(tt, "0 of 0") {
		t.Errorf("texts %q", tt.Texts())
	}
	click(t, tt, "Add item")
	if len(items) != 1 || items[0].Title != "New item" {
		t.Errorf("items %v", items)
	}
	click(t, tt, "Toggle New item")
	if !items[0].Done || !exact(tt, "1 of 1") {
		t.Errorf("items %v texts %q", items, tt.Texts())
	}
}

func TestRetroColumnsKeepIdentity(t *testing.T) {
	columns := []RetroColumn{{Title: "Went well"}, {Title: "To improve"}, {Title: "Action items"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Retro(c, &columns, RetroOptions{}) }), 600, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Add card", func() { columns[0], columns[1] = columns[1], columns[0] })
	if len(columns[0].Cards) != 0 {
		t.Errorf("card went to the column that slid in: %v", columns)
	}
	click(t, tt, "Add card")
	if columns[0].Cards[0] != "New card" {
		t.Errorf("columns %v", columns)
	}
	panics(t, `repeated "A"`, func(c *ui.Context) {
		dup := []RetroColumn{{Title: "A"}, {Title: "A"}}
		Retro(c, &dup, RetroOptions{})
	})
}

func TestVoteIdentityAndCopy(t *testing.T) {
	cards, votes := []string{"Dark mode", "Flaky tests"}, []int{0, 1}
	tt := ui.NewTester(padded(func(c *ui.Context) { Vote(c, cards, &votes, VoteOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "1 vote") || !exact(tt, "Top") {
		t.Errorf("texts %q", tt.Texts())
	}
	swapUnderPress(t, tt, "Vote for Dark mode", func() {
		cards[0], cards[1] = cards[1], cards[0]
		votes[0], votes[1] = votes[1], votes[0]
	})
	if votes[0] != 1 {
		t.Errorf("vote went to the card that slid in: %v %v", cards, votes)
	}
	panics(t, "invalid count", func(c *ui.Context) {
		big := []int{math.MaxInt, 1}
		Vote(c, []string{"A", "B"}, &big, VoteOptions{})
	})
}

func TestStandUpStatus(t *testing.T) {
	people := []StandUpPerson{{Name: "Maya Chen", Yesterday: "Cache", Today: "Badge"},
		{Name: "Daniel Okafor", Yesterday: "Repro", Today: "Guard", Blockers: "Certificates"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { StandUp(c, people, StandUpOptions{}) }), 500, 600)
	settles(t, tt)
	for _, want := range []string{"On track", "Blocked", "No blockers", "Certificates"} {
		if !exact(tt, want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
}

func TestDocCardAuthor(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		DocCard(c, "Offline sync design", "2h ago", DocCardOptions{Author: "Priya Nair"})
	}), 500, 200)
	tt.Frame()
	if !exact(tt, "Edited 2h ago by Priya Nair") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMindMapLimit(t *testing.T) {
	nine := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	added := -1
	panics(t, "at most 8", func(c *ui.Context) { MindMap(c, "Atlas", nine, &added, MindMapOptions{}) })
	eight := nine[:8]
	tt := ui.NewTester(padded(func(c *ui.Context) { MindMap(c, "Atlas", eight, &added, MindMapOptions{}) }), 500, 500)
	settles(t, tt)
	tt.Click("Add branch")
	tt.Frame()
	if added != -1 {
		t.Errorf("full map still added: %d", added)
	}
	ui.Render(func(c *ui.Context) { MindMap(c, "Atlas", []string{"center"}, &added, MindMapOptions{}) }, 500, 500, 1)
}

func TestOKRScoreAndValues(t *testing.T) {
	results := []KeyResult{{Title: "Half", Current: 1, Target: 2}, {Title: "Done", Current: 3, Target: 2}}
	tt := ui.NewTester(padded(func(c *ui.Context) { OKR(c, "Ship", results, OKROptions{}) }), 500, 300)
	tt.Frame()
	if !exact(tt, "75%") || !exact(tt, "3/2") {
		t.Errorf("texts %q", tt.Texts())
	}
	results = []KeyResult{{Title: "Grow", Current: 1e6, Target: 2e6}}
	tt.Frame()
	if !exact(tt, "1000000/2000000") || !exact(tt, "50%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestRisksLevelsAndOverflow(t *testing.T) {
	risks := []Risk{{Title: "A", Likelihood: 3, Impact: 5}, {Title: "B", Likelihood: 3, Impact: 5},
		{Title: "C", Likelihood: 3, Impact: 5}, {Title: "D", Likelihood: 1, Impact: 1}, {Title: "E", Likelihood: 2, Impact: 3}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Risks(c, risks, RisksOptions{}) }), 500, 700)
	tt.Frame()
	for _, want := range []string{"+1", "High", "Medium", "Low"} {
		if !exact(tt, want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	panics(t, "off the 5 by 5 grid", func(c *ui.Context) { Risks(c, []Risk{{Title: "X", Likelihood: 6, Impact: 1}}, RisksOptions{}) })
}

func TestChangesEmptyAndPairs(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Changes(c, nil, ChangesOptions{}) }), 500, 200)
	tt.Frame()
	if !exact(tt, "No changes yet") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "both or neither", func(c *ui.Context) {
		Changes(c, []Change{{Text: "moved", Time: "now", From: "To do"}}, ChangesOptions{})
	})
}

func TestReleaseShipsOnce(t *testing.T) {
	shipped := false
	items := []ReleaseItem{{Kind: ReleaseFixed, Text: "Crash on share"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Release(c, "v2.0.0", items, &shipped, ReleaseOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "Draft") || !exact(tt, "Fixed") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	r, _ := tt.Find("Ship")
	click(t, tt, "Ship")
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !shipped || exact(tt, "Draft") {
		t.Errorf("shipped %v texts %q", shipped, tt.Texts())
	}
	cues(t, before, sound.Success)
}

func TestTeamLoadWeeks(t *testing.T) {
	weeks := []TeamWeek{{Person: "Maya Chen", Loads: []float64{0.5, 1.2}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { TeamLoad(c, weeks, TeamLoadOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "W1") || !exact(tt, "W2") || !exact(tt, "120%") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, `gives "B" 1 weeks, not 2`, func(c *ui.Context) {
		TeamLoad(c, []TeamWeek{{Person: "A", Loads: []float64{0.1, 0.2}}, {Person: "B", Loads: []float64{0.3}}}, TeamLoadOptions{})
	})
	panics(t, "3 week names for 2 weeks", func(c *ui.Context) { TeamLoad(c, weeks, TeamLoadOptions{Weeks: []string{"a", "b", "c"}}) })
}

func TestTemplatesKeepIdentity(t *testing.T) {
	templates := []Template{{Title: "Kanban board", Description: "Lanes"}, {Title: "Bug triage", Description: "Sort"}}
	selected := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Templates(c, templates, &selected) }), 500, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Kanban board", func() { templates[0], templates[1] = templates[1], templates[0] })
	if selected >= 0 && templates[selected].Title != "Kanban board" {
		t.Errorf("selected the template that slid in: %d", selected)
	}
}

func TestQuickAddNote(t *testing.T) {
	title, created := "Write notes tomorrow", false
	var due time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) { QuickAdd(c, &title, &due, &created, QuickAddOptions{}) }), 500, 200)
	settles(t, tt)
	if !strings.Contains(strings.Join(tt.Texts(), "|"), "Due ") {
		t.Errorf("no parsed due chip: %q", tt.Texts())
	}
	click(t, tt, "Add")
	if !strings.Contains(strings.Join(tt.Texts(), "|"), "Added · Due ") {
		t.Errorf("texts %q", tt.Texts())
	}
	title = "Write notes"
	click(t, tt, "Add")
	if !exact(tt, "Added · no due date") {
		t.Errorf("texts %q", tt.Texts())
	}
	click(t, tt, "Task and due")
	tt.Type("!")
	tt.Frame()
	if strings.Contains(strings.Join(tt.Texts(), "|"), "Added") {
		t.Errorf("note survived an edit: %q", tt.Texts())
	}
}

func TestChecklistKeepsIdentity(t *testing.T) {
	items := []ChecklistItem{{ID: "a", Title: "Freeze strings"}, {ID: "b", Title: "Write notes"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Checklist(c, &items, ChecklistOptions{}) }), 500, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Toggle Freeze strings", func() { items[0], items[1] = items[1], items[0] })
	if items[0].Done {
		t.Errorf("toggled the item that slid in: %v", items)
	}
	panics(t, `repeated "a"`, func(c *ui.Context) {
		dup := []ChecklistItem{{ID: "a", Title: "x"}, {ID: "a", Title: "y"}}
		Checklist(c, &dup, ChecklistOptions{})
	})
}

func TestWatchersToggleKeepsFocus(t *testing.T) {
	watching := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Watchers(c, nil, &watching, WatchersOptions{}) }), 500, 200)
	settles(t, tt)
	click(t, tt, "Watch")
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if watching {
		t.Error("Space did not toggle back: the button lost focus")
	}
}

func TestWhiteboardStillDragAndRepeats(t *testing.T) {
	notes := []StickyNote{{ID: "n1", Text: "Edge", X: 0, Y: 0.2}}
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Whiteboard(c, &notes).Changed() {
			changes++
		}
	}), 600, 400)
	settles(t, tt)
	before := baseline()
	r, _ := tt.Find("Edge")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W/2-40, r.Y+r.H/2)
	tt.Release(r.X+r.W/2-40, r.Y+r.H/2)
	tt.Frame()
	if changes != 0 || notes[0].X != 0 {
		t.Errorf("a drag into the wall reported %d changes, note %+v", changes, notes[0])
	}
	cues(t, before)
	panics(t, `repeated "n1"`, func(c *ui.Context) {
		dup := []StickyNote{{ID: "n1", Text: "A"}, {ID: "n1", Text: "B"}}
		Whiteboard(c, &dup)
	})
}

func TestWorkloadLongNameKeepsLoad(t *testing.T) {
	people := []WorkloadPerson{{Name: "Alexandria Catherine Montgomery-Wellington", Role: "Design", Load: 1.15}}
	tt := ui.NewTester(narrow(func(c *ui.Context) { Workload(c, people, WorkloadOptions{}) }), 500, 300)
	settles(t, tt)
	if r, _ := tt.Find("115%"); r.W < 20 || r.X+r.W > 330 {
		t.Errorf("load at %v, cut or outside the card", r)
	}
}

func TestMindMapNodesClearCenter(t *testing.T) {
	added := -1
	branches := []string{"Offline", "Maps", "Sharing", "Widgets"}
	tt := ui.NewTester(narrow(func(c *ui.Context) { MindMap(c, "Atlas 2.0", branches, &added, MindMapOptions{}) }), 500, 500)
	settles(t, tt)
	center, _ := tt.Find("Atlas 2.0")
	for _, b := range branches {
		r, _ := tt.Find(b)
		if r.X < center.X+center.W && center.X < r.X+r.W && r.Y < center.Y+center.H && center.Y < r.Y+r.H {
			t.Errorf("%s at %v covers the center %v", b, r, center)
		}
	}
}

func TestTeamLoadEightWeeksStayReadable(t *testing.T) {
	weeks := []TeamWeek{{Person: "Maya Chen", Loads: []float64{1, 1, 1, 1, 1, 1, 1, 1}}}
	tt := ui.NewTester(narrow(func(c *ui.Context) { TeamLoad(c, weeks, TeamLoadOptions{}) }), 500, 300)
	settles(t, tt)
	a, _ := tt.Find("W1")
	b, _ := tt.Find("W2")
	if b.X-a.X < 48 {
		t.Errorf("week columns %v apart, too narrow to read", b.X-a.X)
	}
}

func TestChecklistIDsLiveApart(t *testing.T) {
	items := []ChecklistItem{{ID: "add", Title: "Write release notes"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Checklist(c, &items, ChecklistOptions{}) }), 500, 300)
	settles(t, tt)
	click(t, tt, "Add item")
	click(t, tt, "Toggle Write release notes")
	if len(items) != 2 || !items[0].Done {
		t.Errorf("items %v", items)
	}
}

func TestFirstItemKeepsFocus(t *testing.T) {
	var items []ChecklistItem
	tt := ui.NewTester(padded(func(c *ui.Context) { Checklist(c, &items, ChecklistOptions{}) }), 500, 300)
	settles(t, tt)
	click(t, tt, "Add item")
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if len(items) != 2 {
		t.Errorf("Space after the first add made %d items: Add lost focus", len(items))
	}
}
