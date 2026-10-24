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
	root := ui.Column(c)
	value := ui.Local(root, "credentials", func() account.Credentials {
		return account.Credentials{Email: "ana@example.com", Password: "password"}
	})
	action := textState(root, "action", "")
	root.Children(func() {
		form, submit := account.SignIn(c, value, []string{"Google"}, action, account.SignInOptions{})
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
	root := ui.Column(c)
	value := ui.Local(root, "signup", func() account.SignUpValue {
		return account.SignUpValue{Name: "Ana", Email: "ana@example.com", Password: "password", Terms: true}
	})
	root.Children(func() {
		form, submit := account.SignUp(c, value, account.SignUpOptions{})
		if submit.Changed() {
			log("created")
		} else if form.Changed() {
			log("terms " + fmt.Sprint(value.Terms))
		}
	})
}

func demoAccountForgotPassword(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	email := textState(root, "email", "ana@example.com")
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
	root := ui.Column(c)
	value := ui.Local(root, "reset", func() account.ResetValue {
		return account.ResetValue{Password: "Password1!", Confirm: "Password1!"}
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
	root := ui.Column(c)
	chosen := textState(root, "provider", "")
	root.Children(func() {
		if account.Providers(c, []string{"Google"}, chosen).Changed() {
			log("provider " + *chosen)
		}
	})
}

func demoAccountMagicLink(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	resent := intState(root, "resent", 0)
	root.Children(func() {
		if account.MagicLink(c, "ana@example.com", resent, account.MagicLinkOptions{}).Changed() {
			log("resent " + fmt.Sprint(*resent))
		}
	})
}

func demoAccountTwoFactor(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	code := textState(root, "code", "123456")
	resent, verified := intState(root, "resent", 0), boolState(root, "verified", false)
	root.Children(func() {
		form, verify := account.TwoFactor(c, code, resent, verified, account.TwoFactorOptions{})
		if verify.Changed() {
			log("verified")
		} else if form.Changed() {
			log("resent " + fmt.Sprint(*resent))
		}
	})
}

func demoAccountBackupCodes(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	copied := intState(root, "copied", -1)
	root.Children(func() {
		if account.BackupCodes(c, []string{"ABCD-1234"}, copied, account.BackupCodesOptions{}).Changed() {
			log("copied " + fmt.Sprint(*copied))
		}
	})
}

func demoAccountPasskey(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	added := boolState(root, "added", false)
	root.Children(func() {
		if account.Passkey(c, added, account.PasskeyOptions{}).Changed() {
			log("added " + fmt.Sprint(*added))
		}
	})
}

func demoAccountSessions(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	sessions := ui.Local(root, "sessions", func() []account.Session {
		return []account.Session{{ID: "s1", Device: "Laptop", Location: "Home", Last: "now", Current: true}}
	})
	root.Children(func() {
		if account.Sessions(c, sessions, account.SessionsOptions{}).Changed() {
			log("sessions " + fmt.Sprint(len(*sessions)))
		}
	})
}

func demoAccountDevicesList(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	devices := ui.Local(root, "devices", func() []account.Device {
		return []account.Device{{ID: "d1", Name: "Laptop", Last: "now"}}
	})
	root.Children(func() {
		if account.DevicesList(c, devices, account.DevicesListOptions{}).Changed() {
			log("devices " + fmt.Sprint(len(*devices)))
		}
	})
}

func demoAccountProfileSettings(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	profile := ui.Local(root, "profile", func() account.Profile {
		return account.Profile{Name: "Ana", Handle: "@ana", Bio: "Designer"}
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
	root := ui.Column(c)
	image := ui.Local(root, "image", func() ui.ImageSource { return photos()[0] })
	change := boolState(root, "change", false)
	root.Children(func() {
		if account.AvatarEditor(c, image, change, account.AvatarEditorOptions{}).Changed() {
			log(fmt.Sprintf("image %v change %v", *image != nil, *change))
		}
	})
}

func demoAccountChangeEmail(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	email := textState(root, "email", "ana@example.com")
	root.Children(func() {
		form, verify := account.ChangeEmail(c, email, account.ChangeEmailOptions{})
		if verify.Changed() {
			log("verify")
		} else if form.Changed() {
			log("emailRunes " + fmt.Sprint(len([]rune(*email))))
		}
	})
}

func demoAccountChangePassword(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	value := ui.Local(root, "change", func() account.PasswordChange {
		return account.PasswordChange{Current: "oldpassword", New: "newpassword"}
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
	steps := []account.SecurityStep{{Title: "Add passkey"}, {Title: "Email verified", Done: true}}
	account.SecurityScore(c, 80, steps, account.SecurityScoreOptions{})
}

func demoAccountAPIKeys(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	keys := ui.Local(root, "keys", func() []account.APIKey {
		return []account.APIKey{{ID: "k1", Name: "CI", Prefix: "ck_live_", Secret: "secret"}}
	})
	root.Children(func() {
		if account.APIKeys(c, keys, account.APIKeysOptions{}).Changed() {
			log("revoked " + fmt.Sprint((*keys)[0].Revoked))
		}
	})
}

func demoAccountPlanCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	change := boolState(root, "change", false)
	plan := account.Plan{Name: "Pro", Price: "$10/month", Features: []string{"1500 components"}}
	root.Children(func() {
		if account.PlanCard(c, plan, true, change, account.PlanCardOptions{}).Changed() {
			log("change " + fmt.Sprint(*change))
		}
	})
}

func demoAccountUsage(c *ui.Context, log func(string)) {
	account.Usage(c, []account.Resource{{Name: "Storage", Unit: "GB", Used: 8, Limit: 10}}, account.UsageOptions{})
}

func demoAccountInvoiceList(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	downloaded := textState(root, "downloaded", "")
	invoices := []account.Invoice{{ID: "i1", Number: "INV-1", Date: "Oct 1", Amount: "$10", URL: "https://example.test/invoice"}}
	root.Children(func() {
		if account.InvoiceList(c, invoices, downloaded, account.InvoiceListOptions{}).Changed() {
			log("downloaded " + *downloaded)
		}
	})
}
