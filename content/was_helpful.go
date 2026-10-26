package content

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Helpful is a page vote.
type Helpful int

// Votes.
const (
	Unrated Helpful = iota
	Helped
	NotHelped
)

// WasHelpfulOptions word the question and answers.
type WasHelpfulOptions struct{ Question, Yes, No, Thanks string }

// WasHelpful asks yes or no about a page; cue Pop.
func WasHelpful(c *ui.Context, vote *Helpful, o WasHelpfulOptions) ui.Element {
	def(&o.Question, "Was this page helpful?")
	def(&o.Yes, "Yes")
	def(&o.No, "No")
	def(&o.Thanks, "Thanks for the feedback.")
	if vote == nil || *vote < Unrated || *vote > NotHelped {
		panic(fmt.Sprintf("content: WasHelpful %q needs a known vote", o.Question))
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(2)).Padding(t.Space(4)).Role(ui.RoleGroup)
	return root.Children(func() {
		line := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center)
		seen := line.LayoutBounds().W
		kit.Settle(c, line, []float32{seen}, func() []float32 { return []float32{line.LayoutBounds().W} })
		// Labels get the line less pill padding and icon.
		room := t.Space(30)
		if seen > 0 {
			room = max(t.Space(6), seen-t.Space(10))
		}
		line.Children(func() {
			ui.Text(c, o.Question).FontSize(t.Rem(0.88)).Bold().Grow(1).Basis(t.Space(40)).MinWidth(0).Wrap()
			c.Key("answers")
			ui.Row(c).Wrap().Gap(t.Space(2)).MinWidth(0).MaxWidth(max(seen, t.Space(30))).Children(func() {
				for _, a := range [2]struct {
					v          Helpful
					text, mark string
				}{{Helped, o.Yes, "thumbs-up"}, {NotHelped, o.No, "thumbs-down"}} {
					on := *vote == a.v
					ink := t.Success
					if a.v == NotHelped {
						ink = t.Danger
					}
					b := kit.Press(ui.ButtonBase(c).Label(a.text).Checked(on).Padding(t.Space(1.5), t.Space(3)).Radius(999))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.14), unit(on)), ui.Snappy))
					b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(on)), ui.Snappy))
					b.Children(func() {
						ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							c.Key(on)
							mk := glyph(c, a.mark, 14).TextColor(t.TextMuted.Mix(ink, unit(on)))
							if on {
								kit.Arrive(mk, ui.Pose{Scale: 0.5, Rotate: -20, Fade: 1}, 0)
							}
							ui.Text(c, a.text).FontSize(t.Rem(0.8)).Bold().TextColor(t.Text.Mix(ink, unit(on))).Shrink(0).MaxWidth(room).SingleLine()
						})
					})
					if b.Clicked() && !on {
						*vote = a.v
						root.SetChanged()
						sound.Play(sound.Pop)
						kit.Log("WasHelpful", "vote", "value", int(a.v))
					}
				}
			})
		})
		if *vote != Unrated {
			c.Key("thanks")
			kit.Arrive(kit.Muted(c, o.Thanks).FontSize(t.Rem(0.78)).FillWidth().Wrap(), ui.Pose{Y: t.Space(1), Fade: 1}, 0)
		}
	})
}
