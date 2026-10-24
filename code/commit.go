package code

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CommitOptions are the copy knobs of Commit.
type CommitOptions struct {
	Copy string
	// Committed words the meta line, "%s" taking author then age.
	Committed string
	FormatAge func(age time.Duration) string
}

// Commit is an author, message, age and hash; cue Pop.
func Commit(c *ui.Context, hash, message, author string, at time.Time, o CommitOptions) (root, copy ui.Element) {
	if hash == "" || message == "" || author == "" || at.IsZero() {
		panic(fmt.Sprintf("code: Commit %q needs a hash, message, author and time", hash))
	}
	def(&o.Copy, "Copy hash")
	def(&o.Committed, "%s committed %s")
	if o.FormatAge == nil {
		o.FormatAge = age
	}
	t := c.Theme()
	root = kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2.5), t.Space(2), t.Space(2.5), t.Space(3)).Role(ui.RoleGroup).Label(message)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Row().AlignItems(ui.Center).Children(func() {
		display.Avatar(c, author, nil).Size(t.Space(8), t.Space(8)).FontSize(t.Rem(0.75)).Shrink(0)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, message).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
			kit.Muted(c, fmt.Sprintf(o.Committed, author, o.FormatAge(c.Now().Sub(at)))).FontSize(t.Rem(0.75)).SingleLine()
		})
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(t.Radius).
			Background(t.Accent.Alpha(0.1)).TextColor(t.Accent).Shrink(0).Children(func() {
			ui.Icon(c, iconCommit).Size(14, 14)
			ui.Text(c, short(hash)).Font("monospace").FontSize(t.Rem(0.8)).Shrink(0)
		})
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
