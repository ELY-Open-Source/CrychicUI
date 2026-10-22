package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/form"
	"github.com/egoist/mygo/ui"
)

var formEntries2 = []Entry{
	{471, "BugReport", "form", "Title, steps, severity and a file.", demoBugReport},
	{472, "InviteForm", "form", "Emails as chips and a role.", demoInviteForm},
	{473, "TypeToConfirm", "form", "Type a name to enable a destructive action.", demoTypeToConfirm},
	{474, "DynamicFields", "form", "Add and remove fields of chosen kinds.", demoDynamicFields},
	{475, "ConditionalField", "form", "A field that springs open under a switch.", demoConditionalField},
	{476, "ActionBar", "form", "Save and Cancel slide up when dirty.", demoActionBar},
	{477, "DirtyBar", "form", "Unsaved changes, with Discard and Save.", demoDirtyBar},
	{478, "ErrorSummary", "form", "Errors that focus their fields.", demoErrorSummary},
	{479, "HelpPopover", "form", "A ? that opens field help.", demoHelpPopover},
	{480, "RequiredLegend", "form", "The * Required note.", demoRequiredLegend},
	{481, "AutosaveForm", "form", "Saves after a pause; status cross-fades.", demoAutosaveForm},
	{482, "ReviewStep", "form", "A read-only summary with Edit links.", demoReviewStep},
	{483, "FileField", "form", "A dropped file with remove and browse.", demoFileField},
	{484, "ImageField", "form", "An image preview with replace and remove.", demoImageField},
	{485, "Consent", "form", "A terms check that enables submit.", demoConsent},
	{486, "ScheduleField", "form", "A date and time pair.", demoScheduleField},
	{487, "NumberRange", "form", "Min and max kept in order.", demoNumberRange},
	{488, "Repeater", "form", "Groups of fields added and removed.", demoRepeater},
	{489, "QuestionBlock", "form", "Prompt, hint and a counted answer.", demoQuestionBlock},
	{490, "FormResult", "form", "A success or error panel replacing the form.", demoFormResult},
}

func demoBugReport(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	v := ui.Local(root, "v", func() form.Bug { return form.Bug{Severity: "Medium"} })
	root.Children(func() {
		f, send, attach := form.BugReport(c, v, []string{"Low", "Medium", "High"}, form.BugReportOptions{})
		if f.Changed() {
			log(fmt.Sprintf("edit: title %d chars, %s", len(v.Title), v.Severity))
		}
		if attach.Clicked() {
			v.Attachment = "/tmp/crash.log"
			log("attach")
		}
		if send.Clicked() {
			log("send")
		}
	})
}

func demoInviteForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	emails := listState(root, "emails")
	role := ui.Local(root, "role", func() string { return "Member" })
	root.Children(func() {
		f, send := form.InviteForm(c, emails, role, []string{"Member", "Admin", "Viewer"}, form.InviteFormOptions{})
		if f.Changed() {
			log(fmt.Sprintf("%d emails, %s", len(*emails), *role))
		}
		if send.Clicked() {
			log(fmt.Sprintf("invite %d", len(*emails)))
		}
	})
}

func demoTypeToConfirm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	typed := ui.Local(root, "typed", func() string { return "" })
	root.Children(func() {
		f, ok := form.TypeToConfirm(c, "acme-prod", typed, form.TypeToConfirmOptions{})
		if f.Changed() {
			log(fmt.Sprintf("typed %d chars", len(*typed)))
		}
		if ok.Clicked() {
			log("deleted")
		}
	})
}

func demoDynamicFields(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	fields := ui.Local(root, "fields", func() []form.DynamicField {
		return []form.DynamicField{{Kind: "Phone"}}
	})
	root.Children(func() {
		if form.DynamicFields(c, fields, []string{"Phone", "Website", "Note"}, form.DynamicFieldsOptions{}).Changed() {
			log(fmt.Sprintf("fields: %d", len(*fields)))
		}
	})
}

