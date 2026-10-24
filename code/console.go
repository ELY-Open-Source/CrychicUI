package code

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LogLevel is a console level.
type LogLevel uint8

// Console levels.
const (
	LevelLog LogLevel = iota
	LevelWarn
	LevelError
	logLevelCount
)

// LogEntry is one console row and where it came from.
type LogEntry struct {
	Level        LogLevel
	Text, Source string
}

// ConsoleOptions are the copy knobs of Console.
type ConsoleOptions struct {
	Log, Warn, Error, Empty string
}

// Console shows log rows filtered by level; cue Toggle.
func Console(c *ui.Context, label string, entries []LogEntry, filter *int, o ConsoleOptions) ui.Element {
	if label == "" || filter == nil || *filter < 0 || *filter > 7 {
		panic(fmt.Sprintf("code: Console %q needs a label and a filter from 0 to 7", label))
	}
	var counts [logLevelCount]int
	for _, e := range entries {
		if e.Level >= logLevelCount || e.Text == "" {
			panic(fmt.Sprintf("code: Console %q needs a known level and text in every entry", label))
		}
		counts[e.Level]++
	}
	def(&o.Log, "Log")
	def(&o.Warn, "Warn")
	def(&o.Error, "Error")
	def(&o.Empty, "No messages")
	t := c.Theme()
	names := [logLevelCount]string{o.Log, o.Warn, o.Error}
	glyphs := [logLevelCount]string{"info", "triangle-alert", "circle-x"}
	inks := [logLevelCount]ui.Color{t.TextMuted, t.Warning, t.Danger}
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		kit.CardHead(c, "square-terminal", label, func() {
			for lv := LogLevel(0); lv < logLevelCount; lv++ {
				c.Key(lv)
				on := *filter&(1<<lv) != 0
				name := names[lv] + " " + strconv.Itoa(counts[lv])
				b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(0.5), t.Space(2)).Radius(999).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(inks[lv].Alpha(0.14), unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(inks[lv].Mix(t.Text, 0.3), unit(on)), ui.Snappy))
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icon.Must(glyphs[lv])).Size(12, 12).Shrink(0)
						ui.Text(c, names[lv]).FontSize(t.Rem(0.75)).Shrink(0)
						kit.Numeric(c, strconv.Itoa(counts[lv])).FontSize(t.Rem(0.7)).Bold().Shrink(0)
					})
				})
				if b.Clicked() {
					*filter ^= 1 << lv
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Console", "filter", "level", int(lv), "on", !on)
				}
			}
		})
		ui.Column(c).FillWidth().Padding(t.Space(1), 0).Children(func() {
			shown := 0
			for i, e := range entries {
				if *filter&(1<<e.Level) == 0 {
					continue
				}
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(3)).
					Background(inks[e.Level].Alpha(0.07 * unit(e.Level != LevelLog))).Transition(rowExit)
				kit.Stagger(c, root, row, shown).Children(func() {
					ui.Icon(c, icon.Must(glyphs[e.Level])).Size(14, 14).TextColor(inks[e.Level]).Shrink(0)
					txt := ui.Text(c, e.Text).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
					if e.Level == LevelError {
						txt.TextColor(t.Danger)
					}
					if e.Source != "" {
						kit.Muted(c, e.Source).Font("monospace").FontSize(t.Rem(0.7)).Shrink(0)
					}
				})
				shown++
			}
			if shown == 0 {
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.85)).Padding(t.Space(4)).AlignSelf(ui.Center).Transition(kit.Enter)
			}
		})
	})
	return root
}
