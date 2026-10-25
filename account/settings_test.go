package account

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDangerZoneReportsCommandsOnly(t *testing.T) {
	action, changes := "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if DangerZone(c, "kiln-co", &action, DangerZoneOptions{}).Changed() {
			changes++
		}
	}), 500, 600)
	settles(t, tt)
	click(t, tt, "Transfer kiln-co")
	click(t, tt, "Delete kiln-co")
	click(t, tt, "Cancel")
	if changes != 1 || action != "transfer" {
		t.Errorf("changes %d action %q, want one transfer", changes, action)
	}
}

func TestResendMarksTheInvite(t *testing.T) {
	invites := []Invite{{ID: "a", Email: "jonas@berg.se", Role: "Editor", Sent: "now"}, {ID: "b", Email: "ana@kiln.co", Role: "Viewer", Sent: "now"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { PendingInvites(c, &invites, PendingInvitesOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Resend ana@kiln.co")
	if invites[0].Resent || !invites[1].Resent {
		t.Errorf("resent flags %v %v", invites[0].Resent, invites[1].Resent)
	}
}

func TestLanguageRegionNamedLikeKeys(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("a region choice collided with a key: %v", r)
		}
	}()
	for _, region := range []string{"region", "language", "note"} {
		lang := "English"
		ui.Render(func(c *ui.Context) {
			LanguagePrefs(c, []string{"English"}, &lang, []string{region}, &region, LanguagePrefsOptions{})
		}, 500, 500, 1)
	}
}

func TestLongAccountMenuReachesSignOut(t *testing.T) {
	items := make([]MenuItem, 40)
	for i := range items {
		items[i] = MenuItem{ID: fmt.Sprint(i), Label: fmt.Sprint("Item ", i), Icon: "settings"}
	}
	chosen := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { AccountMenu(c, "Maya Chen", "m@k.co", items, &chosen, AccountMenuOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Account menu")
	settles(t, tt)
	r, ok := tt.Find("Sign out")
	if !ok || r.H == 0 || r.Y+r.H > 400 {
		t.Fatalf("Sign out out of reach: %v %+v", ok, r)
	}
	click(t, tt, "Sign out")
	if chosen != "signout" {
		t.Errorf("chosen %q", chosen)
	}
}

func TestAppearancePreviewsFollowMode(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode, density := ThemeLight, DensityComfortable
		tt := ui.NewTester(padded(func(c *ui.Context) { Appearance(c, &mode, &density, AppearanceOptions{}) }), 600, 500)
		tt.SetDark(dark)
		settles(t, tt)
		img := tt.Image()
		luma := func(name string) int {
			r, _ := tt.Find(name)
			p := img.RGBAAt(int(r.X+r.W-14), int(r.Y+16))
			return int(p.R) + int(p.G) + int(p.B)
		}
		if light, darker := luma("Light"), luma("Dark"); light <= darker {
			t.Errorf("dark mode %v: Light preview %d not brighter than Dark %d", dark, light, darker)
		}
	}
}
