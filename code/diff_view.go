package code

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// DiffViewOptions tune a DiffView.
type DiffViewOptions struct {
	// Height of the viewer; zero gives the default.
	Height         float32
	Split, Unified string
}

// DiffView is MyGo's diff viewer, split or unified.
func DiffView(c *ui.Context, st *editor.DiffState, split *bool, o DiffViewOptions) ui.Element {
	if st == nil || split == nil {
		panic("code: DiffView needs a state and a layout flag")
	}
	t := c.Theme()
	size("DiffView", &o.Height, t.Space(80))
	def(&o.Split, "Split")
	def(&o.Unified, "Unified")
	if *split {
		st.Layout = editor.DiffSplit
	} else {
		st.Layout = editor.DiffUnified
	}
	root := ui.Column(c).FillWidth().Height(o.Height).Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup)
	seg := func(name string, on bool) ui.Element {
		b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(on))
		bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
		if on {
			bg = t.Accent.Alpha(0.12)
		}
		return b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
			ui.Text(c, name).FontSize(t.Rem(0.85)).SingleLine()
		})
	}
	var one, two ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Background(t.Surface).Children(func() {
			one = seg(o.Split, *split)
			two = seg(o.Unified, !*split)
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
