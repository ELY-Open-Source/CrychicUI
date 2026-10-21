package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FollowButtonOptions tune a FollowButton.
type FollowButtonOptions struct {
	// Follow and Following label the two states.
	Follow, Following string
}

// FollowButton swaps labels with a width transition.
func FollowButton(c *ui.Context, following *bool, o FollowButtonOptions) ui.Element {
	if following == nil {
		panic("button: FollowButton needs state")
	}
	if o.Follow == "" {
		o.Follow = "Follow"
	}
	if o.Following == "" {
		o.Following = "Following"
	}
	t := c.Theme()
	root := textBase(c, o.Follow).Border(1, t.Border).Transition(ui.ElementTransition{Size: true})
	root.Children(func() {
		label, bg, fg := o.Follow, t.Surface, t.Text
		if *following {
			label, bg, fg = o.Following, t.Accent.Alpha(0.14), t.Accent
		}
		root.Label(label).Background(root.AnimateColor("bg", bg, ui.Snappy)).TextColor(fg)
		ui.Text(c, label).SingleLine().ContentTransition(ui.ContentTransition{})
	})
	if root.Clicked() {
		*following = !*following
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("FollowButton", "toggled", "following", *following)
	}
	return root
}
