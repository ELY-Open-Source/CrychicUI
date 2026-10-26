package account

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func panicsWith(t *testing.T, want string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want %q", r, want)
		}
	}()
	ui.Render(view, 400, 400, 1)
}

func TestAccessibility(t *testing.T) {
	var v A11y
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Accessibility(c, &v, AccessibilityOptions{}).Changed() {
			changes++
		}
	}), 420, 520)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Reduce motion")
	if !v.ReduceMotion || v.HighContrast || changes != 1 {
		t.Errorf("motion switch gave %+v after %d changes", v, changes)
	}
	cues(t, before, sound.Toggle)
	before = baseline()
	click(t, tt, "Text size 3")
	if v.TextSize != 2 || changes != 2 {
		t.Errorf("size step gave %+v after %d changes", v, changes)
	}
	cues(t, before, sound.Toggle)
	settles(t, tt)
	before = baseline()
	click(t, tt, "Text size 3")
	if v.TextSize != 2 || changes != 2 {
		t.Errorf("the chosen size changed again: %+v, %d changes", v, changes)
	}
	cues(t, before)
	panicsWith(t, "Accessibility", func(c *ui.Context) { Accessibility(c, &A11y{TextSize: 4}, AccessibilityOptions{}) })
}

func TestSignOut(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, in := SignOut(c, "Maya Chen", SignOutOptions{}); in.Clicked() {
			clicks++
		}
	}), 420, 400)
	settles(t, tt)
	if !tt.HasText("You're signed out") || !tt.HasText("See you soon, Maya Chen.") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Sign in again")
	if clicks != 1 {
		t.Errorf("sign in reported %d clicks", clicks)
	}
	cues(t, before, sound.Tap)
	panicsWith(t, "SignOut needs a name", func(c *ui.Context) { SignOut(c, "", SignOutOptions{}) })
}

func TestVerifyBanner(t *testing.T) {
	resent, changes := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if VerifyBanner(c, "maya@kiln.co", &resent, VerifyBannerOptions{}).Changed() {
			changes++
		}
	}), 520, 300)
	settles(t, tt)
	if !tt.HasText("Verify maya@kiln.co to keep your account secure.") || tt.HasText("Email sent") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Resend email")
	settles(t, tt)
	if resent != 1 || changes != 1 || !tt.HasText("Email sent") {
		t.Errorf("resend gave %d sends, %d changes, texts %q", resent, changes, tt.Texts())
	}
	cues(t, before, sound.Send)
	full := math.MaxInt
	panicsWith(t, "VerifyBanner", func(c *ui.Context) { VerifyBanner(c, "not an email", &resent, VerifyBannerOptions{}) })
	tt = ui.NewTester(padded(func(c *ui.Context) { VerifyBanner(c, "maya@kiln.co", &full, VerifyBannerOptions{}) }), 520, 300)
	settles(t, tt)
	before = baseline()
	click(t, tt, "Resend email")
	if full != math.MaxInt {
		t.Errorf("a full resend count wrapped to %d", full)
	}
	cues(t, before)
}
