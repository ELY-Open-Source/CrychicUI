package game

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestDiceRoll(t *testing.T) {
	faces := []int{1, 2}
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Dice(c, &faces, DiceOptions{Rand: fixed(5)}).Changed() {
			changes++
		}
	}), 400, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Roll")
	cues(t, before, sound.Drop)
	settles(t, tt)
	if !slices.Equal(faces, []int{6, 6}) || changes != 1 || !exact(tt, "Total 12") {
		t.Errorf("faces %v, changes %d: %q", faces, changes, tt.Texts())
	}
	func() {
		defer func() {
			if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), "Rand gave 6") {
				t.Errorf("panic %v, want a bad Rand named", r)
			}
		}()
		f := []int{1}
		bad := ui.NewTester(padded(func(c *ui.Context) { Dice(c, &f, DiceOptions{Rand: fixed(6)}) }), 400, 300)
		bad.Frame()
		_ = bad.Click("Roll")
		bad.Frame()
	}()
	panics(t, "shows 7", func(c *ui.Context) { f := []int{7}; Dice(c, &f, DiceOptions{}) })
}

func TestCoinFlip(t *testing.T) {
	heads := true
	tt := ui.NewTester(padded(func(c *ui.Context) { Coin(c, &heads, CoinOptions{Rand: fixed(1)}) }), 400, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Flip")
	cues(t, before, sound.Pop)
	settles(t, tt)
	if heads || !exact(tt, "Tails 1") || !exact(tt, "Heads 0") || !exact(tt, "Tails") {
		t.Errorf("heads %v: %q", heads, tt.Texts())
	}
}

func TestSpinWheelTicks(t *testing.T) {
	picked := -1
	slices3 := []string{"Gold", "Gem", "Shield", "Key"}
	tt := ui.NewTester(padded(func(c *ui.Context) { SpinWheel(c, slices3, &picked, SpinWheelOptions{Rand: fixed(2)}) }), 400, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Spin")
	for range 200 {
		tt.Advance(16 * time.Millisecond)
	}
	ticks := 0
	for _, q := range heard(before) {
		if q == sound.Tick {
			ticks++
		}
	}
	// Four turns pass at least sixteen slices.
	if picked != 2 || ticks < 16 || !exact(tt, "Landed on Shield") {
		t.Errorf("picked %d, %d ticks: %q", picked, ticks, tt.Texts())
	}
	panics(t, "2 to 12 slices", func(c *ui.Context) { p := -1; SpinWheel(c, []string{"a"}, &p, SpinWheelOptions{}) })
}

func TestScoreboardAdds(t *testing.T) {
	teams := [2]Team{{Name: "Comets", Score: 3}, {Name: "Giants", Score: math.MaxInt}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Scoreboard(c, &teams, ScoreboardOptions{}) }), 420, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add a point to Comets")
	_ = tt.Click("Add a point to Giants")
	tt.Frame()
	cues(t, before, sound.Pop)
	if teams[0].Score != 4 || teams[1].Score != math.MaxInt {
		t.Errorf("scores %+v", teams)
	}
	panics(t, "two teams named", func(c *ui.Context) {
		s := [2]Team{{Name: "A"}, {Name: "A"}}
		Scoreboard(c, &s, ScoreboardOptions{})
	})
}

func TestLeaderboardMoves(t *testing.T) {
	players := []Player{{Name: "Priya", Score: 9, Prev: 3}, {Name: "Maya", Score: 8, Prev: 1}, {Name: "Sofia", Score: 5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Leaderboard(c, players, LeaderboardOptions{}) }), 480, 400)
	settles(t, tt)
	if !exact(tt, "2") || !exact(tt, "New") {
		t.Errorf("moves: %q", tt.Texts())
	}
	panics(t, "high to low", func(c *ui.Context) {
		Leaderboard(c, []Player{{Name: "a", Score: 1}, {Name: "b", Score: 2}}, LeaderboardOptions{})
	})
}

