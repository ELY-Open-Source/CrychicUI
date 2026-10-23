package agent

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProposedDiffOptions name the buttons.
type ProposedDiffOptions struct {
	Apply, Reject string
}

// ProposedDiff shows a diff with Apply and Reject.
func ProposedDiff(c *ui.Context, title, diff string, o ProposedDiffOptions) (root, apply, reject ui.Element) {
	if title == "" || diff == "" {
		panic("agent: ProposedDiff needs a title and a diff")
	}
	def(&o.Apply, "Apply")
	def(&o.Reject, "Reject")
	t := c.Theme()
	root = kit.Card(c).AlignSelf(ui.Start).Width(t.Space(96)).Role(ui.RoleGroup).Label(title)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
			ui.Text(c, title).Bold().Grow(1).SingleLine()
			reject = kit.Press(ui.ButtonBase(c).Label(o.Reject).Padding(t.Space(1), t.Space(3)).Radius(t.Radius).Border(1, t.Border))
			reject.Background(reject.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(reject.Hovered())), ui.Snappy)).
				Children(func() { ui.Text(c, o.Reject).SingleLine() })
			apply = kit.Press(ui.ButtonBase(c).Label(o.Apply).Padding(t.Space(1), t.Space(3)).Radius(t.Radius))
			apply.Background(apply.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(apply.Hovered())), ui.Snappy)).
				TextColor(t.AccentText).Children(func() { ui.Text(c, o.Apply).Bold().SingleLine() })
		})
		ui.Column(c).FillWidth().Clip().Padding(t.Space(2), 0).Children(func() {
			for _, line := range strings.Split(strings.TrimRight(diff, "\n"), "\n") {
				row := ui.Text(c, line).Font("monospace").FontSize(t.Rem(0.85)).Wrap()
				switch {
				case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
					row.TextColor(t.Success)
				case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
					row.TextColor(t.Danger)
				default:
					row.TextColor(t.TextMuted)
				}
			}
		})
	})
	if apply.Clicked() {
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("ProposedDiff", "applied")
	}
	if reject.Clicked() {
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("ProposedDiff", "rejected")
	}
	return root, apply, reject
}
