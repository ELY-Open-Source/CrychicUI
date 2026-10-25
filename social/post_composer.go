package social

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PostComposerOptions tune the composer labels and limit.
type PostComposerOptions struct {
	Placeholder, Media, Remove, Send, Count string
	Limit                                   int
}

// PostComposer edits text with count and media; cue Send.
func PostComposer(c *ui.Context, draft *string, media ui.ImageSource, mediaAdded *bool, o PostComposerOptions) (root, send ui.Element) {
	if draft == nil || mediaAdded == nil || o.Limit < 0 {
		panic("social: PostComposer needs draft, media state, and limit")
	}
	def(&o.Placeholder, "What's happening?")
	def(&o.Media, "Add media")
	def(&o.Remove, "Remove media")
	def(&o.Send, "Post")
	def(&o.Count, "%d/%d")
	if o.Limit == 0 {
		o.Limit = 280
	}
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background)
	root.Border(1, root.AnimateColor("ring", t.Border.Mix(t.Accent, unit(root.FocusWithin())), ui.Snappy))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		in := ui.TextAreaBase(c, draft).Label(o.Placeholder).Placeholder(o.Placeholder).FillWidth().Lines(3, 6).FontSize(t.Rem(0.95))
		if in.Changed() {
			root.SetChanged()
		}
		c.Key("preview")
		ui.Column(c).FillWidth().Children(func() {
			if !*mediaAdded || media == nil {
				return
			}
			frame := ui.Box(c).FillWidth().AspectRatio(16.0 / 9).Radius(t.Radius).Clip().Background(t.SurfaceHover)
			frame.Transition(ui.ElementTransition{Exit: &ui.Motion{Scale: 0.96}})
			kit.Arrive(frame, ui.Pose{Scale: 0.96, Fade: 1}, 0).Children(func() {
				ui.Image(c, media).Fill().Fit(ui.Cover)
				x := kit.Press(ui.ButtonBase(c).Label(o.Remove).Absolute().Top(t.Space(2)).Right(t.Space(2)).Size(t.Space(7), t.Space(7)).
					Radius(999).TextColor(t.Background))
				x.Background(x.AnimateColor("bg", t.Text.Alpha(0.55).Mix(t.Text.Alpha(0.75), unit(x.Hovered())), ui.Snappy))
				x.Children(func() { glyph(c, "x", 14) })
				if x.Clicked() {
					*mediaAdded = false
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("PostComposer", "media", "added", false)
				}
			})
		})
		count := utf8.RuneCountInString(*draft)
		share := float32(count) / float32(o.Limit)
		ink := t.Accent
		switch {
		case count > o.Limit:
			ink = t.Danger
		case share >= 0.9:
			ink = t.Warning
		}
		c.Key("controls")
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).
			Padding(t.Space(2.5), 0, 0, 0).Children(func() {
			c.Key("media")
			add := kit.Press(ui.ButtonBase(c).Label(o.Media).Size(t.Space(8), t.Space(8)).Radius(999).TextColor(t.Accent).Disabled(media == nil || *mediaAdded))
			add.Background(add.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.12), unit(add.Hovered())), ui.Snappy))
			add.Children(func() { glyph(c, "image-plus", 18) })
			if add.Clicked() {
				*mediaAdded = true
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("PostComposer", "media", "added", true)
			}
			ui.Box(c).Grow(1)
			c.Key("count")
			meter := ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0)
			fill := meter.AnimateTo("fill", min(share, 1), ui.Smooth)
			stroke := meter.AnimateColor("ink", ink, ui.Snappy)
			meter.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(0.5)
				p.StrokePath(ring(cx, cy, rad), t.Space(0.5), t.Border)
				if fill > 0 {
					p.StrokePath(ring(cx, cy, rad).Trim(0, fill), t.Space(0.5), stroke)
				}
			}))
			number := kit.Numeric(c, fmt.Sprintf(o.Count, count, o.Limit)).FontSize(t.Rem(0.75)).Shrink(0)
			number.TextColor(number.AnimateColor("ink", t.TextMuted.Mix(ink, unit(ink != t.Accent)), ui.Snappy))
			valid := strings.TrimSpace(*draft) != "" && count <= o.Limit
			c.Key("send")
			send = kit.Press(ui.ButtonBase(c).Label(o.Send).Padding(t.Space(1.5), t.Space(4)).Radius(999).Shrink(0).TextColor(t.AccentText).Disabled(!valid))
			send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy))
			send.Children(func() { ui.Text(c, o.Send).FontSize(t.Rem(0.85)).Bold().Shrink(0) })
		})
	})
	if send.Clicked() {
		send.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("PostComposer", "sent", "runes", utf8.RuneCountInString(*draft))
	}
	return root, send
}
