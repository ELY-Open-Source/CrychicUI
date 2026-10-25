package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/account"
	"github.com/egoist/mygo/ui"
)

var accountEntries = []Entry{
	{ID: 1226, Name: "SignIn", Category: "account", Summary: "Email, password, providers, forgot; cue Tap, Error.", Demo: demoAccountSignIn},
	{ID: 1227, Name: "SignUp", Category: "account", Summary: "Name, email, password, terms; cue Success.", Demo: demoAccountSignUp},
	{ID: 1228, Name: "ForgotPassword", Category: "account", Summary: "Email and send link; cue Send.", Demo: demoAccountForgotPassword},
	{ID: 1229, Name: "ResetPassword", Category: "account", Summary: "New password twice with strength; cue Success.", Demo: demoAccountResetPassword},
	{ID: 1230, Name: "Providers", Category: "account", Summary: "Provider buttons in a column; cue Tap.", Demo: demoAccountProviders},
	{ID: 1231, Name: "MagicLink", Category: "account", Summary: "Check your email state with resend; cue Send.", Demo: demoAccountMagicLink},
	{ID: 1232, Name: "TwoFactor", Category: "account", Summary: "Six-digit code with resend; cue Success.", Demo: demoAccountTwoFactor},
	{ID: 1233, Name: "BackupCodes", Category: "account", Summary: "Codes grid with copy; cue Pop.", Demo: demoAccountBackupCodes},
	{ID: 1234, Name: "Passkey", Category: "account", Summary: "Add passkey prompt; cue Success.", Demo: demoAccountPasskey},
	{ID: 1235, Name: "Sessions", Category: "account", Summary: "Sessions with sign out; cue Close.", Demo: demoAccountSessions},
	{ID: 1236, Name: "DevicesList", Category: "account", Summary: "Devices with last seen and remove; cue Delete.", Demo: demoAccountDevicesList},
	{ID: 1237, Name: "ProfileSettings", Category: "account", Summary: "Name, handle, bio; cue Success.", Demo: demoAccountProfileSettings},
	{ID: 1238, Name: "AvatarEditor", Category: "account", Summary: "Avatar with change and remove; cue Pop.", Demo: demoAccountAvatarEditor},
	{ID: 1239, Name: "ChangeEmail", Category: "account", Summary: "New email with verify; cue Send.", Demo: demoAccountChangeEmail},
	{ID: 1240, Name: "ChangePassword", Category: "account", Summary: "Current and new; cue Success.", Demo: demoAccountChangePassword},
	{ID: 1241, Name: "SecurityScore", Category: "account", Summary: "Score ring with steps.", Demo: demoAccountSecurityScore},
	{ID: 1242, Name: "APIKeys", Category: "account", Summary: "Keys masked with copy and revoke; cue Pop, Delete.", Demo: demoAccountAPIKeys},
	{ID: 1243, Name: "PlanCard", Category: "account", Summary: "Current plan with change; cue Tap.", Demo: demoAccountPlanCard},
	{ID: 1244, Name: "Usage", Category: "account", Summary: "Usage meters by resource.", Demo: demoAccountUsage},
	{ID: 1245, Name: "InvoiceList", Category: "account", Summary: "Invoices with download; cue Tap.", Demo: demoAccountInvoiceList},
}

func demoAccountSignIn(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	value := ui.Local(root, "credentials", func() account.Credentials {
		return account.Credentials{Email: "maya.chen@kiln.co", Password: "glaze-firing-24"}
	})
	action := textState(root, "action", "")
	root.Children(func() {
		form, submit := account.SignIn(c, value, []string{"Google", "Apple"}, action, account.SignInOptions{
			Subtitle: "Sign in to Kiln & Co."})
		if submit.Changed() {
			log("signin")
		} else if form.Changed() {
			if *action != "" {
				log("action " + *action)
			} else {
				log("edit")
			}
		}
	})
}

func demoAccountSignUp(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	value := ui.Local(root, "signup", func() account.SignUpValue {
		return account.SignUpValue{Name: "Maya Chen", Email: "maya.chen@kiln.co", Password: "spring-glaze-24", Terms: true}
	})
	root.Children(func() {
		terms := value.Terms
		form, submit := account.SignUp(c, value, account.SignUpOptions{})
		switch {
		case submit.Changed():
			log("created")
		case form.Changed() && terms != value.Terms:
			log("terms " + fmt.Sprint(value.Terms))
		case form.Changed():
			log(fmt.Sprintf("nameRunes %d emailRunes %d passwordRunes %d",
				len([]rune(value.Name)), len([]rune(value.Email)), len([]rune(value.Password))))
		}
	})
}

