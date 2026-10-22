package overlay

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// searchText holds the query; closing clears it and highlight.
func searchText(root ui.Element, open bool) *string {
	q := ui.Local(root, "query", func() string { return "" })
	if !open {
		*q = ""
		*ui.Local(root, "hi", func() int { return 0 }) = 0
	}
	return q
}

// results filters n labels by query; keys move the highlight.
func results(c *ui.Context, root, field ui.Element, query string, n int, label, group func(i int) string, choose func(i int), empty string) {
	t := c.Theme()
	hi := ui.Local(root, "hi", func() int { return 0 })
	pointed := ui.Local(root, "pointed", func() int { return -1 })
	seen := ui.Local(root, "seen", func() string { return query })
	reveal := query != *seen
	if reveal {
		*hi, *seen = 0, query
	}
	q := strings.ToLower(strings.TrimSpace(query))
	var matches []int
	for i := range n {
		if strings.Contains(strings.ToLower(label(i)), q) {
			matches = append(matches, i)
		}
	}
	*hi = max(0, min(*hi, len(matches)-1))
	switch {
	case field.Shortcut(0, ui.KeyDown):
		*hi, reveal = min(*hi+1, len(matches)-1), true
	case field.Shortcut(0, ui.KeyUp):
		*hi, reveal = max(*hi-1, 0), true
	case field.Submitted() && len(matches) > 0:
		choose(matches[*hi])
		return
	}
	if len(matches) == 0 {
		kit.Muted(c, empty).Padding(t.Space(3)).Key("empty").Transition(kit.Enter)
		return
	}
	ui.Scroll(c).FillWidth().MaxHeight(t.Space(90)).Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
			// Rows first: the highlight settles before any is drawn.
			rows := make([]ui.Element, len(matches))
			now, last := -1, ""
			for k, i := range matches {
				if g := group(i); g != last {
					c.Key("group " + g)
					kit.Muted(c, g).Bold().FontSize(t.Rem(0.75)).Padding(t.Space(2), t.Space(2), t.Space(1), t.Space(2))
					last = g
				}
				c.Key(i)
				rows[k] = kit.Press(ui.ButtonBase(c).Label(label(i)).FillWidth().Padding(t.Space(2), t.Space(3)).Radius(t.Radius))
				if rows[k].Hovered() {
					if now = i; *pointed != i {
						*hi = k
					}
				}
			}
			*pointed = now
			for k, i := range matches {
				b := rows[k]
				b.Background(b.AnimateColor("bg", t.Accent.Alpha(0.12*unit(k == *hi)), ui.Snappy)).
					Children(func() { ui.Text(c, label(i)).FillWidth().SingleLine() })
				if b.Clicked() {
					choose(i)
				}
			}
			if reveal {
				rows[*hi].ScrollIntoView()
			}
		})
	})
}
