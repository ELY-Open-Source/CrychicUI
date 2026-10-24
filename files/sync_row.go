package files

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SyncState is a file's synchronization state.
type SyncState uint8

// Sync states.
const (
	Synced SyncState = iota
	Syncing
	SyncFailed
)

// SyncRowOptions name the sync states.
type SyncRowOptions struct{ Done, Working, Failed string }

// SyncRow shows a file's sync state, spinning while syncing.
func SyncRow(c *ui.Context, item FileItem, state SyncState, o SyncRowOptions) ui.Element {
	validFile("SyncRow", item)
	if state > SyncFailed {
		panic(fmt.Sprintf("files: SyncRow %q needs a known state", item.Name))
	}
	def(&o.Done, "Synced")
	def(&o.Working, "Syncing")
	def(&o.Failed, "Sync failed")
	t := c.Theme()
	name, glyph, ink := o.Done, "cloud-check", t.Success
	switch state {
	case Syncing:
		name, glyph, ink = o.Working, "refresh-cw", t.Accent
	case SyncFailed:
		name, glyph, ink = o.Failed, "cloud-alert", t.Danger
	}
	root := kit.Card(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(2.5)).Role(ui.RoleGroup).Label(item.Name)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Row().AlignItems(ui.Center).Children(func() {
		fileTile(c, item, t.Space(8))
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			ui.Text(c, item.Name).FontSize(t.Rem(0.85)).FillWidth().SingleLine()
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				sizeText(c, item.Size).FontSize(t.Rem(0.75))
				if item.Path != "" {
					kit.Muted(c, "·").FontSize(t.Rem(0.75))
					kit.Muted(c, item.Path).FontSize(t.Rem(0.75)).Shrink(1).MinWidth(0).SingleLine()
				}
			})
		})
		pill := ui.Row(c).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2.5)).Radius(999)
		pill.Background(pill.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).TextColor(pill.AnimateColor("ink", ink, ui.Snappy))
		pill.Children(func() {
			mark := ui.Icon(c, icon.Must(glyph)).Size(14, 14).ContentTransition(ui.ContentTransition{})
			if state == Syncing {
				mark.Rotate(360 * mark.Loop("spin", time.Second, ui.Linear))
			}
			ui.Text(c, name).FontSize(t.Rem(0.75)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
		})
	})
	return root
}
