package media

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
)

// QualityMenuOptions tune a QualityMenu.
type QualityMenuOptions struct{}

// QualityMenu picks a quality; the chosen one is checked.
func QualityMenu(c *ui.Context, label string, names []string, at *int, open *bool, o QualityMenuOptions) ui.Element {
	if label == "" || len(names) == 0 || open == nil {
		panic(fmt.Sprintf("media: QualityMenu %q needs a label, names and state", label))
	}
	for i, n := range names {
		if n == "" || slices.Index(names, n) != i {
			panic(fmt.Sprintf("media: QualityMenu %q names must be distinct and not empty", label))
		}
	}
	chosen("QualityMenu", label, at, 0, len(names))
	b := opener(c, label, names[*at], *open)
	menu(c, b, "QualityMenu", label, names, at, open)
	return b
}
