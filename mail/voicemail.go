package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VoiceMessage is one voicemail and transcript.
type VoiceMessage struct {
	ID, Sender, Time, Transcript string
}

// VoicemailOptions are the voicemail labels.
type VoicemailOptions struct{ Play, Pause string }

// Voicemail plays messages and reveals transcripts; cue Toggle.
func Voicemail(c *ui.Context, messages []VoiceMessage, playing *int, o VoicemailOptions) ui.Element {
	if len(messages) == 0 || playing == nil || *playing < -1 || *playing >= len(messages) {
		panic("mail: Voicemail needs messages and valid state")
	}
	def(&o.Play, "Play")
	def(&o.Pause, "Pause")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, message := range messages {
			if message.ID == "" || message.Sender == "" || message.Time == "" || message.Transcript == "" {
				panic("mail: invalid voicemail")
			}
			c.Key(message.ID)
			card := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
			card.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, message.Sender).Grow(1)
					kit.Muted(c, message.Time)
					label := o.Play
					if *playing == i {
						label = o.Pause
					}
					b := button(c, label).Checked(*playing == i)
					b.Children(func() { ui.Text(c, label) })
					if b.Clicked() {
						if *playing == i {
							*playing = -1
						} else {
							*playing = i
						}
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("Voicemail", "toggled", "index", i, "playing", *playing == i)
					}
				})
				if *playing == i {
					ui.Text(c, message.Transcript).Wrap().Transition(kit.Enter)
				}
			})
		}
	})
	return root
}
