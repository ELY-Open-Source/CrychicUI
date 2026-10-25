package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// APIKey is one masked account API key; Created is optional.
type APIKey struct {
	ID, Name, Prefix, Secret, Created string
	Revoked                           bool
}

// APIKeysOptions word the card; Label takes action and key name.
type APIKeysOptions struct{ Title, Copy, Copied, Revoke, Revoked, Label string }

// APIKeys copies and revokes masked keys; cue Pop, Delete.
func APIKeys(c *ui.Context, keys *[]APIKey, o APIKeysOptions) ui.Element {
	if keys == nil || len(*keys) == 0 {
		panic("account: APIKeys needs keys")
	}
	ids := make([]string, len(*keys))
	for i, key := range *keys {
		if key.Name == "" || key.Prefix == "" || key.Secret == "" {
			panic(fmt.Sprintf("account: APIKeys %q needs a name, prefix and secret", key.ID))
		}
		ids[i] = key.ID
	}
	unique("APIKeys", ids)
	def(&o.Title, "API keys")
	def(&o.Copy, "Copy")
	def(&o.Copied, "Copied")
	def(&o.Revoke, "Revoke")
	def(&o.Revoked, "Revoked")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	copied := ui.Local(root, "copied", func() string { return "" })
	copiedSecret := ui.Local(root, "copiedSecret", func() string { return "" })
	root.Children(func() {
		live := 0
		for _, key := range *keys {
			if !key.Revoked {
				live++
			}
		}
		kit.CardHead(c, "key", o.Title, func() { kit.Chip(c, fmt.Sprint(live), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, key := range *keys {
				c.Key(key.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ink := t.Accent
					if key.Revoked {
						ink = t.TextMuted
					}
					tile(c, "key", ink, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						name := ui.Text(c, key.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						if key.Revoked {
							name.TextColor(t.TextMuted).Strikethrough()
						}
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							kit.Muted(c, key.Prefix+"••••••••").Font("monospace").FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
							if key.Created != "" {
								dot(c)
								kit.Muted(c, key.Created).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
							}
						})
					})
					if key.Revoked {
						kit.Chip(c, o.Revoked, t.Danger)
						return
					}
					done := *copied == key.ID && *copiedSecret == key.Secret
					text, mark := o.Copy, "copy"
					if done {
						text, mark = o.Copied, "check"
					}
					copyButton := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, text, key.Name)).Size(t.Space(8), t.Space(8)).
						Shrink(0).Center().Radius(999))
					copyButton.Background(copyButton.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(copyButton.Hovered())), ui.Snappy))
					copyButton.TextColor(copyButton.AnimateColor("ink", t.TextMuted.Mix(t.Success, unit(done)), ui.Snappy))
					copyButton.Children(func() {
						c.Key(done)
						kit.Arrive(glyph(c, mark, 16), ui.Pose{Scale: 0.5, Fade: 1}, 0)
					})
					if copyButton.Clicked() {
						c.WriteClipboard(key.Secret)
						*copied = key.ID
						*copiedSecret = key.Secret
						root.SetChanged()
						sound.Play(sound.Pop)
						kit.Log("APIKeys", "copied", "index", i)
					}
					revoke := outline(c, fmt.Sprintf(o.Label, o.Revoke, key.Name), o.Revoke, "", t.Danger)
					if revoke.Clicked() {
						(*keys)[i].Revoked = true
						root.SetChanged()
						sound.Play(sound.Delete)
						kit.Log("APIKeys", "revoked", "index", i)
					}
				})
			}
		})
	})
	return root
}
