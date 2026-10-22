package gallery

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/media"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var mediaEntries = []Entry{
	{771, "ImageView", "media", "An image that fades in over a placeholder.", demoImageView},
	{772, "Gallery", "media", "A grid of images opening a lightbox.", demoMediaGallery},
	{773, "Carousel", "media", "Slides with arrows, dots and drags.", demoCarousel},
	{774, "Slideshow", "media", "Slides on a timer with a progress bar.", demoSlideshow},
	{775, "ThumbStrip", "media", "Thumbnails choosing the main image.", demoThumbStrip},
	{776, "ZoomImage", "media", "An image zoomed by wheel and dragged.", demoZoomImage},
	{777, "Compare", "media", "Before and after at a dragged divider.", demoCompare},
	{778, "Cropper", "media", "A crop box dragged and resized.", demoCropper},
	{779, "ImageFilters", "media", "Filter presets on a thumbnail strip.", demoImageFilters},
	{780, "AudioPlayer", "media", "Play and seek a synthesized tune.", demoAudioPlayer},
	{781, "Waveform", "media", "A clip's bars; the played part tinted.", demoWaveform},
	{782, "Spectrum", "media", "Bars moving like a spectrum.", demoSpectrum},
	{783, "PlayerControls", "media", "Previous, play, next, shuffle, repeat.", demoPlayerControls},
	{784, "SeekBar", "media", "A track with buffered and played parts.", demoSeekBar},
	{785, "Volume", "media", "A speaker that mutes and a slider.", demoVolume},
	{786, "Playlist", "media", "Tracks; the playing one moves.", demoPlaylist},
	{787, "NowPlaying", "media", "Cover, title, artist and progress.", demoNowPlaying},
	{788, "MiniBar", "media", "A compact player row.", demoMiniBar},
	{789, "Lyrics", "media", "Lines lit in time; scroll follows.", demoLyrics},
	{790, "PodcastCard", "media", "An episode with length and play.", demoPodcastCard},
	{791, "Recorder", "media", "Record, pause and stop with a timer.", demoRecorder},
	{792, "Piano", "media", "Keys that play notes.", demoPiano},
	{793, "DrumPad", "media", "Pads that play drums and flash.", demoDrumPad},
}

// mediaTracks are sample songs with gallery photos.
func mediaTracks() []media.Track {
	p := photos()
	return []media.Track{
		{Title: "Blue Hour", Artist: "Kai Mori", Length: 3*time.Minute + 24*time.Second, Cover: p[1]},
		{Title: "Dune Road", Artist: "Ana Reyes", Length: 4*time.Minute + 2*time.Second, Cover: p[0]},
		{Title: "Green Room", Artist: "The Ferns", Length: 2*time.Minute + 51*time.Second, Cover: p[2]},
	}
}

func demoImageView(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(12)
	since := ui.Local(root, "since", func() time.Time { return c.Now() })
	root.Children(func() {
		var img *ui.Bitmap
		if wait := 1200*time.Millisecond - c.Now().Sub(*since); wait > 0 {
			c.After(wait)
		} else {
			img = photos()[1]
		}
		media.ImageView(c, "Harbor at dusk", img, media.ImageViewOptions{})
		if button.Outline(c, "Reload").Clicked() {
			*since = c.Now()
			log("reloading")
		}
	})
}

func demoMediaGallery(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	open := intState(root, "open", -1)
	p := photos()
	root.Children(func() {
		if media.Gallery(c, "Trip photos", slices.Concat(p, p), open, media.GalleryOptions{}).Changed() {
			log(fmt.Sprintf("open: %d", *open))
		}
	})
}

func demoCarousel(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	at := intState(root, "at", 0)
	p := photos()
	root.Children(func() {
		slide := func(i int) { ui.Image(c, p[i]).Fill().Fit(ui.Cover) }
		if media.Carousel(c, "Featured", len(p), at, slide, media.CarouselOptions{Height: 220}).Changed() {
			log(fmt.Sprintf("slide %d", *at+1))
		}
	})
}

