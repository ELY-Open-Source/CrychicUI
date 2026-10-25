package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/account"
	"github.com/egoist/mygo/ui"
)

var accountEntries2 = []Entry{
	{ID: 1246, Name: "PaymentMethods", Category: "account", Summary: "Cards with default mark; cue Select.", Demo: demoAccountPaymentMethods},
	{ID: 1247, Name: "TeamMembers", Category: "account", Summary: "Members with roles and remove; cue Delete.", Demo: demoAccountTeamMembers},
	{ID: 1248, Name: "RoleSelect", Category: "account", Summary: "Role with description; cue Select.", Demo: demoAccountRoleSelect},
	{ID: 1249, Name: "PendingInvites", Category: "account", Summary: "Invites with resend and revoke; cue Send.", Demo: demoAccountPendingInvites},
	{ID: 1250, Name: "OrgSwitcher", Category: "account", Summary: "Organizations with logos; cue Select.", Demo: demoAccountOrgSwitcher},
	{ID: 1251, Name: "DangerZone", Category: "account", Summary: "Delete account and transfer; cue Warning.", Demo: demoAccountDangerZone},
	{ID: 1252, Name: "DataExport", Category: "account", Summary: "Request export with status; cue Success.", Demo: demoAccountDataExport},
	{ID: 1253, Name: "Privacy", Category: "account", Summary: "Privacy toggles; cue Toggle.", Demo: demoAccountPrivacy},
	{ID: 1254, Name: "NotifyPrefs", Category: "account", Summary: "Email, push, SMS by topic; cue Toggle.", Demo: demoAccountNotifyPrefs},
	{ID: 1255, Name: "Appearance", Category: "account", Summary: "Theme and density choices; cue Select.", Demo: demoAccountAppearance},
	{ID: 1256, Name: "LanguagePrefs", Category: "account", Summary: "Language and region; cue Select.", Demo: demoAccountLanguagePrefs},
	{ID: 1257, Name: "Accessibility", Category: "account", Summary: "Motion, contrast, text size; cue Toggle.", Demo: demoAccountAccessibility},
	{ID: 1258, Name: "ConnectedApps", Category: "account", Summary: "Apps with disconnect; cue Close.", Demo: demoAccountConnectedApps},
	{ID: 1259, Name: "AuditLog", Category: "account", Summary: "Account events.", Demo: demoAccountAuditLog},
	{ID: 1260, Name: "Onboarding", Category: "account", Summary: "Getting started checklist; cue Toggle.", Demo: demoAccountOnboarding},
	{ID: 1261, Name: "Welcome", Category: "account", Summary: "Welcome screen with next; cue Tap.", Demo: demoAccountWelcome},
	{ID: 1262, Name: "VerifyBanner", Category: "account", Summary: "Verify email banner with resend; cue Send.", Demo: demoAccountVerifyBanner},
	{ID: 1263, Name: "AccountMenu", Category: "account", Summary: "Avatar menu with sign out; cue Open.", Demo: demoAccountAccountMenu},
	{ID: 1264, Name: "Lockout", Category: "account", Summary: "Too many attempts with countdown.", Demo: demoAccountLockout},
	{ID: 1265, Name: "SignOut", Category: "account", Summary: "Signed out state with sign in; cue Tap.", Demo: demoAccountSignOut},
}

func demoAccountPaymentMethods(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	chosen := textState(root, "chosen", "visa")
	cards := []account.Card{{ID: "visa", Brand: "Visa", Last4: "4242", Expiry: "08/28"},
		{ID: "mastercard", Brand: "Mastercard", Last4: "5555", Expiry: "01/27"},
		{ID: "amex", Brand: "Amex", Last4: "0005", Expiry: "11/29"}}
	root.Children(func() {
		if account.PaymentMethods(c, cards, chosen, account.PaymentMethodsOptions{}).Changed() {
			log("default " + *chosen)
		}
	})
}

