package files

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TrashOptions word the trash; Deleted takes the formatted date.
type TrashOptions struct {
	Restore, Empty, Note, Deleted string
	FormatDate                    func(t time.Time) string
}

// Trash lists deleted items to restore; cue Success.
func Trash(c *ui.Context, label string, items []FileItem, restored *int, o TrashOptions) ui.Element {
	if label == "" || restored == nil || *restored < -1 || *restored >= len(items) {
		panic(fmt.Sprintf("files: Trash %q needs a label and a restore index or -1", label))
	}
	for _, f := range items {
		validFile("Trash", f)
	}
	unique("Trash", label, items)
	def(&o.Restore, "Restore")
	def(&o.Empty, "Trash is empty")
	def(&o.Note, "Items are deleted forever after 30 days")
	def(&o.Deleted, "Deleted %s")
	if o.FormatDate == nil {
		o.FormatDate = shortDate
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleList).Label(label)
	root.Children(func() {
		kit.CardHead(c, "trash-2", label, func() {
			kit.Numeric(c, strconv.Itoa(len(items))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		ui.Column(c).FillWidth().Padding(t.Space(1)).Children(func() {
			if len(items) == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(4)).Transition(kit.Enter).Children(func() {
					ui.Icon(c, icon.Must("trash-2")).Size(20, 20).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
				})
			}
			for i, f := range items {
				c.Key(fileKey(f))
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Transition(rowExit)
				kit.Stagger(c, root, row, i).Children(func() {
					ui.Box(c).Opacity(0.6).Children(func() { fileTile(c, f, t.Space(8)) })
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, f.Name).FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						meta := fmt.Sprintf(o.Deleted, o.FormatDate(f.Modified))
						if f.Path != "" {
							meta += " · " + f.Path
						}
						kit.Muted(c, meta).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
					})
					name := o.Restore + " " + f.Name
					b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(o.Restore).Padding(t.Space(1), t.Space(2.5)).Radius(t.Radius).Border(1, t.Border))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy)).Children(func() {
						ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Icon(c, icon.Must("rotate-ccw")).Size(12, 12)
							ui.Text(c, o.Restore).FontSize(t.Rem(0.8)).Shrink(0)
						})
					})
					if b.Clicked() {
						*restored = i
						root.SetChanged()
						sound.Play(sound.Success)
						kit.Log("Trash", "restored", "index", i)
					}
				})
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Background(t.Surface).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Icon(c, icon.Must("info")).Size(12, 12).TextColor(t.TextMuted)
			kit.Muted(c, o.Note).FontSize(t.Rem(0.75))
		})
	})
	return root
}
