package project

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// KeyResult is one objective key result.
type KeyResult struct {
	Title           string
	Current, Target float64
}

// OKROptions word the card; Owner and Period are optional.
type OKROptions struct{ Owner, Period, Label string }

// OKR shows an objective with key result bars.
func OKR(c *ui.Context, objective string, results []KeyResult, o OKROptions) ui.Element {
	if objective == "" || len(results) == 0 {
		panic(fmt.Sprintf("project: OKR %q needs key results", objective))
	}
	def(&o.Label, "KR%d")
	titles := make([]string, len(results))
	score := 0.0
	for i, result := range results {
		if result.Current < 0 || result.Target <= 0 || !finite(result.Current) || !finite(result.Target) {
			panic(fmt.Sprintf("project: OKR %q has an invalid result %q", objective, result.Title))
		}
		titles[i] = result.Title
		score += min(1, result.Current/result.Target)
	}
	unique("OKR", titles)
	score /= float64(len(results))
	t := c.Theme()
	ink := func(v float64) ui.Color {
		switch {
		case v >= 0.7:
			return t.Success
		case v >= 0.4:
			return t.Warning
		}
		return t.Danger
	}
	root := card(c)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start).Children(func() {
				tile(c, "target", t.Accent)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
					ui.Text(c, objective).FontSize(t.Rem(0.9)).Bold().Wrap().FillWidth()
					if o.Owner != "" || o.Period != "" {
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							if o.Owner != "" {
								display.Avatar(c, o.Owner, nil).Size(t.Space(5), t.Space(5)).FontSize(t.Rem(0.5))
								kit.Muted(c, o.Owner).FontSize(t.Rem(0.72)).Shrink(0)
							}
							kit.Muted(c, o.Period).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
						})
					}
				})
				kit.Chip(c, fmt.Sprintf("%d%%", percent(score)), ink(score))
			})
			progress(c, score, ink(score))
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3))
		list.Children(func() {
			for i, result := range results {
				c.Key(result.Title)
				share := min(1, result.Current/result.Target)
				row := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						kit.Muted(c, fmt.Sprintf(o.Label, i+1)).FontSize(t.Rem(0.65)).Bold().Shrink(0).Padding(t.Space(0.25), t.Space(1.5)).
							Radius(999).Background(t.SurfaceHover)
						ui.Text(c, result.Title).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
						value := strconv.FormatFloat(result.Current, 'f', -1, 64) + "/" + strconv.FormatFloat(result.Target, 'f', -1, 64)
						kit.Numeric(c, value).FontSize(t.Rem(0.75)).Bold().Shrink(0).MaxWidth(t.Space(30)).SingleLine().TextColor(t.TextMuted)
					})
					progress(c, share, ink(share))
				})
			}
		})
	})
	return root
}
