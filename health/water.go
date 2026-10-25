package health

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WaterOptions word the card; Size is millilitres per glass.
type WaterOptions struct {
	Title, Of, Glass, Add string
	Size                  int
}

// Water fills glasses toward a goal; cue Pop.
func Water(c *ui.Context, glasses *int, goal int, o WaterOptions) ui.Element {
	def(&o.Title, "Water")
	if o.Size == 0 {
		o.Size = 250
	}
	if glasses == nil || goal < 1 || goal > 12 || *glasses < 0 || *glasses > 12 || o.Size < 1 || o.Size > 2000 {
		panic(fmt.Sprintf("health: Water %q needs a goal and glasses of 0 to 12, and a glass of 1 to 2000 ml", o.Title))
	}
	def(&o.Of, "%s of %s ml")
	def(&o.Glass, "Glass %d")
	def(&o.Add, "Add a glass")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	set := func(n int) {
		*glasses = n
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Water", "glasses", "count", n)
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "droplets", t.Accent, t.Space(11))
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				kit.Numeric(c, fmt.Sprintf(o.Of, number(c, *glasses*o.Size), number(c, goal*o.Size))).FontSize(t.Rem(1.05)).Bold().
					FillWidth().MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{Numeric: true})
			})
		})
		row := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5))
		row.Children(func() {
			for i := range max(goal, *glasses) {
				full := i < *glasses
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Glass, i+1)).Checked(full).Size(t.Space(9), t.Space(12)).Shrink(0))
				level := b.AnimateTo("level", unit(full), ui.Bouncy)
				hover := b.AnimateTo("hover", unit(b.Hovered()), ui.Snappy)
				b.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
					var glass ui.Path
					inset := r.W * 0.12
					glass.MoveTo(r.X+inset*0.5, r.Y+r.H*0.08).LineTo(r.X+r.W-inset*0.5, r.Y+r.H*0.08).
						LineTo(r.X+r.W-inset, r.Y+r.H*0.94).LineTo(r.X+inset, r.Y+r.H*0.94)
					glass.Close()
					fill := r.H * 0.86 * max(0, min(1, level))
					p.Clip(ui.Rect{X: r.X, Y: r.Y + r.H*0.94 - fill, W: r.W, H: fill}, 0, func() { p.FillPath(&glass, t.Accent.Alpha(0.75)) })
					p.StrokePath(&glass, t.Space(0.5), t.Border.Mix(t.Accent, 0.4+0.6*hover))
				}))
				if b.Clicked() {
					n := i + 1
					if full {
						n = i
					}
					set(n)
				}
			}
		})
		c.Key("add")
		add := kit.Press(ui.ButtonBase(c).Label(o.Add).FillWidth().Padding(t.Space(2), 0).Radius(999).Disabled(*glasses >= 12))
		add.Background(add.AnimateColor("bg", t.Accent.Alpha(0.1).Mix(t.Accent.Alpha(0.2), unit(add.Hovered())), ui.Snappy)).TextColor(t.Accent)
		add.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "plus", 16)
				ui.Text(c, o.Add).FontSize(t.Rem(0.82)).Bold().Shrink(0)
			})
		})
		if add.Clicked() && *glasses < 12 {
			set(*glasses + 1)
		}
	})
	return root
}
