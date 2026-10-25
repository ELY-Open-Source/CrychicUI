package devtools

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LogLine is one log entry; Source is optional.
type LogLine struct {
	ID, Time, Message, Source string
	Level                     Level
}

// LogViewOptions word the panel; Levels names Debug to Error.
type LogViewOptions struct {
	Title, Search, Follow, Empty string
	Levels                       [4]string
	Height                       float32
}

// LogView shows leveled logs with search and follow; cue Toggle.
func LogView(c *ui.Context, lines []LogLine, query *string, follow *bool, o LogViewOptions) ui.Element {
	if query == nil || follow == nil || o.Height < 0 {
		panic(fmt.Sprintf("devtools: LogView %q needs query and follow state", o.Title))
	}
	ids := make([]string, len(lines))
	for i, l := range lines {
		if l.Time == "" || l.Message == "" {
			panic(fmt.Sprintf("devtools: LogView line %q needs a time and message", l.ID))
		}
		validLevel("LogView", l.Level)
		ids[i] = l.ID
	}
	if len(lines) > 0 {
		unique("LogView", ids)
	}
	def(&o.Title, "Logs")
	def(&o.Search, "Search logs")
	def(&o.Follow, "Follow")
	def(&o.Empty, "No matching logs")
	for i, name := range [...]string{"DEBUG", "INFO", "WARN", "ERROR"} {
		def(&o.Levels[i], name)
	}
	t := c.Theme()
	if o.Height == 0 {
		o.Height = t.Space(80)
	}
	var shown []LogLine
	for _, l := range lines {
		if lo, _ := kit.IndexFold(l.Message, *query); *query == "" || lo >= 0 {
			shown = append(shown, l)
		}
	}
	root := card(c).Role(ui.RoleGroup)
	state := ui.Local(root, "list", func() ui.ListState { return ui.ListState{FollowAnimation: &ui.Smooth} })
	state.FollowEnd = *follow
	state.Key = func(row int) any { return shown[row].ID }
	root.Children(func() {
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			glyph(c, "scroll-text", 16).TextColor(t.TextMuted)
			ui.Text(c, o.Title).FontSize(t.Rem(0.85)).Bold().Shrink(0)
			kit.Chip(c, fmt.Sprint(len(shown)), t.Accent)
			ui.Box(c).Grow(1)
			frame := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Radius(999).
				Width(t.Space(44)).MinWidth(0)
			frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
			frame.Children(func() {
				glyph(c, "search", 14).TextColor(t.TextMuted)
				in := ui.TextInputBase(c, query).Label(o.Search).Placeholder(o.Search).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0)
				if in.Changed() {
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("LogView", "searched", "runes", utf8.RuneCountInString(*query))
				}
			})
			b := kit.Press(ui.ButtonBase(c).Label(o.Follow).Padding(t.Space(1), t.Space(2.5)).Radius(999).Shrink(0).Checked(*follow))
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, unit(*follow)), ui.Snappy))
			b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(*follow)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(*follow)), ui.Snappy))
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					glyph(c, "arrow-down-to-line", 14)
					ui.Text(c, o.Follow).FontSize(t.Rem(0.75)).Bold().Shrink(0)
				})
			})
			if b.Clicked() {
				*follow = !*follow
				if *follow {
					state.ScrollToEnd()
				}
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("LogView", "follow", *follow)
			}
		})
		if len(shown) == 0 {
			c.Key("empty")
			ui.Column(c).FillWidth().Height(o.Height).Children(func() { empty(c, "search-x", o.Empty) })
			return
		}
		c.Key("lines")
		ui.List(c, state, len(shown), func(i int) {
			l := shown[i]
			ink := levelInk(t, l.Level)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Padding(t.Space(1), t.Space(3))
			if l.Level >= LevelWarn {
				row.Background(ink.Alpha(0.06))
			}
			row.Children(func() {
				kit.Muted(c, l.Time).Font("monospace").FontSize(t.Rem(0.72)).Shrink(0)
				ui.Text(c, o.Levels[l.Level]).Font("monospace").FontSize(t.Rem(0.68)).Bold().TextColor(ink).Width(t.Space(11)).Shrink(0)
				lo, hi := kit.IndexFold(l.Message, *query)
				ui.RichText(c, kit.Marked(t, l.Message, lo, hi)...).Font("monospace").FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).Wrap()
				if l.Source != "" {
					kit.Muted(c, l.Source).Font("monospace").FontSize(t.Rem(0.68)).Shrink(0).MaxWidth(t.Space(24)).SingleLine()
				}
			})
		}).FillWidth().Height(o.Height).Padding(t.Space(1), 0)
	})
	return root
}
