package feedback

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Upload is a file name and progress, 0 to 1.
type Upload struct {
	Name     string
	Progress float64
}

var rowMotion = ui.ElementTransition{Animation: &ui.Snappy, Position: true, Exit: &ui.Motion{Collapse: true}}

// UploadProgress is file rows with bars that come and go.
func UploadProgress(c *ui.Context, uploads []Upload) ui.Element {
	names := make([]string, len(uploads))
	for i, u := range uploads {
		fraction("UploadProgress", u.Progress)
		if u.Name == "" || slices.Contains(names[:i], u.Name) {
			panic("feedback: UploadProgress needs distinct file names")
		}
		names[i] = u.Name
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3))
	return root.Children(func() {
		for _, u := range uploads {
			c.Key(u.Name)
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Transition(rowMotion).Children(func() {
				glyph, tint := "file", t.TextMuted
				if u.Progress == 1 {
					glyph, tint = "circle-check", t.Success
				}
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(tint).ContentTransition(ui.ContentTransition{})
				ui.Column(c).Grow(1).Gap(t.Space(1)).Children(func() {
					ui.Row(c).Children(func() {
						ui.Text(c, u.Name).Grow(1).MaxLines(1)
						kit.Numeric(c, fmt.Sprintf("%.0f%%", u.Progress*100)).TextColor(t.TextMuted).FontSize(t.Rem(0.85))
					})
					ProgressBar(c, u.Name, u.Progress)
				})
			})
		}
	})
}
