package account

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestResetRejectsShortMixedPassword(t *testing.T) {
	value := ResetValue{Password: "Aa1", Confirm: "Aa1"}
	sent := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, submit := ResetPassword(c, &value, ResetPasswordOptions{})
		sent = sent || submit.Changed()
	}), 500, 350)
	settles(t, tt)
	if err := tt.Click("Reset password"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if sent || !tt.HasText("Weak") {
		t.Errorf("short password accepted: sent %v texts %q", sent, tt.Texts())
	}
}

func TestUsageFractions(t *testing.T) {
	resources := []Resource{{Name: "Storage", Unit: "GB", Used: 0.1, Limit: 0.2}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Usage(c, resources, UsageOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("0.1 of 0.2 GB") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestAPIKeyCopyTracksSecret(t *testing.T) {
	keys := []APIKey{{ID: "k1", Name: "CI", Prefix: "ck_", Secret: "first"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { APIKeys(c, &keys, APIKeysOptions{}) }), 500, 250)
	settles(t, tt)
	click(t, tt, "Copy CI")
	keys[0].Secret = "second"
	tt.Frame()
	if !tt.HasText("Copy") || tt.HasText("Copied") || tt.Clipboard() != "first" {
		t.Errorf("stale copy state: %q clipboard %q", tt.Texts(), tt.Clipboard())
	}
}

func TestSignInSubmitKeepsFocusOnError(t *testing.T) {
	value := Credentials{Email: "ana@example.com", Password: "password"}
	action, message := "", ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SignIn(c, &value, nil, &action, SignInOptions{Error: message})
	}), 500, 400)
	settles(t, tt)
	click(t, tt, "Sign in")
	message = "Try again"
	tt.Frame()
	if !tt.Focused("Sign in") {
		t.Error("submit lost focus after error")
	}
}

func TestPlanButtonKeepsFocusOnFeatureInsert(t *testing.T) {
	change := false
	plan := Plan{Name: "Pro", Price: "$10"}
	tt := ui.NewTester(padded(func(c *ui.Context) { PlanCard(c, plan, true, &change, PlanCardOptions{}) }), 500, 350)
	settles(t, tt)
	click(t, tt, "Change plan")
	plan.Features = append(plan.Features, "New feature")
	tt.Frame()
	if !tt.Focused("Change plan") {
		t.Error("change button lost focus after feature insertion")
	}
}

func TestBackupCodesIdentityDuringReorder(t *testing.T) {
	codes := []string{"code-A", "code-B"}
	copied := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { BackupCodes(c, codes, &copied, BackupCodesOptions{}) }), 500, 250)
	settles(t, tt)
	r, _ := tt.Find("code-A")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	codes[0], codes[1] = codes[1], codes[0]
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
	if tt.Clipboard() == "code-B" {
		t.Error("reordered press copied the wrong code")
	}
	if tt.Clipboard() != "" && tt.Clipboard() != "code-A" {
		t.Errorf("unexpected clipboard %q", tt.Clipboard())
	}
}

func TestBackupCodesDuplicateDoesNotLeak(t *testing.T) {
	copied := -1
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("duplicate recovery codes panic: %v", r)
		}
	}()
	ui.Render(func(c *ui.Context) {
		BackupCodes(c, []string{"same-secret", "same-secret"}, &copied, BackupCodesOptions{})
	}, 500, 220, 1)
}
