// Package form holds form layouts, fields and whole forms.
package form

import (
	"regexp"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// primary is an accent button whose cue its form plays.
func primary(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(3.5)).Radius(t.Radius)).
		TextColor(t.AccentText)
	bg := t.Accent
	if b.Pressed() {
		bg = t.AccentPressed
	} else if b.Hovered() {
		bg = t.AccentHover
	}
	return b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() { ui.Text(c, label).SingleLine() })
}

// def sets an empty copy field to its English default.
func def(field *string, english string) {
	if *field == "" {
		*field = english
	}
}

// blank reports whether s holds only spaces.
func blank(s string) bool { return strings.TrimSpace(s) == "" }

// when is msg if bad, else empty.
func when(bad bool, msg string) string {
	if bad {
		return msg
	}
	return ""
}

// edited marks root changed when e changed.
func edited(root, e ui.Element) {
	if e.Changed() {
		root.SetChanged()
	}
}

// entry is a FieldGroup whose problem shows once focus leaves.
func entry(c *ui.Context, label, problem string, build func() ui.Element) ui.Element {
	box := ui.Column(c).FillWidth()
	left := ui.Local(box, "left", func() bool { return false })
	was := ui.Local(box, "focused", func() bool { return false })
	return box.Children(func() {
		FieldGroup(c, label, FieldGroupOptions{Error: when(*left, problem)}, func() {
			in := build()
			if *was && !in.Focused() && !*left {
				*left = true
				c.AnimationFrame()
			}
			*was = in.Focused()
		})
	})
}

// rows keeps a stable key per row while n changes.
type rows struct {
	ids  []int
	next int
}

// keys returns n row keys, renumbering after outside changes.
func (r *rows) keys(n int) []int {
	if len(r.ids) != n {
		r.ids = r.ids[:0]
		for range n {
			r.ids = append(r.ids, r.next)
			r.next++
		}
	}
	return r.ids
}

// add appends a key for a new last row.
func (r *rows) add() {
	r.ids = append(r.ids, r.next)
	r.next++
}

// drop removes row i's key.
func (r *rows) drop(i int) { r.ids = append(r.ids[:i], r.ids[i+1:]...) }

// rowMotion slides rows in and collapses them out.
func rowMotion(c *ui.Context) ui.ElementTransition {
	return ui.ElementTransition{Enter: &ui.Motion{Y: -c.Theme().Space(2)}, Exit: &ui.Motion{Collapse: true},
		Position: true}
}

// unit is 1 when on, else 0.
func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}