func demoAccountTeamMembers(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	members := ui.Local(root, "members", func() []account.Member {
		return []account.Member{
			{ID: "maya", Name: "Maya Chen", Email: "maya.chen@kiln.co", Role: "Owner", Owner: true},
			{ID: "daniel", Name: "Daniel Okafor", Email: "daniel@kiln.co", Role: "Admin"},
			{ID: "priya", Name: "Priya Nair", Email: "priya@kiln.co", Role: "Editor"},
			{ID: "lucas", Name: "Lucas Martin", Email: "lucas@kiln.co", Role: "Viewer"},
			{ID: "sofia", Name: "Sofia Rossi", Email: "sofia@kiln.co", Role: "Editor"},
		}
	})
	root.Children(func() {
		if account.TeamMembers(c, members, account.TeamMembersOptions{}).Changed() {
			log("members " + fmt.Sprint(len(*members)))
		}
	})
}

func demoAccountRoleSelect(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	chosen := textState(root, "role", "Editor")
	roles := []account.Role{{Name: "Admin", Icon: "shield", Description: "Full access, including billing and team settings."},
		{Name: "Editor", Icon: "pencil", Description: "Create and edit products, posts and pages."},
		{Name: "Viewer", Icon: "eye", Description: "Read-only access to everything in the studio."}}
	root.Children(func() {
		if account.RoleSelect(c, "Role for Jonas Berg", roles, chosen).Changed() {
			log("role " + *chosen)
		}
	})
}

func demoAccountPendingInvites(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	invites := ui.Local(root, "invites", func() []account.Invite {
		return []account.Invite{{ID: "jonas", Email: "jonas@bergfurniture.se", Role: "Editor", Sent: "2 days ago"},
			{ID: "ana", Email: "ana.lima@kiln.co", Role: "Viewer", Sent: "yesterday"},
			{ID: "noah", Email: "noah.kim@kiln.co", Role: "Admin", Sent: "5 minutes ago"}}
	})
	root.Children(func() {
		before := len(*invites)
		if account.PendingInvites(c, invites, account.PendingInvitesOptions{}).Changed() {
			if len(*invites) != before {
				log("invites " + fmt.Sprint(len(*invites)))
			} else {
				log("resent")
			}
		}
	})
}

func demoAccountOrgSwitcher(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340)
	current := textState(root, "org", "kiln")
	orgs := []account.Org{{ID: "kiln", Name: "Kiln & Co.", Plan: "Studio plan · 5 seats"},
		{ID: "berg", Name: "Berg Furniture", Plan: "Pro plan · 12 seats"},
		{ID: "potters", Name: "Studio Potters", Plan: "Free plan"}}
	root.Children(func() {
		if account.OrgSwitcher(c, orgs, current, account.OrgSwitcherOptions{}).Changed() {
			log("org " + *current)
		}
	})
}

func demoAccountDangerZone(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	action := textState(root, "action", "")
	root.Children(func() {
		if account.DangerZone(c, "kiln-co", action, account.DangerZoneOptions{}).Changed() {
			log("action " + *action)
		}
	})
}

func demoAccountDataExport(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	status := ui.Local(root, "status", func() account.ExportStatus { return account.ExportIdle })
	started := ui.Local(root, "started", func() time.Time { return time.Time{} })
	progress := 0.0
	if *status == account.ExportPreparing {
		progress = min(1, float64(c.Now().Sub(*started))/float64(4*time.Second))
		if progress >= 1 {
			*status = account.ExportReady
		} else {
			c.After(100 * time.Millisecond)
		}
	}
	root.Children(func() {
		opts := account.DataExportOptions{URL: "https://kiln.co/exports/maya-chen-2026-10-25.zip"}
		if account.DataExport(c, status, progress, opts).Changed() {
			if *status == account.ExportPreparing {
				*started = c.Now()
			}
			log("status " + fmt.Sprint(int(*status)))
		}
	})
}

