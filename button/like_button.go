package button

import (
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LikeButtonOptions tune a LikeButton.
type LikeButtonOptions struct {
	// Label names the action.
	Label string
}

// LikeButton fills and bursts a heart as its count rolls.
func LikeButton(c *ui.Context, liked *bool, count *int, o LikeButtonOptions) ui.Element {
	if liked == nil || count == nil || *count < 0 {
		panic("button: LikeButton needs state and a non-negative count")
	}
	if o.Label == "" {
		o.Label = "Like"
	}
	t := c.Theme()
	root := textBase(c, o.Label).Border(1, t.Border).Background(t.Surface)
	run := ui.Local(root, "run", func() uint64 { return 0 })
	if root.Clicked() {
		*liked = !*liked
		if *liked {
			*count++
			(*run)++
		} else {
			*count--
		}
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("LikeButton", "toggled", "count", *count)
	}
	root.Children(func() {
		glyph := icon.Must("heart")
		if *liked {
			glyph = icon.Filled("heart")
		}
		heart := ui.Icon(c, glyph).Size(16, 16)
		if *liked {
			heart.TextColor(t.Danger)
		}
		if *run > 0 {
			frames := []ui.Keyframe{{Value: 1}, {At: 80 * time.Millisecond, Value: 1.35}, {At: 200 * time.Millisecond, Value: 1}}
			scale := heart.Keyframes("burst", frames, ui.KeyframeOptions{Run: *run})
			heart.Scale(scale, scale)
		}
		heart.ContentTransition(ui.ContentTransition{})
		kit.Numeric(c, strconv.Itoa(*count)).Padding(t.Space(0.5), t.Space(1)).Radius(999).Background(t.Accent.Alpha(0.12))
	})
	return root
}
