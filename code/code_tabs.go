package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CodeFile is one tab of a CodeTabs view.
type CodeFile struct {
	Name, Lang, Code string
}

// CodeTabs shows files as tabs over numbered code.
func CodeTabs(c *ui.Context, files []CodeFile, current *int) ui.Element {
	if len(files) == 0 || current == nil || *current < 0 || *current >= len(files) {
		panic(fmt.Sprintf("code: CodeTabs needs files and a tab, got %d of %d", deref(current), len(files)))
	}
	seen := map[string]bool{}
	for _, f := range files {
		if f.Name == "" || f.Code == "" || seen[f.Name] {
			panic(fmt.Sprintf("code: CodeTabs file %q needs a unique name and code", f.Name))
		}
		seen[f.Name] = true
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Clip().Role(ui.RoleGroup)
	root.Children(func() {
		// The strip scrolls: tabs keep whole names.
		ui.ScrollHorizontal(c).FillWidth().Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).
			FocusGroup(ui.Horizontal).Role(ui.RoleTabList).Children(func() {
			for i, f := range files {
				c.Key(f.Name)
				on := i == *current
				b := kit.Press(ui.ButtonBase(c).Label(f.Name).Padding(t.Space(2), t.Space(3)).Checked(on).Role(ui.RoleTab))
				bg := t.SurfaceHover.Mix(t.Surface, unit(on))
				b.Background(b.AnimateColor("bg", bg, ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(on || b.Hovered())), ui.Snappy))
				line := b.AnimateTo("line", unit(on), ui.Snappy)
				b.Draw(func(p *ui.Painter, r ui.Rect) {
					if line > 0 {
						w := r.W * line
						p.Fill(ui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + r.H - 2, W: w, H: 2}, t.Accent, 0)
					}
				})
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icon.Must(glyphFor(f.Name))).Size(14, 14).Shrink(0)
						ui.Text(c, f.Name).Font("monospace").FontSize(t.Rem(0.8)).Shrink(0)
					})
				})
				if b.Clicked() && !on {
					*current = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("CodeTabs", "selected", "index", i)
				}
			}
		})
		f := files[*current]
		body := ui.Column(c).FillWidth().Key("body." + f.Name)
		kit.Arrive(body, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() { codeBody(c, f.Lang, f.Code, true, nil, nil) })
	})
	return root
}

func deref(p *int) int {
	if p == nil {
		return -1
	}
	return *p
}
