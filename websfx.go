//go:build exclude

package main

import (
	"syscall/js"
)

type WebSfx struct {
	filename string
}

func NewWebSfx(filename string) *WebSfx {
	return &WebSfx{filename}
}

func (wsfx *WebSfx) Play() {
	if wsfx.filename == "assets/ding.wav" {
		js.Global().Call("ding")
	} else {
		js.Global().Call("click")
	}
}
