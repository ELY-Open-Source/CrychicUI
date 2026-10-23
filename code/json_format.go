// Code components for the gallery.
package code

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"io"
	"strings"
)

type JSONFormatOptions struct {
	Format      string
	Minify      string
	FormatError func(err error) string
}

// JSONFormat formats or minifies JSON; cue Success, Error.

func JSONFormat(c *ui.Context, label string, text *string, o JSONFormatOptions) (root, format, minify ui.Element) {
	if label == "" || text == nil {
		panic("code: JSONFormat needs a label and a bound text")
	}
	def(&o.Format, "Format")
	def(&o.Minify, "Minify")
	if o.FormatError == nil {
		o.FormatError = func(err error) string { return "invalid JSON: " + err.Error() }
	}
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(90)).Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	bad := ui.Local(root, "bad", func() string { return "" })
	root.Children(func() {
		in := ui.TextArea(c, text).Label(label).Font("monospace").Lines(3, 8).FillWidth()
		if in.Changed() {
			*bad = ""
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("JSONFormat", "edited")
		}
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			format = kit.Press(ui.ButtonBase(c).Label(o.Format).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText))
			format.Children(func() { ui.Text(c, o.Format).SingleLine() })
			minify = kit.Press(ui.ButtonBase(c).Label(o.Minify).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Border(1, t.Border))
			minify.Children(func() { ui.Text(c, o.Minify).SingleLine() })
		})
		if *bad != "" {
			ui.Text(c, *bad).FontSize(t.Rem(0.85)).TextColor(t.Danger).Wrap()
		}
	})
	reshape := func(pretty bool) {
		var v any
		dec := json.NewDecoder(strings.NewReader(*text))
		dec.UseNumber()
		if err := dec.Decode(&v); err != nil {
			*bad = o.FormatError(err)
			root.SetChanged()
			sound.Play(sound.Error)
			kit.Log("JSONFormat", "invalid")
			return
		}
		var extra any
		if err := dec.Decode(&extra); err != io.EOF {
			*bad = o.FormatError(fmt.Errorf("more than one JSON value"))
			root.SetChanged()
			sound.Play(sound.Error)
			kit.Log("JSONFormat", "invalid")
			return
		}
		var out bytes.Buffer
		src := strings.TrimSpace(*text)
		var err error
		if pretty {
			err = json.Indent(&out, []byte(src), "", "  ")
		} else {
			err = json.Compact(&out, []byte(src))
		}
		if err != nil {
			panic(fmt.Sprintf("code: JSONFormat cannot reshape validated JSON: %v", err))
		}
		*text = out.String()
		*bad = ""
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("JSONFormat", "reshaped", "pretty", pretty)
	}
	if format.Clicked() {
		reshape(true)
	}
	if minify.Clicked() {
		reshape(false)
	}
	return root, format, minify
}

// Base64ToolOptions are the copy knobs of Base64Tool.
