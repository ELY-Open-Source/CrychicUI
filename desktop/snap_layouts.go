package desktop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SnapLayoutsOptions word areas: layout, then area number.
type Snap struct{ Layout, Zone int }

// SnapLayoutsOptions word the areas; Zone formats layout and area.
type SnapLayoutsOptions struct{ Zone string }

// snapZones are areas as fractions: x, y, w, h.
var snapZones = [][][4]float32{
	{{0, 0, 0.5, 1}, {0.5, 0, 0.5, 1}},
	{{0, 0, 0.66, 1}, {0.66, 0, 0.34, 1}},
	{{0, 0, 1.0 / 3, 1}, {1.0 / 3, 0, 1.0 / 3, 1}, {2.0 / 3, 0, 1.0 / 3, 1}},
	{{0, 0, 0.5, 0.5}, {0.5, 0, 0.5, 0.5}, {0, 0.5, 0.5, 0.5}, {0.5, 0.5, 0.5, 0.5}},
}

// SnapLayouts picks an area of a window layout; cue Select.
func SnapLayouts(c *ui.Context, picked *Snap, o SnapLayoutsOptions) ui.Element {
	def(&o.Zone, "Layout %d, area %d")
	if picked == nil || picked.Layout < -1 || picked.Layout >= len(snapZones) ||
		(picked.Layout >= 0 && (picked.Zone < 0 || picked.Zone >= len(snapZones[picked.Layout]))) {
		panic(fmt.Sprintf("desktop: SnapLayouts needs a known layout and area, got %v", picked))
	}
	t := c.Theme()
	w, h, gap := t.Space(26), t.Space(16), t.Space(0.75)
	root := ui.Grid(c).Columns(2).Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius*2).Background(t.SurfaceHover).
		Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.12)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Scale: 0.95, Fade: 1}, 0)
	return root.Children(func() {
		for li, zones := range snapZones {
			c.Key(li)
			frame := ui.Box(c).Size(w, h).Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
			kit.Stagger(c, root, frame, li).Children(func() {
				in, ih := w-2*t.Space(1), h-2*t.Space(1)
				for zi, z := range zones {
					c.Key(zi)
					on := picked.Layout == li && picked.Zone == zi
					b := ui.ButtonBase(c).Label(fmt.Sprintf(o.Zone, li+1, zi+1)).Checked(on).Absolute().
						Left(t.Space(1)+z[0]*in+gap/2).Top(t.Space(1)+z[1]*ih+gap/2).Size(z[2]*in-gap, z[3]*ih-gap).Radius(t.Radius * 0.5)
					b.Background(b.AnimateColor("bg", t.Border.Mix(t.Accent.Alpha(0.5), unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
					if b.Clicked() && !on {
						*picked = Snap{li, zi}
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("SnapLayouts", "picked", "layout", li, "zone", zi)
					}
				}
			})
		}
	})
}
