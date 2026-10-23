// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"strconv"
)

type LogLevel uint8

const (
	LevelLog LogLevel = iota
	LevelWarn
	LevelError
	logLevelCount
)

type LogEntry struct {
	Level LogLevel
	Text  string
}

// ConsoleOptions are the copy knobs of Console.

type ConsoleOptions struct {
	Log   string
	Warn  string
	Error string
}

// Console is log rows with a filter; cue Toggle.

func Console(c *ui.Context, label string, entries []LogEntry, filter *int, o ConsoleOptions) ui.Element {
	if label == "" || filter == nil || *filter < 0 || *filter > 7 {
		panic("code: Console needs a label and a filter from 0 to 7")
	}
	def(&o.Log, "Log")
	def(&o.Warn, "Warn")
	def(&o.Error, "Error")
	t := c.Theme()
	names := []string{o.Log, o.Warn, o.Error}
	glyphs := []string{"info", "triangle-alert", "circle-x"}
	inks := []ui.Color{t.TextMuted, t.Warning, t.Danger}
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(90)).Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		counts := [3]int{}
		for _, e := range entries {
			if e.Level >= logLevelCount || e.Text == "" {
				panic("code: Console needs a known level and text in every entry")
			}
			counts[e.Level]++
		}
		ui.Row(c).Gap(t.Space(1)).Children(func() {
			for lv := 0; lv < 3; lv++ {
				lv := lv
				on := *filter&(1<<lv) != 0
				name := names[lv]
				b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(0.5), t.Space(1.5)).Radius(999).Checked(on))
				bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
				if on {
					bg = inks[lv].Alpha(0.15)
				}
				b.Background(b.AnimateColor("bg", bg, ui.Snappy))
				b.Children(func() {
					ui.Text(c, name+" "+strconv.Itoa(counts[lv])).FontSize(t.Rem(0.75)).SingleLine()
				})
				if b.Clicked() {
					*filter ^= 1 << lv
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Console", "filter", "level", lv)
				}
			}
		})
		list := ui.Column(c).FillWidth()
		list.Children(func() {
			for i, e := range entries {
				if *filter&(1<<e.Level) == 0 {
					continue
				}
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(1))
				row.Children(func() {
					ui.Icon(c, icon.Must(glyphs[e.Level])).Size(14, 14).TextColor(inks[e.Level]).Shrink(0)
					ui.Text(c, e.Text).Font("monospace").FontSize(t.Rem(0.85)).SingleLine()
				})
			}
		})
	})
	return root
}
