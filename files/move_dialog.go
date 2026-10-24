package files

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MoveDialogOptions word the dialog.
type MoveDialogOptions struct {
	Title, Move string
	// To words the chosen target, "%s" taking its path.
	To string
}

// MoveDialog chooses a target folder from a tree; cue Drop.
func MoveDialog(c *ui.Context, label string, folders []string, target *int, o MoveDialogOptions) (root, move ui.Element) {
	if label == "" || len(folders) == 0 || target == nil || *target < -1 || *target >= len(folders) {
		panic(fmt.Sprintf("files: MoveDialog %q needs folders and a target or -1", label))
	}
	seen := map[string]bool{}
	for _, f := range folders {
		if f == "" || seen[f] || strings.HasPrefix(f, "/") || strings.HasSuffix(f, "/") {
			panic(fmt.Sprintf("files: MoveDialog %q needs unique relative folder paths", label))
		}
		seen[f] = true
	}
	def(&o.Title, "Choose a destination")
	def(&o.Move, "Move")
	def(&o.To, "Move to %s")
	t := c.Theme()
	root = kit.Card(c).FillWidth().Role(ui.RoleDialog).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Scale: 0.98, Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "folder-input", o.Title, nil)
		ui.Column(c).FillWidth().Padding(t.Space(1)).FocusGroup(ui.Vertical).Role(ui.RoleTree).Children(func() {
			for i, folder := range folders {
				c.Key(folder)
				depth := strings.Count(folder, "/")
				name := folder[strings.LastIndex(folder, "/")+1:]
				on := i == *target
				row := kit.Press(ui.ButtonBase(c).Label(folder).FillWidth().AlignItems(ui.Stretch).Radius(t.Radius).Checked(on).
					Padding(t.Space(1.25), t.Space(2), t.Space(1.25), t.Space(2)+float32(depth)*t.Space(5)).Role(ui.RoleTreeItem))
				kit.Stagger(c, root, row, i)
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				row.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						glyph, ink := "folder", t.Warning
						if on {
							glyph, ink = "folder-open", t.Accent
						}
						ui.Icon(c, icon.Must(glyph)).Size(16, 16).Shrink(0).TextColor(ink).ContentTransition(ui.ContentTransition{})
						txt := ui.Text(c, name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						if on {
							txt.Bold()
						}
					})
				})
				if row.Clicked() && !on {
					*target = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("MoveDialog", "target", "index", i)
				}
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Background(t.Surface).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			to := ""
			if *target >= 0 {
				to = fmt.Sprintf(o.To, folders[*target])
			}
			kit.Muted(c, to).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
			c.Key("move")
			move = kit.Press(ui.ButtonBase(c).Label(o.Move).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)).Disabled(*target < 0)
			move.Background(move.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(move.Hovered())), ui.Snappy)).TextColor(t.AccentText).
				Children(func() { ui.Text(c, o.Move).Bold().Shrink(0) })
		})
	})
	if move.Clicked() && *target >= 0 {
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("MoveDialog", "moved", "target", *target)
	}
	return root, move
}
