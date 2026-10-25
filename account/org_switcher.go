package account

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Org is one organization; Logo is optional.
type Org struct {
	ID, Name, Plan string
	Logo           ui.ImageSource
}

// OrgSwitcherOptions word the switcher.
type OrgSwitcherOptions struct{ Label, Title string }

// OrgSwitcher switches organizations; cue Open, Select.
func OrgSwitcher(c *ui.Context, orgs []Org, current *string, o OrgSwitcherOptions) ui.Element {
	if len(orgs) == 0 || current == nil {
		panic("account: OrgSwitcher needs organizations and a current one")
	}
	ids := make([]string, len(orgs))
	at := -1
	for i, org := range orgs {
		if org.Name == "" || org.Plan == "" {
			panic(fmt.Sprintf("account: OrgSwitcher %q needs a name and plan", org.ID))
		}
		ids[i] = org.ID
		if org.ID == *current {
			at = i
		}
	}
	unique("OrgSwitcher", ids)
	if at < 0 {
		panic(fmt.Sprintf("account: OrgSwitcher current %q is not an organization", *current))
	}
	def(&o.Label, "Switch organization")
	def(&o.Title, "Organizations")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(o.Label).FillWidth().Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).Expanded(*open))
		trigger.Background(trigger.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, max(unit(trigger.Hovered()), unit(*open))), ui.Snappy))
		trigger.Border(1, trigger.AnimateColor("line", t.Border.Mix(t.Accent, unit(*open)), ui.Snappy))
		trigger.Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
				c.Key(orgs[at].ID)
				kit.Arrive(logo(c, orgs[at], t.Space(9)), ui.Pose{Scale: 0.7, Fade: 1}, 0)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, orgs[at].Name).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
					kit.Muted(c, orgs[at].Plan).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
				})
				glyph(c, "chevrons-up-down", 16).TextColor(t.TextMuted)
			})
		})
		if trigger.Clicked() {
			*open = !*open
			if *open {
				sound.Play(sound.Open)
				kit.Log("OrgSwitcher", "opened")
			}
		}
		ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(64)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).Bold().Padding(t.Space(1), t.Space(2))
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(0.5)).Role(ui.RoleMenu).Children(func() {
				for i, org := range orgs {
					c.Key(org.ID)
					on := org.ID == *current
					item := kit.Press(ui.ButtonBase(c).Label(org.Name).FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
						Role(ui.RoleMenuItem).Checked(on))
					item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy))
					item.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							logo(c, org, t.Space(7))
							ui.Text(c, org.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							glyph(c, "check", 14).TextColor(t.Accent).Opacity(unit(on))
						})
					})
					if item.Clicked() {
						*open = false
						if !on {
							*current = org.ID
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("OrgSwitcher", "switched", "index", i)
						}
					}
				}
			})
		})
	})
	return root
}

// logo is an organization image or its initials.
func logo(c *ui.Context, org Org, size float32) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(t.Radius).Clip().Background(t.Accent.Alpha(0.14)).TextColor(t.Accent)
	return box.Children(func() {
		if org.Logo != nil {
			ui.Image(c, org.Logo).Fill().Fit(ui.Cover)
			return
		}
		letters := ""
		for _, w := range strings.Fields(org.Name) {
			if r := []rune(w)[0]; unicode.IsLetter(r) && len([]rune(letters)) < 2 {
				letters += strings.ToUpper(string(r))
			}
		}
		if letters == "" {
			glyph(c, "building-2", size*0.5)
			return
		}
		ui.Text(c, letters).FontSize(size * 0.38).Bold().Shrink(0)
	})
}
