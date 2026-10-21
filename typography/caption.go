package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Caption is small muted text with cross-fading changes.
func Caption(c *ui.Context, text string) ui.Element {
	return kit.Muted(c, text).FontSize(c.Theme().Rem(0.85)).ContentTransition(ui.ContentTransition{})
}
