package project

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestWhiteboardStoredPositions(t *testing.T) {
	notes := []StickyNote{{ID: "a", Text: "A", X: 0.6, Y: 0.2}, {ID: "b", Text: "B", X: 0.2, Y: 0.6}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Whiteboard(c, &notes) }), 600, 400)
	settles(t, tt)
	r, _ := tt.Find("A")
	if r.X < 200 {
		t.Errorf("stored X ignored: %v", r)
	}
	before := baseline()
	b, _ := tt.Find("B")
	x, y := b.X+b.W/2, b.Y+b.H/2
	tt.Press(x, y)
	tt.Move(x+40, y+20)
	tt.Release(x+40, y+20)
	tt.Frame()
	if notes[0].X != 0.6 || notes[0].Y != 0.2 || notes[1].X <= 0.2 || notes[1].Y <= 0.6 {
		t.Errorf("notes committed incorrectly: %+v", notes)
	}
	cues(t, before, sound.Drop)
}

func TestEpicCompletion(t *testing.T) {
	children := []Task{{ID: "a", Title: "A", Done: true}, {ID: "b", Title: "B"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Epic(c, "Epic", children, EpicOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("1/2 done") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestChecklistEdit(t *testing.T) {
	items := []ChecklistItem{{Title: "Old", Done: true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Checklist(c, &items, ChecklistOptions{}) }), 500, 220)
	settles(t, tt)
	r, _ := tt.Find("Old")
	tt.ClickAt(r.X+5, r.Y+5)
	tt.ClickAt(r.X+5, r.Y+5)
	tt.Type("New")
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if items[0].Title != "New" {
		t.Errorf("title %q", items[0].Title)
	}
}

func TestRisksPosition(t *testing.T) {
	risks := []Risk{{Title: "Delay", Likelihood: 1, Impact: 5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Risks(c, risks) }), 500, 400)
	settles(t, tt)
	first, _ := tt.Find("Delay")
	risks[0].Likelihood, risks[0].Impact = 5, 1
	tt.Frame()
	second, _ := tt.Find("Delay")
	if first.X == second.X || first.Y == second.Y {
		t.Errorf("risk did not move: %v %v", first, second)
	}
}

func TestQuickAddSuffix(t *testing.T) {
	title, created := "Ship UI tomorrow", false
	var due time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		QuickAdd(c, &title, &due, &created, QuickAddOptions{})
	}), 500, 180)
	settles(t, tt)
	click(t, tt, "Add")
	if !created || due.IsZero() {
		t.Errorf("created %v due %v", created, due)
	}
}

func TestVoteNegative(t *testing.T) {
	votes := []int{-1}
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	ui.Render(func(c *ui.Context) { Vote(c, []string{"Idea"}, &votes, VoteOptions{}) }, 500, 180, 1)
}

func TestStaticMotion(t *testing.T) {
	for name, fn := range map[string]func(c *ui.Context){
		"dependencies": func(c *ui.Context) {
			opened := ""
			Dependencies(c, []string{"A"}, nil, &opened, DependenciesOptions{})
		},
		"standup": func(c *ui.Context) {
			StandUp(c, []StandUpPerson{{Name: "Ana", Yesterday: "A", Today: "B"}}, StandUpOptions{})
		},
		"risks":    func(c *ui.Context) { Risks(c, []Risk{{Title: "Delay", Likelihood: 2, Impact: 3}}) },
		"changes":  func(c *ui.Context) { Changes(c, []Change{{Text: "Updated", Time: "now"}}) },
		"teamload": func(c *ui.Context) { TeamLoad(c, []TeamWeek{{Person: "Ana", Loads: []float64{0.5, 0.5, 0.5, 0.5}}}) },
	} {
		t.Run(name, func(t *testing.T) {
			tt := ui.NewTester(padded(func(c *ui.Context) { fn(c) }), 500, 250)
			tt.Frame()
			if !tt.FramePending() {
				t.Error("no arrival motion")
			}
		})
	}
}

