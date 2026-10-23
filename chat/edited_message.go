package chat

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// EditedMessageOptions tune an EditedMessage mark.
type EditedMessageOptions struct {
	Mark string
}

// EditedMessage is an "(edited)" mark whose tip shows the time.
func EditedMessage(c *ui.Context, at time.Time, o EditedMessageOptions) ui.Element {
	if at.IsZero() {
		panic("chat: EditedMessage needs a time")
	}
	def(&o.Mark, "(edited)")
	t := c.Theme()
	full := i18n.FormatDate(c.Region().FormatLocale(c.Locale()), at, i18n.DateOptions{Date: i18n.MediumStyle, Time: i18n.ShortStyle})
	root := kit.Muted(c, o.Mark).FontSize(t.Rem(0.75)).Italic().AlignSelf(ui.Start).Tooltip(full)
	return kit.Arrive(root, ui.Pose{Fade: 1}, 0)
}
