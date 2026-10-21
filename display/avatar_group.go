package display

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AvatarGroupOptions tune an AvatarGroup.
type AvatarGroupOptions struct {
	// Max is how many avatars show before "+N".
	Max int
}

// AvatarGroup overlaps avatars and spreads them on hover.
func AvatarGroup(c *ui.Context, names []string, o AvatarGroupOptions) ui.Element {
	if len(names) == 0 {
		panic("display: AvatarGroup needs names")
	}
	if o.Max == 0 {
		o.Max = 3
	}
	if o.Max < 1 {
		panic("display: AvatarGroup Max must be positive")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center)
	target := -t.Space(2)
	if root.Hovered() {
		target = t.Space(1)
	}
	overlap := root.AnimateTo("spread", target, ui.Snappy)
	return root.Children(func() {
		shown := min(o.Max, len(names))
		for i, name := range names[:shown] {
			c.Key(i)
			a := avatarCircle(c, name, nil, t.Space(8)).Border(1, t.Background)
			if i > 0 {
				a.Margin(0, 0, 0, overlap)
			}
		}
		if len(names) > shown {
			kit.Muted(c, "+"+strconv.Itoa(len(names)-shown)).Padding(t.Space(0.5), t.Space(1)).Radius(999).
				Background(t.Surface).Border(1, t.Border).Margin(0, 0, 0, overlap)
		}
	})
}
