package code

import (
	"crypto/sha256"
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// HashToolOptions name the field and the copy button.
type HashToolOptions struct {
	Text, Copy string
}

// HashTool shows the SHA-256 of text in blocks; digests cross-fade.
func HashTool(c *ui.Context, label string, text *string, o HashToolOptions) (root, copy ui.Element) {
	if label == "" || text == nil {
		panic(fmt.Sprintf("code: HashTool %q needs a label and a bound text", label))
	}
	def(&o.Text, "Text to hash")
	def(&o.Copy, "Copy")
	t := c.Theme()
	var digest string
	root = kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		kit.CardHead(c, "fingerprint", label, func() {
			typography.Bytes(c, int64(len(*text))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		})
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			if ui.TextArea(c, text).Label(o.Text).Placeholder(o.Text).Font("monospace").Lines(2, 5).FillWidth().Changed() {
				edited(root, "HashTool")
			}
			digest = fmt.Sprintf("%x", sha256.Sum256([]byte(*text)))
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(1), t.Space(2), t.Space(3)).
				Background(t.Surface).Border(1, t.Border).Radius(t.Radius).Children(func() {
				ui.Box(c).Grow(1).MinWidth(0).Children(func() {
					// Keyed by digest: old fades out, new fades in.
					c.Key(digest)
					ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Label(digest).
						Transition(ui.ElementTransition{Enter: &ui.Motion{}, Exit: &ui.Motion{}}).Children(func() {
						for i := 0; i < len(digest); i += 8 {
							c.Key(i)
							ink := t.Text
							if i/8%2 == 1 {
								ink = t.TextMuted
							}
							ui.Text(c, digest[i:i+8]).Font("monospace").FontSize(t.Rem(0.8)).TextColor(ink).Shrink(0)
						}
					})
				})
				copy = button.CopyButton(c, o.Copy, digest)
			})
		})
	})
	if copy.Clicked() {
		kit.Log("HashTool", "copied")
	}
	return root, copy
}
