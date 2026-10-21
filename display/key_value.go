package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// KeyValuePair is one key and value row.
type KeyValuePair struct{ Key, Value string }

// KeyValue is one key and value pair; value cross-fades.
func KeyValue(c *ui.Context, key, value string) ui.Element {
	t := c.Theme()
	return ui.Row(c).FillWidth().Gap(t.Space(3)).Children(func() {
		kit.Muted(c, key).Width(t.Rem(10))
		ui.Text(c, value).Grow(1).ContentTransition(ui.ContentTransition{})
	})
}
