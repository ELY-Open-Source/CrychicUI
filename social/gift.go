package social

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// GiftData is one gift and sent count.
type GiftData struct {
	ID, Name, Icon string
	Count          int
}

// GiftOptions are the gift send labels.
type GiftOptions struct{ Send string }

// Gift sends gifts with counts; cue Pop.
func Gift(c *ui.Context, gifts *[]GiftData, selected *int, o GiftOptions) ui.Element {
	if gifts == nil || len(*gifts) == 0 || selected == nil || *selected < -1 || *selected >= len(*gifts) {
		panic("social: Gift needs gifts and selection")
	}
	def(&o.Send, "Send")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, gift := range *gifts {
			if gift.ID == "" || gift.Name == "" || gift.Icon == "" || gift.Count < 0 {
				panic("social: invalid gift")
			}
			c.Key(gift.ID)
			card := button(c, gift.Name).Column().Gap(t.Space(1)).AlignItems(ui.Center).Border(1, t.Border).
				Checked(*selected == i).Disabled(gift.Count == math.MaxInt)
			card.Children(func() {
				ui.Icon(c, icon.Must(gift.Icon)).Size(t.Space(6), t.Space(6)).TextColor(t.Accent)
				ui.Text(c, gift.Name)
				kit.Numeric(c, fmt.Sprint(gift.Count))
				kit.Muted(c, o.Send)
			})
			if card.Clicked() && gift.Count < math.MaxInt {
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
