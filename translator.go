package main

import (
	"strings"
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
	s     style
}

func (tr translation) getTranslation() string {
	return tr.trans
}

func (tr translation) getColor() style {
	return tr.s
}

var master ReplaceMap = map[string]translation{
	"a": {"wóláchííʼ\nant", t1ln},
	"b": {"shash\nbear", t1ln},
	"c": {"mósí\ncat", t1ln},
	"d": {"bįįh\ndeer", t1ln},
	"e": {"dzééh\nelk", t1ln},
	"f": {"mąʼii\nfox", t1ln},
	"g": {"tłʼízí\ngoat", t1ln},
	"h": {"chʼah\nhat", t1ln},
	"i": {"tin\nice", t1ln},
	"j": {"téliichoʼí\njackass", t1ln},
	"k": {"tłʼízí yázhí\nkid", t1ln},
	"l": {"ajáád\nleg", t1ln},
	"m": {"naʼatsʼǫǫsí\nmouse", t1ln},
	"n": {"tsah\nneedle", t1ln},
	"o": {"tłʼohchin\nonion", t1ln},
	"p": {"bisóodi\npig", t1ln},
	"q": {"kʼaaʼ yeiłtįįh\nquiver", t1ln},
	"r": {"gah\nrabbit", t1ln},
	"s": {"dibé\nsheep", t1ln},
	"t": {"dééh\ntea", t1ln},
	"u": {"shidáʼí\nuncle", t1ln},
	"v": {"akʼehdidlíní\nvictor", t1ln},
	"w": {"dlǫ́ʼii\nweasel", t1ln},
	"x": {"ałnáʼázdzoh\ncross", t1ln},
	"y": {"tsáʼásziʼ\nyucca", t1ln},
	"z": {"béésh dootłʼizh\nzinc", t1ln},

	"ant":     {"wóláchííʼ", t1en},
	"bear":    {"shash", t1en},
	"cat":     {"mósí", t1en},
	"deer":    {"bįįh", t1en},
	"elk":     {"dzééh", t1en},
	"fox":     {"mąʼii", t1en},
	"goat":    {"tłʼízí", t1en},
	"hat":     {"chʼah", t1en},
	"ice":     {"tin", t1en},
	"jackass": {"téliichoʼí", t1en},
	"kid":     {"tłʼízí yázhí", t1en},
	"leg":     {"ajáád", t1en},
	"mouse":   {"naʼatsʼǫǫsí", t1en},
	"needle":  {"tsah", t1en},
	"onion":   {"tłʼohchin", t1en},
	"pig":     {"bisóodi", t1en},
	"quiver":  {"kʼaaʼ yeiłtįįh", t1en},
	"rabbit":  {"gah", t1en},
	"sheep":   {"dibé", t1en},
	"tea":     {"dééh", t1en},
	"uncle":   {"shidáʼí", t1en},
	"victor":  {"akʼehdidlíní", t1en},
	"weasel":  {"dlǫ́ʼii", t1en},
	"cross":   {"ałnáʼázdzoh", t1en},
	"yucca":   {"tsáʼásziʼ", t1en},
	"zinc":    {"béésh dootłʼizh", t1en},

	"wóláchííʼ":       {"ant", t1ne},
	"shash":           {"bear", t1ne},
	"mósí":            {"cat", t1ne},
	"bįįh":            {"deer", t1ne},
	"dzééh":           {"elk", t1ne},
	"mąʼii":           {"fox", t1ne},
	"tłʼízí":          {"goat", t1ne},
	"chʼah":           {"hat", t1ne},
	"tin":             {"ice", t1ne},
	"téliichoʼí":      {"jackass", t1ne},
	"tłʼízí yázhí":    {"kid", t1ne},
	"ajáád":           {"leg", t1ne},
	"naʼatsʼǫǫsí":     {"mouse", t1ne},
	"tsah":            {"needle", t1ne},
	"tłʼohchin":       {"onion", t1ne},
	"bisóodi":         {"pig", t1ne},
	"kʼaaʼ yeiłtįįh":  {"quiver", t1ne},
	"gah":             {"rabbit", t1ne},
	"dibé":            {"sheep", t1ne},
	"dééh":            {"tea", t1ne},
	"shidáʼí":         {"uncle", t1ne},
	"akʼehdidlíní":    {"victor", t1ne},
	"dlǫ́ʼii":         {"weasel", t1ne},
	"ałnáʼázdzoh":     {"cross", t1ne},
	"tsáʼásziʼ":       {"yucca", t1ne},
	"béésh dootłʼizh": {"zinc", t1ne},

	"yókeed":      {"ask", t2ne},
	"naakáí":      {"company\n\"Mexican\"", t2ne},
	"hohkááh":     {"come", t2ne},
	"tó nilį́į́h": {"creek", t2ne},
	"iron fish":   {"béésh łóóʼ", t2ne},
	"whale":       {"łóóʼtsoh", t2ne},
	"tsídii":      {"bird", t2ne},
	"mobba yéhé":  {"it transports", t2ne},
}

// Replacer is basically a read only map.
type Replacer interface {
	getText(string) string
	getColor(string) style
}

type ReplaceMap map[string]translation

func (rm ReplaceMap) getText(text string) string {
	return rm[strings.ToLower(text)].getTranslation()
}

func (rm ReplaceMap) getColor(text string) style {
	return rm[strings.ToLower(text)].getColor()
}

func HoverReplace(textS string, rplcr Replacer, rect *Rect, w Window) (pus []*PopUp, dcs []drawCall) {
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
				rplcr.getColor(strings.ToLower(sub)),
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
