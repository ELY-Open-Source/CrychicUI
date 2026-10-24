package agent

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestAgentAvatar(t *testing.T) {
	working := true
	tt := ui.NewTester(padded(func(c *ui.Context) { AgentAvatar(c, "Planner Bot", nil, working) }), 300, 200)
	tt.Frame()
	if !tt.HasText("PB") || !moves(tt, 300*time.Millisecond) {
		t.Errorf("working avatar %q does not spin", tt.Texts())
	}
	working = false
	tt.Frame()
	settles(t, tt, "idle avatar")
	panics(t, "no name", func() { render(func(c *ui.Context) { AgentAvatar(c, "", nil, false) }) })
}

func TestVoiceOrb(t *testing.T) {
	level := 0.2
	var accent ui.Color
	tt := ui.NewTester(padded(func(c *ui.Context) {
		accent = c.Theme().Accent
		VoiceOrb(c, level, VoiceOrbOptions{})
	}), 300, 300)
	tt.Advance(time.Second)
	quiet := inked(tt, accent)
	level = 1
	tt.Advance(time.Second)
	if loud := inked(tt, accent); loud <= quiet*3/2 {
		t.Errorf("core did not grow: %d then %d pixels", quiet, loud)
	}
	if !moves(tt, 300*time.Millisecond) {
		t.Error("orb does not breathe")
	}
	panics(t, "over one", func() { render(func(c *ui.Context) { VoiceOrb(c, 1.5, VoiceOrbOptions{}) }) })
}

// inked counts pixels near ink.
func inked(tt *ui.Tester, ink ui.Color) int {
	img, want, n := tt.Image(), rgba(ink), 0
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if p := img.RGBAAt(x, y); near(p.R, want.R) && near(p.G, want.G) && near(p.B, want.B) {
				n++
			}
		}
	}
	return n
}

func TestMemories(t *testing.T) {
	facts := []Memory{{Fact: "Prefers metric units", Saved: "Saved 2 days ago"}, {Fact: "Works in Lisbon"}}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Memories(c, "Memory", &facts, MemoriesOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "memories")
	if !tt.HasText("2 remembered") || !tt.HasText("Saved 2 days ago") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Forget Prefers metric units")
	if len(facts) != 1 || facts[0].Fact != "Works in Lisbon" || changed != 1 {
		t.Fatalf("facts %v changed %d", facts, changed)
	}
	click(t, tt, "Forget Works in Lisbon")
	settles(t, tt, "empty memories")
	if !tt.HasText("Nothing remembered yet") {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Delete, sound.Delete)
	panics(t, "duplicate", func() {
		dup := []Memory{{Fact: "a"}, {Fact: "a"}}
		render(func(c *ui.Context) { Memories(c, "M", &dup, MemoriesOptions{}) })
	})
}

func TestAgentStatus(t *testing.T) {
	state := Working
	tt := ui.NewTester(padded(func(c *ui.Context) { AgentStatus(c, state, AgentStatusOptions{}) }), 300, 200)
	tt.Frame()
	if !tt.HasText("Working") || !moves(tt, 300*time.Millisecond) {
		t.Errorf("working status %q does not pulse", tt.Texts())
	}
	state = Blocked
	tt.Frame()
	settles(t, tt, "blocked status")
	if !tt.HasText("Blocked") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "unknown state", func() { render(func(c *ui.Context) { AgentStatus(c, runStateCount, AgentStatusOptions{}) }) })
}

func TestRunTimeline(t *testing.T) {
	steps := []RunStep{
		{Name: "Plan", Duration: 2 * time.Second, State: StepDone},
		{Name: "Search", Duration: 9 * time.Second, State: StepDone},
		{Name: "Write", Duration: 4 * time.Second, State: StepRunning},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { RunTimeline(c, "Run", steps, RunTimelineOptions{}) }), 500, 300)
	tt.Frame()
	if !tt.HasText("3 steps · 0:15") || !tt.HasText("0:09") || !moves(tt, 200*time.Millisecond) {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "zero total", func() {
		render(func(c *ui.Context) { RunTimeline(c, "R", []RunStep{{Name: "x"}}, RunTimelineOptions{}) })
	})
}

