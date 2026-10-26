package content

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Heading is an outline entry; Depth is 0 to 2.
type Heading struct {
	Title string
	Depth int
}

// DocTOCOptions word the outline.
type DocTOCOptions struct{ Title string }

// DocTOC is a page outline; cue Select.
func DocTOC(c *ui.Context, items []Heading, active *int, o DocTOCOptions) ui.Element {
	def(&o.Title, "On this page")
	titles := make([]string, len(items))
	for i, h := range items {
		if h.Depth < 0 || h.Depth > 2 {
			panic(fmt.Sprintf("content: DocTOC heading %q needs depth 0 to 2", h.Title))
		}
		titles[i] = h.Title
	}
	unique("DocTOC", titles)
	if len(items) == 0 || active == nil || *active < 0 || *active >= len(items) {
		panic(fmt.Sprintf("content: DocTOC %q needs headings and an active one", o.Title))
	}
	t := c.Theme()
	row := t.Space(7)
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).Bold().FillWidth().MinWidth(0).SingleLine()
		list := ui.Column(c).FillWidth()
		at := list.AnimateTo("marker", float32(*active), ui.Snappy)
		list.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: 1, H: r.H}, t.Border, 0)
			p.Fill(ui.Rect{X: r.X, Y: r.Y + at*row + row*0.2, W: t.Space(0.5), H: row * 0.6}, t.Accent, t.Space(0.25))
		}))
		list.Children(func() {
			for i, h := range items {
				c.Key(strconv.Quote(h.Title))
				on := i == *active
				b := ui.ButtonBase(c).Label(h.Title).FillWidth().Height(row).Padding(0, t.Space(2), 0, t.Space(3)+float32(h.Depth)*t.Space(3)).
					AlignItems(ui.Center).Justify(ui.Start).Checked(on)
				ink := t.TextMuted.Mix(t.Text, unit(b.Hovered())).Mix(t.Accent, unit(on))
				b.Children(func() {
					ui.Text(c, h.Title).FontSize(t.Rem(0.8)).TextColor(b.AnimateColor("ink", ink, ui.Snappy)).Grow(1).MinWidth(0).SingleLine()
				})
				if b.Clicked() && !on {
					*active = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("DocTOC", "active", "index", i)
				}
			}
		})
	})
	return root
}
