package choice

import (
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var shake = []ui.Keyframe{
	{Value: 0}, {At: 60 * time.Millisecond, Value: -6}, {At: 120 * time.Millisecond, Value: 6},
	{At: 180 * time.Millisecond, Value: -4}, {At: 240 * time.Millisecond, Value: 4}, {At: 300 * time.Millisecond, Value: 0},
}

// PickNOptions tune a PickN.
type PickNOptions struct {
	// Refused runs when a choice past n is refused.
	Refused func()
}

// PickN is chips limited to n; extra choices shake.
func PickN(c *ui.Context, label string, options []string, n int, selected *[]string, o PickNOptions) ui.Element {
	distinct("PickN", options)
	subset("PickN", options, selected)
	if label == "" || n < 1 || n > len(options) || len(*selected) > n {
		panic("choice: PickN needs a label and 1 <= n <= options, selection within n")
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).Label(label).FocusGroup(ui.Horizontal)
	runs := ui.Local(root, "runs", func() map[string]uint64 { return map[string]uint64{} })
	return root.Children(func() {
		for _, option := range options {
			on := slices.Contains(*selected, option)
			shown := on
			c.Key(option)
			chip := ui.ToggleBase(c, &on).Label(option).Padding(t.Space(1), t.Space(3)).Radius(999)
			bg, fg := t.Surface, t.Text
			if shown {
				bg, fg = t.Accent, t.AccentText
			} else if chip.Hovered() {
				bg = t.SurfaceHover
			}
			chip.Background(chip.AnimateColor("bg", bg, quick)).TextColor(chip.AnimateColor("fg", fg, quick))
			if r := (*runs)[option]; r > 0 {
				chip.Translate(chip.Keyframes("shake", shake, ui.KeyframeOptions{Run: r}), 0)
			}
			kit.Press(chip).Children(func() { ui.Text(c, option).SingleLine() })
			if !chip.Changed() {
				continue
			}
			if !shown && len(*selected) >= n {
				(*runs)[option]++
				sound.Play(sound.Error)
				kit.Log("PickN", "refused", "limit", n)
				if o.Refused != nil {
					o.Refused()
				}
				continue
			}
			*selected = flip(*selected, option)
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("PickN", "toggled", "count", len(*selected))
		}
	})
}
