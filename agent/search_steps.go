package agent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SearchStep is one phase of a SearchSteps list.
type SearchStep struct {
	Label, Detail string
	Count         int
	State         StepState
}

// SearchSteps shows search phases; the running one spins.
func SearchSteps(c *ui.Context, label string, steps []SearchStep) ui.Element {
	if label == "" || len(steps) == 0 {
		panic(fmt.Sprintf("agent: SearchSteps %q needs a label and steps", label))
	}
	seen := map[string]bool{}
	for _, s := range steps {
		if s.Label == "" || seen[s.Label] || s.Count < 0 || s.State > StepFailed {
			panic(fmt.Sprintf("agent: SearchSteps %q step %q needs a unique label, a count and a known state", label, s.Label))
		}
		seen[s.Label] = true
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).AlignItems(ui.Center).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, s := range steps {
			c.Key(s.Label)
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() { searchStep(c, i, s) })
		}
	})
}

// searchStep draws one phase pill, after a rule unless first.
func searchStep(c *ui.Context, i int, s SearchStep) {
	t := c.Theme()
	if i > 0 {
		ui.Box(c).Size(t.Space(4), 1).Background(t.Border)
	}
	pill := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2.5), t.Space(1), t.Space(1.5)).
		Radius(999).Role(ui.RoleListItem).Label(s.Label)
	kit.Arrive(pill, ui.Pose{X: -t.Space(2), Fade: 1}, time.Duration(i)*60*time.Millisecond)
	on := s.State == StepRunning
	pill.Background(pill.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
	pill.Children(func() {
		mark(c, s.State)
		text := ui.Text(c, s.Label).FontSize(t.Rem(0.85)).Shrink(0)
		if !on {
			text.TextColor(t.TextMuted)
		}
		if s.Count > 0 {
			kit.Numeric(c, strconv.Itoa(s.Count)).FontSize(t.Rem(0.7)).Bold().TextColor(t.Accent).
				Padding(0, t.Space(1.5)).Radius(999).Background(t.Accent.Alpha(0.12))
		}
		if s.Detail != "" {
			kit.Muted(c, s.Detail).FontSize(t.Rem(0.75)).Shrink(0)
		}
	})
}
