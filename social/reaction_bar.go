package social

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Reaction is one reaction and count.
type Reaction struct {
	ID, Name, Icon string
	Count          int
	Active         bool
}

// ReactionBar toggles reactions and counts; cue Pop.
func ReactionBar(c *ui.Context, reactions *[]Reaction) ui.Element {
	if reactions == nil || len(*reactions) == 0 {
		panic("social: ReactionBar needs reactions")
	}
	ids := make([]string, len(*reactions))
	for i, reaction := range *reactions {
		if reaction.Name == "" || reaction.Icon == "" || reaction.Count < 0 || reaction.Active && reaction.Count == 0 {
			panic(fmt.Sprintf("social: ReactionBar reaction %q is invalid", reaction.ID))
		}
		ids[i] = reaction.ID
	}
	unique("ReactionBar", ids)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, reaction := range *reactions {
			c.Key(reaction.ID)
			on := reaction.Active
			b := kit.Press(ui.ButtonBase(c).Label(reaction.Name).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).Shrink(0).
				Padding(t.Space(1), t.Space(2.5)).Radius(999).Checked(on).Disabled(!on && reaction.Count == math.MaxInt))
			hover := unit(b.Hovered())
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, hover).Mix(t.Accent.Alpha(0.14), unit(on)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Alpha(0).Mix(t.Accent, unit(on)), ui.Snappy))
			b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.Accent, unit(on)), ui.Snappy))
			kit.Stagger(c, root, b, i).Children(func() {
				c.Key(on)
				mark := glyph(c, reaction.Icon, 16)
				if on {
					kit.Arrive(mark, ui.Pose{Scale: 1.6, Rotate: -15}, 0)
				}
				kit.Numeric(c, number(c, reaction.Count)).FontSize(t.Rem(0.8)).Bold().Shrink(0)
			})
			if b.Clicked() && (on || reaction.Count < math.MaxInt) {
				(*reactions)[i].Active = !on
				if on {
					(*reactions)[i].Count--
				} else {
					(*reactions)[i].Count++
				}
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("ReactionBar", "toggled", "index", i, "count", (*reactions)[i].Count)
			}
		}
	})
	return root
}
