// Package social holds social networking components.
package social

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// User is one social profile.
type User struct {
	ID, Name, Handle, Bio string
	Avatar, Cover         ui.ImageSource
}

// PostData is one post's content and counts.
type PostData struct {
	ID, Text                string
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

func validUser(u User) {
	if u.ID == "" || u.Name == "" {
		panic("social: user needs id and name")
	}
}

func validPost(p PostData) {
	validUser(p.Author)
	if p.ID == "" || p.Text == "" || p.Likes < 0 || p.Comments < 0 || p.Shares < 0 {
		panic("social: invalid post")
	}
}

func avatar(c *ui.Context, u User, size float32) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Radius(999).Clip().Background(t.Accent.Alpha(0.15)).Center()
	box.Children(func() {
		if u.Avatar != nil {
			ui.Image(c, u.Avatar).Fill().Fit(ui.Cover)
		} else {
			r, _ := utf8.DecodeRuneInString(u.Name)
			ui.Text(c, string(r)).TextColor(t.Accent)
		}
	})
	return box
}
