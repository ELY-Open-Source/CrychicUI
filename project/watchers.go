package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WatchersOptions word the toggle; You joins the avatars while watching.
type WatchersOptions struct{ Watch, Watching, You, One, Many, None string }

// Watchers toggles watching and shows avatars; cue Toggle.
func Watchers(c *ui.Context, people []string, watching *bool, o WatchersOptions) ui.Element {
	if watching == nil {
		panic("project: Watchers needs watching state")
	}
	def(&o.Watch, "Watch")
	def(&o.Watching, "Watching")
	def(&o.You, "You")
	def(&o.One, "1 watcher")
	def(&o.Many, "%d watchers")
	def(&o.None, "No watchers yet")
	unique("Watchers", append([]string{o.You}, people...))
	t := c.Theme()
	names := people
	if *watching {
		names = append([]string{o.You}, people...)
	}
	count := o.None
	switch {
	case len(names) == 1:
		count = o.One
	case len(names) > 1:
		count = fmt.Sprintf(o.Many, len(names))
	}
	root := card(c).Row().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3))
	root.Children(func() {
		if len(names) > 0 {
			display.AvatarGroup(c, names, display.AvatarGroupOptions{Max: 4})
		}
		kit.Muted(c, count).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
		name, mark := o.Watch, "eye"
		if *watching {
			name, mark = o.Watching, "eye-off"
		}
		c.Key("toggle")
		b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1.5), t.Space(3)).Radius(999).Shrink(0).Checked(*watching))
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.14), unit(*watching)), ui.Snappy))
		b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(*watching)), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.Accent, unit(*watching)), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				c.Key(*watching)
				kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.6, Fade: 1}, 0)
				ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().Shrink(0)
			})
		})
		if b.Clicked() {
			*watching = !*watching
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Watchers", "toggled", "watching", *watching)
		}
	})
	return root
}
