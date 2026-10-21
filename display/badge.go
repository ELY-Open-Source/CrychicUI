package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BadgeVariant is a Badge's semantic color.
type BadgeVariant int

// Badge variants.
const (
	BadgeNeutral BadgeVariant = iota
	BadgeAccent
	BadgeSuccess
	BadgeWarning
	BadgeDanger
)

// BadgeOptions tune a Badge.
type BadgeOptions struct {
	// Variant chooses the badge color.
	Variant BadgeVariant
}

// Badge is a small semantic pill popping in.
func Badge(c *ui.Context, text string, o BadgeOptions) ui.Element {
	if o.Variant < BadgeNeutral || o.Variant > BadgeDanger {
		panic("display: Badge variant out of range")
	}
	t := c.Theme()
	bg, fg, border := t.Surface, t.Text, t.Border
	switch o.Variant {
	case BadgeAccent:
		bg, fg, border = t.Accent, t.AccentText, t.Accent
	case BadgeSuccess:
		bg, fg, border = t.Success, t.AccentText, t.Success
	case BadgeWarning:
		bg, fg, border = t.Warning, t.AccentText, t.Warning
	case BadgeDanger:
		bg, fg, border = t.Danger, t.AccentText, t.Danger
	}
	root := ui.Row(c).AlignSelf(ui.Start).Padding(t.Space(0.5), t.Space(2)).Radius(999).
		Background(bg).TextColor(fg).Border(1, border).Children(func() { ui.Text(c, text).FontSize(t.Rem(0.85)) })
	return kit.Arrive(root, ui.Pose{Scale: 0.8, Fade: 1}, 0)
}
