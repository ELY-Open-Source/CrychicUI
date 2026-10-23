// Code components for the gallery.
package code

import (
	"crypto/sha256"
	"fmt"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HashTool shows the SHA-256 of text; digests cross-fade.
func HashTool(c *ui.Context, label string, text *string) ui.Element {
	if label == "" || text == nil {
		panic("code: HashTool needs a label and a bound text")
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(90)).Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		in := ui.TextArea(c, text).Label(label).Font("monospace").Lines(2, 5).FillWidth()
		if in.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("HashTool", "edited")
		}
		sum := sha256.Sum256([]byte(*text))
		box := ui.Column(c).FillWidth().Padding(t.Space(2)).Background(t.Surface).Radius(t.Radius)
		box.Children(func() {
			c.Key(fmt.Sprintf("%x", sum))
			ui.Text(c, fmt.Sprintf("%x", sum)).Font("monospace").FontSize(t.Rem(0.8)).Wrap().
				Transition(ui.ElementTransition{Enter: &ui.Motion{}, Exit: &ui.Motion{}})
		})
	})
	return root
}

// UnixTimeOptions are the copy knobs of UnixTime.
