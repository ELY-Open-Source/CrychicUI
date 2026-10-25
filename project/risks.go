package project

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Risk rates likelihood and impact from 1 to 5.
type Risk struct {
	Title              string
	Likelihood, Impact int
}

// RisksOptions word the card, axes and severity tags.
type RisksOptions struct{ Title, Likelihood, Impact, Low, Medium, High string }

// Risks shows a likelihood by impact grid.
func Risks(c *ui.Context, risks []Risk, o RisksOptions) ui.Element {
	def(&o.Title, "Risks")
	def(&o.Likelihood, "Likelihood")
	def(&o.Impact, "Impact")
	def(&o.Low, "Low")
	def(&o.Medium, "Medium")
	def(&o.High, "High")
	if len(risks) == 0 {
		panic(fmt.Sprintf("project: Risks %q needs risks", o.Title))
	}
	titles := make([]string, len(risks))
	for i, risk := range risks {
		if risk.Likelihood < 1 || risk.Likelihood > 5 || risk.Impact < 1 || risk.Impact > 5 {
			panic(fmt.Sprintf("project: Risks %q has %q off the 5 by 5 grid", o.Title, risk.Title))
		}
		titles[i] = risk.Title
	}
	unique("Risks", titles)
	t := c.Theme()
	level := func(score int) (string, ui.Color) {
		switch {
		case score >= 12:
			return o.High, t.Danger
		case score >= 5:
			return o.Medium, t.Warning
		}
		return o.Low, t.Success
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "shield-alert", o.Title, func() { kit.Chip(c, fmt.Sprint(len(risks)), t.TextMuted) })
		ui.Column(c).FillWidth().Gap(t.Space(0.75)).Padding(t.Space(3)).Children(func() {
			kit.Muted(c, o.Impact).FontSize(t.Rem(0.65)).Bold()
			for impact := 5; impact >= 1; impact-- {
				c.Key(impact)
				ui.Row(c).FillWidth().Gap(t.Space(0.75)).Children(func() {
					kit.Muted(c, fmt.Sprint(impact)).FontSize(t.Rem(0.65)).Width(t.Space(4)).Shrink(0).AlignSelf(ui.Center)
					for likelihood := 1; likelihood <= 5; likelihood++ {
						c.Key(likelihood)
						_, ink := level(likelihood * impact)
						cell := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Height(t.Space(12)).Gap(t.Space(0.5)).Padding(t.Space(0.75)).
							Radius(t.Space(1)).Background(ink.Alpha(0.06 + 0.012*float32(likelihood*impact)))
						cell.Children(func() {
							var here []Risk
							for _, risk := range risks {
								if risk.Likelihood == likelihood && risk.Impact == impact {
									here = append(here, risk)
								}
							}
							for _, risk := range here[:min(len(here), 2)] {
								c.Key(risk.Title)
								chip := ui.Text(c, risk.Title).FontSize(t.Rem(0.6)).Bold().MinWidth(0).SingleLine().Padding(t.Space(0.25), t.Space(1)).
									Radius(999).Background(t.Background).TextColor(ink).Border(1, ink.Alpha(0.4))
								kit.Arrive(chip, ui.Pose{Scale: 0.7, Fade: 1}, 0)
							}
							if len(here) > 2 {
								kit.Muted(c, fmt.Sprintf("+%d", len(here)-2)).FontSize(t.Rem(0.6)).Bold()
							}
						})
					}
				})
			}
			ui.Row(c).FillWidth().Gap(t.Space(0.75)).Children(func() {
				ui.Box(c).Width(t.Space(4)).Shrink(0)
				for likelihood := 1; likelihood <= 5; likelihood++ {
					c.Key(likelihood)
					ui.Row(c).Grow(1).Basis(0).MinWidth(0).Justify(ui.Center).Children(func() {
						kit.Muted(c, fmt.Sprint(likelihood)).FontSize(t.Rem(0.65))
					})
				}
			})
			ui.Row(c).FillWidth().Justify(ui.End).Children(func() { kit.Muted(c, o.Likelihood).FontSize(t.Rem(0.65)).Bold() })
		})
		ranked := slices.Clone(risks)
		slices.SortStableFunc(ranked, func(a, b Risk) int { return cmp.Compare(b.Likelihood*b.Impact, a.Likelihood*a.Impact) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(2), t.Space(2), t.Space(2)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
		list.Children(func() {
			for i, risk := range ranked {
				c.Key(risk.Title)
				name, ink := level(risk.Likelihood * risk.Impact)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(1))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(ink)
					ui.Text(c, risk.Title).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
					kit.Muted(c, fmt.Sprintf("%d×%d", risk.Likelihood, risk.Impact)).FontSize(t.Rem(0.7)).Shrink(0)
					kit.Chip(c, name, ink)
				})
			}
		})
	})
	return root
}
