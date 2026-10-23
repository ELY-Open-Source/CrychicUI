package chat

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Reaction is an emoji count; Mine when I gave it.
type Reaction struct {
	Emoji string
	Count int
	Mine  bool
}

// Reactions shows emoji counts; a tap toggles mine.
func Reactions(c *ui.Context, label string, items *[]Reaction) ui.Element {
	if label == "" || items == nil || len(*items) == 0 {
		panic(fmt.Sprintf("chat: Reactions %q needs a label and items", label))
	}
	for _, r := range *items {
		if r.Emoji == "" || r.Count < 0 || r.Mine && r.Count == 0 {
			panic(fmt.Sprintf("chat: Reactions %q items need an emoji and a count covering Mine", label))
		}
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).Wrap().Role(ui.RoleGroup).Label(label)
	return root.Children(func() {
		for i := range *items {
			c.Key(i)
			r := &(*items)[i]
			b := kit.Press(ui.ButtonBase(c).Label(r.Emoji).Row().Gap(t.Space(1)).AlignItems(ui.Center).Checked(r.Mine).
				Padding(t.Space(0.5), t.Space(2)).Radius(999))
			bg, bd := t.Surface.Mix(t.SurfaceHover, unit(b.Hovered())), t.Border
			if r.Mine {
				bg, bd = t.Accent.Alpha(0.12), t.Accent
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Border(1, b.AnimateColor("bd", bd, ui.Snappy))
			pops := ui.Local(b, "pops", func() int { return 0 })
			if b.Clicked() {
				r.Mine = !r.Mine
				if r.Mine {
					r.Count++
				} else {
					r.Count--
				}
				*pops++
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("Reactions", "toggled", "index", i, "mine", r.Mine)
			}
			b.Children(func() {
				e := ui.Text(c, r.Emoji)
				if *pops > 0 {
					k := 1 + 0.4*e.Keyframes(fmt.Sprint("pop.", *pops), []ui.Keyframe{{Value: 1}, {At: 220 * time.Millisecond, Value: 0, Ease: ui.EaseOut}})
					e.Scale(k, k)
				}
				kit.Numeric(c, strconv.Itoa(r.Count)).FontSize(t.Rem(0.85))
			})
		}
	})
}
