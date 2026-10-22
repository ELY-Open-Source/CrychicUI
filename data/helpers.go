// Package data holds tables, lists, trees and queries.
package data

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// columns checks a table's columns.
func columns(component, label string, cols []Column) {
	if label == "" || len(cols) == 0 {
		panic(fmt.Sprintf("data: %s %q needs columns", component, label))
	}
	for _, col := range cols {
		if col.Title == "" || col.Width < 0 || !finite(col.Width) {
			panic(fmt.Sprintf("data: %s %q columns need titles and finite widths", component, label))
		}
	}
}

// line lays a row of cells after lead.
func line(c *ui.Context, e ui.Element, cols []Column, lead func(), cell func(col int)) ui.Element {
	t := c.Theme()
	return e.FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
		if lead != nil {
			lead()
		}
		for i, col := range cols {
			c.Key(i)
			box := ui.Row(c).AlignItems(ui.Center)
			if col.Width > 0 {
				box.Width(col.Width).Shrink(0)
			} else {
				box.Grow(1)
			}
			if col.End {
				box.Justify(ui.End)
			}
			box.Children(func() { cell(i) })
		}
	})
}

// frame is a table's bordered box.
func frame(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup).Label(label)
}

// hovered is a row tinting on hover or lit.
func hovered(c *ui.Context, lit bool) ui.Element {
	t := c.Theme()
	r := ui.Row(c)
	bg := t.SurfaceHover.Alpha(unit(r.Hovered()))
	if lit {
		bg = t.Accent.Alpha(0.1)
	}
	return r.Background(r.AnimateColor("bg", bg, ui.Snappy))
}

// head shows column titles, sorting when sort is set.
func head(c *ui.Context, root ui.Element, component string, cols []Column, sort *ui.SortOrder, lead func()) {
	t := c.Theme()
	line(c, ui.Row(c).Background(t.Surface), cols, lead, func(i int) {
		title := cols[i].Title
		if sort == nil {
			kit.Muted(c, title).Bold().FontSize(t.Rem(0.8)).SingleLine()
			return
		}
		on := sort.Column == title
		b := kit.Press(ui.ButtonBase(c).Label(title).Gap(t.Space(1)).AlignItems(ui.Center).Radius(t.Radius))
		b.Children(func() {
			ui.Text(c, title).Bold().FontSize(t.Rem(0.8)).TextColor(t.TextMuted).SingleLine()
			arrow := ui.Icon(c, icon.Must("arrow-up")).Size(16, 16).TextColor(t.TextMuted)
			arrow.Rotate(180 * arrow.AnimateTo("turn", unit(on && sort.Descending), ui.Snappy)).
				Opacity(arrow.AnimateTo("show", 0.25+0.75*unit(on), ui.Snappy))
		})
		if b.Clicked() {
			*sort = ui.SortOrder{Column: title, Descending: on && !sort.Descending}
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log(component, "sorted", "column", i, "descending", sort.Descending)
		}
	})
}

// branch is a tree row: indent, arrow, content.
func branch(c *ui.Context, root ui.Element, component, expand, path string, depth int, kids bool, open *ui.Selection[string], content func()) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Padding(0, t.Space(1), 0, t.Space(1)+t.Space(4)*float32(depth)).Children(func() {
		if !kids {
			ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0)
		} else {
			on := open.Has(path)
			b := kit.Press(ui.ButtonBase(c).Label(expand).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(t.Radius)).Expanded(on)
			b.Children(func() { chevron(c, on) })
			if b.Clicked() {
				toggle(open, path)
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log(component, "toggled", "depth", depth, "open", !on)
			}
		}
		content()
	})
}

// kids springs open a branch's children.
func kids(c *ui.Context, open bool, children func()) {
	if open {
		ui.Column(c).Key("kids").FillWidth().Clip().Transition(spring).Children(children)
	}
}

// named checks sibling names are distinct, without "/".
func named(component, label string, names []string) {
	seen := map[string]bool{}
	for _, n := range names {
		if n == "" || strings.Contains(n, "/") || seen[n] {
			panic(fmt.Sprintf("data: %s %q names need text, no /, and no twins", component, label))
		}
		seen[n] = true
	}
}

