package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type TermWindow struct {
	Screen        tcell.Screen
	DrawableRect  *Rect
	width, height int
}

func NewTermWindow(width, height int) *TermWindow {
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
	w := TermWindow{s, &Rect{0, 0, width, height - 1}, width, height}

	DrawOverlay(&w)

	return &w
}

func (w *TermWindow) ChannelEvents() (evChan chan event, quit chan struct{}) {
	evChan = make(chan event)
	quit = make(chan struct{})
	go w.pipeline(evChan, quit)
	return
}

func (w *TermWindow) pipeline(evChan chan<- event, quit chan struct{}) {
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

func (w *TermWindow) GetContent(x, y int) (rune, style) {
	a, _, b, _ := w.Screen.GetContent(x, y)
	return a, tcellToStyle(b)
}

func (w *TermWindow) ResolutionCheck(e *tcell.EventResize, old *VirtualRegion) *VirtualRegion {
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
		DrawOverlay(w)
		old.PasteContent(w.DrawableRect.x, w.DrawableRect.y, w)
	}
	return nil
}

func (w *TermWindow) displayResolutionWarning(width, height int) {
	DrawText("79x20 is min size.\nPlease expand your terminal.", &Rect{0, 0, width, height}, normal, w)
}

func (w *TermWindow) SetContent(x, y int, r rune, s style) {
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

func (w *TermWindow) HideCursor() {
	w.Screen.HideCursor()
}

func (w *TermWindow) ShowCursor(x, y int) {
	w.Screen.ShowCursor(x, y)
}

func (w *TermWindow) Show() {
	w.Screen.Show()
}

func (w *TermWindow) Fini() {
	w.Screen.Fini()
}

func (w *TermWindow) Sync() {
	w.Screen.Sync()
}

func (w *TermWindow) GetWidth() int {
	a, _ := w.Screen.Size()
	return a
}

func (w *TermWindow) GetHeight() int {
	_, a := w.Screen.Size()
	return a
}

func (w *TermWindow) GetDrawingRect() *Rect {
	return w.DrawableRect
}
