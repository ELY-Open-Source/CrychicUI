package devtools

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DeployStep is one pipeline step.
type DeployStep struct {
	Name     string
	State    StepState
	Duration time.Duration
}

// DeployInfo is one deploy of a commit to a stage.
type DeployInfo struct {
	Commit, Message, Author, Branch string
	Stage                           Stage
	Steps                           []DeployStep
}

// DeployOptions word the card; Stages names the four stages.
type DeployOptions struct {
	Running, Ready, Failed, Queued string
	Stages                         [4]string
}

// Deploy is a deploy card with its steps.
func Deploy(c *ui.Context, d DeployInfo, o DeployOptions) ui.Element {
	if d.Commit == "" || d.Message == "" || d.Author == "" || d.Branch == "" || len(d.Steps) == 0 {
		panic(fmt.Sprintf("devtools: Deploy %q needs a commit, message, author, branch and steps", d.Commit))
	}
	validStage("Deploy", d.Stage)
	names := make([]string, len(d.Steps))
	done, failed, running := 0, false, false
	for i, s := range d.Steps {
		if s.State < StepPending || s.State > StepFailed || s.Duration < 0 {
			panic(fmt.Sprintf("devtools: Deploy step %q needs a state and non-negative duration", s.Name))
		}
		names[i] = s.Name
		switch s.State {
		case StepDone:
			done++
		case StepFailed:
			failed = true
		case StepRunning:
			running = true
		}
	}
	unique("Deploy", names)
	def(&o.Running, "Building")
	def(&o.Ready, "Ready")
	def(&o.Failed, "Failed")
	def(&o.Queued, "Queued")
	t := c.Theme()
	status, ink := o.Queued, t.TextMuted
	switch {
	case failed:
		status, ink = o.Failed, t.Danger
	case done == len(d.Steps):
		status, ink = o.Ready, t.Success
	case running || done > 0:
		status, ink = o.Running, t.Accent
	}
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(4)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				tile(c, "rocket", ink, t.Space(9))
				ui.Text(c, d.Message).FontSize(t.Rem(0.9)).Bold().Grow(1).MinWidth(0).Wrap()
				c.Key(status)
				kit.Arrive(kit.Chip(c, status, ink), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			})
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				EnvBadge(c, d.Stage, EnvBadgeOptions{Stages: o.Stages})
				ui.Box(c).Padding(t.Space(0.5), t.Space(1.5)).Radius(t.Radius * 0.6).Background(t.Surface).Children(func() {
					ui.Text(c, d.Commit).Font("monospace").FontSize(t.Rem(0.72)).Bold().Shrink(0)
				})
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "git-branch", 14)
					ui.Text(c, d.Branch).Font("monospace").FontSize(t.Rem(0.72)).Shrink(0)
				})
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					display.Avatar(c, d.Author, nil).Size(t.Space(6), t.Space(6)).FontSize(t.Rem(0.55)).Shrink(0)
					kit.Muted(c, d.Author).FontSize(t.Rem(0.75)).Shrink(0)
				})
			})
			meter(c, float32(done)/float32(len(d.Steps)), ink)
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(3))
		list.Children(func() {
			for i, s := range d.Steps {
				c.Key(s.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(1))
				kit.Stagger(c, list, row, i).Children(func() {
					stepMark(c, s.State)
					name := ui.Text(c, s.Name).FontSize(t.Rem(0.82)).Grow(1).MinWidth(0).SingleLine()
					if s.State == StepPending {
						name.TextColor(t.TextMuted)
					}
					if s.Duration > 0 {
						kit.Muted(c, took(s.Duration)).Font("monospace").FontSize(t.Rem(0.7)).Shrink(0)
					}
				})
			}
		})
	})
	return root
}
