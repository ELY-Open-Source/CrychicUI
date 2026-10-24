package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CardControlsOptions word the card; notes sit under each switch.
type CardControlsOptions struct{ Title, Freeze, FreezeNote, Online, OnlineNote, Limit, Frozen string }

// CardControls switches freeze and online, sets a limit; cue Toggle.
func CardControls(c *ui.Context, frozen, online *bool, limit *float64, currency string, o CardControlsOptions) ui.Element {
	if frozen == nil || online == nil || limit == nil || *limit < 0 || *limit > 10000 || !finite(*limit) || currency == "" {
		panic("finance: CardControls needs states, a 0–10000 limit and a currency")
	}
	def(&o.Title, "Card controls")
	def(&o.Freeze, "Freeze card")
	def(&o.FreezeNote, "Block all new payments")
	def(&o.Online, "Online payments")
	def(&o.OnlineNote, "Allow web and in-app purchases")
	def(&o.Limit, "Spending limit")
	def(&o.Frozen, "Frozen")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "credit-card", o.Title, func() {
			if *frozen {
				kit.Arrive(kit.Chip(c, o.Frozen, t.Accent), ui.Pose{Scale: 0.8, Fade: 1}, 0)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2))
		list.Children(func() {
			for i, item := range []struct {
				name, note, glyph string
				value             *bool
			}{{o.Freeze, o.FreezeNote, "snowflake", frozen}, {o.Online, o.OnlineNote, "globe", online}} {
				c.Key(i)
				row := kit.Press(ui.ButtonBase(c).Label(item.name).Row().FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).
					Padding(t.Space(2)).Radius(t.Radius).Checked(*item.value))
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					tile(c, item.glyph, palette(t, i))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, item.name).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						kit.Muted(c, item.note).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
					})
					toggle(c, *item.value)
				})
				if row.Clicked() {
					*item.value = !*item.value
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("CardControls", "toggled", "control", i, "value", *item.value)
				}
			}
			ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
					tile(c, "gauge", palette(t, 2))
					ui.Text(c, o.Limit).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					money(c, *limit, currency).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
				})
				slider := ui.Slider(c, limit, 0, 10000).Label(o.Limit).FillWidth()
				if kit.Changed(slider, sound.Tick) {
					root.SetChanged()
					kit.Log("CardControls", "limit", "value", *limit)
				}
			})
		})
	})
	return root
}

// toggle draws a switch track whose knob slides.
func toggle(c *ui.Context, on bool) {
	t := c.Theme()
	track := ui.Box(c).Size(t.Space(9), t.Space(5)).Shrink(0).Radius(999)
	track.Background(track.AnimateColor("track", t.SurfaceHover.Mix(t.Accent, unit(on)), ui.Snappy))
	x := track.AnimateTo("knob", unit(on), ui.Snappy)
	track.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		d := r.H - t.Space(1)
		p.Fill(ui.Rect{X: r.X + t.Space(0.5) + (r.W-d-t.Space(1))*x, Y: r.Y + t.Space(0.5), W: d, H: d}, t.Background, 999)
	}))
}
