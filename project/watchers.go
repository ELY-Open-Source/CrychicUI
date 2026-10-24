package project

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WatchersOptions are the watcher labels.
type WatchersOptions struct{ Watch, Watching string }

// Watchers toggles watching and shows avatars; cue Toggle.
func Watchers(c *ui.Context, people []string, watching *bool, o WatchersOptions) ui.Element {
	if len(people) == 0 || watching == nil {
		panic("project: Watchers needs people and state")
	}
	def(&o.Watch, "Watch")
	def(&o.Watching, "Watching")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		name := o.Watch
		if *watching {
			name = o.Watching
		}
		b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1), t.Space(2)).Radius(999).Checked(*watching))
		b.Children(func() { ui.Text(c, name).SingleLine() })
		if b.Clicked() {
			*watching = !*watching
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Watchers", "toggled", "watching", *watching)
		}
		for _, person := range people {
			if person == "" {
				panic("project: Watchers person must not be empty")
			}
			ui.Box(c).Size(t.Space(7), t.Space(7)).Radius(999).Background(t.Accent.Alpha(0.15)).Center().Children(func() {
				first, _ := utf8.DecodeRuneInString(person)
				ui.Text(c, string(first)).TextColor(t.Accent)
			})
		}
	})
	return root
}
