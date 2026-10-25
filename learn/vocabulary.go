package learn

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Word is a term; Tone is Hz, 0 for none.
type Word struct {
	Term, Reading, Meaning, Example string
	Tone                            float64
}

// VocabularyOptions word the list; Played, if set, counts tones.
type VocabularyOptions struct {
	Title, Count, Listen string
	Played               *int
}

// Vocabulary lists words; a row opens its example; cue Tap.
func Vocabulary(c *ui.Context, words []Word, open *int, o VocabularyOptions) ui.Element {
	def(&o.Title, "Vocabulary")
	if open == nil || len(words) == 0 || *open < -1 || *open >= len(words) || o.Played != nil && *o.Played < 0 {
		panic(fmt.Sprintf("learn: Vocabulary %q needs words and an open row in range", o.Title))
	}
	terms := make([]string, len(words))
	for i, w := range words {
		if strings.TrimSpace(w.Meaning) == "" || w.Tone != 0 && !(w.Tone >= 20 && w.Tone <= 20000) {
			panic(fmt.Sprintf("learn: Vocabulary %q needs a meaning and a tone of 0 or 20–20000 Hz", w.Term))
		}
		terms[i] = w.Term
	}
	unique("Vocabulary", terms)
	def(&o.Count, "%d words")
	def(&o.Listen, "Listen to %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	played := ui.Local(root, "played", func() map[string]time.Time { return map[string]time.Time{} })
	root.Children(func() {
		kit.CardHead(c, "languages", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Count, len(words)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2))
		list.Children(func() {
			for i, w := range words {
				on := *open == i
				c.Key(strconv.Quote(w.Term))
				item := ui.Column(c).FillWidth().Radius(t.Radius)
				item.Background(item.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.07), unit(on)), ui.Snappy))
				kit.Stagger(c, list, item, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						b := kit.Press(ui.ButtonBase(c).Label(w.Term).Grow(1).MinWidth(0).Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).
							AlignItems(ui.Stretch).Checked(on))
						b.Background(b.AnimateColor("bg", t.Background.Alpha(0).Mix(t.SurfaceHover, unit(b.Hovered() && !on)), ui.Snappy))
						b.Children(func() {
							ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
								ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
									ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
										ui.Text(c, w.Term).FontSize(t.Rem(0.95)).Bold().Shrink(0).MaxWidth(t.Space(36)).SingleLine()
										if w.Reading != "" {
											kit.Muted(c, w.Reading).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
										}
									})
									kit.Muted(c, w.Meaning).FontSize(t.Rem(0.78)).FillWidth().Wrap()
								})
								chev := glyph(c, "chevron-down", 16).TextColor(t.TextMuted)
								chev.Rotate(chev.AnimateTo("turn", 180*unit(on), ui.Snappy))
							})
						})
						if b.Clicked() {
							*open = i
							if on {
								*open = -1
							}
							root.SetChanged()
							sound.Play(sound.Tap)
							kit.Log("Vocabulary", "opened", "row", *open)
						}
						if w.Tone == 0 {
							return
						}
						c.Key("listen")
						s := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Listen, w.Term)).Size(t.Space(9), t.Space(9)).Center().Radius(999).Shrink(0))
						fresh := c.Now().Sub((*played)[w.Term]) < 300*time.Millisecond
						if fresh {
							c.After(300*time.Millisecond - c.Now().Sub((*played)[w.Term]))
						}
						s.Background(s.AnimateColor("bg", t.Accent.Alpha(0.1).Mix(t.Accent.Alpha(0.22), unit(s.Hovered() || fresh)), ui.Snappy))
						s.TextColor(t.Accent)
						s.Children(func() {
							g := glyph(c, "volume-2", 16)
							k := g.AnimateTo("bump", 1+0.25*unit(fresh), ui.Bouncy)
							g.Scale(k, k)
						})
						if s.Clicked() {
							(*played)[w.Term] = c.Now()
							sound.PlayNote(w.Tone)
							if o.Played != nil && *o.Played < math.MaxInt {
								*o.Played++
								root.SetChanged()
							}
							kit.Log("Vocabulary", "listened", "row", i)
						}
					})
					layout.Collapse(c, on && w.Example != "", func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(0, t.Space(3), t.Space(2.5), t.Space(3)).Children(func() {
							ui.Box(c).Width(t.Space(0.75)).FillHeight().Radius(999).Background(t.Accent.Alpha(0.5))
							ui.Text(c, w.Example).FontSize(t.Rem(0.78)).TextColor(t.Text.Mix(t.TextMuted, 0.3)).Grow(1).MinWidth(0).Wrap()
						})
					})
				})
			}
		})
	})
	return root
}
