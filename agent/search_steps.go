package agent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SearchStep is one step of a SearchSteps list.
type SearchStep struct {
	Label string
	Count int
	Done  bool
}

// SearchSteps shows search phases; the active one spins.
func SearchSteps(c *ui.Context, label string, steps []SearchStep) ui.Element {
	if label == "" || len(steps) == 0 {
		panic(fmt.Sprintf("agent: SearchSteps %q needs a label and steps", label))
	}
	for _, s := range steps {
		if s.Label == "" || s.Count < 0 {
			panic(fmt.Sprintf("agent: SearchSteps %q steps need a label and a count", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleList).Label(label)
	current := -1
	for i, s := range steps {
		if !s.Done && current < 0 {
			current = i
		}
	}
	return root.Children(func() {
		for i, s := range steps {
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			row.Children(func() {
				switch {
				case s.Done:
					ui.Icon(c, icon.Must("circle-check")).Size(16, 16).TextColor(t.Success)
				case i == current:
					ic := ui.Icon(c, icon.Must("rotate-cw")).Size(16, 16).TextColor(t.Accent)
					ic.Rotate(360 * float32(ic.Loop("spin", time.Second, ui.Linear)))
				default:
					ui.Icon(c, icon.Must("circle")).Size(16, 16).TextColor(t.TextMuted)
				}
				text := s.Label
				if s.Count > 0 {
					text += " " + strconv.Itoa(s.Count)
				}
				out := kit.Numeric(c, text).SingleLine()
				switch {
				case s.Done:
					out.TextColor(t.TextMuted)
				case i != current:
					out.TextColor(t.TextMuted)
				}
			})
		}
	})
}
