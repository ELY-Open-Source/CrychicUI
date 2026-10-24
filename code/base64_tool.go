package code

import (
	"encoding/base64"
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// Base64ToolOptions are the copy knobs of Base64Tool.
type Base64ToolOptions struct {
	Encode, Decode, Invalid string
	// YourText and Base64 name the two panes.
	YourText, Base64, Copy string
}

// Base64Tool encodes or decodes between two panes; cue Tick.
func Base64Tool(c *ui.Context, label string, text *string, encode *bool, o Base64ToolOptions) (root, copy ui.Element) {
	if label == "" || text == nil || encode == nil {
		panic(fmt.Sprintf("code: Base64Tool %q needs a label and bound values", label))
	}
	def(&o.Encode, "Encode")
	def(&o.Decode, "Decode")
	def(&o.Invalid, "invalid base64")
	def(&o.YourText, "Text")
	def(&o.Base64, "Base64")
	def(&o.Copy, "Copy")
	t := c.Theme()
	var out string
	var err error
	// convert reads the text after its input applied edits.
	convert := func() {
		err = nil
		if *encode {
			out = base64.StdEncoding.EncodeToString([]byte(*text))
			return
		}
		var b []byte
		b, err = base64.StdEncoding.DecodeString(*text)
		out = string(b)
	}
	from, to := o.YourText, o.Base64
	if !*encode {
		from, to = to, from
	}
	root = kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	var enc, dec ui.Element
	root.Children(func() {
		panelHead(c, "binary", label, func() { enc, dec = pair(c, o.Encode, o.Decode, "lock", "lock-open", *encode) })
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			captioned(c, from, func() {
				if ui.TextArea(c, text).Label(from).Font("monospace").Lines(2, 5).FillWidth().Changed() {
					edited(root, "Base64Tool")
				}
			})
			convert()
			ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.Center).AlignItems(ui.Center).Children(func() {
				typography.Bytes(c, int64(len(*text))).FontSize(t.Rem(0.7)).TextColor(t.TextMuted)
				arrow := ui.Icon(c, icon.Must("arrow-down")).Size(16, 16).TextColor(t.TextMuted)
				arrow.Rotate(180 * arrow.AnimateTo("turn", unit(!*encode), ui.Snappy))
				typography.Bytes(c, int64(len(out))).FontSize(t.Rem(0.7)).TextColor(t.TextMuted)
			})
			fieldLabel(c, to)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(1), t.Space(2), t.Space(3)).
				Background(t.Surface).Border(1, t.Border).Radius(t.Radius).Children(func() {
				switch {
				case err != nil:
					ui.Icon(c, icon.Must("circle-x")).Size(14, 14).Shrink(0).TextColor(t.Danger)
					ui.Text(c, o.Invalid).FontSize(t.Rem(0.85)).TextColor(t.Danger).Grow(1)
				default:
					ui.Text(c, out).Font("monospace").FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap().ContentTransition(ui.ContentTransition{})
				}
				copy = button.CopyButton(c, o.Copy, out).Disabled(err != nil || out == "")
			})
		})
	})
	flip := func(to bool, event string) {
		*encode = to
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Base64Tool", event)
	}
	switch {
	case enc.Clicked() && !*encode:
		flip(true, "encode")
	case dec.Clicked() && *encode:
		flip(false, "decode")
	}
	if copy.Clicked() {
		kit.Log("Base64Tool", "copied", "bytes", len(out))
	}
	return root, copy
}
