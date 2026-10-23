package files

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SyncState is a file's synchronization state.
type SyncState uint8

const (
	Synced SyncState = iota
	Syncing
	SyncFailed
)

// SyncRowOptions are the sync state labels.
type SyncRowOptions struct{ Done, Working, Failed string }

// SyncRow shows a file and its sync state.
func SyncRow(c *ui.Context, item FileItem, state SyncState, o SyncRowOptions) ui.Element {
	validFile("SyncRow", item)
	if state > SyncFailed {
		panic("files: SyncRow needs a valid state")
	}
	def(&o.Done, "Synced")
	def(&o.Working, "Syncing")
	def(&o.Failed, "Sync failed")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).
		Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Icon(c, icon.Must(itemGlyph(item))).Size(16, 16).TextColor(t.TextMuted)
		ui.Text(c, item.Name).SingleLine().Grow(1)
		switch state {
		case Syncing:
			spin := ui.Icon(c, icon.Must("rotate-cw")).Size(16, 16).TextColor(t.Accent)
			spin.Rotate(360 * float32(spin.Loop("spin", time.Second, ui.Linear)))
			kit.Muted(c, o.Working).FontSize(t.Rem(0.8))
		case SyncFailed:
			ui.Icon(c, icon.Must("circle-x")).Size(16, 16).TextColor(t.Danger)
			ui.Text(c, o.Failed).FontSize(t.Rem(0.8)).TextColor(t.Danger)
		default:
			ui.Icon(c, icon.Must("circle-check")).Size(16, 16).TextColor(t.Success)
			kit.Muted(c, o.Done).FontSize(t.Rem(0.8))
		}
	})
	return root
}
