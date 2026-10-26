package marketing

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Face is one customer avatar.
type Face struct {
	Name  string
	Photo *ui.Bitmap
}

// SocialProofOptions: Text takes the count, zero Rating hides stars.
type SocialProofOptions struct {
	Text   string
	Rating float64
}

// SocialProof is overlapping avatars and a loved-by count.
func SocialProof(c *ui.Context, faces []Face, count int, o SocialProofOptions) ui.Element {
	def(&o.Text, "Loved by %s studios")
	if len(faces) == 0 || len(faces) > 8 || count < len(faces) || !(o.Rating >= 0 && o.Rating <= 5) || math.IsNaN(o.Rating) {
		panic(fmt.Sprintf("marketing: SocialProof needs 1 to 8 faces, a count of at least %d and a rating of 0 to 5", len(faces)))
	}
	names := make([]string, len(faces))
	for i, f := range faces {
		names[i] = f.Name
	}
	unique("SocialProof", names)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Justify(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Shrink(0).AlignItems(ui.Center).Children(func() {
			for i, f := range faces {
				c.Key(strconv.Quote(f.Name))
				a := avatar(c, f.Name, f.Photo, 9).Border(t.Space(0.5), t.Background)
				if i > 0 {
					a.MarginStart(-t.Space(2.5))
				}
				kit.Arrive(a, ui.Pose{Scale: 0.5, Fade: 1}, time.Duration(i)*70*time.Millisecond)
			}
		})
		ui.Column(c).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			if o.Rating > 0 {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					rating(c, o.Rating, 14)
					ui.Text(c, decimal(c, o.Rating, 1)).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				})
			}
			kit.Muted(c, fmt.Sprintf(o.Text, number(c, count))).FontSize(t.Rem(0.84)).Shrink(0).MaxWidth(t.Space(70)).Wrap().
				ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root
}
