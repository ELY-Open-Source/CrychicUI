// Command gallery opens the CrychicUI gallery on desktops and browsers.
package main

import (
	"log"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/gallery"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
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
