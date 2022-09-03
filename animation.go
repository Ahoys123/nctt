package main

import (
	"math/rand"
	"strings"
)

type Element interface {
	// Update is called every frame and updates the display.
	// The inputs registered between the last and current
	// frame will be passed in as an array of events.
	Update([]event)

	// Done outputs if the element is finished.
	Done() bool

	// Reset resets the element. Clearing screen artifacts
	// is NOT a requirement for reset. It is called on s
	Reset()
}

//#region WaitForInput

type WaitForNext struct {
	done bool
}

func NewWaitForNext() *WaitForNext {
	return &WaitForNext{false}
}

func (wfn *WaitForNext) Update(ec []event) {
	if wfn.done {
		return
	}

	if len(ec) > 0 { // if input
		switch ev := ec[0].(type) {
		case *keyEvent: // if key
			switch ev.Rune() {
			case ' ': // if space
				wfn.done = true
			}
		}
	}
}

func (wfn *WaitForNext) Done() bool {
	return wfn.done
}

func (wfn *WaitForNext) Reset() {
	wfn.done = false
}

//#endregion WaitForInput

//#region SlowText
type SlowText struct {
	text    []rune
	bound   *Rect
	curChar int
	w       *Window
	done    bool
	offset  int

	click     *SoundEffect
	ding      *SoundEffect
	playSound bool
}

func NewSlowText(text string, bound *Rect, w *Window) *SlowText {
	return &SlowText{[]rune(text), bound, 0, w, false, 0, nil, nil, false}
}

func NewTypewritter(text string, bound *Rect, w *Window, click *SoundEffect, ding *SoundEffect) *SlowText {
	return &SlowText{[]rune(text), bound, 0, w, false, 0, click, ding, true}
}

func (st *SlowText) Update(ec []event) {

	if st.done {
		return
	}

	if len(ec) > 0 { // if input
		switch ev := ec[0].(type) {
		case *keyEvent: // if key
			switch ev.Rune() {
			case ' ': // if space
				if !st.done { // if not displayed
					if st.playSound {
						st.ding.Play()
					}
					st.playSound = false
					for !st.done { // update until it is displayed
						st.Update([]event{})
					}
				} else { // if fully displayed
					st.done = true // we done
				}
			}
		}
	}

	if st.curChar >= len(st.text) {
		// past the last char or no more room, no more animation!
		st.done = true
		if st.playSound {
			st.ding.Play()
		}
		return
	}

	if st.text[st.curChar] == '\t' {
		st.offset += 3
		st.curChar++
		return
	}

	x := (st.curChar + st.offset) % st.bound.w
	// set next cell's content
	if st.text[st.curChar] == '\n' {
		st.offset += st.bound.w - x - 1
		st.curChar++
		if st.playSound {
			st.ding.Play()
		}
		return
	}
	y := (st.curChar + st.offset) / st.bound.w
	if y >= st.bound.h {
		st.done = true
		if st.playSound {
			st.ding.Play()
		}
		return
	}

	if st.playSound && rand.Float32() < 0.5 {
		st.click.Play()
	}
	st.w.SetContent(st.bound.x+x, st.bound.y+y, st.text[st.curChar], normal)
	st.curChar++
}

func (st *SlowText) Done() bool {
	return st.done
}

func (st *SlowText) Reset() {
	w.FillRect(' ', st.bound)

	st.curChar = 0
	st.done = false
	st.offset = 0
	st.playSound = st.click != nil
}

//#endregion SlowText

//#region ConcurrentPlayer

// ConcurrentPlayer plays all elements simultaniously.
// It transparently runs the Update() method on its
// children, and Done() is based on if all of its
// children are all individually done.
type ConcurrentPlayer struct {
	elms []Element
}

func NewConcurrentPlayer(elms []Element) *ConcurrentPlayer {
	return &ConcurrentPlayer{elms}
}

func (cp *ConcurrentPlayer) Update(ec []event) {
	for _, elm := range cp.elms {
		elm.Update(ec)
	}
}

