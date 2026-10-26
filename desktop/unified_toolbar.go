package desktop

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UnifiedToolbarOptions word the toolbar.
type UnifiedToolbarOptions struct{ Subtitle, Placeholder, Back, Forward string }

// UnifiedToolbar has nav, a title and search; cue Tap.
func UnifiedToolbar(c *ui.Context, title string, query *string, o UnifiedToolbarOptions) (root, back, forward ui.Element) {
	def(&o.Placeholder, "Search")
	def(&o.Back, "Back")
	def(&o.Forward, "Forward")
	if strings.TrimSpace(title) == "" || query == nil {
		panic(fmt.Sprintf("desktop: UnifiedToolbar %q needs a title and a query", title))
	}
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).
		Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Role(ui.RoleToolbar)
	kit.Arrive(root, ui.Pose{Y: -t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		lights(c, true)
		nav := func(label, mark string) ui.Element {
			b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1)).Radius(t.Radius))
			b.Background(b.AnimateColor("bg", t.Border.Alpha(0.6*unit(b.Hovered())), ui.Snappy))
			b.Children(func() { glyph(c, mark, 16).TextColor(t.TextMuted) })
			if b.Clicked() {
				sound.Play(sound.Tap)
				kit.Log("UnifiedToolbar", "nav", "to", label)
			}
			return b
		}
		ui.Row(c).Shrink(0).Children(func() {
			back = nav(o.Back, "chevron-left")
			forward = nav(o.Forward, "chevron-right")
		})
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			ui.Text(c, title).FontSize(t.Rem(0.86)).Bold().FillWidth().MinWidth(0).SingleLine()
			if o.Subtitle != "" {
				kit.Muted(c, o.Subtitle).FontSize(t.Rem(0.7)).FillWidth().MinWidth(0).SingleLine()
			}
		})
		c.Key("search")
		frame := ui.Row(c).Basis(t.Space(44)).Shrink(1).MinWidth(t.Space(20)).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1.25), t.Space(2)).
			Radius(t.Radius).Background(t.Background)
		frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
		frame.Children(func() {
			glyph(c, "search", 14).TextColor(t.TextMuted)
			in := ui.TextInputBase(c, query).Label(o.Placeholder).Placeholder(o.Placeholder).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0)
			if in.Changed() {
				root.SetChanged()
				kit.Log("UnifiedToolbar", "query", "runes", utf8.RuneCountInString(*query))
			}
		})
	})
	return root, back, forward
}
