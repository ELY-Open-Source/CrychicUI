package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SelectAction is a bulk selection action.
type SelectAction uint8

const (
	SelectNone SelectAction = iota
	SelectDownload
	SelectDelete
	SelectClear
)

// SelectBarOptions are the bulk action labels.
type SelectBarOptions struct {
	Format, Download, Delete, Clear string
}

// SelectBar shows selection count and actions; cue Tap.
func SelectBar(c *ui.Context, count int, action *SelectAction, o SelectBarOptions) ui.Element {
	if count <= 0 || action == nil || *action > SelectClear {
		panic("files: SelectBar needs a positive count and action")
	}
	def(&o.Format, "%d selected")
	def(&o.Download, "Download")
	def(&o.Delete, "Delete")
	def(&o.Clear, "Clear")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
		Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, fmt.Sprintf(o.Format, count)).Bold().Grow(1)
		for _, item := range []struct {
			name   string
			action SelectAction
		}{{o.Download, SelectDownload}, {o.Delete, SelectDelete}, {o.Clear, SelectClear}} {
			b := kit.Press(ui.ButtonBase(c).Label(item.name).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
				Background(t.AccentText.Alpha(0.12)).TextColor(t.AccentText))
			b.Children(func() { ui.Text(c, item.name).FontSize(t.Rem(0.8)).SingleLine() })
			if b.Clicked() {
				*action = item.action
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("SelectBar", "action", "action", *action, "count", count)
			}
		}
	})
	return root
}
