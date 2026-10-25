package account

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestTeamListsKeepIdentity(t *testing.T) {
	cards := []Card{{ID: "v", Brand: "Visa", Last4: "4242", Expiry: "08/28"}, {ID: "m", Brand: "Mastercard", Last4: "5555", Expiry: "01/27"}}
	chosen := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { PaymentMethods(c, cards, &chosen, PaymentMethodsOptions{}) }), 500, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Mastercard ending 5555", func() { cards[0], cards[1] = cards[1], cards[0] })
	if chosen == "v" {
		t.Error("PaymentMethods chose the card that slid in")
	}

	members := []Member{{ID: "a", Name: "Daniel Okafor", Email: "d@k.co", Role: "Admin"}, {ID: "b", Name: "Priya Nair", Email: "p@k.co", Role: "Editor"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { TeamMembers(c, &members, TeamMembersOptions{}) }), 500, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Remove Priya Nair", func() { members[0], members[1] = members[1], members[0] })
	if len(members) == 1 && members[0].ID == "b" {
		t.Error("TeamMembers removed the row that slid in")
	}

	steps := []Step{{ID: "a", Title: "Add your logo"}, {ID: "b", Title: "Invite your team"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Onboarding(c, &steps, OnboardingOptions{}) }), 500, 500)
	settles(t, tt)
	swapUnderPress(t, tt, "Invite your team", func() { steps[0], steps[1] = steps[1], steps[0] })
	for _, s := range steps {
		if s.ID == "a" && s.Done {
			t.Error("Onboarding ticked the row that slid in")
		}
	}

	settings := []Setting{{ID: "a", Title: "Show online status", Detail: "x"}, {ID: "b", Title: "Share usage data", Detail: "y"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Privacy(c, &settings, PrivacyOptions{}) }), 500, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Share usage data", func() { settings[0], settings[1] = settings[1], settings[0] })
	for _, s := range settings {
		if s.ID == "a" && s.On {
			t.Error("Privacy toggled the row that slid in")
		}
	}
}

func TestTeamRepeatsPanic(t *testing.T) {
	cases := map[string]func(c *ui.Context){
		"PaymentMethods": func(c *ui.Context) {
			k := Card{ID: "a", Brand: "Visa", Last4: "4242", Expiry: "1/28"}
			PaymentMethods(c, []Card{k, k}, new(string), PaymentMethodsOptions{})
		},
		"TeamMembers": func(c *ui.Context) {
			m := []Member{{ID: "a", Name: "A", Email: "a@k.co", Role: "R"}, {ID: "a", Name: "B", Email: "b@k.co", Role: "R"}}
			TeamMembers(c, &m, TeamMembersOptions{})
		},
		"RoleSelect": func(c *ui.Context) {
			r := Role{Name: "Admin", Description: "All"}
			RoleSelect(c, "Role", []Role{r, r}, new(string))
		},
		"PendingInvites": func(c *ui.Context) {
			i := Invite{ID: "a", Email: "a@k.co", Role: "R", Sent: "now"}
			v := []Invite{i, i}
			PendingInvites(c, &v, PendingInvitesOptions{})
		},
		"OrgSwitcher": func(c *ui.Context) {
			o := Org{ID: "a", Name: "Kiln", Plan: "Free"}
			cur := "a"
			OrgSwitcher(c, []Org{o, o}, &cur, OrgSwitcherOptions{})
		},
		"NotifyPrefs": func(c *ui.Context) {
			tp := []Topic{{ID: "a", Name: "Orders"}, {ID: "a", Name: "Comments"}}
			NotifyPrefs(c, &tp, NotifyPrefsOptions{})
		},
		"ConnectedApps": func(c *ui.Context) {
			a := App{ID: "a", Name: "Slack", Access: "x", Since: "y"}
			v := []App{a, a}
			ConnectedApps(c, &v, ConnectedAppsOptions{})
		},
		"AuditLog": func(c *ui.Context) {
			e := AuditEvent{ID: "a", Actor: "Maya", Action: "x", Time: "now"}
			AuditLog(c, []AuditEvent{e, e}, AuditLogOptions{})
		},
		"AccountMenu": func(c *ui.Context) {
			it := MenuItem{ID: "a", Label: "Profile", Icon: "user"}
			AccountMenu(c, "Maya", "m@k.co", []MenuItem{it, it}, new(string), AccountMenuOptions{})
		},
	}
	for name, view := range cases {
		t.Run(name, func(t *testing.T) { panics(t, "account: "+name, view) })
	}
	panics(t, "signout is reserved", func(c *ui.Context) {
		AccountMenu(c, "Maya", "m@k.co", []MenuItem{{ID: "signout", Label: "Out", Icon: "x"}}, new(string), AccountMenuOptions{})
	})
	panics(t, "is not an organization", func(c *ui.Context) {
		cur := "z"
		OrgSwitcher(c, []Org{{ID: "a", Name: "Kiln", Plan: "Free"}}, &cur, OrgSwitcherOptions{})
	})
	panics(t, "ready without a URL", func(c *ui.Context) {
		s := ExportReady
		DataExport(c, &s, 1, DataExportOptions{})
	})
}

