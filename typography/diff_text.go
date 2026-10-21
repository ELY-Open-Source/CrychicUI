package typography

import "github.com/egoist/mygo/ui"

// DiffText shows deleted and inserted text with changing fades.
func DiffText(c *ui.Context, deleted, inserted string) ui.Element {
	t := c.Theme()
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).Children(func() {
		ui.Text(c, deleted).MinWidth(1).MinHeight(t.FontSize).TextColor(t.Danger).Strikethrough().
			ContentTransition(ui.ContentTransition{})
		ui.Text(c, inserted).MinWidth(1).MinHeight(t.FontSize).TextColor(t.Success).
			ContentTransition(ui.ContentTransition{})
	})
}
