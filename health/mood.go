package health

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MoodDay is a past day's mood: 0–4, or -1 unrecorded.
type MoodDay struct {
	Label string
	Mood  int
}

// MoodOptions word the card; Names run from awful to great.
type MoodOptions struct {
	Title, Prompt string
	Names         [5]string
}

var moodFaces = [5]string{"angry", "frown", "meh", "smile", "laugh"}

// Mood picks today's mood below a week; cue Select.
func Mood(c *ui.Context, week []MoodDay, mood *int, o MoodOptions) ui.Element {
	def(&o.Title, "Mood")
	if mood == nil || *mood < -1 || *mood > 4 || len(week) > 7 {
		panic(fmt.Sprintf("health: Mood %q needs a mood of -1 to 4 and at most 7 past days", o.Title))
	}
	labels := make([]string, len(week))
	for i, d := range week {
		if d.Mood < -1 || d.Mood > 4 {
			panic(fmt.Sprintf("health: Mood day %q has mood %d", d.Label, d.Mood))
		}
		labels[i] = d.Label
	}
	unique("Mood", labels)
	def(&o.Prompt, "How are you today?")
	for i, n := range [...]string{"Awful", "Low", "Okay", "Good", "Great"} {
		def(&o.Names[i], n)
	}
	t := c.Theme()
	inks := [5]ui.Color{t.Danger, t.Warning, t.TextMuted, t.Accent, t.Success}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Title).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine()
		if len(week) > 0 {
			days := ui.Row(c).FillWidth().Gap(t.Space(1)).Justify(ui.SpaceBetween)
			days.Children(func() {
				for i, d := range week {
					col := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center)
					kit.Stagger(c, days, col, i).Children(func() {
						if d.Mood < 0 {
							ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Radius(999).Border(1, t.Border)
						} else {
							glyph(c, moodFaces[d.Mood], 24).TextColor(inks[d.Mood])
						}
						kit.Muted(c, d.Label).FontSize(t.Rem(0.62)).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
					})
				}
			})
		}
		kit.Muted(c, o.Prompt).FontSize(t.Rem(0.8)).FillWidth().Wrap()
		pick := ui.Row(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleRadioGroup)
		pick.Children(func() {
			for i, name := range o.Names {
				on := *mood == i
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(name).Role(ui.RoleRadio).Checked(on).Grow(1).Basis(0).MinWidth(0).Column().Center().
					Gap(t.Space(1)).Padding(t.Space(2), 0).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(inks[i].Alpha(0.14), unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(inks[i], unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(inks[i], unit(on)), ui.Snappy))
				b.Children(func() {
					g := glyph(c, moodFaces[i], 26)
					k := g.AnimateTo("pop", 1+0.2*unit(on), ui.Bouncy)
					g.Scale(k, k)
					ui.Text(c, name).FontSize(t.Rem(0.66)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
				})
				if b.Clicked() && !on {
					*mood = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Mood", "picked", "mood", i)
				}
			}
		})
	})
	return root
}
