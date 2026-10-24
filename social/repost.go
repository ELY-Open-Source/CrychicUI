package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RepostOptions are the repost menu labels.
type RepostOptions struct{ Repost, Quote string }

// Repost chooses repost or quote from a menu; cue Pop.
func Repost(c *ui.Context, action *string, o RepostOptions) ui.Element {
	if action == nil {
		panic("social: Repost needs action state")
	}
	def(&o.Repost, "Repost")
	def(&o.Quote, "Quote")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(1)).Role(ui.RoleGroup)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		b := button(c, o.Repost).Checked(*open)
		b.Children(func() { ui.Text(c, o.Repost) })
		if b.Clicked() {
			*open = !*open
			sound.Play(sound.Pop)
			kit.Log("Repost", "menu", "open", *open)
		}
		if *open {
			menu := ui.Column(c).Gap(t.Space(1)).Padding(t.Space(1)).Radius(t.Radius).Border(1, t.Border).Transition(kit.Enter)
			menu.Children(func() {
				for _, item := range []struct {
					id, label string
				}{{"repost", o.Repost}, {"quote", o.Quote}} {
					c.Key(item.id)
					choice := button(c, item.label).FillWidth()
					choice.Children(func() { ui.Text(c, item.label) })
					if choice.Clicked() {
						*action = item.id
						*open = false
						root.SetChanged()
						sound.Play(sound.Pop)
						kit.Log("Repost", "selected", "id", item.id)
					}
				}
			})
		}
	})
	return root
}
