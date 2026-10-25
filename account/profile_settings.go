package account

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Profile is one editable account profile.
type Profile struct {
	Name, Handle, Bio string
}

// ProfileSettingsOptions word the form; Count takes used and limit.
type ProfileSettingsOptions struct {
	Title, Name, Handle, Bio, Save, Saved, Count string
	Limit                                        int
}

// ProfileSettings saves name, handle, and bio; cue Success.
func ProfileSettings(c *ui.Context, profile *Profile, o ProfileSettingsOptions) (root, save ui.Element) {
	if profile == nil || o.Limit < 0 {
		panic("account: ProfileSettings needs a profile and a bio limit of at least 0")
	}
	def(&o.Title, "Profile")
	def(&o.Name, "Name")
	def(&o.Handle, "Handle")
	def(&o.Bio, "Bio")
	def(&o.Save, "Save")
	def(&o.Saved, "Profile saved")
	def(&o.Count, "%d/%d")
	if o.Limit == 0 {
		o.Limit = 160
	}
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	saved := ui.Local(root, "saved", func() bool { return false })
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			face := strings.TrimSpace(profile.Name)
			def(&face, "?")
			display.Avatar(c, face, nil).Size(t.Space(14), t.Space(14)).FontSize(t.Rem(1.1)).Shrink(0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, o.Title).FontSize(t.Rem(1.05)).Bold().FillWidth().MinWidth(0).SingleLine()
				kit.Muted(c, profile.Handle).FontSize(t.Rem(0.8)).FillWidth().MinWidth(0).SingleLine()
			})
		})
		used := utf8.RuneCountInString(profile.Bio)
		var bio ui.Element
		ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
			name := field(c, o.Name, "user", &profile.Name, false)
			handle := field(c, o.Handle, "at-sign", &profile.Handle, false)
			ui.Grid(c).Columns(1).FillWidth().Gap(t.Space(1)).Children(func() {
				frame := ui.Column(c).RowStart(2).ColumnStart(1).FillWidth().MinWidth(0).Gap(t.Space(1)).
					Padding(t.Space(2), t.Space(3)).Radius(t.Radius)
				frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
				frame.Children(func() {
					bio = ui.TextAreaBase(c, &profile.Bio).Label(o.Bio).Placeholder(o.Bio).FontSize(t.Rem(0.88)).FillWidth().Lines(3, 6)
					count := kit.Numeric(c, fmt.Sprintf(o.Count, used, o.Limit)).FontSize(t.Rem(0.72)).AlignSelf(ui.End).Shrink(0)
					count.TextColor(count.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(used > o.Limit)), ui.Snappy))
				})
				ui.Text(c, o.Bio).RowStart(1).ColumnStart(1).FontSize(t.Rem(0.78)).Bold().TextColor(t.TextMuted)
			})
			if name.Changed() || handle.Changed() || bio.Changed() {
				*saved = false
				root.SetChanged()
			}
		})
		valid := strings.TrimSpace(profile.Name) != "" && strings.TrimSpace(profile.Handle) != "" && used <= o.Limit
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				if *saved {
					c.Key("saved")
					note(c, "circle-check", o.Saved, t.Success)
				}
			})
			save = primary(c, o.Save, "check").Disabled(!valid)
		})
		if save.Clicked() && valid {
			*saved = true
			save.SetChanged()
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("ProfileSettings", "saved", "bioRunes", used)
		}
	})
	return root, save
}
