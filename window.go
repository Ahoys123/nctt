package main

type Window interface {
	SetContent(int, int, rune, style)
	GetDrawingRect() *Rect
	GetWidth() int
	GetHeight() int
	GetContent(int, int) (rune, style)
	HideCursor()
	ShowCursor(int, int)
	Show()
	Fini()
	Sync()
}

func DrawBoxAround(r *Rect, s style, w Window) {

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

func DrawText(text string, r *Rect, s style, w Window) {
	DrawTextOffset(text, r, 0, s, w)
}

func MarginRect(x, y, height int, w Window) *Rect {
	dr := w.GetDrawingRect()
	return &Rect{dr.x + x + 1, dr.y + y + 1, dr.w - 1 - (dr.x + x + 1), height}
}

func DrawTextOffset(text string, r *Rect, offset int, s style, w Window) {
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

func FillRect(with rune, r *Rect, w Window) {
	for x := r.x; x < r.x+r.w; x++ {
		for y := r.y; y < r.y+r.h; y++ {
			w.SetContent(x, y, with, normal)
		}
	}
}

func DrawOverlay(w Window) {
	DrawText("[SPACE] to advance", &Rect{0, 0, 18, 1}, option, w)
	DrawText("[ESCAPE] to title", &Rect{w.GetWidth() - 17, 0, 17, 1}, option, w)
}

func DrawDebug(text string, y int, w Window) {
	DrawText(text, &Rect{1, y, 100, 10}, normal, w)
}
