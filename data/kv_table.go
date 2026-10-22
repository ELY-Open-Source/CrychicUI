package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Pair is a key and its value.
type Pair struct{ Key, Value string }

// KVTableOptions tune a KVTable.
type KVTableOptions struct {
	// Copy takes the key.
	Copy string
}

// KVTable lists keys and values to copy.
func KVTable(c *ui.Context, label string, pairs []Pair, o KVTableOptions) ui.Element {
	if label == "" || len(pairs) == 0 {
		panic(fmt.Sprintf("data: KVTable %q needs pairs", label))
	}
	var keys []string
	for _, p := range pairs {
		keys = append(keys, p.Key)
	}
	unique("KVTable", label, keys)
	def(&o.Copy, "Copy %s")
	t := c.Theme()
	root := frame(c, label)
	copied := ui.Local(root, "copied", func() int { return -1 })
	at := ui.Local(root, "at", func() time.Time { return time.Time{} })
	now := c.Now()
	if *copied >= 0 {
		if left := at.Add(1500 * time.Millisecond).Sub(now); left > 0 {
			c.After(left)
		} else {
			*copied = -1
		}
	}
	root.Children(func() {
		for i, p := range pairs {
			c.Key(i)
			hovered(c, false).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
				kit.Muted(c, p.Key).Width(t.Space(36)).Shrink(0).SingleLine()
				ui.Text(c, p.Value).Grow(1).SingleLine()
				name := fmt.Sprintf(o.Copy, p.Key)
				b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
					glyph, ink := "copy", t.TextMuted
					if *copied == i {
						glyph, ink = "check", t.Success
					}
					ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(ink).ContentTransition(ui.ContentTransition{})
				})
				if b.Clicked() {
					c.WriteClipboard(p.Value)
					*copied, *at = i, now
					c.After(1500 * time.Millisecond)
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("KVTable", "copied", "index", i)
				}
			})
		}
	})
	return root
}
