package display

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

// Version is a version pill cross-fading on change.
func Version(c *ui.Context, version string) ui.Element {
	if version == "" {
		panic("display: Version needs a version")
	}
	t := c.Theme()
	return ui.Text(c, "v"+strings.TrimPrefix(version, "v")).Padding(t.Space(0.5), t.Space(2)).Radius(999).
		Background(t.Surface).Border(1, t.Border).ContentTransition(ui.ContentTransition{})
}
