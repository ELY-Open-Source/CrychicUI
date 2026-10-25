package social

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// GiftData is one gift and sent count.
type GiftData struct {
	ID, Name, Icon string
	Count          int
}

// GiftOptions word the send hint; Count formats the grouped total.
type GiftOptions struct{ Send, Count string }

// Gift sends gifts with counts; cue Pop.
func Gift(c *ui.Context, gifts *[]GiftData, selected *int, o GiftOptions) ui.Element {
	if gifts == nil || len(*gifts) == 0 || selected == nil || *selected < -1 || *selected >= len(*gifts) {
		panic("social: Gift needs gifts and selection")
	}
	ids := make([]string, len(*gifts))
	for i, gift := range *gifts {
		if gift.Name == "" || gift.Icon == "" || gift.Count < 0 {
			panic(fmt.Sprintf("social: Gift %q needs a name, icon and count", gift.ID))
		}
		ids[i] = gift.ID
	}
	unique("Gift", ids)
	def(&o.Send, "Send")
	def(&o.Count, "×%s")
	t := c.Theme()
	inks := []ui.Color{t.Danger, t.Warning, t.Accent, t.Success}
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, gift := range *gifts {
			ink := inks[i%len(inks)]
			on := *selected == i
			c.Key(gift.ID)
			b := kit.Press(ui.ButtonBase(c).Label(gift.Name).Column().Gap(t.Space(1)).AlignItems(ui.Center).MinWidth(0).
				Padding(t.Space(3), t.Space(2)).Radius(t.Radius * 1.5).Checked(on).Disabled(gift.Count == math.MaxInt))
			hover := unit(b.Hovered())
			b.Background(b.AnimateColor("bg", t.Background.Mix(ink.Alpha(0.08), max(hover, unit(on))), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(on)), ui.Snappy))
			kit.Stagger(c, root, b, i).Children(func() {
				c.Key(gift.Count)
				tile := ui.Box(c).Size(t.Space(11), t.Space(11)).Center().Radius(999).Background(ink.Alpha(0.14)).TextColor(ink)
				kit.Arrive(tile, ui.Pose{Scale: 1.35, Rotate: -12}, 0).Children(func() { glyph(c, gift.Icon, 22) })
				ui.Text(c, gift.Name).FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).SingleLine().TextAlign(ui.Center)
				kit.Numeric(c, fmt.Sprintf(o.Count, number(c, gift.Count))).FontSize(t.Rem(0.72)).TextColor(ink).
					FillWidth().MinWidth(0).SingleLine().TextAlign(ui.Center)
				kit.Muted(c, o.Send).FontSize(t.Rem(0.68)).Shrink(0)
			})
			if b.Clicked() && gift.Count < math.MaxInt {
				(*gifts)[i].Count++
				*selected = i
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("Gift", "sent", "index", i, "count", (*gifts)[i].Count)
			}
		}
	})
	return root
}