func TestHealthBarTrail(t *testing.T) {
	hp := 96
	tt := ui.NewTester(padded(func(c *ui.Context) { HealthBar(c, hp, 120, HealthBarOptions{}) }), 420, 120)
	settles(t, tt)
	r, _ := tt.Find("HP")
	x, y := int(r.X+r.W*0.65), int(r.Y+r.H-10)
	calm := tt.Image().RGBAAt(x, y)
	hp = 60
	tt.Frame()
	tt.Advance(150 * time.Millisecond)
	trail := tt.Image().RGBAAt(x, y)
	settles(t, tt)
	drained := tt.Image().RGBAAt(x, y)
	if trail == drained || trail == calm {
		t.Errorf("no trail between 50%% and 80%%: before %v, during %v, after %v", calm, trail, drained)
	}
	if !exact(tt, "60 / 120") {
		t.Errorf("count: %q", tt.Texts())
	}
	panics(t, "0 <= hp <= max", func(c *ui.Context) { HealthBar(c, 5, 4, HealthBarOptions{}) })
}

func TestXPBarLevelUpOnce(t *testing.T) {
	level, xp := 11, int64(2900)
	tt := ui.NewTester(padded(func(c *ui.Context) { XPBar(c, level, xp, 3000, XPBarOptions{}) }), 440, 200)
	settles(t, tt)
	before := baseline()
	level, xp = 12, 500
	for range 5 {
		tt.Frame()
	}
	settles(t, tt)
	cues(t, before, sound.Success)
	if !exact(tt, "Level 12") || !exact(tt, "500 / 3,000 XP") {
		t.Errorf("xp: %q", tt.Texts())
	}
	panics(t, "0 <= xp < need", func(c *ui.Context) { XPBar(c, 1, 10, 10, XPBarOptions{}) })
}

func TestInventoryDragSwaps(t *testing.T) {
	slots := []Item{{Name: "Sword", Icon: "sword", Count: 1}, {}, {Name: "Gem", Icon: "gem", Count: 3, Rarity: Rare}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Inventory(c, &slots, InventoryOptions{Columns: 3}) }), 420, 300)
	settles(t, tt)
	a, _ := tt.Find("Sword")
	b, _ := tt.Find("Empty slot 2")
	before := baseline()
	drag(tt, a, b)
	cues(t, before, sound.Drop)
	if slots[0].Name != "" || slots[1].Name != "Sword" || !exact(tt, "2 / 3") {
		t.Errorf("slots %+v: %q", slots, tt.Texts())
	}
	if _, ok := tt.Find("Gem × 3"); !ok {
		t.Error("a stack is not named with its count")
	}
	panics(t, "count of at least 1", func(c *ui.Context) {
		s := []Item{{Name: "x", Icon: "gem"}}
		Inventory(c, &s, InventoryOptions{})
	})
}

func TestItemTipStats(t *testing.T) {
	it := ItemDetail{Name: "Ember sword", Kind: "Sword", Icon: "sword", Rarity: Epic, Level: 12,
		Stats: []Stat{{Name: "Damage", Value: 48, Delta: 12}, {Name: "Durability", Value: 80, Delta: -20}, {Name: "Speed", Value: 1}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ItemTip(c, it, ItemTipOptions{}) }), 400, 400)
	settles(t, tt)
	if !exact(tt, "+12") || !exact(tt, "-20") || exact(tt, "+0") || !exact(tt, "Epic") || !exact(tt, "Requires level 12") {
		t.Errorf("tip: %q", tt.Texts())
	}
	panics(t, "repeated", func(c *ui.Context) {
		ItemTip(c, ItemDetail{Name: "a", Kind: "k", Icon: "gem", Stats: []Stat{{Name: "x"}, {Name: "x"}}}, ItemTipOptions{})
	})
}

func TestAchievementsCount(t *testing.T) {
	list := []Achievement{{Name: "First", Icon: "star", Unlocked: true}, {Name: "Collector", Icon: "gem", Goal: 50, Progress: 34}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Achievements(c, list, AchievementsOptions{}) }), 460, 300)
	settles(t, tt)
	if !exact(tt, "1 of 2 unlocked") || !exact(tt, "34 / 50") {
		t.Errorf("achievements: %q", tt.Texts())
	}
	panics(t, "0 <= progress <= goal", func(c *ui.Context) {
		Achievements(c, []Achievement{{Name: "x", Icon: "gem", Goal: 1, Progress: 2}}, AchievementsOptions{})
	})
}

