package form

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// FieldGroupOptions tune a FieldGroup.
type FieldGroupOptions struct {
	// Help and Error are caller copy; errors mark controls invalid.
	Help, Error string
}

// FieldGroup labels a control; its error shakes in.
func FieldGroup(c *ui.Context, label string, o FieldGroupOptions, control func()) ui.Element {
	if label == "" || control == nil {
		panic("form: FieldGroup needs a label and a control")
	}
	t := c.Theme()
	f := ui.Field(c, label, control).FillWidth().Description(o.Help).Error(o.Error)
	last := ui.Local(f, "error", func() string { return "" })
	run := ui.Local(f, "run", func() uint64 { return 0 })
	if o.Error != *last {
		if *last = o.Error; o.Error != "" {
			(*run)++
		}
	}
	if *run > 0 {
		d := t.Space(1)
		f.Translate(f.Keyframes("shake", []ui.Keyframe{{Value: 0}, {At: 50 * time.Millisecond, Value: -d},
			{At: 100 * time.Millisecond, Value: d}, {At: 160 * time.Millisecond, Value: 0, Ease: ui.EaseOut}},
			ui.KeyframeOptions{Run: *run}), 0)
	}
	return f
}
