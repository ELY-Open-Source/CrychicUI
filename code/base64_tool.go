// Code components for the gallery.
package code

import (
	"encoding/base64"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

type Base64ToolOptions struct {
	Encode   string
	Decode   string
	Invalid  string
	YourText string
}

// Base64Tool encodes and decodes panes; cue Tick.

func Base64Tool(c *ui.Context, label string, text *string, encode *bool, o Base64ToolOptions) ui.Element {
	if label == "" || text == nil || encode == nil {
		panic("code: Base64Tool needs a label and bound values")
	}
	def(&o.Encode, "Encode")
	def(&o.Decode, "Decode")
	def(&o.Invalid, "invalid base64")
	def(&o.YourText, "Text")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(90)).Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).Children(func() {
			seg := func(name string, on bool) ui.Element {
				b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1), t.Space(2.5)).Radius(t.Radius).Checked(on))
				bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
				if on {
					bg = t.Accent.Alpha(0.12)
				}
				b.Background(b.AnimateColor("bg", bg, ui.Snappy))
				return b.Children(func() { ui.Text(c, name).FontSize(t.Rem(0.85)).SingleLine() })
			}
			if seg(o.Encode, *encode).Clicked() && !*encode {
				*encode = true
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("Base64Tool", "encode")
			}
			if seg(o.Decode, !*encode).Clicked() && *encode {
				*encode = false
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("Base64Tool", "decode")
			}
		})
		in := ui.TextArea(c, text).Label(o.YourText).Font("monospace").Lines(2, 5).FillWidth()
		if in.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("Base64Tool", "edited")
		}
		out, err := "", error(nil)
		if *encode {
			out = base64.StdEncoding.EncodeToString([]byte(*text))
		} else {
			var b []byte
			b, err = base64.StdEncoding.DecodeString(*text)
			out = string(b)
		}
		box := ui.Column(c).FillWidth().Padding(t.Space(2)).Background(t.Surface).Radius(t.Radius)
		box.Children(func() {
			if err != nil {
				ui.Text(c, o.Invalid).FontSize(t.Rem(0.85)).TextColor(t.Danger)
			} else if out != "" {
				ui.Text(c, out).Font("monospace").FontSize(t.Rem(0.85)).Wrap().ContentTransition(ui.ContentTransition{})
			} else {
				kit.Muted(c, "…").FontSize(t.Rem(0.85))
			}
		})
	})
	return root
}

// HashTool is the SHA-256 of a text; the digest cross-fades.
