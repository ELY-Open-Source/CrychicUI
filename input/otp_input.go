package input

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// codeInput is the shared boxed one-time-code input.
func codeInput(c *ui.Context, label string, value *string, length int, masked, success bool) ui.Element {
	if value == nil || length <= 0 || length > 12 {
		panic("input: code input needs a value and 1-12 boxes")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	wasFull := ui.Local(root, "full", func() bool { return utf8.RuneCountInString(*value) == length })
	return root.Children(func() {
		c.Key("field")
		in := ui.TextInput(c, value).Label(label).Absolute().Fill().Opacity(0)
		if masked {
			in.Password()
		}
		changedNow := changed(root, in)
		if changedNow {
			clean := strings.Map(func(r rune) rune {
				if unicode.IsDigit(r) {
					return r
				}
				return -1
			}, *value)
			if utf8.RuneCountInString(clean) > length {
				clean = string([]rune(clean)[:length])
			}
			if clean != *value {
				*value = clean
				root.SetChanged()
			}
		}
		full := utf8.RuneCountInString(*value) == length
		if success && changedNow && full && !*wasFull {
			sound.Play(sound.Success)
		}
		*wasFull = full
		runes := []rune(*value)
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			for i := range length {
				box := ui.Box(c).Key(i).Size(t.Space(10), t.Space(10)).Center().Radius(t.Radius).Background(t.Background)
				border := t.Border
				if in.Focused() && i == min(len(runes), length-1) {
					border = t.Accent
				}
				box.Border(1, box.AnimateColor("border", border, ui.Snappy)).Children(func() {
					if i < len(runes) {
						text := string(runes[i])
						if masked {
							text = "•"
						}
						textElement := ui.Text(c, text).ContentTransition(ui.ContentTransition{})
						kit.Arrive(textElement, ui.Pose{Scale: 0.5, Fade: 1}, 0)
					}
				})
			}
		})
	})
}

// OTPInputOptions tune an OTPInput.
type OTPInputOptions struct {
	// Label names the hidden input.
	Label string
}

// OTPInput is n digit boxes with paste and completion cues.
func OTPInput(c *ui.Context, value *string, length int, o OTPInputOptions) ui.Element {
	if o.Label == "" {
		o.Label = "One-time code"
	}
	return codeInput(c, o.Label, value, length, false, true)
}
