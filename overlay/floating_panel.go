package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FloatingPanelOptions tune a FloatingPanel.
type FloatingPanelOptions struct {
	// Height is the area's height; zero gives the default.
	Height float32
}

// FloatingPanel drags anywhere and springs back inside.
func FloatingPanel(c *ui.Context, title string, at *ui.Point, o FloatingPanelOptions, content func()) ui.Element {
	if title == "" || at == nil || content == nil || !finite(at.X) || !finite(at.Y) || !finite(o.Height) || o.Height < 0 {
		panic(fmt.Sprintf("overlay: FloatingPanel %q needs a finite place, height and content", title))
	}
	t := c.Theme()
	h := o.Height
	if h == 0 {
		h = t.Space(80)
	}
	root := ui.Box(c).FillWidth().Height(h).Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
	held := ui.Local(root, "held", func() bool { return false })
	start := ui.Local(root, "start", func() ui.Point { return *at })
	root.Children(func() {
		panel := ui.Column(c).Absolute().Left(0).Top(0).Width(t.Space(60)).Padding(t.Space(3)).Gap(t.Space(2)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).Role(ui.RoleGroup).Label(title)
		room := func() []float32 {
			a, s := root.LayoutBounds(), panel.LayoutBounds()
			return []float32{a.W, max(0, a.W-s.W), max(0, a.H-s.H)}
		}
		used := room()
		kit.Settle(c, root, used, room)
		drag := ui.DraggableOptions{Animation: ui.Bouncy}
		if used[0] > 0 {
			drag.Bounds = &ui.Rect{W: used[1], H: used[2]}
		}
		panel.Draggable(at, drag).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must("grip-horizontal")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, title).Bold().SingleLine()
			})
			content()
		})
		if _, _, down := panel.Dragged(); down && !*held {
			*held, *start = true, *at
		}
		if *held && resting(panel) {
			*held = false
			if *at != *start {
				root.SetChanged()
				sound.Play(sound.Drop)
				kit.Log("FloatingPanel", "moved")
			}
		}
	})
	return root
}
