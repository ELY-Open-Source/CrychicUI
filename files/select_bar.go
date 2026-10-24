package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SelectAction is a bulk selection action.
type SelectAction uint8

// Bulk actions.
const (
	SelectNone SelectAction = iota
	SelectDownload
	SelectDelete
	SelectClear
)

// SelectBarOptions word the bar; Format takes the count.
type SelectBarOptions struct {
	Format, Download, Delete, Clear string
}

// SelectBar is a selection count with actions; cue Tap.
func SelectBar(c *ui.Context, count int, action *SelectAction, o SelectBarOptions) ui.Element {
	if count <= 0 || action == nil || *action > SelectClear {
		panic(fmt.Sprintf("files: SelectBar needs a positive count and a known action, got %d", count))
	}
	def(&o.Format, "%d selected")
	def(&o.Download, "Download")
	def(&o.Delete, "Delete")
	def(&o.Clear, "Clear")
	t := c.Theme()
	bg, fg := inverse(t)
	root := ui.Row(c).AlignSelf(ui.Center).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(2)).
		Radius(999).Background(bg).TextColor(fg).Shadow(0, 1, 3, 0, kit.Shade(0.06)).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).
		Label(fmt.Sprintf(o.Format, count))
	kit.Arrive(root, ui.Pose{Y: t.Space(4), Scale: 0.96, Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, t.Space(2), 0, 0).Children(func() {
			ui.Icon(c, icon.Must("square-check")).Size(16, 16).TextColor(t.Accent)
			ui.Text(c, fmt.Sprintf(o.Format, count)).FontSize(t.Rem(0.85)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		ui.Box(c).Size(1, t.Space(5)).Background(fg.Alpha(0.2))
		for _, it := range []struct {
			name, glyph string
			action      SelectAction
			ink         ui.Color
		}{{o.Download, "download", SelectDownload, fg}, {o.Delete, "trash-2", SelectDelete, t.Danger.Mix(fg, 0.25)}, {o.Clear, "x", SelectClear, fg}} {
			c.Key(it.action)
			b := kit.Press(ui.ButtonBase(c).Label(it.name).Padding(t.Space(1.5), t.Space(2.5)).Radius(999))
			b.Background(b.AnimateColor("bg", fg.Alpha(0.12*unit(b.Hovered())), ui.Snappy)).TextColor(it.ink)
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must(it.glyph)).Size(14, 14).Shrink(0)
					ui.Text(c, it.name).FontSize(t.Rem(0.8)).Shrink(0)
				})
			})
			if b.Clicked() {
				*action = it.action
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("SelectBar", "action", "action", int(it.action), "count", count)
			}
		}
	})
	return root
}

// inverse is the dark bar palette, falling back per field.
func inverse(t *ui.Theme) (bg, fg ui.Color) {
	bg, fg = t.Inverse, t.InverseText
	if bg.A == 0 {
		bg = t.Text
	}
	if fg.A == 0 {
		fg = t.Background
	}
	return bg, fg
}
