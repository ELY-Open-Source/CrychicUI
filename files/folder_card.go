package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FolderCardOptions word the item count, "%d" taking it.
type FolderCardOptions struct{ Item, Items string }

// FolderCard opens a folder showing its item count; cue Open.
func FolderCard(c *ui.Context, name string, count int, opened *bool, o FolderCardOptions) ui.Element {
	if name == "" || count < 0 || opened == nil {
		panic(fmt.Sprintf("files: FolderCard %q needs a name, count and state", name))
	}
	def(&o.Item, "%d item")
	def(&o.Items, "%d items")
	format := o.Items
	if count == 1 {
		format = o.Item
	}
	t := c.Theme()
	root := kit.Lift(kit.Press(ui.ButtonBase(c).Label(name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(3)).Radius(t.Radius).Checked(*opened)))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(root.Hovered())), ui.Snappy))
	root.Border(1, root.AnimateColor("line", t.Border.Mix(t.Warning, unit(*opened)), ui.Snappy))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			glyph := "folder"
			if *opened {
				glyph = "folder-open"
			}
			ui.Box(c).Size(t.Space(11), t.Space(11)).Shrink(0).Center().Radius(t.Radius).Background(t.Warning.Alpha(0.14)).TextColor(t.Warning).
				Children(func() { ui.Icon(c, icon.Must(glyph)).Size(22, 22).ContentTransition(ui.ContentTransition{}) })
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).Bold().FillWidth().SingleLine()
				kit.Muted(c, fmt.Sprintf(format, count)).FontSize(t.Rem(0.8))
			})
			ch := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
			ch.Translate(ch.AnimateTo("nudge", t.Space(1)*unit(root.Hovered()), ui.Snappy), 0)
		})
	})
	if root.Clicked() && !*opened {
		*opened = true
		root.SetChanged()
		sound.Play(sound.Open)
		kit.Log("FolderCard", "opened", "items", count)
	}
	return root
}
