package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/game"
	"github.com/egoist/mygo/ui"
)

var gameEntries = []Entry{
	{ID: 1301, Name: "Dice", Category: "game", Summary: "Die rolls with tumble; cue Drop.", Demo: demoGameDice},
	{ID: 1302, Name: "Coin", Category: "game", Summary: "Coin flips; cue Pop.", Demo: demoGameCoin},
	{ID: 1303, Name: "SpinWheel", Category: "game", Summary: "Wheel spins to a slice; cue Tick per slice.", Demo: demoGameSpinWheel},
	{ID: 1304, Name: "Scoreboard", Category: "game", Summary: "Two scores that roll; cue Pop.", Demo: demoGameScoreboard},
	{ID: 1305, Name: "Leaderboard", Category: "game", Summary: "Ranks with movement arrows.", Demo: demoGameLeaderboard},
	{ID: 1306, Name: "HealthBar", Category: "game", Summary: "HP bar with damage trail.", Demo: demoGameHealthBar},
	{ID: 1307, Name: "XPBar", Category: "game", Summary: "XP bar with level up; cue Success.", Demo: demoGameXPBar},
	{ID: 1308, Name: "Inventory", Category: "game", Summary: "Item slots; drag moves; cue Drop.", Demo: demoGameInventory},
	{ID: 1309, Name: "ItemTip", Category: "game", Summary: "Item card with stats.", Demo: demoGameItemTip},
	{ID: 1310, Name: "Achievements", Category: "game", Summary: "Locked and unlocked list.", Demo: demoGameAchievements},
	{ID: 1311, Name: "Quests", Category: "game", Summary: "Quest log with progress.", Demo: demoGameQuests},
	{ID: 1312, Name: "TurnTimer", Category: "game", Summary: "Ring timer for a turn; cue Warning.", Demo: demoGameTurnTimer},
	{ID: 1313, Name: "Hand", Category: "game", Summary: "Fanned cards; hover lifts; cue Swipe.", Demo: demoGameHand},
	{ID: 1314, Name: "PlayingCard", Category: "game", Summary: "Card that flips; cue Swipe.", Demo: demoGamePlayingCard},
	{ID: 1315, Name: "TicTacToe", Category: "game", Summary: "Board with win line; cue Pop, Success.", Demo: demoGameTicTacToe},
}

func demoGameDice(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	faces := ui.Local(root, "faces", func() []int { return []int{3, 5} })
	root.Children(func() {
		if game.Dice(c, faces, game.DiceOptions{}).Changed() {
			log(fmt.Sprintf("rolled %v", *faces))
		}
	})
}

func demoGameCoin(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340)
	heads := boolState(root, "heads", true)
	root.Children(func() {
		if game.Coin(c, heads, game.CoinOptions{}).Changed() {
			log(fmt.Sprintf("flipped heads=%v", *heads))
		}
	})
}

func demoGameSpinWheel(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	picked := intState(root, "picked", -1)
	prizes := []string{"50 coins", "Gem", "Spin again", "100 coins", "Shield", "Mystery box"}
	root.Children(func() {
		if game.SpinWheel(c, prizes, picked, game.SpinWheelOptions{}).Changed() {
			log("spun to " + prizes[*picked])
		}
	})
}

func demoGameScoreboard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	teams := ui.Local(root, "teams", func() [2]game.Team {
		return [2]game.Team{{Name: "Kiln Comets", Score: 3}, {Name: "Glaze Giants", Score: 2}}
	})
	root.Children(func() {
		if game.Scoreboard(c, teams, game.ScoreboardOptions{Period: "Q3"}).Changed() {
			log(fmt.Sprintf("score %d-%d", teams[0].Score, teams[1].Score))
		}
	})
}

func demoGameLeaderboard(c *ui.Context, log func(string)) {
	players := []game.Player{
		{Name: "Priya Nair", Score: 18_420, Prev: 2},
		{Name: "Maya Chen", Score: 17_905, Prev: 1},
		{Name: "Lucas Martin", Score: 15_310, Prev: 3},
		{Name: "Sofia Rossi", Score: 12_880},
		{Name: "Daniel Okafor", Score: 11_045, Prev: 4},
	}
	ui.Column(c).Width(480).Children(func() { game.Leaderboard(c, players, game.LeaderboardOptions{Title: "Season 4"}) })
}

func demoGameHealthBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380).Gap(12)
	hp := intState(root, "hp", 96)
	root.Children(func() {
		game.HealthBar(c, *hp, 120, game.HealthBarOptions{})
		ui.Row(c).Gap(8).Children(func() {
			if ui.Button(c, "Take 30 damage").Clicked() {
				*hp = max(0, *hp-30)
				log(fmt.Sprintf("hp %d", *hp))
			}
			if ui.Button(c, "Drink potion").Clicked() {
				*hp = min(120, *hp+40)
				log(fmt.Sprintf("hp %d", *hp))
			}
		})
	})
}

func demoGameXPBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440).Gap(12)
	level := intState(root, "level", 11)
	xp := intState(root, "xp", 2_150)
	root.Children(func() {
		game.XPBar(c, *level, int64(*xp), 3_000, game.XPBarOptions{})
		if ui.Button(c, "Finish quest (+600 XP)").Clicked() {
			*xp += 600
			if *xp >= 3_000 {
				*xp -= 3_000
				*level++
			}
			log(fmt.Sprintf("level %d xp %d", *level, *xp))
		}
	})
}

