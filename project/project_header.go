package project

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProjectHeaderOptions set the icon; Subtitle is optional.
type ProjectHeaderOptions struct{ Icon, Subtitle string }

// ProjectHeader shows icon, name, and tabs; cue Select.
func ProjectHeader(c *ui.Context, name string, tab *string, tabs []string, o ProjectHeaderOptions) ui.Element {
	if name == "" || tab == nil || len(tabs) == 0 {
		panic(fmt.Sprintf("project: ProjectHeader %q needs tabs and state", name))
	}
	unique("ProjectHeader", tabs)
	if !slices.Contains(tabs, *tab) {
		panic(fmt.Sprintf("project: ProjectHeader %q has no tab %q", name, *tab))
	}
	def(&o.Icon, "folder-kanban")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(3), t.Space(2), t.Space(3)).Children(func() {
			tile(c, o.Icon, t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(1.05)).Bold().Wrap().FillWidth()
				if o.Subtitle != "" {
					kit.Muted(c, o.Subtitle).FontSize(t.Rem(0.75)).Wrap().FillWidth()
				}
			})
		})
		ui.ScrollHorizontal(c).FillWidth().BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Row(c).Padding(0, t.Space(1.5)).Children(func() {
				for _, item := range tabs {
					c.Key(item)
					on := *tab == item
					b := kit.Press(ui.ButtonBase(c).Label(item).Column().AlignItems(ui.Stretch).Shrink(0).Checked(on))
					b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
					line := b.AnimateColor("line", t.Accent.Alpha(0).Mix(t.Accent, unit(on)), ui.Snappy)
					b.Children(func() {
						ui.Text(c, item).FontSize(t.Rem(0.8)).Bold().Shrink(0).Padding(t.Space(2), t.Space(2.5))
						ui.Box(c).FillWidth().Height(t.Space(0.5)).Radius(999).Background(line)
					})
					if b.Clicked() && !on {
						*tab = item
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("ProjectHeader", "tab", "tab", *tab)
					}
				}
			})
		})
	})
	return root
}
