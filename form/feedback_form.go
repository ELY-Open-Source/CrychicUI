package form

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Feedback is what a FeedbackForm edits; Rating 0 is unrated.
type Feedback struct {
	Rating  int
	Comment string
}

// FeedbackFormOptions tune a FeedbackForm.
type FeedbackFormOptions struct {
	// Rating, Comment and Send are labels.
	Rating, Comment, Send string
}

// FeedbackForm takes stars and a comment; Send enables once rated.
func FeedbackForm(c *ui.Context, value *Feedback, o FeedbackFormOptions) (root, send ui.Element) {
	if value == nil {
		panic("form: FeedbackForm needs a value")
	}
	def(&o.Rating, "Rating")
	def(&o.Comment, "Comment")
	def(&o.Send, "Send feedback")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(4))
	root.Children(func() {
		FieldGroup(c, o.Rating, FieldGroupOptions{}, func() {
			edited(root, choice.StarRating(c, o.Rating, &value.Rating, choice.StarRatingOptions{}))
		})
		FieldGroup(c, o.Comment, FieldGroupOptions{}, func() {
			edited(root, ui.TextArea(c, &value.Comment).FillWidth().Lines(3, 6))
		})
		send = primary(c, o.Send).Disabled(value.Rating == 0)
	})
	if kit.Clicked(send, sound.Send) {
		kit.Log("FeedbackForm", "send", "rating", value.Rating, "comment", len(value.Comment))
	}
	return root, send
}
