package input

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConfirmPasswordOptions tune a ConfirmPassword.
type ConfirmPasswordOptions struct {
	// Password, Confirm and Error are caller copy.
	Password, Confirm, Error string
}

// ConfirmPassword compares two masked fields and draws a match mark.
func ConfirmPassword(c *ui.Context, password, confirm *string, o ConfirmPasswordOptions) ui.Element {
	if password == nil || confirm == nil {
		panic("input: ConfirmPassword needs both values")
	}
	if o.Password == "" {
		o.Password = "Password"
	}
	if o.Confirm == "" {
		o.Confirm = "Confirm password"
	}
	if o.Error == "" {
		o.Error = "Passwords do not match."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	wasMatch := ui.Local(root, "match", func() bool { return false })
	return root.Children(func() {
		first := ui.TextInput(c, password).FillWidth().Label(o.Password).Placeholder(o.Password).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		first.Password()
		changedFirst := changed(root, first)
		second := ui.TextInput(c, confirm).FillWidth().Label(o.Confirm).Placeholder(o.Confirm).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		second.Password()
		changedSecond := changed(root, second)
		changedNow := changedFirst || changedSecond
		match := *password != "" && *password == *confirm
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			if match {
				kit.Arrive(ui.Icon(c, icon.Must("check-circle")).Size(16, 16).TextColor(t.Success), ui.Pose{Scale: 0.5, Fade: 1}, 0)
			} else if *confirm != "" {
				ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85))
			}
		})
		if match != *wasMatch {
			*wasMatch = match
			if match && changedNow {
				sound.Play(sound.Success)
			}
		}
	})
}