// chevron is an arrow turning open.
func chevron(c *ui.Context, open bool) {
	ch := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(c.Theme().TextMuted)
	turn := float32(90)
	if c.Direction() == ui.RightToLeft {
		turn = -90
	}
	ch.Rotate(turn * ch.AnimateTo("turn", unit(open), ui.Snappy))
}

// spring bounces a row's height open and shut.
var spring = ui.ElementTransition{Animation: &ui.Bouncy, Enter: &ui.Motion{Collapse: true}, Exit: &ui.Motion{Collapse: true}}

// checkbox is a box whose mark pops in.
func checkbox(c *ui.Context, label string, on *bool, mixed bool) ui.Element {
	t := c.Theme()
	b := ui.CheckboxBase(c, on).Label(label)
	if mixed {
		b.Mixed()
	}
	lit := *on || mixed
	b.Children(func() {
		mark := ui.Box(c).Size(t.Space(4), t.Space(4)).Center().Radius(t.Radius/2).Border(1, t.Border)
		mark.Background(mark.AnimateColor("bg", t.Background.Mix(t.Accent, unit(lit)), ui.Snappy)).Children(func() {
			glyph := "check"
			if mixed {
				glyph = "minus"
			}
			k := mark.AnimateTo("pop", unit(lit), ui.Bouncy)
			ui.Icon(c, icon.Must(glyph)).Size(12, 12).TextColor(t.AccentText).Scale(k, k).ContentTransition(ui.ContentTransition{})
		})
	})
	return b
}

// height is h, or the default for zero.
func height(c *ui.Context, h float32) float32 {
	if h < 0 || !finite(h) {
		panic("data: Height must be finite and not negative")
	}
	if h == 0 {
		return c.Theme().Space(80)
	}
	return h
}

// unique checks items are non-empty and distinct.
func unique(component, label string, items []string) {
	seen := map[string]bool{}
	for _, it := range items {
		if it == "" || seen[it] {
			panic(fmt.Sprintf("data: %s %q items need distinct text", component, label))
		}
		seen[it] = true
	}
}

// move puts rows before row to.
func move(items []string, rows []int, to int) []string {
	picked, rest, k := []string{}, []string{}, to
	for i, it := range items {
		if slices.Contains(rows, i) {
			picked = append(picked, it)
			if i < to {
				k--
			}
		} else {
			rest = append(rest, it)
		}
	}
	return slices.Concat(rest[:k], picked, rest[k:])
}

// toggle flips k in s.
func toggle[K comparable](s *ui.Selection[K], k K) {
	if s.Has(k) {
		s.Remove(k)
	} else {
		s.Add(k)
	}
}

// def fills an empty copy field.
func def(field *string, english string) {
	if *field == "" {
		*field = english
	}
}

// valid reports whether i indexes n items.
func valid(i, n int) bool { return i >= 0 && i < n }

func finite(v float32) bool { return finite64(float64(v)) }

func finite64(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// scalar is a JSON scalar's text and color.
func scalar(t *ui.Theme, v any) (string, ui.Color) {
	switch v := v.(type) {
	case string:
		return strconv.Quote(v), t.Success
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64), t.Accent
	case bool:
		return strconv.FormatBool(v), t.Warning
	case nil:
		return "null", t.TextMuted
	case map[string]any, []any:
		return "", t.TextMuted
	}
	panic(fmt.Sprintf("data: JSONTree cannot show %T", v))
}

// letters names column i: A, B, … AA.
func letters(i int) string {
	s := ""
	for i++; i > 0; i = (i - 1) / 26 {
		s = string(rune('A'+(i-1)%26)) + s
	}
	return s
}

// ruleCopy names a rule row's parts.
type ruleCopy struct{ field, op, value, remove string }

// ruleRow edits one rule, reporting edits and removal.
func ruleRow(c *ui.Context, fields, ops []string, r *Rule, n int, say ruleCopy) (edited, remove bool) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		edited = choice.Dropdown(c, fmt.Sprintf(say.field, n), fields, &r.Field).Changed()
		edited = choice.Dropdown(c, fmt.Sprintf(say.op, n), ops, &r.Op).Changed() || edited
		edited = ui.TextInput(c, &r.Value).Label(fmt.Sprintf(say.value, n)).Grow(1).Changed() || edited
		name := fmt.Sprintf(say.remove, n)
		b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius))
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("x")).Size(16, 16).TextColor(t.TextMuted)
		})
		remove = b.Clicked()
	})
	return edited, remove
}

