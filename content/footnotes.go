package content

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FootnotesOptions word the list; Back names a back link.
type FootnotesOptions struct{ Title, Back string }

// Footnotes lists notes; every back click reports; cue Tap.
func Footnotes(c *ui.Context, notes []string, back *int, o FootnotesOptions) ui.Element {
	def(&o.Title, "Notes")
	def(&o.Back, "Back to reference %d")
	if len(notes) == 0 || back == nil {
		panic(fmt.Sprintf("content: Footnotes %q needs notes and a back index", o.Title))
	}
	for i, n := range notes {
		if strings.TrimSpace(n) == "" {
			panic(fmt.Sprintf("content: Footnotes %q note %d is blank", o.Title, i+1))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3), 0, 0, 0).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Role(ui.RoleList)
	root.Children(func() {
		kit.Muted(c, o.Title).FontSize(t.Rem(0.74)).Bold().FillWidth().MinWidth(0).SingleLine()
		list := ui.Column(c).FillWidth().Gap(t.Space(1.5))
		list.Children(func() {
			for i, n := range notes {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Role(ui.RoleListItem)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Text(c, strconv.Itoa(i+1)).FontSize(t.Rem(0.66)).Bold().TextColor(t.Accent).Width(t.Space(5)).Shrink(0).
						TextAlign(ui.Center).Padding(t.Space(0.25), 0).Radius(999).Background(t.Accent.Alpha(0.1))
					ui.Text(c, n).FontSize(t.Rem(0.8)).TextColor(t.Text.Alpha(0.86)).LineHeight(1.5).Grow(1).MinWidth(0).Wrap()
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Back, i+1)).Shrink(0).Padding(t.Space(0.75)).Radius(999))
					b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
					b.Children(func() {
						glyph(c, "corner-left-up", 14).TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(b.Hovered())), ui.Snappy))
					})
					if b.Clicked() {
						*back = i
						root.SetChanged()
						sound.Play(sound.Tap)
						kit.Log("Footnotes", "back", "index", i)
					}
				})
			}
		})
	})
	return root
}
