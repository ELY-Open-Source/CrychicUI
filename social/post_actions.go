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
	root := ui.Row(c).FillWidth().Gap(t.Space(1)).Justify(ui.SpaceBetween).Role(ui.RoleGroup)
	root.Children(func() {
		actions := []struct {
			id, label, mark string
			on              bool
			ink             ui.Color
		}{
			{"like", o.Like, "heart", state.Liked, t.Danger},
			{"comment", o.Comment, "message-circle", false, t.Accent},
			{"share", o.Share, "repeat-2", false, t.Success},
			{"save", o.Save, "bookmark", state.Saved, t.Accent},
		}
		for _, action := range actions {
			c.Key(action.id)
			b := kit.Press(ui.ButtonBase(c).Label(action.label).Padding(t.Space(1.5), t.Space(2.5)).Radius(999).Shrink(0).Checked(action.on))
			hover := unit(b.Hovered())
			b.Background(b.AnimateColor("bg", t.Background.Mix(action.ink.Alpha(0.1), max(hover, unit(action.on))), ui.Snappy))
			b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(action.ink, max(hover*0.7, unit(action.on))), ui.Snappy))
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					c.Key(action.on)
					mark := glyph(c, action.mark, 16)
					if action.on {
						kit.Arrive(mark, ui.Pose{Scale: 1.5, Fade: 0.5}, 0)
					}
					ui.Text(c, action.label).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				})
			})
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
