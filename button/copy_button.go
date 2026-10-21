package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CopyButton writes text and briefly shows a check.
func CopyButton(c *ui.Context, label, text string) ui.Element {
	t := c.Theme()
	b := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	copied := ui.Local(b, "copied", func() bool { return false })
	until := ui.Local(b, "until", func() time.Time { return time.Time{} })
	if *copied {
		if !c.Now().Before(*until) {
			*copied = false
		} else {
			c.After(until.Sub(c.Now()))
		}
	}
	if b.Clicked() {
		c.WriteClipboard(text)
		*copied = true
		*until = c.Now().Add(2 * time.Second)
		c.After(2 * time.Second)
		sound.Play(sound.Pop)
		kit.Log("CopyButton", "copied")
	}
	b.Children(func() {
		name := "copy"
		if *copied {
			name = "check"
		}
		ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
		ui.Text(c, label).SingleLine()
	})
	return b
}
