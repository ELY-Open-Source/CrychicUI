package overlay

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShareTarget is a label and Lucide icon.
type ShareTarget struct{ Label, Icon string }

// ShareSheetOptions tune a ShareSheet.
type ShareSheetOptions struct {
	Copy, Copied string
}

// ShareSheet offers targets and a link to copy.
func ShareSheet(c *ui.Context, title, link string, targets []ShareTarget, chosen *int, open *bool, o ShareSheetOptions) ui.Element {
	if title == "" || link == "" || len(targets) == 0 || chosen == nil || open == nil || *chosen < -1 || *chosen >= len(targets) {
		panic(fmt.Sprintf("overlay: ShareSheet %q needs a link, targets, a choice or -1, and open", title))
	}
	for _, s := range targets {
		if s.Label == "" {
			panic(fmt.Sprintf("overlay: ShareSheet %q targets need labels", title))
		}
		icon.Must(s.Icon)
	}
	def(&o.Copy, "Copy link")
	def(&o.Copied, "Copied")
	t := c.Theme()
	root := ui.Box(c)
	copied := ui.Local(root, "copied", func() time.Time { return time.Time{} })
	now := c.Now()
	fresh := !copied.IsZero() && now.Sub(*copied) < 1500*time.Millisecond
	if fresh {
		c.After(copied.Add(1500 * time.Millisecond).Sub(now))
	}
	if side(c, open, title, bottom(c), ui.Pose{Y: t.Space(60)}, func() {
		ui.Text(c, title).Bold().FontSize(t.Rem(1.1))
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).FocusGroup(ui.Horizontal).Children(func() {
			for i, s := range targets {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(s.Label).Width(t.Space(18)).Padding(t.Space(1)).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
					ui.Column(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(t.Space(12), t.Space(12)).Center().Radius(999).Background(t.Accent.Alpha(0.12)).Children(func() {
							ui.Icon(c, icon.Must(s.Icon)).Size(16, 16).TextColor(t.Accent)
						})
						ui.Text(c, s.Label).FontSize(t.Rem(0.8)).SingleLine()
					})
				})
				if b.Clicked() {
					*chosen, *open = i, false
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("ShareSheet", "shared", "index", i)
				}
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).
			Radius(t.Radius).Border(1, t.Border).Children(func() {
			kit.Muted(c, link).Grow(1).SingleLine()
			word := o.Copy
			if fresh {
				word = o.Copied
			}
			b := kit.Press(ui.ButtonBase(c).Label(o.Copy).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius))
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
				ui.Text(c, word).SingleLine().ContentTransition(ui.ContentTransition{})
			})
			if b.Clicked() {
				c.WriteClipboard(link)
				*copied = now
				c.After(1500 * time.Millisecond)
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("ShareSheet", "copied")
			}
		})
	}) {
		root.SetChanged()
	}
	moved(root, *open, "ShareSheet")
	return root
}
