package account

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func exact(tt *ui.Tester, s string) bool { return slices.Contains(tt.Texts(), s) }

func panics(t *testing.T, want string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want %q", r, want)
		}
	}()
	ui.Render(view, 600, 600, 1)
}

// swapUnderPress presses name, reorders, then releases there.
func swapUnderPress(t *testing.T, tt *ui.Tester, name string, reorder func()) {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("%q missing: %q", name, tt.Texts())
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	reorder()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func TestAccountListsKeepIdentity(t *testing.T) {
	providers, chosen := []string{"Google", "Apple"}, ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Providers(c, providers, &chosen, ProvidersOptions{}) }), 500, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Apple", func() { providers[0], providers[1] = providers[1], providers[0] })
	if chosen == "Google" {
		t.Error("Providers chose the row that slid in")
	}

	sessions := []Session{{ID: "a", Device: "iPhone 15", Last: "now"}, {ID: "b", Device: "iPad Air", Last: "now"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Sessions(c, &sessions, SessionsOptions{}) }), 500, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Sign out iPad Air", func() { sessions[0], sessions[1] = sessions[1], sessions[0] })
	if len(sessions) == 1 && sessions[0].ID == "b" {
		t.Errorf("Sessions signed out the row that slid in: %v", sessions)
	}

	keys := []APIKey{{ID: "p", Name: "Production", Prefix: "ck_", Secret: "1"}, {ID: "s", Name: "Staging", Prefix: "ck_", Secret: "2"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { APIKeys(c, &keys, APIKeysOptions{}) }), 500, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Revoke Staging", func() { keys[0], keys[1] = keys[1], keys[0] })
	for _, k := range keys {
		if k.ID == "p" && k.Revoked {
			t.Error("APIKeys revoked the row that slid in")
		}
	}
}

func TestAccountRepeatsPanic(t *testing.T) {
	cases := map[string]func(c *ui.Context){
		"Providers": func(c *ui.Context) { Providers(c, []string{"Google", "Google"}, new(string), ProvidersOptions{}) },
		"SignIn": func(c *ui.Context) {
			SignIn(c, &Credentials{}, []string{"Apple", "Apple"}, new(string), SignInOptions{})
		},
		"Sessions": func(c *ui.Context) {
			s := []Session{{ID: "a", Device: "Mac", Last: "now"}, {ID: "a", Device: "PC", Last: "now"}}
			Sessions(c, &s, SessionsOptions{})
		},
		"DevicesList": func(c *ui.Context) {
			d := []Device{{ID: "a", Name: "Mac", Last: "now"}, {ID: "a", Name: "PC", Last: "now"}}
			DevicesList(c, &d, DevicesListOptions{})
		},
		"APIKeys": func(c *ui.Context) {
			k := []APIKey{{ID: "a", Name: "A", Prefix: "p", Secret: "1"}, {ID: "a", Name: "B", Prefix: "p", Secret: "2"}}
			APIKeys(c, &k, APIKeysOptions{})
		},
		"Usage": func(c *ui.Context) {
			r := Resource{Name: "Storage", Unit: "GB", Used: 1, Limit: 2}
			Usage(c, []Resource{r, r}, UsageOptions{})
		},
		"InvoiceList": func(c *ui.Context) {
			i := Invoice{ID: "a", Number: "INV-1", Date: "Oct 1", Amount: "$1", URL: "https://a"}
			InvoiceList(c, []Invoice{i, i}, new(string), InvoiceListOptions{})
		},
	}
	for name, view := range cases {
		t.Run(name, func(t *testing.T) { panics(t, "account: "+name, view) })
	}
}

