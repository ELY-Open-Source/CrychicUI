package learn

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Award is who completed what, when, and who signs it.
type Award struct {
	Name, Course, Issuer, Signer, ID string
	Date                             time.Time
}

// CertificateOptions word the certificate; Date is a time layout.
type CertificateOptions struct{ Title, Awarded, For, Date, Signed, Number string }

// Certificate shows a completion certificate with a seal.
func Certificate(c *ui.Context, a Award, o CertificateOptions) ui.Element {
	if strings.TrimSpace(a.Name) == "" || strings.TrimSpace(a.Course) == "" || strings.TrimSpace(a.Issuer) == "" || a.Date.IsZero() {
		panic(fmt.Sprintf("learn: Certificate %q needs a name, course, issuer and date", a.Course))
	}
	def(&o.Title, "Certificate of completion")
	def(&o.Awarded, "This certifies that")
	def(&o.For, "has completed")
	def(&o.Date, "January 2, 2006")
	def(&o.Signed, "Signed")
	def(&o.Number, "No. %s")
	t := c.Theme()
	root := card(c).Padding(t.Space(2)).Role(ui.RoleGroup).Label(o.Title)
	root.Children(func() {
		frame := ui.Column(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(5), t.Space(4)).Radius(t.Radius).
			Border(1, t.Warning.Alpha(0.5)).Background(t.Warning.Alpha(0.04))
		frame.Children(func() {
			seal := ui.Box(c).Size(t.Space(14), t.Space(14)).Shrink(0).Center().Radius(999).Background(t.Warning.Alpha(0.16)).
				Border(t.Space(0.5), t.Warning.Alpha(0.6)).TextColor(t.Warning)
			kit.Arrive(seal, ui.Pose{Scale: 0.4, Rotate: -45, Fade: 1}, 200*time.Millisecond).Children(func() { glyph(c, "award", 28) })
			ui.Text(c, strings.ToUpper(o.Title)).FontSize(t.Rem(0.72)).Bold().TextColor(t.Warning.Mix(t.Text, 0.4)).
				FillWidth().TextAlign(ui.Center).Wrap()
			kit.Muted(c, o.Awarded).FontSize(t.Rem(0.76)).FillWidth().TextAlign(ui.Center).Wrap()
			name := ui.Text(c, a.Name).Font("serif").FontSize(t.Rem(1.6)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
			kit.Arrive(name, ui.Pose{Y: t.Space(1.5), Fade: 1}, 350*time.Millisecond)
			ui.Box(c).Width(t.Space(24)).Height(1).Background(t.Warning.Alpha(0.5))
			kit.Muted(c, o.For).FontSize(t.Rem(0.76)).FillWidth().TextAlign(ui.Center).Wrap()
			ui.Text(c, a.Course).FontSize(t.Rem(1)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
			ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.End).Padding(t.Space(3), 0, 0, 0).Children(func() {
				side := func(top, under string) {
					ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, top).FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
						ui.Box(c).FillWidth().Height(1).Background(t.Border)
						kit.Muted(c, under).FontSize(t.Rem(0.66)).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
					})
				}
				side(a.Date.Format(o.Date), a.Issuer)
				if a.Signer != "" {
					side(a.Signer, o.Signed)
				}
			})
			if a.ID != "" {
				kit.Muted(c, fmt.Sprintf(o.Number, a.ID)).Font("monospace").FontSize(t.Rem(0.62)).FillWidth().MinWidth(0).
					TextAlign(ui.Center).SingleLine()
			}
		})
	})
	return root
}
