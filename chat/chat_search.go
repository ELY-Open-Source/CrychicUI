package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChatSearchOptions tune a ChatSearch.
type ChatSearchOptions struct {
	Placeholder, Empty string
}

// ChatSearch finds messages; a click on one jumps to it.
func ChatSearch(c *ui.Context, label string, list []Message, query *string, jump *int, o ChatSearchOptions) (root, field ui.Element) {
	if label == "" || query == nil || jump == nil {
		panic(fmt.Sprintf("chat: ChatSearch %q needs a label, a query and a jump", label))
	}
	messages("ChatSearch", list)
	def(&o.Placeholder, "Search messages")
	def(&o.Empty, "No messages match")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		field = input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder})
		q := strings.ToLower(strings.TrimSpace(*query))
		if q == "" {
			return
		}
		var hits []int
		for i, m := range list {
			if strings.Contains(strings.ToLower(m.Text), q) {
				hits = append(hits, i)
			}
		}
		if len(hits) == 0 {
			kit.Arrive(kit.Muted(c, o.Empty).Padding(t.Space(2)), ui.Pose{Fade: 1}, 0)
			return
		}
		ui.Scroll(c).FillWidth().MaxHeight(t.Space(60)).Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Children(func() {
			for n, i := range hits {
				c.Key(i)
				m := list[i]
				b := kit.Press(ui.ButtonBase(c).Label(m.Text).Column().FillWidth().Gap(t.Space(0.5)).
					Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius))
				kit.Arrive(b, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(n)*20*time.Millisecond)
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
						ui.Text(c, m.Sender).Bold().FontSize(t.Rem(0.8)).Grow(1).SingleLine()
						kit.Muted(c, when(c, m.At)).FontSize(t.Rem(0.75))
					})
					ui.Text(c, m.Text).SingleLine().FontSize(t.Rem(0.9))
				})
				if b.Clicked() {
					*jump = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("ChatSearch", "jumped", "index", i, "hits", len(hits))
				}
			}
		})
	})
	return root, field
}
