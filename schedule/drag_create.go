package schedule

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DragCreateOptions are the drag-create labels.
type DragCreateOptions struct{ Label string }

// DragCreate drags over hours to create an event; cue Pop.
func DragCreate(c *ui.Context, start, end *time.Time, created *bool, o DragCreateOptions) ui.Element {
	if start == nil || start.IsZero() || end == nil || created == nil {
		panic("schedule: DragCreate needs start and state")
	}
	def(&o.Label, "Drag to create")
	t := c.Theme()
	hourH := t.Space(12)
	root := ui.Row(c).FillWidth().Height(2 * hourH).Gap(t.Space(2)).Cursor(ui.CursorGrab).Role(ui.RoleGroup).Label(o.Label)
	visualY := ui.Local(root, "visualY", func() float32 { return 0 })
	moved := ui.Local(root, "moved", func() bool { return false })
	_, dy, held := root.Dragged()
	if held {
		*visualY += dy
		if dy != 0 {
			*moved = true
			*end = start.Add(time.Duration(float32(*visualY)/hourH) * time.Hour)
		}
	}
	span := float32(end.Sub(*start).Hours())
	root.Children(func() {
		ui.Column(c).Width(t.Space(11)).Shrink(0).Children(func() {
			for h := range 2 {
				ui.Box(c).Height(hourH).Children(func() {
					kit.Muted(c, start.Add(time.Duration(h)*time.Hour).Format("15:04")).FontSize(t.Rem(0.65))
				})
			}
		})
		area := ui.Box(c).Grow(1).MinWidth(0).FillHeight().Radius(t.Radius).Clip()
		area.Background(area.AnimateColor("bg", t.SurfaceHover.Alpha(0.4).Mix(t.Accent.Alpha(0.08), unit(held)), ui.Snappy))
		area.Draw(func(p *ui.Painter, r ui.Rect) { p.Line(r.X, r.Y+hourH, r.X+r.W, r.Y+hourH, 1, t.Border) })
		area.Children(func() {
			grow := area.AnimateTo("ghost", max(0, min(span, 2))*hourH, ui.Snappy)
			if grow < 1 {
				ui.Row(c).Absolute().Left(0).Right(0).Top(0).Bottom(0).Center().Gap(t.Space(1.5)).TextColor(t.TextMuted).Children(func() {
					glyph(c, "mouse-pointer-2", 14)
					ui.Text(c, o.Label).FontSize(t.Rem(0.75))
				})
				return
			}
			ghost := ui.Column(c).Absolute().Left(t.Space(1)).Right(t.Space(1)).Top(0).Height(grow).Padding(t.Space(1), t.Space(2)).
				Radius(t.Space(1)).Border(1, t.Accent).Background(t.Accent.Alpha(0.16)).TextColor(t.Accent)
			ghost.Children(func() {
				kit.Numeric(c, start.Format("15:04")+"–"+end.Format("15:04")).FontSize(t.Rem(0.7)).Bold()
			})
		})
	})
	if !held && *moved {
		*moved = false
		*visualY = 0
		if end.After(*start) {
			*created = true
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("DragCreate", "created", "minutes", int(end.Sub(*start).Minutes()))
		}
	}
	return root
}
