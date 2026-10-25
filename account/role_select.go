package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Role is one assignable role; Icon is optional.
type Role struct {
	Name, Description, Icon string
}

// RoleSelect picks a role with its description; cue Select.
func RoleSelect(c *ui.Context, label string, roles []Role, chosen *string) ui.Element {
	if label == "" || len(roles) == 0 || chosen == nil {
		panic(fmt.Sprintf("account: RoleSelect %q needs a label, roles and a choice", label))
	}
	names := make([]string, len(roles))
	for i, r := range roles {
		if r.Description == "" {
			panic(fmt.Sprintf("account: RoleSelect %q role %q needs a description", label, r.Name))
		}
		names[i] = r.Name
	}
	unique("RoleSelect", names)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleRadioGroup)
	root.Children(func() {
		ui.Text(c, label).FontSize(t.Rem(0.78)).Bold().TextColor(t.TextMuted)
		for i, r := range roles {
			on := *chosen == r.Name
			c.Key(r.Name)
			b := kit.Press(ui.ButtonBase(c).Label(r.Name).FillWidth().Padding(t.Space(3)).Radius(t.Radius).Checked(on))
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.06), unit(on)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			kit.Stagger(c, root, b, i).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start).Children(func() {
					mark := r.Icon
					def(&mark, "user")
					ink := t.TextMuted
					if on {
						ink = t.Accent
					}
					tile(c, mark, ink, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, r.Name).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, r.Description).FontSize(t.Rem(0.78)).Wrap().FillWidth()
					})
					ring := ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Radius(999)
					ring.Border(t.Space(0.5), ring.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
					ring.Children(func() {
						dot := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Background(t.Accent)
						k := dot.AnimateTo("dot", unit(on), ui.Bouncy)
						dot.Scale(k, k)
					})
				})
			})
			if b.Clicked() && !on {
				*chosen = r.Name
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("RoleSelect", "selected", "index", i)
			}
		}
	})
	return root
}
