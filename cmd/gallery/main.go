// Command gallery opens the CrychicUI gallery on desktops and browsers.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ELY-Open-Source/CrychicUI/gallery"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	stats := flag.Bool("stats", false, "print registry statistics")
	flag.Parse()
	if *stats {
		fmt.Printf("%d %d\n", len(gallery.Entries()), len(gallery.Categories))
		return
	}
	app := gallery.New("/")
	hook(app)
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "CrychicUI",
			Width:     1280,
			Height:    820,
			MinWidth:  720,
			MinHeight: 480,
			Content:   ui.View(app.View),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