func TestWhiteboardOversizedNote(t *testing.T) {
	notes := []StickyNote{{ID: "a", Text: strings.Repeat("long ", 200), X: 1, Y: 1}}
	tt := ui.NewTester(func(c *ui.Context) { Whiteboard(c, &notes) }, 100, 80)
	settles(t, tt)
	if notes[0].X < 0 || notes[0].Y < 0 || notes[0].X > 1 || notes[0].Y > 1 {
		t.Errorf("invalid position: %+v", notes[0])
	}
}

func TestWhiteboardInsertedNoteSettles(t *testing.T) {
	notes := []StickyNote{{ID: "a", Text: "A", X: 0, Y: 0}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Whiteboard(c, &notes) }), 600, 400)
	settles(t, tt)
	notes = append(notes, StickyNote{ID: "b", Text: "B", X: 1, Y: 1})
	tt.Frame()
	settles(t, tt)
	r, ok := tt.Find("B")
	if !ok || r.X+r.W > 600 || r.Y+r.H > 400 {
		t.Errorf("inserted note clipped: %v found %v notes %+v texts %q", r, ok, notes, tt.Texts())
	}
	notes[1].Text = strings.Repeat("long ", 40)
	notes[1].Y = 1
	tt.Frame()
	settles(t, tt)
	r, ok = tt.Find(notes[1].Text)
	if !ok || r.Y+r.H > 400 {
		t.Errorf("grown note clipped: %v found %v", r, ok)
	}
}

func TestWhiteboardCanvasShrink(t *testing.T) {
	notes := []StickyNote{{ID: "a", Text: "A", X: 0.6, Y: 0.2}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Whiteboard(c, &notes) }), 600, 400)
	settles(t, tt)
	tt.SetSize(200, 400)
	settles(t, tt)
	r, ok := tt.Find("A")
	if !ok || r.X+r.W > 200 {
		t.Errorf("note lost after shrink: %v found %v", r, ok)
	}
}

func TestMindMapRequestsNewBranch(t *testing.T) {
	added := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MindMap(c, "Center", []string{"Branch"}, &added, MindMapOptions{})
	}), 500, 250)
	settles(t, tt)
	click(t, tt, "Add branch")
	if added != 1 {
		t.Errorf("added index %d, want 1", added)
	}
}

func TestQuickAddDateUsesClockLocation(t *testing.T) {
	now := time.Date(2026, 10, 23, 9, 0, 0, 0, time.FixedZone("New York", -4*3600))
	due := parseQuickDue(now, "Ship 2026-10-24")
	if due.Location() != now.Location() || due.YearDay() != 297 || due.Hour() != 0 {
		t.Errorf("due %v in %v", due, due.Location())
	}
}

func TestOKRFractionalValues(t *testing.T) {
	results := []KeyResult{{Title: "Grow", Current: 0.1, Target: 0.2}}
	tt := ui.NewTester(padded(func(c *ui.Context) { OKR(c, "Objective", results) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("0.1/0.2") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestC50CopyOverrides(t *testing.T) {
	children := []Task{{ID: "a", Title: "A", Done: true}, {ID: "b", Title: "B"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Epic(c, "Epic", children, EpicOptions{Done: "%d/%d complete"}) }), 500, 220)
	if !tt.HasText("1/2 complete") {
		t.Errorf("texts %q", tt.Texts())
	}
	card := ui.NewTester(padded(func(c *ui.Context) {
		ProjectCard(c, "Crychic", []string{"Ana"}, 0.5, ProjectCardOptions{Members: "%d people"})
	}), 500, 220)
	if !card.HasText("1 people") {
		t.Errorf("texts %q", card.Texts())
	}
	people := []StandUpPerson{{Name: "Ana", Yesterday: "A", Today: "B", Blockers: "C"}}
	stand := ui.NewTester(padded(func(c *ui.Context) { StandUp(c, people, StandUpOptions{Yesterday: "Y", Today: "T", Blockers: "B"}) }), 500, 220)
	if !stand.HasText("Y: A") || !stand.HasText("T: B") || !stand.HasText("B: C") {
		t.Errorf("texts %q", stand.Texts())
	}
}
