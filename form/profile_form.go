package form

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Profile is what a ProfileForm edits; Avatar may be nil.
type Profile struct {
	Name, Bio string
	Links     []string
	Avatar    *ui.Bitmap
}

// ProfileFormOptions tune a ProfileForm.
type ProfileFormOptions struct {
	// Labels, button copy and the name error.
	Name, Bio, Links, Link, Remove, AddLink, Save, NameError string
}

// ProfileForm edits name, bio and links beside an avatar.
func ProfileForm(c *ui.Context, value *Profile, o ProfileFormOptions) (root, save ui.Element) {
	if value == nil {
		panic("form: ProfileForm needs a value")
	}
	def(&o.Name, "Name")
	def(&o.Bio, "Bio")
	def(&o.Links, "Links")
	def(&o.Link, "Link")
	def(&o.Remove, "Remove link")
	def(&o.AddLink, "Add link")
	def(&o.Save, "Save profile")
	def(&o.NameError, "Enter your name.")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(4))
	ids := ui.Local(root, "rows", func() *rows { return &rows{} })
	keys := (*ids).keys(len(value.Links))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Children(func() {
			if blank(value.Name) && value.Avatar == nil {
				display.IconBox(c, "user", display.IconBoxOptions{}).Radius(999)
			} else {
				display.Avatar(c, value.Name, value.Avatar)
			}
			ui.Box(c).Grow(1).Children(func() {
				entry(c, o.Name, when(blank(value.Name), o.NameError), func() ui.Element {
					in := ui.TextInput(c, &value.Name).FillWidth()
					edited(root, in)
					return in
				})
			})
		})
		FieldGroup(c, o.Bio, FieldGroupOptions{}, func() {
			edited(root, ui.TextArea(c, &value.Bio).FillWidth().Lines(3, 6))
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
			ui.Text(c, o.Links).FontWeight(500)
			gone := -1
			for i, key := range keys {
				ui.Row(c).Key(key).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Transition(rowMotion(c)).Children(func() {
					edited(root, ui.TextInput(c, &value.Links[i]).Label(o.Link).Grow(1))
					if button.IconButton(c, o.Remove, "x").Clicked() {
						gone = i
					}
				})
			}
			if button.Ghost(c, o.AddLink).Clicked() {
				value.Links = append(value.Links, "")
				(*ids).add()
				root.SetChanged()
				kit.Log("ProfileForm", "links", "count", len(value.Links))
			}
			if gone >= 0 {
				value.Links = append(value.Links[:gone], value.Links[gone+1:]...)
				(*ids).drop(gone)
				root.SetChanged()
				kit.Log("ProfileForm", "links", "count", len(value.Links))
			}
		})
		save = primary(c, o.Save).Disabled(blank(value.Name))
	})
	if kit.Clicked(save, sound.Success) {
		kit.Log("ProfileForm", "saved", "links", len(value.Links))
	}
	return root, save
}
