package main

type event interface{}

type mouseEvent struct{ mx, my int }

func (me *mouseEvent) Position() (int, int) {
	return me.mx, me.my
}

type keyEvent struct{ key rune }

func (ke *keyEvent) Rune() rune {
	return ke.key
}

type specialEvent struct{ key specialKey }

func (se *specialEvent) Key() specialKey {
	return se.key
}

type specialKey uint8

const (
	up specialKey = iota
	down
	left
	right
	backspace
	enter
	quit
	reset
)
