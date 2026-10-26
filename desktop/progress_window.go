package desktop

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProgressWindowOptions word the window; Count formats done of total.
type ProgressWindowOptions struct{ Detail, Count, Cancel, Done, Icon string }

// ProgressWindow shows a task's progress with a cancel; cue Tap.
func ProgressWindow(c *ui.Context, title string, done, total int, o ProgressWindowOptions) (root, cancel ui.Element) {
	def(&o.Count, "%d of %d items")
	def(&o.Cancel, "Cancel")
	def(&o.Done, "Done")
	def(&o.Icon, "download")
	if strings.TrimSpace(title) == "" || total <= 0 || done < 0 || done > total {
		panic(fmt.Sprintf("desktop: ProgressWindow %q needs 0 <= done <= total and total > 0, got %d of %d", title, done, total))
	}
	t := c.Theme()
	finished := done == total
	root = window(c).Role(ui.RoleGroup)
	root.Children(func() {
		titleBar(c, title, true, nil)
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start).Padding(t.Space(4)).Children(func() {
			c.Key(finished)
			mark := o.Icon
			if finished {
				mark = "check"
			}
			ink := t.Accent.Mix(t.Success, unit(finished))
			kit.Arrive(tile(c, mark, ink, t.Space(10)), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(2)).Children(func() {
				ui.Text(c, title).FontSize(t.Rem(0.86)).Bold().FillWidth().MinWidth(0).SingleLine()
				if o.Detail != "" {
					kit.Muted(c, o.Detail).FontSize(t.Rem(0.74)).FillWidth().MinWidth(0).SingleLine()
				}
				c.Key("bar")
				bar := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).Background(t.SurfaceHover)
				fill := bar.AnimateTo("fill", float32(done)/float32(total), ui.Smooth)
				paint := bar.AnimateColor("ink", t.Accent.Mix(t.Success, unit(finished)), ui.Snappy)
				bar.Draw(func(p *ui.Painter, r ui.Rect) {
					if fill > 0 {
						p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, paint, r.H/2)
					}
				})
				c.Key("actions")
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					kit.Numeric(c, fmt.Sprintf(o.Count, done, total)).FontSize(t.Rem(0.74)).TextColor(t.TextMuted).Grow(1).MinWidth(0).SingleLine()
					label := o.Cancel
					if finished {
						label = o.Done
					}
					c.Key("cancel")
					cancel = kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1), t.Space(3)).Radius(999).Border(1, t.Border))
					cancel.Background(cancel.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(cancel.Hovered())), ui.Snappy))
					cancel.Children(func() {
						ui.Text(c, label).FontSize(t.Rem(0.76)).Bold().Shrink(0).MaxWidth(t.Space(30)).SingleLine().ContentTransition(ui.ContentTransition{})
					})
					if cancel.Clicked() {
						sound.Play(sound.Tap)
						kit.Log("ProgressWindow", "closed", "finished", finished)
					}
				})
			})
		})
	})
	return root, cancel
}
