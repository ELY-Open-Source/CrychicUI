package code

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// JSONFormatOptions are the copy knobs of JSONFormat.
type JSONFormatOptions struct {
	Format, Minify, Valid, Invalid, Text string
	FormatError                          func(err error) string
}

// JSONFormat formats or minifies JSON; cue Success, Error.
func JSONFormat(c *ui.Context, label string, text *string, o JSONFormatOptions) (root, format, minify ui.Element) {
	if label == "" || text == nil {
		panic(fmt.Sprintf("code: JSONFormat %q needs a label and a bound text", label))
	}
	def(&o.Format, "Format")
	def(&o.Minify, "Minify")
	def(&o.Valid, "Valid")
	def(&o.Invalid, "Invalid")
	def(&o.Text, "JSON")
	if o.FormatError == nil {
		o.FormatError = func(err error) string { return "invalid JSON: " + err.Error() }
	}
	t := c.Theme()
	root = kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	bad := ui.Local(root, "bad", func() string { return "" })
	root.Children(func() {
		panelHead(c, "braces", label, func() {
			typography.Bytes(c, int64(len(*text))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			if json.Valid([]byte(*text)) {
				chip(c, o.Valid, t.Success)
			} else {
				chip(c, o.Invalid, t.Danger)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			if ui.TextArea(c, text).Label(o.Text).Font("monospace").Lines(4, 10).FillWidth().Changed() {
				*bad = ""
				edited(root, "JSONFormat")
			}
			if *bad != "" {
				ui.Row(c).Key("bad").FillWidth().Gap(t.Space(2)).Padding(t.Space(2)).Radius(t.Radius).Background(t.Danger.Alpha(0.08)).
					Transition(kit.Enter).Children(func() {
					ui.Icon(c, icon.Must("circle-x")).Size(14, 14).Shrink(0).TextColor(t.Danger)
					ui.Text(c, *bad).FontSize(t.Rem(0.8)).TextColor(t.Danger).Grow(1).MinWidth(0).Wrap()
				})
			}
			// Keyed: the error row appearing must not move the actions.
			c.Key("actions")
			ui.Row(c).Gap(t.Space(2)).Children(func() {
				format = action(c, o.Format, "wand-sparkles", true)
				minify = action(c, o.Minify, "fold-horizontal", false)
			})
		})
	})
	reshape := func(pretty bool) {
		var v any
		dec := json.NewDecoder(strings.NewReader(*text))
		dec.UseNumber()
		err := dec.Decode(&v)
		if err == nil {
			var extra any
			if dec.Decode(&extra) != io.EOF {
				err = fmt.Errorf("more than one JSON value")
			}
		}
		if err != nil {
			*bad = o.FormatError(err)
			root.SetChanged()
			sound.Play(sound.Error)
			kit.Log("JSONFormat", "invalid")
			return
		}
		var out bytes.Buffer
		src := []byte(strings.TrimSpace(*text))
		if pretty {
			err = json.Indent(&out, src, "", "  ")
		} else {
			err = json.Compact(&out, src)
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
