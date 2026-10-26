package content

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DocVersion is a docs release; one may be Latest.
type DocVersion struct {
	Name, Note string
	Latest     bool
}

// VersionSelectOptions word the picker.
type VersionSelectOptions struct{ Label, Title, Latest string }

// VersionSelect picks a docs version; opening reports; cue Select.
func VersionSelect(c *ui.Context, versions []DocVersion, current *string, o VersionSelectOptions) ui.Element {
	def(&o.Label, "Docs version")
	def(&o.Title, "Versions")
	def(&o.Latest, "Latest")
	names, latest := make([]string, len(versions)), 0
	for i, v := range versions {
		names[i] = v.Name
		if v.Latest {
			latest++
		}
	}
	unique("VersionSelect", names)
	if current == nil || !slices.Contains(names, *current) || latest > 1 {
		panic(fmt.Sprintf("content: VersionSelect %q needs a current version and at most one latest", o.Label))
	}
	at := slices.Index(names, *current)
	t := c.Theme()
	root := ui.Column(c).AlignItems(ui.Start).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	open := ui.Local(root, "open", func() bool { return false })
	tag := func(v DocVersion) {
		if v.Latest {
			ui.Text(c, o.Latest).FontSize(t.Rem(0.62)).Bold().TextColor(t.Success).Padding(t.Space(0.25), t.Space(1.5)).Radius(999).
				Background(t.Success.Alpha(0.12)).Shrink(0).MaxWidth(t.Space(20)).SingleLine()
		}
	}
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(o.Label).Padding(t.Space(1.5), t.Space(3)).Radius(999).Expanded(*open))
		trigger.Background(trigger.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, max(unit(trigger.Hovered()), unit(*open))), ui.Snappy))
		trigger.Border(1, trigger.AnimateColor("line", t.Border.Mix(t.Accent, unit(*open)), ui.Snappy))
		trigger.Children(func() {
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				glyph(c, "git-branch", 14).TextColor(t.Accent)
				c.Key(strconv.Quote(*current))
				kit.Arrive(ui.Text(c, *current).Font("monospace").FontSize(t.Rem(0.8)).Bold().MinWidth(0).MaxWidth(t.Space(36)).SingleLine(),
					ui.Pose{Y: t.Space(1), Fade: 1}, 0)
				tag(versions[at])
				glyph(c, "chevrons-up-down", 14).TextColor(t.TextMuted)
			})
		})
		if trigger.Clicked() {
			*open = !*open
			if *open {
				root.SetChanged()
				sound.Play(sound.Open)
				kit.Log("VersionSelect", "opened")
			}
		}
		ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(56)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).Bold().Padding(t.Space(1), t.Space(2))
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(0.5)).Role(ui.RoleMenu).Children(func() {
				for i, v := range versions {
					c.Key(strconv.Quote(v.Name))
					on := i == at
					item := kit.Press(ui.ButtonBase(c).Label(v.Name).FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
						Role(ui.RoleMenuItem).Checked(on))
					item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy))
					item.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
								ui.Text(c, v.Name).Font("monospace").FontSize(t.Rem(0.82)).Bold().FillWidth().MinWidth(0).SingleLine()
								if v.Note != "" {
									kit.Muted(c, v.Note).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
								}
							})
							tag(v)
							glyph(c, "check", 14).TextColor(t.Accent).Opacity(unit(on))
						})
					})
					if item.Clicked() {
						*open = false
						if !on {
							*current = v.Name
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("VersionSelect", "switched", "index", i)
						}
					}
				}
			})
		})
	})
	return root
}
