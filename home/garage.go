package home

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// GarageOptions word the door; OpenDoor and CloseDoor take the name.
type GarageOptions struct{ Open, Closed, Opening, Closing, OpenDoor, CloseDoor string }

// Garage opens and closes a garage door; cue Toggle.
func Garage(c *ui.Context, name string, open *bool, o GarageOptions) ui.Element {
	if strings.TrimSpace(name) == "" || open == nil {
		panic(fmt.Sprintf("home: Garage %q needs an open state", name))
	}
	def(&o.Open, "Open")
	def(&o.Closed, "Closed")
	def(&o.Opening, "Opening…")
	def(&o.Closing, "Closing…")
	def(&o.OpenDoor, "Open %s")
	def(&o.CloseDoor, "Close %s")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		door := ui.Box(c).Size(t.Space(40), t.Space(28)).Shrink(0)
		k := door.AnimateTo("open", unit(*open), ui.Animation{Duration: 1800 * time.Millisecond, Ease: ui.EaseInOut})
		moving := k > 0.01 && k < 0.99 && (k < unit(*open)-0.005 || k > unit(*open)+0.005)
		state, ink := o.Closed, t.Success
		switch {
		case moving && *open:
			state, ink = o.Opening, t.Warning
		case moving:
			state, ink = o.Closing, t.Warning
		case *open:
			state, ink = o.Open, t.Warning
		}
		door.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			roof := t.Space(5)
			var house ui.Path
			house.MoveTo(r.X, r.Y+roof).LineTo(r.X+r.W/2, r.Y).LineTo(r.X+r.W, r.Y+roof).LineTo(r.X+r.W, r.Y+r.H).
				LineTo(r.X, r.Y+r.H).Close()
			p.FillPath(&house, t.SurfaceHover)
			bay := ui.Rect{X: r.X + r.W*0.15, Y: r.Y + roof + t.Space(3), W: r.W * 0.7, H: r.H - roof - t.Space(3)}
			p.Fill(bay, t.Text.Alpha(0.75), 0)
			drop := bay.H * (1 - k)
			const panels = 4
			for i := 0; i < panels; i++ {
				ph := drop / panels
				y := bay.Y + float32(i)*ph
				p.Fill(ui.Rect{X: bay.X, Y: y, W: bay.W, H: max(0, ph-t.Space(0.5))}, t.Background.Mix(t.Border, 0.5), 0)
			}
		}))
		// Own scope: caller words never meet the button key.
		ui.Column(c).AlignItems(ui.Center).Children(func() {
			c.Key(state)
			kit.Arrive(ui.Text(c, state).FontSize(t.Rem(1)).Bold().TextColor(ink).Shrink(0), ui.Pose{Y: t.Space(1), Fade: 1}, 0)
		})
		label, mark := o.OpenDoor, "chevron-up"
		if *open {
			label, mark = o.CloseDoor, "chevron-down"
		}
		c.Key("door")
		b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(label, name)).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.Text, 0.12*unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, mark, 16)
				ui.Text(c, fmt.Sprintf(label, name)).FontSize(t.Rem(0.82)).Bold().Shrink(0).MaxWidth(t.Space(60)).SingleLine()
			})
		})
		if b.Clicked() {
			*open = !*open
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Garage", "door", "open", *open)
		}
	})
	return root
}
