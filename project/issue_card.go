package project

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Issue is one project issue.
type Issue struct {
	Key, Title, Type string
	Priority         PriorityLevel
}

// IssueCardOptions are the priority labels.
type IssueCardOptions struct{ Urgent, High, Medium, Low string }

// IssueCard shows key, title, type icon, and priority.
func IssueCard(c *ui.Context, issue Issue, o IssueCardOptions) ui.Element {
	if issue.Key == "" || issue.Title == "" || issue.Type == "" || issue.Priority > Low {
		panic("project: IssueCard has an invalid issue")
	}
	def(&o.Urgent, "Urgent")
	def(&o.High, "High")
	def(&o.Medium, "Medium")
	def(&o.Low, "Low")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Icon(c, icon.Must("circle-dot")).Size(16, 16).TextColor(t.Accent)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, issue.Key).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			ui.Text(c, issue.Title).SingleLine()
		})
		ui.Text(c, issue.Type).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		ui.Text(c, priorityText(issue.Priority, o.Urgent, o.High, o.Medium, o.Low)).FontSize(t.Rem(0.75)).TextColor(t.Danger)
	})
	return root
}
