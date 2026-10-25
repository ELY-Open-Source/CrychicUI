package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/game"
	"github.com/egoist/mygo/ui"
)

var gameEntries2 = []Entry{
	{ID: 1316, Name: "Memory", Category: "game", Summary: "Card pairs to match; cue Swipe, Success.", Demo: demoGameMemory},
	{ID: 1317, Name: "SlotMachine", Category: "game", Summary: "Reels spin and stop; cue Tick.", Demo: demoGameSlotMachine},
	{ID: 1318, Name: "Hearts", Category: "game", Summary: "Lives as hearts; lost one breaks.", Demo: demoGameHearts},
	{ID: 1319, Name: "Combo", Category: "game", Summary: "Combo count that bounces.", Demo: demoGameCombo},
	{ID: 1320, Name: "LevelPick", Category: "game", Summary: "Level tiles with stars; cue Select.", Demo: demoGameLevelPick},
	{ID: 1321, Name: "Character", Category: "game", Summary: "Character card with stats.", Demo: demoGameCharacter},
	{ID: 1322, Name: "StatPoints", Category: "game", Summary: "Spend points on stats; cue Tick.", Demo: demoGameStatPoints},
	{ID: 1323, Name: "SkillTree", Category: "game", Summary: "Nodes unlock along links; cue Success.", Demo: demoGameSkillTree},
	{ID: 1324, Name: "Loot", Category: "game", Summary: "Chest opens to items; cue Pop.", Demo: demoGameLoot},
	{ID: 1325, Name: "DailyReward", Category: "game", Summary: "Days with claim; cue Success.", Demo: demoGameDailyReward},
	{ID: 1326, Name: "Lobby", Category: "game", Summary: "Players and ready marks.", Demo: demoGameLobby},
	{ID: 1327, Name: "ReadyCheck", Category: "game", Summary: "Ready and not ready; cue Select.", Demo: demoGameReadyCheck},
	{ID: 1328, Name: "ChessClock", Category: "game", Summary: "Two clocks switching; cue Tick.", Demo: demoGameChessClock},
	{ID: 1329, Name: "SlidePuzzle", Category: "game", Summary: "Tiles slide into the gap; cue Tick, Success.", Demo: demoGameSlidePuzzle},
	{ID: 1330, Name: "Joystick", Category: "game", Summary: "Thumb dragged in a circle; springs home.", Demo: demoGameJoystick},
}

func demoGameMemory(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	faces := []string{"cat", "moon", "anchor", "cherry", "moon", "rocket", "cat", "clover", "anchor", "rocket", "clover", "cherry"}
	matched := ui.Local(root, "matched", func() []bool { return make([]bool, len(faces)) })
	turns := ui.Local(root, "turns", func() int { return 0 })
	root.Children(func() {
		if game.Memory(c, faces, matched, game.MemoryOptions{Columns: 6, Turns: turns}).Changed() {
			n := 0
			for _, m := range *matched {
				if m {
					n++
				}
			}
			log(fmt.Sprintf("turn %d, pairs %d", *turns, n/2))
		}
	})
}

func demoGameSlotMachine(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	reels := ui.Local(root, "reels", func() [3]int { return [3]int{0, 3, 1} })
	symbols := []string{"cherry", "citrus", "grape", "bell", "gem", "clover"}
	root.Children(func() {
		if game.SlotMachine(c, symbols, reels, game.SlotMachineOptions{}).Changed() {
			log(fmt.Sprintf("spun %v", *reels))
		}
	})
}

func demoGameHearts(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340).Gap(12)
	lives := intState(root, "lives", 3)
	root.Children(func() {
		game.Hearts(c, *lives, 5, game.HeartsOptions{})
		ui.Row(c).Gap(8).Children(func() {
			if ui.Button(c, "Lose a life").Clicked() && *lives > 0 {
				*lives--
				log(fmt.Sprintf("lives %d", *lives))
			}
			if ui.Button(c, "Gain a life").Clicked() && *lives < 5 {
				*lives++
				log(fmt.Sprintf("lives %d", *lives))
			}
		})
	})
}

func demoGameCombo(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380).Gap(12)
	count := intState(root, "count", 7)
	root.Children(func() {
		game.Combo(c, *count, game.ComboOptions{})
		ui.Row(c).Gap(8).Children(func() {
			if ui.Button(c, "Hit").Clicked() {
				*count++
				log(fmt.Sprintf("combo %d", *count))
			}
			if ui.Button(c, "Miss").Clicked() {
				*count = 0
				log("combo 0")
			}
		})
	})
}

func demoGameLevelPick(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	picked := intState(root, "picked", 5)
	levels := []game.Level{{Name: "1-1", Stars: 3}, {Name: "1-2", Stars: 3}, {Name: "1-3", Stars: 2}, {Name: "1-4", Stars: 3},
		{Name: "1-5", Stars: 1}, {Name: "1-6"}, {Name: "1-7", Locked: true}, {Name: "1-8", Locked: true}}
	root.Children(func() {
		if game.LevelPick(c, levels, picked, game.LevelPickOptions{Title: "Ember Valley"}).Changed() {
			log("level " + levels[*picked].Name)
		}
	})
}

func demoGameCharacter(c *ui.Context, log func(string)) {
	hero := game.Hero{Name: "Priya Nair", Class: "Ember mage", Level: 14, HP: 84, MaxHP: 120, MP: 62, MaxMP: 90,
		Stats: []game.Stat{{Name: "Strength", Value: 9}, {Name: "Agility", Value: 14}, {Name: "Intellect", Value: 22},
			{Name: "Spirit", Value: 17}, {Name: "Stamina", Value: 12}, {Name: "Luck", Value: 6}}}
	ui.Column(c).Width(380).Children(func() { game.Character(c, hero, game.CharacterOptions{}) })
}