func demoConditionalField(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	on := boolState(root, "on", false)
	city := ui.Local(root, "city", func() string { return "" })
	root.Children(func() {
		if form.ConditionalField(c, "Ship to another address", on, func() {
			form.FieldGroup(c, "City", form.FieldGroupOptions{}, func() {
				if ui.TextInput(c, city).Changed() {
					log(fmt.Sprintf("city: %d chars", len(*city)))
				}
			})
		}).Changed() {
			log(fmt.Sprintf("other address: %v", *on))
		}
	})
}

func demoActionBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Size(420, 260)
	name := ui.Local(root, "name", func() string { return "Ada" })
	saved := ui.Local(root, "saved", func() string { return "Ada" })
	root.Children(func() {
		_, save, cancel := form.ActionBar(c, *name != *saved, form.ActionBarOptions{}, func() {
			ui.Column(c).Padding(12).Children(func() {
				form.FieldGroup(c, "Display name", form.FieldGroupOptions{}, func() {
					if ui.TextInput(c, name).Changed() {
						log(fmt.Sprintf("name: %d chars", len(*name)))
					}
				})
			})
		})
		if save.Clicked() {
			*saved = *name
			log("save")
		}
		if cancel.Clicked() {
			*name = *saved
			log("cancel")
		}
	})
}

func demoDirtyBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(420)
	dirty := boolState(root, "dirty", true)
	root.Children(func() {
		if ui.Button(c, "Make a change").AlignSelf(ui.Start).Clicked() {
			*dirty = true
			log("changed")
		}
		_, discard, save := form.DirtyBar(c, *dirty, form.DirtyBarOptions{})
		if discard.Clicked() {
			*dirty = false
			log("discard")
		}
		if save.Clicked() {
			*dirty = false
			log("save")
		}
	})
}

func demoErrorSummary(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(10).Width(400)
	name := ui.Local(root, "name", func() string { return "" })
	email := ui.Local(root, "email", func() string { return "ada@" })
	nameField := ui.Local(root, "nameField", func() *ui.Handle { return new(ui.Handle) })
	emailField := ui.Local(root, "emailField", func() *ui.Handle { return new(ui.Handle) })
	root.Children(func() {
		var errs []feedback.FormError
		if *name == "" {
			errs = append(errs, feedback.FormError{Message: "Enter your name.", Field: *nameField})
		}
		if !emailOK(*email) {
			errs = append(errs, feedback.FormError{Message: "Finish the email address.", Field: *emailField})
		}
		form.ErrorSummary(c, errs, form.ErrorSummaryOptions{})
		if ui.TextInput(c.Key("name"), name).Label("Name").Bind(*nameField).Changed() {
			log(fmt.Sprintf("name: %d chars, %d errors", len(*name), len(errs)))
		}
		if ui.TextInput(c.Key("email"), email).Label("Email").Bind(*emailField).Changed() {
			log(fmt.Sprintf("email: %d chars, %d errors", len(*email), len(errs)))
		}
	})
}

// emailOK is the demo's own rough check.
func emailOK(s string) bool {
	for i := 1; i < len(s)-1; i++ {
		if s[i] == '@' {
			return true
		}
	}
	return false
}

func demoHelpPopover(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
		ui.Text(c, "Slug")
		if form.HelpPopover(c, "About the slug", "The slug is the last part of the page's link.").Clicked() {
			log("help")
		}
	})
}

func demoRequiredLegend(c *ui.Context, log func(string)) {
	form.RequiredLegend(c, form.RequiredLegendOptions{})
}

func demoAutosaveForm(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	note := ui.Local(root, "note", func() string { return "" })
	edits := ui.Local(root, "edits", func() uint64 { return 0 })
	root.Children(func() {
		if form.AutosaveForm(c, edits, form.AutosaveFormOptions{}, func() {
			form.FieldGroup(c, "Note", form.FieldGroupOptions{}, func() {
				if ui.TextArea(c, note).Lines(3, 6).Changed() {
					*edits++
				}
			})
		}).Changed() {
			log(fmt.Sprintf("saved %d chars", len(*note)))
		}
	})
}

func demoReviewStep(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	edit := intState(root, "edit", -1)
	root.Children(func() {
		if form.ReviewStep(c, []form.ReviewSection{
			{Title: "Contact", Items: []display.KeyValuePair{{Key: "Name", Value: "Ada Lovelace"}, {Key: "Email", Value: "ada@example.com"}}},
			{Title: "Shipping", Items: []display.KeyValuePair{{Key: "City", Value: "London"}, {Key: "Speed", Value: "Express"}}},
		}, edit, form.ReviewStepOptions{}).Changed() {
			log(fmt.Sprintf("edit section %d", *edit+1))
		}
	})
}

