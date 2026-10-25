package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestAccountDemoSecondaryActions(t *testing.T) {
	cases := []struct {
		name, action, want string
	}{{"SignUp", "Accept terms", "terms false"}, {"TwoFactor", "Resend", "resent 1"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := Entries()[tcIndex(tc.name)]
			var logs []string
			tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
			tt.Advance(2 * time.Second)
			if err := tt.Click(tc.action); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if !strings.Contains(strings.Join(logs, "\n"), tc.want) {
				t.Errorf("logs %q, want %q", logs, tc.want)
			}
		})
	}
}

func TestAccountDemoTextEdits(t *testing.T) {
	cases := []struct {
		name, label, want string
	}{{"SignIn", "Email", "edit"}, {"ForgotPassword", "Email", "emailRunes"},
		{"ResetPassword", "New password", "passwordRunes"}, {"ProfileSettings", "Name", "bioRunes"},
		{"ChangeEmail", "New email", "emailRunes"}, {"ChangePassword", "Current password", "currentRunes"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := Entries()[tcIndex(tc.name)]
			var logs []string
			tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
			tt.Advance(2 * time.Second)
			if err := tt.Click(tc.label); err != nil {
				t.Fatal(err)
			}
			tt.Type("x")
			tt.Frame()
			joined := strings.Join(logs, "\n")
			if !strings.Contains(joined, tc.want) {
				t.Errorf("logs %q, want %q", logs, tc.want)
			}
			if strings.Contains(joined, "kiln.co") || strings.Contains(joined, "Cone6&Oxide") {
				t.Errorf("typed text logged: %q", logs)
			}
		})
	}
}

func tcIndex(name string) int {
	for i, e := range Entries() {
		if e.Name == name {
			return i
		}
	}
	panic("missing account demo " + name)
}

func TestAccountDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"SignIn":          {press("Sign in"), "signin"},
		"SignUp":          {press("Create account"), "created"},
		"ForgotPassword":  {press("Send link"), "sent"},
		"ResetPassword":   {press("Reset password"), "reset"},
		"Providers":       {press("GitHub"), "provider GitHub"},
		"MagicLink":       {press("Resend"), "resent 1"},
		"TwoFactor":       {press("Verify"), "verified"},
		"BackupCodes":     {press("7KQ2-M9XD"), "copied 0"},
		"Passkey":         {press("Add passkey"), "added true"},
		"Sessions":        {press("Sign out iPad Air"), "sessions 3"},
		"DevicesList":     {press("Remove Pixel 8"), "devices 3"},
		"ProfileSettings": {press("Save"), "saved"},
		"AvatarEditor":    {press("Remove"), "image false"},
		"ChangeEmail":     {press("Verify"), "verify"},
		"ChangePassword":  {press("Update password"), "changed"},
		"APIKeys":         {press("Revoke Staging"), "revoked 2"},
		"PlanCard":        {press("Change plan"), "change true"},
		"InvoiceList":     {press("Download INV-2026-010"), "downloaded INV-2026-010"},
	}
	quiet := map[string]bool{"SecurityScore": true, "Usage": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "account" {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(2 * time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 20 {
		t.Errorf("checked %d account demos", seen)
	}
}

func TestAccountDemoEventsMatch(t *testing.T) {
	cases := []struct {
		name, action string
		typed        bool
		want, never  string
	}{{"APIKeys", "Copy Staging", false, "copied", "revoked"}, {"SignUp", "Name", true, "nameRunes", "terms"},
		{"TwoFactor", "Six-digit code", true, "codeRunes", "resent"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := Entries()[tcIndex(tc.name)]
			var logs []string
			tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
			tt.Advance(2 * time.Second)
			if err := tt.Click(tc.action); err != nil {
				t.Fatal(err)
			}
			if tc.typed {
				tt.Type("x")
			}
			tt.Frame()
			all := strings.Join(logs, "\n")
			if !strings.Contains(all, tc.want) || strings.Contains(all, tc.never) {
				t.Errorf("logs %q, want %q without %q", logs, tc.want, tc.never)
			}
		})
	}
}
