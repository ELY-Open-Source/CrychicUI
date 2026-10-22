package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// EmojiRatingOptions tune an EmojiRating.
type EmojiRatingOptions struct {
	// Labels name the faces; Terrible to Great when zero.
	Labels [5]string
}

var faces = [5]string{"angry", "frown", "meh", "smile", "laugh"}

// EmojiRating is five faces; the chosen face grows.
func EmojiRating(c *ui.Context, label string, value *int, o EmojiRatingOptions) ui.Element {
	if label == "" || value == nil || *value < 0 || *value > 5 {
		panic("choice: EmojiRating needs a label and a value from 0 to 5")
	}
	if o.Labels == ([5]string{}) {
		o.Labels = [5]string{"Terrible", "Bad", "Okay", "Good", "Great"}
	}
	distinct("EmojiRating", o.Labels[:])
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for i, face := range faces {
			c.Key(i)
			b := ui.RadioBase(c, value, i+1).Label(o.Labels[i]).Padding(t.Space(2)).Radius(999).Tooltip(o.Labels[i])
			on := *value == i+1
			fg, bg := t.TextMuted, ui.Color{}
			if on {
				fg, bg = t.Accent, t.Accent.Alpha(0.12)
			} else if b.Hovered() {
				fg = t.Text
			}
			b.Background(b.AnimateColor("bg", bg, quick)).TextColor(b.AnimateColor("fg", fg, quick))
			b.Children(func() {
				pose := ui.Pose{}
				if on {
					pose.Scale = 1.5
				}
				ui.Icon(c, icon.Must(face)).Size(16, 16).AnimatePose(pose, ui.Bouncy)
			})
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().Gap(t.Space(2))
	if changed {
		root.SetChanged()
		kit.Log("EmojiRating", "rated", "value", *value)
	}
	return root
}
