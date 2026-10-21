package typography

import (
	"strings"
	"unicode"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Monogram is a colored circle of name initials.
func Monogram(c *ui.Context, name string) ui.Element {
	words := strings.Fields(name)
	if len(words) == 0 {
		panic("typography: Monogram needs a name")
	}
	initials := make([]rune, 0, 2)
	var hash uint32
	for _, r := range name {
		hash = hash*31 + uint32(r)
	}
	for _, word := range words {
		for _, r := range word {
			initials = append(initials, unicode.ToUpper(r))
			break
		}
		if len(initials) == 2 {
			break
		}
	}
	t := c.Theme()
	palette := [...]ui.Color{t.Accent, t.Success, t.Warning}
	color := palette[hash%uint32(len(palette))]
	root := ui.Column(c).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Center().Radius(999).
		Background(color.Alpha(0.16)).TextColor(color).Bold().Children(func() {
		ui.Text(c, string(initials))
	})
	return kit.Arrive(root, ui.Pose{Scale: 0.5, Fade: 1}, 0)
}