func demoAccountPrivacy(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	settings := ui.Local(root, "settings", func() []account.Setting {
		return []account.Setting{
			{ID: "search", Title: "Show profile in search", Detail: "Search engines can list your public shop.", Icon: "globe", On: true},
			{ID: "online", Title: "Show online status", Detail: "Teammates see when you are active.", Icon: "eye", On: true},
			{ID: "usage", Title: "Share usage data", Detail: "Help improve Kiln & Co. with anonymous stats.", Icon: "activity"},
			{ID: "ads", Title: "Personalized offers", Detail: "Tailor partner offers to your shop.", Icon: "megaphone"},
		}
	})
	root.Children(func() {
		if account.Privacy(c, settings, account.PrivacyOptions{}).Changed() {
			on := 0
			for _, s := range *settings {
				if s.On {
					on++
				}
			}
			log("on " + fmt.Sprint(on))
		}
	})
}

func demoAccountNotifyPrefs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	topics := ui.Local(root, "topics", func() []account.Topic {
		return []account.Topic{
			{ID: "orders", Name: "Orders", Detail: "New orders and refunds", Email: true, Push: true},
			{ID: "comments", Name: "Comments", Detail: "Replies and mentions", Push: true},
			{ID: "updates", Name: "Product updates", Detail: "Release notes, once a month", Email: true},
			{ID: "security", Name: "Security alerts", Detail: "Sign-ins and password changes", Email: true, Push: true, SMS: true},
		}
	})
	root.Children(func() {
		if account.NotifyPrefs(c, topics, account.NotifyPrefsOptions{}).Changed() {
			n := 0
			for _, tp := range *topics {
				for _, on := range []bool{tp.Email, tp.Push, tp.SMS} {
					if on {
						n++
					}
				}
			}
			log("channels " + fmt.Sprint(n))
		}
	})
}

func demoAccountAppearance(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	mode := ui.Local(root, "mode", func() account.ThemeMode { return account.ThemeSystem })
	density := ui.Local(root, "density", func() account.Density { return account.DensityComfortable })
	root.Children(func() {
		if account.Appearance(c, mode, density, account.AppearanceOptions{}).Changed() {
			log(fmt.Sprintf("theme %d density %d", *mode, *density))
		}
	})
}

func demoAccountLanguagePrefs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	lang, region := textState(root, "lang", "English (US)"), textState(root, "region", "United States")
	langs := []string{"English (US)", "English (UK)", "Deutsch", "Español", "Français", "日本語"}
	regions := []string{"United States", "United Kingdom", "Germany", "Spain", "France", "Japan"}
	root.Children(func() {
		if account.LanguagePrefs(c, langs, lang, regions, region, account.LanguagePrefsOptions{}).Changed() {
			log("language " + *lang + " region " + *region)
		}
	})
}

func demoAccountAccessibility(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	value := ui.Local(root, "a11y", func() account.A11y { return account.A11y{TextSize: 1} })
	root.Children(func() {
		if account.Accessibility(c, value, account.AccessibilityOptions{}).Changed() {
			log(fmt.Sprintf("motion %v contrast %v size %d", value.ReduceMotion, value.HighContrast, value.TextSize))
		}
	})
}

func demoAccountConnectedApps(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	apps := ui.Local(root, "apps", func() []account.App {
		return []account.App{
			{ID: "slack", Name: "Slack", Access: "Posts order alerts", Since: "Since Mar 2", Icon: "message-square"},
			{ID: "drive", Name: "Google Drive", Access: "Reads and writes files", Since: "Since Jan 14", Icon: "hard-drive"},
			{ID: "shopify", Name: "Shopify", Access: "Syncs products and stock", Since: "Since Sep 30", Icon: "shopping-bag"},
			{ID: "figma", Name: "Figma", Access: "Views design files", Since: "Since Jun 8", Icon: "pen-tool"},
		}
	})
	root.Children(func() {
		if account.ConnectedApps(c, apps, account.ConnectedAppsOptions{}).Changed() {
			log("apps " + fmt.Sprint(len(*apps)))
		}
	})
}