func demoGameInventory(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	slots := ui.Local(root, "slots", func() []game.Item {
		s := make([]game.Item, 10)
		s[0] = game.Item{Name: "Ember sword", Icon: "sword", Count: 1, Rarity: game.Epic}
		s[1] = game.Item{Name: "Oak shield", Icon: "shield", Count: 1, Rarity: game.Uncommon}
		s[2] = game.Item{Name: "Healing potion", Icon: "flask-conical", Count: 5, Rarity: game.Common}
		s[4] = game.Item{Name: "Moonstone", Icon: "gem", Count: 2, Rarity: game.Rare}
		s[6] = game.Item{Name: "Crown of ash", Icon: "crown", Count: 1, Rarity: game.Legendary}
		s[7] = game.Item{Name: "Old map", Icon: "map", Count: 1, Rarity: game.Common}
		return s
	})
	root.Children(func() {
		if game.Inventory(c, slots, game.InventoryOptions{}).Changed() {
			log("inventory rearranged")
		}
	})
}

func demoGameItemTip(c *ui.Context, log func(string)) {
	item := game.ItemDetail{Name: "Ember sword", Kind: "One-handed sword", Icon: "sword", Rarity: game.Epic, Level: 12,
		Text: "Forged in the kiln at Hollow Peak. Warm to the touch.",
		Stats: []game.Stat{{Name: "Damage", Value: 48, Delta: 12}, {Name: "Speed", Value: 1, Delta: 0},
			{Name: "Fire", Value: 15, Delta: 15}, {Name: "Durability", Value: 80, Delta: -20}}}
	ui.Column(c).Width(340).Children(func() { game.ItemTip(c, item, game.ItemTipOptions{}) })
}

func demoGameAchievements(c *ui.Context, log func(string)) {
	list := []game.Achievement{
		{Name: "First steps", Text: "Finish the tutorial.", Icon: "footprints", Unlocked: true, When: "Oct 2"},
		{Name: "Dragon slayer", Text: "Defeat the ember drake.", Icon: "flame", Unlocked: true, When: "Oct 19"},
		{Name: "Collector", Text: "Find 50 gems.", Icon: "gem", Goal: 50, Progress: 34},
		{Name: "Untouchable", Text: "Clear a dungeon without damage.", Icon: "shield", Goal: 1},
		{Name: "Cartographer", Text: "Discover every region.", Icon: "map", Goal: 12, Progress: 9},
	}
	ui.Column(c).Width(460).Children(func() { game.Achievements(c, list, game.AchievementsOptions{}) })
}

func demoGameQuests(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	tracked := textState(root, "tracked", "The ember drake")
	quests := []game.Quest{
		{Title: "The ember drake", Reward: "800 XP · Ember sword",
			Steps: []game.QuestStep{{Text: "Speak to the smith", Done: true}, {Text: "Find the cave entrance", Done: true}, {Text: "Defeat the drake"}}},
		{Title: "Lost shipment", Reward: "250 coins",
			Steps: []game.QuestStep{{Text: "Search the old road"}, {Text: "Return the crates"}}},
		{Title: "A gift for Maya", Reward: "Friendship +1",
			Steps: []game.QuestStep{{Text: "Gather moonflowers", Done: true}, {Text: "Deliver the bouquet", Done: true}}},
	}
	root.Children(func() {
		if game.Quests(c, quests, tracked, game.QuestsOptions{}).Changed() {
			log("tracking " + *tracked)
		}
	})
}

func demoGameTurnTimer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	born := ui.Local(root, "born", func() time.Time { return c.Now() })
	root.Children(func() {
		game.TurnTimer(c, "Maya", born.Add(45*time.Second), time.Minute, game.TurnTimerOptions{})
	})
}

func demoGameHand(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).AlignItems(ui.Center)
	picked := intState(root, "picked", -1)
	cards := []game.Card{{Rank: "9", Suit: game.SuitHearts}, {Rank: "J", Suit: game.SuitClubs}, {Rank: "Q", Suit: game.SuitDiamonds},
		{Rank: "K", Suit: game.SuitSpades}, {Rank: "A", Suit: game.SuitHearts}}
	root.Children(func() {
		if game.Hand(c, cards, picked, game.HandOptions{}).Changed() {
			log(fmt.Sprintf("picked %d", *picked))
		}
	})
}

func demoGamePlayingCard(c *ui.Context, log func(string)) {
	root := ui.Row(c).Width(340).Gap(16).Justify(ui.Center)
	a, b := boolState(root, "a", true), boolState(root, "b", false)
	root.Children(func() {
		if game.PlayingCard(c, game.Card{Rank: "A", Suit: game.SuitSpades}, a, game.PlayingCardOptions{}).Changed() {
			log(fmt.Sprintf("ace up=%v", *a))
		}
		if game.PlayingCard(c, game.Card{Rank: "Q", Suit: game.SuitHearts}, b, game.PlayingCardOptions{}).Changed() {
			log(fmt.Sprintf("queen up=%v", *b))
		}
	})
}

func demoGameTicTacToe(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	board := ui.Local(root, "board", func() [9]game.Mark {
		return [9]game.Mark{game.MarkX, game.Blank, game.MarkO, game.Blank, game.MarkX, game.Blank, game.Blank, game.MarkO, game.Blank}
	})
	root.Children(func() {
		if game.TicTacToe(c, board, game.TicTacToeOptions{}).Changed() {
			n := 0
			for _, m := range board {
				if m != game.Blank {
					n++
				}
			}
			log(fmt.Sprintf("marks %d", n))
		}
	})
}
