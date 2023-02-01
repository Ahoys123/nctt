//go:build exclude

package main

import "syscall/js"

type WebWindow struct {
	changeLog map[[2]int]struct {
		r rune
		s string
	}
	doc         js.Value
	DrawingRect *Rect
	w, h        int
	evChan      chan event
}

// TODO: implement in javascript
// getContent(x, y) -> char + style
// setContent(x, y, char, style)
//
// IF ERR: check if w.func words with js.FuncOf
// call onMouseMove(x, y)
// call onKeyPressed(char)
// call onSpecialEvent(string)

func NewWebWindow(w, h int) *WebWindow {
	ww := &WebWindow{make(map[[2]int]struct {
		r rune
		s string
	}), js.Global().Get("document"), &Rect{0, 1, w, h - 1}, w, h, make(chan event)}

	DrawOverlay(ww)
	return ww
}

func (w *WebWindow) SetContent(x, y int, c rune, s style) {
	if !(0 <= x && x < w.w && 0 <= y && y < w.h) {
		return
	}
	w.changeLog[[2]int{x, y}] = struct {
		r rune
		s string
	}{c, s.String()}
}
func (w *WebWindow) GetContent(x, y int) (rune, style) {
	if v, ok := w.changeLog[[2]int{x, y}]; ok {
		return v.r, DeString(v.s)
	}
	v := js.Global().Call("getContent", x, y)
	f := ' '
	for _, c := range v.Index(0).String() {
		f = c
		break
	}
	return f, DeString(v.Index(1).String())
}

// should only be called once; registers js callbacks too
func (w *WebWindow) ChannelEvents() (chan event, chan struct{}) {
	js.Global().Set("onMouseMove", js.FuncOf(w.newMouseEvent))
	js.Global().Set("onKeyPressed", js.FuncOf(w.newKeyEvent))
	js.Global().Set("onSpecialKey", js.FuncOf(w.newSpecialEvent))
	return w.evChan, make(chan struct{})
}

func (w *WebWindow) GetDrawingRect() *Rect {
	return w.DrawingRect
}
func (w *WebWindow) GetWidth() int  { return w.w }
func (w *WebWindow) GetHeight() int { return w.h }

func (w *WebWindow) HideCursor()         {}
func (w *WebWindow) ShowCursor(int, int) {}

func (w *WebWindow) Show() {
	if len(w.changeLog) == 0 {
		return
	}
	for k, v := range w.changeLog {
		js.Global().Call("setContent", k[0], k[1], string(v.r), v.s)
		delete(w.changeLog, k)
	}
	js.Global().Call("show")
}
func (w *WebWindow) Fini() {}
func (w *WebWindow) Sync() {}

func (w *WebWindow) newMouseEvent(this js.Value, args []js.Value) any {
	w.evChan <- &mouseEvent{args[0].Int(), args[1].Int()}
	return nil
}

func (w *WebWindow) newKeyEvent(this js.Value, args []js.Value) any {
	f := ' '
	for _, c := range args[0].String() {
		f = c
		break
	}
	w.evChan <- &keyEvent{f}
	return nil
}

func (w *WebWindow) newSpecialEvent(this js.Value, args []js.Value) any {
	var sk specialKey
	switch args[0].String() {
	case "up":
		sk = up
	case "down":
		sk = down
	case "left":
		sk = left
	case "right":
		sk = right
	case "backspace":
		sk = backspace
	case "enter":
		sk = enter
	case "quit":
		sk = quit
	case "reset":
		sk = reset
	}
	w.evChan <- &specialEvent{sk}
	return nil
}
