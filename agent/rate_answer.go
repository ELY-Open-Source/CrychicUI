package agent

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RateAnswerOptions name the question, actions and thanks.
type RateAnswerOptions struct {
	Question, Like, Dislike, Send, Thanks string
}

// RateAnswer takes a thumb, then a reason, then sends.
func RateAnswer(c *ui.Context, rating *int, reasons []string, reason *int, sent *bool, o RateAnswerOptions) (root, send ui.Element) {
	if rating == nil || *rating < -1 || *rating > 1 || reason == nil || *reason < -1 || *reason >= len(reasons) || sent == nil {
		panic("agent: RateAnswer needs a rating from -1 to 1, a reason or -1 and sent")
	}
	seen := map[string]bool{}
	for _, r := range reasons {
		if r == "" || seen[r] {
			panic("agent: RateAnswer reasons must be unique and not empty")
		}
		seen[r] = true
	}
	def(&o.Question, "Was this answer helpful?")
	def(&o.Like, "Helpful")
	def(&o.Dislike, "Not helpful")
	def(&o.Send, "Send feedback")
	def(&o.Thanks, "Thanks for the feedback")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Role(ui.RoleGroup).Label(o.Question)
	root.Children(func() {
		if *sent {
			done := ui.Row(c).Key("sent").Gap(t.Space(2)).AlignItems(ui.Center).TextColor(t.Success)
			kit.Arrive(done, ui.Pose{Scale: 0.9, Fade: 1}, 0).Children(func() {
				ui.Icon(c, icon.Must("circle-check")).Size(16, 16)
				ui.Text(c, o.Thanks).Bold()
			})
			return
		}
		ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Question).FontSize(t.Rem(0.9)).Grow(1).MinWidth(0).Wrap()
			up := round(c, o.Like, "thumbs-up").Checked(*rating == 1)
			down := round(c, o.Dislike, "thumbs-down").Checked(*rating == -1)
			up.TextColor(up.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(*rating == 1)), ui.Snappy))
			down.TextColor(down.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(*rating == -1)), ui.Snappy))
			v := 0
			switch {
			case up.Clicked():
				v = 1
			case down.Clicked():
				v = -1
			}
			if v != 0 {
				vote(rating, v)
				*reason = -1
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("RateAnswer", "rating", *rating)
			}
		})
		if *rating == 0 {
			return
		}
		if len(reasons) > 0 {
			chips := ui.Row(c).Key(*rating).FillWidth().Wrap().Gap(t.Space(1.5)).FocusGroup(ui.Horizontal)
			kit.Arrive(chips, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
				for i, r := range reasons {
					c.Key(r)
					on := i == *reason
					b := kit.Press(ui.ButtonBase(c).Label(r).Padding(t.Space(1), t.Space(2.5)).Radius(999).Checked(on))
					b.Border(1, b.AnimateColor("ring", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.12), unit(on)).Mix(t.SurfaceHover, 0.5*unit(b.Hovered() && !on)), ui.Snappy))
					b.Children(func() { ui.Text(c, r).FontSize(t.Rem(0.85)).Shrink(0).MaxWidth(t.Space(60)).Wrap() })
					if b.Clicked() && !on {
						*reason = i
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("RateAnswer", "reason", "index", i)
					}
				}
			})
		}
		ready := len(reasons) == 0 || *reason >= 0
		send = kit.Press(ui.ButtonBase(c).Label(o.Send).AlignSelf(ui.End).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)).Disabled(!ready)
		send.Background(send.AnimateColor("bg", t.SurfaceHover.Mix(t.Accent.Mix(t.AccentHover, unit(send.Hovered())), unit(ready)), ui.Snappy)).
			TextColor(send.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(ready)), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("send")).Size(16, 16)
			ui.Text(c, o.Send).Bold().Shrink(0)
		})
		kit.Arrive(send, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	})
	if !*sent && *rating != 0 && send.Clicked() && (len(reasons) == 0 || *reason >= 0) {
		*sent = true
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
