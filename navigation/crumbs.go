package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CrumbsOptions tune Crumbs.
type CrumbsOptions struct {
	// Max shows at most this many crumbs; 4 when zero.
	Max int
	// More names the menu holding the hidden crumbs.
	More string
}

// Crumbs is a path; its middle folds away.
func Crumbs(c *ui.Context, label string, path []string, chosen *int, o CrumbsOptions) ui.Element {
	if o.Max == 0 {
		o.Max = 4
	}
	if label == "" || len(path) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(path) || o.Max < 3 {
		panic(fmt.Sprintf("navigation: Crumbs %q needs a path, a chosen index or -1, and Max >= 3", label))
	}
	def(&o.More, "More")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	pick := func(i int) {
		*chosen = i
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("Crumbs", "picked", "index", i)
	}
	shown := make([]int, 0, len(path))
	for i := range path {
		if len(path) <= o.Max || i == 0 || i >= len(path)-(o.Max-2) {
			shown = append(shown, i)
		}
	}
	root.Children(func() {
		for n, i := range shown {
			c.Key(i)
			if n > 0 {
				ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
			}
			if n == 1 && i > 1 {
				more := kit.Press(ui.ButtonBase(c).Label(o.More).Tooltip(o.More).Padding(0, t.Space(1.5)).Radius(t.Radius))
				more.Background(more.AnimateColor("bg", t.SurfaceHover.Alpha(unit(more.Hovered())), ui.Snappy)).Children(func() {
					ui.Text(c, "…")
				})
				more.Menu(func(m *ui.Menu) {
					for j := 1; j < i; j++ {
						if m.Item(path[j]).Chosen() {
							pick(j)
						}
					}
				})
				ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
			}
			if i == len(path)-1 {
				ui.Text(c, path[i]).Bold().SingleLine().ContentTransition(ui.ContentTransition{})
				continue
			}
			b := kit.Press(ui.ButtonBase(c).Label(path[i]).Padding(0, t.Space(1)).Radius(t.Radius))
			fg := t.TextMuted
			if b.Hovered() {
				fg = t.Accent
			}
			b.TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Children(func() { ui.Text(c, path[i]).SingleLine() })
			if b.Clicked() {
				pick(i)
			}
		}
	})
	return root
}
