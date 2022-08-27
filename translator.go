package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

/*
type Translator map[string]string

type navajo struct {
	word        string
	translation string
}

func (n *navajo) getWord() string {
	return n.word
}

func (n *navajo) getTranslation() string {
	return n.translation
}

func toCode(toTranslate string) string {
	tr := ""
	for i, char := range toTranslate {
		if letters[char] != nil {
			tr += strings.Title(letters[char].getTranslation())
			if i != len(toTranslate)-1 {
				tr += "-"
			}
		}
	}
	return tr
}
*/

type translation struct {
	trans string
	color tcell.Color
}

func (tr translation) getTranslation() string {
	return tr.trans
}

func (tr translation) getColor() tcell.Color {
	return tr.color
}

var master ReplaceMap = map[string]translation{
	"a": {"wóláchííʼ\nant", tcell.ColorGoldenrod},
	"b": {"shash\nbear", tcell.ColorGoldenrod},
	"c": {"mósí\ncat", tcell.ColorGoldenrod},
	"d": {"bįįh\ndeer", tcell.ColorGoldenrod},
	"e": {"dzééh\nelk", tcell.ColorGoldenrod},
	"f": {"mąʼii\nfox", tcell.ColorGoldenrod},
	"g": {"tłʼízí\ngoat", tcell.ColorGoldenrod},
	"h": {"chʼah\nhat", tcell.ColorGoldenrod},
	"i": {"tin\nice", tcell.ColorGoldenrod},
	"j": {"téliichoʼí\njackass", tcell.ColorGoldenrod},
	"k": {"tłʼízí yázhí\nkid", tcell.ColorGoldenrod},
	"l": {"ajáád\nleg", tcell.ColorGoldenrod},
	"m": {"naʼatsʼǫǫsí\nmouse", tcell.ColorGoldenrod},
	"n": {"tsah\nneedle", tcell.ColorGoldenrod},
	"o": {"tłʼohchin\nonion", tcell.ColorGoldenrod},
	"p": {"bisóodi\npig", tcell.ColorGoldenrod},
	"q": {"kʼaaʼ yeiłtįįh\nquiver", tcell.ColorGoldenrod},
	"r": {"gah\nrabbit", tcell.ColorGoldenrod},
	"s": {"dibé\nsheep", tcell.ColorGoldenrod},
	"t": {"dééh\ntea", tcell.ColorGoldenrod},
	"u": {"shidáʼí\nuncle", tcell.ColorGoldenrod},
	"v": {"akʼehdidlíní\nvictor", tcell.ColorGoldenrod},
	"w": {"dlǫ́ʼii\nweasel", tcell.ColorGoldenrod},
	"x": {"ałnáʼázdzoh\ncross", tcell.ColorGoldenrod},
	"y": {"tsáʼásziʼ\nyucca", tcell.ColorGoldenrod},
	"z": {"béésh dootłʼizh\nzinc", tcell.ColorGoldenrod},

	"ant":     {"wóláchííʼ", tcell.ColorPink},
	"bear":    {"shash", tcell.ColorPink},
	"cat":     {"mósí", tcell.ColorPink},
	"deer":    {"bįįh", tcell.ColorPink},
	"elk":     {"dzééh", tcell.ColorPink},
	"fox":     {"mąʼii", tcell.ColorPink},
	"goat":    {"tłʼízí", tcell.ColorPink},
	"hat":     {"chʼah", tcell.ColorPink},
	"ice":     {"tin", tcell.ColorPink},
	"jackass": {"téliichoʼí", tcell.ColorPink},
	"kid":     {"tłʼízí yázhí", tcell.ColorPink},
	"leg":     {"ajáád", tcell.ColorPink},
	"mouse":   {"naʼatsʼǫǫsí", tcell.ColorPink},
	"needle":  {"tsah", tcell.ColorPink},
	"onion":   {"tłʼohchin", tcell.ColorPink},
	"pig":     {"bisóodi", tcell.ColorPink},
	"quiver":  {"kʼaaʼ yeiłtįįh", tcell.ColorPink},
	"rabbit":  {"gah", tcell.ColorPink},
	"sheep":   {"dibé", tcell.ColorPink},
	"tea":     {"dééh", tcell.ColorPink},
	"uncle":   {"shidáʼí", tcell.ColorPink},
	"victor":  {"akʼehdidlíní", tcell.ColorPink},
	"weasel":  {"dlǫ́ʼii", tcell.ColorPink},
	"cross":   {"ałnáʼázdzoh", tcell.ColorPink},
	"yucca":   {"tsáʼásziʼ", tcell.ColorPink},
	"zinc":    {"béésh dootłʼizh", tcell.ColorPink},

	"wóláchííʼ":       {"ant", tcell.ColorRed},
	"shash":           {"bear", tcell.ColorRed},
	"mósí":            {"cat", tcell.ColorRed},
	"bįįh":            {"deer", tcell.ColorRed},
	"dzééh":           {"elk", tcell.ColorRed},
	"mąʼii":           {"fox", tcell.ColorRed},
	"tłʼízí":          {"goat", tcell.ColorRed},
	"chʼah":           {"hat", tcell.ColorRed},
	"tin":             {"ice", tcell.ColorRed},
	"téliichoʼí":      {"jackass", tcell.ColorRed},
	"tłʼízí yázhí":    {"kid", tcell.ColorRed},
	"ajáád":           {"leg", tcell.ColorRed},
	"naʼatsʼǫǫsí":     {"mouse", tcell.ColorRed},
	"tsah":            {"needle", tcell.ColorRed},
	"tłʼohchin":       {"onion", tcell.ColorRed},
	"bisóodi":         {"pig", tcell.ColorRed},
	"kʼaaʼ yeiłtįįh":  {"quiver", tcell.ColorRed},
	"gah":             {"rabbit", tcell.ColorRed},
	"dibé":            {"sheep", tcell.ColorRed},
	"dééh":            {"tea", tcell.ColorRed},
	"shidáʼí":         {"uncle", tcell.ColorRed},
	"akʼehdidlíní":    {"victor", tcell.ColorRed},
	"dlǫ́ʼii":         {"weasel", tcell.ColorRed},
	"ałnáʼázdzoh":     {"cross", tcell.ColorRed},
	"tsáʼásziʼ":       {"yucca", tcell.ColorRed},
	"béésh dootłʼizh": {"zinc", tcell.ColorRed},

	"yókeed":      {"ask", tcell.ColorBlue},
	"naakáí":      {"company\n\"Mexican\"", tcell.ColorBlue},
	"hohkááh":     {"come", tcell.ColorBlue},
	"tó nilį́į́h": {"creek", tcell.ColorBlue},
	"iron fish":   {"béésh łóóʼ", tcell.ColorBlue},
	"whale":       {"łóóʼtsoh", tcell.ColorBlue},
	"tsídii":      {"bird", tcell.ColorBlue},
	"mobba yéhé":  {"it transports", tcell.ColorBlue},
}

