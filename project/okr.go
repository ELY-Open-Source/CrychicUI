package project

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// KeyResult is one objective key result.
type KeyResult struct {
	Title           string
	Current, Target float64
}

// OKR shows an objective with key result bars.
func OKR(c *ui.Context, objective string, results []KeyResult) ui.Element {
	if objective == "" || len(results) == 0 {
		panic("project: OKR needs objective and results")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, objective).Bold()
		for i, result := range results {
			if result.Title == "" || result.Current < 0 || result.Target <= 0 || !finite(result.Current) || !finite(result.Target) {
				panic("project: OKR has an invalid result")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, result.Title).Grow(1)
					value := strconv.FormatFloat(result.Current, 'g', -1, 64) + "/" + strconv.FormatFloat(result.Target, 'g', -1, 64)
					kit.Numeric(c, value).TextColor(t.TextMuted)
				})
				progress(c, min(1, result.Current/result.Target), t.Accent)
			})
		}
	})
	return root
}
