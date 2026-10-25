package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RepostOptions word the trigger and menu; hints are optional.
type RepostOptions struct{ Repost, Quote, RepostHint, QuoteHint string }

// Repost chooses repost or quote from a menu; cue Pop.
func Repost(c *ui.Context, action *string, o RepostOptions) ui.Element {
	if action == nil {
		panic("social: Repost needs action state")
	}
	def(&o.Repost, "Repost")
	def(&o.Quote, "Quote")
	def(&o.RepostHint, "Share it with your followers")
	def(&o.QuoteHint, "Add your thoughts")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(1.5)).AlignItems(ui.Start).Role(ui.RoleGroup)
	open := ui.Local(root, "open", func() bool { return false })
	done := *action != ""
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(o.Repost).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(3)).Radius(999).Checked(*open))
		hover := unit(b.Hovered() || *open)
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.Success.Alpha(0.12), max(hover, unit(done))), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Success, max(hover, unit(done))), ui.Snappy))
		b.Children(func() {
			c.Key(*action)
			kit.Arrive(glyph(c, "repeat-2", 16), ui.Pose{Rotate: -180, Scale: 0.6}, 0)
			ui.Text(c, o.Repost).FontSize(t.Rem(0.82)).Bold().Shrink(0)
		})
		if b.Clicked() {
			*open = !*open
			sound.Play(sound.Pop)
			kit.Log("Repost", "menu", "open", *open)
		}
		c.Key("menu")
		ui.Column(c).Children(func() {
			if !*open {
				return
			}
			menu := ui.Column(c).Width(t.Space(64)).Padding(t.Space(1)).Radius(t.Radius).Border(1, t.Border).Background(t.Background).
				Shadow(0, 8, 24, 0, t.Text.Alpha(0.14)).Transition(kit.Enter)
			menu.Children(func() {
				for _, item := range []struct{ id, label, hint, mark string }{
					{"repost", o.Repost, o.RepostHint, "repeat-2"}, {"quote", o.Quote, o.QuoteHint, "pen-line"},
				} {
					c.Key(item.id)
					on := *action == item.id
					choice := kit.Press(ui.ButtonBase(c).Label(item.label).Row().FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).
						Padding(t.Space(2)).Radius(t.Radius).Checked(on))
					choice.Background(choice.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(choice.Hovered())), ui.Snappy))
					choice.Children(func() {
						ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Background(t.Success.Alpha(0.12)).
							TextColor(t.Success).Children(func() { glyph(c, item.mark, 16) })
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							ui.Text(c, item.label).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
							kit.Muted(c, item.hint).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
						})
						c.Key(on)
						if on {
							kit.Arrive(glyph(c, "check", 14).TextColor(t.Success), ui.Pose{Scale: 0.4, Fade: 1}, 0)
						}
					})
					if choice.Clicked() {
						*action = item.id
						*open = false
						root.SetChanged()
						sound.Play(sound.Pop)
						kit.Log("Repost", "selected", "id", item.id)
					}
				}
			})
		})
	})
	return root
}
