package marketing

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Opening is one open role.
type Opening struct{ Title, Team, Place, Kind string }

// CareersOptions word the list; Team takes the filter.
type CareersOptions struct {
	Title, All, Empty, Count, One string
	Team                          *string
}

// Careers lists open roles by team; cue Tap.
func Careers(c *ui.Context, roles []Opening, opened *string, o CareersOptions) ui.Element {
	def(&o.Title, "Open roles")
	def(&o.All, "All teams")
	def(&o.Empty, "No open roles right now")
	def(&o.Count, "%d open roles")
	def(&o.One, "%d open role")
	titles, teams := make([]string, len(roles)), []string{}
	for i, r := range roles {
		if strings.TrimSpace(r.Team) == "" || strings.TrimSpace(r.Place) == "" || strings.TrimSpace(r.Kind) == "" {
			panic(fmt.Sprintf("marketing: Careers %q role %q needs a team, a place and a kind", o.Title, r.Title))
		}
		titles[i] = r.Title
		if !slices.Contains(teams, r.Team) {
			teams = append(teams, r.Team)
		}
	}
	unique("Careers", titles)
	if opened == nil || *opened != "" && !slices.Contains(titles, *opened) {
		panic(fmt.Sprintf("marketing: Careers %q needs an opened role among its roles", o.Title))
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	team := o.Team
	if team == nil {
		team = ui.Local(root, "team", func() string { return "" })
	}
	if *team != "" && !slices.Contains(teams, *team) {
		*team = ""
		root.SetChanged()
		kit.Log("Careers", "filter reset")
	}
	var shown []Opening
	for _, r := range roles {
		if *team == "" || r.Team == *team {
			shown = append(shown, r)
		}
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			glyph(c, "briefcase", 16).TextColor(t.TextMuted)
			ui.Text(c, o.Title).FontSize(t.Rem(0.98)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Chip(c, counted(len(shown), o.One, o.Count), t.Accent)
		})
		c.Key("filters")
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
			for i, tm := range append([]string{""}, teams...) {
				name, on := tm, tm == *team
				if tm == "" {
					name = o.All
				}
				c.Key(strconv.Quote(tm))
				b := kit.Press(ui.ButtonBase(c).Label(name).Checked(on).Shrink(0).Padding(t.Space(1), t.Space(3)).Radius(999))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
				b.Children(func() {
					ink := t.TextMuted.Mix(t.Text, unit(b.Hovered()))
					if on {
						ink = t.AccentText
					}
					ui.Text(c, name).FontSize(t.Rem(0.76)).Bold().TextColor(ink).Shrink(0).MaxWidth(t.Space(40)).SingleLine()
				})
				if b.Clicked() && !on {
					*team = tm
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Careers", "filtered", "index", i)
				}
			}
		})
		c.Key("list")
		list := ui.Column(c).FillWidth().Gap(t.Space(1))
		list.Children(func() {
			if len(shown) == 0 {
				ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2)).Padding(t.Space(6), 0).Children(func() {
					glyph(c, "briefcase", 28).TextColor(t.TextMuted.Alpha(0.6))
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.84)).Shrink(0)
				})
				return
			}
			for i, r := range shown {
				on := r.Title == *opened
				c.Key(strconv.Quote(r.Title))
				b := kit.Press(ui.ButtonBase(c).Label(r.Title).Checked(on).FillWidth().Row().Gap(t.Space(3)).AlignItems(ui.Center).
					Padding(t.Space(3)).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, 0.08*unit(on)), ui.Snappy))
				b.Border(1, t.Border)
				nudge := b.AnimateTo("nudge", unit(b.Hovered()), ui.Snappy)
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Text(c, r.Title).FontSize(t.Rem(0.88)).Bold().TextColor(t.Text.Mix(t.Accent, unit(on))).FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().MinWidth(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, r.Team).FontSize(t.Rem(0.72)).Bold().TextColor(t.Accent).Shrink(0).MaxWidth(t.Space(30)).SingleLine()
							dot(c)
							glyph(c, "map-pin", 12).TextColor(t.TextMuted)
							kit.Muted(c, r.Place).FontSize(t.Rem(0.72)).MinWidth(0).SingleLine()
							dot(c)
							kit.Muted(c, r.Kind).FontSize(t.Rem(0.72)).Shrink(0).SingleLine()
						})
					})
					glyph(c, "arrow-up-right", 16).TextColor(t.TextMuted).Translate(t.Space(0.5)*nudge, -t.Space(0.5)*nudge)
				})
				if b.Clicked() {
					sound.Play(sound.Tap)
					kit.Log("Careers", "opened", "index", i)
					if !on {
						*opened = r.Title
						root.SetChanged()
					}
				}
			}
		})
	})
	return root
}
