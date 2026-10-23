package files

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UploadItem is one upload and its completion fraction.
type UploadItem struct {
	Name     string
	Progress float64
}

// UploadListOptions are the upload labels.
type UploadListOptions struct{ Cancel, Empty string }

// UploadList shows upload bars with cancel; cue Close.
func UploadList(c *ui.Context, label string, uploads *[]UploadItem, o UploadListOptions) ui.Element {
	if label == "" || uploads == nil {
		panic("files: UploadList needs a label and uploads")
	}
	def(&o.Cancel, "Cancel")
	def(&o.Empty, "No uploads")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	cancel := -1
	root.Children(func() {
		if len(*uploads) == 0 {
			kit.Muted(c, o.Empty).Padding(t.Space(2))
		}
		for i, u := range *uploads {
			if u.Name == "" || u.Progress < 0 || u.Progress > 1 || math.IsNaN(u.Progress) || math.IsInf(u.Progress, 0) {
				panic(fmt.Sprintf("files: UploadList invalid upload %q", u.Name))
			}
			c.Key(i)
			row := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			row.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must(glyph(u.Name))).Size(16, 16).TextColor(t.TextMuted)
					ui.Text(c, u.Name).SingleLine().Grow(1)
					kit.Numeric(c, fmt.Sprintf("%.0f%%", u.Progress*100)).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
					name := fmt.Sprintf("%s %s", o.Cancel, u.Name)
					b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(6), t.Space(6)).Center().Radius(999))
					b.Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
					if b.Clicked() {
						cancel = i
					}
				})
				bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
				w := bar.AnimateTo("progress", float32(u.Progress), ui.Smooth)
				bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, t.Accent, 999)
				}))
			})
		}
	})
	if cancel >= 0 {
		*uploads = append((*uploads)[:cancel], (*uploads)[cancel+1:]...)
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("UploadList", "canceled", "left", len(*uploads))
	}
	return root
}
