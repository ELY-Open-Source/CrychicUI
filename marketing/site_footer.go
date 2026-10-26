package marketing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Link is a labeled URL.
type Link struct{ Label, URL string }

// FooterColumn is a titled group of links.
type FooterColumn struct {
	Title string
	Links []Link
}

// Social is an icon link to a profile.
type Social struct{ Label, Icon, URL string }

// SiteFooterOptions word the footer; Opened takes each opened URL.
type SiteFooterOptions struct {
	Tagline, Legal string
	Socials        []Social
	Opened         *string
}

// SiteFooter shows links, socials and legal text; cue Tap.
func SiteFooter(c *ui.Context, brand string, columns []FooterColumn, o SiteFooterOptions) ui.Element {
	if strings.TrimSpace(brand) == "" || len(columns) == 0 {
		panic(fmt.Sprintf("marketing: SiteFooter %q needs a brand and link columns", brand))
	}
	titles := make([]string, len(columns))
	for i, col := range columns {
		labels := make([]string, len(col.Links))
		for j, l := range col.Links {
			if l.URL == "" {
				panic(fmt.Sprintf("marketing: SiteFooter link %q needs a URL", l.Label))
			}
			labels[j] = l.Label
		}
		unique("SiteFooter", labels)
		titles[i] = col.Title
	}
	unique("SiteFooter", titles)
	names := make([]string, len(o.Socials))
	for i, s := range o.Socials {
		if s.Icon == "" || s.URL == "" {
			panic(fmt.Sprintf("marketing: SiteFooter social %q needs an icon and a URL", s.Label))
		}
		names[i] = s.Label
	}
	unique("SiteFooter", names)
	t := c.Theme()
	var root ui.Element
	open := func(url string, args ...any) {
		c.OpenURL(url)
		if o.Opened != nil {
			*o.Opened = url
		}
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("SiteFooter", "opened", args...)
	}
	root = ui.Column(c).FillWidth().Gap(t.Space(5)).Padding(t.Space(6), t.Space(5)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).
		Role(ui.RoleGroup)
	root.Children(func() {
		top := ui.Row(c).FillWidth().Wrap().Gap(t.Space(6))
		top.Children(func() {
			ui.Column(c).Grow(2).Basis(t.Space(50)).MinWidth(0).Gap(t.Space(2)).Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					tile(c, "sparkles", t.Accent, t.Space(8))
					ui.Text(c, brand).FontSize(t.Rem(1)).Bold().Shrink(0).MaxWidth(t.Space(60)).SingleLine()
				})
				if o.Tagline != "" {
					kit.Muted(c, o.Tagline).FontSize(t.Rem(0.8)).FillWidth().Wrap()
				}
				ui.Row(c).Wrap().Gap(t.Space(1.5)).Padding(t.Space(1), 0, 0, 0).Children(func() {
					for i, s := range o.Socials {
						c.Key(strconv.Quote(s.Label))
						if arrow(c, s.Label, s.Icon).Clicked() {
							open(s.URL, "social", i)
						}
					}
				})
			})
			for i, col := range columns {
				c.Key(strconv.Quote(col.Title))
				group := ui.Column(c).Grow(1).Basis(t.Space(28)).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Start)
				kit.Stagger(c, top, group, i).Children(func() {
					ui.Text(c, col.Title).FontSize(t.Rem(0.78)).Bold().FillWidth().MinWidth(0).Padding(0, 0, t.Space(1), 0).SingleLine()
					for j, l := range col.Links {
						c.Key(strconv.Quote(l.Label))
						b := ui.ButtonBase(c).Label(l.Label).AlignItems(ui.Start).Padding(t.Space(0.5), 0)
						b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(b.Hovered())), ui.Snappy))
						b.Children(func() {
							ui.Text(c, l.Label).FontSize(t.Rem(0.8)).MinWidth(0).SingleLine()
						})
						if b.Clicked() {
							open(l.URL, "column", i, "link", j)
						}
					}
				})
			}
		})
		if o.Legal != "" {
			kit.Muted(c, o.Legal).FontSize(t.Rem(0.72)).FillWidth().Padding(t.Space(4), 0, 0, 0).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Wrap()
		}
	})
	return root
}
