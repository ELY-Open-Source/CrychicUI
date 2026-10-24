package files

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StarredOptions name the star toggles.
type StarredOptions struct{ Star, Unstar string }

// Starred lists files with popping star toggles; cue Toggle.
func Starred(c *ui.Context, label string, list []FileItem, starred *[]bool, o StarredOptions) ui.Element {
	validFiles("Starred", label, list)
	unique("Starred", label, list)
	if starred == nil || len(*starred) != len(list) {
		panic(fmt.Sprintf("files: Starred %q needs one flag per file", label))
	}
	def(&o.Star, "Star")
	def(&o.Unstar, "Unstar")
	t := c.Theme()
	count := 0
	for _, s := range *starred {
		if s {
			count++
		}
	}
	root := kit.Card(c).FillWidth().Role(ui.RoleList).Label(label)
	root.Children(func() {
		kit.CardHead(c, "star", label, func() {
			kit.Numeric(c, strconv.Itoa(count)).FontSize(t.Rem(0.75)).TextColor(t.Warning).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
		})
		ui.Column(c).FillWidth().Padding(t.Space(1)).Children(func() {
			for i, f := range list {
				c.Key(fileKey(f))
				on := (*starred)[i]
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(2)).Radius(t.Radius)
				kit.Stagger(c, root, row, i).Children(func() {
					fileTile(c, f, t.Space(7))
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, f.Name).FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						if f.Path != "" {
							kit.Muted(c, f.Path).FontSize(t.Rem(0.7)).FillWidth().SingleLine()
						}
					})
					name := o.Star + " " + f.Name
					if on {
						name = o.Unstar + " " + f.Name
					}
					b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Checked(on))
					b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Warning.Alpha(0.14), unit(on)), ui.Snappy))
					b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Warning, unit(on)), ui.Snappy))
					pop := b.AnimateTo("pop", unit(on), ui.Bouncy)
					b.Children(func() {
						star := ui.Icon(c, icon.Must("star")).Size(16, 16)
						star.Scale(1+0.15*pop*(1-pop)*4, 1+0.15*pop*(1-pop)*4)
					})
					if b.Clicked() {
						(*starred)[i] = !on
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("Starred", "toggled", "index", i, "starred", !on)
					}
				})
			}
		})
	})
	return root
}
