package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HoldToRecordOptions name the hints.
type HoldToRecordOptions struct {
	Hint, CancelHint string
}

// HoldToRecord records while held; sliding away cancels.
func HoldToRecord(c *ui.Context, label string, recording *bool, onCancel func(), o HoldToRecordOptions) ui.Element {
	if label == "" || recording == nil || onCancel == nil {
		panic(fmt.Sprintf("chat: HoldToRecord %q needs a label, recording and a cancel func", label))
	}
	def(&o.Hint, "Hold to record")
	def(&o.CancelHint, "Release to cancel")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(label)
	held := ui.Local(root, "held", func() bool { return false })
	slid := ui.Local(root, "slid", func() float32 { return 0 })
	cancelled := ui.Local(root, "cancelled", func() bool { return false })
	var btn ui.Element
	root.Children(func() {
		btn = kit.Press(ui.ButtonBase(c).Label(label).Tooltip(o.Hint).Size(t.Space(11), t.Space(11)).Center().Radius(999).
			Checked(*recording))
		fg, bg := t.Text, t.SurfaceHover.Alpha(unit(btn.Hovered()))
		if *recording {
			fg, bg = t.AccentText, t.Danger
		}
		btn.TextColor(btn.AnimateColor("fg", fg, ui.Snappy)).Background(btn.AnimateColor("bg", bg, ui.Snappy)).
			Children(func() { ui.Icon(c, icon.Must("mic")).Size(16, 16) })
		if *recording {
			hint := kit.Muted(c, o.CancelHint).FontSize(t.Rem(0.8))
			kit.Arrive(hint, ui.Pose{Fade: 1}, 0)
		}
	})
	btn.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind == ui.InputPointerDown && ev.Button == 0 {
			*held = true
		}
		return false
	})
	_, dy, down := btn.Dragged()
	if down {
		*slid += dy
	}
	switch {
	case *held && !*recording && !*cancelled:
		*recording = true
		*slid = 0
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("HoldToRecord", "recording", true)
	case *held && *recording && !*cancelled && *slid < -t.Space(14):
		*cancelled = true
		*recording = false
		onCancel()
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("HoldToRecord", "cancelled")
	case *held && !down:
		*held = false
		*slid = 0
		if *recording {
			*recording = false
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("HoldToRecord", "recording", false)
		}
		*cancelled = false
	}
	return root
}
