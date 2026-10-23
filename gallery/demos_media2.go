package gallery

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/media"
	"github.com/egoist/mygo/ui"
)

var mediaEntries2 = []Entry{
	{794, "Metronome", "media", "A tempo with a swinging arm; ticks play.", demoMetronome},
	{795, "Equalizer", "media", "Band sliders under a redrawn curve.", demoEqualizer},
	{796, "PhotoGrid", "media", "Photos in mixed sizes; tiles lift.", demoPhotoGrid},
	{797, "AlbumCard", "media", "A cover whose play button rises on hover.", demoAlbumCard},
	{798, "CoverFlow", "media", "Covers turned in 3D around the chosen one.", demoCoverFlow},
	{799, "Stories", "media", "Slides under segment bars; a tap skips.", demoStories},
	{800, "Collage", "media", "One large image beside the rest.", demoCollage},
	{801, "Polaroids", "media", "Tilted prints; the top one lifts.", demoPolaroids},
	{802, "Filmstrip", "media", "Frames on film under a playhead.", demoFilmstrip},
	{803, "Scrubber", "media", "A timeline with a preview frame.", demoScrubber},
	{804, "Captions", "media", "Caption lines timed to a clock.", demoCaptions},
	{805, "Chapters", "media", "Chapters that jump the clock.", demoChapters},
	{806, "SpeedMenu", "media", "Playback speeds in a menu.", demoSpeedMenu},
	{807, "QualityMenu", "media", "Quality choices with a check.", demoQualityMenu},
	{808, "Tuner", "media", "A dial over a sliding band.", demoTuner},
	{809, "VideoFrame", "media", "A 16:9 poster with play and controls.", demoVideoFrame},
	{810, "Turntable", "media", "A record that spins; the arm swings on.", demoTurntable},
	{811, "Cassette", "media", "A tape whose reels turn.", demoCassette},
	{812, "Visualizer", "media", "Bars ringed around a cover.", demoVisualizer},
	{813, "MediaGrid", "media", "Photos and clips with type badges.", demoMediaGrid},
	{814, "ImageStack", "media", "Fanned images; a click brings the next.", demoImageStack},
	{815, "AvatarUpload", "media", "An avatar framed in a circle with zoom.", demoAvatarUpload},
}

func demoMetronome(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	bpm, playing := numberState(root, "bpm", 96), boolState(root, "playing", false)
	root.Children(func() {
		if media.Metronome(c, "Metronome", bpm, playing, media.MetronomeOptions{}).Changed() {
			log(fmt.Sprintf("%.0f BPM, playing %v", *bpm, *playing))
		}
	})
}

func demoEqualizer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	gains := ui.Local(root, "gains", func() []float64 { return []float64{3, 1, 0, -2, 4} })
	root.Children(func() {
		if media.Equalizer(c, "Equalizer", []string{"60 Hz", "250 Hz", "1 kHz", "4 kHz", "12 kHz"}, *gains, media.EqualizerOptions{}).Changed() {
			log(fmt.Sprintf("gains %v dB", *gains))
		}
	})
}

func demoPhotoGrid(c *ui.Context, log func(string)) {
	p := photos()
	ui.Column(c).Width(420).Children(func() {
		media.PhotoGrid(c, "Trip photos", slices.Concat(p, p, p[:1]), media.PhotoGridOptions{})
	})
}

func demoAlbumCard(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(16)
	on := [2]*bool{boolState(root, "a", false), boolState(root, "b", false)}
	root.Children(func() {
		for i, k := range mediaTracks()[:2] {
			c.Key(i)
			if media.AlbumCard(c, k.Title, k.Artist, k.Cover, on[i], media.AlbumCardOptions{}).Changed() {
				log(fmt.Sprintf("%s playing %v", k.Title, *on[i]))
			}
		}
	})
}

func demoCoverFlow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	at := intState(root, "at", 2)
	p := photos()
	root.Children(func() {
		if media.CoverFlow(c, "Albums", slices.Concat(p, p[:2]), at, media.CoverFlowOptions{}).Changed() {
			log(fmt.Sprintf("cover %d", *at+1))
		}
	})
}

func demoStories(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(260)
	at := intState(root, "at", 0)
	p := photos()
	root.Children(func() {
		slide := func(i int) { ui.Image(c, p[i]).Fill().Fit(ui.Cover) }
		if media.Stories(c, "Stories", len(p), at, slide, media.StoriesOptions{}).Changed() {
			log(fmt.Sprintf("story %d", *at+1))
		}
	})
}

func demoCollage(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() { media.Collage(c, "Moodboard", photos(), media.CollageOptions{}) })
}

func demoPolaroids(c *ui.Context, log func(string)) {
	media.Polaroids(c, "Prints", photos(), []string{"Dunes", "Harbor", "Ferns"}, media.PolaroidsOptions{})
}

func demoFilmstrip(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	pos := ui.Local(root, "pos", func() time.Duration { return 4 * time.Second })
	p := photos()
	root.Children(func() {
		media.Filmstrip(c, "Reel", slices.Concat(p, p), *pos, 12*time.Second, media.FilmstripOptions{})
		if button.Outline(c, "Advance 2 s").Clicked() {
			*pos = (*pos + 2*time.Second) % (12 * time.Second)
			log(fmt.Sprintf("at %v", *pos))
		}
	})
}

