package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HashtagOptions format the grouped post count.
type HashtagOptions struct{ Posts string }

// Hashtag opens a tag chip with count; cue Tap.
func Hashtag(c *ui.Context, tag string, count int, opened *bool, o HashtagOptions) ui.Element {
	if tag == "" || count < 0 || opened == nil {
		panic(fmt.Sprintf("social: Hashtag %q needs a count and state", tag))
	}
	def(&o.Posts, "%s posts")
	t := c.Theme()
	root := ui.Row(c).FillWidth()
	room := root.LayoutBounds().W
	kit.Settle(c, root, []float32{room}, func() []float32 { return []float32{root.LayoutBounds().W} })
	// The count takes two fifths at most.
	counted := room * 0.4
	named := max(0, room-counted-t.Space(14))
	root.Children(func() {
		chip := kit.Press(ui.ButtonBase(c).Label(tag).Row().Gap(t.Space(2)).AlignItems(ui.Center).
			Padding(t.Space(1), t.Space(3), t.Space(1), t.Space(1)).Radius(999).Checked(*opened))
		hover := unit(chip.Hovered())
		chip.Background(chip.AnimateColor("bg", t.Accent.Alpha(0.08).Mix(t.Accent.Alpha(0.16), max(hover, unit(*opened))), ui.Snappy))
		kit.Arrive(chip, ui.Pose{Scale: 0.9, Fade: 1}, 0).Children(func() {
			ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText).Children(func() {
				glyph(c, "hash", 14)
			})
			ui.Text(c, tag).FontSize(t.Rem(0.85)).Bold().TextColor(t.Accent).Shrink(0).MaxWidth(named).SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Posts, number(c, count))).FontSize(t.Rem(0.78)).TextColor(t.TextMuted).Shrink(0).
				MaxWidth(counted).SingleLine()
		})
		if chip.Clicked() {
			*opened = true
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("Hashtag", "opened", "count", count)
		}
	})
	return root
}
