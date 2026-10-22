package data

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Ranked is a label and a value, zero or more.
type Ranked struct {
	Label string
	Value float64
}

// RankListOptions tune a RankList.
type RankListOptions struct {
	// Format prints a value; "%g" when empty.
	Format string
}

// RankList ranks rows with bars by value; bars grow.
func RankList(c *ui.Context, label string, items []Ranked, o RankListOptions) ui.Element {
	if label == "" || len(items) == 0 {
		panic(fmt.Sprintf("data: RankList %q needs items", label))
	}
	var names []string
	top := 0.0
	for _, it := range items {
		if !finite64(it.Value) || it.Value < 0 {
			panic(fmt.Sprintf("data: RankList %q values must be finite and not negative", label))
		}
		names = append(names, it.Label)
		top = max(top, it.Value)
	}
	unique("RankList", label, names)
	def(&o.Format, "%g")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, it := range items {
			c.Key(it.Label)
			frac := float32(0)
			if top > 0 {
				frac = float32(it.Value / top)
			}
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				kit.Numeric(c, strconv.Itoa(i+1)).TextColor(t.TextMuted).Width(t.Space(5)).Shrink(0)
				ui.Text(c, it.Label).Width(t.Space(28)).Shrink(0).SingleLine()
				bar(c, it.Label, frac, t.Accent)
				kit.Numeric(c, fmt.Sprintf(o.Format, it.Value)).Width(t.Space(14)).Shrink(0).TextAlign(ui.End)
			})
		}
	})
}
