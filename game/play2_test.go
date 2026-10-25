package game

import (
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestMemoryPairs(t *testing.T) {
	faces := []string{"cat", "moon", "cat", "moon"}
	matched := make([]bool, 4)
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Memory(c, faces, &matched, MemoryOptions{}).Changed() {
			changes++
		}
	}), 400, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Card 1")
	click(t, tt, "Card 2")
	if !exact(tt, "1 moves") || matched[0] {
		t.Fatalf("mismatch: %q", tt.Texts())
	}
	tt.Advance(time.Second)
	if _, ok := tt.Find("Card 1"); !ok {
		t.Fatal("a missed pair stayed face up")
	}
	click(t, tt, "Card 1")
	click(t, tt, "Card 3")
	click(t, tt, "Card 2")
	click(t, tt, "Card 4")
	cues(t, before, sound.Swipe, sound.Swipe, sound.Swipe, sound.Success, sound.Swipe, sound.Success)
	settles(t, tt)
	if !slices.Equal(matched, []bool{true, true, true, true}) || changes != 2 || !exact(tt, "All pairs in 3 moves") {
		t.Errorf("matched %v, changes %d: %q", matched, changes, tt.Texts())
	}
	panics(t, "appears 3 times", func(c *ui.Context) {
		m := make([]bool, 4)
		Memory(c, []string{"a", "a", "a", "b"}, &m, MemoryOptions{})
	})
}

func TestSlotMachineStops(t *testing.T) {
	reels := [3]int{0, 1, 2}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SlotMachine(c, []string{"cherry", "bell", "gem"}, &reels, SlotMachineOptions{Rand: fixed(1)})
	}), 400, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Spin")
	if err := tt.Click("Spin"); err == nil {
		tt.Frame()
	}
	for range 150 {
		tt.Advance(16 * time.Millisecond)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	if reels != [3]int{1, 1, 1} || !exact(tt, "Three of a kind!") {
		t.Errorf("reels %v: %q", reels, tt.Texts())
	}
}

func TestHeartsBreak(t *testing.T) {
	lives := 3
	tt := ui.NewTester(padded(func(c *ui.Context) { Hearts(c, lives, 5, HeartsOptions{}) }), 300, 120)
	settles(t, tt)
	if tt.FramePending() {
		t.Fatal("hearts at rest still animate")
	}
	if _, ok := tt.Find("3 of 5 lives"); !ok {
		t.Errorf("not named: %q", tt.Texts())
	}
	lives = 2
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("a lost heart did not break")
	}
	settles(t, tt)
	if tt.FramePending() {
		t.Error("the break never ends")
	}
	panics(t, "0 <= lives <= max", func(c *ui.Context) { Hearts(c, 6, 5, HeartsOptions{}) })
}

func TestComboBest(t *testing.T) {
	count := 9
	tt := ui.NewTester(padded(func(c *ui.Context) { Combo(c, count, ComboOptions{}) }), 400, 200)
	settles(t, tt)
	count = 12
	tt.Frame()
	if !tt.FramePending() {
		t.Error("a rising combo did not bounce")
	}
	count = 0
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "Best 12") || !exact(tt, "×0") {
		t.Errorf("combo: %q", tt.Texts())
	}
	panics(t, "is negative", func(c *ui.Context) { Combo(c, -1, ComboOptions{}) })
}

func TestLevelPickPicks(t *testing.T) {
	levels := []Level{{Name: "1-1", Stars: 3}, {Name: "1-2", Stars: 1}, {Name: "1-3", Locked: true}}
	picked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { LevelPick(c, levels, &picked, LevelPickOptions{}) }), 420, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Level 1-2")
	_ = tt.Click("Level 1-3, locked")
	tt.Frame()
	cues(t, before, sound.Select)
	if picked != 1 || !exact(tt, "4 / 9") {
		t.Errorf("picked %d: %q", picked, tt.Texts())
	}
	swapUnderPress(t, tt, "Level 1-1", func() { levels[0], levels[1] = levels[1], levels[0] })
	if picked == 0 {
		t.Error("a press on 1-1 picked the level that slid under it")
	}
	panics(t, "none while locked", func(c *ui.Context) {
		p := -1
		LevelPick(c, []Level{{Name: "x", Stars: 2, Locked: true}}, &p, LevelPickOptions{})
	})
}

