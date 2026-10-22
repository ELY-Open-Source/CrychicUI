package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// DueChipOptions tune a DueChip.
type DueChipOptions struct {
	// Soon is the days ahead that warn; 2 when zero.
	Soon int
}

// DueChip shows a due date tinted by its nearness.
func DueChip(c *ui.Context, due time.Time, o DueChipOptions) ui.Element {
	if due.IsZero() || o.Soon < 0 {
		panic(fmt.Sprintf("datetime: DueChip needs a due date and Soon >= 0, not %d", o.Soon))
	}
	if o.Soon == 0 {
		o.Soon = 2
	}
	t := c.Theme()
	now := c.Now()
	c.After(tomorrow(now).Sub(now))
	bg, edge, fg := t.Surface, t.Border, t.Text
	switch n := days(now, due); {
	case n < 0:
		bg, edge, fg = t.Danger, t.Danger, t.AccentText
	case n <= o.Soon:
		bg, edge, fg = t.Warning, t.Warning, t.AccentText
	}
	text := i18n.FormatDate(locale(c), due, i18n.DateOptions{Skeleton: "MMMd"})
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(1.5)).Padding(t.Space(0.5), t.Space(2.5)).
		Radius(999).Role(ui.RoleText).Label(text)
	root.Background(root.AnimateColor("bg", bg, ui.Smooth)).Border(1, root.AnimateColor("edge", edge, ui.Smooth)).
		TextColor(root.AnimateColor("fg", fg, ui.Smooth))
	return root.Children(func() {
		ui.Icon(c, icon.Must("calendar")).Size(16, 16)
		ui.Text(c, text).FontSize(t.Rem(0.85))
	})
}
