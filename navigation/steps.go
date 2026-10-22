package navigation

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Steps are numbered steps joined by filling lines.
func Steps(c *ui.Context, label string, steps []string, current *int) ui.Element {
	if label == "" || len(steps) < 2 || current == nil || !valid(*current, len(steps)) {
		panic(fmt.Sprintf("navigation: Steps %q needs two steps and a current one", label))
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().AlignItems(ui.Start).FocusGroup(ui.Horizontal).Label(label)
	shown := *current
	return root.Children(func() {
		for i, name := range steps {
			c.Key(i)
			if i > 0 {
				line := ui.Row(c).Grow(1).Height(2).Margin(t.Space(4)-1, t.Space(1), 0, t.Space(1)).Radius(999).Background(t.Border)
				k := line.AnimateTo("fill", unit(i <= shown), ui.Smooth)
				line.Clip().Children(func() {
					ui.Box(c).Grow(k).FillHeight().Background(t.Accent)
					ui.Box(c).Grow(1 - k)
				})
			}
			b := kit.Press(ui.ButtonBase(c).Label(name).Checked(i == shown)).Gap(t.Space(1.5))
			done, on := i < shown, i == shown
			bg, edge, fg := t.Background, t.Border, t.TextMuted
			switch {
			case done:
				bg, edge, fg = t.Accent, t.Accent, t.AccentText
			case on:
				edge, fg = t.Accent, t.Accent
			}
			b.Children(func() {
				ui.Column(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					dot := ui.Box(c).Size(t.Space(8), t.Space(8)).Center().Radius(999).Border(2, b.AnimateColor("edge", edge, ui.Snappy))
					dot.Background(dot.AnimateColor("bg", bg, ui.Snappy)).TextColor(dot.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
						if done {
							ui.Icon(c, icon.Must("check")).Size(16, 16).ContentTransition(ui.ContentTransition{})
						} else {
							kit.Numeric(c, strconv.Itoa(i+1)).Bold().ContentTransition(ui.ContentTransition{})
						}
					})
					ui.Text(c, name).FontSize(t.Rem(0.85)).SingleLine().TextColor(b.AnimateColor("name", t.TextMuted.Mix(t.Text, unit(on)), ui.Snappy))
				})
			})
			if kit.Clicked(b, sound.Select) && i != *current {
				*current = i
				root.SetChanged()
				kit.Log("Steps", "picked", "index", i)
			}
		}
	})
}
