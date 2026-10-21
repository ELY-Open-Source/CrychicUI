package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SlideConfirmOptions tune a SlideConfirm.
type SlideConfirmOptions struct {
	// Width is the track width in DIPs.
	Width float32
}

// SlideConfirm confirms only when its thumb reaches the end.
func SlideConfirm(c *ui.Context, label string, o SlideConfirmOptions) ui.Element {
	if o.Width == 0 {
		o.Width = 256
	}
	t := c.Theme()
	thumbSize, inset := t.Space(8), t.Space(1)
	limit := o.Width - 2*inset - thumbSize
	if limit <= 0 {
		panic("button: SlideConfirm width is too small")
	}
	root := ui.Row(c).AlignSelf(ui.Start).Width(o.Width).Height(t.Space(10)).Padding(inset).Gap(t.Space(1)).
		AlignItems(ui.Center).Radius(999).Background(t.Surface).Border(1, t.Border).Clip()
	offset := ui.Local(root, "offset", func() float32 { return 0 })
	target := ui.Local(root, "target", func() float32 { return 0 })
	dragging := ui.Local(root, "dragging", func() bool { return false })
	confirmed := ui.Local(root, "confirmed", func() bool { return false })
	root.Children(func() {
		thumb := ui.Column(c).Size(thumbSize, thumbSize).Center().Radius(999).TextColor(t.AccentText)
		thumb.Disabled(*confirmed)
		if !*confirmed {
			if dx, _, ok := thumb.Dragged(); ok {
				if !*dragging {
					*target = *offset
				}
				*dragging = true
				*target = max(0, min(limit, *target+dx))
				*offset = thumb.Jump("slide", *target)
			} else {
				animation := ui.Snappy
				if vx, _, ok := thumb.DragEnded(); ok {
					*dragging = false
					*target = 0
					if *offset >= limit*0.8 {
						*target = limit
					}
					animation.Spring.Velocity = vx
				}
				*offset = thumb.AnimateTo("slide", *target, animation)
			}
		} else {
			*offset = limit
		}
		thumb.Translate(*offset, 0)
		bg := t.Accent
		if *confirmed {
			bg = t.Success
		}
		thumb.Background(thumb.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
			name := "chevron-right"
			if *confirmed {
				name = "check"
			}
			ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
		})
		ui.Text(c, label).Grow(1).TextAlign(ui.Center)
	})
	if !*confirmed && !*dragging && *target == limit && *offset >= limit-0.5 {
		*confirmed = true
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("SlideConfirm", "confirmed")
	}
	return root
}
