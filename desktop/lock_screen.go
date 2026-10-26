package desktop

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LockScreenOptions dress the screen; Clock and Date are layouts.
type LockScreenOptions struct {
	Photo                                *ui.Bitmap
	Clock, Date, Password, Unlock, Wrong string
	Welcome                              string
}

// LockScreen unlocks when check passes; cue Success.
func LockScreen(c *ui.Context, user string, password *string, check func(string) bool, unlocked *bool, o LockScreenOptions) ui.Element {
	def(&o.Clock, "15:04")
	def(&o.Date, "Monday, January 2")
	def(&o.Password, "Password")
	def(&o.Unlock, "Unlock")
	def(&o.Wrong, "Wrong password. Try again.")
	def(&o.Welcome, "Welcome back")
	if strings.TrimSpace(user) == "" || password == nil || check == nil || unlocked == nil {
		panic(fmt.Sprintf("desktop: LockScreen %q needs a password, a check and an unlocked flag", user))
	}
	t := c.Theme()
	now := minute(c)
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(8), t.Space(5)).Radius(t.Radius * 2).Clip().
		Background(t.Accent.Mix(t.Background, 0.82)).Role(ui.RoleGroup).Label(user)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	wrong := ui.Local(root, "wrong", func() bool { return false })
	shook := ui.Local(root, "shook", func() time.Time { return time.Time{} })
	kick := c.Now().Sub(*shook) < 60*time.Millisecond
	if kick {
		c.After(60*time.Millisecond - c.Now().Sub(*shook))
	}
	return root.Children(func() {
		kit.Numeric(c, now.Format(o.Clock)).FontSize(t.Rem(3.2)).Bold().Shrink(0).MaxWidth(t.Space(80)).SingleLine()
		ui.Text(c, now.Format(o.Date)).FontSize(t.Rem(0.9)).TextColor(t.Text.Alpha(0.7)).Shrink(0).MaxWidth(t.Space(80)).SingleLine()
		ui.Box(c).Height(t.Space(6))
		avatar(c, user, o.Photo, 16)
		ui.Text(c, user).FontSize(t.Rem(0.95)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
		if *unlocked {
			c.Key("open")
			kit.Arrive(ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center), ui.Pose{Scale: 0.9, Fade: 1}, 0).Children(func() {
				glyph(c, "unlock", 16).TextColor(t.Success)
				ui.Text(c, o.Welcome).FontSize(t.Rem(0.84)).Bold().TextColor(t.Success).Shrink(0).MaxWidth(t.Space(60)).SingleLine()
			})
			return
		}
		var in, enter ui.Element
		submit := func() {
			if check(*password) {
				*unlocked, *wrong = true, false
				sound.Play(sound.Success)
				kit.Log("LockScreen", "unlocked")
			} else {
				*wrong, *shook = true, c.Now()
				sound.Play(sound.Error)
				kit.Log("LockScreen", "refused", "runes", utf8.RuneCountInString(*password))
			}
			*password = ""
			root.SetChanged()
		}
		c.Key("field")
		field := ui.Row(c).Width(t.Space(64)).MaxWidth(t.Space(64)).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).
			Radius(999).Background(t.Background.Alpha(0.85))
		field.Border(1, field.AnimateColor("ring", t.Border.Mix(t.Accent, unit(field.FocusWithin())).Mix(t.Danger, unit(*wrong)), ui.Snappy))
		field.Translate(field.AnimateTo("shake", t.Space(3)*unit(kick), ui.Bouncy), 0).Children(func() {
			in = ui.TextInputBase(c, password).Label(o.Password).Placeholder(o.Password).Password().FontSize(t.Rem(0.86)).Grow(1).MinWidth(0)
			if in.Changed() {
				*wrong = false
				root.SetChanged()
				kit.Log("LockScreen", "typed", "runes", utf8.RuneCountInString(*password))
			}
			enter = kit.Press(ui.ButtonBase(c).Label(o.Unlock).Size(t.Space(7), t.Space(7)).Center().Radius(999).Disabled(*password == ""))
			enter.Background(enter.AnimateColor("bg", t.Accent.Mix(t.Text, 0.12*unit(enter.Hovered())), ui.Snappy))
			enter.Children(func() { glyph(c, "arrow-right", 14).TextColor(t.AccentText) })
		})
		if (enter.Clicked() || in.Submitted()) && *password != "" {
			submit()
		}
		c.Key("note")
		ui.Text(c, o.Wrong).FontSize(t.Rem(0.76)).TextColor(t.Danger).Opacity(root.AnimateTo("wrong", unit(*wrong), ui.Snappy)).
			MaxWidth(t.Space(64)).TextAlign(ui.Center).Wrap()
	})
}
