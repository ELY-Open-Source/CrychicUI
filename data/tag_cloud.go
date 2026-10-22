package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Tag is a word and its weight, zero or more.
type Tag struct {
	Text   string
	Weight float64
}

// TagCloud sizes words by weight.
func TagCloud(c *ui.Context, label string, tags []Tag, chosen *int) ui.Element {
	if label == "" || len(tags) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(tags) {
		panic(fmt.Sprintf("data: TagCloud %q needs tags and a choice or -1", label))
	}
	var words []string
	lo, hi := tags[0].Weight, tags[0].Weight
	for _, tg := range tags {
		if !finite64(tg.Weight) || tg.Weight < 0 {
			panic(fmt.Sprintf("data: TagCloud %q weights must be finite and not negative", label))
		}
		words = append(words, tg.Text)
		lo, hi = min(lo, tg.Weight), max(hi, tg.Weight)
	}
	unique("TagCloud", label, words)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleList).Label(label).FocusGroup(ui.Horizontal)
	shown := *chosen
	return root.Children(func() {
		for i, tg := range tags {
			c.Key(tg.Text)
			k := float32(0.5)
			if hi > lo {
				k = float32((tg.Weight - lo) / (hi - lo))
			}
			on := i == shown
			b := kit.Press(ui.ButtonBase(c).Label(tg.Text).Padding(t.Space(0.5), t.Space(1.5)).Radius(t.Radius).Checked(on))
			b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.Accent, unit(on || b.Hovered())), ui.Snappy)).
				Background(b.AnimateColor("bg", t.Accent.Alpha(0.12*unit(on)), ui.Snappy))
			kit.Arrive(b, ui.Pose{Scale: 0.8, Fade: 1}, time.Duration(i)*20*time.Millisecond).Children(func() {
				ui.Text(c, tg.Text).FontSize(t.Rem(0.8 + k)).Bold()
			})
			if b.Clicked() {
				*chosen = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("TagCloud", "chose", "index", i)
			}
		}
	})
}
