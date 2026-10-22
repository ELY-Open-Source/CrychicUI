package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Card is a title, text and Lucide icon.
type Card struct{ Title, Text, Icon string }

// CardList stacks cards that lift on hover.
func CardList(c *ui.Context, label string, cards []Card, chosen *int) ui.Element {
	if label == "" || len(cards) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(cards) {
		panic(fmt.Sprintf("data: CardList %q needs cards and a choice or -1", label))
	}
	for _, card := range cards {
		if card.Title == "" {
			panic(fmt.Sprintf("data: CardList %q cards need titles", label))
		}
		icon.Must(card.Icon)
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleList).Label(label)
	shown := *chosen
	root.Children(func() {
		for i, card := range cards {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(card.Title).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(4)).
				Radius(t.Radius).Background(t.Background).Checked(i == shown))
			lift := b.AnimateTo("lift", unit(b.Hovered()), ui.Snappy)
			b.Translate(0, -t.Space(1)*lift).Shadow(0, 1, 3, 0, kit.Shade(0.06)).
				Border(1, b.AnimateColor("edge", t.Border.Mix(t.Accent, unit(i == shown)), ui.Snappy))
			kit.Arrive(b, ui.Pose{Y: t.Space(3), Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(t.Accent.Alpha(0.12)).Children(func() {
					ui.Icon(c, icon.Must(card.Icon)).Size(16, 16).TextColor(t.Accent)
				})
				ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, card.Title).Bold()
					if card.Text != "" {
						kit.Muted(c, card.Text)
					}
				})
			})
			if b.Clicked() {
				*chosen = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("CardList", "chose", "index", i)
			}
		}
	})
	return root
}
