package code

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// DiffViewOptions tune a DiffView.
type DiffViewOptions struct {
	// Height of the viewer; zero gives the default.
	Height float32
	// Name is the file compared.
	Name           string
	Split, Unified string
}

// DiffView is MyGo's diff viewer under a file bar.
func DiffView(c *ui.Context, st *editor.DiffState, split *bool, o DiffViewOptions) ui.Element {
	if st == nil || split == nil {
		panic("code: DiffView needs a state and a layout flag")
	}
	t := c.Theme()
	size("DiffView", &o.Height, t.Space(80))
	def(&o.Name, "Changes")
	def(&o.Split, "Split")
	def(&o.Unified, "Unified")
	st.Layout = editor.DiffUnified
	if *split {
		st.Layout = editor.DiffSplit
	}
	added, deleted := st.Stats()
	root := ui.Column(c).FillWidth().Height(o.Height).Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup).Label(o.Name)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	var one, two ui.Element
	root.Children(func() {
		codeHeader(c, o.Name, "", func() {
			kit.Numeric(c, "+"+strconv.Itoa(added)).FontSize(t.Rem(0.75)).Bold().TextColor(t.Success)
			kit.Numeric(c, "−"+strconv.Itoa(deleted)).FontSize(t.Rem(0.75)).Bold().TextColor(t.Danger)
			track := ui.Row(c).Padding(t.Space(0.5)).Radius(999).Background(t.Surface).Border(1, t.Border).FocusGroup(ui.Horizontal)
			x := track.AnimateTo("thumb", unit(!*split), ui.Snappy)
			track.Draw(func(p *ui.Painter, r ui.Rect) {
				pad := t.Space(0.5) + 1
				w := (r.W - 2*pad) / 2
				p.Fill(ui.Rect{X: r.X + pad + w*x, Y: r.Y + pad, W: w, H: r.H - 2*pad}, t.Background, 999)
			})
			track.Children(func() {
				one = diffSegment(c, o.Split, "columns-2", *split)
				two = diffSegment(c, o.Unified, "rows-2", !*split)
			})
		})
		editor.DiffView(c, st).Grow(1).FillWidth()
	})
	pick := func(to bool, name string) {
		if *split == to {
			return
		}
		*split = to
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("DiffView", "layout", name)
	}
	if one.Clicked() {
		pick(true, "split")
	}
	if two.Clicked() {
		pick(false, "unified")
	}
	return root
}

// diffSegment is one half of the layout switch.
func diffSegment(c *ui.Context, name, glyph string, on bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(name).Width(t.Space(22)).Padding(t.Space(1), t.Space(2)).Radius(999).Checked(on))
	b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(on || b.Hovered())), ui.Snappy))
	return b.Children(func() {
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(14, 14).Shrink(0)
			ui.Text(c, name).FontSize(t.Rem(0.8)).Shrink(0)
		})
	})
}
