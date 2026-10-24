package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// IssueAction is an issue header action.
type IssueAction uint8

const (
	IssueNone IssueAction = iota
	IssueEdit
	IssueClose
)

// IssueHeaderOptions are the issue actions and states.
type IssueHeaderOptions struct{ Edit, Close, Todo, Doing, Done string }

// IssueHeader shows key, title, status, and actions.
func IssueHeader(c *ui.Context, issue Issue, status TaskStatus, action *IssueAction, o IssueHeaderOptions) ui.Element {
	if issue.Key == "" || issue.Title == "" || issue.Type == "" || issue.Priority > Low || status > Done || action == nil || *action > IssueClose {
		panic("project: IssueHeader has invalid data")
	}
	def(&o.Edit, "Edit")
	def(&o.Close, "Close")
	def(&o.Todo, "To do")
	def(&o.Doing, "Doing")
	def(&o.Done, "Done")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, issue.Key).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			ui.Text(c, issue.Title).Bold()
		})
		ui.Text(c, statusText(status, o.Todo, o.Doing, o.Done)).TextColor(t.Accent)
		for _, item := range []struct {
			name   string
			action IssueAction
		}{{o.Edit, IssueEdit}, {o.Close, IssueClose}} {
			b := kit.Press(ui.ButtonBase(c).Label(item.name).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
			b.Children(func() { ui.Text(c, item.name).FontSize(t.Rem(0.75)) })
			if b.Clicked() {
				*action = item.action
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("IssueHeader", "action", "action", *action)
			}
		}
	})
	return root
}