func demoSlideshow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	at := intState(root, "at", 0)
	p := photos()
	root.Children(func() {
		slide := func(i int) { ui.Image(c, p[i]).Fill().Fit(ui.Cover) }
		if media.Slideshow(c, "Highlights", len(p), at, slide, media.SlideshowOptions{Height: 220}).Changed() {
			log(fmt.Sprintf("slide %d", *at+1))
		}
	})
}

func demoThumbStrip(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	at := intState(root, "at", 0)
	root.Children(func() {
		if media.ThumbStrip(c, "Listing photos", photos(), at, media.ThumbStripOptions{}).Changed() {
			log(fmt.Sprintf("photo %d", *at+1))
		}
	})
}

func demoZoomImage(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	z := numberState(root, "zoom", 1)
	root.Children(func() {
		if media.ZoomImage(c, "Site map", photos()[2], z, media.ZoomImageOptions{Height: 240}).Changed() {
			log(fmt.Sprintf("zoom %.2f", *z))
		}
	})
}

func demoCompare(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	split := numberState(root, "split", 0.5)
	p := photos()
	root.Children(func() {
		if media.Compare(c, "Color grade", p[0], p[1], split, media.CompareOptions{}).Changed() {
			log(fmt.Sprintf("split %.2f", *split))
		}
	})
}

func demoCropper(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	crop := ui.Local(root, "crop", func() media.Crop { return media.Crop{X: 0.2, Y: 0.2, W: 0.5, H: 0.5} })
	root.Children(func() {
		if media.Cropper(c, "Cover photo", photos()[0], crop, media.CropperOptions{}).Changed() {
			log(fmt.Sprintf("crop %.2f, %.2f, %.2f × %.2f", crop.X, crop.Y, crop.W, crop.H))
		}
	})
}

func demoImageFilters(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	f := ui.Local(root, "filter", func() media.Filter { return media.Original })
	root.Children(func() {
		if media.ImageFilters(c, "Photo looks", photos()[1], f, media.ImageFiltersOptions{Height: 200}).Changed() {
			log(fmt.Sprintf("filter %d", *f))
		}
	})
}

func demoAudioPlayer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	pos := ui.Local(root, "pos", func() time.Duration { return 0 })
	playing := boolState(root, "playing", false)
	tune := media.Clip{Notes: []float64{261.63, 261.63, 392, 392, 440, 440, 392, 0}, Beat: 400 * time.Millisecond}
	root.Children(func() {
		if media.AudioPlayer(c, "Lullaby", tune, pos, playing, media.AudioPlayerOptions{}).Changed() {
			log(fmt.Sprintf("playing %v at %s", *playing, pos.Round(100*time.Millisecond)))
		}
	})
}

func demoWaveform(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	played := numberState(root, "played", 0.3)
	levels := make([]float64, 48)
	for i := range levels {
		levels[i] = 0.2 + 0.8*math.Abs(math.Sin(float64(i)*0.45)*math.Cos(float64(i)*0.13))
	}
	root.Children(func() {
		media.Waveform(c, "Voice memo", levels, *played, media.WaveformOptions{})
		if button.Outline(c, "Play 10% more").Clicked() {
			*played = math.Mod(*played+0.1, 1.0001)
			log(fmt.Sprintf("played %.0f%%", *played*100))
		}
	})
}

func demoSpectrum(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	playing := boolState(root, "playing", false)
	root.Children(func() {
		media.Spectrum(c, "Live audio", *playing, media.SpectrumOptions{})
		if button.PlayButton(c, playing, button.PlayButtonOptions{}).Changed() {
			log(fmt.Sprintf("playing %v", *playing))
		}
	})
}

func demoPlayerControls(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	playing, shuffle := boolState(root, "playing", false), boolState(root, "shuffle", false)
	rep := ui.Local(root, "repeat", func() media.Repeat { return media.RepeatOff })
	root.Children(func() {
		box, prev, next := media.PlayerControls(c, "Player", playing, shuffle, rep, media.PlayerControlsOptions{})
		if box.Changed() {
			log(fmt.Sprintf("playing %v, shuffle %v, repeat %d", *playing, *shuffle, *rep))
		}
		if prev.Clicked() {
			log("previous")
		}
		if next.Clicked() {
			log("next")
		}
	})
}

func demoSeekBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	pos := ui.Local(root, "pos", func() time.Duration { return 42 * time.Second })
	root.Children(func() {
		if media.SeekBar(c, "Episode position", pos, 90*time.Second, 3*time.Minute, media.SeekBarOptions{}).Changed() {
			log(fmt.Sprintf("seek %s", *pos))
		}
	})
}

func demoVolume(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	level := numberState(root, "level", 60)
	root.Children(func() {
		if media.Volume(c, "Volume", level, media.VolumeOptions{}).Changed() {
			log(fmt.Sprintf("volume %.0f", *level))
		}
	})
}

func demoPlaylist(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	playing := intState(root, "playing", -1)
	root.Children(func() {
		if media.Playlist(c, "Up next", mediaTracks(), playing, media.PlaylistOptions{}).Changed() {
			log(fmt.Sprintf("playing track %d", *playing+1))
		}
	})
}

func demoNowPlaying(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	pos := ui.Local(root, "pos", func() time.Duration { return 72 * time.Second })
	track := mediaTracks()[0]
	root.Children(func() {
		media.NowPlaying(c, track, *pos, media.NowPlayingOptions{})
		if button.Outline(c, "Skip 30 s").Clicked() {
			*pos = min(track.Length, *pos+30*time.Second)
			log(fmt.Sprintf("at %s", *pos))
		}
	})
}

func demoMiniBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	playing := boolState(root, "playing", false)
	track := mediaTracks()[1]
	root.Children(func() {
		if media.MiniBar(c, track, playing, 95*time.Second, media.MiniBarOptions{}).Changed() {
			log(fmt.Sprintf("playing %v", *playing))
		}
	})
}

func demoLyrics(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	pos := ui.Local(root, "pos", func() time.Duration { return 0 })
	words := []string{"The harbor lights are low", "The tide is coming in", "I count the boats by name",
		"And wait for you again", "The gulls have gone to sleep", "The rope is cold and wet",
		"I hum the song you knew", "The one I can't forget"}
	var lines []media.Line
	for i, w := range words {
		lines = append(lines, media.Line{At: time.Duration(i) * 4 * time.Second, Text: w})
	}
	root.Children(func() {
		media.Lyrics(c, "Harbor song", lines, *pos, media.LyricsOptions{Height: 180})
		if button.Outline(c, "Next line").Clicked() {
			*pos = (*pos + 4*time.Second) % (time.Duration(len(lines)) * 4 * time.Second)
			log(fmt.Sprintf("at %s", *pos))
		}
	})
}

func demoPodcastCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	playing := boolState(root, "playing", false)
	ep := media.Track{Title: "How lighthouses kept time", Artist: "Night Talk", Length: 41*time.Minute + 40*time.Second, Cover: photos()[1]}
	root.Children(func() {
		if media.PodcastCard(c, ep, playing, media.PodcastCardOptions{}).Changed() {
			log(fmt.Sprintf("playing %v", *playing))
		}
	})
}

func demoRecorder(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	state := ui.Local(root, "state", func() media.Take { return media.Idle })
	root.Children(func() {
		if media.Recorder(c, "Voice note", state, 0.55, media.RecorderOptions{}).Changed() {
			log(fmt.Sprintf("state %d", *state))
		}
	})
}

func demoPiano(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	note := intState(root, "note", -1)
	root.Children(func() {
		if media.Piano(c, "Piano", note, media.PianoOptions{}).Changed() {
			log(fmt.Sprintf("note %d", *note))
		}
	})
}

func demoDrumPad(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	hit := intState(root, "hit", -1)
	pads := []sound.Piece{sound.Kick, sound.Snare, sound.Hat, sound.Clap, sound.Tom}
	root.Children(func() {
		if media.DrumPad(c, "Drums", pads, hit, media.DrumPadOptions{}).Changed() {
			log(fmt.Sprintf("hit %s", pads[*hit]))
		}
	})
}
