package finance

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Payee is one saved payee; Detail is an account hint.
type Payee struct{ Name, Detail string }

// PayeesOptions word the card.
type PayeesOptions struct{ Title, Placeholder, Empty string }

// Payees filters saved payees and picks one; cue Select.
func Payees(c *ui.Context, payees []Payee, query *string, selected *int, o PayeesOptions) ui.Element {
	if len(payees) == 0 || query == nil || selected == nil || *selected < -1 || *selected >= len(payees) {
		panic("finance: Payees needs payees, a query and a selection")
	}
	for i, p := range payees {
		if p.Name == "" || slices.IndexFunc(payees, func(x Payee) bool { return x.Name == p.Name }) != i {
			panic(fmt.Sprintf("finance: Payees name %q must be unique and not empty", p.Name))
		}
	}
	def(&o.Title, "Payees")
	def(&o.Placeholder, "Search payees")
	def(&o.Empty, "No payees match")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "users", o.Title, func() { kit.Chip(c, strconv.Itoa(len(payees)), t.TextMuted) })
		ui.Row(c).FillWidth().Padding(t.Space(2), t.Space(2), 0, t.Space(2)).Children(func() {
			if input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder}).Grow(1).Changed() {
				root.SetChanged()
				kit.Log("Payees", "query", "runes", utf8.RuneCountInString(*query))
			}
		})
		q := strings.TrimSpace(*query)
		hits := 0
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5)).FocusGroup(ui.Vertical).Role(ui.RoleList)
		list.Children(func() {
			for i, p := range payees {
				lo, hi := kit.IndexFold(p.Name, q)
				if q != "" && lo < 0 {
					continue
				}
				c.Key(p.Name)
				on := i == *selected
				row := kit.Press(ui.ButtonBase(c).Label(p.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(2)).
					Radius(t.Radius).Checked(on))
				kit.Stagger(c, list, row, hits)
				hits++
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				row.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						display.Avatar(c, p.Name, nil)
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							ui.RichText(c, kit.Marked(t, p.Name, lo, hi)...).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
							if p.Detail != "" {
								kit.Muted(c, p.Detail).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
							}
						})
						check := ui.Icon(c, icon.Must("circle-check")).Size(18, 18).Shrink(0).TextColor(t.Accent)
						check.Opacity(check.AnimateTo("on", unit(on), ui.Snappy))
					})
				})
				if row.Clicked() && !on {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Payees", "selected", "index", i)
				}
			}
			if hits == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(4)).Transition(kit.Enter).
					Children(func() {
						ui.Icon(c, icon.Must("search-x")).Size(20, 20).TextColor(t.TextMuted)
						kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
					})
			}
		})
	})
	return root
}
