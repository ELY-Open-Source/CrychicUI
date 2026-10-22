package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/form"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/egoist/mygo/ui"
)

var formEntries = []Entry{
	{451, "FormCard", "form", "Fields on a card with title and actions.", demoFormCard},
	{452, "FormSection", "form", "A titled group of fields with a description.", demoFormSection},
	{453, "FormRow", "form", "Fields side by side, stacking when narrow.", demoFormRow},
	{454, "FieldGroup", "form", "Label, control, help and a shaking error.", demoFieldGroup},
	{455, "FieldBox", "form", "A fieldset with a legend, disabled as a group.", demoFieldBox},
	{456, "InlineForm", "form", "An input and a submit in one row.", demoInlineForm},
	{457, "Wizard", "form", "Steps with Back and Next; pages slide.", demoWizard},
	{458, "StepHeader", "form", "Numbered steps with a filling connector.", demoStepHeader},
	{459, "Survey", "form", "One question at a time.", demoSurvey},
	{460, "ContactForm", "form", "Name, email and message, validated.", demoContactForm},
	{461, "LoginBox", "form", "Email, password and remember me.", demoLoginBox},
	{462, "FeedbackForm", "form", "Stars and a comment.", demoFeedbackForm},
	{463, "NewsletterForm", "form", "Email and subscribe; thanks swaps in.", demoNewsletterForm},
	{464, "AdvancedSearch", "form", "Field, operator and value rules.", demoAdvancedSearch},
	{465, "FilterPanel", "form", "Faceted checks with counts.", demoFilterPanel},
	{466, "SettingsForm", "form", "Grouped switch and select rows.", demoSettingsForm},
	{467, "ProfileForm", "form", "Avatar, name, bio and links.", demoProfileForm},
	{468, "AddressForm", "form", "An address with a country select.", demoAddressForm},
	{469, "PaymentForm", "form", "Card, expiry, CVC and name.", demoPaymentForm},
	{470, "ShippingForm", "form", "Address and delivery speed cards.", demoShippingForm},
}

func demoFormCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	name := ui.Local(root, "name", func() string { return "" })
	root.Children(func() {
		_, save, cancel := form.FormCard(c, "Workspace", form.FormCardOptions{Description: "Shown to your team."}, func() {
			form.FieldGroup(c, "Name", form.FieldGroupOptions{}, func() {
				if ui.TextInput(c, name).Changed() {
					log(fmt.Sprintf("name: %d chars", len(*name)))
				}
			})
		})
		if save.Clicked() {
			log(fmt.Sprintf("save: %d chars", len(*name)))
		}
		if cancel.Clicked() {
			log("cancel")
		}
	})
}

func demoFormSection(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	city := ui.Local(root, "city", func() string { return "" })
	root.Children(func() {
		form.FormSection(c, "Billing", form.FormSectionOptions{Description: "Where we send invoices."}, func() {
			form.FieldGroup(c, "City", form.FieldGroupOptions{}, func() {
				if ui.TextInput(c, city).Changed() {
					log(fmt.Sprintf("city: %d chars", len(*city)))
				}
			})
		})
	})
}

func demoFormRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).FillWidth()
	first := ui.Local(root, "first", func() string { return "" })
	last := ui.Local(root, "last", func() string { return "" })
	field := func(label string, v *string) func() {
		return func() {
			form.FieldGroup(c, label, form.FieldGroupOptions{}, func() {
				if ui.TextInput(c, v).Changed() {
					log(fmt.Sprintf("%s: %d chars", label, len(*v)))
				}
			})
		}
	}
	root.Children(func() { form.FormRow(c, []func(){field("First name", first), field("Last name", last)}) })
}

func demoFieldGroup(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	handle := ui.Local(root, "handle", func() string { return "ada" })
	root.Children(func() {
		problem := ""
		if len(*handle) < 3 {
			problem = "Use three letters or more."
		}
		form.FieldGroup(c, "Handle", form.FieldGroupOptions{Help: "Your public name.", Error: problem}, func() {
			if ui.TextInput(c, handle).Changed() {
				log(fmt.Sprintf("handle: %d chars", len(*handle)))
			}
		})
	})
}