func demoAccountAuditLog(c *ui.Context, log func(string)) {
	events := []account.AuditEvent{
		{ID: "e5", Actor: "Maya Chen", Action: "signed in from Portland, OR", Time: "Today, 09:12", Kind: account.AuditSignIn},
		{ID: "e4", Actor: "Maya Chen", Action: "turned on two-step verification", Time: "Yesterday, 18:40", Kind: account.AuditSecurity},
		{ID: "e3", Actor: "Daniel Okafor", Action: "changed the shop name to Kiln & Co.", Time: "Oct 22, 14:05", Kind: account.AuditChange},
		{ID: "e2", Actor: "Maya Chen", Action: "upgraded to the Studio plan", Time: "Oct 12, 10:31", Kind: account.AuditBilling},
		{ID: "e1", Actor: "Priya Nair", Action: "signed in from Seattle, WA", Time: "Oct 11, 08:57", Kind: account.AuditSignIn},
	}
	ui.Column(c).Width(440).Children(func() { account.AuditLog(c, events, account.AuditLogOptions{}) })
}

func demoAccountOnboarding(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	steps := ui.Local(root, "steps", func() []account.Step {
		return []account.Step{
			{ID: "logo", Title: "Add your logo", Detail: "Shown on receipts and your shop header.", Done: true},
			{ID: "product", Title: "Create your first product", Done: true},
			{ID: "payment", Title: "Connect a payment method", Detail: "Get paid out every Friday."},
			{ID: "team", Title: "Invite your team"},
			{ID: "publish", Title: "Publish your shop"},
		}
	})
	root.Children(func() {
		if account.Onboarding(c, steps, account.OnboardingOptions{}).Changed() {
			n := 0
			for _, s := range *steps {
				if s.Done {
					n++
				}
			}
			log("done " + fmt.Sprint(n))
		}
	})
}

func demoAccountWelcome(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	page := intState(root, "page", 0)
	slides := []account.Slide{
		{Icon: "sparkles", Title: "Welcome to Kiln & Co.", Text: "Your studio's shop, orders and team in one place."},
		{Icon: "package", Title: "Add your products", Text: "Photos, prices and stock in a few taps."},
		{Icon: "users", Title: "Bring your team", Text: "Invite makers and give each the right role."},
	}
	root.Children(func() {
		if account.Welcome(c, slides, page, account.WelcomeOptions{}).Changed() {
			log("page " + fmt.Sprint(*page))
		}
	})
}

func demoAccountVerifyBanner(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	resent := intState(root, "resent", 0)
	root.Children(func() {
		if account.VerifyBanner(c, "maya.chen@kiln.co", resent, account.VerifyBannerOptions{}).Changed() {
			log("resent " + fmt.Sprint(*resent))
		}
	})
}

func demoAccountAccountMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).AlignItems(ui.End)
	chosen := textState(root, "chosen", "")
	items := []account.MenuItem{{ID: "profile", Label: "Profile", Icon: "user-round"}, {ID: "settings", Label: "Settings", Icon: "settings"},
		{ID: "billing", Label: "Billing", Icon: "credit-card"}, {ID: "shortcuts", Label: "Keyboard shortcuts", Icon: "keyboard"}}
	root.Children(func() {
		if account.AccountMenu(c, "Maya Chen", "maya.chen@kiln.co", items, chosen, account.AccountMenuOptions{}).Changed() {
			log("chose " + *chosen)
		}
	})
}

func demoAccountLockout(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	from := ui.Local(root, "from", func() time.Time { return c.Now().Add(-time.Minute) })
	root.Children(func() { account.Lockout(c, *from, from.Add(5*time.Minute), account.LockoutOptions{}) })
}

func demoAccountSignOut(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Children(func() {
		if _, signIn := account.SignOut(c, "Maya", account.SignOutOptions{}); signIn.Changed() {
			log("signin")
		}
	})
}