// Replacer is basically a read only map.
type Replacer interface {
	getText(string) string
	getColor(string) tcell.Color
}

type ReplaceMap map[string]translation

func (rm ReplaceMap) getText(text string) string {
	return rm[strings.ToLower(text)].getTranslation()
}

func (rm ReplaceMap) getColor(text string) tcell.Color {
	return rm[strings.ToLower(text)].getColor()
}

func HoverReplace(textS string, rplcr Replacer, rect *Rect, w *Window) (pus []*PopUp, dcs []drawCall) {
	text := []rune(textS)
	start := 0
	highlight := false
	offset := 0 // where to draw text; - is before "true" pos, + is after
	nextOffset := -1
	var r rune
	for i := 0; i <= len(text); i++ {
		if i < len(text) {
			r = text[i]
		}
		if (i == len(text) || (r == '{' || r == '}')) && (start-i <= 1) {
			sub := string(text[start:i]) // substring

			dcs = append(dcs, drawCall{
				sub,
				rect,
				start + offset,
				tcell.StyleDefault.Foreground(rplcr.getColor(strings.ToLower(sub))),
			})

			if highlight {
				for _, v := range GetDrawingRect(sub, rect, start+offset) {
					dictVal := rplcr.getText(strings.ToLower(sub))
					width, height := GetDimensions(dictVal)

					pus = append(pus, NewPopUp(
						dictVal,
						width,
						height,
						v,
						w,
					))
				}
			}

			offset += nextOffset
			nextOffset = -1
			start = i + 1
			highlight = !highlight
		} else if r == '\n' {
			x := ((i + offset + 2 + nextOffset) % rect.w)
			if x != 0 {
				nextOffset += rect.w - x
			}
		} else if r == '\t' {
			offset -= 1
			nextOffset += 4
		}
	}

	return pus, dcs
}
