package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FormError is a message and its field's handle.
type FormError struct {
	Message string
	Field   *ui.Handle
}

// ErrorListOptions tune an ErrorList.
type ErrorListOptions struct {
	// Title heads the list; "Fix these fields" when empty.
	Title string
}

// ErrorList lists field errors; a click focuses the field.
func ErrorList(c *ui.Context, errors []FormError, o ErrorListOptions) ui.Element {
	if len(errors) == 0 {
		panic("feedback: ErrorList needs errors")
	}
	for _, e := range errors {
		if e.Message == "" || e.Field == nil {
			panic("feedback: ErrorList errors need a message and a field")
		}
	}
	if o.Title == "" {
		o.Title = "Fix these fields"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3)).Radius(t.Radius).
		Background(t.Danger.Alpha(0.06)).Border(1, t.Danger.Alpha(0.3)).Role(ui.RoleStatus).Label(o.Title)
	kit.Arrive(root, ui.Pose{Y: -t.Space(1.5), Fade: 1}, 0)
	return root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).Margin(0, 0, t.Space(1), 0).Children(func() {
			ui.Icon(c, icon.Must("circle-alert")).Size(16, 16).TextColor(t.Danger)
			ui.Text(c, o.Title).Bold()
		})
		for i, e := range errors {
			c.Key(i)
			b := ui.ButtonBase(c).Label(e.Message).Justify(ui.Start).Gap(t.Space(2)).Padding(t.Space(1), t.Space(2)).Radius(t.Radius)
			fg := t.Text.Mix(t.Danger, 0.4)
			if b.Hovered() {
				fg = t.Danger
			}
			b.TextColor(b.AnimateColor("fg", fg, quick)).Children(func() {
				ui.Icon(c, icon.Must("arrow-right")).Size(16, 16)
				ui.Text(c, e.Message).Underline()
			})
			if kit.Clicked(b, sound.Tap) {
				e.Field.Focus(c)
				kit.Log("ErrorList", "focused", "index", i)
			}
		}
	})
}
