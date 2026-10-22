package overlay

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BottomSheetOptions tune a BottomSheet.
type BottomSheetOptions struct {
	// Heights are snap heights, rising; nil gives two.
	Heights []float32
}

// BottomSheet is dragged between snap heights, or down to close.
func BottomSheet(c *ui.Context, title string, open *bool, snap *int, o BottomSheetOptions, content func()) ui.Element {
	t := c.Theme()
	hs := o.Heights
	if hs == nil {
		hs = []float32{t.Space(50), t.Space(100)}
	}
	if title == "" || open == nil || snap == nil || content == nil || !valid(*snap, len(hs)) {
		panic(fmt.Sprintf("overlay: BottomSheet %q needs open, a snap and content", title))
	}
	for i, h := range hs {
		if h <= 0 || math.IsInf(float64(h), 0) || math.IsNaN(float64(h)) || (i > 0 && h <= hs[i-1]) {
			panic(fmt.Sprintf("overlay: BottomSheet %q heights must be finite, positive and rising", title))
		}
	}
	top := hs[len(hs)-1]
	root := ui.Box(c)
	drag := ui.Local(root, "drag", func() ui.Point { return ui.Point{Y: top - hs[*snap]} })
	last := ui.Local(root, "snap", func() int { return *snap })
	if *snap != *last || !*open {
		*last, *drag = *snap, ui.Point{Y: top - hs[*snap]}
	}
	// Rest at a snap height, or fully down to close.
	nearest := func(p ui.Point) ui.Point {
		best := top
		for _, h := range hs {
			if math.Abs(float64(top-h-p.Y)) < math.Abs(float64(best-p.Y)) {
				best = top - h
			}
		}
		return ui.Point{Y: best}
	}
	shut := func() {
		*open = false
		root.SetChanged()
	}
	settle := func(to int) {
		*snap, *last = to, to
		*drag = ui.Point{Y: top - hs[to]}
		root.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("BottomSheet", "snapped", "index", to)
	}
	if *open {
		ui.Overlay(c, func() {
			back := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Justify(ui.End).Modal()
			kit.Arrive(back.Background(kit.Shade(0.3)), ui.Pose{Fade: 1}, 0)
			back.Children(func() {
				sheet := ui.Column(c).FillWidth().Height(top).Padding(t.Space(2), t.Space(5), t.Space(5), t.Space(5)).Gap(t.Space(3)).
					Radius(t.Radius*2).Background(t.Background).Shadow(0, -4, 24, 0, kit.Shade(0.2)).Role(ui.RoleDialog).Label(title).
					Draggable(drag, ui.DraggableOptions{Axis: ui.Vertical, Snap: nearest, Animation: ui.Snappy, Bounds: &ui.Rect{H: top}})
				kit.Arrive(sheet, ui.Pose{Y: top}, 0).Children(func() {
					grip := kit.Press(ui.ButtonBase(c).Label(title).FillWidth().Center().Padding(t.Space(1)))
					grip.Children(func() { ui.Box(c).Size(t.Space(10), t.Space(1)).Radius(999).Background(t.Border) })
					switch {
					case grip.Shortcut(0, ui.KeyUp) && *snap < len(hs)-1:
						settle(*snap + 1)
					case grip.Shortcut(0, ui.KeyDown) && *snap > 0:
						settle(*snap - 1)
					}
					ui.Text(c, title).Bold().FontSize(t.Rem(1.1))
					content()
				})
				if resting(sheet) && drag.Y != top-hs[*snap] {
					if drag.Y >= top {
						sound.Play(sound.Swipe)
						shut()
					} else {
						for i, h := range hs {
							if drag.Y == top-h {
								settle(i)
							}
						}
					}
				}
				if sheet.PressedOutside() || back.OverlayShortcut(0, ui.KeyEscape) {
					shut()
				}
			})
		})
	}
	moved(root, *open, "BottomSheet")
	return root
}
