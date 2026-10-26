package content

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CodeSample is one language's code under a tab name.
type CodeSample struct{ Name, Lang, Code string }

// DocCodeTabsOptions word the copy button.
type DocCodeTabsOptions struct{ Copy string }

// DocCodeTabs shows code per language tab; cue Select.
func DocCodeTabs(c *ui.Context, samples []CodeSample, picked *int, o DocCodeTabsOptions) (root, copy ui.Element) {
	def(&o.Copy, "Copy code")
	names := make([]string, len(samples))
	for i, s := range samples {
		if strings.TrimSpace(s.Code) == "" {
			panic(fmt.Sprintf("content: DocCodeTabs sample %q needs code", s.Name))
		}
		names[i] = s.Name
	}
	unique("DocCodeTabs", names)
	if len(samples) == 0 || picked == nil || *picked < 0 || *picked >= len(samples) {
		panic("content: DocCodeTabs needs samples and a picked one")
	}
	t := c.Theme()
	s := samples[*picked]
	root = card(c).Clip().Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.ScrollHorizontal(c).Grow(1).MinWidth(0).FocusGroup(ui.Horizontal).Role(ui.RoleTabList).Children(func() {
				ui.Row(c).Padding(0, t.Space(2)).Children(func() {
					for i, smp := range samples {
						c.Key(strconv.Quote(smp.Name))
						on := i == *picked
						b := ui.ButtonBase(c).Label(smp.Name).Role(ui.RoleTab).Checked(on).Padding(t.Space(2.5), t.Space(3), 0, t.Space(3))
						b.Children(func() {
							ui.Column(c).Gap(t.Space(2)).AlignItems(ui.Stretch).Children(func() {
								ink := t.TextMuted.Mix(t.Text, unit(b.Hovered())).Mix(t.Accent, unit(on))
								ui.Text(c, smp.Name).FontSize(t.Rem(0.8)).Bold().TextColor(b.AnimateColor("ink", ink, ui.Snappy)).Shrink(0)
								ui.Box(c).FillWidth().Height(t.Space(0.5)).Radius(999).Background(b.AnimateColor("line", t.Accent.Alpha(unit(on)), ui.Snappy))
							})
						})
						if b.Clicked() && !on {
							*picked = i
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("DocCodeTabs", "picked", "index", i)
						}
					}
				})
			})
			ui.Row(c).Shrink(0).Padding(0, t.Space(2)).Children(func() {
				c.Key(strconv.Quote(s.Code))
				copy = button.CopyButton(c, o.Copy, s.Code)
			})
		})
		c.Key(strconv.Quote(s.Name))
		body := ui.Column(c).FillWidth().Padding(t.Space(3), 0)
		kit.Arrive(body, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
			h := newHighlighter(c, s.Lang)
			state := uint64(0)
			for i, line := range strings.Split(strings.TrimRight(s.Code, "\n"), "\n") {
				c.Key(i)
				var spans []ui.Span
				spans, state = h.line(state, line)
				ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.8)).LineHeight(1.6).FillWidth().Padding(0, t.Space(4)).Wrap()
			}
		})
	})
	return root, copy
}
