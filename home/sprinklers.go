package home

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Zone is one sprinkler zone and its run length.
type Zone struct {
	Name    string
	Minutes int
}

// SprinklersOptions word the zones; Run and Stop take the name.
type SprinklersOptions struct{ Title, Run, Stop, Length, Left string }

// Sprinklers runs one zone with a countdown; cue Toggle.
func Sprinklers(c *ui.Context, zones []Zone, running *string, o SprinklersOptions) ui.Element {
	def(&o.Title, "Sprinklers")
	def(&o.Run, "Run %s")
	def(&o.Stop, "Stop %s")
	def(&o.Length, "%d min")
	def(&o.Left, "%s left")
	names := make([]string, len(zones))
	for i, z := range zones {
		if z.Minutes < 1 || z.Minutes > 120 {
			panic(fmt.Sprintf("home: Sprinklers zone %q needs 1 to 120 minutes", z.Name))
		}
		names[i] = z.Name
	}
	if len(zones) == 0 || running == nil || *running != "" && !slices.Contains(names, *running) {
		panic(fmt.Sprintf("home: Sprinklers %q needs zones and a running name among them", o.Title))
	}
	unique("Sprinklers", names)
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	ends := ui.Local(root, "ends", func() map[string]time.Time { return map[string]time.Time{} })
	now := c.Now()
	var left int64
	if *running == "" && len(*ends) > 0 {
		*ends = map[string]time.Time{}
	}
	if *running != "" {
		z := zones[slices.Index(names, *running)]
		end, ok := (*ends)[z.Name]
		if !ok {
			end = now.Add(time.Duration(z.Minutes) * time.Minute)
			*ends = map[string]time.Time{z.Name: end}
		}
		if left = countdown(c, end); left == 0 {
			*running, *ends = "", map[string]time.Time{}
			root.SetChanged()
			kit.Log("Sprinklers", "finished", "zone", slices.Index(names, z.Name))
		}
	}
	root.Children(func() {
		kit.CardHead(c, "sprout", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5), 0)
		list.Children(func() {
			for i, z := range zones {
				on := z.Name == *running
				c.Key(strconv.Quote(z.Name))
				row := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2.5), t.Space(3)).Role(ui.RoleListItem)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.Accent, 0.06*unit(on)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						tile(c, "droplets", t.Accent, t.Space(9), on)
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
							ui.Text(c, z.Name).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).SingleLine()
							meta := fmt.Sprintf(o.Length, z.Minutes)
							if on {
								meta = fmt.Sprintf(o.Left, fmt.Sprintf("%d:%02d", left/60, left%60))
							}
							kit.Numeric(c, meta).FontSize(t.Rem(0.72)).TextColor(t.TextMuted.Mix(t.Accent, unit(on))).FillWidth().MinWidth(0).
								SingleLine().ContentTransition(ui.ContentTransition{Numeric: true})
						})
						label, mark := o.Run, "play"
						if on {
							label, mark = o.Stop, "square"
						}
						c.Key("act")
						b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(label, z.Name)).Shrink(0).Padding(t.Space(1.5), t.Space(3)).Radius(999).
							Border(1, t.Border))
						b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
						b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
						b.Children(func() {
							c.Key(mark)
							kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.6, Fade: 1}, 0)
						})
						if b.Clicked() {
							if on {
								*running, *ends = "", map[string]time.Time{}
							} else {
								*running = z.Name
								*ends = map[string]time.Time{z.Name: now.Add(time.Duration(z.Minutes) * time.Minute)}
							}
							root.SetChanged()
							sound.Play(sound.Toggle)
							kit.Log("Sprinklers", "run", "zone", i, "on", !on)
						}
					})
					if on {
						c.Key("bar")
						share := float32(left) / float32(z.Minutes*60)
						meter(c, share, t.Accent)
					}
				})
			}
		})
	})
	return root
}
