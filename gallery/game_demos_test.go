package gallery

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestGameDemosLog(t *testing.T) {
	press := func(names ...string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, name := range names {
				if err := tt.Click(name); err != nil {
					return err
				}
				tt.Frame()
			}
			return nil
		}
	}
	drag := func(from, to string, dx float32) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			a, ok := tt.Find(from)
			if !ok {
				return fmt.Errorf("no %q", from)
			}
			b := a
			b.X += dx
			if to != "" {
				if b, ok = tt.Find(to); !ok {
					return fmt.Errorf("no %q", to)
				}
			}
			tt.Press(a.X+a.W/2, a.Y+a.H/2)
			tt.Move(b.X+b.W/2, b.Y+b.H/2)
			tt.Release(b.X+b.W/2, b.Y+b.H/2)
			tt.Frame()
			return nil
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Dice":        {press("Roll"), "rolled"},
		"Coin":        {press("Flip"), "flipped"},
		"SpinWheel":   {press("Spin"), "spun to"},
		"Scoreboard":  {press("Add a point to Kiln Comets"), "score 4-2"},
		"HealthBar":   {press("Take 30 damage"), "hp 66"},
		"XPBar":       {press("Finish quest (+600 XP)"), "level 11 xp 2750"},
		"Inventory":   {drag("Ember sword", "Empty slot 4", 0), "inventory rearranged"},
		"Quests":      {press("Track Lost shipment"), "tracking Lost shipment"},
		"Hand":        {press("K of spades"), "picked 3"},
		"PlayingCard": {press("A of spades"), "ace up=false"},
		"TicTacToe":   {press("Cell 2"), "marks 5"},
		"Memory":      {press("Card 1", "Card 7"), "turn 1, pairs 0\nturn 2, pairs 1"},
		"SlotMachine": {press("Spin"), "spun"},
		"Hearts":      {press("Lose a life"), "lives 2"},
		"Combo":       {press("Hit"), "combo 8"},
		"LevelPick":   {press("Level 1-2"), "level 1-2"},
		"StatPoints":  {press("Raise Strength"), "points left 4"},
		"SkillTree":   {press("Learn Fireball"), "learned 2 points 2"},
		"Loot":        {press("Open chest"), "chest opened"},
		"DailyReward": {press("Claim"), "claimed 4"},
		"ReadyCheck":  {press("Ready"), "ready true"},
		"ChessClock":  {press("End White's move"), "turn 1"},
		"SlidePuzzle": {press("Tile 7"), "tiles [1 2 3 4 5 6 7 0 8]"},
		"Joystick":    {drag("Joystick", "", 80), "stick"},
	}
	quiet := map[string]bool{"Leaderboard": true, "ItemTip": true, "Achievements": true, "TurnTimer": true, "Character": true, "Lobby": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "game" {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(2 * time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 30 {
		t.Errorf("checked %d game demos", seen)
	}
}
