package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// textBase returns the shared text-button geometry.
func textBase(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	return kit.Press(ui.ButtonBase(c).Label(label).AlignSelf(ui.Start).Gap(t.Space(1.5)).Padding(t.Space(1.5), t.Space(3.5)).Radius(t.Radius))
}

// unit converts a boolean to an animation target.
func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}
