package marketing

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HowStep is one numbered step.
type HowStep struct{ Title, Text, Icon string }

// HowItWorksOptions word the section.
type HowItWorksOptions struct{ Title, Text string }

// HowItWorks shows two to six numbered steps.
func HowItWorks(c *ui.Context, steps []HowStep, o HowItWorksOptions) ui.Element {
	if len(steps) < 2 || len(steps) > 6 {
		panic(fmt.Sprintf("marketing: HowItWorks %q needs 2 to 6 steps, got %d", o.Title, len(steps)))
	}
	titles := make([]string, len(steps))
	for i, s := range steps {
		if strings.TrimSpace(s.Icon) == "" {
			panic(fmt.Sprintf("marketing: HowItWorks step %q needs an icon", s.Title))
		}
		titles[i] = s.Title
	}
	unique("HowItWorks", titles)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		grid := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3))
		grid.Children(func() {
			for i, s := range steps {
				c.Key(strconv.Quote(s.Title))
				cell := surface(c).Grow(1).Basis(t.Space(40)).MinWidth(0).Gap(t.Space(2.5)).Padding(t.Space(4))
				kit.Stagger(c, grid, cell, i).Children(func() {
					ui.Row(c).FillWidth().AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
						n := ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText)
						kit.Arrive(n, ui.Pose{Scale: 0.4, Fade: 1}, time.Duration(i)*90*time.Millisecond+150*time.Millisecond).Children(func() {
							ui.Text(c, strconv.Itoa(i+1)).FontSize(t.Rem(0.86)).Bold()
						})
						glyph(c, s.Icon, 18).TextColor(t.TextMuted)
					})
					ui.Text(c, s.Title).FontSize(t.Rem(0.92)).Bold().FillWidth().Wrap()
					if s.Text != "" {
						kit.Muted(c, s.Text).FontSize(t.Rem(0.8)).FillWidth().Wrap()
					}
				})
			}
		})
	})
	return root
}
