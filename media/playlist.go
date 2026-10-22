package media

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Track is one song or episode; Cover may be nil.
type Track struct {
	Title, Artist string
	Length        time.Duration
	Cover         *ui.Bitmap
}

// PlaylistOptions tune a Playlist.
type PlaylistOptions struct{}

// Playlist lists tracks; the playing one shows moving bars.
func Playlist(c *ui.Context, label string, tracks []Track, playing *int, o PlaylistOptions) ui.Element {
	if label == "" || len(tracks) == 0 {
		panic(fmt.Sprintf("media: Playlist %q needs a label and tracks", label))
	}
	for _, k := range tracks {
		tracked("Playlist", k)
	}
	chosen("Playlist", label, playing, -1, len(tracks))
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Label(label)
	rows := make([]ui.Element, len(tracks))
	root.Children(func() {
		for i := range tracks {
			c.Key(i)
			rows[i] = kit.Press(ui.ButtonBase(c).Label(tracks[i].Title).FillWidth().Row().Gap(t.Space(3)).
				AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius))
			kit.Arrive(rows[i], ui.Pose{Fade: 1, Y: t.Space(1.5)}, time.Duration(i)*30*time.Millisecond)
		}
	})
	for i, r := range rows {
		if r.Clicked() && i != *playing {
			*playing = i
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("Playlist", "played", "index", i)
		}
	}
	for i, r := range rows {
		on := i == *playing
		bg := t.SurfaceHover.Alpha(unit(r.Hovered()))
		if on {
			bg = t.Accent.Alpha(0.1)
		}
		r.Checked(on).Background(r.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
			ui.Box(c).Width(t.Space(5)).Shrink(0).Center().Children(func() {
				if !on {
					c.Key("n")
					kit.Numeric(c, fmt.Sprint(i+1)).TextColor(t.TextMuted)
					return
				}
				c.Key("eq")
				eq := ui.Row(c).Height(t.Space(4)).Gap(t.Space(0.5)).AlignItems(ui.End)
				ph := float64(eq.Loop("eq", 900*time.Millisecond, ui.Linear))
				eq.Children(func() {
					for k := range 3 {
						h := 0.3 + 0.7*math.Abs(math.Sin(2*math.Pi*(ph+float64(k)/3)))
						ui.Box(c).Width(t.Space(0.75)).Height(t.Space(4) * float32(h)).Radius(999).Background(t.Accent)
					}
				})
			})
			fg := t.Text
			if on {
				fg = t.Accent
			}
			ui.Column(c).Grow(1).Children(func() {
				ui.Text(c, tracks[i].Title).SingleLine().Bold().TextColor(r.AnimateColor("fg", fg, ui.Snappy))
				kit.Muted(c, tracks[i].Artist).SingleLine().FontSize(t.Rem(0.8))
			})
			kit.Numeric(c, clock(tracks[i].Length)).TextColor(t.TextMuted)
		})
	}
	return root
}
