package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestAccount2DemosLog(t *testing.T) {
	press := func(names ...string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, name := range names {
				if err := tt.Click(name); err != nil {
					return err
				}
				tt.Frame()
			}
			return nil
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"PaymentMethods": {press("Mastercard ending 5555"), "default mastercard"},
		"TeamMembers":    {press("Remove Lucas Martin"), "members 4"},
		"RoleSelect":     {press("Viewer"), "role Viewer"},
		"PendingInvites": {press("Revoke noah.kim@kiln.co"), "invites 2"},
		"OrgSwitcher":    {press("Switch organization", "Studio Potters"), "org potters"},
		"DangerZone":     {press("Transfer kiln-co"), "action transfer"},
		"DataExport":     {press("Request export"), "status 1"},
		"Privacy":        {press("Show online status"), "on 1"},
		"NotifyPrefs":    {press("Push Product updates"), "channels 8"},
		"Appearance":     {press("Dark"), "theme 1 density 1"},
		"LanguagePrefs":  {press("Language", "Deutsch"), "language Deutsch"},
		"Accessibility":  {press("Reduce motion"), "motion true"},
		"ConnectedApps":  {press("Disconnect Figma"), "apps 3"},
		"Onboarding":     {press("Connect a payment method"), "done 3"},
		"Welcome":        {press("Next"), "page 1"},
		"VerifyBanner":   {press("Resend email"), "resent 1"},
		"AccountMenu":    {press("Account menu", "Settings"), "chose settings"},
		"SignOut":        {press("Sign in again"), "signin"},
	}
	quiet := map[string]bool{"AuditLog": true, "Lockout": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "account" || e.ID < 1246 {
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
		t.Errorf("checked %d second account demos", seen)
	}
}
