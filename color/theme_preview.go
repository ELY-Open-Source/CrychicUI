package color

import "github.com/egoist/mygo/ui"

// ThemePreview draws a small window in th; colors animate.
func ThemePreview(c *ui.Context, label string, th *ui.Theme) ui.Element {
	if label == "" || th == nil {
		panic("color: ThemePreview needs a label and a theme")
	}
	t := c.Theme()
	tint := func(e ui.Element, key string, col ui.Color) ui.Element {
		return e.Background(e.AnimateColor(key, col, ui.Smooth))
	}
	bar := func(key string, w float32, col ui.Color) {
		tint(ui.Box(c).Size(w, t.Space(1.5)).Radius(999), key, col)
	}
	root := ui.Column(c).Size(t.Space(64), t.Space(40)).Radius(t.Radius).Clip().Role(ui.RoleImage).Label(label)
	tint(root, "bg", th.Background).Border(1, root.AnimateColor("edge", th.Border, ui.Smooth))
	return root.Children(func() {
		top := ui.Row(c).FillWidth().Height(t.Space(6)).Gap(t.Space(1)).AlignItems(ui.Center).PaddingX(t.Space(2))
		tint(top, "top", th.Surface).Children(func() {
			for i, col := range []ui.Color{th.Danger, th.Warning, th.Success} {
				c.Key(i)
				tint(ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999), "dot", col)
			}
		})
		ui.Row(c).FillWidth().Grow(1).Children(func() {
			side := ui.Column(c).Width(t.Space(16)).FillHeight().Gap(t.Space(2)).Padding(t.Space(2))
			tint(side, "side", th.Surface).Children(func() {
				bar("s1", t.Space(10), th.TextMuted)
				bar("s2", t.Space(8), th.TextMuted)
				bar("s3", t.Space(9), th.Accent)
			})
			ui.Column(c).Grow(1).Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
				bar("h", t.Space(24), th.Text)
				bar("b1", t.Space(36), th.TextMuted)
				bar("b2", t.Space(30), th.TextMuted)
				btn := ui.Box(c).Size(t.Space(16), t.Space(5)).Radius(th.Radius).Center()
				tint(btn, "btn", th.Accent).Children(func() { bar("label", t.Space(8), th.AccentText) })
			})
		})
	})
}
