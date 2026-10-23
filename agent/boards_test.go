package agent

import (
	"image/color"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestAgentAvatar(t *testing.T) {
	working := true
	tt := ui.NewTester(padded(func(c *ui.Context) { AgentAvatar(c, "Ana Reyes", nil, working) }), 300, 200)
	tt.Frame()
	if !moves(tt, 200*time.Millisecond) {
		t.Error("ring does not pulse")
	}
	working = false
	tt.Frame()
	settles(t, tt, "avatar")
	panics(t, "empty name", func() { render(func(c *ui.Context) { AgentAvatar(c, "", nil, false) }) })
}

func TestVoiceOrb(t *testing.T) {
	level := 0.3
	var acc color.RGBA
	tt := ui.NewTester(padded(func(c *ui.Context) {
		acc = rgba(c.Theme().Accent)
		VoiceOrb(c, level)
	}), 300, 200)
	settles(t, tt, "orb")
	r := find(t, tt, "Voice level")
	if !at(tt, r.X+r.W/2, r.Y+r.H/2, acc) {
		t.Error("orb not drawn")
	}
	level = 0.9
	tt.Frame()
	if !tt.FramePending() {
		t.Error("orb does not grow")
	}
	settles(t, tt, "orb")
	panics(t, "over one", func() { render(func(c *ui.Context) { VoiceOrb(c, 1.5) }) })
	panics(t, "negative", func() { render(func(c *ui.Context) { VoiceOrb(c, -0.5) }) })
}

func TestMemories(t *testing.T) {
	facts := []string{"Likes green tea", "Works from Duluth", "Reviews on Fridays"}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Memories(c, "Memories", &facts, MemoriesOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "memories")
	before := baseline()
	r := find(t, tt, "Likes green tea")
	tt.ClickAt(r.X+r.W+30, r.Y+r.H/2)
	tt.Frame()
	if len(facts) != 2 || facts[0] != "Works from Duluth" || changed != 1 {
		t.Fatalf("facts %v changed %d", facts, changed)
	}
	cues(t, before, sound.Delete)
	click(t, tt, "Forget")
	click(t, tt, "Forget")
	if len(facts) != 0 || !tt.HasText("Nothing remembered yet") {
		t.Errorf("facts %v texts %q", facts, tt.Texts())
	}
	panics(t, "empty fact", func() {
		f := []string{""}
		render(func(c *ui.Context) { Memories(c, "x", &f, MemoriesOptions{}) })
	})
}

func TestAgentStatus(t *testing.T) {
	state := Working
	tt := ui.NewTester(padded(func(c *ui.Context) { AgentStatus(c, state, AgentStatusOptions{}) }), 300, 200)
	tt.Frame()
	if !tt.HasText("Working") {
		t.Errorf("texts %q", tt.Texts())
	}
	if !moves(tt, 200*time.Millisecond) {
		t.Error("working dot does not pulse")
	}
	state = Done
	tt.Frame()
	settles(t, tt, "status")
	if !tt.HasText("Done") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "unknown state", func() { render(func(c *ui.Context) { AgentStatus(c, 9, AgentStatusOptions{}) }) })
}

func TestRunTimeline(t *testing.T) {
	steps := []RunStep{{Name: "Search", Duration: 2 * time.Second}, {Name: "Read", Duration: 500 * time.Millisecond}, {Name: "Write", Duration: 3 * time.Second}}
	tt := ui.NewTester(padded(func(c *ui.Context) { RunTimeline(c, "Run", steps) }), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bars do not grow")
	}
	settles(t, tt, "timeline")
	if !tt.HasText("0:02") || !tt.HasText("0:00") || !tt.HasText("0:03") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "all zero", func() {
		render(func(c *ui.Context) { RunTimeline(c, "x", []RunStep{{Name: "a"}}) })
	})
	panics(t, "negative duration", func() {
		render(func(c *ui.Context) { RunTimeline(c, "x", []RunStep{{Name: "a", Duration: -time.Second}}) })
	})
}