func demoScrubber(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	pos := ui.Local(root, "pos", func() time.Duration { return 0 })
	root.Children(func() {
		if media.Scrubber(c, "Timeline", photos(), pos, time.Minute, media.ScrubberOptions{}).Changed() {
			log(fmt.Sprintf("at %v", *pos))
		}
	})
}

func demoCaptions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	pos := ui.Local(root, "pos", func() time.Duration { return time.Second })
	lines := []media.Caption{
		{From: 0, To: 2 * time.Second, Text: "The tide comes in at dusk."},
		{From: 2 * time.Second, To: 4 * time.Second, Text: "Boats turn their lamps on."},
		{From: 4 * time.Second, To: 6 * time.Second, Text: "Then the harbor goes quiet."},
	}
	root.Children(func() {
		media.Captions(c, "Captions", lines, *pos, media.CaptionsOptions{})
		if button.Outline(c, "Next line").Clicked() {
			*pos = (*pos + 2*time.Second) % (6 * time.Second)
			log(fmt.Sprintf("at %v", *pos))
		}
	})
}

func demoChapters(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	pos := ui.Local(root, "pos", func() time.Duration { return 75 * time.Second })
	chs := []media.Chapter{{Title: "Cold open", At: 0}, {Title: "The harbor", At: time.Minute}, {Title: "Night shift", At: 3 * time.Minute}}
	root.Children(func() {
		if media.Chapters(c, "Chapters", chs, pos, media.ChaptersOptions{}).Changed() {
			log(fmt.Sprintf("jumped to %v", *pos))
		}
	})
}

func demoSpeedMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	speed, open := numberState(root, "speed", 1), boolState(root, "open", false)
	root.Children(func() {
		if media.SpeedMenu(c, "Speed", []float64{0.5, 0.75, 1, 1.25, 1.5, 2}, speed, open, media.SpeedMenuOptions{}).Changed() {
			log(fmt.Sprintf("speed %g×, open %v", *speed, *open))
		}
	})
}

func demoQualityMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at, open := intState(root, "at", 0), boolState(root, "open", false)
	names := []string{"Auto", "1080p", "720p", "480p"}
	root.Children(func() {
		if media.QualityMenu(c, "Quality", names, at, open, media.QualityMenuOptions{}).Changed() {
			log(fmt.Sprintf("quality %s, open %v", names[*at], *open))
		}
	})
}

func demoTuner(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	freq := numberState(root, "freq", 98)
	root.Children(func() {
		if media.Tuner(c, "Station", freq, 87.5, 108, media.TunerOptions{}).Changed() {
			log(fmt.Sprintf("%.1f MHz", *freq))
		}
	})
}

func demoVideoFrame(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	pos := ui.Local(root, "pos", func() time.Duration { return 0 })
	playing := boolState(root, "playing", false)
	root.Children(func() {
		if media.VideoFrame(c, "Harbor at dusk", photos()[1], pos, 90*time.Second, playing, media.VideoFrameOptions{}).Changed() {
			log(fmt.Sprintf("playing %v at %v", *playing, pos.Truncate(time.Second)))
		}
	})
}

// deck shows a deck above a play button.
func deck(c *ui.Context, log func(string), view func(on bool)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	on := boolState(root, "on", false)
	root.Children(func() {
		view(*on)
		name := "Play"
		if *on {
			name = "Stop"
		}
		if button.Outline(c, name).Clicked() {
			*on = !*on
			log(fmt.Sprintf("playing %v", *on))
		}
	})
}

func demoTurntable(c *ui.Context, log func(string)) {
	deck(c, log, func(on bool) { media.Turntable(c, "Turntable", photos()[2], on, media.TurntableOptions{}) })
}

func demoCassette(c *ui.Context, log func(string)) {
	deck(c, log, func(on bool) { media.Cassette(c, "Cassette", "Side A · Summer", on, media.CassetteOptions{}) })
}

func demoVisualizer(c *ui.Context, log func(string)) {
	deck(c, log, func(on bool) { media.Visualizer(c, "Visualizer", photos()[0], on, media.VisualizerOptions{}) })
}

func demoMediaGrid(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	picked := intState(root, "picked", -1)
	p := photos()
	items := []media.MediaItem{
		{Name: "Dunes", Image: p[0]}, {Name: "Harbor reel", Image: p[1], Length: 42 * time.Second},
		{Name: "Ferns", Image: p[2]}, {Name: "Night drive", Image: p[1], Length: 2*time.Minute + 5*time.Second},
	}
	root.Children(func() {
		if media.MediaGrid(c, "Library", items, picked, media.MediaGridOptions{}).Changed() {
			log("picked " + items[*picked].Name)
		}
	})
}

func demoImageStack(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	top := intState(root, "top", 0)
	root.Children(func() {
		if media.ImageStack(c, "Photo stack", photos(), top, media.ImageStackOptions{}).Changed() {
			log(fmt.Sprintf("photo %d on top", *top+1))
		}
	})
}

func demoAvatarUpload(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	f := ui.Local(root, "framing", func() media.Framing { return media.Framing{Zoom: 1.5} })
	root.Children(func() {
		box, save := media.AvatarUpload(c, "Profile photo", photos()[1], f, media.AvatarUploadOptions{})
		if box.Changed() {
			log(fmt.Sprintf("zoom %.1f, pan %.2f, %.2f", f.Zoom, f.X, f.Y))
		}
		if save.Clicked() {
			log("saved")
		}
	})
}
