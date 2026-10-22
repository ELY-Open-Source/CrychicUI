package motion

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PullRefreshOptions tune a PullRefresh.
type PullRefreshOptions struct {
	// Pull, Release and Busy name each stage; English when empty.
	Pull, Release, Busy string
}

// PullRefresh sets refreshing when content is pulled past a threshold.
func PullRefresh(c *ui.Context, refreshing *bool, o PullRefreshOptions, content func()) ui.Element {
	if refreshing == nil || content == nil {
		panic("motion: PullRefresh needs a refreshing pointer and content")
	}
	if o.Pull == "" {
		o.Pull = "Pull to refresh"
	}
	if o.Release == "" {
		o.Release = "Release to refresh"
	}
	if o.Busy == "" {
		o.Busy = "Refreshing…"
	}
	t := c.Theme()
	reach, hold, give := t.Space(12), t.Space(10), t.Space(60)
	root := ui.Column(c).ClipY()
	pull := ui.Local(root, "pull", func() float32 { return 0 })
	return root.Children(func() {
		body := ui.Column(c).Background(t.Background)
		var y float32
		if _, dy, held := body.Dragged(); held && !*refreshing {
			*pull = max(*pull+dy, 0)
			y = body.Jump("y", ui.Rubberband(*pull, 0, 0, give))
		} else {
			if *pull > 0 && !*refreshing && ui.Rubberband(*pull, 0, 0, give) >= reach {
				*refreshing = true
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("PullRefresh", "refresh")
			}
			*pull = 0
			target := float32(0)
			if *refreshing {
				target = hold
			}
			y = body.AnimateTo("y", target, ui.Snappy)
		}
		body.Translate(0, y).Children(content)
		text, glyph, turn := o.Pull, "arrow-down", float32(0)
		switch {
		case *refreshing:
			text, glyph = o.Busy, "loader-circle"
		case y >= reach:
			text, turn = o.Release, 180
		}
		ui.Row(c).Absolute().Top(0).Start(0).End(0).Height(max(y, 0)).ClipY().Gap(t.Space(2)).Center().
			TextColor(t.TextMuted).Children(func() {
			i := ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
			if *refreshing {
				i.Transform(ui.Rotation(i.Loop("spin", 900*time.Millisecond, ui.Linear) * 360))
			} else {
				i.Transform(ui.Rotation(i.AnimateTo("turn", turn, ui.Snappy)))
			}
			ui.Text(c, text).FontSize(t.Rem(0.85))
		})
	})
}
