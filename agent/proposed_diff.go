package agent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProposedDiffOptions name the buttons and outcomes.
type ProposedDiffOptions struct {
	Apply, Reject, Applied, Rejected string
}

// ProposedDiff shows a diff with Apply and Reject.
func ProposedDiff(c *ui.Context, path, diff string, decision *Decision, o ProposedDiffOptions) (root, apply, reject ui.Element) {
	if path == "" || diff == "" || decision == nil || *decision > Denied {
		panic(fmt.Sprintf("agent: ProposedDiff %q needs a path, a diff and a decision", path))
	}
	def(&o.Apply, "Apply")
	def(&o.Reject, "Reject")
	def(&o.Applied, "Applied")
	def(&o.Rejected, "Rejected")
	lines := strings.Split(strings.TrimRight(diff, "\n"), "\n")
	kinds := diffKinds(lines)
	added, removed := 0, 0
	for _, k := range kinds {
		switch k {
		case '+':
			added++
		case '-':
			removed++
		}
	}
	t := c.Theme()
	root = kit.Card(c).FillWidth().Clip().Role(ui.RoleGroup).Label(path)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Icon(c, icon.Must("file-code")).Size(16, 16).TextColor(t.TextMuted)
			ui.Text(c, path).Font("monospace").FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Numeric(c, "+"+strconv.Itoa(added)).FontSize(t.Rem(0.8)).TextColor(t.Success)
			kit.Numeric(c, "−"+strconv.Itoa(removed)).FontSize(t.Rem(0.8)).TextColor(t.Danger)
		})
		ui.Column(c).FillWidth().Padding(t.Space(1), 0).Children(func() {
			for i, line := range lines {
				c.Key(i)
				bg, ink := ui.Color{}, t.TextMuted
				switch kinds[i] {
				case '+':
					bg, ink = t.Success.Alpha(0.1), t.Text
				case '-':
					bg, ink = t.Danger.Alpha(0.1), t.Text
				}
				ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(0, t.Space(3)).Background(bg).Children(func() {
					kit.Numeric(c, strconv.Itoa(i+1)).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.TextMuted).
						Width(t.Space(5)).Shrink(0)
					ui.Text(c, line).Font("monospace").FontSize(t.Rem(0.8)).TextColor(ink).Grow(1).MinWidth(0).Wrap()
				})
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.End).Padding(t.Space(2), t.Space(3)).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			if *decision != Undecided {
				text, glyph, ink := o.Applied, "circle-check", t.Success
				if *decision == Denied {
					text, glyph, ink = o.Rejected, "circle-x", t.Danger
				}
				done := ui.Row(c).Key(*decision).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(ink)
				kit.Arrive(done, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(16, 16)
					ui.Text(c, text).Bold()
				})
				return
			}
			reject = kit.Press(ui.ButtonBase(c).Label(o.Reject).Padding(t.Space(1), t.Space(3)).Radius(t.Radius).Border(1, t.Border))
			reject.Background(reject.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(reject.Hovered())), ui.Snappy)).
				Children(func() { ui.Text(c, o.Reject).Shrink(0) })
			apply = kit.Press(ui.ButtonBase(c).Label(o.Apply).Padding(t.Space(1), t.Space(3)).Radius(t.Radius))
			apply.Background(apply.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(apply.Hovered())), ui.Snappy)).
				TextColor(t.AccentText).Children(func() { ui.Text(c, o.Apply).Bold().Shrink(0) })
		})
	})
	switch {
	case apply.Clicked():
		*decision = Allowed
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("ProposedDiff", "applied", "added", added, "removed", removed)
	case reject.Clicked():
		*decision = Denied
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("ProposedDiff", "rejected")
	}
	return root, apply, reject
}

// diffKinds marks added and removed lines; file headers stay plain.
func diffKinds(lines []string) []byte {
	out, inHunk := make([]byte, len(lines)), false
	for i, l := range lines {
		out[i] = ' '
		switch {
		case strings.HasPrefix(l, "@@"):
			inHunk = true
		case !inHunk && (strings.HasPrefix(l, "+++") || strings.HasPrefix(l, "---")):
		case strings.HasPrefix(l, "+"):
			out[i] = '+'
		case strings.HasPrefix(l, "-"):
			out[i] = '-'
		}
	}
	return out
}
