package overlay

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// searchText holds the query; closing clears it and highlight.
func searchText(root ui.Element, open bool) *string {
	q := ui.Local(root, "query", func() string { return "" })
	if !open {
		*q = ""
		forget(root)
	}
	return q
}

// forget moves the highlight back to the first match.
func forget(root ui.Element) { *ui.Local(root, "hi", func() int { return 0 }) = 0 }

// suggest pops results under a focused caller field.
func suggest(c *ui.Context, root, field ui.Element, label, query string, n int, name, glyph func(i int) string, open *bool, choose func(i int), empty string) {
	if !*open {
		forget(root)
		return
	}
	t := c.Theme()
	ui.Overlay(c, func() {
		panel := ui.Column(c).AttachTo(field, ui.AnchorBottomLeft, ui.AnchorTopLeft).KeepFocus().Margin(t.Space(1), 0, 0, 0).
			Width(t.Space(64)).Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border).
			Shadow(0, 8, 24, 0, kit.Shade(0.16)).Role(ui.RoleList).Label(label)
		kit.Arrive(panel, ui.Pose{Y: -t.Space(1), Fade: 1}, 0).Children(func() {
			results(c, root, field, query, n, name, func(int) string { return "" }, glyph, choose, empty)
		})
		if panel.PressedOutside() || field.Shortcut(0, ui.KeyEscape) {
			*open = false
			root.SetChanged()
		}
	})
}

// results filters n labels by query; keys move the highlight.
func results(c *ui.Context, root, field ui.Element, query string, n int, label, group, glyph func(i int) string, choose func(i int), empty string) {
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
					Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							if glyph != nil && glyph(i) != "" {
								ui.Icon(c, icon.Must(glyph(i))).Size(16, 16)
							}
							ui.Text(c, label(i)).Grow(1).SingleLine()
						})
					})
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
