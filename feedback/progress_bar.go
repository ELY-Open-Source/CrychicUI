package feedback

import "github.com/egoist/mygo/ui"

// ProgressBar is a bar whose fill eases to the value.
func ProgressBar(c *ui.Context, label string, value float64) ui.Element {
	if label == "" {
		panic("feedback: ProgressBar needs a label")
	}
	v := fraction("ProgressBar", value)
	t := c.Theme()
	rad := t.Space(0.75)
	root := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(rad).Background(t.Border).Clip().
		Role(ui.RoleProgress).Label(label).Range(0, 1, value)
	return root.Children(func() {
		fill := ui.Box(c).Absolute().Start(0).Top(0).Bottom(0).Radius(rad).Background(t.Accent)
		fill.WidthPercent(fill.AnimateTo("w", v*100, ui.Smooth))
	})
}
