package display

import (
	"fmt"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestBadge(t *testing.T) {
	for _, variant := range []BadgeVariant{BadgeNeutral, BadgeAccent, BadgeSuccess, BadgeWarning, BadgeDanger} {
		t.Run(fmt.Sprint(variant), func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) { Badge(c, "Live", BadgeOptions{Variant: variant}) }, 300, 200)
			settles(t, tt, "Badge")
			find(t, tt, "Live")
		})
	}
	panics(t, "Badge variant", func() {
		ui.Render(func(c *ui.Context) { Badge(c, "Live", BadgeOptions{Variant: BadgeVariant(9)}) }, 10, 10, 1)
	})
}

func TestDotBadge(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		DotBadge(c, DotBadgeOptions{}, func() { ui.Text(c, "Inbox") })
	}, 300, 200)
	settles(t, tt, "DotBadge")
	find(t, tt, "Inbox")
}

func TestCountBadge(t *testing.T) {
	count := 2
	tt := ui.NewTester(func(c *ui.Context) {
		CountBadge(c, count, func() { ui.Text(c, "Inbox") })
	}, 300, 200)
	count = 120
	tt.Frame()
	settles(t, tt, "CountBadge")
	find(t, tt, "99+")
	panics(t, "CountBadge count", func() { ui.Render(func(c *ui.Context) { CountBadge(c, -1, func() {}) }, 10, 10, 1) })
}

func TestStatusBadge(t *testing.T) {
	live := false
	tt := ui.NewTester(func(c *ui.Context) { StatusBadge(c, "Deploy", live) }, 300, 200)
	find(t, tt, "Deploy")
	live = true
	tt.Frame()
	if !tt.FramePending() {
		t.Error("live dot did not pulse")
	}
	live = false
	tt.Frame()
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Error("offline dot kept pulsing")
	}
}

func TestTag(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Tag(c, "Design", TagOptions{Icon: "tag"}) }, 300, 200)
	settles(t, tt, "Tag")
	find(t, tt, "Design")
}

func TestRemovableTag(t *testing.T) {
	visible, changed := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if RemovableTag(c, "Docs", &visible).Changed() {
			changed++
		}
	}, 300, 200)
	before := baseline()
	if err := tt.Click("Remove Docs"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "RemovableTag exit")
	if visible || changed != 1 || !appended(before, sound.Delete) || tt.HasText("Docs") {
		t.Errorf("visible %v, changed %d, texts %q", visible, changed, tt.Texts())
	}
	panics(t, "RemovableTag state", func() { ui.Render(func(c *ui.Context) { RemovableTag(c, "Docs", nil) }, 10, 10, 1) })
}

func TestChip(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Chip(c, "Ada", testBitmap()) }, 300, 200)
	r := find(t, tt, "Ada")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "Chip hover")
	panics(t, "Chip avatar", func() { ui.Render(func(c *ui.Context) { Chip(c, "Ada", nil) }, 10, 10, 1) })
}

func TestAvatar(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Avatar(c, "Crychic UI", nil) }, 300, 200)
	settles(t, tt, "Avatar initials")
	find(t, tt, "CU")
	panics(t, "Avatar name", func() { ui.Render(func(c *ui.Context) { Avatar(c, " ", nil) }, 10, 10, 1) })
}

func TestAvatarImage(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Avatar(c, "Ada", testBitmap()) }, 300, 200)
	settles(t, tt, "Avatar image")
	if tt.HasText("AD") {
		t.Error("image avatar still shows initials")
	}
}

func TestAvatarLateImage(t *testing.T) {
	var image *ui.Bitmap
	tt := ui.NewTester(func(c *ui.Context) { Avatar(c, "Ada", image) }, 300, 200)
	settles(t, tt, "Avatar initials")
	image = testBitmap()
	tt.Frame()
	settles(t, tt, "Avatar image fade")
	if tt.HasText("AD") {
		t.Error("initials remain after the image arrived")
	}
}

func TestAvatarGroup(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		AvatarGroup(c, []string{"Ada Lovelace", "Grace Hopper", "Alan Turing", "Katherine Johnson"}, AvatarGroupOptions{Max: 3})
	}, 400, 200)
	settles(t, tt, "AvatarGroup arrival")
	more := find(t, tt, "+1")
	second := find(t, tt, "GH")
	tt.Move(second.X+second.W/2, second.Y+second.H/2)
	settles(t, tt, "AvatarGroup spread")
	if got := find(t, tt, "+1").X; got <= more.X {
		t.Errorf("group did not spread: %v to %v", more.X, got)
	}
}

func TestAvatarStatus(t *testing.T) {
	status := PresenceOffline
	tt := ui.NewTester(func(c *ui.Context) { AvatarStatus(c, "Ada", nil, status) }, 300, 200)
	settles(t, tt, "AvatarStatus arrival")
	status = PresenceBusy
	tt.Frame()
	settles(t, tt, "AvatarStatus color")
	panics(t, "AvatarStatus status", func() {
		ui.Render(func(c *ui.Context) { AvatarStatus(c, "Ada", nil, PresenceStatus(9)) }, 10, 10, 1)
	})
}

func TestPersona(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Persona(c, "Ada Lovelace", "Engineer", nil) }, 400, 200)
	settles(t, tt, "Persona")
	find(t, tt, "Ada Lovelace")
	find(t, tt, "Engineer")
}
