package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Tone is a message's kind.
type Tone uint8

// The four tones.
const (
	ToneInfo Tone = iota
	ToneSuccess
	ToneWarning
	ToneDanger
)

// alertBox is a tinted box: tone icon, title, body.
func alertBox(c *ui.Context, tone Tone, title, body string, end func()) ui.Element {
	if title == "" && body == "" {
		panic("feedback: an alert needs a title or a body")
	}
	t := c.Theme()
	col, glyph := look(t, tone)
	root := ui.Row(c).FillWidth().AlignItems(ui.Start).Gap(t.Space(3)).Padding(t.Space(3)).Radius(t.Radius).
		Background(col.Alpha(0.08)).Border(1, col.Alpha(0.3)).Role(ui.RoleStatus).Label(title)
	return root.Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(col).Margin(t.Space(0.5), 0, 0, 0)
		ui.Column(c).Grow(1).Gap(t.Space(1)).Children(func() {
			if title != "" {
				ui.Text(c, title).Bold()
			}
			if body != "" {
				ui.Text(c, body).TextColor(t.Text.Mix(col, 0.2))
			}
		})
		if end != nil {
			end()
		}
	})
}

// Alert is a tinted message box that slides in.
func Alert(c *ui.Context, tone Tone, title, body string) ui.Element {
	return kit.Arrive(alertBox(c, tone, title, body, nil), ui.Pose{Y: -6, Fade: 1}, 0)
}
