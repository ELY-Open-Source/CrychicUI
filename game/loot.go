package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LootOptions word the chest; Count takes an item count.
type LootOptions struct{ Open, Count string }

// Loot is a chest that opens to items; cue Pop.
func Loot(c *ui.Context, items []Item, opened *bool, o LootOptions) ui.Element {
	def(&o.Open, "Open chest")
	if opened == nil || len(items) == 0 || len(items) > 8 {
		panic(fmt.Sprintf("game: Loot %q needs a state and 1 to 8 items", o.Open))
	}
	names := make([]string, len(items))
	for i, it := range items {
		validRarity("Loot", it.Name, it.Rarity)
		if it.Icon == "" || it.Count < 1 {
			panic(fmt.Sprintf("game: Loot %q needs an icon and a count of at least 1", it.Name))
		}
		names[i] = it.Name
	}
	unique("Loot", names)
	def(&o.Count, "× %d")
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		best := Common
		for _, it := range items {
			best = max(best, it.Rarity)
		}
		ink := rarityInk(t, best)
		c.Key(*opened)
		chest := ui.Box(c).Size(t.Space(20), t.Space(20)).Shrink(0).Center().Radius(t.Radius * 1.5).Background(ink.Alpha(0.12)).TextColor(ink)
		if *opened {
			kit.Arrive(chest, ui.Pose{Scale: 1.25, Rotate: -8, Fade: 0.3}, 0)
		}
		chest.Children(func() {
			mark := "package"
			if *opened {
				mark = "package-open"
			}
			glyph(c, mark, 40)
		})
		if !*opened {
			b := kit.Press(ui.ButtonBase(c).Label(o.Open).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0))
			b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "key-round", 16)
					ui.Text(c, o.Open).FontSize(t.Rem(0.85)).Bold().Shrink(0)
				})
			})
			if b.Clicked() {
				*opened = true
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("Loot", "opened", "items", len(items))
			}
			return
		}
		ui.Row(c).Wrap().Gap(t.Space(2)).Justify(ui.Center).Children(func() {
			for i, it := range items {
				tone := rarityInk(t, it.Rarity)
				c.Key(it.Name)
				cell := ui.Column(c).Width(t.Space(20)).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Border(1, tone.Alpha(0.6)).Background(tone.Alpha(0.08))
				kit.Arrive(cell, ui.Pose{Y: -t.Space(10), Scale: 0.4, Fade: 1}, 150*time.Millisecond+time.Duration(i)*110*time.Millisecond).
					Children(func() {
						glyph(c, it.Icon, 22).TextColor(tone)
						ui.Text(c, it.Name).FontSize(t.Rem(0.68)).Bold().FillWidth().TextAlign(ui.Center).SingleLine()
						if it.Count > 1 {
							kit.Muted(c, fmt.Sprintf(o.Count, it.Count)).Font("monospace").FontSize(t.Rem(0.62)).Shrink(0)
						}
					})
			}
		})
	})
	return root
}
