package overlay

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ActionSheetOptions tune an ActionSheet.
type ActionSheetOptions struct {
	Cancel string
}

// ActionSheet is a bottom list of actions with Cancel.
func ActionSheet(c *ui.Context, title string, actions []MenuItem, chosen *int, open *bool, o ActionSheetOptions) ui.Element {
	items("ActionSheet", title, actions, chosen)
	if open == nil {
		panic("overlay: ActionSheet needs open")
	}
	def(&o.Cancel, "Cancel")
	t := c.Theme()
	root := ui.Box(c)
	row := func(label string) ui.Element {
		b := kit.Press(ui.ButtonBase(c).Label(label).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius))
		return b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
	}
	if side(c, open, title, bottom(c), ui.Pose{Y: t.Space(60)}, func() {
		kit.Muted(c, title).FontSize(t.Rem(0.8)).AlignSelf(ui.Center)
		ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Children(func() {
			for i, a := range actions {
				c.Key(i)
				b := row(a.Label)
				b.Children(func() {
					if a.Icon != "" {
						ui.Icon(c, icon.Must(a.Icon)).Size(16, 16)
					}
					ui.Text(c, a.Label).Grow(1).SingleLine()
				})
				if b.Clicked() {
					*chosen, *open = i, false
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("ActionSheet", "chose", "index", i)
				}
			}
		})
		cancel := row(o.Cancel).Center().Background(t.SurfaceHover)
		cancel.Children(func() { ui.Text(c, o.Cancel).Bold() })
		if cancel.Clicked() {
			*open = false
			root.SetChanged()
			sound.Play(sound.Close)
		}
	}) {
		root.SetChanged()
	}
	moved(root, *open, "ActionSheet")
	return root
}
