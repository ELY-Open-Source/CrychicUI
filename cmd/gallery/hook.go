//go:build !(js && wasm)

package main

import "github.com/ELY-Open-Source/CrychicUI/gallery"

func hook(*gallery.App) {}
