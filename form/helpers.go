// Package form holds form layouts, fields and whole forms.
package form

import (
	"regexp"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
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

// plain is a quiet button; its form plays the cue.
func plain(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius))
	bg := t.Surface.Alpha(0)
	if b.Hovered() {
		bg = t.Surface
	}
	return b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() { ui.Text(c, label).SingleLine() })
}

// cross is a remove button; its form plays the cue.
func cross(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius))
	ink := t.TextMuted
	if b.Hovered() {
		ink = t.Danger
	}
	return b.TextColor(b.AnimateColor("ink", ink, ui.Snappy)).Children(func() {
		ui.Icon(c, icon.Must("x")).Size(16, 16)
	})
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

// bar is a bottom button row sliding up when shown.
func bar(c *ui.Context, shown bool, lead func(), back, ahead string) (backButton, aheadButton ui.Element) {
	t := c.Theme()
	row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(4)).
		Background(t.Background).Border(1, t.Border).Disabled(!shown)
	h := row.LayoutBounds().H
	kit.Settle(c, row, []float32{h}, func() []float32 { return []float32{row.LayoutBounds().H} })
	row.Translate(0, row.AnimateTo("slide", h*(1-unit(shown)), ui.Snappy))
	row.Children(func() {
		ui.Box(c).Grow(1).Children(lead)
		backButton = plain(c, back)
		aheadButton = primary(c, ahead)
	})
	return backButton, aheadButton
}

// take keeps the first dropped file and cues Drop.
func take(root ui.Element, path *string, files []string, name string) {
	if len(files) == 0 {
		return
	}
	*path = files[0]
	root.SetChanged()
	sound.Play(sound.Drop)
	kit.Log(name, "dropped", "count", len(files))
}
