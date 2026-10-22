package overlay

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConfirmPopoverOptions tune a ConfirmPopover.
type ConfirmPopoverOptions struct {
	Question, Yes, No string
}

// ConfirmPopover asks to confirm; yes reports Yes.
func ConfirmPopover(c *ui.Context, anchor ui.Element, open *bool, o ConfirmPopoverOptions) (root, yes ui.Element) {
	if open == nil {
		panic("overlay: ConfirmPopover needs open")
	}
	def(&o.Question, "Are you sure?")
	def(&o.Yes, "Yes")
	def(&o.No, "No")
	t := c.Theme()
	root = ui.Box(c)
	answer := func(said string) {
		*open = false
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("ConfirmPopover", said)
	}
	if ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.Margin(t.Space(1.5), 0, 0, 0).Padding(t.Space(3)).Gap(t.Space(2)).Radius(t.Radius).Background(t.Background).
			Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
		kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
		ui.Text(c, o.Question).Bold()
		ui.Row(c).Gap(t.Space(2)).Justify(ui.End).Children(func() {
			if quiet(c, o.No).Clicked() {
				answer("no")
			}
			yes = primary(c, o.Yes, t.Danger, t.AccentText)
			if yes.Clicked() {
				answer("yes")
			}
		})
	}).Dismissed() {
		root.SetChanged()
	}
	moved(root, *open, "ConfirmPopover")
	return root, yes
}
