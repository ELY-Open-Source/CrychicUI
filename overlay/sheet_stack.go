package overlay

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SheetStackOptions tune a SheetStack.
type SheetStackOptions struct {
	// Close names each sheet's close button.
	Close string
}

// SheetStack stacks depth sheets, lower ones scaled back.
func SheetStack(c *ui.Context, titles []string, depth *int, page func(i int), o SheetStackOptions) ui.Element {
	if len(titles) == 0 || depth == nil || page == nil || *depth < 0 || *depth > len(titles) {
		panic("overlay: SheetStack needs titles, a depth 0..len(titles) and a page")
	}
	for _, s := range titles {
		if s == "" {
			panic("overlay: SheetStack titles need text")
		}
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	last := ui.Local(root, "depth", func() int { return 0 })
	pop := func() {
		*depth--
		root.SetChanged()
	}
	if *depth > 0 {
		ui.Overlay(c, func() {
			back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Modal()
			kit.Arrive(back.Background(kit.Shade(0.3)), ui.Pose{Fade: 1}, 0)
			back.Children(func() {
				for i := range *depth {
					under := float32(*depth - 1 - i)
					c.Key(i)
					sheet := ui.Column(c).Absolute().Left(0).Right(0).Top(t.Space(16)).Bottom(0).Padding(t.Space(5)).Gap(t.Space(3)).
						Radius(t.Radius*2).Background(t.Background).Shadow(0, -4, 24, 0, kit.Shade(0.2)).Role(ui.RoleDialog).Label(titles[i]).
						TransformOrigin(0.5, 0).Disabled(under > 0).Transition(ui.ElementTransition{Duration: 220 * time.Millisecond, Exit: &ui.Motion{Y: t.Space(100)}})
					s := 1 - 0.05*sheet.AnimateTo("scale", under, ui.Snappy)
					sheet.Translate(0, -t.Space(3)*sheet.AnimateTo("lift", under, ui.Snappy)).Scale(s, s)
					kit.Arrive(sheet, ui.Pose{Y: t.Space(100)}, 0).Children(func() {
						if heading(c, titles[i], o.Close).Clicked() {
							pop()
						}
						page(i)
					})
				}
				if back.OverlayShortcut(0, ui.KeyEscape) {
					pop()
				}
			})
		})
	}
	switch {
	case *depth > *last:
		sound.Play(sound.Open)
		kit.Log("SheetStack", "depth", "count", *depth)
	case *depth < *last:
		sound.Play(sound.Close)
		kit.Log("SheetStack", "depth", "count", *depth)
	}
	*last = *depth
	return root
}
