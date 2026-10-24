package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ArtifactKind picks an artifact's glyph.
type ArtifactKind uint8

// Artifact kinds.
const (
	Document ArtifactKind = iota
	Code
	Image
	Sheet
	Web
)

var artifactGlyphs = [...]string{"file-text", "file-code", "file-image", "file-spreadsheet", "globe"}

// ArtifactOptions name the open button.
type ArtifactOptions struct {
	Open string
}

// Artifact is a card of a generated file with open.
func Artifact(c *ui.Context, name string, kind ArtifactKind, meta string, o ArtifactOptions) (root, open ui.Element) {
	if name == "" || kind > Web {
		panic(fmt.Sprintf("agent: Artifact %q needs a name and a known kind", name))
	}
	def(&o.Open, "Open")
	t := c.Theme()
	root = kit.Card(c).FillWidth().Clip().Role(ui.RoleGroup).Label(name)
	kit.Lift(root)
	root.Children(func() {
		preview := ui.Box(c).FillWidth().Height(t.Space(20)).Center().Background(t.Accent.Alpha(0.06))
		preview.Children(func() {
			ic := ui.Icon(c, icon.Must(artifactGlyphs[kind])).Size(t.Space(10), t.Space(10)).TextColor(t.Accent)
			s := ic.AnimateTo("hover", 1+0.08*unit(root.Hovered()), ui.Snappy)
			ic.Scale(s, s)
		})
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).Bold().SingleLine()
				if meta != "" {
					kit.Muted(c, meta).FontSize(t.Rem(0.8)).SingleLine()
				}
			})
			open = kit.Press(ui.ButtonBase(c).Label(o.Open).Row().Gap(t.Space(1)).AlignItems(ui.Center).
				Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius))
			open.Background(open.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(open.Hovered())), ui.Snappy)).TextColor(t.AccentText).
				Children(func() {
					ui.Text(c, o.Open).Bold().Shrink(0)
					ui.Icon(c, icon.Must("arrow-up-right")).Size(16, 16)
				})
		})
	})
	if open.Clicked() {
		root.SetChanged()
		sound.Play(sound.Open)
		kit.Log("Artifact", "opened", "kind", kind)
	}
	return root, open
}