func demoAccountForgotPassword(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	email := textState(root, "email", "maya.chen@kiln.co")
	root.Children(func() {
		form, send := account.ForgotPassword(c, email, account.ForgotPasswordOptions{})
		if send.Changed() {
			log("sent")
		} else if form.Changed() {
			log("emailRunes " + fmt.Sprint(len([]rune(*email))))
		}
	})
}

func demoAccountResetPassword(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	value := ui.Local(root, "reset", func() account.ResetValue {
		return account.ResetValue{Password: "Cone6&Oxide", Confirm: "Cone6&Oxide"}
	})
	root.Children(func() {
		form, submit := account.ResetPassword(c, value, account.ResetPasswordOptions{})
		if submit.Changed() {
			log("reset")
		} else if form.Changed() {
			log(fmt.Sprintf("passwordRunes %d confirmRunes %d", len([]rune(value.Password)), len([]rune(value.Confirm))))
		}
	})
}

func demoAccountProviders(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	chosen := textState(root, "provider", "")
	root.Children(func() {
		providers := []string{"Google", "GitHub", "Apple", "Microsoft"}
		if account.Providers(c, providers, chosen, account.ProvidersOptions{}).Changed() {
			log("provider " + *chosen)
		}
	})
}

func demoAccountMagicLink(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	resent := intState(root, "resent", 0)
	root.Children(func() {
		if account.MagicLink(c, "maya.chen@kiln.co", resent, account.MagicLinkOptions{}).Changed() {
			log("resent " + fmt.Sprint(*resent))
		}
	})
}

func demoAccountTwoFactor(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	code := textState(root, "code", "482915")
	resent, verified := intState(root, "resent", 0), boolState(root, "verified", false)
	root.Children(func() {
		sent := *resent
		form, verify := account.TwoFactor(c, code, resent, verified, account.TwoFactorOptions{})
		switch {
		case verify.Changed():
			log("verified")
		case form.Changed() && sent != *resent:
			log("resent " + fmt.Sprint(*resent))
		case form.Changed():
			log("codeRunes " + fmt.Sprint(len([]rune(*code))))
		}
	})
}

func demoAccountBackupCodes(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	copied := intState(root, "copied", -1)
	codes := []string{"7KQ2-M9XD", "P4WT-8LNC", "B6HR-2ZQF", "XN3V-K7MJ", "T9FD-4CWE", "H2LY-6RPA"}
	root.Children(func() {
		if account.BackupCodes(c, codes, copied, account.BackupCodesOptions{}).Changed() {
			log("copied " + fmt.Sprint(*copied))
		}
	})
}

func demoAccountPasskey(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	added := boolState(root, "added", false)
	root.Children(func() {
		if account.Passkey(c, added, account.PasskeyOptions{}).Changed() {
			log("added " + fmt.Sprint(*added))
		}
	})
}

func demoAccountSessions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	sessions := ui.Local(root, "sessions", func() []account.Session {
		return []account.Session{
			{ID: "mac", Device: "MacBook Pro", Location: "Portland, OR", Last: "Active now", Icon: "laptop", Current: true},
			{ID: "iphone", Device: "iPhone 15", Location: "Portland, OR", Last: "2 hours ago", Icon: "smartphone"},
			{ID: "ipad", Device: "iPad Air", Location: "Seattle, WA", Last: "Yesterday", Icon: "tablet"},
			{ID: "pc", Device: "Windows PC", Location: "Chicago, IL", Last: "3 days ago", Icon: "monitor"},
		}
	})
	root.Children(func() {
		if account.Sessions(c, sessions, account.SessionsOptions{}).Changed() {
			log("sessions " + fmt.Sprint(len(*sessions)))
		}
	})
}

func demoAccountDevicesList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	devices := ui.Local(root, "devices", func() []account.Device {
		return []account.Device{
			{ID: "imac", Name: "Studio iMac", Last: "today, 09:12", Icon: "monitor"},
			{ID: "air", Name: "MacBook Air", Last: "yesterday", Icon: "laptop"},
			{ID: "pixel", Name: "Pixel 8", Last: "Oct 19", Icon: "smartphone"},
			{ID: "ipad", Name: "Kiln room iPad", Last: "Sep 30", Icon: "tablet"},
		}
	})
	root.Children(func() {
		if account.DevicesList(c, devices, account.DevicesListOptions{}).Changed() {
			log("devices " + fmt.Sprint(len(*devices)))
		}
	})
}

func demoAccountProfileSettings(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	profile := ui.Local(root, "profile", func() account.Profile {
		return account.Profile{Name: "Maya Chen", Handle: "@mayachen", Bio: "Product designer at Kiln & Co. Lamps, mugs and slow mornings."}
	})
	root.Children(func() {
		form, save := account.ProfileSettings(c, profile, account.ProfileSettingsOptions{})
		if save.Changed() {
			log("saved")
		} else if form.Changed() {
			log("bioRunes " + fmt.Sprint(len([]rune(profile.Bio))))
		}
	})
}

