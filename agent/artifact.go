package agent

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ArtifactOptions name the open button.
type ArtifactOptions struct {
	Open string
}

// Artifact is a card of a generated file with open.
func Artifact(c *ui.Context, name, kind string, o ArtifactOptions) (root, open ui.Element) {
	if name == "" || kind == "" {
		panic("agent: Artifact needs a name and a kind")
	}
	def(&o.Open, "Open")
	t := c.Theme()
	root = kit.Card(c).AlignSelf(ui.Start).Width(t.Space(76))
	kit.Lift(root)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(t.Accent.Alpha(0.12)).
				TextColor(t.Accent).Children(func() { ui.Icon(c, icon.Must("file")).Size(16, 16) })
			ui.Column(c).Grow(1).Children(func() {
				ui.Text(c, name).Bold().SingleLine()
				kit.Muted(c, kind).FontSize(t.Rem(0.8)).SingleLine()
			})
			open = kit.Press(ui.ButtonBase(c).Label(o.Open).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Border(1, t.Accent).TextColor(t.Accent))
			open.Background(open.AnimateColor("bg", t.Accent.Alpha(0.12*unit(open.Hovered())), ui.Snappy)).
				Children(func() { ui.Text(c, o.Open).Bold().SingleLine() })
		})
	})
	if open.Clicked() {
		root.SetChanged()
		sound.Play(sound.Open)
		kit.Log("Artifact", "opened")
	}
	return root, open
}
