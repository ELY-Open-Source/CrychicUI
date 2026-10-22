package form

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Bug is what a BugReport edits; Severity is a choice.
type Bug struct{ Title, Steps, Severity, Attachment string }

// BugReportOptions tune a BugReport.
type BugReportOptions struct {
	// Labels and button copy; Attach tunes the file field.
	Title, Steps, Severity, Attachment, Send string
	Attach                                   FileFieldOptions
}

// BugReport takes title, steps, severity and an attached file.
func BugReport(c *ui.Context, value *Bug, severities []string, o BugReportOptions) (root, send, attach ui.Element) {
	if value == nil || len(severities) == 0 {
		panic("form: BugReport needs a value and severities")
	}
	def(&o.Title, "Title")
	def(&o.Steps, "Steps to reproduce")
	def(&o.Severity, "Severity")
	def(&o.Attachment, "Attachment")
	def(&o.Send, "Send report")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(4))
	root.Children(func() {
		FieldGroup(c, o.Title, FieldGroupOptions{}, func() { edited(root, ui.TextInput(c, &value.Title).FillWidth()) })
		FieldGroup(c, o.Steps, FieldGroupOptions{}, func() {
			edited(root, ui.TextArea(c, &value.Steps).FillWidth().Lines(4, 8))
		})
		FieldGroup(c, o.Severity, FieldGroupOptions{}, func() {
			edited(root, choice.Dropdown(c, o.Severity, severities, &value.Severity))
		})
		var file ui.Element
		file, attach = FileField(c, o.Attachment, &value.Attachment, o.Attach)
		edited(root, file)
		send = primary(c, o.Send).Disabled(blank(value.Title) || blank(value.Steps))
	})
	if kit.Clicked(send, sound.Send) {
		kit.Log("BugReport", "send", "steps", len(value.Steps), "attached", value.Attachment != "")
	}
	return root, send, attach
}
