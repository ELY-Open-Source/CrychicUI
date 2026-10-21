package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VoiceButtonOptions tune a VoiceButton.
type VoiceButtonOptions struct {
	// Label names the action.
	Label string
}

// VoiceButton toggles recording and pulses rings while on.
func VoiceButton(c *ui.Context, recording *bool, o VoiceButtonOptions) ui.Element {
	if recording == nil {
		panic("button: VoiceButton needs state")
	}
	if o.Label == "" {
		o.Label = "Voice"
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Tooltip(o.Label).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Center().
		Radius(999).Background(t.Surface).Border(1, t.Border).Checked(*recording))
	if *recording {
		phase := root.Loop("rings", 1600*time.Millisecond, ui.Linear)
		root.DrawOver(func(p *ui.Painter, r ui.Rect) {
			for i := range 2 {
				at := float32(i) * 0.5
				k := phase + at
				if k >= 1 {
					k--
				}
				grow := 10 * k
				p.Stroke(ui.Rect{X: r.X - grow, Y: r.Y - grow, W: r.W + 2*grow, H: r.H + 2*grow},
					t.Accent.Alpha(0.35*(1-k)), 999, 1)
			}
		})
	}
	root.Children(func() {
		name := "mic"
		if *recording {
			name = "mic-off"
		}
		ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
	})
	if root.Clicked() {
		*recording = !*recording
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("VoiceButton", "toggled", "recording", *recording)
	}
	return root
}
