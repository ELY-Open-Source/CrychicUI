package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// SyncState is synced, syncing or failed.
type SyncState uint8

// The three sync states.
const (
	Synced SyncState = iota
	Syncing
	SyncFailed
)

// SyncStatusOptions tune a SyncStatus; empty copy reads English.
type SyncStatusOptions struct {
	Synced, Syncing, Failed string
}

// SyncStatus is a sync icon and text; syncing spins.
func SyncStatus(c *ui.Context, state SyncState, o SyncStatusOptions) ui.Element {
	if o.Synced == "" {
		o.Synced = "Synced"
	}
	if o.Syncing == "" {
		o.Syncing = "Syncing"
	}
	if o.Failed == "" {
		o.Failed = "Sync failed"
	}
	t := c.Theme()
	var text, glyph string
	var tint ui.Color
	switch state {
	case Synced:
		text, glyph, tint = o.Synced, "cloud-check", t.Success
	case Syncing:
		text, glyph, tint = o.Syncing, "refresh-cw", t.Accent
	case SyncFailed:
		text, glyph, tint = o.Failed, "cloud-alert", t.Danger
	default:
		panic("feedback: unknown SyncState")
	}
	root := ui.Row(c).Gap(t.Space(1.5)).AlignSelf(ui.Start).Role(ui.RoleStatus).Label(text)
	return root.Children(func() {
		i := ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
		i.TextColor(i.AnimateColor("tint", tint, quick))
		if state == Syncing {
			i.Rotate(i.Loop("spin", time.Second, ui.Linear) * 360)
		}
		ui.Text(c, text).FontSize(t.Rem(0.85)).TextColor(t.TextMuted).ContentTransition(ui.ContentTransition{})
	})
}
