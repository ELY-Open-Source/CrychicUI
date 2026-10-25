package account

import (
	"crypto/sha256"
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BackupCodesOptions word the card; Copied takes the code number.
type BackupCodesOptions struct{ Title, Subtitle, Copied string }

// BackupCodes copies recovery codes; cue Pop.
func BackupCodes(c *ui.Context, codes []string, copied *int, o BackupCodesOptions) ui.Element {
	if len(codes) == 0 || copied == nil || *copied < -1 || *copied >= len(codes) {
		panic("account: BackupCodes needs codes and state")
	}
	def(&o.Title, "Backup codes")
	def(&o.Subtitle, "Each code works once. Keep them somewhere safe.")
	def(&o.Copied, "Code %02d copied")
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		head(c, "key-square", o.Title, o.Subtitle)
		grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2))
		seen := map[[sha256.Size]byte]int{}
		grid.Children(func() {
			for i, code := range codes {
				if code == "" {
					panic("account: BackupCodes code must not be empty")
				}
				digest := sha256.Sum256([]byte(code))
				occurrence := seen[digest]
				seen[digest] = occurrence + 1
				c.Key(fmt.Sprintf("%x.%d", digest, occurrence))
				on := *copied == i
				b := kit.Press(ui.ButtonBase(c).Label(code).Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).MinWidth(0).Checked(on))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Success.Alpha(0.12), unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Success, unit(on)), ui.Snappy))
				kit.Stagger(c, grid, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						kit.Muted(c, fmt.Sprintf("%02d", i+1)).FontSize(t.Rem(0.7)).Shrink(0)
						ui.Text(c, code).Font("monospace").FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
						c.Key(on)
						if on {
							kit.Arrive(glyph(c, "check", 14).TextColor(t.Success), ui.Pose{Scale: 0.4, Fade: 1}, 0)
						} else {
							glyph(c, "copy", 14).TextColor(t.TextMuted)
						}
					})
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
		if *copied >= 0 {
			c.Key("copied")
			kit.Arrive(kit.Muted(c, fmt.Sprintf(o.Copied, *copied+1)).FontSize(t.Rem(0.75)), ui.Pose{Fade: 1}, 0)
		}
	})
	return root
}
