package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Repeat is a player's repeat mode.
type Repeat uint8

const (
	RepeatOff Repeat = iota
	RepeatAll
	RepeatOne
	repeatCount
)

// PlayerControlsOptions tune a PlayerControls.
type PlayerControlsOptions struct {
	Previous, Next, Shuffle, Repeat, Play string
}

// PlayerControls has previous, play, next, shuffle and repeat.
func PlayerControls(c *ui.Context, label string, playing, shuffle *bool, repeat *Repeat, o PlayerControlsOptions) (root, prev, next ui.Element) {
	if label == "" || playing == nil || shuffle == nil || repeat == nil || *repeat >= repeatCount {
		panic(fmt.Sprintf("media: PlayerControls %q needs a label and play, shuffle and repeat state", label))
	}
	def(&o.Previous, "Previous")
	def(&o.Next, "Next")
	def(&o.Shuffle, "Shuffle")
	def(&o.Repeat, "Repeat")
	def(&o.Play, "Play")
	t := c.Theme()
	root = ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Label(label)
	mode := func(name string) ui.Element {
		return kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(9), t.Space(9)).Shrink(0).Center().Radius(999))
	}
	var sh, rp ui.Element
	root.Children(func() {
		sh = mode(o.Shuffle)
		prev = round(c, o.Previous, "skip-back", t.Space(9))
		if button.PlayButton(c, playing, button.PlayButtonOptions{Label: o.Play}).AlignSelf(ui.Center).Changed() {
			root.SetChanged()
		}
		next = round(c, o.Next, "skip-forward", t.Space(9))
		rp = mode(o.Repeat)
	})
	if sh.Clicked() {
		*shuffle = !*shuffle
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("PlayerControls", "shuffle", "on", *shuffle)
	}
	if rp.Clicked() {
		*repeat = (*repeat + 1) % repeatCount
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("PlayerControls", "repeat", "mode", int(*repeat))
	}
	for _, s := range []string{"previous", "next"} {
		b := prev
		if s == "next" {
			b = next
		}
		if kit.Clicked(b, sound.Toggle) {
			kit.Log("PlayerControls", s)
		}
	}
	style := func(b ui.Element, on bool, glyph string) {
		fg, bg := t.TextMuted, t.SurfaceHover.Alpha(unit(b.Hovered()))
		if on {
			fg, bg = t.Accent, t.Accent.Alpha(0.12)
		}
		b.Checked(on).TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Background(b.AnimateColor("bg", bg, ui.Snappy)).
			Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{}) })
	}
	style(sh, *shuffle, "shuffle")
	if *repeat == RepeatOne {
		style(rp, true, "repeat-1")
	} else {
		style(rp, *repeat == RepeatAll, "repeat")
	}
	return root, prev, next
}
