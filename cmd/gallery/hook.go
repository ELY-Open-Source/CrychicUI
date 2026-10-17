//go:build !(js && wasm)

package main

import "git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/gallery"

func hook(*gallery.App) {}