func TestAgentBoard(t *testing.T) {
	agents := []BoardAgent{{Name: "Alpha", Status: Idle}, {Name: "Beta", Status: Working}}
	tt := ui.NewTester(padded(func(c *ui.Context) { AgentBoard(c, "Board", agents, AgentBoardOptions{}) }), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("cards do not arrive")
	}
	tt.Advance(2 * time.Second)
	a, b := find(t, tt, "Alpha"), find(t, tt, "Beta")
	if b.X >= a.X {
		t.Errorf("working Beta at %v should lead idle Alpha at %v", b.X, a.X)
	}
	agents[1].Status = Done
	tt.Frame()
	if !tt.FramePending() {
		t.Error("cards do not move")
	}
	settles(t, tt, "board")
	a, b = find(t, tt, "Alpha"), find(t, tt, "Beta")
	if a.X >= b.X {
		t.Errorf("after the change Alpha at %v should lead Beta at %v", a.X, b.X)
	}
	own := ui.NewTester(padded(func(c *ui.Context) {
		AgentBoard(c, "Tableau", []BoardAgent{{Name: "A", Status: Idle}}, AgentBoardOptions{Status: AgentStatusOptions{Idle: "Inactif"}})
	}), 500, 300)
	own.Frame()
	if !own.HasText("Inactif") {
		t.Errorf("texts %q", own.Texts())
	}
	panics(t, "duplicate names", func() {
		render(func(c *ui.Context) { AgentBoard(c, "x", []BoardAgent{{Name: "a"}, {Name: "a"}}, AgentBoardOptions{}) })
	})
}

func TestContextFiles(t *testing.T) {
	files := []string{"deck.md", "brief.pdf", "notes.txt"}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ContextFiles(c, "Files", &files, ContextFilesOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("chips do not pop in")
	}
	settles(t, tt, "chips")
	before := baseline()
	click(t, tt, "Remove")
	if len(files) != 2 || files[0] != "brief.pdf" || changed != 1 {
		t.Fatalf("files %v changed %d", files, changed)
	}
	click(t, tt, "Remove")
	click(t, tt, "Remove")
	if len(files) != 0 || changed != 3 {
		t.Errorf("files %v changed %d after removing all", files, changed)
	}
	cues(t, before, sound.Delete, sound.Delete, sound.Delete)
	panics(t, "empty file", func() {
		f := []string{""}
		render(func(c *ui.Context) { ContextFiles(c, "x", &f, ContextFilesOptions{}) })
	})
}

func TestSlashCommands(t *testing.T) {
	commands := []string{"/deploy", "/review", "/test"}
	chosen, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if SlashCommands(c, "Commands", commands, &chosen).Changed() {
			changed++
		}
	}), 400, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("list does not rise in")
	}
	settles(t, tt, "commands")
	before := baseline()
	click(t, tt, "/deploy")
	if chosen != 0 || changed != 1 {
		t.Fatalf("chosen %d changed %d", chosen, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "no slash", func() {
		render(func(c *ui.Context) { SlashCommands(c, "x", []string{"deploy"}, &chosen) })
	})
}

func TestSystemPrompt(t *testing.T) {
	text, changed := "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if SystemPrompt(c, "System prompt", &text, SystemPromptOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "editor")
	if !tt.HasText("≈ 0 tokens") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "System prompt")
	tt.Type("You write launch decks.")
	tt.Frame()
	if changed == 0 || !tt.HasText("≈ 6 tokens") {
		t.Errorf("changed %d texts %q", changed, tt.Texts())
	}
	cues(t, before, sound.Tick)
	panics(t, "nil text", func() { render(func(c *ui.Context) { SystemPrompt(c, "x", nil, SystemPromptOptions{}) }) })
}

func TestTemperature(t *testing.T) {
	value, changed := 0.5, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Temperature(c, &value, TemperatureOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "temperature")
	if !tt.HasText("Precise") || !tt.HasText("Creative") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Frame()
	if value != 1 || changed != 1 {
		t.Fatalf("value %v changed %d", value, changed)
	}
	cues(t, before, sound.Tick)
	panics(t, "over one", func() {
		v := 1.5
		render(func(c *ui.Context) { Temperature(c, &v, TemperatureOptions{}) })
	})
}
