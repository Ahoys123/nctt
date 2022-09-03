package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

var w *TermWindow

var scene Element

func main() {
	w = NewTermWindow(79, 20)
	scene = GetDemoScene(w) //GetMainScene(w)

	evChan, quit := w.ChannelEvents()
	run(w, evChan, quit)
}

func run(w *TermWindow, evChan chan event, cquit chan struct{}) {
	ticker := time.NewTicker(time.Millisecond * 100)
	inputs := []event{}
	var saved *VirtualRegion

updateloop:
	for {
		select {
		case <-ticker.C:

			if saved == nil {
				scene.Update(inputs)
				inputs = nil
			}

			w.Show()
		case e := <-evChan:
			inputs = append(inputs, e)
			switch ev := e.(type) {
			case *tcell.EventResize:
				saved = w.ResolutionCheck(ev, saved)
				w.Sync()
			case *specialEvent:
				switch ev.Key() {
				case quit:
					w.Fini()
					close(cquit)
					break updateloop
				case reset:
					scene.Reset()
					FillRect(' ', w.DrawableRect, w)
				}
			}
		}
	}
}

// GetDemoScene gets the demo scene, made for Mrs. Andres & Ms. Burritto
// as of April 10, 2022.
func GetDemoScene(w *TermWindow) Element {
	click := NewSoundEffect("assets/click.wav")
	ding := NewSoundEffect("assets/ding.wav")
	width, height := w.width, w.height

	return NewDiscretePlayer([]Element{

		NewSequentialPlayer([]Element{
			NewHoverText("                  __+--+__,\n                ,/        +-;\n               /            \\\n              |          .___|\n              |       ,_-+  |     ^\n              `\\____--+      \\    ||\n       ____     \\          <^   ^_LL,\n     _/^   \\-;___;-_     ,__;  /|__ |\n    / `- - _-L_     \\    -+___|     =)\n   /_     |    `.    |__/     '-____=)\n  /./    /|      \\    ,___+--/    /\n |  |   / `\\      +--/         ,-+\n/__/   |   `\\             .__-/\n|      |     `-___ __-+--+\nL______;              |\n       \\               \\\nArt by Kelsala",
				&Rect{(width/2 + 38) / 2, (height-17)/2 + 1, 100, 100}, ReplaceMap{}, w),
			NewHoverText("[{TOP SECRET}]", MarginRect((width-40-12)/2, height/2-2, 1, w), ReplaceMap{"top secret": {"         ________    |^|_.\n    __--+        \\___|   |\n  _|                     |___,\n /     Navajo Nation         |_ \n/            ._,               +--|^;\n\\       ,_---+ |     (Naabeehó      )\n|       |   <^=__      Bináhásdzo)   \\_,\n |.|^|  |       _|                     |\n     |  |______-                ,_____/`\n     |                  <\\      |\n     |___________,    .__|`|_   .\\\n                 U|-__|      `|_/\n                           .____,\n                         ,_|    |\n                         |____. |\n                              |_|\nArt by Kelsala", t2ne}}, w),
			NewTypewritter("How to (Navajo) Code Talk", MarginRect((width-40-25)/2, height/2-1, 1, w), w, click, ding),
			NewTypewritter("Press [SPACE] to start!", MarginRect((width-40-23)/2, height/2, 1, w), w, click, ding),
			NewWaitForNext(),
		}),

		NewSequentialPlayer([]Element{
			NewTypewritter("\tWelcome Private! You've been conscripted into the army. Due to your\nbackground, you have been assigned to a top secret group; the Navajo Code\nTalkers.",
				MarginRect(0, 0, 3, w), w, click, ding),
			NewTypewritter("How to navigate:\n\t• Press [SPACE] to advance and speed up text\n\t• Press [ESC] to reset the program.",
				MarginRect(0, 5, 3, w), w, click, ding),
			NewHoverText("\t• Use the mouse to hover over {colored text} for helpful tips.",
				MarginRect(0, 8, 1, w), ReplaceMap{"colored text": {" You found me! ", option}}, w),
			NewTypewritter("These commands are also found at the top of the screen.",
				MarginRect(0, 9, 1, w), w, click, ding),
			NewWaitForNext(),
		}),

		NewSequentialPlayer([]Element{
			NewTypewritter("LESSON 0:    WHO?",
				MarginRect(0, 0, 1, w), w, click, ding),
			NewTypewritter("\tYou may be asking who the Code Talkers are. Well, they are Native\nAmerican soldiers who transmit encoded messages through their native\nlanguage. Many languages are used, but the most common, and the one you will learn, is the Navajo Language, spoken in Northeastern Arizona and\nNorthwestern New Mexico.",
				MarginRect(0, 2, 5, w), w, click, ding),
			NewWaitForNext(),
		}),

		NewSequentialPlayer([]Element{
			NewTypewritter("LESSON 1:    INTRODUCTION",
				MarginRect(0, 0, 1, w), w, click, ding),
			NewTypewritter("Code talking consists of two types of code; Type 1 and Type 2.",
				MarginRect(0, 2, 1, w), w, click, ding),
			NewTypewritter("\tThe former is much like spelling out a word with words that start with\nthe same letter; \"Tab; T as in tea, A as in ant, B as in bear\". These words\n(tea, ant, bear) are then directly translated to their Navajo equivalents.",
				MarginRect(0, 4, 3, w), w, click, ding),
			NewHoverText("\tThe latter is straight translations from English to Navajo for common\nmilitary words. For words that don't exist in Navajo, like \"battleship\",\nanalogies like {whale} are used.",
				MarginRect(0, 7, 3, w), master, w),
			NewTypewritter("Let's start with Type 1.",
				MarginRect(0, 11, 1, w), w, click, ding),
			NewWaitForNext(),
		}),

		NewSequentialPlayer([]Element{
			NewTypewritter("LESSON 2:    TYPE 1 CODE",
				MarginRect(0, 0, 1, w), w, click, ding),
			NewHoverText("    Type 1 code is a simple alphabet substitution. You substitute each letterwith a word that begins with that letter. {T}{a}{b} in Type 1 code, therefore,\nwould be \"tea ant bear\", but then translated to Navajo; \"{dééh} {wóláchííʼ}\n{shash}\".",
				MarginRect(0, 2, 4, w), master, w),
			NewTypewritter("Let's try a longer example. Remember, hovering over colored gives you\nhints! Red text is Type 1 code.",
				MarginRect(0, 6, 2, w), w, click, ding),
			NewHoverText("What does {shash} {wóláchííʼ} {tsah} {wóláchííʼ} {tsah} {wóláchííʼ} spell? [TYPE THE LETTERS]",
				MarginRect(0, 9, 1, w), master, w),
			NewChecker(
				NewTextInput(MarginRect(0, 10, 1, w)),
				[]string{"banana"},
				NewTypewritter("Good Job! ",
					MarginRect(0, 11, 1, w), w, click, ding),
				NewTypewritter("Try again! Hover over the colored text to see the english translations.\nTranslate it to a english word!",
					MarginRect(0, 11, 1, w), w, click, ding),
			),
			NewWaitForNext(),
		}),

		NewSequentialPlayer([]Element{
			NewTypewritter("LESSON 3:    TYPE 2 CODE",
				MarginRect(0, 0, 1, w), w, click, ding),
			NewTypewritter("    Type 2 code is more like what you would expect from a code made from\nanother language. These are specific military terms used to speed up\ncommunication. Some terms don't have Navajo equivalents, and so descriptive",
				MarginRect(0, 2, 3, w), w, click, ding),
			NewHoverText("analogies are used. For example, \"submarine\" is an \"{iron fish}\".",
				MarginRect(0, 5, 1, w), master, w),
			NewHoverText("What might \"{tsídii} {mobba yéhé}\" mean? Blue text is Type 2 code.",
				MarginRect(0, 7, 1, w), master, w),
			NewChecker(
				NewOptions([]string{"cruiser", "bomber", "aircraft carrier"},
					MarginRect(0, 9, 0, w), w),
				[]string{"aircraft carrier"},
				NewTypewritter("Nice job! A thing that carries birds (planes) is an aircraft carrier!",
					MarginRect(0, 15, 1, w), w, click, ding),
				NewTypewritter("Try again! What might transport tsídii's?",
					MarginRect(0, 15, 1, w), w, click, ding),
			),
			NewWaitForNext(),
		}),

		NewSequentialPlayer([]Element{
			NewTypewritter("LESSON 4:    FINAL TEST",
				MarginRect(0, 0, 1, w), w, click, ding),
			NewTypewritter("\tAlright private! You've shown great progress. You should (theoretically)\nbe equipped to decode any text, and encrypt too, as long as you have a\ndictionary.",
				MarginRect(0, 2, 3, w), w, click, ding),
			NewTypewritter("This will be your final test: a combination of both Type 1 and 2 text. See ifyou can figure it out the instructions for Company B!",
				MarginRect(0, 6, 2, w), w, click, ding),
			NewHoverText("{Yókeed} {naakáí} {shash} {dééh} {tłʼohchin} {hohkááh} {dééh} {tłʼohchin} {tó nilį́į́h}.",
				MarginRect(0, 8, 1, w), master, w),
			NewChecker(
				NewTextInput(MarginRect(0, 9, 2, w)),
				[]string{"ask company b to come to creek", "ask company b to come to the creek"},
				NewTypewritter("You passed! Good job.",
					MarginRect(0, 10, 1, w), w, click, ding),
				NewTypewritter("Try again. Hover over the text to get translations! Red is Type 1, Blue is\nType 2.",
					MarginRect(0, 10, 2, w), w, click, ding),
			),
			NewWaitForNext(),
		}),

		NewSequentialPlayer([]Element{
			NewTypewritter("LESSON 5:    CONGRATS!",
				MarginRect(0, 0, 1, w), w, click, ding),
			NewTypewritter("\tGreat job! You've passed with flying colors. Now that you've been\ntrained, we'll see you on the battlefield!",
				MarginRect(0, 2, 2, w), w, click, ding),
			NewTypewritter("\tThe Navajo Code Talkers went on to serve in the US Marine Corps\nthroughout World War II, becoming a vital part of the war effort.\nThe code was much faster and reliable than other electronic codes at the\ntime- taking minutes rather than hours- and was one of the only codes never\nto be cracked by the Axis powers.",
				MarginRect(0, 5, 5, w), w, click, ding),
			NewTypewritter("\tUsed on all major island battles, from Guadalcanal to Iwo Jima to\nOkinawa, the talkers were classified for use in potential other wars until\n1968. Their contributions made the Navajo language more well known, and were partially responsible for inspiring new schools on the Navajo reservation\nthat teach Navajo language and culture to this day.",
				MarginRect(0, 10, 5, w), w, click, ding),
			NewHoverText("{\tThanks for playing! Please press [SPACE] or [ESC] to reset the simulationfor the next player once you're done reading. Thank you!}",
				MarginRect(0, 16, 2, w), ReplaceMap{"\tthanks for playing! please press [space] or [esc] to reset the simulationfor the next player once you're done reading. thank you!": {" Ahéheeʼ! ", t2ne}}, w),
			NewWaitForNext(),
		}),
	})
}
