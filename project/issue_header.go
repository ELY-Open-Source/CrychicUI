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

// IssueHeaderOptions are the issue actions, states and priorities.
type IssueHeaderOptions struct{ Edit, Close, Todo, Doing, Done, Urgent, High, Medium, Low string }

// IssueHeader shows key, title, status, and actions.
func IssueHeader(c *ui.Context, issue Issue, status TaskStatus, action *IssueAction, o IssueHeaderOptions) ui.Element {
	validIssue("IssueHeader", issue)
	if status > Done || action == nil || *action > IssueClose {
		panic("project: IssueHeader has invalid data")
	}
	def(&o.Edit, "Edit")
	def(&o.Close, "Close")
	def(&o.Todo, "To do")
	def(&o.Doing, "Doing")
	def(&o.Done, "Done")
	def(&o.Urgent, "Urgent")
	def(&o.High, "High")
	def(&o.Medium, "Medium")
	def(&o.Low, "Low")
	t := c.Theme()
	kind, kindInk := issueLook(t, issue.Type)
	_, stateInk := statusLook(t, status)
	_, levelInk := priorityLook(t, issue.Priority)
	root := card(c).Gap(t.Space(2.5)).Padding(t.Space(3))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			glyph(c, kind, 14).TextColor(kindInk)
			ui.Text(c, issue.Key).Font("monospace").FontSize(t.Rem(0.75)).Bold().TextColor(t.TextMuted).Shrink(0)
			kit.Muted(c, "· "+issue.Type).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
			for i, item := range []struct {
				name, glyph string
				action      IssueAction
				ink         ui.Color
			}{{o.Edit, "pencil", IssueEdit, t.Text}, {o.Close, "circle-check", IssueClose, t.Success}} {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(item.name).Shrink(0).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy)).Border(1, t.Border)
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(item.ink).Children(func() {
						glyph(c, item.glyph, 14)
						ui.Text(c, item.name).FontSize(t.Rem(0.75)).Bold().Shrink(0)
					})
				})
				if b.Clicked() {
					*action = item.action
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("IssueHeader", "action", "action", item.action)
				}
			}
		})
		ui.Text(c, issue.Title).FontSize(t.Rem(1.15)).Bold().Wrap().FillWidth()
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
			c.Key(status)
			kit.Arrive(tag(c, statusText(status, o.Todo, o.Doing, o.Done), stateInk), ui.Pose{Scale: 0.8, Fade: 1}, 0)
			tag(c, priorityText(issue.Priority, o.Urgent, o.High, o.Medium, o.Low), levelInk)
		})
	})
	return root
}
