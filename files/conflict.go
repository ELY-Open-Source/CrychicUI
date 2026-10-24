package files

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConflictChoice is a conflict resolution.
type ConflictChoice uint8

// Conflict resolutions.
const (
	ConflictNone ConflictChoice = iota
	KeepMine
	KeepTheirs
	KeepBoth
)

// ConflictOptions word the two copies and the choices.
type ConflictOptions struct {
	Yours, Theirs, Mine, KeepTheirs, Both string
	FormatDate                            func(t time.Time) string
}

// Conflict compares two copies, keeping one or both; cue Select.
func Conflict(c *ui.Context, label string, mine, theirs FileItem, choice *ConflictChoice, o ConflictOptions) ui.Element {
	validFile("Conflict", mine)
	validFile("Conflict", theirs)
	if label == "" || choice == nil || *choice > KeepBoth {
		panic(fmt.Sprintf("files: Conflict %q needs a label and a known choice", label))
	}
	def(&o.Yours, "Your copy")
	def(&o.Theirs, "Their copy")
	def(&o.Mine, "Keep mine")
	def(&o.KeepTheirs, "Keep theirs")
	def(&o.Both, "Keep both")
	if o.FormatDate == nil {
		o.FormatDate = func(t time.Time) string { return t.Format("Jan 2, 15:04") }
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "git-compare", label, nil)
		ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			copyCard := func(key, name string, f FileItem, on bool) {
				c.Key(key)
				card := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1.5)).Padding(t.Space(2.5)).Radius(t.Radius)
				card.Background(card.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
				card.Border(1, card.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				card.Children(func() {
					kit.Muted(c, name).FontSize(t.Rem(0.7)).Bold()
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						fileTile(c, f, t.Space(7))
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							ui.Text(c, f.Name).FontSize(t.Rem(0.8)).Bold().FillWidth().SingleLine()
							ui.Row(c).Gap(t.Space(1)).Children(func() {
								sizeText(c, f.Size).FontSize(t.Rem(0.7))
								kit.Muted(c, "·").FontSize(t.Rem(0.7))
								kit.Muted(c, o.FormatDate(f.Modified)).FontSize(t.Rem(0.7)).Shrink(0)
							})
						})
					})
				})
			}
			copyCard("mine", o.Yours, mine, *choice == KeepMine || *choice == KeepBoth)
			copyCard("theirs", o.Theirs, theirs, *choice == KeepTheirs || *choice == KeepBoth)
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(0, t.Space(3), t.Space(3), t.Space(3)).FocusGroup(ui.Horizontal).Role(ui.RoleRadioGroup).Children(func() {
			for _, it := range []struct {
				choice      ConflictChoice
				name, glyph string
			}{{KeepMine, o.Mine, "user"}, {KeepTheirs, o.KeepTheirs, "cloud"}, {KeepBoth, o.Both, "copy"}} {
				c.Key(it.choice)
				on := *choice == it.choice
				b := kit.Press(ui.ButtonBase(c).Label(it.name).Grow(1).Basis(0).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on).Role(ui.RoleRadio))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icon.Must(it.glyph)).Size(14, 14).Shrink(0)
						ui.Text(c, it.name).FontSize(t.Rem(0.8)).Bold().Shrink(0)
					})
				})
				if b.Clicked() && !on {
					*choice = it.choice
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Conflict", "resolved", "choice", int(it.choice))
				}
			}
		})
	})
	return root
}
