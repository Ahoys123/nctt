package main

import "github.com/gdamore/tcell/v2"

type Rect struct {
	x, y, w, h int
}

func (r *Rect) Contains(x, y int) bool {
	return r.x <= x && x < r.x+r.w && r.y <= y && y < r.y+r.h
}

func GetDrawingRect(text string, r *Rect, offset int) (tr []*Rect) {
	textRune := []rune(text)
	x := offset % r.w
	startX := x
	ci := 0
	tr = append(tr, &Rect{startX + r.x, (offset / r.w) + r.y, 0, 1})

	i := 0
	for ; i < len(textRune); i++ {

		x = (i + offset) % r.w
		if textRune[i] == '\n' {
			if x != 0 {
				offset += r.w - x
			}
			offset -= 1
			i++
			continue
		} else if textRune[i] == '\t' {
			offset += 4
			i++
			continue
		}

		y := (i + offset) / r.w

		// new line, new rect!
		if i != 0 && x == 0 {
			tr[ci].w = r.w - startX
			ci++
			startX = 0
			tr = append(tr, &Rect{r.x, y + r.y, 0, 1})
		}

		if y >= r.h {
			return
		}
	}

	tr[ci].w = x + 1 - startX

	return tr
}

func GetDimensions(text string) (w, h int) {
	textRune := []rune(text)
	longestLine := 0
	w, h = 0, 1
	for i := 0; i < len(textRune); i++ {
		if textRune[i] == '\n' {
			h++
			if w < longestLine {
				w = longestLine
			}
			longestLine = 0
			continue
		}

		if textRune[i] == '\t' {
			longestLine += 4
			continue
		}

		longestLine++
	}

	if w < longestLine {
		w = longestLine
	}

	return w, h
}

type VirtualRegion [][]*pixel

type loc struct {
	x, y int
	pixel
}

type pixel struct {
	mainc rune
	combc []rune
	style tcell.Style
}

func CopyContent(r *Rect, w *Window) VirtualRegion {

	vr := make(VirtualRegion, r.h)
	for y := range vr {
		vr[y] = make([]*pixel, r.w)
		for x := range vr[y] {
			mainc, combc, style, _ := w.Screen.GetContent(x+r.x, y+r.y)
			vr[y][x] = &pixel{mainc, combc, style}
		}
	}

	return vr

}

func (vr VirtualRegion) GetContent(x, y int) *pixel {
	return vr[y][x]
}

func (vr VirtualRegion) PasteContent(x, y int, w *Window) {
	for ny := 0; ny < len(vr); ny++ {
		for nx := 0; nx < len(vr[ny]); nx++ {
			px := vr.GetContent(nx, ny)
			w.Screen.SetContent(x+nx, y+ny, px.mainc, px.combc, px.style)
		}
	}
}

func (vr VirtualRegion) GetDifferences(initial VirtualRegion) (tr []loc) {

	for y := range vr {
		for x, px := range vr[y] {
			if !px.Equals(initial.GetContent(x, y)) {
				tr = append(tr, loc{x, y, *px})
			}
		}
	}

	return tr
}

func (vr VirtualRegion) Equals(other VirtualRegion) bool {
	for y := range vr {
		for x, px := range vr[y] {
			if !px.Equals(other.GetContent(x, y)) {
				return false
			}
		}
	}
	return true
}

func (px *pixel) Equals(other *pixel) bool {
	if len(px.combc) != len(other.combc) {
		return false
	}
	for i, comb := range px.combc {
		if other.combc[i] != comb {
			return false
		}
	}
	return px.mainc == other.mainc
}
