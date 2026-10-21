package input

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TextFieldOptions tune a TextField.
type TextFieldOptions struct {
	// Placeholder, Help and Error are caller copy.
	Placeholder, Help, Error string
}

// TextField is a labeled input with a shaking error.
func TextField(c *ui.Context, label string, value *string, o TextFieldOptions) ui.Element {
	if value == nil {
		panic("input: TextField needs a value")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	last := ui.Local(root, "error", func() string { return o.Error })
	run := ui.Local(root, "run", func() uint64 { return 0 })
	if o.Error != *last {
		*last = o.Error
		(*run)++
	}
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		in := ui.TextInputBase(c, value).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).
			Radius(t.Radius).Background(t.Background)
		border := t.Border
		if o.Error != "" {
			border = t.Danger
		} else if in.Focused() {
			border = t.Accent
		}
		in.Border(1, in.AnimateColor("border", border, ui.Snappy))
		changed(root, in)
		if o.Error != "" {
			err := ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85))
			shake := err.Keyframes("shake", []ui.Keyframe{
				{Value: 0}, {At: 50 * time.Millisecond, Value: -4}, {At: 100 * time.Millisecond, Value: 4},
				{At: 160 * time.Millisecond, Value: 0, Ease: ui.EaseOut},
			}, ui.KeyframeOptions{Run: *run})
			err.Translate(shake, 0)
		} else if o.Help != "" {
			kit.Muted(c, o.Help).FontSize(t.Rem(0.85))
		}
	})
}
