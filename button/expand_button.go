package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ExpandButtonOptions tune an ExpandButton.
type ExpandButtonOptions struct {
	// Label names the action.
	Label string
}

// ExpandButton turns its chevron one hundred eighty degrees.
func ExpandButton(c *ui.Context, expanded *bool, o ExpandButtonOptions) ui.Element {
	if expanded == nil {
		panic("button: ExpandButton needs state")
	}
	if o.Label == "" {
		o.Label = "Expand"
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Tooltip(o.Label).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Center().
		Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Expanded(*expanded))
	root.Children(func() {
		chevron := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16)
		chevron.Rotate(chevron.AnimateTo("turn", 180*unit(*expanded), ui.Snappy))
	})
	if root.Clicked() {
		*expanded = !*expanded
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("ExpandButton", "toggled", "expanded", *expanded)
	}
	return root
}
