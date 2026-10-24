package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PostActionsOptions are the post action labels.
type PostActionsOptions struct{ Like, Comment, Share, Save string }

// PostActions toggles like and save, comment and share; cue Pop.
func PostActions(c *ui.Context, state *PostState, o PostActionsOptions) ui.Element {
	if state == nil {
		panic("social: PostActions needs state")
	}
	def(&o.Like, "Like")
	def(&o.Comment, "Comment")
	def(&o.Share, "Share")
	def(&o.Save, "Save")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		actions := []struct {
			id, label string
			checked   bool
		}{{"like", o.Like, state.Liked}, {"comment", o.Comment, false}, {"share", o.Share, false}, {"save", o.Save, state.Saved}}
		for _, action := range actions {
			c.Key(action.id)
			b := button(c, action.label).Checked(action.checked)
			b.Children(func() { ui.Text(c, action.label) })
			if b.Clicked() {
				switch action.id {
				case "like":
					state.Liked = !state.Liked
				case "save":
					state.Saved = !state.Saved
				default:
					state.Action = action.id
				}
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("PostActions", "action", "id", action.id)
			}
		}
	})
	return root
}
