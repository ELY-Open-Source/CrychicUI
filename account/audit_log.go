package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AuditKind is what an audit event touched.
type AuditKind int

const (
	AuditSignIn AuditKind = iota
	AuditChange
	AuditSecurity
	AuditBilling
)

// AuditEvent is one account event.
type AuditEvent struct {
	ID, Actor, Action, Time string
	Kind                    AuditKind
}

// AuditLogOptions word the card.
type AuditLogOptions struct{ Title, Empty string }

// AuditLog lists account events on a timeline.
func AuditLog(c *ui.Context, events []AuditEvent, o AuditLogOptions) ui.Element {
	ids := make([]string, len(events))
	for i, e := range events {
		if e.Actor == "" || e.Action == "" || e.Time == "" || e.Kind < AuditSignIn || e.Kind > AuditBilling {
			panic(fmt.Sprintf("account: AuditLog %q needs an actor, action, time and kind", e.ID))
		}
		ids[i] = e.ID
	}
	if len(events) > 0 {
		unique("AuditLog", ids)
	}
	def(&o.Title, "Audit log")
	def(&o.Empty, "No account events yet")
	t := c.Theme()
	looks := []struct {
		mark string
		ink  ui.Color
	}{{"log-in", t.Accent}, {"pencil", t.TextMuted}, {"shield-alert", t.Warning}, {"credit-card", t.Success}}
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "history", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(3), t.Space(4))
		list.Children(func() {
			if len(events) == 0 {
				empty(c, "history", o.Empty)
			}
			for i, e := range events {
				look := looks[e.Kind]
				c.Key(e.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Width(t.Space(8)).Shrink(0).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Background(look.ink.Alpha(0.14)).
							TextColor(look.ink).Children(func() { glyph(c, look.mark, 14) })
						if i < len(events)-1 {
							ui.Box(c).Width(1).Grow(1).MinHeight(t.Space(3)).Background(t.Border)
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(1), 0, t.Space(3), 0).Children(func() {
						ui.RichText(c, []ui.Span{{Text: e.Actor, Weight: 700}, {Text: " " + e.Action}}...).FontSize(t.Rem(0.85)).Wrap().FillWidth()
						kit.Muted(c, e.Time).FontSize(t.Rem(0.75))
					})
				})
			}
		})
	})
	return root
}
