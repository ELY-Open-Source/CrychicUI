package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PianoOptions tune a Piano.
type PianoOptions struct {
	// From is the lowest white MIDI key; zero gives C4.
	From int
	// Keys counts white keys; zero gives 14.
	Keys int
	// Height of the keys; zero gives the default.
	Height float32
}

// Piano keys play notes and press down.
func Piano(c *ui.Context, label string, note *int, o PianoOptions) ui.Element {
	if o.From == 0 {
		o.From = 60
	}
	if o.Keys == 0 {
		o.Keys = 14
	}
	if label == "" || o.From < 16 || !white(o.From) || o.Keys < 2 {
		panic(fmt.Sprintf("media: Piano %q needs a label, a white lowest key from 16 and two keys", label))
	}
	var whites []int
	for m := o.From; len(whites) < o.Keys; m++ {
		if white(m) {
			whites = append(whites, m)
		}
	}
	top := whites[len(whites)-1]
	if top > 127 {
		panic(fmt.Sprintf("media: Piano %q keys must end by MIDI 127", label))
	}
	chosen("Piano", label, note, -1, top+1)
	if *note >= 0 && *note < o.From {
		panic(fmt.Sprintf("media: Piano %q note must be -1 or on the keys", label))
	}
	t := c.Theme()
	size("Piano", label, &o.Height, t.Space(36))
	root := ui.Box(c).FillWidth().Height(o.Height).Label(label)
	n := float32(len(whites))
	var notes []int
	var keys []ui.Element
	root.Children(func() {
		ui.Row(c).Fill().Gap(t.Space(0.5)).Children(func() {
			for _, m := range whites {
				c.Key(m)
				notes = append(notes, m)
				keys = append(keys, ui.ButtonBase(c).Label(named(m)).Grow(1).Basis(0).FillHeight().
					Radius(t.Radius/2).Border(1, t.Border))
			}
		})
		for i, m := range whites[:len(whites)-1] {
			if white(m + 1) {
				continue
			}
			c.Key(m + 1)
			notes = append(notes, m+1)
			keys = append(keys, ui.ButtonBase(c).Label(named(m+1)).Absolute().Top(0).HeightPercent(62).
				LeftPercent((float32(i+1)-0.3)/n*100).WidthPercent(0.6/n*100).Radius(t.Radius/2))
		}
	})
	for i, e := range keys {
		m := notes[i]
		run := ui.Local(e, "run", func() int { return 0 })
		if strike(e) {
			*run++
			*note = m
			root.SetChanged()
			sound.PlayNote(pitch(m))
			kit.Log("Piano", "played", "note", m)
		}
		depth := max(e.AnimateTo("down", unit(e.Pressed()), ui.Snappy), flash(e, *run))
		face := t.Background
		if !white(m) {
			face = t.Text
		}
		e.Background(face.Mix(t.Accent, 0.35*depth)).Translate(0, depth*t.Space(0.75))
	}
	return root
}
