package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Window struct {
	Screen        tcell.Screen
	DrawableRect  *Rect
	width, height int
}

func NewWindow(width, height int) *Window {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.EnableMouse()
	s.SetStyle(defStyle)
	s.Clear()
	w := Window{s, &Rect{0, 1, width, height - 1}, width, height}

	w.DrawOverlay()

	return &w
}

func (w *Window) ChannelEvents() (evChan chan tcell.Event, quit chan struct{}) {
	evChan = make(chan tcell.Event)
	quit = make(chan struct{})
	go w.Screen.ChannelEvents(evChan, quit)
	return
}

func (w *Window) DrawBoxAround(r *Rect, style tcell.Style) {

	w.Screen.SetContent(r.x-1, r.y-1, '•', nil, style)
	w.Screen.SetContent(r.x+r.w, r.y-1, '•', nil, style)
	w.Screen.SetContent(r.x-1, r.y+r.h, '•', nil, style)
	w.Screen.SetContent(r.x+r.w, r.y+r.h, '•', nil, style)

	for x := r.x; x < r.x+r.w; x++ {
		w.Screen.SetContent(x, r.y-1, '-', nil, style)
		w.Screen.SetContent(x, r.y+r.h, '-', nil, style)
	}
	for y := r.y; y < r.y+r.h; y++ {
		w.Screen.SetContent(r.x-1, y, '|', nil, style)
		w.Screen.SetContent(r.x+r.w, y, '|', nil, style)
	}
}

func (w *Window) DrawText(text string, r *Rect, style tcell.Style) {
	w.DrawTextOffset(text, r, 0, style)
}

func (w *Window) MarginRect(x, y, height int) *Rect {
	return &Rect{w.DrawableRect.x + x + 1, w.DrawableRect.y + y + 1, w.DrawableRect.w - 1 - (w.DrawableRect.x + x + 1), height}
}

func (w *Window) DrawTextOffset(text string, r *Rect, offset int, style tcell.Style) {
	textRune := []rune(text)
	curChar := 0

	for curChar < len(textRune) {

		x := (curChar + offset) % r.w
		// set next cell's content
		if textRune[curChar] == '\n' {
			if x != 0 {
				offset += r.w - x
			}
			offset -= 1
			curChar++
			continue
		}
		if textRune[curChar] == '\t' {
			offset += 4
			curChar++
			continue
		}
		y := (curChar + offset) / r.w
		if y >= r.h {
			return
		}

		w.Screen.SetContent(r.x+x, r.y+y, textRune[curChar], nil, style)
		curChar++
	}
}

func (w *Window) FillRect(with rune, r *Rect) {
	for x := r.x; x < r.x+r.w; x++ {
		for y := r.y; y < r.y+r.h; y++ {
			w.Screen.SetContent(x, y, with, nil, tcell.StyleDefault)
		}
	}
}

func (w *Window) DrawOverlay() {
	w.DrawText("[SPACE] to advance", &Rect{0, 0, 18, 1}, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	w.DrawText("[ESCAPE] to title", &Rect{w.width - 17, 0, 17, 1}, tcell.StyleDefault.Foreground(tcell.ColorYellow))
}

func (w *Window) DrawDebug(text string, y int) {
	w.DrawText(text, &Rect{1, y, 100, 10}, tcell.StyleDefault.Foreground(tcell.ColorPurple))
}

func (w *Window) ResolutionCheck(e *tcell.EventResize, old *VirtualRegion) *VirtualRegion {
	width, height := e.Size()
	if width < w.width || height < w.height {
		if old != nil {
			return old
		}
		vr := CopyContent(w.DrawableRect, w)
		w.Screen.Clear()
		w.displayResolutionWarning(width, height)
		return &vr
	} else if old != nil {
		w.Screen.Clear()
		w.DrawOverlay()
		old.PasteContent(w.DrawableRect.x, w.DrawableRect.y, w)
	}
	return nil
}

func (w *Window) displayResolutionWarning(width, height int) {
	w.DrawText("79x20 is min size.\nPlease expand your terminal.", &Rect{0, 0, width, height}, tcell.StyleDefault)
}
