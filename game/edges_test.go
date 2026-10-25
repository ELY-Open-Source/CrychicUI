package game

import (
	"fmt"
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func rectOf(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("%q missing: %q", name, tt.Texts())
	}
	return r
}

func TestLobbySeatKeyName(t *testing.T) {
	ui.Render(func(c *ui.Context) {
		Lobby(c, []Seat{{Name: "game.emptySeat:1", Host: true}}, 2, LobbyOptions{})
	}, 400, 400, 1)
}

func TestChessClockKeepsUsedTime(t *testing.T) {
	left := [2]time.Duration{5 * time.Minute, 5 * time.Minute}
	turn := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { ChessClock(c, &left, &turn, ChessClockOptions{}) }), 420, 300)
	tt.Advance(10 * time.Second)
	turn = -1
	tt.Frame()
	if left[0] != 4*time.Minute+50*time.Second || !exact(tt, "4:50") {
		t.Errorf("paused by the caller: %v, %q", left[0], tt.Texts())
	}
}

func TestSpinWheelFollowsPick(t *testing.T) {
	names := []string{"Gold", "Silver", "Bronze", "Iron"}
	picked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { SpinWheel(c, names, &picked, SpinWheelOptions{}) }), 400, 500)
	settles(t, tt)
	picked = 2
	tt.Frame()
	settles(t, tt)
	top, at := "", float32(math.Inf(1))
	for _, n := range names {
		if r := rectOf(t, tt, n); r.Y < at {
			top, at = n, r.Y
		}
	}
	if top != "Bronze" {
		t.Errorf("%s sits under the pointer, want Bronze", top)
	}
	names = append(names, "Potion", "Sword")
	picked = 1
	tt.Frame()
	settles(t, tt)
	top, at = "", float32(math.Inf(1))
	for _, n := range names {
		if r := rectOf(t, tt, n); r.Y < at {
			top, at = n, r.Y
		}
	}
	if top != "Silver" {
		t.Errorf("after two more slices %s sits under the pointer, want Silver", top)
	}
	names = names[:4]
	tt.Frame()
	settles(t, tt)
	top, at = "", float32(math.Inf(1))
	for _, n := range names {
		if r := rectOf(t, tt, n); r.Y < at {
			top, at = n, r.Y
		}
	}
	if top != "Silver" {
		t.Errorf("with the pick kept and the count changed, %s sits under the pointer, want Silver", top)
	}
}

func TestStatPointsRefundCap(t *testing.T) {
	stats := []Stat{{Name: "Strength", Value: 5}}
	points := 1
	tt := ui.NewTester(padded(func(c *ui.Context) { StatPoints(c, &stats, &points, StatPointsOptions{}) }), 400, 300)
	settles(t, tt)
	click(t, tt, "Raise Strength")
	points = math.MaxInt
	tt.Frame()
	_ = tt.Click("Lower Strength")
	tt.Frame()
	if points != math.MaxInt || stats[0].Value != 6 {
		t.Errorf("refund past MaxInt: points %d, value %d", points, stats[0].Value)
	}
}

func TestTurnTimerHugeTurn(t *testing.T) {
	total := time.Duration(math.MaxInt64)
	before := baseline()
	var end time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if end.IsZero() {
			end = c.Now().Add(total)
		}
		TurnTimer(c, "Maya Chen", end, total, TurnTimerOptions{})
	}), 400, 300)
	tt.Frame()
	if got := heard(before); slices.Contains(got, sound.Warning) {
		t.Errorf("a 292-year turn warned at once: %v", got)
	}
}

func TestSlotMachineSpinKeepsFocus(t *testing.T) {
	reels := [3]int{0, 0, 0}
	symbols := []string{"cherry", "citrus", "grape", "bell"}
	tt := ui.NewTester(padded(func(c *ui.Context) { SlotMachine(c, symbols, &reels, SlotMachineOptions{Rand: fixed(3)}) }), 400, 400)
	settles(t, tt)
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	reels = [3]int{0, 1, 0}
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if reels != [3]int{3, 3, 3} {
		t.Errorf("Space after the win notice left: reels %v", reels)
	}
}

