package social

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
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
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, reaction := range *reactions {
			if reaction.ID == "" || reaction.Name == "" || reaction.Icon == "" || reaction.Count < 0 || reaction.Active && reaction.Count == 0 {
				panic("social: invalid reaction")
			}
			c.Key(reaction.ID)
			b := button(c, reaction.Name).Row().Gap(t.Space(1)).AlignItems(ui.Center).Radius(999).
				Checked(reaction.Active).Disabled(!reaction.Active && reaction.Count == math.MaxInt)
			b.Background(b.AnimateColor("reaction", t.Surface.Mix(t.Accent.Alpha(0.15), unit(reaction.Active)), ui.Snappy))
			b.Children(func() {
				ui.Icon(c, icon.Must(reaction.Icon)).Size(16, 16).TextColor(t.Accent)
				kit.Numeric(c, fmt.Sprint(reaction.Count))
			})
			if b.Clicked() && (reaction.Active || reaction.Count < math.MaxInt) {
				(*reactions)[i].Active = !reaction.Active
				if reaction.Active {
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