func TestCharacterCard(t *testing.T) {
	h := Hero{Name: "Priya Nair", Class: "Ember mage", Level: 14, HP: 84, MaxHP: 120, MP: 62, MaxMP: 90,
		Stats: []Stat{{Name: "Intellect", Value: 22}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Character(c, h, CharacterOptions{}) }), 400, 400)
	settles(t, tt)
	if !exact(tt, "Lv 14") || !exact(tt, "84 / 120") || !exact(tt, "62 / 90") || !exact(tt, "22") {
		t.Errorf("character: %q", tt.Texts())
	}
	panics(t, "0 <= mp <= max", func(c *ui.Context) {
		Character(c, Hero{Name: "a", Class: "b", Level: 1, MaxHP: 1, MP: 3, MaxMP: 2}, CharacterOptions{})
	})
}

func TestStatPointsFloor(t *testing.T) {
	stats := []Stat{{Name: "Strength", Value: 9}}
	points := 1
	tt := ui.NewTester(padded(func(c *ui.Context) { StatPoints(c, &stats, &points, StatPointsOptions{}) }), 420, 200)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Raise Strength")
	_ = tt.Click("Raise Strength")
	tt.Frame()
	click(t, tt, "Lower Strength")
	_ = tt.Click("Lower Strength")
	tt.Frame()
	cues(t, before, sound.Tick, sound.Tick)
	if stats[0].Value != 9 || points != 1 {
		t.Errorf("value %d, points %d", stats[0].Value, points)
	}
}

func TestSkillTreeUnlocks(t *testing.T) {
	skills := []Skill{{ID: "spark", Name: "Spark", Icon: "zap", Cost: 1},
		{ID: "fire", Name: "Fireball", Icon: "flame", Requires: []string{"spark"}, Cost: 1},
		{ID: "inferno", Name: "Inferno", Icon: "sun", Requires: []string{"fire"}, Cost: 1}}
	learned := []string{"spark"}
	points := 1
	tt := ui.NewTester(padded(func(c *ui.Context) { SkillTree(c, skills, &learned, &points, SkillTreeOptions{}) }), 460, 500)
	settles(t, tt)
	before := baseline()
	_ = tt.Click("Learn Inferno")
	tt.Frame()
	click(t, tt, "Learn Fireball")
	_ = tt.Click("Learn Inferno")
	tt.Frame()
	cues(t, before, sound.Success)
	if !slices.Equal(learned, []string{"spark", "fire"}) || points != 0 {
		t.Errorf("learned %v, points %d", learned, points)
	}
	panics(t, "form a loop", func(c *ui.Context) {
		l, p := []string{}, 0
		SkillTree(c, []Skill{{ID: "a", Name: "A", Icon: "zap", Cost: 1, Requires: []string{"b"}},
			{ID: "b", Name: "B", Icon: "zap", Cost: 1, Requires: []string{"a"}}}, &l, &p, SkillTreeOptions{})
	})
	panics(t, "unknown skill", func(c *ui.Context) {
		l, p := []string{}, 0
		SkillTree(c, []Skill{{ID: "a", Name: "A", Icon: "zap", Cost: 1, Requires: []string{"z"}}}, &l, &p, SkillTreeOptions{})
	})
}

func TestLootOpens(t *testing.T) {
	opened := false
	items := []Item{{Name: "Moonstone", Icon: "gem", Count: 3, Rarity: Rare}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Loot(c, items, &opened, LootOptions{}) }), 400, 400)
	settles(t, tt)
	if exact(tt, "Moonstone") {
		t.Fatal("items show before the chest opens")
	}
	before := baseline()
	click(t, tt, "Open chest")
	cues(t, before, sound.Pop)
	settles(t, tt)
	if !opened || !exact(tt, "Moonstone") || !exact(tt, "× 3") {
		t.Errorf("opened %v: %q", opened, tt.Texts())
	}
}

