package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RateAnswerOptions name the actions.
type RateAnswerOptions struct {
	Like, Dislike, Send string
}

// RateAnswer takes a thumb, then a reason, then sends.
func RateAnswer(c *ui.Context, rating *int, reasons []string, reason *int, o RateAnswerOptions) (root, send ui.Element) {
	if rating == nil || *rating < -1 || *rating > 1 || reason == nil || *reason < -1 || *reason >= len(reasons) {
		panic("agent: RateAnswer needs a rating from -1 to 1 and a reason or -1")
	}
	for _, r := range reasons {
		if r == "" {
			panic("agent: RateAnswer reasons must not be empty")
		}
	}
	def(&o.Like, "Helpful")
	def(&o.Dislike, "Not helpful")
	def(&o.Send, "Send feedback")
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).Children(func() {
			up := round(c, o.Like, "thumbs-up").Checked(*rating == 1)
			down := round(c, o.Dislike, "thumbs-down").Checked(*rating == -1)
			upInk, downInk := t.TextMuted, t.TextMuted
			if *rating == 1 {
				upInk = t.Accent
			}
			if *rating == -1 {
				downInk = t.Accent
			}
			up.TextColor(up.AnimateColor("ink", upInk, ui.Snappy))
			down.TextColor(down.AnimateColor("ink", downInk, ui.Snappy))
			if up.Clicked() {
				vote(rating, 1)
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("RateAnswer", "rating", *rating)
			}
			if down.Clicked() {
				vote(rating, -1)
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("RateAnswer", "rating", *rating)
			}
		})
		if *rating != 0 && len(reasons) > 0 {
			chips := ui.Row(c).Wrap().Gap(t.Space(1.5)).FocusGroup(ui.Horizontal)
			kit.Arrive(chips, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
				for i, r := range reasons {
					c.Key(i)
					on := i == *reason
					b := kit.Press(ui.ButtonBase(c).Label(r).Padding(t.Space(1), t.Space(2)).Radius(999).Border(1, t.Border).Checked(on))
					bg := t.Surface.Mix(t.SurfaceHover, unit(b.Hovered()))
					if on {
						bg = t.Accent.Alpha(0.12)
					}
					b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
						ui.Text(c, r).FontSize(t.Rem(0.85)).SingleLine()
					})
					if b.Clicked() && !on {
						*reason = i
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("RateAnswer", "reason", "index", i)
					}
				}
			})
		}
		if *rating != 0 {
			ready := len(reasons) == 0 || *reason >= 0
			send = kit.Press(ui.ButtonBase(c).Label(o.Send).AlignSelf(ui.Start).Padding(t.Space(1.5), t.Space(4)).Radius(t.Radius)).Disabled(!ready)
			send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy)).
				TextColor(t.AccentText).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must("send")).Size(16, 16)
					ui.Text(c, o.Send).Bold().SingleLine()
				})
			})
			kit.Arrive(send, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
		}
	})
	if *rating != 0 && send.Clicked() && (len(reasons) == 0 || *reason >= 0) {
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("RateAnswer", "sent", "rating", *rating)
	}
	return root, send
}

// vote sets or clears a rating.
func vote(rating *int, v int) {
	if *rating == v {
		*rating = 0
	} else {
		*rating = v
	}
}
