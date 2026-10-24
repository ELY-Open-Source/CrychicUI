package code

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PRState is a pull request's state.
type PRState uint8

// Pull request states.
const (
	PROpen PRState = iota
	PRMerged
	PRClosed
	prStateCount
)

// CheckState is a CI check's state.
type CheckState uint8

// Check states.
const (
	CheckPending CheckState = iota
	CheckPass
	CheckFail
	checkStateCount
)

// Check is one CI check and how long it ran.
type Check struct {
	Name     string
	State    CheckState
	Duration time.Duration
}

func checkOne(name string, k Check) {
	if k.Name == "" || k.State >= checkStateCount || k.Duration < 0 {
		panic(fmt.Sprintf("code: %s check %q needs a name, a known state and a duration", name, k.Name))
	}
}

// PR is a pull request.
type PR struct {
	Number                    int
	Title, Author, Head, Base string
	State                     PRState
	Added, Removed            int
	Checks                    []Check
	Reviewers                 []string
}

// PullRequestOptions are the copy knobs of PullRequest.
type PullRequestOptions struct {
	Open, Merged, Closed string
	// Passing words the checks, "%d" taking passed then total.
	Passing string
}

// PullRequest is a title, status, branches, checks and reviewers card.
func PullRequest(c *ui.Context, pr PR, o PullRequestOptions) ui.Element {
	if pr.Title == "" || pr.Author == "" || pr.Head == "" || pr.Base == "" || pr.Number < 1 || pr.State >= prStateCount ||
		pr.Added < 0 || pr.Removed < 0 {
		panic(fmt.Sprintf("code: PullRequest %q needs a number, author, branches, a known state and counts", pr.Title))
	}
	for _, k := range pr.Checks {
		checkOne("PullRequest", k)
	}
	for _, r := range pr.Reviewers {
		if r == "" {
			panic(fmt.Sprintf("code: PullRequest %q needs a name for every reviewer", pr.Title))
		}
	}
	def(&o.Open, "Open")
	def(&o.Merged, "Merged")
	def(&o.Closed, "Closed")
	def(&o.Passing, "%d of %d passing")
	t := c.Theme()
	name, glyph, ink := o.Open, "git-pull-request", t.Success
	switch pr.State {
	case PRMerged:
		name, glyph, ink = o.Merged, "git-merge", t.Accent
	case PRClosed:
		name, glyph, ink = o.Closed, "git-pull-request-closed", t.Danger
	}
	root := kit.Card(c).FillWidth().Label(pr.Title)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Children(func() {
			ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Background(ink.Alpha(0.12)).TextColor(ink).
				Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16) })
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
					ui.Text(c, pr.Title).Bold().Grow(1).MinWidth(0).Wrap()
					kit.Muted(c, "#"+strconv.Itoa(pr.Number)).Shrink(0)
				})
				ui.Row(c).Wrap().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					kit.Muted(c, pr.Author).FontSize(t.Rem(0.8)).Shrink(0)
					branchChip(c, pr.Head)
					ui.Icon(c, icon.Must("arrow-right")).Size(12, 12).TextColor(t.TextMuted)
					branchChip(c, pr.Base)
				})
			})
			kit.Chip(c, name, ink)
		})
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Background(t.Surface).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			if len(pr.Checks) > 0 {
				checkTally(c, pr.Checks, o.Passing)
			}
			ui.Box(c).Grow(1)
			kit.Numeric(c, "+"+strconv.Itoa(pr.Added)).FontSize(t.Rem(0.75)).Bold().TextColor(t.Success).Shrink(0)
			kit.Numeric(c, "−"+strconv.Itoa(pr.Removed)).FontSize(t.Rem(0.75)).Bold().TextColor(t.Danger).Shrink(0)
			if len(pr.Reviewers) > 0 {
				display.AvatarGroup(c, pr.Reviewers, display.AvatarGroupOptions{})
			}
		})
	})
}

// branchChip is a branch name in a mono pill.
func branchChip(c *ui.Context, name string) {
	t := c.Theme()
	ui.Text(c, name).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.Accent).Shrink(0).
		Padding(0, t.Space(1.5)).Radius(t.Radius).Background(t.Accent.Alpha(0.1))
}

// checkTally is a state bar and passing count; pending pulses.
func checkTally(c *ui.Context, checks []Check, format string) {
	t := c.Theme()
	var n [checkStateCount]int
	for _, k := range checks {
		n[k.State]++
	}
	bar := ui.Box(c).Size(t.Space(20), t.Space(1.5)).Shrink(0).Radius(999).Background(t.Border)
	pulse := float32(1)
	if n[CheckPending] > 0 {
		pulse = 0.45 + 0.55*float32(1-2*abs(bar.Loop("pulse", 1200*time.Millisecond, ui.Linear)-0.5))
	}
	inks := [checkStateCount]ui.Color{CheckPending: t.Warning.Alpha(pulse), CheckPass: t.Success, CheckFail: t.Danger}
	bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		x := r.X
		for _, s := range []CheckState{CheckPass, CheckFail, CheckPending} {
			if w := r.W * float32(n[s]) / float32(len(checks)); w > 0 {
				p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, inks[s], 0)
				x += w
			}
		}
	}))
	kit.Muted(c, fmt.Sprintf(format, n[CheckPass], len(checks))).FontSize(t.Rem(0.75)).Shrink(0)
}
