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

func (w *Window) ChannelEvents() (evChan chan event, quit chan struct{}) {
	evChan = make(chan event)
	quit = make(chan struct{})
	go pipeline(evChan, quit)
	return
}

func pipeline(evChan chan<- event, quit chan struct{}) {
	ichan := make(chan tcell.Event)
	iquit := make(chan struct{})
	go w.Screen.ChannelEvents(ichan, iquit)
	for {
		select {
		case <-quit:
			iquit <- struct{}{}
			close(evChan)
			close(quit)
			return
		case e := <-ichan:
			evChan <- tcellToEvent(e)
		}
	}
}

func tcellToEvent(e tcell.Event) event {
	switch ev := e.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			return &keyEvent{ev.Rune()}
		case tcell.KeyUp:
			return &specialEvent{up}
		case tcell.KeyDown:
			return &specialEvent{down}
		case tcell.KeyLeft:
			return &specialEvent{left}
		case tcell.KeyRight:
			return &specialEvent{right}
		case tcell.KeyBackspace2, tcell.KeyBackspace:
			return &specialEvent{backspace}
		case tcell.KeyEnter:
			return &specialEvent{enter}
		case tcell.KeyCtrlC:
			return &specialEvent{quit}
		case tcell.KeyESC:
			return &specialEvent{reset}
		}
	case *tcell.EventMouse:
		x, y := ev.Position()
		return &mouseEvent{x, y}
	}
	return &keyEvent{}
}

func (w *Window) DrawBoxAround(r *Rect, s style) {

	w.SetContent(r.x-1, r.y-1, '•', s)
	w.SetContent(r.x+r.w, r.y-1, '•', s)
	w.SetContent(r.x-1, r.y+r.h, '•', s)
	w.SetContent(r.x+r.w, r.y+r.h, '•', s)

	for x := r.x; x < r.x+r.w; x++ {
		w.SetContent(x, r.y-1, '-', s)
		w.SetContent(x, r.y+r.h, '-', s)
	}
	for y := r.y; y < r.y+r.h; y++ {
		w.SetContent(r.x-1, y, '|', s)
		w.SetContent(r.x+r.w, y, '|', s)
	}
}

func (w *Window) DrawText(text string, r *Rect, s style) {
	w.DrawTextOffset(text, r, 0, s)
}

func (w *Window) MarginRect(x, y, height int) *Rect {
	return &Rect{w.DrawableRect.x + x + 1, w.DrawableRect.y + y + 1, w.DrawableRect.w - 1 - (w.DrawableRect.x + x + 1), height}
}

func (w *Window) DrawTextOffset(text string, r *Rect, offset int, s style) {
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

		w.SetContent(r.x+x, r.y+y, textRune[curChar], s)
		curChar++
	}
}

func (w *Window) FillRect(with rune, r *Rect) {
	for x := r.x; x < r.x+r.w; x++ {
		for y := r.y; y < r.y+r.h; y++ {
			w.SetContent(x, y, with, normal)
		}
	}
}

func (w *Window) DrawOverlay() {
	w.DrawText("[SPACE] to advance", &Rect{0, 0, 18, 1}, option)
	w.DrawText("[ESCAPE] to title", &Rect{w.width - 17, 0, 17, 1}, option)
}

func (w *Window) DrawDebug(text string, y int) {
	w.DrawText(text, &Rect{1, y, 100, 10}, normal)
}

func (w *Window) GetContent(x, y int) (rune, style) {
	a, _, b, _ := w.Screen.GetContent(x, y)
	return a, tcellToStyle(b)
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
	w.DrawText("79x20 is min size.\nPlease expand your terminal.", &Rect{0, 0, width, height}, normal)
}

func (w *Window) SetContent(x, y int, r rune, s style) {
	var st tcell.Style
	switch s {
	case normal:
		st = tcell.StyleDefault
	case option:
		st = tcell.StyleDefault.Foreground(tcell.ColorYellow)

	case popup:
		st = tcell.StyleDefault.Foreground(tcell.ColorTurquoise)

	case popupBox:
		st = tcell.StyleDefault.Foreground(tcell.ColorDarkTurquoise)

	case t1ln:
		st = tcell.StyleDefault.Foreground(tcell.ColorGoldenrod)

	case t1en:
		st = tcell.StyleDefault.Foreground(tcell.ColorPink)

	case t1ne:
		st = tcell.StyleDefault.Foreground(tcell.ColorRed)

	case t2ne:
		st = tcell.StyleDefault.Foreground(tcell.ColorBlue)
	}

	w.Screen.SetContent(x, y, r, nil, st)
}

type style uint8

const (
	normal style = iota
	option
	popup
	popupBox
	t1ln
	t1en
	t1ne
	t2ne
)

func tcellToStyle(s tcell.Style) style {
	switch s {
	case tcell.StyleDefault:
		return normal
	case tcell.StyleDefault.Foreground(tcell.ColorYellow):
		return option
	case tcell.StyleDefault.Foreground(tcell.ColorTurquoise):
		return popup
	case tcell.StyleDefault.Foreground(tcell.ColorDarkTurquoise):
		return popupBox
	case tcell.StyleDefault.Foreground(tcell.ColorGoldenrod):
		return t1ln
	case tcell.StyleDefault.Foreground(tcell.ColorPink):
		return t1en
	case tcell.StyleDefault.Foreground(tcell.ColorRed):
		return t1ne
	case tcell.StyleDefault.Foreground(tcell.ColorBlue):
		return t2ne
	}
	return normal
}

func (w *Window) HideCursor() {
	w.Screen.HideCursor()
}

func (w *Window) ShowCursor(x, y int) {
	w.Screen.ShowCursor(x, y)
}

func (w *Window) Show() {
	w.Screen.Show()
}

func (w *Window) Fini() {
	w.Screen.Fini()
}

func (w *Window) Sync() {
	w.Screen.Sync()
}
