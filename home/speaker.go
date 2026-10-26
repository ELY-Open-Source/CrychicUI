package home

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Track is what plays; nil Art shows a note.
type Track struct {
	Title, Artist string
	Art           *ui.Bitmap
}

// SpeakerOptions word the speaker controls.
type SpeakerOptions struct{ Room, Play, Pause, Volume string }

// Speaker plays a track and sets volume; cue Tick.
func Speaker(c *ui.Context, name string, tr Track, playing *bool, volume *float64, o SpeakerOptions) ui.Element {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(tr.Title) == "" || playing == nil || volume == nil || !(*volume >= 0 && *volume <= 1) {
		panic(fmt.Sprintf("home: Speaker %q needs a track title, a playing state and a volume of 0 to 1", name))
	}
	def(&o.Play, "Play")
	def(&o.Pause, "Pause")
	def(&o.Volume, "Volume")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
			glyph(c, "speaker", 14)
			ui.Text(c, name).FontSize(t.Rem(0.72)).Bold().Grow(1).MinWidth(0).SingleLine()
			if o.Room != "" {
				ui.Text(c, o.Room).FontSize(t.Rem(0.72)).MinWidth(0).SingleLine()
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			art := ui.Box(c).Size(t.Space(14), t.Space(14)).Shrink(0).Center().Radius(t.Radius).Clip().Background(t.Accent.Alpha(0.12)).
				TextColor(t.Accent)
			art.Children(func() {
				if tr.Art != nil {
					ui.Image(c, tr.Art).Fill().Fit(ui.Cover)
				} else {
					glyph(c, "music", 22)
				}
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, tr.Title).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).SingleLine()
				if tr.Artist != "" {
					kit.Muted(c, tr.Artist).FontSize(t.Rem(0.76)).FillWidth().MinWidth(0).SingleLine()
				}
			})
			eq := ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0)
			var beat float32
			if *playing {
				beat = eq.Loop("eq", 900*time.Millisecond, ui.Linear)
			}
			on := eq.AnimateTo("on", unit(*playing), ui.Smooth)
			eq.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				w := r.W / 5
				for i := 0; i < 3; i++ {
					h := 0.25 + 0.75*on*float32(math.Abs(math.Sin(math.Pi*(float64(beat)+float64(i)/3))))
					p.Fill(ui.Rect{X: r.X + float32(i)*2*w, Y: r.Y + r.H*(1-h), W: w, H: r.H * h}, t.Accent, w/2)
				}
			}))
			label, mark := o.Play, "play"
			if *playing {
				label, mark = o.Pause, "pause"
			}
			c.Key("play")
			b := round(c, label, mark, t.Space(11), t.Accent, t.AccentText)
			if b.Clicked() {
				*playing = !*playing
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Speaker", "playing", "on", *playing)
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			mark := "volume-2"
			switch {
			case *volume == 0:
				mark = "volume-x"
			case *volume < 0.5:
				mark = "volume-1"
			}
			c.Key(mark)
			glyph(c, mark, 16).TextColor(t.TextMuted)
			c.Key("volume")
			slider := ui.Slider(c, volume, 0, 1).Step(0.05).Label(o.Volume).Grow(1).MinWidth(0)
			if slider.Changed() {
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("Speaker", "volume", "percent", int(math.Round(*volume*100)))
			}
		})
	})
	return root
}
