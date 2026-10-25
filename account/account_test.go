package account

import (
	"image"
	"strings"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestSignIn(t *testing.T) {
	value := Credentials{Email: "ana@example.com", Password: "password"}
	action := ""
	sent := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, submit := SignIn(c, &value, []string{"Google"}, &action, SignInOptions{})
		sent = sent || submit.Changed()
	}), 500, 450)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Sign in")
	if action != "signin" || !sent {
		t.Errorf("action %q sent %v", action, sent)
	}
	cues(t, before, sound.Tap)
}

func TestSignInError(t *testing.T) {
	value := Credentials{}
	action, message := "", ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SignIn(c, &value, nil, &action, SignInOptions{Error: message})
	}), 500, 400)
	settles(t, tt)
	before := baseline()
	message = "Invalid credentials"
	tt.Frame()
	if !tt.HasText(message) {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Error)
}

func TestSignUp(t *testing.T) {
	value := SignUpValue{Name: "Ana", Email: "ana@example.com", Password: "password", Terms: true}
	sent := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, submit := SignUp(c, &value, SignUpOptions{})
		sent = sent || submit.Changed()
	}), 500, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Create account")
	if !sent {
		t.Error("sign up not reported")
	}
	cues(t, before, sound.Success)
}

func TestForgotPassword(t *testing.T) {
	email, sent := "ana@example.com", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, send := ForgotPassword(c, &email, ForgotPasswordOptions{})
		sent = sent || send.Changed()
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Send link")
	if !sent {
		t.Error("link not sent")
	}
	cues(t, before, sound.Send)
}

func TestResetPassword(t *testing.T) {
	value := ResetValue{Password: "Password1!", Confirm: "Password1!"}
	sent := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, submit := ResetPassword(c, &value, ResetPasswordOptions{})
		sent = sent || submit.Changed()
	}), 500, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Reset password")
	if !sent || !tt.HasText("Strong") {
		t.Errorf("sent %v texts %q", sent, tt.Texts())
	}
	cues(t, before, sound.Success)
}

func TestProviders(t *testing.T) {
	chosen := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Providers(c, []string{"Google"}, &chosen, ProvidersOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Google")
	if chosen != "Google" {
		t.Errorf("chosen %q", chosen)
	}
	cues(t, before, sound.Tap)
}

func TestMagicLink(t *testing.T) {
	resent := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { MagicLink(c, "ana@example.com", &resent, MagicLinkOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Resend")
	if resent != 1 {
		t.Errorf("resent %d", resent)
	}
	cues(t, before, sound.Send)
}

func TestTwoFactor(t *testing.T) {
	code, resent, verified, reported := "123456", 0, false, false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, verify := TwoFactor(c, &code, &resent, &verified, TwoFactorOptions{})
		reported = reported || verify.Changed()
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Verify")
	if !verified || !reported {
		t.Errorf("verified %v reported %v", verified, reported)
	}
	cues(t, before, sound.Success)
}

func TestBackupCodes(t *testing.T) {
	copied := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { BackupCodes(c, []string{"ABCD-1234"}, &copied, BackupCodesOptions{}) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "ABCD-1234")
	if copied != 0 || tt.Clipboard() != "ABCD-1234" {
		t.Errorf("copied %d clipboard %q", copied, tt.Clipboard())
	}
	cues(t, before, sound.Pop)
}

func TestPasskey(t *testing.T) {
	added := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Passkey(c, &added, PasskeyOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add passkey")
	if !added {
		t.Error("passkey not added")
	}
	cues(t, before, sound.Success)
}

func TestSessions(t *testing.T) {
	sessions := []Session{{ID: "s1", Device: "Laptop", Last: "now"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Sessions(c, &sessions, SessionsOptions{}) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Sign out Laptop")
	settles(t, tt)
	if len(sessions) != 0 || !tt.HasText("No other sessions") {
		t.Errorf("sessions %v texts %q", sessions, tt.Texts())
	}
	cues(t, before, sound.Close)
}

func TestDevicesList(t *testing.T) {
	devices := []Device{{ID: "d1", Name: "Laptop", Last: "now"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { DevicesList(c, &devices, DevicesListOptions{}) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Remove Laptop")
	if len(devices) != 0 {
		t.Errorf("devices %v", devices)
	}
	cues(t, before, sound.Delete)
}

func TestProfileSettings(t *testing.T) {
	profile, saved := Profile{Name: "Ana", Handle: "@ana"}, false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, save := ProfileSettings(c, &profile, ProfileSettingsOptions{})
		saved = saved || save.Changed()
	}), 500, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Save")
	if !saved {
		t.Error("profile not saved")
	}
	cues(t, before, sound.Success)
}

func TestAvatarEditor(t *testing.T) {
	var image ui.ImageSource = ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2)))
	change := false
	tt := ui.NewTester(padded(func(c *ui.Context) { AvatarEditor(c, &image, &change, AvatarEditorOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Remove")
	if image != nil {
		t.Error("avatar not removed")
	}
	cues(t, before, sound.Pop)
}

func TestChangeEmail(t *testing.T) {
	email, sent := "ana@example.com", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, verify := ChangeEmail(c, &email, ChangeEmailOptions{})
		sent = sent || verify.Changed()
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Verify")
	if !sent {
		t.Error("verification not sent")
	}
	cues(t, before, sound.Send)
}

func TestChangePassword(t *testing.T) {
	value := PasswordChange{Current: "oldpassword", New: "newpassword"}
	sent := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, submit := ChangePassword(c, &value, ChangePasswordOptions{})
		sent = sent || submit.Changed()
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Update password")
	if !sent {
		t.Error("password not changed")
	}
	cues(t, before, sound.Success)
}

func TestSecurityScore(t *testing.T) {
	steps := []SecurityStep{{Title: "Add passkey"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { SecurityScore(c, 80, steps, SecurityScoreOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("80") || !tt.HasText("Add passkey") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestAPIKeys(t *testing.T) {
	keys := []APIKey{{ID: "k1", Name: "CI", Prefix: "ck_", Secret: "secret"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { APIKeys(c, &keys, APIKeysOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Copy CI")
	if tt.Clipboard() != "secret" || strings.Contains(strings.Join(tt.Texts(), "\n"), "secret") {
		t.Errorf("clipboard or texts leaked: %q", tt.Texts())
	}
	click(t, tt, "Revoke CI")
	if !keys[0].Revoked {
		t.Error("key not revoked")
	}
	cues(t, before, sound.Pop, sound.Delete)
}

func TestPlanCard(t *testing.T) {
	change := false
	plan := Plan{Name: "Pro", Price: "$10", Features: []string{"All"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { PlanCard(c, plan, true, &change, PlanCardOptions{}) }), 500, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Change plan")
	if !change || !tt.HasText("Current plan") {
		t.Errorf("change %v texts %q", change, tt.Texts())
	}
	cues(t, before, sound.Tap)
}

func TestUsage(t *testing.T) {
	resources := []Resource{{Name: "Storage", Unit: "GB", Used: 8, Limit: 10}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Usage(c, resources, UsageOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("8 of 10 GB") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestInvoiceList(t *testing.T) {
	downloaded := ""
	invoices := []Invoice{{ID: "i1", Number: "INV-1", Date: "Oct 1", Amount: "$10", URL: "https://example.test/i"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { InvoiceList(c, invoices, &downloaded, InvoiceListOptions{}) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Download INV-1")
	if downloaded != "i1" {
		t.Errorf("downloaded %q", downloaded)
	}
	cues(t, before, sound.Tap)
}
