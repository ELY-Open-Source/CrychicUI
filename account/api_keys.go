package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// APIKey is one masked account API key.
type APIKey struct {
	ID, Name, Prefix, Secret string
	Revoked                  bool
}

// APIKeysOptions are the API key labels.
type APIKeysOptions struct{ Copy, Copied, Revoke, Revoked string }

// APIKeys copies and revokes masked keys; cue Pop, Delete.
func APIKeys(c *ui.Context, keys *[]APIKey, o APIKeysOptions) ui.Element {
	if keys == nil || len(*keys) == 0 {
		panic("account: APIKeys needs keys")
	}
	def(&o.Copy, "Copy")
	def(&o.Copied, "Copied")
	def(&o.Revoke, "Revoke")
	def(&o.Revoked, "Revoked")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleList)
	copied := ui.Local(root, "copied", func() string { return "" })
	copiedSecret := ui.Local(root, "copiedSecret", func() string { return "" })
	root.Children(func() {
		for i, key := range *keys {
			if key.ID == "" || key.Name == "" || key.Prefix == "" || key.Secret == "" {
				panic("account: invalid API key")
			}
			c.Key(key.ID)
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, key.Name)
					kit.Muted(c, key.Prefix+"••••").Font("monospace")
				})
				copyLabel := o.Copy
				if *copied == key.ID && *copiedSecret == key.Secret {
					copyLabel = o.Copied
				}
				copyButton := button(c, key.ID+" "+copyLabel).Disabled(key.Revoked)
				copyButton.Children(func() { ui.Text(c, copyLabel) })
				if copyButton.Clicked() && !key.Revoked {
					c.WriteClipboard(key.Secret)
					*copied = key.ID
					*copiedSecret = key.Secret
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("APIKeys", "copied", "index", i)
				}
				revokeLabel := o.Revoke
				if key.Revoked {
					revokeLabel = o.Revoked
				}
				revoke := button(c, key.ID+" "+revokeLabel).Disabled(key.Revoked)
				revoke.Children(func() { ui.Text(c, revokeLabel) })
				if revoke.Clicked() && !key.Revoked {
					(*keys)[i].Revoked = true
					root.SetChanged()
					sound.Play(sound.Delete)
					kit.Log("APIKeys", "revoked", "index", i)
				}
			})
		}
	})
	return root
}
