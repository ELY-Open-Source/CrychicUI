package code

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LineHighlight shows numbered code with marked lines that pulse once.
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
	order := slices.Sorted(slices.Values(marked))
	t := c.Theme()
	root := ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Role(ui.RoleGroup).Label(lang)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		codeBody(c, lang, code, true, func(i int, row ui.Element) {
			if !seen[i+1] {
				return
			}
			k := slices.Index(order, i+1)
			pulse := row.Keyframes("crychic.pulse", []ui.Keyframe{{Value: 1}, {At: 900 * time.Millisecond, Value: 0, Ease: ui.EaseOut}},
				ui.KeyframeOptions{Delay: time.Duration(k) * 90 * time.Millisecond})
			row.Background(t.Warning.Alpha(0.12 + 0.3*pulse)).Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: 3, H: r.H}, t.Warning, 0)
			})
		}, nil)
	})
}