func (cp *ConcurrentPlayer) Done() bool {
	for _, elm := range cp.elms {
		if !elm.Done() {
			return false
		}
	}
	return true
}

func (cp *ConcurrentPlayer) Reset() {
	for _, elm := range cp.elms {
		elm.Reset()
	}
}

//#endregion ConcurrentPlayer

//#region DiscretePlayer

type DiscretePlayer struct {
	elms        []Element
	curElmIndex int
	done        bool
}

func NewDiscretePlayer(elms []Element) *DiscretePlayer {
	return &DiscretePlayer{elms, 0, false}
}

func (dp *DiscretePlayer) Update(ec []event) {

	if dp.elms[dp.curElmIndex].Done() {
		dp.elms[dp.curElmIndex].Reset()
		dp.curElmIndex = (dp.curElmIndex + 1) % len(dp.elms)
	} else {
		dp.elms[dp.curElmIndex].Update(ec)
	}
}

func (dp *DiscretePlayer) Done() bool {
	return dp.done
}

func (dp *DiscretePlayer) Reset() {
	dp.curElmIndex = 0
	dp.done = false
	for _, elm := range dp.elms {
		elm.Reset()
	}
}

//#endregion DiscretePlayer

//#region SequentialPlayer

type SequentialPlayer struct {
	*DiscretePlayer
}

func NewSequentialPlayer(elms []Element) *SequentialPlayer {
	return &SequentialPlayer{NewDiscretePlayer(elms)}
}

func (sp *SequentialPlayer) Update(ec []event) {
	for elmIndex := 0; elmIndex <= sp.curElmIndex; elmIndex++ {
		sp.elms[elmIndex].Update(ec)
	}

	if !sp.done && sp.elms[sp.curElmIndex].Done() {
		if sp.curElmIndex < len(sp.elms)-1 {
			sp.curElmIndex++
		} else {
			sp.done = true
		}
	}
}

//#endregion SequentialPlayer

//#region Popup

type PopUp struct {
	text          string
	bound         *Rect
	width, height int

	w          *Window
	hovBox     *Rect
	repContent VirtualRegion

	showNext bool
	mx, my   int
}

func NewPopUp(text string, width, height int, bound *Rect, w *Window) *PopUp {
	return &PopUp{text, bound, width, height, w, nil, nil, false, 0, 0}
}

func (pu *PopUp) Update(ec []event) {

	if pu.showNext {
		screenW, screenH := pu.w.width, pu.w.height
		pu.mx += 2
		pu.my += 2
		if pu.mx+pu.width >= screenW {
			pu.mx = screenW - pu.width - 2
		}
		if pu.my+pu.height >= screenH {
			pu.my = screenH - pu.height - 2
		}

		pu.hovBox = &Rect{pu.mx, pu.my, pu.width, pu.height}
		pu.repContent = CopyContent(&Rect{pu.hovBox.x - 1, pu.hovBox.y - 1, pu.hovBox.w + 2, pu.hovBox.h + 2}, pu.w)
		w.FillRect(' ', pu.hovBox)
		w.DrawBoxAround(pu.hovBox, popupBox)
		w.DrawText(pu.text, pu.hovBox, popupBox)
		pu.showNext = false
	}

	if len(ec) > 0 { // if input
		switch ev := ec[0].(type) {
		case *mouseEvent: // if mouse
			if pu.hovBox != nil {
				pu.repContent.PasteContent(pu.hovBox.x-1, pu.hovBox.y-1, pu.w)
				pu.hovBox = nil
				pu.repContent = nil
			}

			pu.showNext = false

			pu.mx, pu.my = ev.Position()
			if !pu.bound.Contains(pu.mx, pu.my) {
				return
			}

			pu.showNext = true
		}
	}
}

func (pu *PopUp) Done() bool {
	return true
}

