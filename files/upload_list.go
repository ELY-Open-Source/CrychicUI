package files

import (
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UploadItem is one upload, its size and completion fraction.
type UploadItem struct {
	Name     string
	Size     int64
	Progress float64
}

// UploadListOptions word the uploads.
type UploadListOptions struct {
	Cancel, Empty, Of, Done string
}

// uploaded is the bytes sent, exact when done, never overflowing.
func uploaded(u UploadItem) int64 {
	f := float64(u.Size) * u.Progress
	if f >= float64(u.Size) {
		return u.Size
	}
	return int64(f)
}

// UploadList shows upload bars with cancel; cue Close.
func UploadList(c *ui.Context, label string, uploads *[]UploadItem, o UploadListOptions) ui.Element {
	if label == "" || uploads == nil {
		panic(fmt.Sprintf("files: UploadList %q needs a label and uploads", label))
	}
	seen := map[string]bool{}
	total := 0.0
	for _, u := range *uploads {
		if u.Name == "" || seen[u.Name] || u.Size < 0 || !(u.Progress >= 0 && u.Progress <= 1) || math.IsInf(u.Progress, 0) {
			panic(fmt.Sprintf("files: UploadList %q upload %q needs a unique name, a size and progress in [0,1]", label, u.Name))
		}
		seen[u.Name] = true
		total += u.Progress
	}
	def(&o.Cancel, "Cancel")
	def(&o.Empty, "No uploads")
	def(&o.Of, "of")
	def(&o.Done, "Done")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleList).Label(label)
	cancel := -1
	root.Children(func() {
		kit.CardHead(c, "cloud-upload", label, func() {
			if n := len(*uploads); n > 0 {
				kit.Numeric(c, strconv.Itoa(int(total/float64(n)*100+0.5))+"%").FontSize(t.Rem(0.75)).Bold().TextColor(t.Accent).
					ContentTransition(ui.ContentTransition{Numeric: true})
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).Children(func() {
			if len(*uploads) == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(3)).Transition(kit.Enter).Children(func() {
					ui.Icon(c, icon.Must("cloud-upload")).Size(20, 20).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
				})
			}
			for i, u := range *uploads {
				c.Key(u.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5)).Radius(t.Radius).Transition(rowExit)
				kit.Stagger(c, root, row, i).Children(func() {
					fileTile(c, FileItem{Name: u.Name}, t.Space(8))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, u.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							kit.Numeric(c, strconv.Itoa(int(u.Progress*100))+"%").FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Shrink(0)
						})
						ink := t.Accent
						if u.Progress == 1 {
							ink = t.Success
						}
						progressBar(c, u.Progress, ink)
						ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
							sizeText(c, uploaded(u)).FontSize(t.Rem(0.7))
							kit.Muted(c, o.Of).FontSize(t.Rem(0.7))
							sizeText(c, u.Size).FontSize(t.Rem(0.7))
						})
					})
					if u.Progress == 1 {
						ui.Icon(c, icon.Must("circle-check")).Size(18, 18).Shrink(0).TextColor(t.Success).Label(o.Done + " " + u.Name)
						return
					}
					name := o.Cancel + " " + u.Name
					b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999))
					b.Background(b.AnimateColor("bg", t.Danger.Alpha(0.12*unit(b.Hovered())), ui.Snappy)).
						TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(b.Hovered())), ui.Snappy)).
						Children(func() { ui.Icon(c, icon.Must("x")).Size(14, 14) })
					if b.Clicked() {
						cancel = i
					}
				})
			}
		})
	})
	if cancel >= 0 {
		*uploads = slices.Delete(slices.Clone(*uploads), cancel, cancel+1)
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("UploadList", "canceled", "left", len(*uploads))
	}
	return root
}
