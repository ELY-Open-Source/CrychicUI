package motion

import "github.com/egoist/mygo/ui"

// UnderlineGrow grows an underline under text while hovered.
func UnderlineGrow(c *ui.Context, text string) ui.Element {
	if text == "" {
		panic("motion: UnderlineGrow needs text")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Padding(0, 0, t.Space(1), 0)
	on := root.Hovered()
	return root.Children(func() {
		ui.Text(c, text)
		line := ui.Box(c).Absolute().Start(0).Bottom(0).Height(t.Space(0.5)).Radius(t.Space(0.25)).Background(t.Accent)
		line.WidthPercent(line.AnimateTo("grow", 100*unit(on), ui.Snappy))
	})
}
