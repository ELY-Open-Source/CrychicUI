package code

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ErrorOverlayOptions name the close button.
type ErrorOverlayOptions struct {
	Close string
}

var overlayFade = ui.ElementTransition{Duration: 180 * time.Millisecond, Exit: &ui.Motion{}}

// ErrorOverlay slides a compile error panel down over dimmed content.
func ErrorOverlay(c *ui.Context, title, detail string, open *bool, o ErrorOverlayOptions, content func()) ui.Element {
	if title == "" || open == nil || content == nil {
		panic(fmt.Sprintf("code: ErrorOverlay %q needs a title, open and content", title))
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c).FillWidth().Clip().Radius(t.Radius)
	height := ui.Local(root, "height", func() float32 { return 0 })
	if !*open {
		*height = 0
	}
	root.MinHeight(*height + t.Space(4))
	var shut ui.Element
	root.Children(func() {
		// Content gets its own scope beside the overlay's keys.
		ui.Box(c).Key("content").FillWidth().Children(content)
		if !*open {
			return
		}
		scrim := ui.Box(c).Key("scrim").Absolute().Top(0).Left(0).Right(0).Bottom(0).Background(t.Background.Alpha(0.7)).Transition(overlayFade)
		kit.Arrive(scrim, ui.Pose{Fade: 1}, 0)
		panel := ui.Column(c).Key("panel").Absolute().Top(t.Space(2)).Left(t.Space(2)).Right(t.Space(2)).
			Background(t.Background).Border(1, t.Danger.Alpha(0.5)).Radius(t.Radius).Clip().Shadow(0, 8, 24, 0, kit.Shade(0.16)).
			Role(ui.RoleDialog).Label(title).Transition(ui.ElementTransition{Duration: 180 * time.Millisecond, Exit: &ui.Motion{Y: -t.Space(6)}})
		kit.Arrive(panel, ui.Pose{Y: -t.Space(6), Fade: 1}, 0)
		panel.Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(3)).
				Background(t.Danger.Alpha(0.08)).Children(func() {
				ui.Icon(c, icon.Must("circle-x")).Size(16, 16).Shrink(0).TextColor(t.Danger)
				ui.Text(c, title).Bold().TextColor(t.Danger).Grow(1).MinWidth(0).SingleLine()
				shut = round(c, o.Close, "x").Size(t.Space(7), t.Space(7))
			})
			if detail != "" {
				ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(3)).Children(func() {
					for i, line := range linesOf(detail) {
						c.Key(i)
						ink := t.TextMuted
						if i == 0 {
							ink = t.Text
						}
						ui.Text(c, untab(line)).Font("monospace").FontSize(t.Rem(0.8)).TextColor(ink).Wrap()
					}
				})
			}
		})
		// The root settles: the panel's exit copy outlives this pass.
		kit.Settle(c, root, []float32{*height}, func() []float32 { return []float32{panel.LayoutBounds().H} })
		if h := panel.LayoutBounds().H; h > 0 {
			*height = h
		}
	})
	if *open && shut.Clicked() {
		*open = false
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("ErrorOverlay", "closed")
	}
	return root
}
