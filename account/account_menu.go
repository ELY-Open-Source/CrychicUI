package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MenuItem is one account menu entry.
type MenuItem struct {
	ID, Label, Icon string
}

// AccountMenuOptions word the menu; SignOut sets chosen to "signout".
type AccountMenuOptions struct{ Label, SignOut string }

// AccountMenu is an avatar menu with sign out; cue Open.
func AccountMenu(c *ui.Context, name, email string, items []MenuItem, chosen *string, o AccountMenuOptions) ui.Element {
	if name == "" || email == "" || chosen == nil {
		panic(fmt.Sprintf("account: AccountMenu %q needs a name, email and choice", name))
	}
	ids := make([]string, len(items))
	for i, it := range items {
		if it.Label == "" || it.Icon == "" || it.ID == "signout" {
			panic(fmt.Sprintf("account: AccountMenu item %q needs a label and icon, and signout is reserved", it.ID))
		}
		ids[i] = it.ID
	}
	if len(items) > 0 {
		unique("AccountMenu", ids)
	}
	def(&o.Label, "Account menu")
	def(&o.SignOut, "Sign out")
	t := c.Theme()
	root := ui.Row(c).Role(ui.RoleGroup)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(o.Label).Padding(t.Space(0.75)).Radius(999).Expanded(*open))
		trigger.Background(trigger.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, max(unit(trigger.Hovered()), unit(*open))), ui.Snappy))
		trigger.Border(1, trigger.AnimateColor("line", t.Border.Mix(t.Accent, unit(*open)), ui.Snappy))
		trigger.Children(func() {
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, t.Space(1.5), 0, 0).Children(func() {
				display.Avatar(c, name, nil).Size(t.Space(8), t.Space(8)).FontSize(t.Rem(0.7)).Shrink(0)
				chev := glyph(c, "chevron-down", 14).TextColor(t.TextMuted)
				chev.Rotate(180 * chev.AnimateTo("turn", unit(*open), ui.Snappy))
			})
		})
		if trigger.Clicked() {
			*open = !*open
			if *open {
				sound.Play(sound.Open)
				kit.Log("AccountMenu", "opened")
			}
		}
		ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).Width(t.Space(64)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(1, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Children(func() {
				display.Avatar(c, name, nil).Size(t.Space(10), t.Space(10)).FontSize(t.Rem(0.8)).Shrink(0)
				ui.Column(c).Grow(1).MinWidth(0).Children(func() {
					ui.Text(c, name).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
					kit.Muted(c, email).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
				})
			})
			ui.Box(c).FillWidth().Height(1).Background(t.Border)
			entry := func(label, mark string, ink ui.Color) bool {
				item := kit.Press(ui.ButtonBase(c).Label(label).FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
					Role(ui.RoleMenuItem))
				item.Background(item.AnimateColor("bg", ink.Alpha(0).Mix(ink.Alpha(0.1), unit(item.Hovered())), ui.Snappy))
				item.TextColor(ink).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						glyph(c, mark, 16)
						ui.Text(c, label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					})
				})
				return item.Clicked()
			}
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(48)).Gap(t.Space(0.5)).Role(ui.RoleMenu).Children(func() {
				for i, it := range items {
					c.Key(it.ID)
					if entry(it.Label, it.Icon, t.Text) {
						*chosen, *open = it.ID, false
						root.SetChanged()
						sound.Play(sound.Tap)
						kit.Log("AccountMenu", "chose", "index", i)
					}
				}
			})
			ui.Box(c).FillWidth().Height(1).Background(t.Border)
			if entry(o.SignOut, "log-out", t.Danger) {
				*chosen, *open = "signout", false
				root.SetChanged()
				sound.Play(sound.Close)
				kit.Log("AccountMenu", "signedOut")
			}
		})
	})
	return root
}
