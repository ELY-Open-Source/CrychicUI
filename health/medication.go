package health

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Dose is one scheduled dose; At is its time label.
type Dose struct {
	Name, Amount, At string
	Taken            bool
}

// MedicationOptions word the list; Mark takes name and time.
type MedicationOptions struct{ Title, Count, Mark, Done string }

// Medication lists doses with taken marks; cue Toggle.
func Medication(c *ui.Context, doses *[]Dose, o MedicationOptions) ui.Element {
	def(&o.Title, "Medication")
	if doses == nil || len(*doses) == 0 {
		panic(fmt.Sprintf("health: Medication %q needs doses", o.Title))
	}
	keys := make([]string, len(*doses))
	taken := 0
	for i, d := range *doses {
		if strings.TrimSpace(d.Name) == "" || strings.TrimSpace(d.At) == "" {
			panic(fmt.Sprintf("health: Medication dose %d needs a name and a time", i))
		}
		keys[i] = strconv.Quote(d.Name) + " " + strconv.Quote(d.At)
		if d.Taken {
			taken++
		}
	}
	unique("Medication", keys)
	def(&o.Count, "%d of %d taken")
	def(&o.Mark, "%s at %s")
	def(&o.Done, "All taken")
	t := c.Theme()
	all := taken == len(*doses)
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "pill", o.Title, func() {
			c.Key(all)
			if all {
				kit.Arrive(kit.Chip(c, o.Done, t.Success), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			} else {
				kit.Chip(c, fmt.Sprintf(o.Count, taken, len(*doses)), t.Accent)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2))
		list.Children(func() {
			for i, d := range *doses {
				c.Key(keys[i])
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.Success.Alpha(0.06), unit(d.Taken)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					tile(c, "pill", t.Accent.Mix(t.TextMuted, unit(d.Taken)), t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						name := ui.Text(c, d.Name).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).SingleLine()
						if d.Taken {
							name.TextColor(t.TextMuted)
						}
						if d.Amount != "" {
							kit.Muted(c, d.Amount).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
						}
					})
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Shrink(0).TextColor(t.TextMuted).Children(func() {
						glyph(c, "clock", 12)
						ui.Text(c, d.At).FontSize(t.Rem(0.72)).Shrink(0).MaxWidth(t.Space(18)).SingleLine()
					})
					c.Key("mark")
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Mark, d.Name, d.At)).Checked(d.Taken).Size(t.Space(8), t.Space(8)).
						Shrink(0).Center().Radius(999))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Success, unit(d.Taken)), ui.Snappy))
					b.Border(t.Space(0.5), b.AnimateColor("line", t.Border.Mix(t.Success, unit(d.Taken)), ui.Snappy)).TextColor(t.AccentText)
					b.Children(func() {
						if d.Taken {
							c.Key("check")
							kit.Arrive(glyph(c, "check", 16), ui.Pose{Scale: 0.3, Fade: 1}, 0)
						}
					})
					if b.Clicked() {
						(*doses)[i].Taken = !d.Taken
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("Medication", "marked", "dose", i, "taken", !d.Taken)
					}
				})
			}
		})
	})
	return root
}
