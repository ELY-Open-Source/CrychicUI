package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuotaWarningOptions are the quota labels.
type QuotaWarningOptions struct {
	Format, Upgrade string
}

// QuotaWarning shows near-full usage with upgrade; cue Warning.
func QuotaWarning(c *ui.Context, used, total int64, upgraded *bool, o QuotaWarningOptions) ui.Element {
	if used < 0 || total <= 0 || used > total || upgraded == nil {
		panic("files: QuotaWarning needs valid usage and state")
	}
	def(&o.Format, "%.0f%% full")
	def(&o.Upgrade, "Upgrade")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).
		Background(t.Warning.Alpha(0.12)).Border(1, t.Warning).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must("triangle-alert")).Size(16, 16).TextColor(t.Warning)
			ui.Text(c, fmt.Sprintf(o.Format, 100*float64(used)/float64(total))).Bold().Grow(1)
			if !*upgraded {
				b := kit.Press(ui.ButtonBase(c).Label(o.Upgrade).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
					Background(t.Warning).TextColor(t.Background))
				b.Children(func() { ui.Text(c, o.Upgrade).SingleLine() })
				if b.Clicked() {
					*upgraded = true
					root.SetChanged()
					sound.Play(sound.Warning)
					kit.Log("QuotaWarning", "upgraded")
				}
			}
		})
		progressBar(c, float64(used)/float64(total), t.Warning)
	})
	return root
}
