package devtools

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Secret is one stored secret; Updated is optional.
type Secret struct{ Name, Value, Updated string }

// SecretsOptions word the panel; Reveal and Hide take the name.
type SecretsOptions struct{ Title, Reveal, Hide, Empty string }

// mask hides a secret without telling its length.
const mask = "••••••••••••"

// Secrets lists masked secrets; reveal reports a change; cue Toggle.
func Secrets(c *ui.Context, secrets []Secret, o SecretsOptions) ui.Element {
	def(&o.Title, "Secrets")
	names := make([]string, len(secrets))
	for i, s := range secrets {
		if s.Value == "" {
			panic(fmt.Sprintf("devtools: Secrets %q needs a value", s.Name))
		}
		names[i] = s.Name
	}
	if len(secrets) > 0 {
		unique("Secrets", names)
	}
	def(&o.Reveal, "Reveal %s")
	def(&o.Hide, "Hide %s")
	def(&o.Empty, "No secrets")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	shown := ui.Local(root, "shown", func() map[string]bool { return map[string]bool{} })
	root.Children(func() {
		kit.CardHead(c, "lock", o.Title, func() { kit.Chip(c, fmt.Sprint(len(secrets)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(secrets) == 0 {
				empty(c, "lock", o.Empty)
			}
			for i, s := range secrets {
				open := (*shown)[s.Name]
				c.Key(s.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ink, mark := t.TextMuted, "lock"
					if open {
						ink, mark = t.Warning, "lock-open"
					}
					c.Key(open)
					kit.Arrive(tile(c, mark, ink, t.Space(8)), ui.Pose{Scale: 0.8, Fade: 1}, 0)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.75)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, s.Name).Font("monospace").FontSize(t.Rem(0.78)).Bold().Shrink(0).MaxWidth(t.Space(40)).SingleLine()
							if s.Updated != "" {
								dot(c)
								kit.Muted(c, s.Updated).FontSize(t.Rem(0.68)).MinWidth(0).SingleLine()
							}
						})
						value := mask
						if open {
							value = s.Value
						}
						c.Key(open)
						v := ui.Text(c, value).Font("monospace").FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
						kit.Arrive(v, ui.Pose{X: -t.Space(1), Fade: 1}, 0)
						if !open {
							v.TextColor(t.TextMuted)
						}
					})
					label := fmt.Sprintf(o.Reveal, s.Name)
					if open {
						label = fmt.Sprintf(o.Hide, s.Name)
					}
					b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(8), t.Space(8)).Center().Radius(999).Shrink(0).Checked(open))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Warning.Alpha(0.15), unit(open)), ui.Snappy))
					b.Children(func() {
						eye := "eye"
						if open {
							eye = "eye-off"
						}
						glyph(c, eye, 16).TextColor(t.TextMuted)
					})
					if b.Clicked() {
						(*shown)[s.Name] = !open
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("Secrets", "revealed", "index", i, "shown", !open)
					}
				})
			}
		})
	})
	return root
}
