package display

import "github.com/egoist/mygo/ui"

// Avatar is an image or initials in a circle.
func Avatar(c *ui.Context, name string, image *ui.Bitmap) ui.Element {
	return avatarCircle(c, name, image, c.Theme().Space(10))
}
