// Package social holds social networking components.
package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// User is one social profile.
type User struct {
	ID, Name, Handle, Bio string
	Avatar, Cover         ui.ImageSource
}

// PostData is a post; Likes excludes the viewer.
type PostData struct {
	ID, Text, Time          string
	Author                  User
	Media                   ui.ImageSource
	Likes, Comments, Shares int
}

// PostState is one post's interaction state.
type PostState struct {
	Liked, Saved bool
	Action       string
}

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func button(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	return kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
}

func validUser(component string, u User) {
	if u.ID == "" || u.Name == "" {
		panic(fmt.Sprintf("social: %s user %q (%q) needs id and name", component, u.Name, u.ID))
	}
}

func validPost(component string, p PostData) {
	validUser(component, p.Author)
	if p.ID == "" || p.Text == "" || p.Likes < 0 || p.Comments < 0 || p.Shares < 0 {
		panic(fmt.Sprintf("social: %s post %q needs id, text and counts", component, p.ID))
	}
}

// unique panics on an empty or repeated id.
func unique(component string, ids []string) {
	seen := map[string]bool{}
	for _, id := range ids {
		if id == "" || seen[id] {
			panic(fmt.Sprintf("social: %s has an empty or repeated %q", component, id))
		}
		seen[id] = true
	}
}

// avatar is a photo or colored initials, size DIP wide.
func avatar(c *ui.Context, u User, size float32) ui.Element {
	t := c.Theme()
	if u.Avatar == nil {
		return display.Avatar(c, u.Name, nil).Size(size, size).FontSize(size * 0.36).Shrink(0).AlignSelf(ui.Center)
	}
	box := ui.Box(c).Size(size, size).Shrink(0).Radius(999).Clip().Background(t.SurfaceHover).AlignSelf(ui.Center)
	box.Children(func() { ui.Image(c, u.Avatar).Fill().Fit(ui.Cover) })
	return kit.Arrive(box, ui.Pose{Fade: 1}, 0)
}

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius).Border(1, t.Border).Background(t.Background).Clip()
	return kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// glyph is a Lucide icon size DIP square.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// ring is a circle starting at the top, clockwise.
func ring(cx, cy, r float32) *ui.Path {
	const k = 0.5522847498
	var p ui.Path
	p.MoveTo(cx, cy-r)
	p.CubeTo(cx+r*k, cy-r, cx+r, cy-r*k, cx+r, cy)
	p.CubeTo(cx+r, cy+r*k, cx+r*k, cy+r, cx, cy+r)
	p.CubeTo(cx-r*k, cy+r, cx-r, cy+r*k, cx-r, cy)
	p.CubeTo(cx-r, cy-r*k, cx-r*k, cy-r, cx, cy-r)
	return p.Close()
}

// number groups an integer for the current locale.
func number(c *ui.Context, n any) string {
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), n, i18n.NumberOptions{Integer: true})
}

// follow is a pill: filled to follow, outlined once following.
func follow(c *ui.Context, label, text string, on bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1), t.Space(3)).Radius(999).Shrink(0).Checked(on))
	hover := unit(b.Hovered())
	b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, hover).Mix(t.Background.Mix(t.SurfaceHover, hover), unit(on)), ui.Snappy))
	b.TextColor(b.AnimateColor("ink", t.AccentText.Mix(t.Text, unit(on)), ui.Snappy))
	b.Border(1, b.AnimateColor("line", t.Accent.Mix(t.Border, unit(on)), ui.Snappy))
	return b.Children(func() {
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			c.Key(on)
			if on {
				kit.Arrive(glyph(c, "check", 14), ui.Pose{Scale: 0.4, Fade: 1}, 0)
			} else {
				glyph(c, "user-plus", 14)
			}
			ui.Text(c, text).FontSize(t.Rem(0.8)).Bold().Shrink(0)
		})
	})
}

// byline is an author's avatar, name, handle and time.
func byline(c *ui.Context, u User, when string, size float32) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
		avatar(c, u, size)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
			ui.Text(c, u.Name).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
				ui.Text(c, u.Handle).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
				if when != "" {
					ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.TextMuted)
					ui.Text(c, when).FontSize(t.Rem(0.75)).Shrink(0)
				}
			})
		})
	})
}
