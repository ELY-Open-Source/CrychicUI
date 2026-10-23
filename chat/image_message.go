package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/overlay"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ImageMessageOptions tune an ImageMessage.
type ImageMessageOptions struct {
	// Width of the thumbnail; zero gives the default.
	Width              float32
	Back, Forth, Close string
}

// ImageMessage is an image bubble that opens large.
func ImageMessage(c *ui.Context, label string, image *ui.Bitmap, open *bool, o ImageMessageOptions) ui.Element {
	if label == "" || image == nil || open == nil {
		panic(fmt.Sprintf("chat: ImageMessage %q needs a label, an image and open", label))
	}
	t := c.Theme()
	size("ImageMessage", &o.Width, t.Space(60))
	def(&o.Back, "Previous")
	def(&o.Forth, "Next")
	def(&o.Close, "Close")
	w, h := image.Size()
	root := kit.Press(ui.ButtonBase(c).Label(label).AlignSelf(ui.Start).Width(o.Width).Radius(t.Radius*2).Clip().Border(1, t.Border))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	var large ui.Element
	root.Children(func() {
		img := ui.Image(c, image).FillWidth().AspectRatio(float32(w) / float32(h)).Fit(ui.Cover)
		k := 1 + 0.04*img.AnimateTo("zoom", unit(root.Hovered()), ui.Snappy)
		img.Scale(k, k)
		large = overlay.Lightbox(c, label, []*ui.Bitmap{image}, nil, new(int), open, overlay.LightboxOptions{Back: o.Back, Forth: o.Forth, Close: o.Close})
	})
	switch {
	case !*open && kit.Clicked(root, sound.Open):
		*open = true
		root.SetChanged()
		kit.Log("ImageMessage", "open", true)
	case large.Changed():
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("ImageMessage", "open", *open)
	}
	return root
}
