package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// OfflineStateOptions tune an OfflineState; empty copy reads English.
type OfflineStateOptions struct {
	Title, Text, Retry string
}

// OfflineState is a pulsing offline icon, message and Retry.
func OfflineState(c *ui.Context, o OfflineStateOptions) (root, retry ui.Element) {
	if o.Title == "" {
		o.Title = "You are offline"
	}
	if o.Text == "" {
		o.Text = "Check your connection, then try again."
	}
	if o.Retry == "" {
		o.Retry = "Retry"
	}
	root, retry = state(c, icon.Must("wifi-off"), c.Theme().TextMuted, o.Title, o.Text, o.Retry, true, func(i ui.Element) {
		k := i.Loop("pulse", 1600*time.Millisecond, ui.Bounce(ui.EaseInOut))
		i.Opacity(0.45 + 0.55*k)
	})
	if kit.Clicked(retry, sound.Tap) {
		kit.Log("OfflineState", "retry")
	}
	return root, retry
}
