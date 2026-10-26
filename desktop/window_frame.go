package desktop

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

// WindowFrameOptions mark an inactive window.
type WindowFrameOptions struct{ Inactive bool }

// WindowFrame is a window with a title bar around body.
func WindowFrame(c *ui.Context, title string, o WindowFrameOptions, body func()) ui.Element {
	if strings.TrimSpace(title) == "" || body == nil {
		panic(fmt.Sprintf("desktop: WindowFrame %q needs a title and a body", title))
	}
	t := c.Theme()
	root := window(c).Role(ui.RoleGroup).Label(title)
	if o.Inactive {
		root.Shadow(0, 2, 8, 0, t.Text.Alpha(0.08))
	}
	return root.Children(func() {
		titleBar(c, title, !o.Inactive, nil)
		c.Key("body")
		ui.Column(c).FillWidth().Padding(t.Space(4)).Children(body)
	})
}
