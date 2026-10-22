package navigation

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Step is a step's title and description.
type Step struct{ Title, Detail string }

// StepperNav is vertical steps with descriptions; lines fill.
func StepperNav(c *ui.Context, label string, steps []Step, current *int) ui.Element {
	if label == "" || len(steps) < 2 || current == nil || !valid(*current, len(steps)) {
		panic(fmt.Sprintf("navigation: StepperNav %q needs two steps and a current one", label))
	}
	for _, s := range steps {
		if s.Title == "" {
			panic(fmt.Sprintf("navigation: StepperNav %q steps need titles", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(80)).FocusGroup(ui.Vertical).Label(label)
	shown := *current
	return root.Children(func() {
		for i, s := range steps {
			c.Key(i)
			done, on := i < shown, i == shown
			b := kit.Press(ui.ButtonBase(c).Label(s.Title).FillWidth().Radius(t.Radius).Checked(on))
			b.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch).Children(func() {
					ui.Column(c).AlignItems(ui.Center).Children(func() {
						bg, edge, fg := t.Background, t.Border, t.TextMuted
						switch {
						case done:
							bg, edge, fg = t.Accent, t.Accent, t.AccentText
						case on:
							edge, fg = t.Accent, t.Accent
						}
						dot := ui.Box(c).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999).Border(2, b.AnimateColor("edge", edge, ui.Snappy))
						dot.Background(dot.AnimateColor("bg", bg, ui.Snappy)).TextColor(dot.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
							if done {
								ui.Icon(c, icon.Must("check")).Size(16, 16).ContentTransition(ui.ContentTransition{})
							} else {
								kit.Numeric(c, strconv.Itoa(i+1)).Bold().FontSize(t.Rem(0.8)).ContentTransition(ui.ContentTransition{})
							}
						})
						if i < len(steps)-1 {
							line := ui.Column(c).Grow(1).Width(2).MinHeight(t.Space(4)).Margin(t.Space(1), 0).Radius(999).Background(t.Border).Clip()
							k := line.AnimateTo("fill", unit(done), ui.Smooth)
							line.Children(func() {
								ui.Box(c).Grow(k).FillWidth().Background(t.Accent)
								ui.Box(c).Grow(1 - k)
							})
						}
					})
					ui.Column(c).Grow(1).Gap(t.Space(0.5)).Padding(t.Space(1), 0, t.Space(4), 0).Children(func() {
						title := ui.Text(c, s.Title).Bold().SingleLine()
						title.TextColor(title.AnimateColor("fg", t.TextMuted.Mix(t.Text, unit(done || on)), ui.Snappy))
						if s.Detail != "" {
							kit.Muted(c, s.Detail).FontSize(t.Rem(0.85))
						}
					})
				})
			})
			if kit.Clicked(b, sound.Select) && i != *current {
				*current = i
				root.SetChanged()
				kit.Log("StepperNav", "picked", "index", i)
			}
		}
	})
}
