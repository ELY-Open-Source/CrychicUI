package agent

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AnswerActionsOptions name the actions.
type AnswerActionsOptions struct {
	Copy, Retry, Like, Dislike string
}

// AnswerActions are copy, retry, like and dislike.
func AnswerActions(c *ui.Context, text string, vote *int, o AnswerActionsOptions) (root, copy, retry ui.Element) {
	if text == "" || vote == nil || *vote < -1 || *vote > 1 {
		panic("agent: AnswerActions needs text and a vote from -1 to 1")
	}
	def(&o.Copy, "Copy")
	def(&o.Retry, "Retry")
	def(&o.Like, "Like")
	def(&o.Dislike, "Dislike")
	t := c.Theme()
	root = ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(0.5)).Role(ui.RoleGroup)
	var like, dislike ui.Element
	root.Children(func() {
		copy = button.CopyButton(c, o.Copy, text)
		retry = round(c, o.Retry, "rotate-cw")
		like = round(c, o.Like, "thumbs-up").Checked(*vote == 1)
		dislike = round(c, o.Dislike, "thumbs-down").Checked(*vote == -1)
	})
	if retry.Clicked() {
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("AnswerActions", "retry")
	}
	pick := func(v int, name string) {
		if *vote == v {
			*vote = 0
		} else {
			*vote = v
		}
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("AnswerActions", name, "vote", *vote)
	}
	if like.Clicked() {
		pick(1, "like")
	}
	if dislike.Clicked() {
		pick(-1, "dislike")
	}
	up, down := t.TextMuted, t.TextMuted
	if *vote == 1 {
		up = t.Accent
	}
	if *vote == -1 {
		down = t.Accent
	}
	like.TextColor(like.AnimateColor("ink", up, ui.Snappy))
	dislike.TextColor(dislike.AnimateColor("ink", down, ui.Snappy))
	return root, copy, retry
}
