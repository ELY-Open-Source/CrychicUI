// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"time"
)

// CommitOptions are the copy knobs of Commit.
type CommitOptions struct {
	Copy      string
	FormatAge func(age time.Duration) string
}

// Commit is a commit row; cue Pop on copy.
func Commit(c *ui.Context, hash, message, author string, at time.Time, o CommitOptions) (root, copy ui.Element) {
	if hash == "" || message == "" || author == "" || at.IsZero() {
		panic("code: Commit needs a hash, message, author and time")
	}
	def(&o.Copy, "Copy hash")
	if o.FormatAge == nil {
		o.FormatAge = age
	}
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(message)
	root.Children(func() {
		ui.Text(c, short(hash)).Font("monospace").FontSize(t.Rem(0.85)).TextColor(t.Accent).
			Padding(t.Space(0.5), t.Space(1.5)).Background(t.Accent.Alpha(0.1)).Radius(t.Radius)
		ui.Text(c, message).SingleLine().Grow(1)
		kit.Muted(c, author).FontSize(t.Rem(0.8)).Shrink(0)
		kit.Muted(c, o.FormatAge(c.Now().Sub(at))).FontSize(t.Rem(0.8)).Shrink(0)
		copy = copyButton(c, o.Copy, hash)
	})
	return root, copy
}

// short is the seven-rune hash prefix.

func short(hash string) string {
	r := []rune(hash)
	if len(r) > 7 {
		return string(r[:7])
	}
	return hash
}

// PRState is a pull request state.