func demoAccountAvatarEditor(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	image := ui.Local(root, "image", func() ui.ImageSource { return photos()[0] })
	change := boolState(root, "change", false)
	root.Children(func() {
		if account.AvatarEditor(c, image, change, account.AvatarEditorOptions{}).Changed() {
			log(fmt.Sprintf("image %v change %v", *image != nil, *change))
		}
	})
}

func demoAccountChangeEmail(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	email := textState(root, "email", "maya@kiln.studio")
	root.Children(func() {
		form, verify := account.ChangeEmail(c, email, account.ChangeEmailOptions{Current: "Now maya.chen@kiln.co"})
		if verify.Changed() {
			log("verify")
		} else if form.Changed() {
			log("emailRunes " + fmt.Sprint(len([]rune(*email))))
		}
	})
}

func demoAccountChangePassword(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	value := ui.Local(root, "change", func() account.PasswordChange {
		return account.PasswordChange{Current: "glaze-firing-24", New: "cone6-oxidation"}
	})
	root.Children(func() {
		form, submit := account.ChangePassword(c, value, account.ChangePasswordOptions{})
		if submit.Changed() {
			log("changed")
		} else if form.Changed() {
			log(fmt.Sprintf("currentRunes %d newRunes %d", len([]rune(value.Current)), len([]rune(value.New))))
		}
	})
}

func demoAccountSecurityScore(c *ui.Context, log func(string)) {
	steps := []account.SecurityStep{{Title: "Use a strong password", Done: true}, {Title: "Verify your email", Done: true},
		{Title: "Turn on two-step verification", Done: true}, {Title: "Add a passkey"}, {Title: "Review signed-in devices"}}
	ui.Column(c).Width(420).Children(func() { account.SecurityScore(c, 72, steps, account.SecurityScoreOptions{}) })
}

func demoAccountAPIKeys(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	keys := ui.Local(root, "keys", func() []account.APIKey {
		return []account.APIKey{
			{ID: "prod", Name: "Production", Prefix: "ck_live_", Secret: "ck_live_9fQ2xLm7", Created: "Created Mar 4"},
			{ID: "stage", Name: "Staging", Prefix: "ck_test_", Secret: "ck_test_4hN8pRw1", Created: "Created Jun 18"},
			{ID: "ci", Name: "CI pipeline", Prefix: "ck_live_", Secret: "ck_live_2bV6tKs3", Created: "Created Jan 9", Revoked: true},
		}
	})
	revoked := func() int {
		n := 0
		for _, k := range *keys {
			if k.Revoked {
				n++
			}
		}
		return n
	}
	root.Children(func() {
		before := revoked()
		if account.APIKeys(c, keys, account.APIKeysOptions{}).Changed() {
			if after := revoked(); after != before {
				log("revoked " + fmt.Sprint(after))
			} else {
				log("copied")
			}
		}
	})
}

func demoAccountPlanCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	change := boolState(root, "change", false)
	plan := account.Plan{Name: "Studio", Price: "$24", Period: "/ month", Note: "Renews Nov 12, 2026",
		Features: []string{"5 team seats", "50 GB of storage", "Custom domain for your shop", "Priority support"}}
	root.Children(func() {
		if account.PlanCard(c, plan, true, change, account.PlanCardOptions{}).Changed() {
			log("change " + fmt.Sprint(*change))
		}
	})
}

func demoAccountUsage(c *ui.Context, log func(string)) {
	resources := []account.Resource{
		{Name: "Storage", Unit: "GB", Icon: "hard-drive", Used: 42.6, Limit: 50},
		{Name: "Bandwidth", Unit: "GB", Icon: "cloud", Used: 312, Limit: 500},
		{Name: "Team seats", Unit: "seats", Icon: "users", Used: 5, Limit: 5},
		{Name: "Product listings", Unit: "items", Icon: "package", Used: 86, Limit: 250},
	}
	ui.Column(c).Width(420).Children(func() { account.Usage(c, resources, account.UsageOptions{}) })
}

func demoAccountInvoiceList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	downloaded := textState(root, "downloaded", "")
	invoice := func(n int, date string, paid bool) account.Invoice {
		number := fmt.Sprintf("INV-2026-%03d", n)
		return account.Invoice{ID: number, Number: number, Date: date, Amount: "$24.00", Paid: paid,
			URL: "https://billing.kiln.co/invoices/" + number + ".pdf"}
	}
	invoices := []account.Invoice{invoice(10, "Oct 1, 2026", false), invoice(9, "Sep 1, 2026", true),
		invoice(8, "Aug 1, 2026", true), invoice(7, "Jul 1, 2026", true)}
	root.Children(func() {
		if account.InvoiceList(c, invoices, downloaded, account.InvoiceListOptions{}).Changed() {
			log("downloaded " + *downloaded)
		}
	})
}
