package account

import (
	"crypto/sha256"
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BackupCodesOptions are the backup code labels.
type BackupCodesOptions struct{ Copied string }

// BackupCodes copies recovery codes; cue Pop.
func BackupCodes(c *ui.Context, codes []string, copied *int, o BackupCodesOptions) ui.Element {
	if len(codes) == 0 || copied == nil || *copied < -1 || *copied >= len(codes) {
		panic("account: BackupCodes needs codes and state")
	}
	def(&o.Copied, "Copied")
	t := c.Theme()
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	seen := map[[sha256.Size]byte]int{}
	root.Children(func() {
		for i, code := range codes {
			if code == "" {
				panic("account: backup code must not be empty")
			}
			digest := sha256.Sum256([]byte(code))
			occurrence := seen[digest]
			seen[digest] = occurrence + 1
			c.Key(fmt.Sprintf("%x.%d", digest, occurrence))
			b := button(c, code).Column().Gap(t.Space(1)).AlignItems(ui.Center).Border(1, t.Border).Checked(*copied == i)
			b.Children(func() {
				ui.Text(c, code).Font("monospace")
				if *copied == i {
					kit.Muted(c, o.Copied)
				}
			})
			if b.Clicked() {
				c.WriteClipboard(code)
				*copied = i
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("BackupCodes", "copied", "index", i)
			}
		}
	})
	return root
}