func TestSessionsCurrentAndEmpty(t *testing.T) {
	sessions := []Session{{ID: "a", Device: "MacBook Pro", Last: "now", Current: true}, {ID: "b", Device: "iPad Air", Last: "now"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Sessions(c, &sessions, SessionsOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "This device") || tt.HasText("Sign out MacBook Pro") {
		t.Errorf("current session: %q", tt.Texts())
	}
	devices := []Device{{ID: "d", Name: "Pixel 8", Last: "Oct 19"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { DevicesList(c, &devices, DevicesListOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "Last seen Oct 19") {
		t.Errorf("seen: %q", tt.Texts())
	}
	click(t, tt, "Remove Pixel 8")
	settles(t, tt)
	if len(devices) != 0 || !exact(tt, "No remembered devices") {
		t.Errorf("empty: %v %q", devices, tt.Texts())
	}
}

func TestNotesClearOnEdit(t *testing.T) {
	email := "maya.chen@kiln.co"
	tt := ui.NewTester(padded(func(c *ui.Context) { ForgotPassword(c, &email, ForgotPasswordOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Send link")
	if !exact(tt, "Link sent. Check your inbox.") {
		t.Fatalf("no sent note: %q", tt.Texts())
	}
	click(t, tt, "Email")
	tt.Type("x")
	tt.Frame()
	if exact(tt, "Link sent. Check your inbox.") {
		t.Error("ForgotPassword kept its note after an edit")
	}

	next := "maya@kiln.studio"
	tt = ui.NewTester(padded(func(c *ui.Context) { ChangeEmail(c, &next, ChangeEmailOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Verify")
	if !exact(tt, "We sent a confirmation link to maya@kiln.studio.") {
		t.Fatalf("no sent note: %q", tt.Texts())
	}
	next = "maya@kiln.shop"
	tt.Frame()
	if tt.HasText("We sent a confirmation link") {
		t.Error("ChangeEmail kept its note for another address")
	}
}

func TestResetMismatchAndTwoFactorVerified(t *testing.T) {
	value := ResetValue{Password: "Cone6&Oxide", Confirm: "Cone6&Oxid"}
	tt := ui.NewTester(padded(func(c *ui.Context) { ResetPassword(c, &value, ResetPasswordOptions{}) }), 500, 500)
	settles(t, tt)
	if !exact(tt, "Passwords do not match") || !exact(tt, "Strong") {
		t.Errorf("reset: %q", tt.Texts())
	}
	code, resent, verified := "482915", 0, false
	tt = ui.NewTester(padded(func(c *ui.Context) { TwoFactor(c, &code, &resent, &verified, TwoFactorOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Verify")
	if !exact(tt, "Verified") {
		t.Errorf("two factor: %q", tt.Texts())
	}
}

func TestProfileBioLimit(t *testing.T) {
	profile := Profile{Name: "Maya Chen", Handle: "@mayachen", Bio: strings.Repeat("a", 161)}
	saved := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, save := ProfileSettings(c, &profile, ProfileSettingsOptions{})
		saved = saved || save.Changed()
	}), 500, 600)
	settles(t, tt)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if saved || !exact(tt, "161/160") {
		t.Errorf("over-limit bio saved %v: %q", saved, tt.Texts())
	}
	profile.Bio = "Potter"
	tt.Frame()
	click(t, tt, "Save")
	if !saved || !exact(tt, "Profile saved") {
		t.Errorf("save: %v %q", saved, tt.Texts())
	}
}

func TestSecurityBands(t *testing.T) {
	steps := []SecurityStep{{Title: "Strong password", Done: true}, {Title: "Add a passkey"}}
	for score, want := range map[int]string{49: "At risk", 50: "Fair", 79: "Fair", 80: "Strong"} {
		tt := ui.NewTester(padded(func(c *ui.Context) { SecurityScore(c, score, steps, SecurityScoreOptions{}) }), 500, 400)
		settles(t, tt)
		if !exact(tt, want) || !exact(tt, "1 of 2 done") {
			t.Errorf("score %d: %q", score, tt.Texts())
		}
	}
}

func TestUsagePercents(t *testing.T) {
	cases := []struct {
		used, limit float64
		want        string
		over        bool
	}{{0.29, 1, "29%", false}, {0.996, 1, "99%", false}, {5, 5, "100%", false}, {6, 5, "120%", true}}
	for _, tc := range cases {
		r := []Resource{{Name: "Storage", Unit: "GB", Used: tc.used, Limit: tc.limit}}
		tt := ui.NewTester(padded(func(c *ui.Context) { Usage(c, r, UsageOptions{}) }), 500, 300)
		settles(t, tt)
		if !exact(tt, tc.want) || exact(tt, "Over limit") != tc.over {
			t.Errorf("%v of %v: %q", tc.used, tc.limit, tt.Texts())
		}
	}
	panics(t, `Usage "Storage" share overflows`, func(c *ui.Context) {
		Usage(c, []Resource{{Name: "Storage", Unit: "GB", Used: 1e308, Limit: 1e-308}}, UsageOptions{})
	})
}

func TestCopiedAndResentNotes(t *testing.T) {
	copied := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		BackupCodes(c, []string{"7KQ2-M9XD", "P4WT-8LNC"}, &copied, BackupCodesOptions{})
	}), 500, 400)
	settles(t, tt)
	click(t, tt, "P4WT-8LNC")
	if !exact(tt, "Code 02 copied") {
		t.Errorf("backup: %q", tt.Texts())
	}
	resent := 0
	tt = ui.NewTester(padded(func(c *ui.Context) { MagicLink(c, "maya.chen@kiln.co", &resent, MagicLinkOptions{}) }), 500, 500)
	settles(t, tt)
	click(t, tt, "Resend")
	if !exact(tt, "Sent again ×1") {
		t.Errorf("magic link: %q", tt.Texts())
	}
	downloaded := ""
	invoices := []Invoice{{ID: "a", Number: "INV-1", Date: "Oct 1", Amount: "$24.00", URL: "https://a"},
		{ID: "b", Number: "INV-2", Date: "Sep 1", Amount: "$24.00", URL: "https://b", Paid: true}}
	tt = ui.NewTester(padded(func(c *ui.Context) { InvoiceList(c, invoices, &downloaded, InvoiceListOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "Due") || !exact(tt, "Paid") {
		t.Errorf("invoice status: %q", tt.Texts())
	}
}
