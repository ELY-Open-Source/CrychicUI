package form

import (
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/egoist/mygo/ui"
)

// FormRow sets fields side by side, stacking them when narrow.
func FormRow(c *ui.Context, fields []func()) ui.Element {
	if len(fields) == 0 {
		panic("form: FormRow needs fields")
	}
	return layout.Responsive(c, float32(len(fields))*c.Theme().Rem(14), fields)
}
