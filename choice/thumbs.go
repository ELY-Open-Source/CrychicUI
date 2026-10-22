package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThumbsOptions tune a Thumbs.
type ThumbsOptions struct {
	// Up and Down name the buttons; "Helpful", "Not helpful".
	Up, Down string
}

// Thumbs is votes of -1, 0 or 1; thumbs bounce.
func Thumbs(c *ui.Context, label string, value *int, o ThumbsOptions) ui.Element {
	if label == "" || value == nil || *value < -1 || *value > 1 {
		panic("choice: Thumbs needs a label and a value from -1 to 1")
	}
	if o.Up == "" {
		o.Up = "Helpful"
	}
	if o.Down == "" {
		o.Down = "Not helpful"
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(1)).AlignSelf(ui.Start).Label(label).FocusGroup(ui.Horizontal)
	runs := ui.Local(root, "runs", func() uint64 { return 0 })
	return root.Children(func() {
		for _, vote := range []int{1, -1} {
			name, glyph, tint := o.Up, "thumbs-up", t.Success
			if vote < 0 {
				name, glyph, tint = o.Down, "thumbs-down", t.Danger
			}
			on := *value == vote
			c.Key(vote)
			b := ui.ButtonBase(c).Label(name).Checked(on).Size(t.Space(8), t.Space(8)).Radius(t.Radius)
			bg, fg := ui.Color{}, t.TextMuted
			if on {
				bg, fg = tint.Alpha(0.14), tint
			} else if b.Hovered() {
				bg, fg = t.SurfaceHover, t.Text
			}
			b.Background(b.AnimateColor("bg", bg, quick)).TextColor(b.AnimateColor("fg", fg, quick))
			b.Children(func() {
				i := ui.Icon(c, icon.Must(glyph)).Size(16, 16)
				if on && *runs > 0 {
					k := i.Keyframes("bounce", bounce, ui.KeyframeOptions{Run: *runs})
					i.Scale(k, k)
				}
			})
			if kit.Clicked(b, sound.Select) {
				*value = vote
				if on {
					*value = 0
				}
				*runs++
				root.SetChanged()
				kit.Log("Thumbs", "voted", "value", *value)
			}
		}
	})
}
