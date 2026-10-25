package social

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PollChoice is one poll answer and vote count.
type PollChoice struct {
	Text  string
	Votes int
}

// PollPostOptions word the poll; Vote, Votes take the total.
type PollPostOptions struct {
	Left, Closed, Vote, Votes string
	Remaining                 func(time.Duration) string
}

// PollPost collects votes with time left; cue Select.
func PollPost(c *ui.Context, question string, choices *[]PollChoice, ends time.Time, voted *int, o PollPostOptions) ui.Element {
	if question == "" || choices == nil || len(*choices) < 2 || ends.IsZero() || voted == nil || *voted < -1 || *voted >= len(*choices) {
		panic(fmt.Sprintf("social: PollPost %q needs choices, an end and a vote", question))
	}
	def(&o.Left, "%s left")
	def(&o.Closed, "Final results")
	def(&o.Vote, "%s vote")
	def(&o.Votes, "%s votes")
	if o.Remaining == nil {
		o.Remaining = shortRemaining
	}
	total := uint64(0)
	texts := make([]string, len(*choices))
	for i, choice := range *choices {
		votes := uint64(choice.Votes)
		if choice.Text == "" || choice.Votes < 0 || votes > ^uint64(0)-total {
			panic(fmt.Sprintf("social: PollPost %q has an invalid choice %q", question, choice.Text))
		}
		total += votes
		texts[i] = choice.Text
	}
	unique("PollPost", texts)
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	remaining := max(0, ends.Sub(c.Now()))
	if remaining > 0 {
		c.After(time.Second)
	}
	atCapacity := total == ^uint64(0)
	results := *voted >= 0 || remaining == 0
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
				glyph(c, "chart-bar", 16).TextColor(t.Accent)
				ui.Text(c, question).FontSize(t.Rem(0.92)).Bold().Grow(1).MinWidth(0).Wrap()
			})
			list := ui.Column(c).FillWidth().Gap(t.Space(1.5))
			list.Children(func() {
				for i, choice := range *choices {
					c.Key(choice.Text)
					share := float32(0)
					if total > 0 {
						share = float32(float64(choice.Votes) / float64(total))
					}
					mine := *voted == i
					b := kit.Press(ui.ButtonBase(c).Label(choice.Text).FillWidth().AlignItems(ui.Stretch).Radius(t.Radius).Clip().Checked(mine).
						Disabled(results || atCapacity || choice.Votes == math.MaxInt))
					b.Border(1, b.AnimateColor("line", t.Accent.Mix(t.Border, unit(results)), ui.Snappy))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.08), unit(b.Hovered() && !results)), ui.Snappy))
					fill := b.AnimateTo("votes", share*unit(results), ui.Smooth)
					ink := t.Accent.Alpha(0.22)
					if !mine {
						ink = t.SurfaceHover
					}
					b.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
						p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, ink, 0)
					}))
					kit.Stagger(c, list, b, i).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
							label := ui.Text(c, choice.Text).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine().
								TextColor(t.Accent.Mix(t.Text, unit(results)))
							if mine || !results {
								label.Bold()
							}
							c.Key(mine)
							if mine {
								kit.Arrive(glyph(c, "circle-check", 16).TextColor(t.Accent), ui.Pose{Scale: 0.4, Fade: 1}, 0)
							}
							if results {
								kit.Numeric(c, fmt.Sprintf("%.0f%%", share*100)).FontSize(t.Rem(0.8)).Bold().Shrink(0)
							}
						})
					})
					if b.Clicked() && !results && !atCapacity && choice.Votes < math.MaxInt {
						(*choices)[i].Votes++
						*voted = i
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("PollPost", "voted", "index", i, "count", (*choices)[i].Votes)
					}
				}
			})
			votes := o.Votes
			if total == 1 {
				votes = o.Vote
			}
			left := o.Closed
			if remaining > 0 {
				left = fmt.Sprintf(o.Left, o.Remaining(remaining))
			}
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
				kit.Numeric(c, fmt.Sprintf(votes, number(c, total))).FontSize(t.Rem(0.75)).Shrink(0)
				ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.TextMuted)
				glyph(c, "clock", 12)
				ui.Text(c, left).FontSize(t.Rem(0.75)).Shrink(0)
			})
		})
	})
	return root
}

// shortRemaining is a duration in its largest whole unit.
func shortRemaining(d time.Duration) string {
	switch {
	case d >= 24*time.Hour:
		return fmt.Sprintf("%dd", int64(d/(24*time.Hour)))
	case d >= time.Hour:
		return fmt.Sprintf("%dh", int64(d/time.Hour))
	case d >= time.Minute:
		return fmt.Sprintf("%dm", int64(d/time.Minute))
	}
	return fmt.Sprintf("%ds", int64((d+time.Second-1)/time.Second))
}
