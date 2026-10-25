package health

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BodyMapOptions word the figure; Names run head to left leg.
type BodyMapOptions struct {
	Title, Prompt string
	Names         [7]string
}

// BodyMap picks a region on a front-facing figure; cue Select.
func BodyMap(c *ui.Context, region *int, o BodyMapOptions) ui.Element {
	def(&o.Title, "Where does it hurt?")
	if region == nil || *region < -1 || *region > 6 {
		panic(fmt.Sprintf("health: BodyMap %q needs a region of -1 to 6", o.Title))
	}
	def(&o.Prompt, "Tap an area")
	for i, n := range [...]string{"Head", "Chest", "Abdomen", "Right arm", "Left arm", "Right leg", "Left leg"} {
		def(&o.Names[i], n)
	}
	t := c.Theme()
	u := t.Space(1)
	// Front view: the person's right is the viewer's left.
	shapes := [7][5]float32{
		{17, 0, 10, 10, 999}, {12, 12, 20, 14, 3}, {13, 27, 18, 13, 3},
		{4, 12, 7, 26, 3}, {33, 12, 7, 26, 3}, {13, 41, 8.5, 28, 3}, {22.5, 41, 8.5, 28, 3},
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Title).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine()
		figure := ui.Box(c).Size(44*u, 70*u).Shrink(0).Role(ui.RoleRadioGroup)
		figure.Children(func() {
			for i, s := range shapes {
				on := *region == i
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(o.Names[i]).Role(ui.RoleRadio).Checked(on).Absolute().Left(s[0]*u).Top(s[1]*u).
					Size(s[2]*u, s[3]*u).Radius(min(s[4]*u, 999)))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Mix(t.Border, unit(b.Hovered())).Mix(t.Danger.Alpha(0.75), unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Danger, unit(on)), ui.Snappy))
				kit.Arrive(b, ui.Pose{Scale: 0.8, Fade: 1}, 0)
				if b.Clicked() && !on {
					*region = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("BodyMap", "picked", "region", i)
				}
			}
		})
		c.Key(*region)
		if *region >= 0 {
			kit.Arrive(kit.Chip(c, o.Names[*region], t.Danger), ui.Pose{Scale: 0.7, Fade: 1}, 0)
		} else {
			kit.Muted(c, o.Prompt).FontSize(t.Rem(0.78)).Shrink(0)
		}
	})
	return root
}