func TestAgentBoard(t *testing.T) {
	agents := []BoardAgent{{Name: "Planner", Task: "Outline", Status: Working}, {Name: "Writer", Status: Idle}}
	tt := ui.NewTester(padded(func(c *ui.Context) { AgentBoard(c, "Team", agents, AgentBoardOptions{}) }), 900, 400)
	tt.Advance(time.Second)
	working, _ := tt.Find("Planner")
	agents[0].Status = Done
	tt.Frame()
	tt.Advance(time.Second)
	done, _ := tt.Find("Planner")
	if done.X <= working.X {
		t.Errorf("card stayed in its column: %v then %v", working, done)
	}
	if !tt.HasText("None") || !tt.HasText("Outline") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "duplicate", func() {
		render(func(c *ui.Context) { AgentBoard(c, "T", []BoardAgent{{Name: "a"}, {Name: "a"}}, AgentBoardOptions{}) })
	})
}

func TestContextFiles(t *testing.T) {
	files := []ContextFile{{Name: "brief.pdf", Size: 248_000}, {Name: "main.go", Size: 1200}}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ContextFiles(c, "Files", &files, ContextFilesOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "files")
	before := baseline()
	click(t, tt, "Remove brief.pdf")
	if len(files) != 1 || files[0].Name != "main.go" || changed != 1 {
		t.Errorf("files %v changed %d", files, changed)
	}
	cues(t, before, sound.Delete)
	if fileGlyph("a.go") != "file-code" || fileGlyph("a.PNG") != "file-image" || fileGlyph("a.md") != "file-text" {
		t.Error("file glyphs")
	}
}

func TestSlashCommands(t *testing.T) {
	cmds := []Command{{Name: "/deploy", Description: "Ship to production", Icon: "rocket"}, {Name: "/review", Icon: "eye"}}
	chosen, query, changed := "", "/de", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if SlashCommands(c, "Commands", cmds, query, &chosen, SlashCommandsOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "commands")
	if tt.HasText("/review") || !tt.HasText("Ship to production") {
		t.Fatalf("filter %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "/deploy")
	if chosen != "/deploy" || changed != 1 {
		t.Errorf("chosen %q changed %d", chosen, changed)
	}
	cues(t, before, sound.Select)
	query = "/zz"
	tt.Frame()
	if !tt.HasText("No matching commands") {
		t.Errorf("empty %q", tt.Texts())
	}
}

func TestSystemPrompt(t *testing.T) {
	text, changed := "Be brief.", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if SystemPrompt(c, "System prompt", &text, SystemPromptOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "system prompt")
	if !tt.HasText("≈ 3 tokens") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "System prompt")
	tt.Type(" Cite sources.")
	tt.Frame()
	if text != "Be brief. Cite sources." || changed == 0 || !tt.HasText("≈ 6 tokens") {
		t.Errorf("text %q changed %d texts %q", text, changed, tt.Texts())
	}
	for i, cue := range sound.Recent() {
		if i >= 32-changed && cue != sound.Tick || len(before) != 32 {
			t.Errorf("cues %v after %d edits", sound.Recent(), changed)
			break
		}
	}
}

func TestTemperature(t *testing.T) {
	value, changed := 0.5, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Temperature(c, &value, TemperatureOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "temperature")
	if !tt.HasText("Balanced") || !tt.HasText("0.50") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Frame()
	if value != 1 || changed != 1 || !tt.HasText("1.00") {
		t.Fatalf("value %v changed %d texts %q", value, changed, tt.Texts())
	}
	settles(t, tt, "creative")
	if tt.HasText("Balanced") {
		t.Errorf("mood stayed: %q", tt.Texts())
	}
	cues(t, before, sound.Tick)
	panics(t, "over one", func() {
		v := 1.5
		render(func(c *ui.Context) { Temperature(c, &v, TemperatureOptions{}) })
	})
}