func TestDailyRewardClaims(t *testing.T) {
	rewards := []Reward{{Icon: "coins", Amount: "100"}, {Icon: "gem", Amount: "Gem"}}
	claimed := 1
	tt := ui.NewTester(padded(func(c *ui.Context) { DailyReward(c, rewards, &claimed, DailyRewardOptions{}) }), 480, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Claim")
	_ = tt.Click("All claimed")
	tt.Frame()
	cues(t, before, sound.Success)
	if claimed != 2 || !exact(tt, "2-day streak") || !exact(tt, "All claimed") {
		t.Errorf("claimed %d: %q", claimed, tt.Texts())
	}
}

func TestLobbySeats(t *testing.T) {
	players := []Seat{{Name: "Maya", Ready: true, Host: true}, {Name: "Priya"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Lobby(c, players, 4, LobbyOptions{}) }), 440, 400)
	settles(t, tt)
	waiting := 0
	for _, s := range tt.Texts() {
		if s == "Waiting for a player" {
			waiting++
		}
	}
	if waiting != 2 || !exact(tt, "1 of 4 ready") || !exact(tt, "Host") {
		t.Errorf("lobby: %q", tt.Texts())
	}
	panics(t, "2 hosts", func(c *ui.Context) {
		Lobby(c, []Seat{{Name: "a", Host: true}, {Name: "b", Host: true}}, 2, LobbyOptions{})
	})
}

func TestReadyCheckToggles(t *testing.T) {
	others := []Seat{{Name: "Maya", Ready: true}}
	ready := false
	tt := ui.NewTester(padded(func(c *ui.Context) { ReadyCheck(c, others, &ready, ReadyCheckOptions{}) }), 440, 300)
	settles(t, tt)
	if !exact(tt, "1 of 2 ready") {
		t.Fatalf("start: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Ready")
	cues(t, before, sound.Select)
	settles(t, tt)
	if !ready || !exact(tt, "Everyone is ready") {
		t.Errorf("ready %v: %q", ready, tt.Texts())
	}
}

func TestChessClockRuns(t *testing.T) {
	left := [2]time.Duration{5 * time.Minute, 12 * time.Second}
	turn := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { ChessClock(c, &left, &turn, ChessClockOptions{}) }), 420, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "End White's move")
	if turn != 1 {
		t.Fatalf("turn %d", turn)
	}
	tt.Advance(2500 * time.Millisecond)
	if !exact(tt, "0:09.5") {
		t.Errorf("black after 2.5 s: %q", tt.Texts())
	}
	click(t, tt, "End Black's move")
	cues(t, before, sound.Tick, sound.Tick)
	if turn != 0 || left[1] > 9600*time.Millisecond || left[1] < 9400*time.Millisecond {
		t.Errorf("turn %d, black left %v", turn, left[1])
	}
	turn = 1
	tt.Frame()
	tt.Advance(11 * time.Second)
	if !exact(tt, "Time") {
		t.Errorf("flag: %q", tt.Texts())
	}
}

func TestSlidePuzzleSolves(t *testing.T) {
	tiles := []int{1, 2, 3, 4, 5, 6, 0, 7, 8}
	tt := ui.NewTester(padded(func(c *ui.Context) { SlidePuzzle(c, &tiles, SlidePuzzleOptions{}) }), 400, 400)
	settles(t, tt)
	before := baseline()
	_ = tt.Click("Tile 1")
	tt.Frame()
	click(t, tt, "Tile 7")
	click(t, tt, "Tile 8")
	cues(t, before, sound.Tick, sound.Success)
	settles(t, tt)
	if !slices.Equal(tiles, []int{1, 2, 3, 4, 5, 6, 7, 8, 0}) || !exact(tt, "Solved in 2 moves") {
		t.Errorf("tiles %v: %q", tiles, tt.Texts())
	}
	panics(t, "cannot be solved", func(c *ui.Context) {
		s := []int{2, 1, 3, 4, 5, 6, 7, 8, 0}
		SlidePuzzle(c, &s, SlidePuzzleOptions{})
	})
}

func TestJoystickSprings(t *testing.T) {
	var x, y float32
	tt := ui.NewTester(padded(func(c *ui.Context) { Joystick(c, &x, &y, JoystickOptions{}) }), 300, 300)
	settles(t, tt)
	r, _ := tt.Find("Joystick")
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	before := baseline()
	tt.Press(cx, cy)
	tt.Move(cx+200, cy)
	tt.Frame()
	if x < 0.99 || x > 1 || y != 0 {
		t.Errorf("held far right: x %v y %v", x, y)
	}
	tt.Release(cx+200, cy)
	tt.Frame()
	cues(t, before, sound.Tick)
	if x != 0 || y != 0 {
		t.Errorf("released at %v, %v", x, y)
	}
}
