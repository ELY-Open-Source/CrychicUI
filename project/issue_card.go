package project

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Issue is one project issue; Assignee and Points are optional.
type Issue struct {
	Key, Title, Type, Assignee string
	Priority                   PriorityLevel
	Points                     int
}

// IssueCardOptions are the priority labels.
type IssueCardOptions struct{ Urgent, High, Medium, Low string }

// IssueCard shows key, title, type icon, and priority.
func IssueCard(c *ui.Context, issue Issue, o IssueCardOptions) ui.Element {
	validIssue("IssueCard", issue)
	def(&o.Urgent, "Urgent")
	def(&o.High, "High")
	def(&o.Medium, "Medium")
	def(&o.Low, "Low")
	t := c.Theme()
	kind, kindInk := issueLook(t, issue.Type)
	level, levelInk := priorityLook(t, issue.Priority)
	root := kit.Lift(kit.Card(c).FillWidth().Row().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).Role(ui.RoleGroup))
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, 0.5*unit(root.Hovered())), ui.Snappy))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() {
		tile(c, kind, kindInk)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, issue.Key).Font("monospace").FontSize(t.Rem(0.7)).Bold().TextColor(t.TextMuted).Shrink(0)
				kit.Muted(c, issue.Type).FontSize(t.Rem(0.7)).Shrink(0)
			})
			ui.Text(c, issue.Title).FontSize(t.Rem(0.85)).Bold().SingleLine().MinWidth(0)
		})
		ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(levelInk).Children(func() {
			glyph(c, level, 14)
			ui.Text(c, priorityText(issue.Priority, o.Urgent, o.High, o.Medium, o.Low)).FontSize(t.Rem(0.7)).Bold().Shrink(0)
		})
		if issue.Points > 0 {
			kit.Chip(c, fmt.Sprint(issue.Points), t.TextMuted)
		}
		if issue.Assignee != "" {
			display.Avatar(c, issue.Assignee, nil).Size(t.Space(6), t.Space(6)).FontSize(t.Rem(0.6)).Tooltip(issue.Assignee)
		}
	})
	return root
}

func validIssue(component string, issue Issue) {
	if issue.Key == "" || issue.Title == "" || issue.Type == "" || issue.Priority > Low || issue.Points < 0 {
		panic(fmt.Sprintf("project: %s has an invalid issue %q", component, issue.Key))
	}
}

// issueLook is an issue type's icon and ink.
func issueLook(t *ui.Theme, kind string) (string, ui.Color) {
	switch strings.ToLower(kind) {
	case "bug":
		return "bug", t.Danger
	case "story":
		return "bookmark", t.Success
	case "task":
		return "square-check", t.Accent
	case "epic":
		return "zap", t.Warning
	}
	return "circle-dot", t.TextMuted
}
