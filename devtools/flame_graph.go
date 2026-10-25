package devtools

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Frame is one stack frame with its sample count.
type Frame struct {
	Name     string
	Value    int64
	Children []Frame
}

// FlameGraphOptions word the panel; Share takes name and percent.
type FlameGraphOptions struct{ Title, Samples, Share string }

// FlameGraph stacks frames by width; picks a name; cue Select.
func FlameGraph(c *ui.Context, root Frame, selected *string, o FlameGraphOptions) ui.Element {
	if selected == nil {
		panic("devtools: FlameGraph needs a selection")
	}
	check(root)
	def(&o.Title, "CPU profile")
	def(&o.Samples, "%d samples")
	def(&o.Share, "%s · %.1f%% of samples")
	t := c.Theme()
	box := card(c).Role(ui.RoleGroup)
	pick := ""
	box.Children(func() {
		kit.CardHead(c, "flame", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Samples, root.Value), t.Accent) })
		ui.Row(c).FillWidth().Padding(t.Space(3)).Children(func() { flame(c, root, 0, *selected, &pick) })
		if *selected != "" {
			var total int64
			sum(root, *selected, &total)
			c.Key(*selected)
			share := float64(total) / float64(root.Value) * 100
			kit.Arrive(kit.Muted(c, fmt.Sprintf(o.Share, *selected, share)).FontSize(t.Rem(0.78)).
				Padding(0, t.Space(3), t.Space(3), t.Space(3)), ui.Pose{Fade: 1}, 0)
		}
	})
	if pick != "" && pick != *selected {
		*selected = pick
		box.SetChanged()
		sound.Play(sound.Select)
		kit.Log("FlameGraph", "selected")
	}
	return box
}

// flame draws f over its children, sized by value.
func flame(c *ui.Context, f Frame, depth int, selected string, pick *string) {
	t := c.Theme()
	h := fnv.New32a()
	h.Write([]byte(f.Name))
	ink := t.Warning.Mix(t.Danger, float32(h.Sum32()%100)/100*0.7)
	on := f.Name == selected
	ui.Column(c).Grow(float32(f.Value)).Basis(0).MinWidth(0).Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(f.Name).FillWidth().MinWidth(0).Height(t.Space(6)).Padding(0, t.Space(1)).
			Radius(t.Radius*0.4).Border(1, t.Background).Checked(on))
		b.Background(b.AnimateColor("bg", ink.Alpha(0.3).Mix(ink.Alpha(0.5), unit(b.Hovered())).Mix(t.Accent.Alpha(0.45), unit(on)), ui.Snappy))
		kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(depth)*60*time.Millisecond)
		b.Children(func() {
			ui.Text(c, f.Name).Font("monospace").FontSize(t.Rem(0.65)).FillWidth().MinWidth(0).SingleLine()
		})
		if b.Clicked() {
			*pick = f.Name
		}
		if len(f.Children) == 0 {
			return
		}
		ui.Row(c).FillWidth().Children(func() {
			rest := f.Value
			seen := map[string]int{}
			for _, k := range f.Children {
				c.Key(fmt.Sprintf("%q#%d", k.Name, seen[k.Name]))
				seen[k.Name]++
				flame(c, k, depth+1, selected, pick)
				rest -= k.Value
			}
			if rest > 0 {
				c.Key("rest")
				ui.Box(c).Grow(float32(rest)).Basis(0)
			}
		})
	})
}

// check panics on frames that cannot be drawn.
func check(f Frame) {
	if f.Name == "" || f.Value <= 0 {
		panic(fmt.Sprintf("devtools: FlameGraph frame %q needs a name and a positive value", f.Name))
	}
	var kids int64
	for _, k := range f.Children {
		check(k)
		kids += k.Value
		if kids > f.Value || kids < 0 {
			panic(fmt.Sprintf("devtools: FlameGraph frame %q has children over its value", f.Name))
		}
	}
}

// sum adds the values of frames named name, outermost only.
func sum(f Frame, name string, total *int64) {
	if f.Name == name {
		*total += f.Value
		return
	}
	for _, k := range f.Children {
		sum(k, name, total)
	}
}
