// Code components for the gallery.
package code

import (
	"fmt"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"time"
)

type PRState uint8

const (
	PROpen PRState = iota
	PRMerged
	PRClosed
	prStateCount
)

type CheckState uint8

const (
	CheckPending CheckState = iota
	CheckPass
	CheckFail
	checkStateCount
)

type Check struct {
	Name  string
	State CheckState
}

func checkOne(name string, k Check) {
	if k.Name == "" || k.State >= checkStateCount {
		panic(fmt.Sprintf("code: %s needs a name and a known state in every check", name))
	}
}

// PullRequestOptions are the copy knobs of PullRequest.
type PullRequestOptions struct {
	Open    string
	Merged  string
	Closed  string
	Passing string
}

// PullRequest is a title, status, checks and reviewers card.
func PullRequest(c *ui.Context, title string, state PRState, checks []Check, reviewers []string, o PullRequestOptions) ui.Element {
	if title == "" || state >= prStateCount {
		panic("code: PullRequest needs a title and a known state")
	}
	def(&o.Open, "Open")
	def(&o.Merged, "Merged")
	def(&o.Closed, "Closed")
	def(&o.Passing, "%d of %d passing")
	for _, k := range checks {
		checkOne("PullRequest", k)
	}
	t := c.Theme()
	name, ink := o.Open, t.Success
	if state == PRMerged {
		name, ink = o.Merged, t.Accent
	} else if state == PRClosed {
		name, ink = o.Closed, t.Danger
	}
	root := kit.Card(c).AlignSelf(ui.Start).Width(t.Space(90)).Label(title)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, iconPR).Size(16, 16).TextColor(ink)
			ui.Text(c, title).Bold().Grow(1).Wrap()
			ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().TextColor(ink).
				Padding(t.Space(0.5), t.Space(1.5)).Background(ink.Alpha(0.12)).Radius(999)
		})
		if len(checks) > 0 {
			passing := 0
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			row.Children(func() {
				for i, k := range checks {
					c.Key(i)
					if k.State == CheckPass {
						passing++
					}
					glyph, cink := iconCheck, t.Success
					if k.State == CheckFail {
						glyph, cink = iconX, t.Danger
					} else if k.State == CheckPending {
						glyph, cink = iconClock, t.TextMuted
					}
					e := ui.Icon(c, glyph).Size(14, 14).TextColor(cink).Tooltip(k.Name)
					if k.State == CheckPending {
						e.Opacity(0.5 + 0.5*e.Loop("pulse", 1200*time.Millisecond, ui.Linear))
					}
				}
				kit.Muted(c, fmt.Sprintf(o.Passing, passing, len(checks))).FontSize(t.Rem(0.8)).Grow(1)
			})
		}
		if len(reviewers) > 0 {
			ui.Row(c).Gap(t.Space(1)).Children(func() {
				for _, r := range reviewers {
					if r == "" {
						panic("code: PullRequest needs a name for every reviewer")
					}
					c.Key(r)
					avatarChip(c, r)
				}
			})
		}
	})
}

// Checks are CI checks with spinners and results.