func (pu *PopUp) Reset() {
	if pu.hovBox != nil {
		w.FillRect(' ', &Rect{pu.hovBox.x - 1, pu.hovBox.y - 1, pu.hovBox.w + 2, pu.hovBox.h + 2})
		pu.hovBox = nil
		pu.repContent = nil
	}

	pu.showNext = false
}

//#endregion PopUp

//#region HoverText

type HoverText struct {
	// raw data

	// text given by user; words in {brackets}
	// will be run through dictonary, pop up will
	// show outputed text.
	rawText string
	dict    Replacer
	w       *Window

	// autogenerated
	hovs      []*PopUp // hoverregions
	drawCalls []drawCall
	done      bool
}

type drawCall struct {
	text   string
	rect   *Rect
	offset int
	style  style
}

func (dc *drawCall) Draw(w *Window) {
	w.DrawTextOffset(dc.text, dc.rect, dc.offset, dc.style)
}

func NewHoverText(text string, r *Rect, dict Replacer, w *Window) *HoverText {
	hovs, segs := HoverReplace(text, dict, r, w)
	return &HoverText{text, dict, w, hovs, segs, false}
}

func (ht *HoverText) Update(ec []event) {
	if !ht.done {
		// draw all text
		for _, htdc := range ht.drawCalls {
			htdc.Draw(ht.w)
		}

		ht.done = true
	}

	// set up all popup regions
	for _, pu := range ht.hovs {
		pu.Update(ec)
	}
}

func (ht *HoverText) Done() bool {
	return ht.done
}

func (ht *HoverText) Reset() {
	ht.done = false
	for _, pu := range ht.hovs {
		pu.Reset()
	}
	for _, dc := range ht.drawCalls {
		w.FillRect(' ', dc.rect)
	}
}

//#endregion HoverText

//#region Checker

type Checkable interface {
	Element
	Selection() string
}

type Checker struct {
	chk          Checkable
	correct      []string
	right, wrong Element

	state int // 0 = nothing, 1 = right input displaying, 2 = wrong input displaying
	done  bool
}

func NewChecker(chk Checkable, correct []string, right, wrong Element) *Checker {
	thing := []string{}
	for _, v := range correct {
		thing = append(thing, strings.ReplaceAll(strings.ToLower(v), " ", ""))
	}
	return &Checker{chk, thing, right, wrong, 0, false}
}

func (chkr *Checker) Update(ec []event) {
	if chkr.done {
		return
	}

	switch chkr.state {
	case 1: // right displaying state
		if chkr.right.Done() {
			chkr.state = 0
			chkr.done = true
			w.HideCursor()
		}
		chkr.right.Update(ec)
		return
	case 2: // wrong displaying state
		if !chkr.wrong.Done() {
			chkr.wrong.Update(ec)
		} else {
			chkr.state = 0
			chkr.wrong.Reset()
			chkr.chk.Reset()
		}
		return
	}

	if !chkr.chk.Done() {
		chkr.chk.Update(ec)
		return
	}

	// chk is done, we can actually check now!
	if contains(chkr.correct, strings.ReplaceAll(strings.ToLower(chkr.chk.Selection()), " ", "")) {
		chkr.state = 1
	} else {
		chkr.state = 2
	}
}

func contains(x []string, y string) bool {
	for _, v := range x {
		if v == y {
			return true
		}
	}
	return false
}

func (chkr *Checker) Done() bool {
	return chkr.done
}

func (chkr *Checker) Reset() {
	w.HideCursor()
	chkr.chk.Reset()
	chkr.right.Reset()
	chkr.wrong.Reset()
	chkr.state = 0
	chkr.done = false
}

//#endregion Checker

//#region Options

type Options struct {
	options []string
	w       *Window

	selected  int
	done      bool
	drawCalls []drawCall
}

func NewOptions(options []string, r *Rect, w *Window) *Options {
	drawCalls := []drawCall{}
	row := r.y
	for _, option := range options {
		width, height := GetDimensions(option)
		drawCalls = append(drawCalls, drawCall{option, &Rect{r.x + 2, row, width, height}, 0, normal})
		row += height + 1
	}
	return &Options{options, w, 0, false, drawCalls}
}

