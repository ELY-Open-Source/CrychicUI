package dashboard

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RangeBarOptions are the range labels.
type RangeBarOptions struct{ Compare, Previous, Next string }

// RangeBar edits date range and compare; cue Select.
func RangeBar(c *ui.Context, start, end *time.Time, compare *bool, o RangeBarOptions) ui.Element {
	if start == nil || end == nil || start.After(*end) || compare == nil {
		panic("dashboard: RangeBar has invalid data")
	}
	def(&o.Compare, "Compare")
	def(&o.Previous, "Previous")
	def(&o.Next, "Next")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		for _, action := range []struct {
			name, glyph string
			days        int
		}{{o.Previous, "chevron-left", -7}, {o.Next, "chevron-right", 7}} {
			b := kit.Press(ui.ButtonBase(c).Label(action.name).Tooltip(action.name).Size(t.Space(7), t.Space(7)).Center().Radius(999))
			b.Children(func() { ui.Icon(c, icon.Must(action.glyph)).Size(16, 16) })
			if b.Clicked() {
				*start, *end = start.AddDate(0, 0, action.days), end.AddDate(0, 0, action.days)
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("RangeBar", "shifted", "days", action.days)
			}
		}
		ui.Text(c, start.Format("2006-01-02")+" – "+end.Format("2006-01-02")).FontSize(t.Rem(0.8)).Grow(1)
		b := kit.Press(ui.ButtonBase(c).Label(o.Compare).Padding(t.Space(1), t.Space(2)).Radius(999).
			Border(1, t.Border).Checked(*compare))
		b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*compare)), ui.Snappy))
		b.Children(func() { ui.Text(c, o.Compare).FontSize(t.Rem(0.8)) })
		if b.Clicked() {
			*compare = !*compare
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("RangeBar", "compare", "value", *compare)
		}
	})
	return root
}
