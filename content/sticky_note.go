package content

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StickyNoteOptions dress a note; zero Tilt leans -2 degrees.
type StickyNoteOptions struct {
	Author string
	Tint   int
	Tilt   float32
}

// StickyNote is a tinted paper note set at a tilt.
func StickyNote(c *ui.Context, text string, o StickyNoteOptions) ui.Element {
	if o.Tilt == 0 {
		o.Tilt = -2
	}
	if strings.TrimSpace(text) == "" || o.Tint < 0 || o.Tint > 3 || math.IsNaN(float64(o.Tilt)) || o.Tilt < -8 || o.Tilt > 8 {
		panic(fmt.Sprintf("content: StickyNote %q needs text, tint 0 to 3 and tilt -8 to 8", text))
	}
	t := c.Theme()
	ink := [4]ui.Color{t.Warning, t.Success, t.Accent, t.Danger}[o.Tint]
	root := ui.Column(c).Width(t.Space(56)).MinHeight(t.Space(44)).Gap(t.Space(3)).Padding(t.Space(5), t.Space(4), t.Space(4), t.Space(4)).
		Radius(t.Radius*0.5).Background(t.Background.Mix(ink, 0.22)).Shadow(0, 2, 6, 0, t.Text.Alpha(0.1)).Rotate(o.Tilt).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Scale: 0.9, Rotate: o.Tilt * 3, Fade: 1}, 0)
	return root.Children(func() {
		ui.Box(c).Width(t.Space(14)).Height(t.Space(3)).AlignSelf(ui.Center).Margin(-t.Space(6.5), 0, 0, 0).Radius(t.Space(0.5)).
			Background(t.Background.Alpha(0.6)).Border(1, t.Border.Alpha(0.4))
		ui.Text(c, text).FontSize(t.Rem(0.9)).LineHeight(1.5).TextColor(t.Text).FillWidth().Grow(1).Wrap()
		if o.Author != "" {
			ui.Text(c, o.Author).FontSize(t.Rem(0.72)).Bold().TextColor(t.Text.Alpha(0.6)).FillWidth().MinWidth(0).TextAlign(ui.End).SingleLine()
		}
	})
}