func demoFileField(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	path := ui.Local(root, "path", func() string { return "" })
	root.Children(func() {
		f, browse := form.FileField(c, "Attachment", path, form.FileFieldOptions{})
		if f.Changed() {
			log(fmt.Sprintf("file set: %v", *path != ""))
		}
		if browse.Clicked() {
			*path = "/Users/ada/report.pdf"
			log("browse")
		}
	})
}

func demoImageField(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	path := ui.Local(root, "path", func() string { return "photo-1.png" })
	n := intState(root, "n", 0)
	root.Children(func() {
		var preview *ui.Bitmap
		if *path != "" {
			preview = photos()[*n%len(photos())]
		}
		f, pick := form.ImageField(c, "Cover photo", path, preview, form.ImageFieldOptions{})
		if f.Changed() {
			log(fmt.Sprintf("image set: %v", *path != ""))
		}
		if pick.Clicked() {
			*n++
			*path = fmt.Sprintf("photo-%d.png", *n+1)
			log("pick")
		}
	})
}

func demoConsent(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	agreed := boolState(root, "agreed", false)
	root.Children(func() {
		f, submit := form.Consent(c, "I accept the terms of service", agreed, form.ConsentOptions{})
		if f.Changed() {
			log(fmt.Sprintf("agreed: %v", *agreed))
		}
		if submit.Clicked() {
			log("continue")
		}
	})
}

func demoScheduleField(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	at := ui.Local(root, "at", func() time.Time { return time.Date(2026, 11, 2, 9, 30, 0, 0, time.Local) })
	root.Children(func() {
		if form.ScheduleField(c, "Starts", at, form.ScheduleFieldOptions{}).Changed() {
			log(at.Format("Jan 2 15:04"))
		}
	})
}

func demoNumberRange(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	low := ui.Local(root, "low", func() float64 { return 20 })
	high := ui.Local(root, "high", func() float64 { return 80 })
	root.Children(func() {
		if form.NumberRange(c, "Price", low, high, 0, 500, 5, form.NumberRangeOptions{}).Changed() {
			log(fmt.Sprintf("%g to %g", *low, *high))
		}
	})
}

// guest is one group the Repeater demo edits.
type guest struct{ Name, Email string }

func demoRepeater(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	guests := ui.Local(root, "guests", func() []guest { return []guest{{}} })
	root.Children(func() {
		if form.Repeater(c, guests, func(i int, g *guest) {
			form.FieldGroup(c, fmt.Sprintf("Guest %d name", i+1), form.FieldGroupOptions{}, func() {
				if ui.TextInput(c, &g.Name).Changed() {
					log(fmt.Sprintf("guest %d name: %d chars", i+1, len(g.Name)))
				}
			})
		}, form.RepeaterOptions{}).Changed() {
			log(fmt.Sprintf("guests: %d", len(*guests)))
		}
	})
}

func demoQuestionBlock(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	answer := ui.Local(root, "answer", func() string { return "" })
	root.Children(func() {
		if form.QuestionBlock(c, "What should we build next?", answer,
			form.QuestionBlockOptions{Hint: "A sentence or two is plenty."}).Changed() {
			log(fmt.Sprintf("answer: %d chars", len(*answer)))
		}
	})
}

func demoFormResult(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	outcome := ui.Local(root, "outcome", func() form.Outcome { return form.Pending })
	tries := intState(root, "tries", 0)
	root.Children(func() {
		_, action := form.FormResult(c, *outcome, form.FormResultOptions{}, func() {
			if ui.Button(c, "Send").AlignSelf(ui.Start).Clicked() {
				*tries++
				*outcome = form.Failed
				if *tries > 1 {
					*outcome = form.Succeeded
				}
				log(fmt.Sprintf("send, try %d", *tries))
			}
		})
		if action.Clicked() {
			*outcome = form.Pending
			log("back to the form")
		}
	})
}
