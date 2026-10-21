package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PlayButtonOptions tune a PlayButton.
type PlayButtonOptions struct {
	// Label names the action.
	Label string
}

// PlayButton crosses play and pause icons when toggled.
func PlayButton(c *ui.Context, playing *bool, o PlayButtonOptions) ui.Element {
	if playing == nil {
		panic("button: PlayButton needs state")
	}
	if o.Label == "" {
		o.Label = "Play"
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Tooltip(o.Label).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Center().
		Radius(999).Background(t.Surface).Border(1, t.Border).Checked(*playing))
	root.Children(func() {
		name := "play"
		if *playing {
			name = "pause"
		}
		ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
	})
	if root.Clicked() {
		*playing = !*playing
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("PlayButton", "toggled", "playing", *playing)
	}
	return root
}