func (op *Options) Update(ec []event) {
	if op.done {
		return
	}

	if len(ec) > 0 { // if input
		switch ev := ec[0].(type) {
		case *specialEvent: // if key
			switch ev.Key() {
			case up, left:
				op.changeIndex(-1)
			case down, right:
				op.changeIndex(1)
			case enter:
				op.done = true
			}

		case *keyEvent:
			if k := ev.Rune(); k == ' ' || k == 'c' {
				op.done = true
			}
		}
	}

	// draw options
	for _, dc := range op.drawCalls {
		dc.Draw(op.w)
	}

	// draw > <
	selRect := op.drawCalls[op.selected].rect
	w.SetContent(selRect.x-2, selRect.y, '>', option)
	w.SetContent(selRect.x+selRect.w+1, selRect.y, '<', option)
}

func (op *Options) changeIndex(by int) {
	selRect := op.drawCalls[op.selected].rect
	w.SetContent(selRect.x-2, selRect.y, ' ', normal)
	w.SetContent(selRect.x+selRect.w+1, selRect.y, ' ', normal)

	if op.selected+by < 0 {
		op.selected = (op.selected + by + len(op.options)) % len(op.options)
		return
	}
	op.selected = (op.selected + by) % len(op.options)
}

func (op *Options) Done() bool {
	return op.done
}

func (op *Options) Reset() {
	selRect := op.drawCalls[op.selected].rect
	w.SetContent(selRect.x-2, selRect.y, ' ', normal)
	w.SetContent(selRect.x+selRect.w+1, selRect.y, ' ', normal)

	op.done = false
	op.selected = 0
	for _, dc := range op.drawCalls {
		w.FillRect(' ', dc.rect)
	}
}

func (op *Options) Selection() string {
	if !op.done {
		return ""
	}
	return op.options[op.selected]
}

//#endregion Options

//#region UserInput

type TextInput struct {
	uir      *Rect
	userText []rune

	selectionReturn string
	done            bool
	curmx           int
	click, ding     *SoundEffect
}

func NewTextInput(uir *Rect) *TextInput {
	return NewTypewritterInput(uir, nil, nil)
}

func NewTypewritterInput(uir *Rect, click, ding *SoundEffect) *TextInput {
	return &TextInput{uir, nil, "", false, 0, click, ding}
}

func (ti *TextInput) Update(ec []event) {
	if ti.done {
		w.HideCursor()
		return
	}

	if len(ec) > 0 { // if input
		switch ev := ec[0].(type) {
		case *keyEvent: // if key
			ti.selectionReturn = ""
			r := ev.Rune()
			ti.userText = append(ti.userText, r)

			if ti.click != nil {
				ti.click.Play()
			}

			w.FillRect(' ', ti.uir)
			w.DrawText(string(ti.userText), ti.uir, normal)
			ti.curmx++
			w.ShowCursor(ti.uir.x+ti.curmx, ti.uir.y)
		case *specialEvent:
			switch ev.Key() {
			case backspace:
				ti.selectionReturn = ""
				if len(ti.userText) == 0 {
					break
				}

				ti.userText = ti.userText[:len(ti.userText)-1]
				ti.curmx--

				w.FillRect(' ', ti.uir)
				w.DrawText(string(ti.userText), ti.uir, normal)
				w.ShowCursor(ti.uir.x+ti.curmx, ti.uir.y)
			case enter:
				ti.selectionReturn = string(ti.userText)
				if ti.ding != nil {
					ti.ding.Play()
				}
				ti.done = true
			}
		}
	}
}

func (ti *TextInput) Done() bool {
	return ti.done
}

func (ti *TextInput) Selection() string {
	return ti.selectionReturn
}

func (ti *TextInput) Reset() {
	w.HideCursor()
	w.FillRect(' ', ti.uir)
	ti.userText = nil
	ti.curmx = 0
	ti.done = false
}

//#endregion UserInput
