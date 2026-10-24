// Package account holds account and billing components.
package account

import (
	"math"
	"net/mail"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func validEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func button(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	return kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
}

func primary(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	return button(c, label).Background(t.Accent).TextColor(t.AccentText)
}