func TestOwnerStaysAndInvitesResendOnce(t *testing.T) {
	members := []Member{{ID: "a", Name: "Maya Chen", Email: "m@k.co", Role: "Owner", Owner: true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { TeamMembers(c, &members, TeamMembersOptions{}) }), 500, 300)
	settles(t, tt)
	if tt.HasText("Remove Maya Chen") {
		t.Error("the owner can be removed")
	}
	invites := []Invite{{ID: "a", Email: "jonas@berg.se", Role: "Editor", Sent: "now"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { PendingInvites(c, &invites, PendingInvitesOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Resend jonas@berg.se")
	click(t, tt, "Resend jonas@berg.se")
	if !exact(tt, "Sent") {
		t.Errorf("no sent state: %q", tt.Texts())
	}
	cues(t, before, sound.Send)
	click(t, tt, "Revoke jonas@berg.se")
	settles(t, tt)
	if len(invites) != 0 || !exact(tt, "No pending invites") {
		t.Errorf("revoke: %v %q", invites, tt.Texts())
	}
}

func TestDangerZoneNeedsTypedName(t *testing.T) {
	action := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { DangerZone(c, "kiln-co", &action, DangerZoneOptions{}) }), 500, 600)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Delete kiln-co")
	if _, ok := tt.Find("Type kiln-co to confirm"); !ok {
		t.Fatalf("delete did not arm: %q", tt.Texts())
	}
	click(t, tt, "Type kiln-co to confirm")
	tt.Type("kiln")
	tt.Frame()
	if err := tt.Click("Delete forever"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if action != "" {
		t.Fatalf("deleted with a partial name: %q", action)
	}
	click(t, tt, "Type kiln-co to confirm")
	tt.Type("-co")
	tt.Frame()
	click(t, tt, "Delete forever")
	if action != "delete" {
		t.Errorf("action %q", action)
	}
	cues(t, before, sound.Delete)
}

func TestDataExportAndMenus(t *testing.T) {
	status := ExportIdle
	tt := ui.NewTester(padded(func(c *ui.Context) {
		DataExport(c, &status, 0.4, DataExportOptions{URL: "https://kiln.co/x.zip"})
	}), 500, 400)
	settles(t, tt)
	click(t, tt, "Request export")
	tt.Advance(time.Second)
	if status != ExportPreparing || !exact(tt, "40%") {
		t.Fatalf("preparing: %v %q", status, tt.Texts())
	}
	status = ExportReady
	tt.Frame()
	settles(t, tt)
	click(t, tt, "Download ZIP")
	if urls := tt.OpenedURLs(); len(urls) != 1 || urls[0] != "https://kiln.co/x.zip" {
		t.Errorf("opened %v", urls)
	}

	chosen := ""
	items := []MenuItem{{ID: "settings", Label: "Settings", Icon: "settings"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { AccountMenu(c, "Maya Chen", "m@k.co", items, &chosen, AccountMenuOptions{}) }), 500, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Account menu")
	click(t, tt, "Sign out")
	if chosen != "signout" {
		t.Errorf("chosen %q", chosen)
	}
	cues(t, before, sound.Open, sound.Close)

	orgs := []Org{{ID: "k", Name: "Kiln & Co.", Plan: "Studio"}, {ID: "b", Name: "Berg Furniture", Plan: "Pro"}}
	current := "k"
	tt = ui.NewTester(padded(func(c *ui.Context) { OrgSwitcher(c, orgs, &current, OrgSwitcherOptions{}) }), 500, 500)
	settles(t, tt)
	if !exact(tt, "KC") {
		t.Errorf("initials skip symbols: %q", tt.Texts())
	}
	click(t, tt, "Switch organization")
	click(t, tt, "Berg Furniture")
	if current != "b" {
		t.Errorf("current %q", current)
	}
}

func TestOnboardingWelcomeAndLockout(t *testing.T) {
	steps := []Step{{ID: "a", Title: "Add your logo", Done: true}, {ID: "b", Title: "Invite your team"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Onboarding(c, &steps, OnboardingOptions{}) }), 500, 500)
	settles(t, tt)
	click(t, tt, "Invite your team")
	if !exact(tt, "You're all set") || !exact(tt, "2 of 2 done") {
		t.Errorf("onboarding: %q", tt.Texts())
	}

	page := 0
	slides := []Slide{{Icon: "sparkles", Title: "Welcome", Text: "Hi"}, {Icon: "users", Title: "Team", Text: "Invite"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Welcome(c, slides, &page, WelcomeOptions{}) }), 500, 500)
	settles(t, tt)
	click(t, tt, "Next")
	if page != 1 || tt.HasText("Skip") {
		t.Fatalf("page %d texts %q", page, tt.Texts())
	}
	click(t, tt, "Get started")
	if page != 2 || !exact(tt, "Welcome aboard") {
		t.Errorf("finish: %d %q", page, tt.Texts())
	}

	var start time.Time
	tt = ui.NewTester(padded(func(c *ui.Context) {
		if start.IsZero() {
			start = c.Now()
		}
		Lockout(c, start, start.Add(5*time.Minute), LockoutOptions{})
	}), 500, 500)
	tt.Frame()
	if !exact(tt, "05:00") {
		t.Fatalf("start: %q", tt.Texts())
	}
	tt.Advance(61500 * time.Millisecond)
	if !exact(tt, "03:59") {
		t.Errorf("after 61s: %q", tt.Texts())
	}
	tt.Advance(5 * time.Minute)
	if !exact(tt, "You can try again now") {
		t.Errorf("unlocked: %q", tt.Texts())
	}

	start = time.Time{}
	tt = ui.NewTester(padded(func(c *ui.Context) {
		if start.IsZero() {
			start = c.Now()
		}
		Lockout(c, start, start.Add(90*time.Second+time.Nanosecond), LockoutOptions{})
	}), 500, 500)
	tt.Frame()
	if !exact(tt, "01:31") {
		t.Errorf("a part second rounds up: %q", tt.Texts())
	}
}