func demoFieldBox(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(360)
	off := boolState(root, "off", false)
	mail := boolState(root, "mail", true)
	root.Children(func() {
		if choice.Toggle(c, "Lock notifications", off).Changed() {
			log(fmt.Sprintf("locked: %v", *off))
		}
		form.FieldBox(c, "Notifications", form.FieldBoxOptions{Disabled: *off}, func() {
			if ui.Checkbox(c, mail, "Weekly summary").Changed() {
				log(fmt.Sprintf("summary: %v", *mail))
			}
		})
	})
}

func demoInlineForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	email := ui.Local(root, "email", func() string { return "" })
	root.Children(func() {
		f, send := form.InlineForm(c, "Invite by email", email, form.InlineFormOptions{Placeholder: "name@team.com",
			Submit: "Invite"})
		if f.Changed() {
			log(fmt.Sprintf("email: %d chars", len(*email)))
		}
		if send.Clicked() {
			log(fmt.Sprintf("invite: %d chars", len(*email)))
		}
	})
}

func demoWizard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	at := intState(root, "at", 0)
	root.Children(func() {
		w, finish := form.Wizard(c, []string{"Account", "Team", "Done"}, at, func(i int) {
			ui.Text(c, []string{"Create your account.", "Invite your team.", "All set."}[i])
		}, form.WizardOptions{})
		if w.Changed() {
			log(fmt.Sprintf("step %d", *at+1))
		}
		if finish.Clicked() {
			log("finish")
		}
	})
}

func demoStepHeader(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(440)
	at := intState(root, "at", 1)
	root.Children(func() {
		form.StepHeader(c, []string{"Cart", "Address", "Pay"}, *at)
		if ui.Button(c, "Next step").AlignSelf(ui.Start).Clicked() {
			*at = (*at + 1) % 3
			log(fmt.Sprintf("step %d", *at+1))
		}
	})
}

func demoSurvey(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	answers := ui.Local(root, "answers", func() []int { return []int{-1, -1} })
	qs := []form.Question{{Text: "How often do you ship?", Options: []string{"Daily", "Weekly", "Monthly"}},
		{Text: "Team size?", Options: []string{"Just me", "2 to 10", "More"}}}
	root.Children(func() {
		s, back := form.Survey(c, qs, answers, form.SurveyOptions{})
		if s.Changed() {
			log(fmt.Sprintf("answers: %v", *answers))
		}
		if back.Clicked() {
			log("back")
		}
	})
}

func demoContactForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	v := ui.Local(root, "v", func() form.Contact { return form.Contact{} })
	root.Children(func() {
		f, send := form.ContactForm(c, v, form.ContactFormOptions{})
		if f.Changed() {
			log(fmt.Sprintf("edit: %d, %d, %d chars", len(v.Name), len(v.Email), len(v.Message)))
		}
		if send.Clicked() {
			log(fmt.Sprintf("send: %d chars", len(v.Message)))
		}
	})
}

func demoLoginBox(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	v := ui.Local(root, "v", func() form.Login { return form.Login{} })
	failed := boolState(root, "failed", false)
	root.Children(func() {
		problem := ""
		if *failed {
			problem = "That password is wrong."
		}
		f, submit := form.LoginBox(c, v, form.LoginBoxOptions{Error: problem})
		if f.Changed() {
			log(fmt.Sprintf("edit: email %d chars, remember %v", len(v.Email), v.Remember))
		}
		if submit.Clicked() {
			*failed = !*failed
			log(fmt.Sprintf("sign in, remember %v", v.Remember))
		}
	})
}

func demoFeedbackForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	v := ui.Local(root, "v", func() form.Feedback { return form.Feedback{} })
	root.Children(func() {
		f, send := form.FeedbackForm(c, v, form.FeedbackFormOptions{})
		if f.Changed() {
			log(fmt.Sprintf("rating %d", v.Rating))
		}
		if send.Clicked() {
			log(fmt.Sprintf("send: %d stars", v.Rating))
		}
	})
}

func demoNewsletterForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	email := ui.Local(root, "email", func() string { return "" })
	root.Children(func() {
		f, sub := form.NewsletterForm(c, email, form.NewsletterFormOptions{})
		if f.Changed() {
			log(fmt.Sprintf("email: %d chars", len(*email)))
		}
		if sub.Clicked() {
			log("subscribed")
		}
	})
}

func demoAdvancedSearch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(560)
	rules := ui.Local(root, "rules", func() []form.Rule { return []form.Rule{{Field: "Title", Op: "contains"}} })
	root.Children(func() {
		s, search := form.AdvancedSearch(c, []string{"Title", "Author", "Tag"}, rules, form.AdvancedSearchOptions{})
		if s.Changed() {
			log(fmt.Sprintf("rules: %d", len(*rules)))
		}
		if search.Clicked() {
			log(fmt.Sprintf("search with %d rules", len(*rules)))
		}
	})
}

func demoFilterPanel(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(300)
	picked := ui.Local(root, "picked", func() []form.Filter { return nil })
	facets := []form.Facet{{Title: "Color", Options: []string{"Red", "Blue", "Green"}, Counts: []int{12, 8, 3}},
		{Title: "Size", Options: []string{"S", "M", "L"}, Counts: []int{4, 9, 6}}}
	root.Children(func() {
		if form.FilterPanel(c, facets, picked, form.FilterPanelOptions{}).Changed() {
			log(fmt.Sprintf("filters: %d", len(*picked)))
		}
	})
}

func demoSettingsForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	sounds := boolState(root, "sounds", true)
	mail := boolState(root, "mail", false)
	theme := ui.Local(root, "theme", func() string { return "System" })
	root.Children(func() {
		if form.SettingsForm(c, []form.SettingGroup{
			{Title: "General", Settings: []form.Setting{{Title: "Sounds", Description: "Play cues on actions.", On: sounds},
				{Title: "Theme", Choice: theme, Choices: []string{"System", "Light", "Dark"}}}},
			{Title: "Email", Settings: []form.Setting{{Title: "Digest", Description: "A weekly summary.", On: mail}}},
		}).Changed() {
			log(fmt.Sprintf("sounds %v, theme %s, digest %v", *sounds, *theme, *mail))
		}
	})
}

func demoProfileForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	v := ui.Local(root, "v", func() form.Profile { return form.Profile{Name: "Ada Lovelace"} })
	root.Children(func() {
		f, save := form.ProfileForm(c, v, form.ProfileFormOptions{})
		if f.Changed() {
			log(fmt.Sprintf("edit: name %d chars, %d links", len(v.Name), len(v.Links)))
		}
		if save.Clicked() {
			log(fmt.Sprintf("saved with %d links", len(v.Links)))
		}
	})
}

func demoAddressForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	addr := ui.Local(root, "addr", func() input.Address { return input.Address{} })
	country := ui.Local(root, "country", func() string { return "Canada" })
	root.Children(func() {
		if form.AddressForm(c, addr, country, []string{"Canada", "Japan", "Germany"}, form.AddressFormOptions{}).Changed() {
			log(fmt.Sprintf("%s, street %d chars", *country, len(addr.Street)))
		}
	})
}

func demoPaymentForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	v := ui.Local(root, "v", func() form.Payment { return form.Payment{} })
	root.Children(func() {
		f, pay := form.PaymentForm(c, v, form.PaymentFormOptions{})
		if f.Changed() {
			log(fmt.Sprintf("edit: card number %d chars", len(v.Number)))
		}
		if pay.Clicked() {
			log("pay")
		}
	})
}

func demoShippingForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	v := ui.Local(root, "v", func() form.Shipping { return form.Shipping{Speed: "Standard"} })
	root.Children(func() {
		if form.ShippingForm(c, v, []form.Speed{{Name: "Standard", Detail: "Free, 5 to 7 days"},
			{Name: "Express", Detail: "$12, next day"}}, form.ShippingFormOptions{}).Changed() {
			log(fmt.Sprintf("speed %s", v.Speed))
		}
	})
}
