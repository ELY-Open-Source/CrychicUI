//go:build js && wasm

package main

import (
	"syscall/js"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/gallery"
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/sound"
)

// hook exports window.__crychic for Playwright.
func hook(app *gallery.App) {
	obj := js.Global().Get("Object").New()
	js.Global().Get("Object").Call("defineProperty", obj, "ready", map[string]any{
		"get": js.FuncOf(func(js.Value, []js.Value) any { return app.Snapshot().Built }),
	})
	obj.Set("page", js.FuncOf(func(js.Value, []js.Value) any { return app.Snapshot().Page }))
	obj.Set("stats", js.FuncOf(func(js.Value, []js.Value) any {
		return map[string]any{"components": len(gallery.Entries()), "categories": len(gallery.Categories)}
	}))
	obj.Set("recent", js.FuncOf(func(js.Value, []js.Value) any {
		var out []any
		for _, c := range sound.Recent() {
			out = append(out, c.String())
		}
		return out
	}))
	js.Global().Set("__crychic", obj)
}