func demoGameStatPoints(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	points := intState(root, "points", 5)
	stats := ui.Local(root, "stats", func() []game.Stat {
		return []game.Stat{{Name: "Strength", Value: 9}, {Name: "Agility", Value: 14}, {Name: "Intellect", Value: 22}, {Name: "Stamina", Value: 12}}
	})
	root.Children(func() {
		if game.StatPoints(c, stats, points, game.StatPointsOptions{}).Changed() {
			log(fmt.Sprintf("points left %d", *points))
		}
	})
}

func demoGameSkillTree(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	points := intState(root, "points", 3)
	learned := ui.Local(root, "learned", func() []string { return []string{"spark"} })
	skills := []game.Skill{
		{ID: "spark", Name: "Spark", Icon: "zap", Cost: 1},
		{ID: "fireball", Name: "Fireball", Icon: "flame", Requires: []string{"spark"}, Cost: 1},
		{ID: "ward", Name: "Ward", Icon: "shield", Requires: []string{"spark"}, Cost: 1},
		{ID: "inferno", Name: "Inferno", Icon: "sun", Requires: []string{"fireball"}, Cost: 2},
		{ID: "aegis", Name: "Aegis", Icon: "shield-half", Requires: []string{"ward", "fireball"}, Cost: 2},
	}
	root.Children(func() {
		if game.SkillTree(c, skills, learned, points, game.SkillTreeOptions{}).Changed() {
			log(fmt.Sprintf("learned %d points %d", len(*learned), *points))
		}
	})
}

func demoGameLoot(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	opened := boolState(root, "opened", false)
	items := []game.Item{
		{Name: "Moonstone", Icon: "gem", Count: 3, Rarity: game.Rare},
		{Name: "Healing potion", Icon: "flask-conical", Count: 2, Rarity: game.Common},
		{Name: "Gold", Icon: "coins", Count: 250, Rarity: game.Uncommon},
		{Name: "Crown of ash", Icon: "crown", Count: 1, Rarity: game.Legendary},
	}
	root.Children(func() {
		if game.Loot(c, items, opened, game.LootOptions{}).Changed() {
			log("chest opened")
		}
		if *opened && ui.Button(c, "New chest").Clicked() {
			*opened = false
			log("chest reset")
		}
	})
}

func demoGameDailyReward(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	claimed := intState(root, "claimed", 3)
	rewards := []game.Reward{{Icon: "coins", Amount: "100"}, {Icon: "flask-conical", Amount: "Potion"}, {Icon: "coins", Amount: "250"},
		{Icon: "gem", Amount: "Gem"}, {Icon: "coins", Amount: "500"}, {Icon: "key-round", Amount: "Key"}, {Icon: "crown", Amount: "Chest"}}
	root.Children(func() {
		if game.DailyReward(c, rewards, claimed, game.DailyRewardOptions{}).Changed() {
			log(fmt.Sprintf("claimed %d", *claimed))
		}
	})
}

func demoGameLobby(c *ui.Context, log func(string)) {
	players := []game.Seat{{Name: "Maya Chen", Ready: true, Host: true}, {Name: "Daniel Okafor", Ready: true},
		{Name: "Priya Nair"}, {Name: "Lucas Martin", Ready: true}}
	ui.Column(c).Width(440).Children(func() { game.Lobby(c, players, 6, game.LobbyOptions{Title: "Hollow Peak raid"}) })
}

func demoGameReadyCheck(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	ready := boolState(root, "ready", false)
	others := []game.Seat{{Name: "Maya Chen", Ready: true}, {Name: "Daniel Okafor", Ready: true}, {Name: "Lucas Martin"}, {Name: "Sofia Rossi", Ready: true}}
	root.Children(func() {
		if game.ReadyCheck(c, others, ready, game.ReadyCheckOptions{}).Changed() {
			log(fmt.Sprintf("ready %v", *ready))
		}
	})
}

func demoGameChessClock(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	left := ui.Local(root, "left", func() [2]time.Duration { return [2]time.Duration{5 * time.Minute, 4*time.Minute + 38*time.Second} })
	turn := intState(root, "turn", -1)
	root.Children(func() {
		if game.ChessClock(c, left, turn, game.ChessClockOptions{}).Changed() {
			log(fmt.Sprintf("turn %d", *turn))
		}
	})
}

func demoGameSlidePuzzle(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340)
	tiles := ui.Local(root, "tiles", func() []int { return []int{1, 2, 3, 4, 5, 6, 0, 7, 8} })
	root.Children(func() {
		if game.SlidePuzzle(c, tiles, game.SlidePuzzleOptions{}).Changed() {
			log(fmt.Sprintf("tiles %v", *tiles))
		}
	})
}

func demoGameJoystick(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340).AlignItems(ui.Center).Gap(12)
	x, y := ui.Local(root, "x", func() float32 { return 0 }), ui.Local(root, "y", func() float32 { return 0 })
	root.Children(func() {
		if game.Joystick(c, x, y, game.JoystickOptions{}).Changed() {
			log(fmt.Sprintf("stick %.2f %.2f", *x, *y))
		}
		ui.Text(c, fmt.Sprintf("x %.2f · y %.2f", *x, *y)).Font("monospace").FontSize(13)
	})
}
