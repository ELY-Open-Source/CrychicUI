package marketing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Versus is one feature row; Us and Them mark it.
type Versus struct {
	Feature  string
	Us, Them bool
}

// VersusTableOptions word the table; Yes and No name the marks.
type VersusTableOptions struct{ Title, Yes, No string }

// VersusTable compares us with them.
func VersusTable(c *ui.Context, us, them string, rows []Versus, o VersusTableOptions) ui.Element {
	def(&o.Yes, "Yes")
	def(&o.No, "No")
	if strings.TrimSpace(us) == "" || strings.TrimSpace(them) == "" || us == them || len(rows) == 0 {
		panic(fmt.Sprintf("marketing: VersusTable %q needs two different names and rows", o.Title))
	}
	features := make([]string, len(rows))
	for i, r := range rows {
		features[i] = r.Feature
	}
	unique("VersusTable", features)
	t := c.Theme()
	col := t.Space(20)
	root := card(c).Clip().Role(ui.RoleTable)
	mark := func(has, ours bool) {
		ui.Box(c).Width(col).Shrink(0).AlignSelf(ui.Stretch).Center().Background(t.Accent.Alpha(0.06 * unit(ours))).Children(func() {
			name, ink := "x", t.TextMuted
			if has {
				name, ink = "check", t.Success
			}
			label := o.No
			if has {
				label = o.Yes
			}
			glyph(c, name, 16).TextColor(ink).Label(label)
		})
	}
	root.Children(func() {
		if o.Title != "" {
			ui.Text(c, o.Title).FontSize(t.Rem(1)).Bold().FillWidth().Padding(t.Space(4), t.Space(4), t.Space(2), t.Space(4)).Wrap()
		}
		ui.Row(c).FillWidth().AlignItems(ui.Stretch).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Box(c).Grow(1).MinWidth(0)
			for i, name := range [2]string{us, them} {
				ours := i == 0
				ink := t.TextMuted
				if ours {
					ink = t.Accent
				}
				ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().TextColor(ink).Width(col).Shrink(0).TextAlign(ui.Center).
					Padding(t.Space(2.5), t.Space(1)).Background(t.Accent.Alpha(0.06 * unit(ours))).SingleLine()
			}
		})
		list := ui.Column(c).FillWidth()
		list.Children(func() {
			for i, r := range rows {
				c.Key(strconv.Quote(r.Feature))
				row := ui.Row(c).FillWidth().AlignItems(ui.Stretch).Role(ui.RoleRow)
				if i < len(rows)-1 {
					row.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
				}
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Text(c, r.Feature).FontSize(t.Rem(0.82)).Grow(1).MinWidth(0).Padding(t.Space(2.5), t.Space(4)).Wrap()
					mark(r.Us, true)
					mark(r.Them, false)
				})
			}
		})
	})
	return root
}
