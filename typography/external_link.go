package typography

import "github.com/egoist/mygo/ui"

// ExternalLink is a link with an arrow nudging on hover.
func ExternalLink(c *ui.Context, label, url string) ui.Element {
	return link(c, label, url, true)
}