func TestPlayingCardInterruptAndStill(t *testing.T) {
	up := false
	tt := ui.NewTester(padded(func(c *ui.Context) { PlayingCard(c, Card{Rank: "Q", Suit: SuitHearts}, &up, PlayingCardOptions{}) }), 300, 300)
	settles(t, tt)
	click(t, tt, "Face-down card")
	tt.Advance(100 * time.Millisecond)
	click(t, tt, "Q of hearts")
	for range 6 {
		if exact(tt, "Q") {
			t.Fatal("a flip turned back before edge-on showed the face")
		}
		tt.Advance(60 * time.Millisecond)
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	settles(t, tt)
	click(t, tt, "Face-down card")
	if !exact(tt, "Q") {
		t.Error("with Reduce Motion the face waits")
	}
}

func TestHandReservesFan(t *testing.T) {
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	cards := make([]Card, len(ranks))
	for i, r := range ranks {
		cards[i] = Card{Rank: r, Suit: SuitClubs}
	}
	picked := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Hand(c, cards, &picked, HandOptions{})
		ui.Text(c, "After the hand")
	}), 900, 600)
	settles(t, tt)
	after := rectOf(t, tt, "After the hand")
	for _, r := range ranks {
		if b := rectOf(t, tt, r+" of clubs"); b.Y+b.H > after.Y+0.5 {
			t.Errorf("%s of clubs reaches %.1f, past the next line at %.1f", r, b.Y+b.H, after.Y)
		}
	}
}

func TestSkillTreeWideRow(t *testing.T) {
	skills := make([]Skill, 7)
	for i := range skills {
		skills[i] = Skill{ID: fmt.Sprint(i), Name: fmt.Sprintf("Skill %d", i+1), Icon: "star", Cost: 1}
	}
	var learned []string
	points := 3
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(300).Children(func() { SkillTree(c, skills, &learned, &points, SkillTreeOptions{}) })
	}), 320, 400)
	settles(t, tt)
	a, b := rectOf(t, tt, "Learn Skill 1"), rectOf(t, tt, "Learn Skill 2")
	if a.X < 10 || a.X+a.W > b.X {
		t.Errorf("first two skills overlap or leave the card: %v %v", a, b)
	}
}

func TestBigNumbersStayInside(t *testing.T) {
	teams := [2]Team{{"Kiln", math.MaxInt}, {"Glaze", math.MaxInt - 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(300).Children(func() { Scoreboard(c, &teams, ScoreboardOptions{}) })
	}), 320, 400)
	settles(t, tt)
	l, r := rectOf(t, tt, fmt.Sprint(math.MaxInt)), rectOf(t, tt, fmt.Sprint(math.MaxInt-1))
	if l.X < 10 || l.X+l.W > r.X || r.X+r.W > 310 {
		t.Errorf("scores %v %v overlap or leave the card", l, r)
	}
	tt = ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(300).Children(func() {
			Leaderboard(c, []Player{{Name: "Maya Chen", Score: 42, Prev: math.MaxInt}}, LeaderboardOptions{})
		})
	}), 320, 400)
	settles(t, tt)
	if m, s := rectOf(t, tt, fmt.Sprint(math.MaxInt-1)), rectOf(t, tt, "42"); m.X+m.W > s.X {
		t.Errorf("movement %v runs into the score %v", m, s)
	}
	tt = ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(300).Children(func() { XPBar(c, math.MaxInt, 1, 2, XPBarOptions{}) })
	}), 320, 400)
	settles(t, tt)
	n, xp := rectOf(t, tt, fmt.Sprint(math.MaxInt)), rectOf(t, tt, "1 / 2 XP")
	if n.W > 44 || xp.W < 20 || xp.X+xp.W > 310 {
		t.Errorf("level %v leaves its badge or XP %v has no room", n, xp)
	}
}

func TestHeartsBreakInside(t *testing.T) {
	lives := 3
	tt := ui.NewTester(padded(func(c *ui.Context) { Hearts(c, lives, 3, HeartsOptions{}) }), 300, 200)
	settles(t, tt)
	lives = 2
	tt.Frame()
	row := rectOf(t, tt, "2 of 3 lives")
	for _, at := range []time.Duration{100, 100, 100, 100} {
		tt.Advance(at * time.Millisecond)
		img := tt.Image()
		for y := int(row.Y+row.H) + 1; y < int(row.Y+row.H)+30; y++ {
			for x := int(row.X); x < int(row.X+row.W); x++ {
				if p := img.RGBAAt(x, y); int(p.R)-int(p.G) > 8 {
					t.Fatalf("the break paints below the row at %d,%d", x, y)
				}
			}
		}
	}
}

func TestJoystickNaN(t *testing.T) {
	panics(t, "needs x and y", func(c *ui.Context) {
		x, y := float32(math.NaN()), float32(0)
		Joystick(c, &x, &y, JoystickOptions{})
	})
}
