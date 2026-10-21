package typography

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var keyPressFrames = []ui.Keyframe{
	{Value: 0},
	{At: 60 * time.Millisecond, Value: 1, Ease: ui.EaseOut},
	{At: 160 * time.Millisecond, Value: 0, Ease: ui.EaseOut},
}

// keyCap draws one cap and replays its press.
func keyCap(c *ui.Context, label string, run uint64, delay time.Duration) ui.Element {
	t := c.Theme()
	cap := ui.Row(c).AlignSelf(ui.Start).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Background(t.Surface).
		BorderWidth(1, 1, 2, 1).BorderColor(t.Border).Font("monospace").FontSize(t.Rem(0.85))
	if run > 0 {
		p := cap.Keyframes("press", keyPressFrames, ui.KeyframeOptions{Run: run, Delay: delay})
		cap.Translate(0, 2*p)
	}
	return cap.Children(func() { ui.Text(c, label) })
}

// pressRun advances when pressed rises.
func pressRun(e ui.Element, pressed bool) uint64 {
	run := ui.Local(e, "run", func() uint64 { return 0 })
	was := ui.Local(e, "pressed", func() bool { return false })
	if pressed && !*was {
		(*run)++
	}
	*was = pressed
	return *run
}

// Kbd is a key cap pressing on caller matches.
func Kbd(c *ui.Context, label string, pressed bool) ui.Element {
	root := ui.Box(c)
	run := pressRun(root, pressed)
	return root.Children(func() { keyCap(c, label, run, 0) })
}
