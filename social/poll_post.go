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

// PollPostOptions are the poll time labels.
type PollPostOptions struct{ Left, Closed string }

// PollPost collects votes with time left; cue Select.
func PollPost(c *ui.Context, question string, choices *[]PollChoice, ends time.Time, voted *int, o PollPostOptions) ui.Element {
	if question == "" || choices == nil || len(*choices) < 2 || ends.IsZero() || voted == nil || *voted < -1 || *voted >= len(*choices) {
		panic("social: PollPost needs valid question, choices, time, and state")
	}
	def(&o.Left, "%s left")
	def(&o.Closed, "Closed")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	remaining := max(0, ends.Sub(c.Now()))
	if remaining > 0 {
		c.After(time.Second)
	}
	root.Children(func() {
		ui.Text(c, question).Bold()
		total := uint64(0)
		for _, choice := range *choices {
			votes := uint64(choice.Votes)
			if choice.Text == "" || choice.Votes < 0 || votes > ^uint64(0)-total {
				panic("social: invalid poll choice")
			}
			total += votes
		}
		atCapacity := total == ^uint64(0)
		for i, choice := range *choices {
			c.Key(choice.Text)
			fraction := float32(0)
			if total > 0 {
				fraction = float32(float64(choice.Votes) / float64(total))
			}
			b := button(c, choice.Text).Column().FillWidth().Gap(t.Space(1)).AlignItems(ui.Start).
				Checked(*voted == i).Disabled(*voted >= 0 || remaining == 0 || atCapacity || choice.Votes == math.MaxInt)
			b.Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, choice.Text).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%.0f%%", fraction*100)).TextColor(t.TextMuted)
				})
				bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
				w := bar.AnimateTo("votes", fraction, ui.Smooth)
				bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, t.Accent, 999)
				}))
			})
			if b.Clicked() && *voted < 0 && remaining > 0 && !atCapacity && choice.Votes < math.MaxInt {
				(*choices)[i].Votes++
				*voted = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("PollPost", "voted", "index", i, "count", (*choices)[i].Votes)
			}
		}
		label := o.Closed
		if remaining > 0 {
			label = fmt.Sprintf(o.Left, remaining.Round(time.Second).String())
		}
		kit.Muted(c, label)
	})
	return root
}