func TestQuestsTrack(t *testing.T) {
	quests := []Quest{{Title: "Drake", Steps: []QuestStep{{Text: "a", Done: true}, {Text: "b"}}},
		{Title: "Gift", Steps: []QuestStep{{Text: "c", Done: true}}}}
	tracked := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Quests(c, quests, &tracked, QuestsOptions{}) }), 460, 500)
	settles(t, tt)
	if !exact(tt, "1 of 2") || !exact(tt, "Complete") {
		t.Errorf("quests: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Track Drake")
	click(t, tt, "Stop tracking Drake")
	cues(t, before, sound.Select, sound.Select)
	if tracked != "" {
		t.Errorf("tracked %q", tracked)
	}
	panics(t, "needs steps", func(c *ui.Context) { tr := ""; Quests(c, []Quest{{Title: "x"}}, &tr, QuestsOptions{}) })
}

func TestTurnTimerWarnsOnce(t *testing.T) {
	var start time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if start.IsZero() {
			start = c.Now()
		}
		TurnTimer(c, "Maya", start.Add(30*time.Second), time.Minute, TurnTimerOptions{})
	}), 420, 200)
	tt.Frame()
	if !exact(tt, "0:30") {
		t.Fatalf("start: %q", tt.Texts())
	}
	before := baseline()
	tt.Advance(25 * time.Second)
	cues(t, before, sound.Warning)
	tt.Advance(6 * time.Second)
	if !exact(tt, "Time's up") {
		t.Errorf("after deadline: %q", tt.Texts())
	}
	cues(t, before, sound.Warning)
}

func TestHandPicks(t *testing.T) {
	cards := []Card{{Rank: "9", Suit: SuitHearts}, {Rank: "K", Suit: SuitSpades}}
	picked := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Hand(c, cards, &picked, HandOptions{}) }), 420, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "K of spades")
	if picked != 1 {
		t.Fatalf("picked %d", picked)
	}
	click(t, tt, "K of spades")
	cues(t, before, sound.Swipe, sound.Swipe)
	if picked != -1 {
		t.Errorf("second click left %d", picked)
	}
	swapUnderPress(t, tt, "9 of hearts", func() { cards[0], cards[1] = cards[1], cards[0] })
	if picked == 0 {
		t.Error("a press on the 9 picked the king after the swap")
	}
	panics(t, "repeated", func(c *ui.Context) { p := -1; Hand(c, []Card{{Rank: "A"}, {Rank: "A"}}, &p, HandOptions{}) })
}

func TestPlayingCardFlips(t *testing.T) {
	up := false
	tt := ui.NewTester(padded(func(c *ui.Context) { PlayingCard(c, Card{Rank: "A", Suit: SuitSpades}, &up, PlayingCardOptions{}) }), 300, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Face-down card")
	cues(t, before, sound.Swipe)
	tt.Advance(100 * time.Millisecond)
	if exact(tt, "A") {
		t.Error("the face shows before the card turns edge-on")
	}
	settles(t, tt)
	if !up || !exact(tt, "A") {
		t.Errorf("up %v: %q", up, tt.Texts())
	}
	if _, ok := tt.Find("A of spades"); !ok {
		t.Error("the face-up card is not named")
	}
}

func TestTicTacToeWins(t *testing.T) {
	board := [9]Mark{MarkX, Blank, MarkO, Blank, MarkX, Blank, Blank, MarkO, Blank}
	tt := ui.NewTester(padded(func(c *ui.Context) { TicTacToe(c, &board, TicTacToeOptions{}) }), 400, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Cell 9")
	cues(t, before, sound.Success)
	settles(t, tt)
	if board[8] != MarkX || !exact(tt, "X wins") {
		t.Errorf("board %v: %q", board, tt.Texts())
	}
	if _, ok := tt.Find("Cell 1, X"); !ok {
		t.Error("a taken cell is not named with its mark")
	}
	click(t, tt, "New game")
	if board != [9]Mark{} {
		t.Errorf("reset left %v", board)
	}
	panics(t, "X moves first", func(c *ui.Context) { b := [9]Mark{MarkO}; TicTacToe(c, &b, TicTacToeOptions{}) })
}
