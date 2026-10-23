package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FolderCardOptions are the folder labels.
type FolderCardOptions struct{ Item, Items string }

// FolderCard opens a folder with its item count; cue Open.
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
	root := kit.Lift(ui.ButtonBase(c).Label(name).Row().AlignSelf(ui.Start).Width(t.Space(52)).Gap(t.Space(3)).AlignItems(ui.Center).
		Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Background(t.Surface).Checked(*opened))
	root.Children(func() {
		ui.Icon(c, icon.Must("folder")).Size(t.Space(9), t.Space(9)).TextColor(t.Accent)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, name).Bold().SingleLine()
			kit.Muted(c, fmt.Sprintf(format, count)).FontSize(t.Rem(0.8))
		})
		ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
	})
	if root.Clicked() && !*opened {
		*opened = true
		root.SetChanged()
		sound.Play(sound.Open)
		kit.Log("FolderCard", "opened", "items", count)
	}
	return root
}
