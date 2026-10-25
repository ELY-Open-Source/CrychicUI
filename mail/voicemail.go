package mail

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VoiceMessage is one voicemail and transcript; Length is optional.
type VoiceMessage struct {
	ID, Sender, Time, Transcript string
	Length                       time.Duration
}

// VoicemailOptions word the card and play buttons.
type VoicemailOptions struct{ Title, Play, Pause string }

// Voicemail plays messages and reveals transcripts; cue Toggle.
func Voicemail(c *ui.Context, messages []VoiceMessage, playing *int, o VoicemailOptions) ui.Element {
	if len(messages) == 0 || playing == nil || *playing < -1 || *playing >= len(messages) {
		panic("mail: Voicemail needs messages and valid state")
	}
	ids := make([]string, len(messages))
	for i, m := range messages {
		if m.Sender == "" || m.Time == "" || m.Transcript == "" || m.Length < 0 {
			panic(fmt.Sprintf("mail: Voicemail message %q needs sender, time, and transcript", m.ID))
		}
		ids[i] = m.ID
	}
	unique("Voicemail", ids)
	def(&o.Title, "Voicemail")
	def(&o.Play, "Play")
	def(&o.Pause, "Pause")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "voicemail", o.Title, func() { kit.Chip(c, fmt.Sprint(len(messages)), t.TextMuted) })
		list := ui.Column(c).FillWidth()
		list.Children(func() {
			for i, m := range messages {
				c.Key(m.ID)
				on := *playing == i
				item := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2.5), t.Space(3))
				if i < len(messages)-1 {
					item.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
				}
				kit.Stagger(c, list, item, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						label, mark := o.Play, "play"
						if on {
							label, mark = o.Pause, "pause"
						}
						b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(9), t.Space(9)).Shrink(0).Center().Radius(999).
							TextColor(t.AccentText).Checked(on))
						b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy))
						b.Children(func() {
							c.Key(mark)
							kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.5, Fade: 1}, 0)
						})
						if b.Clicked() {
							if on {
								*playing = -1
							} else {
								*playing = i
							}
							root.SetChanged()
							sound.Play(sound.Toggle)
							kit.Log("Voicemail", "toggled", "index", i, "playing", *playing == i)
						}
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
							ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
								ui.Text(c, m.Sender).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
								kit.Muted(c, m.Time).FontSize(t.Rem(0.7)).Shrink(0)
							})
							wave(c, m, on)
						})
						if m.Length > 0 {
							kit.Muted(c, clock(m.Length)).FontSize(t.Rem(0.72)).Shrink(0)
						}
					})
					layout.Collapse(c, on, func() {
						ui.Text(c, m.Transcript).FontSize(t.Rem(0.8)).Wrap().FillWidth().Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).
							Background(t.SurfaceHover.Alpha(0.6))
					})
				})
			}
		})
	})
	return root
}

// wave draws a message's bars; playing sweeps them in accent.
func wave(c *ui.Context, m VoiceMessage, on bool) {
	t := c.Theme()
	bars := ui.Box(c).FillWidth().Height(t.Space(6))
	length := m.Length
	if length <= 0 {
		length = 8 * time.Second
	}
	head := float32(0)
	if on {
		head = bars.Loop("play", length, ui.Linear)
	}
	h := fnv.New64a()
	h.Write([]byte(m.ID))
	seed := h.Sum64()
	bars.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		const n = 32
		w := r.W / n
		s := seed
		for k := range n {
			s = s*6364136223846793005 + 1442695040888963407
			level := 0.25 + 0.75*float32(s>>40)/float32(1<<24)
			bh := r.H * level
			ink := t.Border
			if on && float32(k)/n < head {
				ink = t.Accent
			}
			p.Fill(ui.Rect{X: r.X + float32(k)*w + w*0.2, Y: r.Y + (r.H-bh)/2, W: w * 0.6, H: bh}, ink, w*0.3)
		}
	}))
}
