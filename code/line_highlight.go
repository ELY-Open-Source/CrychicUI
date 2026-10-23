package code

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// LineHighlight shows code with marked lines pulsing once.
func LineHighlight(c *ui.Context, lang, code string, marked []int) ui.Element {
	if code == "" {
		panic("code: LineHighlight needs code")
	}
	n := len(linesOf(code))
	seen := map[int]bool{}
	for _, m := range marked {
		if m < 1 || m > n || seen[m] {
			panic(fmt.Sprintf("code: LineHighlight marks must be unique lines from 1 to %d", n))
		}
		seen[m] = true
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
		Padding(t.Space(2), 0).Role(ui.RoleGroup).Label(lang)
	return root.Children(func() {
		codeLines(c, lang, code, func(i int, e ui.Element) {
			if !seen[i+1] {
				return
			}
			pulse := e.Keyframes(fmt.Sprint("pulse.", i), []ui.Keyframe{{Value: 1}, {At: 700 * time.Millisecond, Value: 0, Ease: ui.EaseOut}})
			e.FillWidth().Padding(0, t.Space(2)).Background(t.Warning.Alpha(0.35 * pulse))
		})
	})
}
