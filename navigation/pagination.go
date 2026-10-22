package navigation

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PaginationOptions tune a Pagination.
type PaginationOptions struct {
	// Prev, Next name arrows; Page names a page.
	Prev, Next, Page string
}

// Pagination lists pages with gaps; the current pops.
func Pagination(c *ui.Context, label string, page *int, pages int, o PaginationOptions) ui.Element {
	if label == "" || page == nil || pages < 1 || *page < 1 || *page > pages {
		panic(fmt.Sprintf("navigation: Pagination %q needs a page from 1 to %d", label, pages))
	}
	def(&o.Prev, "Previous")
	def(&o.Next, "Next")
	def(&o.Page, "Page %d")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(1)).FocusGroup(ui.Horizontal).Label(label)
	shown := *page
	to := func(p int) {
		*page = p
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Pagination", "paged", "page", p)
	}
	return root.Children(func() {
		if pageButton(c, o.Prev, "chevron-left", false, shown == 1).Clicked() {
			to(shown - 1)
		}
		for n, p := range slots(shown, pages) {
			if p == 0 {
				c.Key(-n)
				ui.Text(c, "…").TextColor(t.TextMuted).Width(t.Space(6)).TextAlign(ui.Center)
				continue
			}
			c.Key(p)
			b := pageButton(c, fmt.Sprintf(o.Page, p), "", p == shown, false).Checked(p == shown)
			pop(b, p == shown)
			b.Children(func() { kit.Numeric(c, strconv.Itoa(p)) })
			if b.Clicked() && p != shown {
				to(p)
			}
		}
		if pageButton(c, o.Next, "chevron-right", false, shown == pages).Clicked() {
			to(shown + 1)
		}
	})
}

// pageButton is a square page or arrow button.
func pageButton(c *ui.Context, name, glyph string, on, off bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(name).Size(t.Space(9), t.Space(9)).Shrink(0).Center().Radius(t.Radius)).Disabled(off)
	bg, fg := t.SurfaceHover.Alpha(0), t.Text
	switch {
	case on:
		bg, fg = t.Accent, t.AccentText
	case b.Hovered():
		bg = t.SurfaceHover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy))
	if glyph != "" {
		b.Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16) })
	}
	return b
}

// slots lists pages around page; 0 is a gap.
func slots(page, pages int) []int {
	keep := []int{1, pages, page - 1, page, page + 1}
	if page <= 4 {
		keep = append(keep, 2, 3, 4, 5)
	}
	if page >= pages-3 {
		keep = append(keep, pages-4, pages-3, pages-2, pages-1)
	}
	if pages <= 7 {
		keep = keep[:0]
		for p := 1; p <= pages; p++ {
			keep = append(keep, p)
		}
	}
	slices.Sort(keep)
	var out []int
	for _, p := range slices.Compact(keep) {
		if p < 1 || p > pages {
			continue
		}
		if len(out) > 0 && p > out[len(out)-1]+1 {
			out = append(out, 0)
		}
		out = append(out, p)
	}
	return out
}

// pop springs e from small whenever on turns true.
func pop(e ui.Element, on bool) {
	run := ui.Local(e, "run", func() uint64 { return 0 })
	was := ui.Local(e, "was", func() bool { return on })
	if on && !*was {
		*run++
	}
	*was = on
	if *run == 0 {
		return
	}
	k := e.Keyframes("crychic.pop", []ui.Keyframe{{Value: 1}, {At: time.Millisecond, Value: 0.8},
		{At: 260 * time.Millisecond, Value: 1, Ease: ui.EaseOut}}, ui.KeyframeOptions{Run: *run})
	e.Scale(k, k)
}
