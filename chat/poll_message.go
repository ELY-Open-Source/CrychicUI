package chat

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PollMessageOptions tune a PollMessage.
type PollMessageOptions struct {
	// Total words the vote count, "%d" taking it.
	Total string
}

// PollMessage shows options with vote bars; bars grow.
func PollMessage(c *ui.Context, question string, options []string, votes *[]int, chosen *int, o PollMessageOptions) ui.Element {
	if question == "" || len(options) == 0 || votes == nil || len(*votes) != len(options) || chosen == nil || *chosen < -1 || *chosen >= len(options) {
		panic(fmt.Sprintf("chat: PollMessage %q needs a question, options, votes each and a choice or -1", question))
	}
	if *chosen >= 0 && (*votes)[*chosen] == 0 {
		panic(fmt.Sprintf("chat: PollMessage %q votes must cover the choice", question))
	}
	total := 0
	for _, v := range *votes {
		if v < 0 {
			panic(fmt.Sprintf("chat: PollMessage %q votes must not be negative", question))
		}
		total += v
	}
	def(&o.Total, "%d votes")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(76)).Gap(t.Space(1.5)).Padding(t.Space(3)).Radius(t.Radius*2).
		Background(t.Surface).Border(1, t.Border).FocusGroup(ui.Vertical).Role(ui.RoleGroup).Label(question)
	g := root.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 500 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	root.Children(func() {
		ui.Text(c, question).Bold().Wrap()
		for i, opt := range options {
			if opt == "" {
				panic(fmt.Sprintf("chat: PollMessage %q options must not be empty", question))
			}
			c.Key(i)
			frac := float32(0)
			if total > 0 {
				frac = float32((*votes)[i]) / float32(total)
			}
			b := kit.Press(ui.ButtonBase(c).Label(opt).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(i == *chosen))
			bar := b.AnimateTo(fmt.Sprint("bar.", i), frac, ui.Smooth) * g
			bg := t.Background.Mix(t.SurfaceHover, unit(b.Hovered()))
			if i == *chosen {
				bg = t.Accent.Alpha(0.1)
			}
			b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Border(1, t.Border)
			b.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
				if bar > 0 {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * bar, H: r.H}, t.Accent.Alpha(0.18), t.Radius)
				}
			}))
			b.Children(func() {
				if i == *chosen {
					ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Accent)
				}
				ui.Text(c, opt).Grow(1).SingleLine()
				pct := 0
				if total > 0 {
					pct = int(float64((*votes)[i])*100/float64(total) + 0.5)
				}
				kit.Numeric(c, strconv.Itoa(pct)+"%").FontSize(t.Rem(0.85)).TextColor(t.TextMuted)
			})
			if b.Clicked() && i != *chosen {
				if *chosen >= 0 {
					(*votes)[*chosen]--
				}
				(*votes)[i]++
				*chosen = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("PollMessage", "voted", "index", i)
			}
		}
		kit.Muted(c, fmt.Sprintf(o.Total, total)).FontSize(t.Rem(0.8)).ContentTransition(ui.ContentTransition{})
	})
	return root
}