// rules checks rules use offered fields and operators.
func rules(component, label string, fields, ops []string, rs []Rule) {
	if len(fields) == 0 || len(ops) == 0 {
		panic(fmt.Sprintf("data: %s %q needs fields and operators", component, label))
	}
	unique(component, label, fields)
	unique(component, label, ops)
	for _, r := range rs {
		if !slices.Contains(fields, r.Field) || !slices.Contains(ops, r.Op) {
			panic(fmt.Sprintf("data: %s %q rules need an offered field and operator", component, label))
		}
	}
}

// marked highlights the query's words in s.
func marked(t *ui.Theme, s, query string) []ui.Span {
	low := []rune(s)
	for i, r := range low {
		low[i] = unicode.ToLower(r)
	}
	hit := make([]bool, len(low))
	for _, w := range strings.Fields(query) {
		wr := []rune(w)
		for i, r := range wr {
			wr[i] = unicode.ToLower(r)
		}
		for i := 0; i+len(wr) <= len(low); i++ {
			if slices.Equal(low[i:i+len(wr)], wr) {
				for k := i; k < i+len(wr); k++ {
					hit[k] = true
				}
			}
		}
	}
	rs := []rune(s)
	var spans []ui.Span
	for i := 0; i < len(rs); {
		j := i
		for j < len(rs) && hit[j] == hit[i] {
			j++
		}
		sp := ui.Span{Text: string(rs[i:j])}
		if hit[i] {
			sp.Background, sp.Weight, sp.Radius = t.Warning.Alpha(0.3), 600, t.Space(0.5)
		}
		spans = append(spans, sp)
		i = j
	}
	return spans
}

// pill is a small ghost button; caller cues.
func pill(c *ui.Context, label string) bool {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c.Key(label)).Label(label).Padding(t.Space(1), t.Space(2.5)).Radius(t.Radius))
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).TextColor(t.Accent).
		Children(func() { ui.Text(c, label).FontSize(t.Rem(0.85)) })
	return b.Clicked()
}

// actionsValid checks actions have labels and known icons.
func actionsValid(component string, actions []Action) {
	if len(actions) == 0 {
		panic(fmt.Sprintf("data: %s needs actions", component))
	}
	for _, a := range actions {
		if a.Label == "" {
			panic(fmt.Sprintf("data: %s actions need labels", component))
		}
		icon.Must(a.Icon)
	}
}

// slides glides a list row from its old index i.
func slides(c *ui.Context, row ui.Element, i int) ui.Element {
	slot := ui.Local(row, "slot", func() int { return i })
	from := ui.Local(row, "from", func() float32 { return 0 })
	start := ui.Local(row, "start", func() time.Time { return time.Time{} })
	h := row.LayoutBounds().H
	kit.Settle(c, row, []float32{h}, func() []float32 { return []float32{row.LayoutBounds().H} })
	off := float32(0)
	if k := float32(c.Now().Sub(*start)) / float32(260*time.Millisecond); k < 1 && !c.Preferences().ReduceMotion {
		off = *from * (1 - ui.EaseOut(k))
	}
	if *slot != i {
		// Start from where the row shows now.
		*from, *start, *slot = off+float32(*slot-i)*h, c.Now(), i
		off = *from
	}
	if off != 0 {
		c.AnimationFrame()
	}
	return row.Translate(0, off)
}

// body builds n rows, scrolling within h if set.
func body(c *ui.Context, root ui.Element, list *ui.ListState, n int, h float32, row func(i int)) ui.Element {
	if h < 0 || !finite(h) {
		panic("data: Height must be finite and not negative")
	}
	if h > 0 {
		root.Height(h)
		return ui.List(c, list, n, row).FillWidth().Grow(1)
	}
	// An unsized ui.List clips rows it has not measured yet.
	return ui.Column(c).FillWidth().Children(func() {
		for i := range n {
			if list.Key != nil {
				c.Key(list.Key(i))
			} else {
				c.Key(i)
			}
			ui.Column(c).FillWidth().Children(func() { row(i) })
		}
	})
}
