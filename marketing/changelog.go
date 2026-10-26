package marketing

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ChangeKind sorts a release note.
type ChangeKind int

const (
	ChangeNew ChangeKind = iota
	ChangeImproved
	ChangeFixed
)

// Change is one release note.
type Change struct {
	Kind ChangeKind
	Text string
}

// Release is one version's notes.
type Release struct {
	Version, Title string
	Date           time.Time
	Changes        []Change
}

// ChangelogOptions word the log; Kinds name New, Improved, Fixed.
type ChangelogOptions struct {
	Title, Date, Latest string
	Kinds               [3]string
}

// Changelog lists release notes, newest first.
func Changelog(c *ui.Context, releases []Release, o ChangelogOptions) ui.Element {
	def(&o.Date, "Jan 2, 2006")
	def(&o.Latest, "Latest")
	for i, d := range [3]string{"New", "Improved", "Fixed"} {
		def(&o.Kinds[i], d)
	}
	if len(releases) == 0 {
		panic(fmt.Sprintf("marketing: Changelog %q needs releases", o.Title))
	}
	versions := make([]string, len(releases))
	for i, r := range releases {
		if r.Date.IsZero() || len(r.Changes) == 0 {
			panic(fmt.Sprintf("marketing: Changelog release %q needs a date and changes", r.Version))
		}
		for _, ch := range r.Changes {
			if ch.Kind < ChangeNew || ch.Kind > ChangeFixed || strings.TrimSpace(ch.Text) == "" {
				panic(fmt.Sprintf("marketing: Changelog release %q has a change without a known kind or text", r.Version))
			}
		}
		versions[i] = r.Version
	}
	unique("Changelog", versions)
	t := c.Theme()
	inks := [3]ui.Color{t.Success, t.Accent, t.Warning}
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, "")
		list := ui.Column(c).FillWidth().Gap(t.Space(3))
		list.Children(func() {
			for i, r := range releases {
				c.Key(strconv.Quote(r.Version))
				row := surface(c).Row().Gap(t.Space(4)).Padding(t.Space(4)).AlignItems(ui.Start)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Width(t.Space(22)).Shrink(0).Gap(t.Space(1)).Children(func() {
						ink := t.TextMuted
						if i == 0 {
							ink = t.Accent
						}
						ui.Row(c).AlignSelf(ui.Start).MinWidth(0).MaxWidth(t.Space(22)).Padding(t.Space(0.25), t.Space(2)).Radius(999).Background(ink.Alpha(0.12)).Children(func() {
							ui.Text(c, r.Version).FontSize(t.Rem(0.7)).Bold().TextColor(ink).MinWidth(0).SingleLine()
						})
						kit.Muted(c, r.Date.Format(o.Date)).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
						if i == 0 {
							ui.Text(c, o.Latest).FontSize(t.Rem(0.68)).Bold().TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine()
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(2)).Children(func() {
						if r.Title != "" {
							ui.Text(c, r.Title).FontSize(t.Rem(0.92)).Bold().FillWidth().Wrap()
						}
						for _, ch := range r.Changes {
							ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
								ink := inks[ch.Kind]
								ui.Text(c, o.Kinds[ch.Kind]).FontSize(t.Rem(0.66)).Bold().TextColor(ink).Width(t.Space(18)).Shrink(0).
									TextAlign(ui.Center).Padding(t.Space(0.5), t.Space(1)).Radius(999).Background(ink.Alpha(0.12)).SingleLine()
								ui.Text(c, ch.Text).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
							})
						}
					})
				})
			}
		})
	})
	return root
}
